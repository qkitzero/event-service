package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	authv1 "github.com/qkitzero/auth-service/gen/go/auth/v1"
	eventv1 "github.com/qkitzero/event-service/gen/go/event/v1"
	appevent "github.com/qkitzero/event-service/internal/application/event"
	apiauth "github.com/qkitzero/event-service/internal/infrastructure/api/auth"
	apiuser "github.com/qkitzero/event-service/internal/infrastructure/api/user"
	"github.com/qkitzero/event-service/internal/infrastructure/db"
	infraevent "github.com/qkitzero/event-service/internal/infrastructure/event"
	grpcevent "github.com/qkitzero/event-service/internal/interface/grpc/event"
	userv1 "github.com/qkitzero/user-service/gen/go/user/v1"
)

const shutdownTimeout = 15 * time.Second

type config struct {
	Env             string
	Port            string
	DBHost          string
	DBUser          string
	DBPassword      string
	DBName          string
	DBPort          string
	DBSSLMode       string
	AuthServiceHost string
	AuthServicePort string
	UserServiceHost string
	UserServicePort string
}

func loadConfig() (config, error) {
	env := os.Getenv("ENV")
	if env == "" {
		env = "development"
	}
	cfg := config{Env: env}
	required := []struct {
		key string
		dst *string
	}{
		{"PORT", &cfg.Port},
		{"DB_HOST", &cfg.DBHost},
		{"DB_USER", &cfg.DBUser},
		{"DB_PASSWORD", &cfg.DBPassword},
		{"DB_NAME", &cfg.DBName},
		{"DB_PORT", &cfg.DBPort},
		{"DB_SSL_MODE", &cfg.DBSSLMode},
		{"AUTH_SERVICE_HOST", &cfg.AuthServiceHost},
		{"AUTH_SERVICE_PORT", &cfg.AuthServicePort},
		{"USER_SERVICE_HOST", &cfg.UserServiceHost},
		{"USER_SERVICE_PORT", &cfg.UserServicePort},
	}
	var missing []string
	for _, r := range required {
		v := os.Getenv(r.key)
		if v == "" {
			missing = append(missing, r.key)
			continue
		}
		*r.dst = v
	}
	if len(missing) > 0 {
		return cfg, fmt.Errorf("missing required env vars: %s", strings.Join(missing, ", "))
	}
	return cfg, nil
}

func main() {
	if err := run(); err != nil {
		log.Fatalf("event-service: %v", err)
	}
}

func run() error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	gormDB, err := db.Init(cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort, cfg.DBSSLMode)
	if err != nil {
		return fmt.Errorf("db init: %w", err)
	}
	sqlDB, err := gormDB.DB()
	if err != nil {
		return fmt.Errorf("db handle: %w", err)
	}
	defer sqlDB.Close()

	listener, err := net.Listen("tcp", ":"+cfg.Port)
	if err != nil {
		return fmt.Errorf("listen: %w", err)
	}

	var dialOpt grpc.DialOption
	switch cfg.Env {
	case "production":
		dialOpt = grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, ""))
	default:
		dialOpt = grpc.WithTransportCredentials(insecure.NewCredentials())
	}

	authConn, err := grpc.NewClient(cfg.AuthServiceHost+":"+cfg.AuthServicePort, dialOpt)
	if err != nil {
		return fmt.Errorf("auth client: %w", err)
	}
	defer authConn.Close()

	userConn, err := grpc.NewClient(cfg.UserServiceHost+":"+cfg.UserServicePort, dialOpt)
	if err != nil {
		return fmt.Errorf("user client: %w", err)
	}
	defer userConn.Close()

	server := grpc.NewServer()

	authServiceClient := authv1.NewAuthServiceClient(authConn)
	userServiceClient := userv1.NewUserServiceClient(userConn)
	eventRepository := infraevent.NewEventRepository(gormDB)

	_ = apiauth.NewAuthService(authServiceClient)
	userService := apiuser.NewUserService(userServiceClient)
	eventUsecase := appevent.NewEventUsecase(userService, eventRepository)

	healthServer := health.NewServer()
	eventHandler := grpcevent.NewEventHandler(eventUsecase)

	grpc_health_v1.RegisterHealthServer(server, healthServer)
	eventv1.RegisterEventServiceServer(server, eventHandler)

	healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)
	healthServer.SetServingStatus("event", grpc_health_v1.HealthCheckResponse_SERVING)

	if cfg.Env == "development" {
		reflection.Register(server)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	serveErr := make(chan error, 1)
	go func() {
		log.Printf("gRPC server listening on %s", listener.Addr().String())
		serveErr <- server.Serve(listener)
	}()

	select {
	case err := <-serveErr:
		if err != nil {
			return fmt.Errorf("grpc serve: %w", err)
		}
		return nil
	case <-ctx.Done():
		log.Println("shutdown signal received, starting graceful stop")
		healthServer.Shutdown()

		stopped := make(chan struct{})
		go func() {
			server.GracefulStop()
			close(stopped)
		}()

		select {
		case <-stopped:
			log.Println("gRPC server stopped gracefully")
		case <-time.After(shutdownTimeout):
			log.Printf("graceful stop timed out after %s, forcing stop", shutdownTimeout)
			server.Stop()
			<-stopped
		}
		return nil
	}
}

package main

import (
	"log"
	"net"

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
	"github.com/qkitzero/event-service/util"
	userv1 "github.com/qkitzero/user-service/gen/go/user/v1"
)

func main() {
	db, err := db.Init(
		util.GetEnv("DB_HOST", ""),
		util.GetEnv("DB_USER", ""),
		util.GetEnv("DB_PASSWORD", ""),
		util.GetEnv("DB_NAME", ""),
		util.GetEnv("DB_PORT", ""),
		util.GetEnv("DB_SSL_MODE", ""),
	)
	if err != nil {
		log.Fatal(err)
	}

	listener, err := net.Listen("tcp", ":"+util.GetEnv("PORT", ""))
	if err != nil {
		log.Fatal(err)
	}

	authTarget := util.GetEnv("AUTH_SERVICE_HOST", "") + ":" + util.GetEnv("AUTH_SERVICE_PORT", "")
	userTarget := util.GetEnv("USER_SERVICE_HOST", "") + ":" + util.GetEnv("USER_SERVICE_PORT", "")

	var opts grpc.DialOption
	switch util.GetEnv("ENV", "development") {
	case "production":
		opts = grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, ""))
	default:
		opts = grpc.WithTransportCredentials(insecure.NewCredentials())
	}

	authConn, err := grpc.NewClient(authTarget, opts)
	if err != nil {
		log.Fatal(err)
	}
	defer authConn.Close()

	userConn, err := grpc.NewClient(userTarget, opts)
	if err != nil {
		log.Fatal(err)
	}
	defer userConn.Close()

	server := grpc.NewServer()

	authServiceClient := authv1.NewAuthServiceClient(authConn)
	userServiceClient := userv1.NewUserServiceClient(userConn)
	eventRepository := infraevent.NewEventRepository(db)

	_ = apiauth.NewAuthService(authServiceClient)
	userService := apiuser.NewUserService(userServiceClient)
	eventUsecase := appevent.NewEventUsecase(userService, eventRepository)

	healthServer := health.NewServer()
	eventHandler := grpcevent.NewEventHandler(eventUsecase)

	grpc_health_v1.RegisterHealthServer(server, healthServer)
	eventv1.RegisterEventServiceServer(server, eventHandler)

	healthServer.SetServingStatus("event", grpc_health_v1.HealthCheckResponse_SERVING)

	if util.GetEnv("ENV", "development") == "development" {
		reflection.Register(server)
	}

	if err = server.Serve(listener); err != nil {
		log.Fatal(err)
	}
}

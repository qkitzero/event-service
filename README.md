# Event Service

[![Go Test](https://github.com/qkitzero/event-service/actions/workflows/test.yml/badge.svg)](https://github.com/qkitzero/event-service/actions/workflows/test.yml)
[![codecov](https://codecov.io/gh/qkitzero/event-service/graph/badge.svg?token=hoFlz6L4aQ)](https://codecov.io/gh/qkitzero/event-service)
[![Buf CI](https://github.com/qkitzero/event-service/actions/workflows/buf-ci.yaml/badge.svg)](https://github.com/qkitzero/event-service/actions/workflows/buf-ci.yaml)

- Microservices Architecture
- gRPC
- gRPC Gateway
- Buf ([buf.build/qkitzero/event-service](https://buf.build/qkitzero/event-service))
- Clean Architecture
- Docker
- Test
- Codecov
- Cloud Build
- Cloud Run

```mermaid
flowchart TD
    subgraph gcp[GCP]
        secret_manager[Secret Manager]

        subgraph cloud_build[Cloud Build]
            build_event_service(Build event-service)
            push_event_service(Push event-service)
            deploy_event_service(Deploy event-service)

            build_event_service_gateway(Build event-service-gateway)
            push_event_service_gateway(Push event-service-gateway)
            deploy_event_service_gateway(Deploy event-service-gateway)
        end


        subgraph artifact_registry[Artifact Registry]
            event_service_image[(event-service image)]
            event_service_gateway_image[(event-service-gateway image)]
        end

        subgraph cloud_run[Cloud Run]
            event_service(Event Service)
            event_service_gateway(Event Service Gateway)
        end
    end

    subgraph external[External]
        event_db[(Event DB)]
        user_service(User Service)
    end

    build_event_service --> push_event_service --> event_service_image
    build_event_service_gateway --> push_event_service_gateway --> event_service_gateway_image

    event_service_image --> deploy_event_service --> event_service
    event_service_gateway_image --> deploy_event_service_gateway --> event_service_gateway

    secret_manager --> deploy_event_service
    secret_manager --> deploy_event_service_gateway

    event_service_gateway --> event_service
    event_service --> event_db
    event_service --> user_service
```

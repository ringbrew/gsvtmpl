# gsv Framework Development Guide

## 1. Overview

This document serves as the official guide for developing Go backend services using the `gsv` framework. `gsv` is designed to follow **Clean Architecture** principles and organizes code in a **Domain-Driven** manner.

This guide details the project architecture, directory structure, and core development workflows. It is intended to provide a clear context for both **developers and AI programming assistants** to facilitate daily development tasks.

> **Gsv also provides a command-line toolkit. If you want to use any of Gsv's commands, please make sure you are in the root directory of your project.**

### 1.1 Core Project Architecture

The `gsv` project structure adheres to the separation of concerns, primarily divided into three layers:

*   **Domain Layer**: The core business logic layer. It contains business Entities, Use Cases, and Repository Interfaces. This layer does not depend on any external implementations.
*   **Delivery Layer**: The entry and exit points of the application. It is responsible for dispatching external requests (e.g., HTTP, gRPC) or events (e.g., Cron jobs, Message Queues) to the corresponding Domain Use Cases for processing.
*   **Infrastructure Layer**: Provides implementations for external dependencies required by the project, such as database connections (DAO), Redis clients, and gRPC clients. These instances are provided to the Domain layer via dependency injection.

### 1.2 Directory Structure

A typical `gsv` project directory structure is as follows:

```
📁 my-gsv-project/
├── 📁 build/             # Build-related files (Dockerfile, docker-compose.yml)
├── 📁 cmd/               # Main application entrypoints (main.go), can have multiple
├── 📁 export/            # Definitions to be exposed externally (e.g., generated gRPC code)
├── 📁 internal/          # Core business logic of the project
│   ├── 📁 conf/          # Business configuration definitions (config.go)
│   ├── 📁 delivery/      # Delivery layer implementations (HTTP, gRPC, Cron, MQ)
│   └── 📁 domain/        # Domain layer implementations (business core)
├── 📁 openapi/           # Auto-generated Swagger/OpenAPI JSON files
├── 📁 proto/             # .proto definition files for gRPC
├── 📄 .gitignore
├── 📄 go.mod
├── 📄 go.sum
├── 📄 Makefile
└── 📄 README.md
```

## 2. Domain Layer (`internal/domain`)

The Domain layer is the heart of the business logic, where all business use cases are implemented.

### 2.1 Domain Context (`domain/context.go`)

This file defines the `UseCaseContext` struct, which acts as the application's runtime **Dependency Injection (DI) container**. All infrastructure layer instances (e.g., database, Redis, external service clients) are initialized and managed here.

**Code Example:**
```golang
// internal/domain/context.go
package domain

import (
	// ... imports
	"github.com/ringbrew/gsv/cli"
	"github.com/ringbrew/gsv/discovery"
)

// UseCaseContext holds all shared resources needed by use cases.
type UseCaseContext struct {
	Config    conf.Config
	Signal    context.Context
	cancel    context.CancelFunc
	WaitGroup sync.WaitGroup
	Redis     *redis.Client
	MysqlDao  *dao.DAO
	// gRPC client initialized via gsv/cli
	JobCli    job.ServiceClient
}

// NewUseCaseContext initializes and returns a singleton UseCaseContext.
func NewUseCaseContext(c conf.Config) *UseCaseContext {
	// ... (initialization logic)

	// Inject DAO
	dsc.MysqlDao = dao.NewDataAccess(...)
	
	// Inject Redis Client
	dsc.Redis = redis.NewClient(...)

	// Initialize gRPC client using gsv/cli
	opts := cli.Classic()
	cc, err := cli.NewClient(discovery.Scheme("go-job"), opts)
	if err != nil {
		log.Fatal(err.Error())
	}
	dsc.JobCli = job.NewServiceClient(cc.Conn())
	
	// ... (other initializations)
	return dsc
}
```

### 2.2 Domain Implementation (`domain/<domain_name>/`)

Each business domain resides in a separate directory. The recommended structure includes entities, use cases, and a repository.

**Quick Start**: Use the `gsv` CLI tool to rapidly create a domain skeleton.
```shell
gsv gen domain <domain_name>
```
For example, `gsv gen domain example` will create the following structure:

#### `example.go` (Entity Definition)
```golang
package example

type Example struct {
	Id   int
	// ... other fields
}
```

#### `usecase.go` (Use Case Implementation)
A use case encapsulates a specific business process, calling the repository to persist data.
```golang
package example

import "xxx/internal/domain"

type UseCase struct {
	ctx  *domain.UseCaseContext
	repo *repo
}

func NewUseCase(ctx *domain.UseCaseContext) *UseCase {
	return &UseCase{
		ctx:  ctx,
		repo: newRepo(), // repo is typically stateless
	}
}

func (uc *UseCase) Create(ctx context.Context, example *Example) error {
	return uc.repo.Create(ctx, example)
}
// ... other use case methods: Update, Delete, Get, Query
```

#### `repo.go` (Repository Implementation)
The repository is responsible for interacting with data storage. In simple scenarios, it can directly implement data access logic.
```golang
package example

// The repo struct is typically stateless.
type repo struct {}

func newRepo() *repo {
	return &repo{}
}

// Create persists an Example entity.
func (r *repo) Create(ctx context.Context, example *Example) error {
	// Implement database insertion logic here.
	// For example: return domain.GetUseCaseContext().MysqlDao.Create(example)
	return nil
}
// ... other data access methods
```

## 3. Delivery Layer (`internal/delivery`)

The Delivery layer acts as the bridge
between the outside world and the domain use cases.

### 3.1 Bootstrap File (`delivery/bootstrap.go`)

This file is the assembly point of the Delivery layer. It is responsible for:

1. Choosing the server type (`server.GRPC` or `server.HTTP`).
2. Configuring server-level middleware/interceptors.
3. Registering all delivery services through `ServiceList`.

```golang
package delivery

import (
	"github.com/ringbrew/gsv/server"
	"github.com/ringbrew/gsv/service"
	"{{moduleName}}/internal/domain"
)

// NewServer creates a gsv server instance.
func NewServer(ctx *domain.UseCaseContext) server.Server {
	opt := server.Classic()
	// configure ports, host, grpc-gateway, middleware, interceptors, etc.
	return server.NewServer(server.GRPC, &opt)
}

// ServiceList registers all services that implement the service.Service interface.
func ServiceList(ctx *domain.UseCaseContext) []service.Service {
	return []service.Service{
		NewExampleGrpcService(ctx),
		NewExampleHttpService(ctx),
		// ... add more services here
	}
}
```

**Recommended rule**:

Keep all business logic in the Domain layer. The Delivery layer should focus on request parsing, protocol conversion, authentication/authorization hooks, and invoking domain use cases.

### 3.2 gRPC Service

In `gsv`, a gRPC delivery usually follows this flow:

1. Define the `.proto` file in `proto/`.
2. Ensure `go_package` points to the `export/` package.
3. Run `gsv gen grpc`.
4. Implement the generated service by translating request/response models and calling the domain use case.
5. Register the service in `delivery/bootstrap.go`.


#### 3.2.1 Define proto first

```proto
syntax = "proto3";
package example;

option go_package = "{{moduleName}}/export/example";

service Service {
    rpc Create(Example) returns (OpResp){};
}
```

After the proto is ready, generate code:

```shell
gsv gen grpc
```

This command typically generates:

1. gRPC-related code in `export/`.
2. A delivery service skeleton in `internal/delivery/<name>/`.

#### 3.2.2 Implement the generated service

The generated service should stay thin and delegate the real business work to `internal/domain/<domain_name>`.

```golang
package example

import (
	"context"
	"github.com/ringbrew/gsv/service"
	pb "{{moduleName}}/export/example"
	"{{moduleName}}/internal/domain"
	domainexample "{{moduleName}}/internal/domain/example"
	"google.golang.org/grpc"
)

type Service struct {
	pb.UnimplementedServiceServer
	ctx  *domain.UseCaseContext
	desc service.Description
}

func NewService(ctx *domain.UseCaseContext) service.Service {
	return &Service{
		ctx: ctx,
		desc: service.Description{
			Valid:           true,
			GrpcServiceDesc: []grpc.ServiceDesc{pb.Service_ServiceDesc},
		},
	}
}

func (s *Service) Description() service.Description {
	return s.desc
}

func (s *Service) Create(ctx context.Context, req *pb.Example) (*pb.OpResp, error) {
	uc := domainexample.NewUseCase(s.ctx)

	model := &domainexample.Example{
		// map proto fields to domain entity
	}

	if err := uc.Create(ctx, model); err != nil {
		return nil, err
	}

	return &pb.OpResp{Code: 200, Msg: "Success"}, nil
}
```

#### 3.2.3 `grpc-gateway` vs host-only gRPC

There are two common gRPC exposure styles:

1. **Host-only gRPC**: only exposes the native gRPC service.
2. **gRPC + grpc-gateway**: exposes both gRPC and HTTP endpoints for the same RPC.

To expose an RPC through `grpc-gateway`, add `google.api.http` options in the proto:

```proto
rpc ExportHostAccount(ExportHostRequest) returns (ExportHostResponse) {
    option (google.api.http) = {
        post: "url"
        body: "*"
    };
};
```

If an RPC does not define a `google.api.http` option, it remains a normal gRPC-only endpoint.

At the server level, `grpc-gateway` must also be enabled in `delivery/bootstrap.go`:

```golang
func NewServer(ctx *domain.UseCaseContext) server.Server {
	opt := server.Classic()
	opt.EnableGrpcGateway = true
	return server.NewServer(server.GRPC, &opt)
}
```

Besides the proto annotation and server switch, the delivery service itself must also register the gateway handler in `service.Description`. Otherwise, the HTTP mapping will not be mounted.

```golang
func NewService(ctx *domain.UseCaseContext) service.Service {
	return &Service{
		ctx: ctx,
		desc: service.Description{
			Valid:           true,
			GrpcServiceDesc: []grpc.ServiceDesc{pb.Service_ServiceDesc},
			GrpcGateway: []service.GatewayRegister{
				pb.RegisterServiceHandler,
			},
		},
	}
}
```

So enabling `grpc-gateway` requires all three parts:

1. Add `google.api.http` to the RPC in the proto.
2. Enable `opt.EnableGrpcGateway = true` in `delivery/bootstrap.go`.
3. Register `GrpcGateway` in the service `Description()`.

If you only want host-only gRPC, keep `GrpcServiceDesc` and omit `GrpcGateway`.

#### 3.2.4 Register gRPC interceptors in bootstrap

gRPC middleware is configured on the server options, typically through `UnaryInterceptors`:

```golang
func NewServer(ctx *domain.UseCaseContext) server.Server {
	opt := server.Classic()
	opt.EnableGrpcGateway = true
	opt.UnaryInterceptors = append(opt.UnaryInterceptors, errbiz.ErrWrapInterceptor())
	return server.NewServer(server.GRPC, &opt)
}
```

Use this hook for concerns such as:

1. Error wrapping.
2. Authentication.
3. Logging and tracing.
4. Metrics and request auditing.

### 3.3 HTTP Service

Compared with gRPC delivery, HTTP delivery in `gsv` is usually organized around two parts:

1. `handler.go`: define request handlers and routes.
2. `service.http.impl.go`: assemble the handler routes into a `service.Service`.

#### 3.3.1 Recommended workflow

1. Create an HTTP delivery module, or generate one with `gsv gen http <service_name>`.
2. Implement the request handling logic in `handler.go`.
3. Expose routes through `HttpRoute()`.
4. Assemble them in `service.http.impl.go`.
5. Register the service in `delivery/bootstrap.go`.
6. Register HTTP middleware in `delivery/bootstrap.go` when needed.

#### 3.3.2 Define `Handler` and routes

```golang
package auth

import (
	"net/http"

	"github.com/ringbrew/gsv/service"
	"{{moduleName}}/internal/domain"
)

type Handler struct {
	ctx *domain.UseCaseContext
}

func NewHandler(ctx *domain.UseCaseContext) *Handler {
	return &Handler{ctx: ctx}
}

func (h *Handler) Sign(w http.ResponseWriter, r *http.Request) {
	// 1. parse request
	// 2. call domain or remote client
	// 3. write response
}

func (h *Handler) HttpRoute() []service.HttpRoute {
	return []service.HttpRoute{
		service.NewHttpRoute(http.MethodPost, "/sign", h.Sign, service.HttpMeta{
			Remark: "Sign in",
		}),
	}
}
```

#### 3.3.3 Assemble routes in `service.http.impl.go`

```golang
package auth

import (
	"{{moduleName}}/internal/domain"
	"github.com/ringbrew/gsv/service"
)

type Service struct {
	ctx  *domain.UseCaseContext
	desc service.Description
}

func NewService(ctx *domain.UseCaseContext) service.Service {
	s := &Service{
		ctx: ctx,
		desc: service.Description{
			Valid: true,
		},
	}

	handler := NewHandler(ctx)
	s.desc.HttpRoute = append(s.desc.HttpRoute, handler.HttpRoute()...)
	return s
}
```

This split keeps route definitions and HTTP handling code close together while still exposing a unified `service.Service` to the framework.

#### 3.3.4 Register HTTP middleware in bootstrap

HTTP middleware is also configured in `delivery/bootstrap.go`, for example:

```golang
func NewServer(ctx *domain.UseCaseContext) server.Server {
	opt := server.Classic()
	opt.Port = ctx.Config.Port.Api
	opt.HttpMiddleware = append(opt.HttpMiddleware, cors.AllowAll(), auth.NewMiddleware(ctx))
	return server.NewServer(server.HTTP, &opt)
}
```

Typical HTTP middleware responsibilities include:

1. CORS.
2. Authentication and token parsing.
3. Permission checks.
4. Request context injection.

When the service starts, the framework combines:

1. The server created by `NewServer`.
2. The service list returned by `ServiceList`.
3. The middleware/interceptor chain configured in bootstrap.

That makes `delivery/bootstrap.go` the central place to understand how the entire Delivery layer is exposed.

#### 3.3.5 Recommended HTTP middleware structure

For HTTP projects, middleware is usually implemented as an independent file such as `middlware.go` and then registered in bootstrap.

```golang
package auth

import (
	"net/http"

    "github.com/ringbrew/gsv/server"
	"{{moduleName}}/internal/domain"
)

type Middleware struct {
	ctx *domain.UseCaseContext
}

func NewMiddleware(ctx *domain.UseCaseContext) server.Handler {
	return &Middleware{ctx: ctx}
}

func (m *Middleware) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	// 1. parse token or headers
	// 2. validate identity / permission
	// 3. inject context
	next(rw, r)
}
```

This keeps HTTP concerns separated cleanly:

1. `handler.go` handles routes and request/response conversion.
2. `service.http.impl.go` assembles routes as a service.
3. `middlware.go` handles cross-cutting concerns such as auth and context injection.

### 3.4 Recommended Delivery Development Process

When adding a new business capability, the recommended order is:

1. Implement or extend the business use case in `internal/domain`.
2. Decide the exposure style in Delivery:
   - Use gRPC when the capability is primarily internal service-to-service communication.
   - Use gRPC + `grpc-gateway` when the same RPC also needs an HTTP entry.
   - Use a pure HTTP delivery when the interface is designed directly around HTTP handlers.
3. Implement the delivery adapter:
   - For gRPC: define proto, run `gsv gen grpc`, then implement the generated service.
   - For HTTP: create `handler.go` and `service.http.impl.go`, then register routes.
4. Register the delivery service in `ServiceList`.
5. Register interceptors or middleware in `NewServer` when needed.
6. Keep protocol conversion in Delivery and business rules in Domain.

### 3.5 Practical Checklist

Before considering a new module complete, verify the following:

1. The business logic is implemented in `internal/domain/<domain_name>`.
2. The delivery code only performs protocol adaptation and orchestration.
3. The service has been registered in `internal/delivery/bootstrap.go`.
4. Required middleware or interceptors have been registered in `NewServer`.
5. If `grpc-gateway` is needed, both the proto `google.api.http` option and `opt.EnableGrpcGateway = true` and `opt.GrpcGateway` are configured.
6. If it is an HTTP service, routes are exposed through `HttpRoute()` and assembled in `service.http.impl.go`.
package app

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	auth_interceptors "github.com/neracastle/auth/pkg/user_v1/auth/grpc-interceptors"
	"github.com/neracastle/go-libs/pkg/closer"
	"github.com/neracastle/go-libs/pkg/sys/logger"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	grpc_server "github.com/neracastle/chat-server/internal/grpc-server"
	"github.com/neracastle/chat-server/internal/grpc-server/interceptors"
	"github.com/neracastle/chat-server/internal/tracer"
	"github.com/neracastle/chat-server/pkg/chat_v1"
)

type App struct {
	grpc          *grpc.Server
	srvProvider   *serviceProvider
	traceExporter *otlptrace.Exporter
}

func NewApp(ctx context.Context) *App {
	app := &App{srvProvider: newServiceProvider()}
	app.init(ctx)
	app.initTracing(ctx, "chat-service")
	return app
}

func (a *App) init(ctx context.Context) {
	lg := logger.SetupLogger(a.srvProvider.Config().Env)
	a.grpc = grpc.NewServer(
		grpc.Creds(insecure.NewCredentials()),
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
		grpc.ChainUnaryInterceptor(
			interceptors.NewLoggerInterceptor(lg),
			auth_interceptors.NewAccessInterceptor([]string{
				chat_v1.ChatV1_Create_FullMethodName,
				chat_v1.ChatV1_Delete_FullMethodName,
			}, a.srvProvider.Config().SecretKey),
		),
	)

	reflection.Register(a.grpc)
	chat_v1.RegisterChatV1Server(a.grpc, grpc_server.NewServer(lg, a.srvProvider.ChatService(ctx)))
}

func (a *App) initTracing(ctx context.Context, serviceName string) {
	//экспортер в jaeger
	exporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(a.srvProvider.Config().Trace.JaegerGRPCAddress))
	if err != nil {
		log.Fatalf("failed to create trace exporter: %v", err)
	}
	a.traceExporter = exporter

	//собиратель трейсов
	r, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(semconv.SchemaURL, semconv.ServiceNameKey.String(serviceName)),
	)
	if err != nil {
		log.Fatalf("failed to create trace provider: %v", err)
	}

	traceProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter, sdktrace.WithExportTimeout(time.Second*time.Duration(a.srvProvider.Config().Trace.BatchTimeout))),
		sdktrace.WithResource(r))

	//пробрасываем провайдер для исп. в других местах приложения
	tracer.Init(traceProvider.Tracer(serviceName))
	//регистрируем глобально
	prop := propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})
	otel.SetTracerProvider(traceProvider)
	otel.SetTextMapPropagator(prop)
}

func (a *App) Start() error {
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()

	conn, err := net.Listen("tcp", fmt.Sprintf("%s:%d", a.srvProvider.Config().GRPC.Host, a.srvProvider.Config().GRPC.Port))
	if err != nil {
		return err
	}

	log.Printf("ChatAPI service started on %s:%d\n", a.srvProvider.Config().GRPC.Host, a.srvProvider.Config().GRPC.Port)

	closer.Add(func() error {
		a.grpc.GracefulStop()
		return nil
	})

	if err = a.grpc.Serve(conn); err != nil {
		return err
	}

	return nil
}

package main

import (
	"fmt"
	"github.com/KirillGreenev/crypto-rate-service/pkg/tracing"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/KirillGreenev/crypto-rate-service/cmd/config"
	crl "github.com/KirillGreenev/crypto-rate-service/internal/controller/grpc"
	pb "github.com/KirillGreenev/crypto-rate-service/internal/controller/grpc/proto"
	"github.com/KirillGreenev/crypto-rate-service/internal/repository/api"
	"github.com/KirillGreenev/crypto-rate-service/internal/repository/postgres"
	"github.com/KirillGreenev/crypto-rate-service/internal/service"
	"github.com/KirillGreenev/crypto-rate-service/pkg/debug"
	"github.com/KirillGreenev/crypto-rate-service/pkg/logger"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	cf := config.LoadConfig()
	logger.BuildLogger(cf.LoggerLevel)

	shutdown := tracing.InitTracer("otel-collector:4317")
	defer shutdown()

	db := config.NewPostgres().InitDB(cf).Migrate()
	defer db.Close()

	log := logger.Logger()

	serverGRPC := grpc.NewServer(grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()))
	pb.RegisterRatesServiceServer(serverGRPC,
		crl.NewRatesServiceGRPC(service.NewRatesServiceImpl(api.NewGarantexApiImpl(),
			postgres.NewPostgesRepositoryImpl(db.GetDB()))))

	l, err := net.Listen("tcp", fmt.Sprintf(":%v", cf.GRPCServerPort))
	defer l.Close()

	if err != nil {
		log.Fatal("listen error", zap.Error(err))
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Info(fmt.Sprintf("debug service started on the port %v", cf.DebugServerPort))
		debug.Run(fmt.Sprintf(":%v", cf.DebugServerPort))
	}()

	go func() {
		log.Info(fmt.Sprintf("The crypto-rate-service server gRPC is running on port %v", cf.GRPCServerPort))
		if err = serverGRPC.Serve(l); err != nil {
			log.Fatal("Accept error", zap.Error(err))
		}
	}()

	<-stop
	serverGRPC.GracefulStop()
	log.Info("Server stopped")
}

package main

import (
	"fmt"
	proto "github.com/DoktorGhost/external-api/src/go/pkg/grpc/clients/api/grpc/protobuf/books_v1"
	"github.com/DoktorGhost/golibrary-books/config"
	"github.com/DoktorGhost/golibrary-books/internal/app"
	"github.com/DoktorGhost/golibrary-books/internal/delivery/controllers/handlers"
	"github.com/DoktorGhost/golibrary-books/internal/delivery/grpc/grpcUC"
	deliveryServ "github.com/DoktorGhost/golibrary-books/internal/delivery/grpc/server"
	"github.com/DoktorGhost/golibrary-books/internal/delivery/http/server"
	"github.com/DoktorGhost/platform/logger"
	"github.com/DoktorGhost/platform/storage/psg"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	//инициализация логгера
	log, err := logger.GetLogger()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer log.Sync()

	//загрзка переменных окружения
	viper.AutomaticEnv()

	// Конвертируем в конфигурацию, которую ожидает InitStorage
	psgConfig := config.ConvertToPsgDBConfig(config.LoadConfig().DBConfig)

	// Инициализируем подключение к БД
	var pgsqlConnector *pgxpool.Pool

	for i := 0; i < 5; i++ {
		pgsqlConnector, err = psg.InitStorage(psgConfig)
		if err != nil {
			log.Error(err.Error())
			time.Sleep(5 * time.Second)
			continue
		}
		log.Info("соединение с БД установлено")
		break
	}

	defer pgsqlConnector.Close()

	r := handlers.SetupRoutes(app.Init(pgsqlConnector).UseCaseProvider)

	//старт grpc сервера
	lis, err := net.Listen("tcp", ":"+config.LoadConfig().ProviderConfig.Provider_port)
	if err != nil {
		log.Fatal("failed to listen: %v", "err", err)
	}

	grpcServer := grpc.NewServer()
	bookGRPCServer := grpcUC.NewBookGRPCServer(app.Init(pgsqlConnector).UseCaseProvider.BookUseCase)

	proto.RegisterBooksServiceServer(grpcServer, bookGRPCServer)
	reflection.Register(grpcServer)

	grpcSrv := deliveryServ.NewGRPCServer(lis, grpcServer)

	grpcSrv.Serve()
	log.Info("Grpc-server started")

	//старт http-сервера
	httpServer := server.NewHttpServer(r, ":"+config.LoadConfig().ProviderConfig.Http_port)
	httpServer.Serve()
	log.Info("http-server started")

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGTERM, syscall.SIGINT)

	select {
	case killSignal := <-interrupt:
		log.Info("Выключение сервера", "signal", killSignal)
	case err = <-httpServer.Notify():
		log.Error("Ошибка HTTP сервера", "error", err)
	case err = <-grpcSrv.Notify():
		log.Error("Ошибка GRPC сервера", "error", err)
	}

	httpServer.Shutdown()
	grpcSrv.Shutdown()

}

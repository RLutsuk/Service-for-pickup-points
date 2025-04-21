package main

import (
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"net"
	"os"

	"github.com/fatih/color"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"

	pickupPointDev "github.com/RLutsuk/Service-for-pickup-points/app/internal/pickup_point/delivery"
	pickupPointRep "github.com/RLutsuk/Service-for-pickup-points/app/internal/pickup_point/repository"
	pickupPointUC "github.com/RLutsuk/Service-for-pickup-points/app/internal/pickup_point/usecase"
	grpcPP "github.com/RLutsuk/Service-for-pickup-points/app/proto/grpc"
	pickuppoint "github.com/RLutsuk/Service-for-pickup-points/app/proto/pickuppoint"

	productDev "github.com/RLutsuk/Service-for-pickup-points/app/internal/product/delivery"
	productRep "github.com/RLutsuk/Service-for-pickup-points/app/internal/product/repository"
	productUC "github.com/RLutsuk/Service-for-pickup-points/app/internal/product/usecase"

	receptionDev "github.com/RLutsuk/Service-for-pickup-points/app/internal/reception/delivery"
	receptionRep "github.com/RLutsuk/Service-for-pickup-points/app/internal/reception/repository"
	receptionUC "github.com/RLutsuk/Service-for-pickup-points/app/internal/reception/usecase"

	userDev "github.com/RLutsuk/Service-for-pickup-points/app/internal/user/delivery"
	userRep "github.com/RLutsuk/Service-for-pickup-points/app/internal/user/repository"
	userUC "github.com/RLutsuk/Service-for-pickup-points/app/internal/user/usecase"
	echoSwagger "github.com/swaggo/echo-swagger"

	middleware "github.com/RLutsuk/Service-for-pickup-points/app/internal/middleware"
	monitoring "github.com/RLutsuk/Service-for-pickup-points/app/monitoring"

	_ "github.com/RLutsuk/Service-for-pickup-points/docs"
	"github.com/labstack/echo/v4"
)

// @title Service for Avito pick-up point
// @version 1.1
// @description API server for handling product reception at pick-up points
// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {

	serverAddress := os.Getenv("SERVER_ADDRESS")
	postgresConn := os.Getenv("POSTGRES_CONN")
	postgresUser := os.Getenv("POSTGRES_USERNAME")
	postgresPassword := os.Getenv("POSTGRES_PASSWORD")
	postgresHost := os.Getenv("POSTGRES_HOST")
	postgresPort := os.Getenv("POSTGRES_PORT")
	postgresDB := os.Getenv("POSTGRES_DATABASE")

	var config string

	if postgresConn != "" {
		config = postgresConn
	} else {
		config = fmt.Sprintf("host=%s user=%s password=%s database=%s port=%s",
			postgresHost, postgresUser, postgresPassword, postgresDB, postgresPort)
	}

	// for local testing
	// if postgresHost == "" {
	// 	config = "host=localhost user=db_pg password=db_postgres database=db_pps port=5432 sslmode=disable"
	// 	serverAddress = ":8080"
	// }

	db, err := sql.Open("postgres", config)
	defer db.Close()
	if err != nil {
		log.Fatal(err)
	}

	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)

	pickupPointDB := pickupPointRep.New(db)
	receptionDB := receptionRep.New(db)
	productDB := productRep.New(db)
	userDB := userRep.New(db)

	loggerUC := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		// AddSource: true,
	}))
	slog.SetDefault(loggerUC)

	pickupPointUC := pickupPointUC.New(pickupPointDB, loggerUC)
	receptionUC := receptionUC.New(receptionDB, pickupPointDB, loggerUC)
	productUC := productUC.New(productDB, receptionDB, pickupPointDB, loggerUC)
	userUC := userUC.New(userDB, loggerUC)

	e := echo.New()

	loggerDel := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		// AddSource: true,
	}))
	slog.SetDefault(loggerDel)

	pickupPointDev.NewDelivery(e, pickupPointUC, loggerDel)
	receptionDev.NewDelivery(e, receptionUC, loggerDel)
	productDev.NewDelivery(e, productUC, loggerDel)
	userDev.NewDelivery(e, userUC, loggerDel)

	promClient := echo.New()
	promClient.GET("/metrics", echo.WrapHandler(promhttp.Handler()))
	monitoring.Init()

	go func() {
		promClient.Logger.Fatal(promClient.Start(":9000"))
	}()

	go startGRPCServer(pickupPointUC)

	e.Use(middleware.PrometheusMiddleware)

	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.Logger.Fatal(e.Start(serverAddress))
}

func startGRPCServer(pickupPointUC pickupPointUC.UseCaseI) {
	lis, err := net.Listen("tcp", ":3000")
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()
	srv := grpcPP.NewGRPCServer(pickupPointUC)
	pickuppoint.RegisterPPServiceServer(s, srv)

	green := color.New(color.FgGreen).SprintFunc()
	fmt.Printf("%s %s\n", "⇨ grpc server started on", green("[::]:3000"))
	if err := s.Serve(lis); err != nil {
		log.Fatal(err)
	}
}

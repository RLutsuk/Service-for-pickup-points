package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"

	pickupPointDev "github.com/RLutsuk/Service-for-pickup-points/app/internal/pickup_point/delivery"
	pickupPointRep "github.com/RLutsuk/Service-for-pickup-points/app/internal/pickup_point/repository"
	pickupPointUC "github.com/RLutsuk/Service-for-pickup-points/app/internal/pickup_point/usecase"

	productDev "github.com/RLutsuk/Service-for-pickup-points/app/internal/product/delivery"
	productRep "github.com/RLutsuk/Service-for-pickup-points/app/internal/product/repository"
	productUC "github.com/RLutsuk/Service-for-pickup-points/app/internal/product/usecase"

	receptionDev "github.com/RLutsuk/Service-for-pickup-points/app/internal/reception/delivery"
	receptionRep "github.com/RLutsuk/Service-for-pickup-points/app/internal/reception/repository"
	receptionUC "github.com/RLutsuk/Service-for-pickup-points/app/internal/reception/usecase"

	userDev "github.com/RLutsuk/Service-for-pickup-points/app/internal/user/delivery"
	userRep "github.com/RLutsuk/Service-for-pickup-points/app/internal/user/repository"
	userUC "github.com/RLutsuk/Service-for-pickup-points/app/internal/user/usecase"

	"github.com/labstack/echo"
)

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

	//потом убрать !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
	if postgresHost == "" {
		config = "host=localhost user=db_pg password=db_postgres database=db_pps port=5432 sslmode=disable"
		serverAddress = ":8080"
	}

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

	pickupPointUC := pickupPointUC.New(pickupPointDB)
	receptionUC := receptionUC.New(receptionDB, pickupPointDB)
	productUC := productUC.New(productDB, receptionDB, pickupPointDB)
	userUC := userUC.New(userDB)

	e := echo.New()

	pickupPointDev.NewDelivery(e, pickupPointUC)
	receptionDev.NewDelivery(e, receptionUC)
	productDev.NewDelivery(e, productUC)
	userDev.NewDelivery(e, userUC)

	e.Logger.Fatal(e.Start(serverAddress))
}

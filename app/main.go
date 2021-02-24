package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/labstack/echo"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"

	fileRepo "github.com/rknizzle/rkmesh/file/repository/s3"
	_modelHTTPController "github.com/rknizzle/rkmesh/model/controller/http"
	_modelRepo "github.com/rknizzle/rkmesh/model/repository/postgres"
	_modelService "github.com/rknizzle/rkmesh/model/service"
)

func init() {
	viper.SetConfigFile(`config.json`)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

func main() {
	dbHost := viper.GetString(`database.host`)
	dbPort := viper.GetString(`database.port`)
	dbUser := viper.GetString(`database.user`)
	dbPass := viper.GetString(`database.pass`)
	dbName := viper.GetString(`database.name`)

	connection := fmt.Sprintf(
		`host=%s port=%s user=%s
		password=%s dbname=%s sslmode=disable`,
		dbHost, dbPort, dbUser, dbPass, dbName)

	dbConn, err := sql.Open("postgres", connection)

	// run database migrations
	driver, err := postgres.WithInstance(dbConn, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
	}
	mig, err := migrate.NewWithDatabaseInstance(
		"file://./migrations",
		"postgres", driver)
	if err != nil {
		log.Fatal(err)
	}
	mig.Steps(2)

	err = dbConn.Ping()
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		err := dbConn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	e := echo.New()
	m := _modelRepo.NewPostgresModelRepository(dbConn)

	timeoutContext := time.Duration(viper.GetInt("context.timeout")) * time.Second

	host := viper.GetString(`filestore.host`)
	region := viper.GetString(`filestore.region`)
	access := viper.GetString(`filestore.access`)
	secret := viper.GetString(`filestore.secret`)
	mbucket := viper.GetString(`model_bucket`)

	sess, err := session.NewSession(&aws.Config{
		Credentials:      credentials.NewStaticCredentials(access, secret, ""),
		Region:           aws.String(region),
		Endpoint:         aws.String(host),
		DisableSSL:       aws.Bool(true),
		S3ForcePathStyle: aws.Bool(true),
	})

	modelFileRepo := fileRepo.NewS3FileRepository(sess, mbucket)
	s := _modelService.NewModelService(m, modelFileRepo, timeoutContext)
	_modelHTTPController.NewModelHandler(e, s)

	log.Fatal(e.Start(viper.GetString("server.address")))
}

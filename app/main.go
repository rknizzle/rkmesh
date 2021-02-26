package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
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

	"github.com/rknizzle/rkmesh/domain"
	fileRepo "github.com/rknizzle/rkmesh/file/repository/s3"
	_modelHTTPController "github.com/rknizzle/rkmesh/model/controller/http"
	_modelRepo "github.com/rknizzle/rkmesh/model/repository/postgres"
	_modelService "github.com/rknizzle/rkmesh/model/service"
)

func init() {
	viper.SetConfigFile(`config.json`)
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("Failed to load config: %s\n", err.Error())
		os.Exit(1)
	}
}

func main() {
	dbConn, err := connectToDatabase(
		viper.GetString(`database.host`),
		viper.GetString(`database.port`),
		viper.GetString(`database.user`),
		viper.GetString(`database.pass`),
		viper.GetString(`database.name`),
	)
	if err != nil {
		fmt.Printf("Failed to connect to database: %s\n", err.Error())
		os.Exit(1)
	}

	runDatabaseMigrations(dbConn)

	defer func() {
		err := dbConn.Close()
		if err != nil {
			fmt.Printf("Failed to close database connection: %s\n", err.Error())
			os.Exit(1)
		}
	}()

	e := echo.New()
	m := _modelRepo.NewPostgresModelRepository(dbConn)

	modelFileStorage, err := connectToFileStorage(
		viper.GetString(`filestore.host`),
		viper.GetString(`filestore.region`),
		viper.GetString(`filestore.access`),
		viper.GetString(`filestore.secret`),
		viper.GetString(`model_bucket`),
	)
	if err != nil {
		fmt.Printf("Failed to connect to file storage: %s\n", err.Error())
		os.Exit(1)
	}

	timeoutContext := time.Duration(viper.GetInt("context.timeout")) * time.Second
	s := _modelService.NewModelService(m, modelFileStorage, timeoutContext)
	_modelHTTPController.NewModelHandler(e, s)

	log.Fatal(e.Start(viper.GetString("server.address")))
}

func connectToDatabase(dbHost, dbPort, dbUser, dbPass, dbName string) (*sql.DB, error) {
	connection := fmt.Sprintf(
		`host=%s port=%s user=%s
		password=%s dbname=%s sslmode=disable`,
		dbHost, dbPort, dbUser, dbPass, dbName)

	dbConn, err := sql.Open("postgres", connection)
	if err != nil {
		return nil, err
	}

	err = dbConn.Ping()
	if err != nil {
		return nil, err
	}
	return dbConn, nil
}

func runDatabaseMigrations(dbConn *sql.DB) error {
	driver, err := postgres.WithInstance(dbConn, &postgres.Config{})
	if err != nil {
		return err
	}
	mig, err := migrate.NewWithDatabaseInstance(
		"file://./migrations",
		"postgres", driver)
	if err != nil {
		return err
	}

	mig.Steps(2)
	return nil
}

func connectToFileStorage(host, region, access, secret, bucket string) (domain.FileRepository, error) {
	sess, err := session.NewSession(&aws.Config{
		Credentials:      credentials.NewStaticCredentials(access, secret, ""),
		Region:           aws.String(region),
		Endpoint:         aws.String(host),
		DisableSSL:       aws.Bool(true),
		S3ForcePathStyle: aws.Bool(true),
	})
	if err != nil {
		return nil, err
	}

	fileStorage := fileRepo.NewS3FileRepository(sess, bucket)
	return fileStorage, nil
}

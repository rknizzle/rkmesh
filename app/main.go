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

	"github.com/rknizzle/rkmesh/auth"
	"github.com/rknizzle/rkmesh/domain"
	"github.com/rknizzle/rkmesh/filestore"
	"github.com/rknizzle/rkmesh/model"
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
	m := model.NewPostgresModelRepository(dbConn)

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
	s := model.NewModelService(m, modelFileStorage, timeoutContext)
	model.NewModelHandler(e, s)

	// auth handling
	userRepo := auth.NewPostgresUserRepository(dbConn)
	authService := auth.NewAuthService(userRepo, timeoutContext)
	auth.NewAuthHandler(e, authService)

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

func connectToFileStorage(host, region, access, secret, bucket string) (domain.Filestore, error) {
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

	fileStorage := filestore.NewS3Filestore(sess, bucket)
	return fileStorage, nil
}

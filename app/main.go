package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"

	"github.com/rknizzle/rkmesh/auth"
	"github.com/rknizzle/rkmesh/domain"
	"github.com/rknizzle/rkmesh/filestore"
	"github.com/rknizzle/rkmesh/model"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("Failed to load config: %s\n", err.Error())
		os.Exit(1)
	}
}

func main() {
	dbConn, err := connectToDatabase(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_DATABASE"),
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
	timeoutInt, err := strconv.Atoi(os.Getenv("GENERIC_TIMEOUT"))
	if err != nil {
		fmt.Printf("Failed to parse timeout value: %s\n", err.Error())
		os.Exit(1)

	}
	timeoutContext := time.Duration(timeoutInt) * time.Second

	// auth handling
	userRepo := auth.NewPostgresUserRepository(dbConn)
	authService := auth.NewAuthService(userRepo, timeoutContext)
	auth.NewAuthHandler(e, authService)

	// models handling
	m := model.NewPostgresModelRepository(dbConn)

	modelFileStorage, err := connectToFileStorage(
		os.Getenv("FS_HOST"),
		os.Getenv("FS_REGION"),
		os.Getenv("FS_ACCESS_KEY"),
		os.Getenv("FS_SECRET_KEY"),
		os.Getenv("MODEL_BUCKET"),
	)
	if err != nil {
		fmt.Printf("Failed to connect to file storage: %s\n", err.Error())
		os.Exit(1)
	}

	s := model.NewModelService(m, modelFileStorage, timeoutContext)

	// Require a valid JWT token to access any /models routes
	modelRoutes := e.Group("/models")
	modelRoutes.Use(middleware.JWT([]byte(os.Getenv("JWT_SECRET_KEY"))))
	model.NewModelHandler(modelRoutes, s)

	log.Fatal(e.Start(":" + os.Getenv("PORT")))
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

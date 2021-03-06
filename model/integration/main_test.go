package integration

import (
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/rknizzle/rkmesh/filestore"
	"github.com/rknizzle/rkmesh/model"
	"github.com/rknizzle/rkmesh/testFileStore"
	"github.com/rknizzle/rkmesh/testdb"
)

// Sets up the application in a test state and ready to be hit with requests for integration tests.
// The application will be connected to a test database and a test file/object store. The handler
// and test storage instances are defined in global variables that are declared in model_test.go
func TestMain(m *testing.M) {
	var dbConn *sql.DB
	var err error
	tdb, dbConn, err = testdb.Open()
	if err != nil {
		fmt.Printf("Failed to connect to test database: %s\n", err.Error())
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

	// create the model repo which handles the interactions with the models database table
	mRepo := model.NewPostgresModelRepository(dbConn)

	// initialize the test file storage and get the filestore session that the app will connect
	// to during integration tests
	var sess *session.Session
	var bucket string
	_, sess, bucket, err = testFileStore.InitTestFileStore()
	if err != nil {
		fmt.Printf("Failed to connect to test file storage: %s\n", err.Error())
		os.Exit(1)
	}

	mFilestore := filestore.NewS3Filestore(sess, bucket)

	// create the model service
	timeoutContext := time.Duration(10) * time.Second
	s := model.NewModelService(mRepo, mFilestore, timeoutContext)

	// save the model handler to a global variable that will be used in the integration tests
	mHandler = model.ModelHandler{Service: s}

	exitVal := m.Run()
	os.Exit(exitVal)
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

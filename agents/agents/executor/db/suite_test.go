package db_test

import (
	"database/sql"
	"fmt"
	"github.com/Flaque/filet"
	. "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/synapsecns/sanguine/agents/agents/executor/db"
	"github.com/synapsecns/sanguine/agents/agents/executor/db/datastore/sql/mysql"
	"github.com/synapsecns/sanguine/agents/agents/executor/db/datastore/sql/sqlite"
	"github.com/synapsecns/sanguine/agents/agents/executor/metadata"
	"github.com/synapsecns/sanguine/core"
	"github.com/synapsecns/sanguine/core/metrics"
	"github.com/synapsecns/sanguine/core/metrics/localmetrics"
	"github.com/synapsecns/sanguine/core/testsuite"
	"gorm.io/gorm/schema"
	"os"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

type DBSuite struct {
	*testsuite.TestSuite
	dbs              []db.ExecutorDB
	logIndex         atomic.Int64
	metrics          metrics.Handler
	mysqlTablePrefix string
}

// NewEventDBSuite creates a new EventDBSuite.
func NewEventDBSuite(tb testing.TB) *DBSuite {
	tb.Helper()
	return &DBSuite{
		TestSuite: testsuite.NewTestSuite(tb),
		dbs:       []db.ExecutorDB{},
	}
}

func (t *DBSuite) SetupSuite() {
	t.TestSuite.SetupSuite()

	localmetrics.SetupTestJaeger(t.GetSuiteContext(), t.T())

	var err error
	t.metrics, err = metrics.NewByType(t.GetSuiteContext(), metadata.BuildInfo(), metrics.Jaeger)
	Nil(t.T(), err)
}

func (t *DBSuite) SetupTest() {
	t.TestSuite.SetupTest()

	t.logIndex.Store(0)

	sqliteStore, err := sqlite.NewSqliteStore(t.GetTestContext(), filet.TmpDir(t.T(), ""), t.metrics, false)
	Nil(t.T(), err)

	t.dbs = []db.ExecutorDB{sqliteStore}
	t.setupMysqlDB()
}

// connString gets the connection string.
func (t *DBSuite) connString(dbname string) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", core.GetEnv("MYSQL_USER", "root"), os.Getenv("MYSQL_PASSWORD"), core.GetEnv("MYSQL_HOST", "127.0.0.1"), core.GetEnvInt("MYSQL_PORT", 3306), dbname)
}

func (t *DBSuite) setupMysqlDB() {
	// skip if mysql test disabled, this really only needs to be run in ci

	// skip if mysql test disabled
	if os.Getenv("ENABLE_MYSQL_TEST") == "" {
		return
	}
	// sets up the conn string to the default database
	connString := t.connString(os.Getenv("MYSQL_DATABASE"))
	// sets up the mysql db
	testDB, err := sql.Open("mysql", connString)
	Nil(t.T(), err)
	// close the db once the connection is done
	defer func() {
		Nil(t.T(), testDB.Close())
	}()

	t.mysqlTablePrefix = fmt.Sprintf("test%d_%d_", t.GetTestID(), time.Now().Unix())

	// override the naming strategy to prevent tests from messing with each other.
	// todo this should be solved via a proper teardown process or transactions.
	mysql.NamingStrategy = schema.NamingStrategy{
		TablePrefix: t.mysqlTablePrefix,
	}

	mysql.MaxIdleConns = 10

	// create the sql store
	mysqlStore, err := mysql.NewMysqlStore(t.GetTestContext(), connString, t.metrics, false)
	Nil(t.T(), err)
	// add the db
	t.dbs = append(t.dbs, mysqlStore)
}

func (t *DBSuite) RunOnAllDBs(testFunc func(testDB db.ExecutorDB, tablePrefix string)) {
	t.T().Helper()

	wg := sync.WaitGroup{}
	for i, testDB := range t.dbs {
		var tablePrefix string
		wg.Add(1)
		// capture the value
		i := i
		// Mysql Check
		if i == 1 {
			tablePrefix = t.mysqlTablePrefix
		}
		go func(testDB db.ExecutorDB, tablePrefix string) {
			defer wg.Done()
			testFunc(testDB, tablePrefix)
		}(testDB, tablePrefix)
	}
	wg.Wait()
}

// TestDBSuite tests the db suite.
func TestEventDBSuite(t *testing.T) {
	suite.Run(t, NewEventDBSuite(t))
}

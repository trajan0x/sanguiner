package indexer_test

import (
	"github.com/synapsecns/sanguine/core/metrics"
	"github.com/synapsecns/sanguine/core/metrics/localmetrics"
	"github.com/synapsecns/sanguine/services/scribe/metadata"
	"testing"
	"time"

	"github.com/Flaque/filet"
	. "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/synapsecns/sanguine/core/testsuite"
	"github.com/synapsecns/sanguine/ethergo/signer/signer/localsigner"
	"github.com/synapsecns/sanguine/ethergo/signer/wallet"
	"github.com/synapsecns/sanguine/services/scribe/db"
	"github.com/synapsecns/sanguine/services/scribe/db/datastore/sql/sqlite"
	"github.com/synapsecns/sanguine/services/scribe/testutil"
)

type IndexerSuite struct {
	*testsuite.TestSuite
	testDB  db.EventDB
	manager *testutil.DeployManager
	wallet  wallet.Wallet
	signer  *localsigner.Signer
	metrics metrics.Handler
}

// NewIndexerSuite creates a new indexer test suite.
func NewIndexerSuite(tb testing.TB) *IndexerSuite {
	tb.Helper()
	return &IndexerSuite{
		TestSuite: testsuite.NewTestSuite(tb),
	}
}

// SetupTest sets up the test suite.
func (b *IndexerSuite) SetupTest() {
	b.TestSuite.SetupTest()
	b.SetTestTimeout(time.Minute * 3)
	sqliteStore, err := sqlite.NewSqliteStore(b.GetTestContext(), filet.TmpDir(b.T(), ""), b.metrics, false)
	Nil(b.T(), err)
	b.testDB = sqliteStore
	b.manager = testutil.NewDeployManager(b.T())
	b.wallet, err = wallet.FromRandom()
	Nil(b.T(), err)
	b.signer = localsigner.NewSigner(b.wallet.PrivateKey())
}

func (x *IndexerSuite) SetupSuite() {
	x.TestSuite.SetupSuite()
	localmetrics.SetupTestJaeger(x.GetSuiteContext(), x.T())

	var err error
	x.metrics, err = metrics.NewByType(x.GetSuiteContext(), metadata.BuildInfo(), metrics.Jaeger)
	Nil(x.T(), err)
}

// TestIndexerSuite tests the indexer suite.
func TestIndexerSuite(t *testing.T) {
	suite.Run(t, NewIndexerSuite(t))
}

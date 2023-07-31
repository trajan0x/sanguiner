package db_test

import (
	"math/big"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/ethereum/go-ethereum/common"
	. "github.com/stretchr/testify/assert"
	"github.com/synapsecns/sanguine/agents/agents/guard/db"
	"github.com/synapsecns/sanguine/agents/types"
)

func (t *DBSuite) TestGetUpdateAgentStatusParameters() {
	t.RunOnAllDBs(func(testDB db.GuardDB) {
		guardAddress := common.BigToAddress(big.NewInt(gofakeit.Int64()))

		addressA := common.BigToAddress(big.NewInt(gofakeit.Int64()))
		addressB := common.BigToAddress(big.NewInt(gofakeit.Int64()))
		addressC := common.BigToAddress(big.NewInt(gofakeit.Int64()))

		agentRootA := common.BigToHash(big.NewInt(gofakeit.Int64()))
		agentRootB := common.BigToHash(big.NewInt(gofakeit.Int64()))
		agentRootC := common.BigToHash(big.NewInt(gofakeit.Int64()))

		// Insert three rows into the `AgentTree` table.
		err := testDB.StoreAgentTree(
			t.GetTestContext(),
			agentRootA,
			addressA,
			gofakeit.Uint64(),
			[][32]byte{{gofakeit.Uint8()}},
		)
		Nil(t.T(), err)
		err = testDB.StoreAgentTree(
			t.GetTestContext(),
			agentRootB,
			addressB,
			gofakeit.Uint64(),
			[][32]byte{{gofakeit.Uint8()}},
		)
		Nil(t.T(), err)
		err = testDB.StoreAgentTree(
			t.GetTestContext(),
			agentRootC,
			addressC,
			gofakeit.Uint64(),
			[][32]byte{{gofakeit.Uint8()}},
		)
		Nil(t.T(), err)

		// Insert three rows into `Dispute`, two will have matching agent address to `AgentTree` rows and with status `Resolved`.
		err = testDB.StoreDispute(
			t.GetTestContext(),
			big.NewInt(gofakeit.Int64()),
			types.Resolved,
			guardAddress,
			gofakeit.Uint32(),
			addressA,
		)
		Nil(t.T(), err)
		err = testDB.StoreDispute(
			t.GetTestContext(),
			big.NewInt(gofakeit.Int64()),
			types.Resolved,
			guardAddress,
			gofakeit.Uint32(),
			addressB,
		)
		Nil(t.T(), err)
		err = testDB.StoreDispute(
			t.GetTestContext(),
			big.NewInt(gofakeit.Int64()),
			types.Opened,
			guardAddress,
			gofakeit.Uint32(),
			addressC,
		)
		Nil(t.T(), err)

		// Get the matching agent tree from the database.
		agentTrees, err := testDB.GetUpdateAgentStatusParameters(t.GetTestContext())
		Nil(t.T(), err)

		Equal(t.T(), 2, len(agentTrees))
	})
}

func (t *DBSuite) TestGetLatestConfirmedSummitBlockNumber() {
	t.RunOnAllDBs(func(testDB db.GuardDB) {
		chainIDA := gofakeit.Uint32()
		chainIDB := gofakeit.Uint32()

		agentRootA := common.BigToHash(big.NewInt(gofakeit.Int64()))
		agentRootB := common.BigToHash(big.NewInt(gofakeit.Int64()))

		err := testDB.StoreAgentRoot(
			t.GetTestContext(),
			agentRootA,
			chainIDA,
			gofakeit.Uint64(),
		)
		Nil(t.T(), err)
		err = testDB.StoreAgentRoot(
			t.GetTestContext(),
			agentRootB,
			chainIDB,
			gofakeit.Uint64(),
		)
		Nil(t.T(), err)

		err = testDB.StoreAgentTree(
			t.GetTestContext(),
			agentRootA,
			common.BigToAddress(big.NewInt(gofakeit.Int64())),
			5,
			[][32]byte{{gofakeit.Uint8()}},
		)
		Nil(t.T(), err)
		err = testDB.StoreAgentTree(
			t.GetTestContext(),
			agentRootB,
			common.BigToAddress(big.NewInt(gofakeit.Int64())),
			10,
			[][32]byte{{gofakeit.Uint8()}},
		)
		Nil(t.T(), err)

		blockNumberA, err := testDB.GetLatestConfirmedSummitBlockNumber(t.GetTestContext(), chainIDA)
		Nil(t.T(), err)

		Equal(t.T(), uint64(5), blockNumberA)

		blockNumberB, err := testDB.GetLatestConfirmedSummitBlockNumber(t.GetTestContext(), chainIDB)
		Nil(t.T(), err)

		Equal(t.T(), uint64(10), blockNumberB)
	})
}
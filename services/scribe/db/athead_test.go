package db_test

import (
	"github.com/brianvoe/gofakeit/v6"
	"github.com/ethereum/go-ethereum/common"
	. "github.com/stretchr/testify/assert"
	"github.com/synapsecns/sanguine/services/scribe/db"
	scribeTypes "github.com/synapsecns/sanguine/services/scribe/types"
	"math/big"
	"time"
)

func (t *DBSuite) TestUnconfirmedLogsQuery() {
	t.RunOnAllDBs(func(testDB db.EventDB) {
		chainID := gofakeit.Uint32()
		contractAddress := common.BigToAddress(big.NewInt(gofakeit.Int64()))
		const confirmedBlockHeight = 100
		const headBlock = 110
		for i := 1; i <= confirmedBlockHeight; i++ {
			txHash := common.BigToHash(big.NewInt(gofakeit.Int64()))
			log := t.MakeRandomLog(txHash)
			log.BlockNumber = uint64(i)
			log.Address = contractAddress
			// For testing, all confirmed txs will have an index of 1
			log.Index = 1
			err := testDB.StoreLogs(t.GetTestContext(), chainID, log)
			Nil(t.T(), err)
		}
		err := testDB.StoreLastIndexed(t.GetTestContext(), contractAddress, chainID, confirmedBlockHeight, scribeTypes.IndexingConfirmed)
		Nil(t.T(), err)

		// For testing, having the same txhash for all unconfirmed blocks.
		for i := confirmedBlockHeight + 1; i <= headBlock; i++ {
			txHash := common.BigToHash(big.NewInt(gofakeit.Int64()))

			log := t.MakeRandomLog(txHash)
			log.BlockNumber = uint64(i)
			log.TxHash = common.BigToHash(big.NewInt(gofakeit.Int64()))
			log.Address = contractAddress
			// For testing, all unconfirmed txs will have an index of 0
			log.Index = 0
			err := testDB.StoreLogsAtHead(t.GetTestContext(), chainID, log)
			Nil(t.T(), err)
		}

		logFilter := db.LogFilter{
			ChainID:         chainID,
			ContractAddress: contractAddress.String(),
		}
		logs, err := testDB.RetrieveLogsFromHeadRangeQuery(t.GetTestContext(), logFilter, 0, headBlock, 1)
		Nil(t.T(), err)
		Equal(t.T(), 100, len(logs))
		if len(logs) == 100 {
			Equal(t.T(), uint(0), logs[0].Index)
			// Check block range
			Equal(t.T(), uint64(110), logs[0].BlockNumber)
			Equal(t.T(), uint64(11), logs[99].BlockNumber)
			// check threshold of confirmed vs unconfirmed
			Equal(t.T(), uint(1), logs[10].Index)
			Equal(t.T(), uint(0), logs[9].Index)
		}
		logs, err = testDB.RetrieveLogsFromHeadRangeQuery(t.GetTestContext(), logFilter, 0, headBlock, 2)
		Nil(t.T(), err)

		Equal(t.T(), 10, len(logs))
		if len(logs) == 10 {
			// Check that these are confirmed logs
			Equal(t.T(), uint(1), logs[0].Index)
		}
	})
}

func (t *DBSuite) TestFlushLog() {
	t.RunOnAllDBs(func(testDB db.EventDB) {
		chainID := gofakeit.Uint32()
		contractAddress := common.BigToAddress(big.NewInt(gofakeit.Int64()))
		const deleteUpToBlock = 110
		const desiredBlockHeight = 200
		for i := 1; i <= deleteUpToBlock; i++ {
			txHash := common.BigToHash(big.NewInt(gofakeit.Int64()))
			log := t.MakeRandomLog(txHash)
			log.BlockNumber = uint64(i)
			log.Address = contractAddress

			// For testing, all to delete txs will have an index of 1
			log.Index = 1
			err := testDB.StoreLogsAtHead(t.GetTestContext(), chainID, log)
			Nil(t.T(), err)
		}
		time.Sleep(1 * time.Second)
		deleteTimestamp := time.Now().UnixNano()
		for i := deleteUpToBlock + 1; i <= desiredBlockHeight; i++ {
			txHash := common.BigToHash(big.NewInt(gofakeit.Int64()))

			log := t.MakeRandomLog(txHash)
			log.BlockNumber = uint64(i)
			log.TxHash = common.BigToHash(big.NewInt(gofakeit.Int64()))
			log.Address = contractAddress
			// For testing, all no delete txs will have an index of 0
			log.Index = 0
			err := testDB.StoreLogsAtHead(t.GetTestContext(), chainID, log)
			Nil(t.T(), err)
		}
		logFilter := db.LogFilter{
			ChainID:         chainID,
			ContractAddress: contractAddress.String(),
		}
		logs, err := testDB.RetrieveLogsFromHeadRangeQuery(t.GetTestContext(), logFilter, 0, desiredBlockHeight, 1)
		Nil(t.T(), err)
		Equal(t.T(), 100, len(logs))
		if 100 == len(logs) {
			Equal(t.T(), uint64(desiredBlockHeight), logs[0].BlockNumber)
		}
		err = testDB.FlushFromHeadTables(t.GetTestContext(), deleteTimestamp)
		Nil(t.T(), err)
		logs, err = testDB.RetrieveLogsFromHeadRangeQuery(t.GetTestContext(), logFilter, 0, desiredBlockHeight, 1)
		Nil(t.T(), err)
		Equal(t.T(), 90, len(logs))
		if len(logs) == 90 {
			// Check that the earliest log has a timestamp of 110
			Equal(t.T(), uint(0), logs[0].Index)
			Equal(t.T(), uint64(desiredBlockHeight), logs[0].BlockNumber)
		}
	})
}

func (t *DBSuite) TestUnconfirmedReceiptsQuery() {
	t.RunOnAllDBs(func(testDB db.EventDB) {
		chainID := gofakeit.Uint32()
		contractAddress := common.BigToAddress(big.NewInt(gofakeit.Int64()))
		const confirmedBlockHeight = 100
		const headBlock = 110
		for i := 1; i <= confirmedBlockHeight; i++ {
			txHash := common.BigToHash(big.NewInt(gofakeit.Int64()))
			receipt := t.MakeRandomReceipt(txHash)
			receipt.BlockNumber = big.NewInt(int64(i))
			receipt.ContractAddress = contractAddress
			// For testing, all confirmed receipts will have a status of 1
			receipt.Status = 1
			err := testDB.StoreReceipt(t.GetTestContext(), chainID, receipt)
			Nil(t.T(), err)
		}
		err := testDB.StoreLastIndexed(t.GetTestContext(), contractAddress, chainID, confirmedBlockHeight, scribeTypes.IndexingConfirmed)
		Nil(t.T(), err)

		// For testing, having the same txhash for all unconfirmed blocks.
		for i := confirmedBlockHeight + 1; i <= headBlock; i++ {
			txHash := common.BigToHash(big.NewInt(gofakeit.Int64()))

			receipt := t.MakeRandomReceipt(txHash)
			receipt.BlockNumber = big.NewInt(int64(i))
			receipt.ContractAddress = contractAddress
			// For testing, all confirmed receipts will have a status of 1
			receipt.Status = 0
			err := testDB.StoreReceiptAtHead(t.GetTestContext(), chainID, receipt)
			Nil(t.T(), err)
		}

		receiptFilter := db.ReceiptFilter{
			ChainID:         chainID,
			ContractAddress: contractAddress.String(),
		}
		receipts, err := testDB.RetrieveReceiptsFromHeadRangeQuery(t.GetTestContext(), receiptFilter, 0, headBlock, 1)
		Nil(t.T(), err)
		Equal(t.T(), 100, len(receipts))
		if len(receipts) == 100 {
			Equal(t.T(), uint64(0), receipts[0].Status)
			// Check block range
			Equal(t.T(), uint64(110), receipts[0].BlockNumber.Uint64())
			Equal(t.T(), uint64(11), receipts[99].BlockNumber.Uint64())
			// check threshold of confirmed vs unconfirmed
			Equal(t.T(), uint64(1), receipts[10].Status)
			Equal(t.T(), uint64(0), receipts[9].Status)
		}
		receipts, err = testDB.RetrieveReceiptsFromHeadRangeQuery(t.GetTestContext(), receiptFilter, 0, headBlock, 2)
		Nil(t.T(), err)

		Equal(t.T(), 10, len(receipts))
		if len(receipts) == 10 {
			// Check that these are confirmed logs
			Equal(t.T(), uint64(1), receipts[0].Status)
		}
	})
}
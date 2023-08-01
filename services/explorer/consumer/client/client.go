// Code generated by github.com/Yamashou/gqlgenc, DO NOT EDIT.

package client

import (
	"context"
	"net/http"

	"github.com/Yamashou/gqlgenc/client"
	"github.com/synapsecns/sanguine/services/explorer/consumer/client/model"
)

type Client struct {
	Client *client.Client
}

func NewClient(cli *http.Client, baseURL string, options ...client.HTTPRequestOption) *Client {
	return &Client{Client: client.NewClient(cli, baseURL, options...)}
}

type Query struct {
	Logs                     []*model.Log         "json:\"logs\" graphql:\"logs\""
	LogsRange                []*model.Log         "json:\"logsRange\" graphql:\"logsRange\""
	Receipts                 []*model.Receipt     "json:\"receipts\" graphql:\"receipts\""
	ReceiptsRange            []*model.Receipt     "json:\"receiptsRange\" graphql:\"receiptsRange\""
	Transactions             []*model.Transaction "json:\"transactions\" graphql:\"transactions\""
	TransactionsRange        []*model.Transaction "json:\"transactionsRange\" graphql:\"transactionsRange\""
	BlockTime                *int                 "json:\"blockTime\" graphql:\"blockTime\""
	LastStoredBlockNumber    *int                 "json:\"lastStoredBlockNumber\" graphql:\"lastStoredBlockNumber\""
	FirstStoredBlockNumber   *int                 "json:\"firstStoredBlockNumber\" graphql:\"firstStoredBlockNumber\""
	LastConfirmedBlockNumber *int                 "json:\"lastConfirmedBlockNumber\" graphql:\"lastConfirmedBlockNumber\""
	TxSender                 *string              "json:\"txSender\" graphql:\"txSender\""
	LastIndexed              *int                 "json:\"lastIndexed\" graphql:\"lastIndexed\""
	LogCount                 *int                 "json:\"logCount\" graphql:\"logCount\""
	ReceiptCount             *int                 "json:\"receiptCount\" graphql:\"receiptCount\""
	BlockTimeCount           *int                 "json:\"blockTimeCount\" graphql:\"blockTimeCount\""
	LogsAtHeadRange          []*model.Log         "json:\"logsAtHeadRange\" graphql:\"logsAtHeadRange\""
	ReceiptsAtHeadRange      []*model.Receipt     "json:\"receiptsAtHeadRange\" graphql:\"receiptsAtHeadRange\""
	TransactionsAtHeadRange  []*model.Transaction "json:\"transactionsAtHeadRange\" graphql:\"transactionsAtHeadRange\""
}
type GetLogsRange struct {
	Response []*struct {
		ContractAddress string   "json:\"contract_address\" graphql:\"contract_address\""
		ChainID         int      "json:\"chain_id\" graphql:\"chain_id\""
		Topics          []string "json:\"topics\" graphql:\"topics\""
		Data            string   "json:\"data\" graphql:\"data\""
		BlockNumber     int      "json:\"block_number\" graphql:\"block_number\""
		TxHash          string   "json:\"tx_hash\" graphql:\"tx_hash\""
		TxIndex         int      "json:\"tx_index\" graphql:\"tx_index\""
		BlockHash       string   "json:\"block_hash\" graphql:\"block_hash\""
		Index           int      "json:\"index\" graphql:\"index\""
		Removed         bool     "json:\"removed\" graphql:\"removed\""
	} "json:\"response\" graphql:\"response\""
}
type GetTransactions struct {
	Response []*struct {
		ChainID   int    "json:\"chain_id\" graphql:\"chain_id\""
		TxHash    string "json:\"tx_hash\" graphql:\"tx_hash\""
		Protected bool   "json:\"protected\" graphql:\"protected\""
		Type      int    "json:\"type\" graphql:\"type\""
		Data      string "json:\"data\" graphql:\"data\""
		Gas       int    "json:\"gas\" graphql:\"gas\""
		GasPrice  int    "json:\"gas_price\" graphql:\"gas_price\""
		GasTipCap string "json:\"gas_tip_cap\" graphql:\"gas_tip_cap\""
		GasFeeCap string "json:\"gas_fee_cap\" graphql:\"gas_fee_cap\""
		Value     string "json:\"value\" graphql:\"value\""
		Nonce     int    "json:\"nonce\" graphql:\"nonce\""
		To        string "json:\"to\" graphql:\"to\""
		Timestamp int    "json:\"timestamp\" graphql:\"timestamp\""
		Sender    string "json:\"sender\" graphql:\"sender\""
	} "json:\"response\" graphql:\"response\""
}
type GetBlockTime struct {
	Response *int "json:\"response\" graphql:\"response\""
}
type GetLastStoredBlockNumber struct {
	Response *int "json:\"response\" graphql:\"response\""
}
type GetFirstStoredBlockNumber struct {
	Response *int "json:\"response\" graphql:\"response\""
}
type GetTxSender struct {
	Response *string "json:\"response\" graphql:\"response\""
}
type GetLastIndexed struct {
	Response *int "json:\"response\" graphql:\"response\""
}
type GetLastConfirmedBlockNumber struct {
	Response *int "json:\"response\" graphql:\"response\""
}
type GetLogCount struct {
	Response *int "json:\"response\" graphql:\"response\""
}
type GetReceiptCount struct {
	Response *int "json:\"response\" graphql:\"response\""
}
type GetBlockTimeCount struct {
	Response *int "json:\"response\" graphql:\"response\""
}

const GetLogsRangeDocument = `query GetLogsRange ($chain_id: Int!, $start_block: Int!, $end_block: Int!, $page: Int!, $contract_address: String) {
	response: logsRange(chain_id: $chain_id, start_block: $start_block, end_block: $end_block, page: $page, contract_address: $contract_address) {
		contract_address
		chain_id
		topics
		data
		block_number
		tx_hash
		tx_index
		block_hash
		index
		removed
	}
}
`

func (c *Client) GetLogsRange(ctx context.Context, chainID int, startBlock int, endBlock int, page int, contractAddress *string, httpRequestOptions ...client.HTTPRequestOption) (*GetLogsRange, error) {
	vars := map[string]interface{}{
		"chain_id":         chainID,
		"start_block":      startBlock,
		"end_block":        endBlock,
		"page":             page,
		"contract_address": contractAddress,
	}

	var res GetLogsRange
	if err := c.Client.Post(ctx, "GetLogsRange", GetLogsRangeDocument, &res, vars, httpRequestOptions...); err != nil {
		return nil, err
	}

	return &res, nil
}

const GetTransactionsDocument = `query GetTransactions ($chain_id: Int!, $page: Int!, $tx_hash: String) {
	response: transactions(chain_id: $chain_id, page: $page, tx_hash: $tx_hash) {
		chain_id
		tx_hash
		protected
		type
		data
		gas
		gas_price
		gas_tip_cap
		gas_fee_cap
		value
		nonce
		to
		timestamp
		sender
	}
}
`

func (c *Client) GetTransactions(ctx context.Context, chainID int, page int, txHash *string, httpRequestOptions ...client.HTTPRequestOption) (*GetTransactions, error) {
	vars := map[string]interface{}{
		"chain_id": chainID,
		"page":     page,
		"tx_hash":  txHash,
	}

	var res GetTransactions
	if err := c.Client.Post(ctx, "GetTransactions", GetTransactionsDocument, &res, vars, httpRequestOptions...); err != nil {
		return nil, err
	}

	return &res, nil
}

const GetBlockTimeDocument = `query GetBlockTime ($chain_id: Int!, $block_number: Int!) {
	response: blockTime(chain_id: $chain_id, block_number: $block_number)
}
`

func (c *Client) GetBlockTime(ctx context.Context, chainID int, blockNumber int, httpRequestOptions ...client.HTTPRequestOption) (*GetBlockTime, error) {
	vars := map[string]interface{}{
		"chain_id":     chainID,
		"block_number": blockNumber,
	}

	var res GetBlockTime
	if err := c.Client.Post(ctx, "GetBlockTime", GetBlockTimeDocument, &res, vars, httpRequestOptions...); err != nil {
		return nil, err
	}

	return &res, nil
}

const GetLastStoredBlockNumberDocument = `query GetLastStoredBlockNumber ($chain_id: Int!) {
	response: lastStoredBlockNumber(chain_id: $chain_id)
}
`

func (c *Client) GetLastStoredBlockNumber(ctx context.Context, chainID int, httpRequestOptions ...client.HTTPRequestOption) (*GetLastStoredBlockNumber, error) {
	vars := map[string]interface{}{
		"chain_id": chainID,
	}

	var res GetLastStoredBlockNumber
	if err := c.Client.Post(ctx, "GetLastStoredBlockNumber", GetLastStoredBlockNumberDocument, &res, vars, httpRequestOptions...); err != nil {
		return nil, err
	}

	return &res, nil
}

const GetFirstStoredBlockNumberDocument = `query GetFirstStoredBlockNumber ($chain_id: Int!) {
	response: firstStoredBlockNumber(chain_id: $chain_id)
}
`

func (c *Client) GetFirstStoredBlockNumber(ctx context.Context, chainID int, httpRequestOptions ...client.HTTPRequestOption) (*GetFirstStoredBlockNumber, error) {
	vars := map[string]interface{}{
		"chain_id": chainID,
	}

	var res GetFirstStoredBlockNumber
	if err := c.Client.Post(ctx, "GetFirstStoredBlockNumber", GetFirstStoredBlockNumberDocument, &res, vars, httpRequestOptions...); err != nil {
		return nil, err
	}

	return &res, nil
}

const GetTxSenderDocument = `query GetTxSender ($chain_id: Int!, $tx_hash: String!) {
	response: txSender(chain_id: $chain_id, tx_hash: $tx_hash)
}
`

func (c *Client) GetTxSender(ctx context.Context, chainID int, txHash string, httpRequestOptions ...client.HTTPRequestOption) (*GetTxSender, error) {
	vars := map[string]interface{}{
		"chain_id": chainID,
		"tx_hash":  txHash,
	}

	var res GetTxSender
	if err := c.Client.Post(ctx, "GetTxSender", GetTxSenderDocument, &res, vars, httpRequestOptions...); err != nil {
		return nil, err
	}

	return &res, nil
}

const GetLastIndexedDocument = `query GetLastIndexed ($chain_id: Int!, $contract_address: String!) {
	response: lastIndexed(chain_id: $chain_id, contract_address: $contract_address)
}
`

func (c *Client) GetLastIndexed(ctx context.Context, chainID int, contractAddress string, httpRequestOptions ...client.HTTPRequestOption) (*GetLastIndexed, error) {
	vars := map[string]interface{}{
		"chain_id":         chainID,
		"contract_address": contractAddress,
	}

	var res GetLastIndexed
	if err := c.Client.Post(ctx, "GetLastIndexed", GetLastIndexedDocument, &res, vars, httpRequestOptions...); err != nil {
		return nil, err
	}

	return &res, nil
}

const GetLastConfirmedBlockNumberDocument = `query GetLastConfirmedBlockNumber ($chain_id: Int!) {
	response: lastConfirmedBlockNumber(chain_id: $chain_id)
}
`

func (c *Client) GetLastConfirmedBlockNumber(ctx context.Context, chainID int, httpRequestOptions ...client.HTTPRequestOption) (*GetLastConfirmedBlockNumber, error) {
	vars := map[string]interface{}{
		"chain_id": chainID,
	}

	var res GetLastConfirmedBlockNumber
	if err := c.Client.Post(ctx, "GetLastConfirmedBlockNumber", GetLastConfirmedBlockNumberDocument, &res, vars, httpRequestOptions...); err != nil {
		return nil, err
	}

	return &res, nil
}

const GetLogCountDocument = `query GetLogCount ($chain_id: Int!, $contract_address: String!) {
	response: logCount(chain_id: $chain_id, contract_address: $contract_address)
}
`

func (c *Client) GetLogCount(ctx context.Context, chainID int, contractAddress string, httpRequestOptions ...client.HTTPRequestOption) (*GetLogCount, error) {
	vars := map[string]interface{}{
		"chain_id":         chainID,
		"contract_address": contractAddress,
	}

	var res GetLogCount
	if err := c.Client.Post(ctx, "GetLogCount", GetLogCountDocument, &res, vars, httpRequestOptions...); err != nil {
		return nil, err
	}

	return &res, nil
}

const GetReceiptCountDocument = `query GetReceiptCount ($chain_id: Int!) {
	response: receiptCount(chain_id: $chain_id)
}
`

func (c *Client) GetReceiptCount(ctx context.Context, chainID int, httpRequestOptions ...client.HTTPRequestOption) (*GetReceiptCount, error) {
	vars := map[string]interface{}{
		"chain_id": chainID,
	}

	var res GetReceiptCount
	if err := c.Client.Post(ctx, "GetReceiptCount", GetReceiptCountDocument, &res, vars, httpRequestOptions...); err != nil {
		return nil, err
	}

	return &res, nil
}

const GetBlockTimeCountDocument = `query GetBlockTimeCount ($chain_id: Int!) {
	response: blockTimeCount(chain_id: $chain_id)
}
`

func (c *Client) GetBlockTimeCount(ctx context.Context, chainID int, httpRequestOptions ...client.HTTPRequestOption) (*GetBlockTimeCount, error) {
	vars := map[string]interface{}{
		"chain_id": chainID,
	}

	var res GetBlockTimeCount
	if err := c.Client.Post(ctx, "GetBlockTimeCount", GetBlockTimeCountDocument, &res, vars, httpRequestOptions...); err != nil {
		return nil, err
	}

	return &res, nil
}

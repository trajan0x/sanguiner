// Code generated by github.com/Yamashou/gqlgenc, DO NOT EDIT.

package client

import (
	"context"
	"net/http"

	"github.com/Yamashou/gqlgenc/client"
	"github.com/synapsecns/sanguine/services/explorer/graphql/server/graph/model"
)

type Client struct {
	Client *client.Client
}

func NewClient(cli *http.Client, baseURL string, options ...client.HTTPRequestOption) *Client {
	return &Client{Client: client.NewClient(cli, baseURL, options...)}
}

type Query struct {
	BridgeTransactions     []*model.BridgeTransaction      "json:\"bridgeTransactions\" graphql:\"bridgeTransactions\""
	MessageBusTransactions []*model.MessageBusTransaction  "json:\"messageBusTransactions\" graphql:\"messageBusTransactions\""
	CountByChainID         []*model.TransactionCountResult "json:\"countByChainId\" graphql:\"countByChainId\""
	CountByTokenAddress    []*model.TokenCountResult       "json:\"countByTokenAddress\" graphql:\"countByTokenAddress\""
	AddressRanking         []*model.AddressRanking         "json:\"addressRanking\" graphql:\"addressRanking\""
	AmountStatistic        *model.ValueResult              "json:\"amountStatistic\" graphql:\"amountStatistic\""
	DailyStatistics        *model.DailyResult              "json:\"dailyStatistics\" graphql:\"dailyStatistics\""
	DailyStatisticsByChain []*model.DateResultByChain      "json:\"dailyStatisticsByChain\" graphql:\"dailyStatisticsByChain\""
	RankedChainIDsByVolume []*model.VolumeByChainID        "json:\"rankedChainIDsByVolume\" graphql:\"rankedChainIDsByVolume\""
}
type GetBridgeTransactions struct {
	Response []*struct {
		FromInfo *struct {
			ChainID        *int     "json:\"chainID\" graphql:\"chainID\""
			Address        *string  "json:\"address\" graphql:\"address\""
			TxnHash        *string  "json:\"txnHash\" graphql:\"txnHash\""
			Value          *string  "json:\"value\" graphql:\"value\""
			FormattedValue *float64 "json:\"formattedValue\" graphql:\"formattedValue\""
			USDValue       *float64 "json:\"USDValue\" graphql:\"USDValue\""
			TokenAddress   *string  "json:\"tokenAddress\" graphql:\"tokenAddress\""
			TokenSymbol    *string  "json:\"tokenSymbol\" graphql:\"tokenSymbol\""
			BlockNumber    *int     "json:\"blockNumber\" graphql:\"blockNumber\""
			Time           *int     "json:\"time\" graphql:\"time\""
			FormattedTime  *string  "json:\"formattedTime\" graphql:\"formattedTime\""
		} "json:\"fromInfo\" graphql:\"fromInfo\""
		ToInfo *struct {
			ChainID        *int     "json:\"chainID\" graphql:\"chainID\""
			Address        *string  "json:\"address\" graphql:\"address\""
			TxnHash        *string  "json:\"txnHash\" graphql:\"txnHash\""
			Value          *string  "json:\"value\" graphql:\"value\""
			FormattedValue *float64 "json:\"formattedValue\" graphql:\"formattedValue\""
			USDValue       *float64 "json:\"USDValue\" graphql:\"USDValue\""
			TokenAddress   *string  "json:\"tokenAddress\" graphql:\"tokenAddress\""
			TokenSymbol    *string  "json:\"tokenSymbol\" graphql:\"tokenSymbol\""
			BlockNumber    *int     "json:\"blockNumber\" graphql:\"blockNumber\""
			Time           *int     "json:\"time\" graphql:\"time\""
			FormattedTime  *string  "json:\"formattedTime\" graphql:\"formattedTime\""
		} "json:\"toInfo\" graphql:\"toInfo\""
		Kappa       *string "json:\"kappa\" graphql:\"kappa\""
		Pending     *bool   "json:\"pending\" graphql:\"pending\""
		SwapSuccess *bool   "json:\"swapSuccess\" graphql:\"swapSuccess\""
	} "json:\"response\" graphql:\"response\""
}
type GetCountByChainID struct {
	Response []*struct {
		Count   *int "json:\"count\" graphql:\"count\""
		ChainID *int "json:\"chainID\" graphql:\"chainID\""
	} "json:\"response\" graphql:\"response\""
}
type GetCountByTokenAddress struct {
	Response []*struct {
		ChainID      *int    "json:\"chainID\" graphql:\"chainID\""
		TokenAddress *string "json:\"tokenAddress\" graphql:\"tokenAddress\""
		Count        *int    "json:\"count\" graphql:\"count\""
	} "json:\"response\" graphql:\"response\""
}
type GetAddressRanking struct {
	Response []*struct {
		Address *string "json:\"address\" graphql:\"address\""
		Count   *int    "json:\"count\" graphql:\"count\""
	} "json:\"response\" graphql:\"response\""
}
type GetRankedChainIDsByVolume struct {
	Response []*struct {
		ChainID *int     "json:\"chainID\" graphql:\"chainID\""
		Total   *float64 "json:\"total\" graphql:\"total\""
	} "json:\"response\" graphql:\"response\""
}
type GetAmountStatistic struct {
	Response *struct {
		Value *string "json:\"value\" graphql:\"value\""
	} "json:\"response\" graphql:\"response\""
}
type GetDailyStatisticsByChain struct {
	Response []*struct {
		Date      *string  "json:\"date\" graphql:\"date\""
		Ethereum  *float64 "json:\"ethereum\" graphql:\"ethereum\""
		Optimism  *float64 "json:\"optimism\" graphql:\"optimism\""
		Cronos    *float64 "json:\"cronos\" graphql:\"cronos\""
		Bsc       *float64 "json:\"bsc\" graphql:\"bsc\""
		Polygon   *float64 "json:\"polygon\" graphql:\"polygon\""
		Fantom    *float64 "json:\"fantom\" graphql:\"fantom\""
		Boba      *float64 "json:\"boba\" graphql:\"boba\""
		Metis     *float64 "json:\"metis\" graphql:\"metis\""
		Moonbeam  *float64 "json:\"moonbeam\" graphql:\"moonbeam\""
		Moonriver *float64 "json:\"moonriver\" graphql:\"moonriver\""
		Klaytn    *float64 "json:\"klaytn\" graphql:\"klaytn\""
		Arbitrum  *float64 "json:\"arbitrum\" graphql:\"arbitrum\""
		Avalanche *float64 "json:\"avalanche\" graphql:\"avalanche\""
		Dfk       *float64 "json:\"dfk\" graphql:\"dfk\""
		Aurora    *float64 "json:\"aurora\" graphql:\"aurora\""
		Harmony   *float64 "json:\"harmony\" graphql:\"harmony\""
		Canto     *float64 "json:\"canto\" graphql:\"canto\""
		Total     *float64 "json:\"total\" graphql:\"total\""
	} "json:\"response\" graphql:\"response\""
}
type GetDailyStatistics struct {
	Response *struct {
		Total       *float64 "json:\"total\" graphql:\"total\""
		DateResults []*struct {
			Date  *string  "json:\"date\" graphql:\"date\""
			Total *float64 "json:\"total\" graphql:\"total\""
		} "json:\"dateResults\" graphql:\"dateResults\""
		Type *model.DailyStatisticType "json:\"type\" graphql:\"type\""
	} "json:\"response\" graphql:\"response\""
}
type GetMessageBusTransactions struct {
	Response []*struct {
		FromInfo *struct {
			ChainID              *int    "json:\"chainID\" graphql:\"chainID\""
			ChainName            *string "json:\"chainName\" graphql:\"chainName\""
			DestinationChainID   *int    "json:\"destinationChainID\" graphql:\"destinationChainID\""
			DestinationChainName *string "json:\"destinationChainName\" graphql:\"destinationChainName\""
			ContractAddress      *string "json:\"contractAddress\" graphql:\"contractAddress\""
			TxnHash              *string "json:\"txnHash\" graphql:\"txnHash\""
			Message              *string "json:\"message\" graphql:\"message\""
			BlockNumber          *int    "json:\"blockNumber\" graphql:\"blockNumber\""
			Time                 *int    "json:\"time\" graphql:\"time\""
			FormattedTime        *string "json:\"formattedTime\" graphql:\"formattedTime\""
		} "json:\"fromInfo\" graphql:\"fromInfo\""
		ToInfo *struct {
			ChainID         *int    "json:\"chainID\" graphql:\"chainID\""
			ChainName       *string "json:\"chainName\" graphql:\"chainName\""
			ContractAddress *string "json:\"contractAddress\" graphql:\"contractAddress\""
			TxnHash         *string "json:\"txnHash\" graphql:\"txnHash\""
			Message         *string "json:\"message\" graphql:\"message\""
			BlockNumber     *int    "json:\"blockNumber\" graphql:\"blockNumber\""
			Time            *int    "json:\"time\" graphql:\"time\""
			FormattedTime   *string "json:\"formattedTime\" graphql:\"formattedTime\""
		} "json:\"toInfo\" graphql:\"toInfo\""
		MessageID *string "json:\"messageID\" graphql:\"messageID\""
		Pending   *bool   "json:\"pending\" graphql:\"pending\""
	} "json:\"response\" graphql:\"response\""
}

const GetBridgeTransactionsDocument = `query GetBridgeTransactions ($chainID: [Int], $address: String, $maxAmount: Int, $minAmount: Int, $startTime: Int, $endTime: Int, $txHash: String, $kappa: String, $pending: Boolean, $page: Int, $tokenAddress: [String]) {
	response: bridgeTransactions(chainID: $chainID, address: $address, maxAmount: $maxAmount, minAmount: $minAmount, startTime: $startTime, endTime: $endTime, txnHash: $txHash, kappa: $kappa, pending: $pending, page: $page, tokenAddress: $tokenAddress) {
		fromInfo {
			chainID
			address
			txnHash
			value
			formattedValue
			USDValue
			tokenAddress
			tokenSymbol
			blockNumber
			time
			formattedTime
		}
		toInfo {
			chainID
			address
			txnHash
			value
			formattedValue
			USDValue
			tokenAddress
			tokenSymbol
			blockNumber
			time
			formattedTime
		}
		kappa
		pending
		swapSuccess
	}
}
`

func (c *Client) GetBridgeTransactions(ctx context.Context, chainID []*int, address *string, maxAmount *int, minAmount *int, startTime *int, endTime *int, txHash *string, kappa *string, pending *bool, page *int, tokenAddress []*string, httpRequestOptions ...client.HTTPRequestOption) (*GetBridgeTransactions, error) {
	vars := map[string]interface{}{
		"chainID":      chainID,
		"address":      address,
		"maxAmount":    maxAmount,
		"minAmount":    minAmount,
		"startTime":    startTime,
		"endTime":      endTime,
		"txHash":       txHash,
		"kappa":        kappa,
		"pending":      pending,
		"page":         page,
		"tokenAddress": tokenAddress,
	}

	var res GetBridgeTransactions
	if err := c.Client.Post(ctx, "GetBridgeTransactions", GetBridgeTransactionsDocument, &res, vars, httpRequestOptions...); err != nil {
		return nil, err
	}

	return &res, nil
}

const GetCountByChainIDDocument = `query GetCountByChainId ($chainID: Int, $address: String, $direction: Direction, $hours: Int) {
	response: countByChainId(chainID: $chainID, address: $address, direction: $direction, hours: $hours) {
		count
		chainID
	}
}
`

func (c *Client) GetCountByChainID(ctx context.Context, chainID *int, address *string, direction *model.Direction, hours *int, httpRequestOptions ...client.HTTPRequestOption) (*GetCountByChainID, error) {
	vars := map[string]interface{}{
		"chainID":   chainID,
		"address":   address,
		"direction": direction,
		"hours":     hours,
	}

	var res GetCountByChainID
	if err := c.Client.Post(ctx, "GetCountByChainId", GetCountByChainIDDocument, &res, vars, httpRequestOptions...); err != nil {
		return nil, err
	}

	return &res, nil
}

const GetCountByTokenAddressDocument = `query GetCountByTokenAddress ($chainID: Int, $address: String, $direction: Direction, $hours: Int) {
	response: countByTokenAddress(chainID: $chainID, address: $address, direction: $direction, hours: $hours) {
		chainID
		tokenAddress
		count
	}
}
`

func (c *Client) GetCountByTokenAddress(ctx context.Context, chainID *int, address *string, direction *model.Direction, hours *int, httpRequestOptions ...client.HTTPRequestOption) (*GetCountByTokenAddress, error) {
	vars := map[string]interface{}{
		"chainID":   chainID,
		"address":   address,
		"direction": direction,
		"hours":     hours,
	}

	var res GetCountByTokenAddress
	if err := c.Client.Post(ctx, "GetCountByTokenAddress", GetCountByTokenAddressDocument, &res, vars, httpRequestOptions...); err != nil {
		return nil, err
	}

	return &res, nil
}

const GetAddressRankingDocument = `query GetAddressRanking ($hours: Int) {
	response: addressRanking(hours: $hours) {
		address
		count
	}
}
`

func (c *Client) GetAddressRanking(ctx context.Context, hours *int, httpRequestOptions ...client.HTTPRequestOption) (*GetAddressRanking, error) {
	vars := map[string]interface{}{
		"hours": hours,
	}

	var res GetAddressRanking
	if err := c.Client.Post(ctx, "GetAddressRanking", GetAddressRankingDocument, &res, vars, httpRequestOptions...); err != nil {
		return nil, err
	}

	return &res, nil
}

const GetRankedChainIDsByVolumeDocument = `query GetRankedChainIDsByVolume ($duration: Duration) {
	response: rankedChainIDsByVolume(duration: $duration) {
		chainID
		total
	}
}
`

func (c *Client) GetRankedChainIDsByVolume(ctx context.Context, duration *model.Duration, httpRequestOptions ...client.HTTPRequestOption) (*GetRankedChainIDsByVolume, error) {
	vars := map[string]interface{}{
		"duration": duration,
	}

	var res GetRankedChainIDsByVolume
	if err := c.Client.Post(ctx, "GetRankedChainIDsByVolume", GetRankedChainIDsByVolumeDocument, &res, vars, httpRequestOptions...); err != nil {
		return nil, err
	}

	return &res, nil
}

const GetAmountStatisticDocument = `query GetAmountStatistic ($type: StatisticType!, $platform: Platform, $duration: Duration, $chainID: Int, $address: String, $tokenAddress: String) {
	response: amountStatistic(type: $type, duration: $duration, platform: $platform, chainID: $chainID, address: $address, tokenAddress: $tokenAddress) {
		value
	}
}
`

func (c *Client) GetAmountStatistic(ctx context.Context, typeArg model.StatisticType, platform *model.Platform, duration *model.Duration, chainID *int, address *string, tokenAddress *string, httpRequestOptions ...client.HTTPRequestOption) (*GetAmountStatistic, error) {
	vars := map[string]interface{}{
		"type":         typeArg,
		"platform":     platform,
		"duration":     duration,
		"chainID":      chainID,
		"address":      address,
		"tokenAddress": tokenAddress,
	}

	var res GetAmountStatistic
	if err := c.Client.Post(ctx, "GetAmountStatistic", GetAmountStatisticDocument, &res, vars, httpRequestOptions...); err != nil {
		return nil, err
	}

	return &res, nil
}

const GetDailyStatisticsByChainDocument = `query GetDailyStatisticsByChain ($chainID: Int, $type: DailyStatisticType, $duration: Duration) {
	response: dailyStatisticsByChain(chainID: $chainID, type: $type, duration: $duration) {
		date
		ethereum
		optimism
		cronos
		bsc
		polygon
		fantom
		boba
		metis
		moonbeam
		moonriver
		klaytn
		arbitrum
		avalanche
		dfk
		aurora
		harmony
		canto
		total
	}
}
`

func (c *Client) GetDailyStatisticsByChain(ctx context.Context, chainID *int, typeArg *model.DailyStatisticType, duration *model.Duration, httpRequestOptions ...client.HTTPRequestOption) (*GetDailyStatisticsByChain, error) {
	vars := map[string]interface{}{
		"chainID":  chainID,
		"type":     typeArg,
		"duration": duration,
	}

	var res GetDailyStatisticsByChain
	if err := c.Client.Post(ctx, "GetDailyStatisticsByChain", GetDailyStatisticsByChainDocument, &res, vars, httpRequestOptions...); err != nil {
		return nil, err
	}

	return &res, nil
}

const GetDailyStatisticsDocument = `query GetDailyStatistics ($chainID: Int, $type: DailyStatisticType, $platform: Platform, $days: Int) {
	response: dailyStatistics(chainID: $chainID, type: $type, days: $days, platform: $platform) {
		total
		dateResults {
			date
			total
		}
		type
	}
}
`

func (c *Client) GetDailyStatistics(ctx context.Context, chainID *int, typeArg *model.DailyStatisticType, platform *model.Platform, days *int, httpRequestOptions ...client.HTTPRequestOption) (*GetDailyStatistics, error) {
	vars := map[string]interface{}{
		"chainID":  chainID,
		"type":     typeArg,
		"platform": platform,
		"days":     days,
	}

	var res GetDailyStatistics
	if err := c.Client.Post(ctx, "GetDailyStatistics", GetDailyStatisticsDocument, &res, vars, httpRequestOptions...); err != nil {
		return nil, err
	}

	return &res, nil
}

const GetMessageBusTransactionsDocument = `query GetMessageBusTransactions ($chainID: [Int], $contractAddress: String, $startTime: Int, $endTime: Int, $txHash: String, $messageID: String, $pending: Boolean, $page: Int) {
	response: messageBusTransactions(chainID: $chainID, contractAddress: $contractAddress, startTime: $startTime, endTime: $endTime, txnHash: $txHash, messageID: $messageID, pending: $pending, page: $page) {
		fromInfo {
			chainID
			chainName
			destinationChainID
			destinationChainName
			contractAddress
			txnHash
			message
			blockNumber
			time
			formattedTime
		}
		toInfo {
			chainID
			chainName
			contractAddress
			txnHash
			message
			blockNumber
			time
			formattedTime
		}
		messageID
		pending
	}
}
`

func (c *Client) GetMessageBusTransactions(ctx context.Context, chainID []*int, contractAddress *string, startTime *int, endTime *int, txHash *string, messageID *string, pending *bool, page *int, httpRequestOptions ...client.HTTPRequestOption) (*GetMessageBusTransactions, error) {
	vars := map[string]interface{}{
		"chainID":         chainID,
		"contractAddress": contractAddress,
		"startTime":       startTime,
		"endTime":         endTime,
		"txHash":          txHash,
		"messageID":       messageID,
		"pending":         pending,
		"page":            page,
	}

	var res GetMessageBusTransactions
	if err := c.Client.Post(ctx, "GetMessageBusTransactions", GetMessageBusTransactionsDocument, &res, vars, httpRequestOptions...); err != nil {
		return nil, err
	}

	return &res, nil
}

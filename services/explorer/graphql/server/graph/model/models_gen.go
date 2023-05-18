// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"fmt"
	"io"
	"strconv"
)

type MessageType interface {
	IsMessageType()
}

type AddressChainRanking struct {
	ChainID *int `json:"chain_id"`
	Count   *int `json:"count"`
}

type AddressDailyCount struct {
	Date  *string `json:"date"`
	Count *int    `json:"count"`
}

type AddressData struct {
	BridgeVolume *float64               `json:"bridgeVolume"`
	BridgeFees   *float64               `json:"bridgeFees"`
	BridgeTxs    *int                   `json:"bridgeTxs"`
	SwapVolume   *float64               `json:"swapVolume"`
	SwapFees     *float64               `json:"swapFees"`
	SwapTxs      *int                   `json:"swapTxs"`
	Rank         *int                   `json:"rank"`
	EarliestTx   *int                   `json:"earliestTx"`
	ChainRanking []*AddressChainRanking `json:"chainRanking"`
	DailyData    []*AddressDailyCount   `json:"dailyData"`
}

// AddressRanking gives the amount of transactions that occurred for a specific address across all chains.
type AddressRanking struct {
	Address *string `json:"address"`
	Count   *int    `json:"count"`
}

// BridgeTransaction represents an entire bridge transaction, including both
// to and from transactions. If a `from` transaction does not have a corresponding
// `to` transaction, `pending` will be true.
type BridgeTransaction struct {
	FromInfo    *PartialInfo `json:"fromInfo"`
	ToInfo      *PartialInfo `json:"toInfo"`
	Kappa       *string      `json:"kappa"`
	Pending     *bool        `json:"pending"`
	SwapSuccess *bool        `json:"swapSuccess"`
}

// DateResult is a given statistic for a given date.
type DateResult struct {
	Date  *string  `json:"date"`
	Total *float64 `json:"total"`
}

// DateResult is a given statistic for a given date.
type DateResultByChain struct {
	Date      *string  `json:"date"`
	Ethereum  *float64 `json:"ethereum"`
	Optimism  *float64 `json:"optimism"`
	Cronos    *float64 `json:"cronos"`
	Bsc       *float64 `json:"bsc"`
	Polygon   *float64 `json:"polygon"`
	Fantom    *float64 `json:"fantom"`
	Boba      *float64 `json:"boba"`
	Metis     *float64 `json:"metis"`
	Moonbeam  *float64 `json:"moonbeam"`
	Moonriver *float64 `json:"moonriver"`
	Klaytn    *float64 `json:"klaytn"`
	Arbitrum  *float64 `json:"arbitrum"`
	Avalanche *float64 `json:"avalanche"`
	Dfk       *float64 `json:"dfk"`
	Aurora    *float64 `json:"aurora"`
	Harmony   *float64 `json:"harmony"`
	Canto     *float64 `json:"canto"`
	Dogechain *float64 `json:"dogechain"`
	Total     *float64 `json:"total"`
}

type HeroType struct {
	Recipient string `json:"recipient"`
	HeroID    string `json:"heroID"`
}

func (HeroType) IsMessageType() {}

// HistoricalResult is a given statistic for dates.
type HistoricalResult struct {
	Total       *float64              `json:"total"`
	DateResults []*DateResult         `json:"dateResults"`
	Type        *HistoricalResultType `json:"type"`
}

type Leaderboard struct {
	Address      *string  `json:"address"`
	VolumeUsd    *float64 `json:"volumeUSD"`
	Fees         *float64 `json:"fees"`
	Txs          *int     `json:"txs"`
	Rank         *int     `json:"rank"`
	AvgVolumeUsd *float64 `json:"avgVolumeUSD"`
}

type MessageBusTransaction struct {
	FromInfo  *PartialMessageBusInfo `json:"fromInfo"`
	ToInfo    *PartialMessageBusInfo `json:"toInfo"`
	Pending   *bool                  `json:"pending"`
	MessageID *string                `json:"messageID"`
}

// PartialInfo is a transaction that occurred on one chain.
type PartialInfo struct {
	ChainID            *int     `json:"chainID"`
	DestinationChainID *int     `json:"destinationChainID"`
	Address            *string  `json:"address"`
	TxnHash            *string  `json:"txnHash"`
	Value              *string  `json:"value"`
	FormattedValue     *float64 `json:"formattedValue"`
	USDValue           *float64 `json:"USDValue"`
	TokenAddress       *string  `json:"tokenAddress"`
	TokenSymbol        *string  `json:"tokenSymbol"`
	BlockNumber        *int     `json:"blockNumber"`
	Time               *int     `json:"time"`
	FormattedTime      *string  `json:"formattedTime"`
}

type PartialMessageBusInfo struct {
	ChainID              *int        `json:"chainID"`
	ChainName            *string     `json:"chainName"`
	DestinationChainID   *int        `json:"destinationChainID"`
	DestinationChainName *string     `json:"destinationChainName"`
	ContractAddress      *string     `json:"contractAddress"`
	TxnHash              *string     `json:"txnHash"`
	Message              *string     `json:"message"`
	MessageType          MessageType `json:"messageType"`
	BlockNumber          *int        `json:"blockNumber"`
	Time                 *int        `json:"time"`
	FormattedTime        *string     `json:"formattedTime"`
	RevertedReason       *string     `json:"revertedReason"`
}

type PetType struct {
	Recipient string `json:"recipient"`
	PetID     string `json:"petID"`
	Name      string `json:"name"`
}

func (PetType) IsMessageType() {}

type TearType struct {
	Recipient string `json:"recipient"`
	Amount    string `json:"amount"`
}

func (TearType) IsMessageType() {}

// TokenCountResult gives the amount of transactions that occurred for a specific token, separated by chain ID.
type TokenCountResult struct {
	ChainID      *int    `json:"chainID"`
	TokenAddress *string `json:"tokenAddress"`
	Count        *int    `json:"count"`
}

// TransactionCountResult gives the amount of transactions that occurred for a specific chain ID.
type TransactionCountResult struct {
	ChainID *int `json:"chainID"`
	Count   *int `json:"count"`
}

type UnknownType struct {
	Known bool `json:"known"`
}

func (UnknownType) IsMessageType() {}

// ValueResult is a value result of either USD or numeric value.
type ValueResult struct {
	Value *string `json:"value"`
}

type VolumeByChainID struct {
	ChainID *int     `json:"chainID"`
	Total   *float64 `json:"total"`
}

type DailyStatisticType string

const (
	DailyStatisticTypeVolume       DailyStatisticType = "VOLUME"
	DailyStatisticTypeTransactions DailyStatisticType = "TRANSACTIONS"
	DailyStatisticTypeAddresses    DailyStatisticType = "ADDRESSES"
	DailyStatisticTypeFee          DailyStatisticType = "FEE"
)

var AllDailyStatisticType = []DailyStatisticType{
	DailyStatisticTypeVolume,
	DailyStatisticTypeTransactions,
	DailyStatisticTypeAddresses,
	DailyStatisticTypeFee,
}

func (e DailyStatisticType) IsValid() bool {
	switch e {
	case DailyStatisticTypeVolume, DailyStatisticTypeTransactions, DailyStatisticTypeAddresses, DailyStatisticTypeFee:
		return true
	}
	return false
}

func (e DailyStatisticType) String() string {
	return string(e)
}

func (e *DailyStatisticType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = DailyStatisticType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid DailyStatisticType", str)
	}
	return nil
}

func (e DailyStatisticType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type Direction string

const (
	DirectionIn  Direction = "IN"
	DirectionOut Direction = "OUT"
)

var AllDirection = []Direction{
	DirectionIn,
	DirectionOut,
}

func (e Direction) IsValid() bool {
	switch e {
	case DirectionIn, DirectionOut:
		return true
	}
	return false
}

func (e Direction) String() string {
	return string(e)
}

func (e *Direction) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Direction(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Direction", str)
	}
	return nil
}

func (e Direction) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type Duration string

const (
	DurationPastDay     Duration = "PAST_DAY"
	DurationPastMonth   Duration = "PAST_MONTH"
	DurationPast3Months Duration = "PAST_3_MONTHS"
	DurationPast6Months Duration = "PAST_6_MONTHS"
	DurationPastYear    Duration = "PAST_YEAR"
	DurationAllTime     Duration = "ALL_TIME"
)

var AllDuration = []Duration{
	DurationPastDay,
	DurationPastMonth,
	DurationPast3Months,
	DurationPast6Months,
	DurationPastYear,
	DurationAllTime,
}

func (e Duration) IsValid() bool {
	switch e {
	case DurationPastDay, DurationPastMonth, DurationPast3Months, DurationPast6Months, DurationPastYear, DurationAllTime:
		return true
	}
	return false
}

func (e Duration) String() string {
	return string(e)
}

func (e *Duration) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Duration(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Duration", str)
	}
	return nil
}

func (e Duration) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type HistoricalResultType string

const (
	HistoricalResultTypeBridgevolume HistoricalResultType = "BRIDGEVOLUME"
	HistoricalResultTypeTransactions HistoricalResultType = "TRANSACTIONS"
	HistoricalResultTypeAddresses    HistoricalResultType = "ADDRESSES"
)

var AllHistoricalResultType = []HistoricalResultType{
	HistoricalResultTypeBridgevolume,
	HistoricalResultTypeTransactions,
	HistoricalResultTypeAddresses,
}

func (e HistoricalResultType) IsValid() bool {
	switch e {
	case HistoricalResultTypeBridgevolume, HistoricalResultTypeTransactions, HistoricalResultTypeAddresses:
		return true
	}
	return false
}

func (e HistoricalResultType) String() string {
	return string(e)
}

func (e *HistoricalResultType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = HistoricalResultType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid HistoricalResultType", str)
	}
	return nil
}

func (e HistoricalResultType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type Platform string

const (
	PlatformAll        Platform = "ALL"
	PlatformSwap       Platform = "SWAP"
	PlatformBridge     Platform = "BRIDGE"
	PlatformMessageBus Platform = "MESSAGE_BUS"
)

var AllPlatform = []Platform{
	PlatformAll,
	PlatformSwap,
	PlatformBridge,
	PlatformMessageBus,
}

func (e Platform) IsValid() bool {
	switch e {
	case PlatformAll, PlatformSwap, PlatformBridge, PlatformMessageBus:
		return true
	}
	return false
}

func (e Platform) String() string {
	return string(e)
}

func (e *Platform) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Platform(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Platform", str)
	}
	return nil
}

func (e Platform) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type StatisticType string

const (
	StatisticTypeMeanVolumeUsd     StatisticType = "MEAN_VOLUME_USD"
	StatisticTypeMedianVolumeUsd   StatisticType = "MEDIAN_VOLUME_USD"
	StatisticTypeTotalVolumeUsd    StatisticType = "TOTAL_VOLUME_USD"
	StatisticTypeMeanFeeUsd        StatisticType = "MEAN_FEE_USD"
	StatisticTypeMedianFeeUsd      StatisticType = "MEDIAN_FEE_USD"
	StatisticTypeTotalFeeUsd       StatisticType = "TOTAL_FEE_USD"
	StatisticTypeCountTransactions StatisticType = "COUNT_TRANSACTIONS"
	StatisticTypeCountAddresses    StatisticType = "COUNT_ADDRESSES"
)

var AllStatisticType = []StatisticType{
	StatisticTypeMeanVolumeUsd,
	StatisticTypeMedianVolumeUsd,
	StatisticTypeTotalVolumeUsd,
	StatisticTypeMeanFeeUsd,
	StatisticTypeMedianFeeUsd,
	StatisticTypeTotalFeeUsd,
	StatisticTypeCountTransactions,
	StatisticTypeCountAddresses,
}

func (e StatisticType) IsValid() bool {
	switch e {
	case StatisticTypeMeanVolumeUsd, StatisticTypeMedianVolumeUsd, StatisticTypeTotalVolumeUsd, StatisticTypeMeanFeeUsd, StatisticTypeMedianFeeUsd, StatisticTypeTotalFeeUsd, StatisticTypeCountTransactions, StatisticTypeCountAddresses:
		return true
	}
	return false
}

func (e StatisticType) String() string {
	return string(e)
}

func (e *StatisticType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = StatisticType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid StatisticType", str)
	}
	return nil
}

func (e StatisticType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

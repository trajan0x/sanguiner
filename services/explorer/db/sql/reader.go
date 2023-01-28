package sql

import (
	"context"
	"fmt"
	"github.com/synapsecns/sanguine/services/explorer/graphql/server/graph/model"
)

/*╔══════════════════════════════════════════════════════════════════════╗*\
▏*║                        Generic Read Functions                        ║*▕
\*╚══════════════════════════════════════════════════════════════════════╝*/

// GetUint64 gets a uint64 from a given query.
func (s *Store) GetUint64(ctx context.Context, query string) (uint64, error) {
	var res int64

	dbTx := s.db.WithContext(ctx).Raw(query).Find(&res)
	if dbTx.Error != nil {
		return 0, fmt.Errorf("failed to read bridge event: %w", dbTx.Error)
	}

	return uint64(res), nil
}

// GetFloat64 gets a float64 from a given query.
func (s *Store) GetFloat64(ctx context.Context, query string) (float64, error) {
	var res float64
	dbTx := s.db.WithContext(ctx).Raw(query).Find(&res)
	if dbTx.Error != nil {
		return 0, fmt.Errorf("failed to read bridge event: %w", dbTx.Error)
	}

	return res, nil
}

// GetStringArray returns a string array for a given query.
func (s *Store) GetStringArray(ctx context.Context, query string) ([]string, error) {
	var res []string

	dbTx := s.db.WithContext(ctx).Raw(query + " SETTINGS readonly=1").Find(&res)
	if dbTx.Error != nil {
		return nil, fmt.Errorf("failed to read bridge event: %w", dbTx.Error)
	}

	return res, nil
}

// GetBridgeEvent returns a bridge event.
func (s *Store) GetBridgeEvent(ctx context.Context, query string) (*BridgeEvent, error) {
	var res BridgeEvent

	dbTx := s.db.WithContext(ctx).Raw(query).Find(&res)
	if dbTx.Error != nil {
		return nil, fmt.Errorf("failed to read bridge event: %w", dbTx.Error)
	}

	return &res, nil
}

// GetBridgeEvents returns bridge events.
func (s *Store) GetBridgeEvents(ctx context.Context, query string) ([]BridgeEvent, error) {
	var res []BridgeEvent
	dbTx := s.db.WithContext(ctx).Raw(query).Find(&res)
	if dbTx.Error != nil {
		return nil, fmt.Errorf("failed to read bridge event: %w", dbTx.Error)
	}

	return res, nil
}

// GetAllBridgeEvents returns bridge events.
func (s *Store) GetAllBridgeEvents(ctx context.Context, query string) ([]HybridBridgeEvent, error) {
	var res []HybridBridgeEvent
	dbTx := s.db.WithContext(ctx).Raw(query).Find(&res)
	if dbTx.Error != nil {
		return nil, fmt.Errorf("failed to read bridge event: %w", dbTx.Error)
	}
	return res, nil
}

// GetAllMessageBusEvents returns message bus events.
func (s *Store) GetAllMessageBusEvents(ctx context.Context, query string) ([]HybridMessageBusEvent, error) {
	var res []HybridMessageBusEvent
	dbTx := s.db.WithContext(ctx).Raw(query).Find(&res)
	if dbTx.Error != nil {
		return nil, fmt.Errorf("failed to read message bus event: %w", dbTx.Error)
	}
	return res, nil
}

// GetTxCounts returns Tx counts.
func (s *Store) GetTxCounts(ctx context.Context, query string) ([]*model.TransactionCountResult, error) {
	var res []*model.TransactionCountResult
	dbTx := s.db.WithContext(ctx).Raw(query + " SETTINGS readonly=1").Find(&res)

	if dbTx.Error != nil {
		return nil, fmt.Errorf("failed to read bridge event: %w", dbTx.Error)
	}

	return res, nil
}

// GetTokenCounts returns Tx counts.
func (s *Store) GetTokenCounts(ctx context.Context, query string) ([]*model.TokenCountResult, error) {
	var res []*model.TokenCountResult
	dbTx := s.db.WithContext(ctx).Raw(query + " SETTINGS readonly=1").Find(&res)
	if dbTx.Error != nil {
		return nil, fmt.Errorf("failed to read bridge event: %w", dbTx.Error)
	}

	return res, nil
}

// GetDateResults returns the dya by day data.
func (s *Store) GetDateResults(ctx context.Context, query string) ([]*model.DateResult, error) {
	var res []*model.DateResult
	dbTx := s.db.WithContext(ctx).Raw(query + " SETTINGS readonly=1").Scan(&res)

	if dbTx.Error != nil {
		return nil, fmt.Errorf("failed to get date results: %w", dbTx.Error)
	}

	return res, nil
}

// GetAddressRanking gets AddressRanking for a given query.
func (s *Store) GetAddressRanking(ctx context.Context, query string) ([]*model.AddressRanking, error) {
	var res []*model.AddressRanking

	dbTx := s.db.WithContext(ctx).Raw(query + " SETTINGS readonly=1").Scan(&res)
	if dbTx.Error != nil {
		return nil, fmt.Errorf("failed to read bridge event: %w", dbTx.Error)
	}
	if len(res) == 0 {
		return nil, nil
	}

	return res, nil
}

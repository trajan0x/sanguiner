// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	context "context"
	big "math/big"

	common "github.com/ethereum/go-ethereum/common"

	db "github.com/synapsecns/sanguine/ethergo/submitter/db"

	mock "github.com/stretchr/testify/mock"
)

// Service is an autogenerated mock type for the Service type
type Service struct {
	mock.Mock
}

// DBTransaction provides a mock function with given fields: ctx, f
func (_m *Service) DBTransaction(ctx context.Context, f db.TransactionFunc) error {
	ret := _m.Called(ctx, f)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, db.TransactionFunc) error); ok {
		r0 = rf(ctx, f)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetAllTXAttemptByStatus provides a mock function with given fields: ctx, fromAddress, chainID, matchStatuses
func (_m *Service) GetAllTXAttemptByStatus(ctx context.Context, fromAddress common.Address, chainID *big.Int, matchStatuses ...db.Status) ([]db.TX, error) {
	_va := make([]interface{}, len(matchStatuses))
	for _i := range matchStatuses {
		_va[_i] = matchStatuses[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, fromAddress, chainID)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 []db.TX
	if rf, ok := ret.Get(0).(func(context.Context, common.Address, *big.Int, ...db.Status) []db.TX); ok {
		r0 = rf(ctx, fromAddress, chainID, matchStatuses...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]db.TX)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, common.Address, *big.Int, ...db.Status) error); ok {
		r1 = rf(ctx, fromAddress, chainID, matchStatuses...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetNonceAttemptsByStatus provides a mock function with given fields: ctx, fromAddress, chainID, nonce, matchStatuses
func (_m *Service) GetNonceAttemptsByStatus(ctx context.Context, fromAddress common.Address, chainID *big.Int, nonce uint64, matchStatuses ...db.Status) ([]db.TX, error) {
	_va := make([]interface{}, len(matchStatuses))
	for _i := range matchStatuses {
		_va[_i] = matchStatuses[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, fromAddress, chainID, nonce)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 []db.TX
	if rf, ok := ret.Get(0).(func(context.Context, common.Address, *big.Int, uint64, ...db.Status) []db.TX); ok {
		r0 = rf(ctx, fromAddress, chainID, nonce, matchStatuses...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]db.TX)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, common.Address, *big.Int, uint64, ...db.Status) error); ok {
		r1 = rf(ctx, fromAddress, chainID, nonce, matchStatuses...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetNonceForChainID provides a mock function with given fields: ctx, fromAddress, chainID
func (_m *Service) GetNonceForChainID(ctx context.Context, fromAddress common.Address, chainID *big.Int) (uint64, error) {
	ret := _m.Called(ctx, fromAddress, chainID)

	var r0 uint64
	if rf, ok := ret.Get(0).(func(context.Context, common.Address, *big.Int) uint64); ok {
		r0 = rf(ctx, fromAddress, chainID)
	} else {
		r0 = ret.Get(0).(uint64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, common.Address, *big.Int) error); ok {
		r1 = rf(ctx, fromAddress, chainID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetNonceStatus provides a mock function with given fields: ctx, fromAddress, chainID, nonce
func (_m *Service) GetNonceStatus(ctx context.Context, fromAddress common.Address, chainID *big.Int, nonce uint64) (db.Status, error) {
	ret := _m.Called(ctx, fromAddress, chainID, nonce)

	var r0 db.Status
	if rf, ok := ret.Get(0).(func(context.Context, common.Address, *big.Int, uint64) db.Status); ok {
		r0 = rf(ctx, fromAddress, chainID, nonce)
	} else {
		r0 = ret.Get(0).(db.Status)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, common.Address, *big.Int, uint64) error); ok {
		r1 = rf(ctx, fromAddress, chainID, nonce)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetTXS provides a mock function with given fields: ctx, fromAddress, chainID, statuses
func (_m *Service) GetTXS(ctx context.Context, fromAddress common.Address, chainID *big.Int, statuses ...db.Status) ([]db.TX, error) {
	_va := make([]interface{}, len(statuses))
	for _i := range statuses {
		_va[_i] = statuses[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, fromAddress, chainID)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 []db.TX
	if rf, ok := ret.Get(0).(func(context.Context, common.Address, *big.Int, ...db.Status) []db.TX); ok {
		r0 = rf(ctx, fromAddress, chainID, statuses...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]db.TX)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, common.Address, *big.Int, ...db.Status) error); ok {
		r1 = rf(ctx, fromAddress, chainID, statuses...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MarkAllBeforeOrAtNonceReplacedOrConfirmed provides a mock function with given fields: ctx, signer, chainID, nonce
func (_m *Service) MarkAllBeforeOrAtNonceReplacedOrConfirmed(ctx context.Context, signer common.Address, chainID *big.Int, nonce uint64) error {
	ret := _m.Called(ctx, signer, chainID, nonce)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, common.Address, *big.Int, uint64) error); ok {
		r0 = rf(ctx, signer, chainID, nonce)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// PutTXS provides a mock function with given fields: ctx, txs
func (_m *Service) PutTXS(ctx context.Context, txs ...db.TX) error {
	_va := make([]interface{}, len(txs))
	for _i := range txs {
		_va[_i] = txs[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, ...db.TX) error); ok {
		r0 = rf(ctx, txs...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewService interface {
	mock.TestingT
	Cleanup(func())
}

// NewService creates a new instance of Service. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewService(t mockConstructorTestingTNewService) *Service {
	mock := &Service{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

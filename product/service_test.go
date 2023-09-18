package product_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/silvan-talos/shipping"
	"github.com/silvan-talos/shipping/mock"
	"github.com/silvan-talos/shipping/product"
)

func TestService_CalculatePacksConfiguration(t *testing.T) {
	tests := map[string]struct {
		qty         uint64
		packs       shipping.PackRepository
		expectedRes []shipping.PackConfig
		expectedErr error
	}{
		"example1": {
			qty:   1,
			packs: &mock.PackRepository{},
			expectedRes: []shipping.PackConfig{
				{
					Count: 1,
					Size:  250,
				},
			},
		},
		"example2": {
			qty:   250,
			packs: &mock.PackRepository{},
			expectedRes: []shipping.PackConfig{
				{
					Count: 1,
					Size:  250,
				},
			},
		},
		"example3": {
			qty:   251,
			packs: &mock.PackRepository{},
			expectedRes: []shipping.PackConfig{
				{
					Count: 1,
					Size:  500,
				},
			},
		},
		"example4": {
			qty:   501,
			packs: &mock.PackRepository{},
			expectedRes: []shipping.PackConfig{
				{
					Count: 1,
					Size:  500,
				},
				{
					Count: 1,
					Size:  250,
				},
			},
		},
		"example5": {
			qty:   12001,
			packs: &mock.PackRepository{},
			expectedRes: []shipping.PackConfig{
				{
					Count: 2,
					Size:  5000,
				},
				{
					Count: 1,
					Size:  2000,
				},
				{
					Count: 1,
					Size:  250,
				},
			},
		},
		"scenario1": {
			qty:   751,
			packs: &mock.PackRepository{},
			expectedRes: []shipping.PackConfig{
				{
					Count: 1,
					Size:  1000,
				},
			},
		},
		"scenario2": {
			qty: 251,
			packs: &mock.PackRepository{
				GetByProductIDFn: func(ctx context.Context, productID uint64) ([]uint64, error) {
					return []uint64{100, 250, 300, 500, 1000}, nil
				},
			},
			expectedRes: []shipping.PackConfig{
				{
					Count: 1,
					Size:  300,
				},
			},
		},
		"scenario3": {
			qty: 281,
			packs: &mock.PackRepository{
				GetByProductIDFn: func(ctx context.Context, productID uint64) ([]uint64, error) {
					return []uint64{210, 250, 260, 280, 500, 1000, 2000, 5000}, nil
				},
			},
			expectedRes: []shipping.PackConfig{
				{
					Count: 2,
					Size:  210,
				},
			},
		},
		"scenario4": {
			qty: 301,
			packs: &mock.PackRepository{
				GetByProductIDFn: func(ctx context.Context, productID uint64) ([]uint64, error) {
					return []uint64{200, 300, 500, 1000, 2000, 5000}, nil
				},
			},
			expectedRes: []shipping.PackConfig{
				{
					Count: 2,
					Size:  200,
				},
			},
		},
		"scenario5": {
			qty: 1251,
			packs: &mock.PackRepository{
				GetByProductIDFn: func(ctx context.Context, productID uint64) ([]uint64, error) {
					return []uint64{250, 600}, nil
				},
			},
			expectedRes: []shipping.PackConfig{
				{
					Count: 2,
					Size:  600,
				},
				{
					Count: 1,
					Size:  250,
				},
			},
		},
		"scenario6": {
			qty: 1499,
			packs: &mock.PackRepository{
				GetByProductIDFn: func(ctx context.Context, productID uint64) ([]uint64, error) {
					return []uint64{250, 600}, nil
				},
			},
			expectedRes: []shipping.PackConfig{
				{
					Count: 6,
					Size:  250,
				},
			},
		},
		"configurationNotFound_returnErrNotFound": {
			qty: 1,
			packs: &mock.PackRepository{
				GetByProductIDFn: func(ctx context.Context, productID uint64) ([]uint64, error) {
					return nil, shipping.ErrNotFound
				}},
			expectedErr: shipping.ErrNotFound,
		},
		"failedToGetConfiguration_internalError": {
			qty: 1,
			packs: &mock.PackRepository{
				GetByProductIDFn: func(ctx context.Context, productID uint64) ([]uint64, error) {
					return nil, errors.New("failed to get config")
				}},
			expectedErr: shipping.InternalServerErr,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			args := product.ServiceArgs{
				Packs: tc.packs,
			}
			s := product.NewService(args)
			res, err := s.CalculatePacksConfiguration(context.Background(), 1, tc.qty)
			require.Equal(t, tc.expectedErr, err, "errors must match")
			if err == nil {
				require.Equal(t, tc.expectedRes, res, "results must match")
			}
		})
	}
}

func TestService_UpdatePacksConfiguration(t *testing.T) {
	tests := map[string]struct {
		config      []uint64
		packs       shipping.PackRepository
		expectedErr error
	}{
		"configEmpty_invalidRequest": {
			config:      []uint64{},
			packs:       &mock.PackRepository{},
			expectedErr: product.ErrInvalidConfig,
		},
		"productIdNotFound_returnErrNotFound": {
			config: []uint64{100, 200},
			packs: &mock.PackRepository{
				UpdateConfigFn: func(ctx context.Context, productID uint64, config []uint64) error {
					return shipping.ErrNotFound
				},
			},
			expectedErr: shipping.ErrNotFound,
		},
		"failedToUpdateConfiguration_returnInternalError": {
			config: []uint64{100, 200},
			packs: &mock.PackRepository{
				UpdateConfigFn: func(ctx context.Context, productID uint64, config []uint64) error {
					return errors.New("failed to update config")
				},
			},
			expectedErr: shipping.InternalServerErr,
		},
		"updateConfig_successful": {
			config:      []uint64{250},
			packs:       &mock.PackRepository{},
			expectedErr: nil,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			args := product.ServiceArgs{
				Packs: tc.packs,
			}
			s := product.NewService(args)
			err := s.UpdatePacksConfiguration(context.Background(), 1, tc.config)
			require.Equal(t, tc.expectedErr, err)
		})
	}
}

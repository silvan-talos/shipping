package mock

import (
	"context"
)

type PackRepository struct {
	GetByProductIDFn func(ctx context.Context, productID uint64) ([]uint64, error)
}

func (pr *PackRepository) GetByProductID(ctx context.Context, productID uint64) ([]uint64, error) {
	if pr.GetByProductIDFn != nil {
		return pr.GetByProductIDFn(ctx, productID)
	}
	return []uint64{250, 500, 1000, 2000, 5000}, nil
}

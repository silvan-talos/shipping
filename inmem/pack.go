package inmem

import (
	"context"
	"sync"

	"github.com/silvan-talos/shipping"
)

func NewPackRepository() shipping.PackRepository {
	return &packRepository{
		configs: make(map[uint64][]uint64),
	}
}

type packRepository struct {
	mtx     sync.RWMutex
	configs map[uint64][]uint64
}

func (pr *packRepository) GetByProductID(_ context.Context, productID uint64) ([]uint64, error) {
	pr.mtx.RLock()
	defer pr.mtx.RUnlock()
	config, ok := pr.configs[productID]
	if !ok {
		// return default configuration if none exists
		return []uint64{250, 500, 1000, 2000, 5000}, nil
	}
	return config, nil
}


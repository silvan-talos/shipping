package product

import (
	"context"
	"errors"
	"log"
	"sort"

	"github.com/silvan-talos/shipping"
)

var (
	ErrInvalidConfig = errors.New("invalid config: config cannot be empty")
)

type Service interface {
	CalculatePacksConfiguration(ctx context.Context, id, qty uint64) ([]shipping.PackConfig, error)
	UpdatePacksConfiguration(ctx context.Context, id uint64, config []uint64) error
}

type service struct {
	packs shipping.PackRepository
}

func NewService(args ServiceArgs) Service {
	err := shipping.Validate.Struct(args)
	if err != nil {
		log.Fatal("failed to create product service, err:", err)
	}

	return &service{
		packs: args.Packs,
	}
}

type ServiceArgs struct {
	Packs shipping.PackRepository `validate:"required"`
}

func (s *service) CalculatePacksConfiguration(ctx context.Context, id, quantity uint64) ([]shipping.PackConfig, error) {
	packSizes, err := s.packs.GetByProductID(ctx, id)
	if err != nil {
		if errors.Is(err, shipping.ErrNotFound) {
			log.Println("no config found for product id:", id)
			return nil, err
		}
		log.Println("failed to get packs config, err:", err)
		return nil, shipping.InternalServerErr
	}
	conf, minOverhead := overheadAlgorithm(int64(quantity), packSizes)
	packConf, overhead := divisionAlgorithm(int64(quantity), packSizes)
	// choose better solution based on configuration accuracy
	if overhead > minOverhead {
		packConf = map[uint64]int64{conf.Size: conf.Count}
	}
	// sort by pack size desc
	keys := make([]uint64, 0, len(packConf))
	for k := range packConf {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] > keys[j]
	})
	// convert packs to meaningful struct
	packs := make([]shipping.PackConfig, 0, len(packConf))
	for _, k := range keys {
		if packConf[k] > 0 {
			packs = append(packs, shipping.PackConfig{
				Count: packConf[k],
				Size:  k,
			})
		}
	}
	return packs, nil
}

func (s *service) UpdatePacksConfiguration(ctx context.Context, id uint64, config []uint64) error {
	if len(config) == 0 {
		return ErrInvalidConfig
	}
	err := s.packs.UpdateConfig(ctx, id, config)
	if err != nil {
		if errors.Is(err, shipping.ErrNotFound) {
			log.Println("no product found for the specified ID, id:", id)
			return shipping.ErrNotFound
		}
		log.Println("error updating configuration, err:", err)
		return shipping.InternalServerErr
	}
	return nil
}

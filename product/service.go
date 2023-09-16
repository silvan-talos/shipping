package product

import (
	"context"
	"errors"
	"log"
	"sort"

	"github.com/silvan-talos/shipping"
)

type Service interface {
	CalculatePacksConfiguration(ctx context.Context, id, qty uint64) ([]shipping.PackConfig, error)
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
	qty := int64(quantity)
	packSizes, err := s.packs.GetByProductID(ctx, id)
	if err != nil {
		if errors.Is(err, shipping.ErrNotFound) {
			log.Println("no config found for product id:", id)
			return nil, err
		}
		log.Println("failed to get packs config, err:", err)
		return nil, shipping.InternalServerErr
	}
	// sort packSizes descending
	sort.Slice(packSizes, func(i, j int) bool {
		return packSizes[i] > packSizes[j]
	})
	log.Printf("pack sizes {%v}\n", packSizes)
	packConf := make(map[uint64]int64)
	for _, size := range packSizes {
		if qty >= int64(size) {
			packConf[size] += qty / int64(size)
			qty %= int64(size)
		}
	}
	// if there is leftover quantity, add one pack of the smallest size
	if qty > 0 {
		packConf[packSizes[len(packSizes)-1]]++
	}
	optimizePacks(packSizes, packConf)
	// convert packs to meaningful struct
	packs := make([]shipping.PackConfig, 0, len(packConf))
	for size, count := range packConf {
		if count > 0 {
			packs = append(packs, shipping.PackConfig{
				Count: count,
				Size:  size,
			})
		}
	}
	return packs, nil
}

// optimizePacks creates an optimal amount of packages by merging smaller packages into bigger ones if possible
func optimizePacks(packSizes []uint64, packs map[uint64]int64) {
	for i := len(packSizes) - 1; i > 0; i-- {
		totalSmallerAmount := sumQuantity(packSizes[i:], packs)
		if totalSmallerAmount >= int64(packSizes[i]) {
			packs[packSizes[i]] += totalSmallerAmount / int64(packSizes[i])
			zeroSubsequent(totalSmallerAmount, packSizes[i+1:], packs)
		}
	}
}

// sumQuantity calculates the sum of size*number of packs for the provided packSizes
func sumQuantity(packSizes []uint64, packs map[uint64]int64) int64 {
	var sum int64 = 0
	for i := len(packSizes) - 1; i > 0; i-- {
		sum += packs[packSizes[i]] * int64(packSizes[i])
	}
	return sum
}

// zeroSubsequent updates values of smaller pack sizes after a merge
func zeroSubsequent(qty int64, packSizes []uint64, packs map[uint64]int64) {
	for _, size := range packSizes {
		if packs[size]*int64(size) <= qty {
			qty -= packs[size] * int64(size)
			packs[size] = 0
		}
	}
}

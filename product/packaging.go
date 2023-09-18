package product

import (
	"math"
	"sort"

	"github.com/silvan-talos/shipping"
)

// overheadAlgorithm returns a configuration based on min items to send in min pack count
func overheadAlgorithm(qty int64, packSizes []uint64) (shipping.PackConfig, int64) {
	// sort packSizes asc
	sort.Slice(packSizes, func(i, j int) bool {
		return packSizes[i] < packSizes[j]
	})
	overheads := make(map[uint64]int64)
	packQuantities := make(map[uint64]int64)
	for _, packSize := range packSizes {
		packQuantities[packSize] = int64(math.Ceil(float64(qty) / float64(packSize)))
	}
	var minOh int64 = math.MaxInt64
	for size, amt := range packQuantities {
		overheads[size] = (int64(size) * amt) - qty
		if overheads[size] < minOh {
			minOh = overheads[size]
		}
	}
	var minPackSize int64 = math.MaxInt64
	sameOhPacks := make(map[uint64]int64)
	for size, oh := range overheads {
		if oh == minOh {
			sameOhPacks[size] = packQuantities[size]
			if packQuantities[size] < minPackSize {
				minPackSize = packQuantities[size]
			}
		}
	}
	var res shipping.PackConfig
	for size, count := range sameOhPacks {
		if count == minPackSize {
			res.Size = size
			res.Count = count
			break
		}
	}
	return res, minOh
}

// divisionAlgorithm creates a configuration based on bigger size first
func divisionAlgorithm(qty int64, packSizes []uint64) (map[uint64]int64, int64) {
	initialQty := qty
	// sort sizes descending
	sort.Slice(packSizes, func(i, j int) bool {
		return packSizes[i] > packSizes[j]
	})
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
	// calculate overhead
	var s int64 = 0
	for size, count := range packConf {
		s += int64(size) * count
	}
	overhead := s - initialQty
	return packConf, overhead
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

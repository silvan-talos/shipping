package product

import (
	"math"
	"sort"
)

// overheadAlgorithm returns a configuration based on min items to send in min pack count
func overheadAlgorithm(qty int64, packSizes []uint64) (map[uint64]int64, int64) {
	// sort packSizes asc
	sort.Slice(packSizes, func(i, j int) bool {
		return packSizes[i] < packSizes[j]
	})
	overheads := make(map[uint64]int64)
	packQuantities := make(map[uint64]map[uint64]int64)
	for i, packSize := range packSizes {
		packQuantities[packSize] = make(map[uint64]int64)
		packQuantities[packSize][packSize] = qty / int64(packSize)
		remainingQty := qty % int64(packSize)
		for j := 0; j <= i && remainingQty != 0; j++ {
			if remainingQty <= int64(packSizes[j]) {
				packQuantities[packSize][packSizes[j]]++
				break
			}
		}
	}
	var minOh int64 = math.MaxInt64
	for size, cfg := range packQuantities {
		overheads[size] = calculateOverhead(qty, cfg)
		if overheads[size] < minOh {
			minOh = overheads[size]
		}
	}
	var minPackSize int64 = math.MaxInt64
	sameOhPacks := make(map[uint64]map[uint64]int64)
	for size, oh := range overheads {
		if oh == minOh {
			sameOhPacks[size] = packQuantities[size]
			if packSize := calculateSize(packQuantities[size]); packSize < minPackSize {
				minPackSize = packSize
			}
		}
	}
	var res = make(map[uint64]int64)
	for _, cfg := range sameOhPacks {
		if calculateSize(cfg) == minPackSize {
			res = cfg
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
	return packConf, calculateOverhead(initialQty, packConf)
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

func calculateOverhead(qty int64, config map[uint64]int64) int64 {
	var s int64 = 0
	for size, count := range config {
		s += int64(size) * count
	}
	overhead := s - qty
	return overhead
}

func calculateSize(conf map[uint64]int64) int64 {
	var s int64 = 0
	for _, count := range conf {
		s += count
	}
	return s
}

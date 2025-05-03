package models

import (
	"crypto/rand"
	"fmt"
	"homework/pkg/errors"
	"io"
	"math/big"
	"sort"
)

// generateUniqueRandomNumbers generates 'count' unique random numbers
// within the range [min, max] using crypto/rand for security.
func generateUniqueRandomNumbers(randReader io.Reader, count, min, max int) ([]int, error) {
	if count > (max - min + 1) {
		return nil, errors.New("cannot generate more unique numbers than available in range")
	}
	if count <= 0 || min > max {
		return nil, errors.New("invalid count or range")
	}

	numbers := make(map[int]struct{})
	result := make([]int, 0, count)
	attempts := 0
	maxAttempts := count * 10 // Safety break

	for len(result) < count && attempts < maxAttempts {
		attempts++
		nBig, err := rand.Int(randReader, big.NewInt(int64(max-min+1)))
		if err != nil {
			return nil, fmt.Errorf("failed to generate random number: %w", err)
		}
		num := int(nBig.Int64()) + min

		if _, exists := numbers[num]; !exists {
			numbers[num] = struct{}{}
			result = append(result, num)
		}
	}

	if len(result) < count {
		return nil, errors.New("failed to generate sufficient unique numbers within attempts limit")
	}

	sort.Ints(result) // Keep winning numbers sorted
	return result, nil
}

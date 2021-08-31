package utils

import (
	"errors"
	"math"
)

// ReverseMap is the method for reversing map (exchange corresponding keys and values).
func ReverseMap(m map[string]string) map[string]string {
	if len(m) == 0 {
		return m
	}

	reversed := make(map[string]string, len(m))

	for key, val := range m {
		reversed[val] = key
	}

	return reversed
}

// FilterSlice is the method for filtering slice.
func FilterSlice(slice []string, unwanted []string) []string {
	if len(slice) == 0 || len(unwanted) == 0 {
		return slice
	}

	filtered := make([]string, 0)

	unwantedSet := make(map[string]struct{})
	for _, el := range unwanted {
		unwantedSet[el] = struct{}{}
	}

	for _, el := range slice {
		if _, ok := unwantedSet[el]; !ok {
			filtered = append(filtered, el)
		}
	}

	return filtered
}

// SplitToBatches is the method for splitting slice of string to batches.
func SplitToBatches(elements []string, batchSize uint) ([][]string, error) {
	if len(elements) == 0 {
		return nil, errors.New("slice must not be empty")
	}

	if batchSize == 0 {
		return nil, errors.New("batch size cannot be equal to 0")
	}

	if int(batchSize) >= len(elements) {
		return [][]string{elements}, nil
	}

	batches := make([][]string, int(math.Ceil(float64(len(elements))/float64(batchSize))))

	for i := 0; i < cap(batches); i++ {
		if start, end := i*int(batchSize), (i+1)*int(batchSize); end < len(elements) {
			batches[i] = elements[start:end]
		} else {
			batches[i] = elements[start:]
		}
	}

	return batches, nil
}

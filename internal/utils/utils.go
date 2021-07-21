package utils

import (
	"errors"
	"math"
)

func ReverseMap(m map[string]string) map[string]string {
	if m == nil || len(m) == 0 {
		return m
	}

	reversed := make(map[string]string, len(m))

	for key, val := range m {
		reversed[val] = key
	}

	return reversed
}

func FilterSlice(slice []string, unwanted []string) []string {
	if len(slice) == 0 || len(unwanted) == 0 {
		return slice
	}

	filtered := make([]string, 0)

	blackList := make(map[string]struct{})
	for _, el := range unwanted {
		blackList[el] = struct{}{}
	}

	for _, el := range slice {
		if _, ok := blackList[el]; !ok {
			filtered = append(filtered, el)
		}
	}

	return filtered
}

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

	batches := make([][]string, 0,  int(math.Ceil(float64(len(elements)) / float64(batchSize))))

	for i := 0; i < cap(batches); i++ {
		var batch []string
		if start, end := i*int(batchSize), (i+1)*int(batchSize); end < len(elements) {
			batch = elements[start:end]
		} else {
			batch = elements[start:]
		}
		batches = append(batches, batch)
	}

	return batches, nil
}
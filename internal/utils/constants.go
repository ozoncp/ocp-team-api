package utils

// SearchType is the type of Full Text Search (FTS): plaintext-oriented or phrase-oriented
type SearchType uint8

const (
	Plain  SearchType = 0
	Phrase SearchType = 1
)

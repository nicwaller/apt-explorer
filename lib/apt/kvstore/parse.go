package kvstore

import (
	"io"
)

type Block struct {
	SingleValues map[string]string
	MultiValues  map[string][]string
}

func ParseImmediate(input io.Reader) ([]Block, error) {
	blocks := make([]Block, 0)
	err := Parse(input, func(block Block) {
		blocks = append(blocks, block)
	})
	if err != nil {
		return nil, err
	}
	return blocks, nil
}

// FIXME: PERF: use some pointers to be faster
// I don't like this callback interface very much. Maybe wrap it with an iterator?
func Parse(input io.Reader, gotBlock func(Block)) error {
	var lastKey string
	// TODO: store a compact numeric encoding of the key, not the whole key string
	singleMap := make(map[string]string)
	arrayMap := make(map[string][]string)
	arr := make([]string, 0)
	var lastError error
	send := func(b Block) {
		gotBlock(b)
		singleMap = make(map[string]string)
		arrayMap = make(map[string][]string)
		arr = make([]string, 0)
	}
	Tokenize(input, func(kind TokenType, tok []byte) {
		switch kind {
		case tSingleKey, tArrayKey:
			if len(arr) > 0 {
				arrayMap[lastKey] = arr
				arr = make([]string, 0)
			}
			lastKey = string(tok)
		case tSingleVal:
			singleMap[lastKey] = string(tok)
		case tArrayVal:
			arr = append(arr, string(tok))
		case tEndOfBlock, tEOF:
			if len(arr) > 0 {
				arrayMap[lastKey] = arr
				arr = make([]string, 0)
			}
			if len(singleMap) > 0 || len(arrayMap) > 0 {
				send(Block{
					SingleValues: singleMap,
					MultiValues:  arrayMap,
				})
			}
		}
	}, func(e error) {
		lastError = e
	})
	return lastError
}

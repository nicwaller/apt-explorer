package kvstore

import (
	"errors"
	"io"
)

type DfaState uint8

const (
	sBegin              DfaState = iota
	sScanningKey        DfaState = iota
	sChoice             DfaState = iota
	sExpectingValue     DfaState = iota
	sScanningValue      DfaState = iota
	sExpectingArray     DfaState = iota
	sScanningArrayValue DfaState = iota
	sIsItOver           DfaState = iota
)

type TokenType uint8

const (
	tSingleKey  TokenType = iota
	tSingleVal  TokenType = iota
	tArrayKey   TokenType = iota
	tArrayVal   TokenType = iota
	tEndOfBlock TokenType = iota
	tEOF        TokenType = iota
)

func Tokenize(input io.Reader, gotToken func(kind TokenType, tok []byte), gotError func(e error)) {
	token := make([]byte, 4096) // some of those dependency chains can be quite large
	tokI := 0
	buf := make([]byte, 4096)
	state := sBegin

	// TODO: maybe report row/col where error occurred, include some context?
	for {
		n, err := input.Read(buf)
		if err == io.EOF {
			gotToken(tEOF, []byte{})
			// NOTE: if the input doesn't properly end with \n then the last record will be lost
			break
		}
		if err != nil {
			gotError(err)
			return
		}
		var c byte
		for i := 0; i < n; i++ {
			c = buf[i]
			switch state {
			case sBegin:
				switch true {
				case c >= 'A' && c <= 'Z':
					state = sScanningKey
					token[tokI] = c
					tokI++
				default:
					gotError(errors.New("key must start with uppercase letter"))
					return
				}
			case sScanningKey:
				switch c {
				case ':':
					state = sChoice
				case ' ', '\n':
					gotError(errors.New("key must be immediately followed by a colon \":\""))
					return
				default:
					token[tokI] = c
					tokI++
					continue
				}
			case sChoice:
				// well it turns out that some files put a trailing space after List:
				// and that sucks. -NW
				switch c {
				case ' ':
					continue
				case '\n':
					gotToken(tArrayKey, token[:tokI])
					tokI = 0
					state = sExpectingArray
				default:
					gotToken(tSingleKey, token[:tokI])
					token[0] = c
					tokI = 1
					state = sScanningValue
				}
			case sExpectingValue:
				// Yes, ":" is totally permitted in the value field (eg. 1:2.0.21)
				switch c {
				case ' ':
					continue
				case '\n':
					gotError(errors.New("expected exactly 1 space between key and value"))
					return
				default:
					state = sScanningValue
					token[tokI] = c
					tokI++
					continue
				}
			case sScanningValue:
				// Q: is it permissible for values to contain a literal ":"?
				switch c {
				case '\n':
					state = sIsItOver
					gotToken(tSingleVal, token[:tokI])
					tokI = 0
				default:
					token[tokI] = c
					tokI++
					continue
				}
			case sExpectingArray:
				switch c {
				case ' ':
					state = sScanningArrayValue
				case '\n':
					gotError(errors.New("empty arrays are not permitted"))
					return
				default:
					gotError(errors.New("lines for array values must start with a single space"))
					return
				}
			case sScanningArrayValue:
				switch c {
				case '\n':
					state = sIsItOver
					gotToken(tArrayVal, token[:tokI])
					tokI = 0
				default:
					token[tokI] = c
					tokI++
					continue
				}
			case sIsItOver:
				switch true {
				case c == ' ':
					state = sScanningArrayValue
				case c == '\n':
					gotToken(tEndOfBlock, []byte{})
					state = sBegin
				case c >= 'A' && c <= 'Z':
					state = sScanningKey
					token[tokI] = c
					tokI++
				default:
					gotError(errors.New("key must start with uppercase letter"))
					return
				}
			}
		}
	}
}

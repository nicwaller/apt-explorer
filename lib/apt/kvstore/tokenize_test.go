package kvstore

import (
	"fmt"
	"strings"
	"testing"
)

type Token struct {
	kind  TokenType
	value string
}

func Tokens(input string) ([]Token, error) {
	toks := make([]Token, 0)
	sr := strings.NewReader(input)
	var eResult error
	Tokenize(sr, func(kind TokenType, tok []byte) {
		if kind == tEOF {
			return
		}
		toks = append(toks, Token{
			kind:  kind,
			value: string(tok),
		})
	}, func(e error) {
		eResult = e
	})
	return toks, eResult
}

func TestOkEmptyString(t *testing.T) {
	toks, e := Tokens("")
	if e != nil {
		t.Errorf("Unexpected Error: %v\n", e)
	}
	if len(toks) != 0 {
		fmt.Println(toks)
		t.Error("Got a record but expected none")
	}
}

func TestOkWhitespace(t *testing.T) {
	toks, e := Tokens("      ")
	if e == nil {
		t.Error("Expected error but got none")
	}
	if e.Error() != "key must start with uppercase letter" {
		t.Error("Wrong error message text")
	}
	if len(toks) != 0 {
		t.Error("Got a record but expected none")
	}
}

func TestRejectEmptyLines(t *testing.T) {
	toks, e := Tokens("\n \n\n \n\n \n")
	if e == nil {
		t.Error("Expected error but got none")
	}
	if e.Error() != "key must start with uppercase letter" {
		t.Error("Wrong error message text")
	}
	if len(toks) != 0 {
		t.Error("Expected zero tokens")
	}
}

func TestRejectInvalidLine(t *testing.T) {
	toks, e := Tokens("ThisLineNotValid")
	if e != nil {
		t.Errorf("Unexpected Error: %v\n", e)
	}
	if len(toks) != 0 {
		t.Error("Got a record but expected none")
	}
}

func TestOkSingleLine(t *testing.T) {
	toks, e := Tokens("Key: Value\n")
	if e != nil {
		t.Errorf("Unexpected Error: %v\n", e)
	}
	if len(toks) != 2 {
		t.Error("Expected exactly 2 tokens")
	}
	if toks[0].kind != tSingleKey {
		t.Error("Wrong token type")
	}
	if toks[0].value != "Key" {
		t.Error("Wrong token value")
	}
	if toks[1].kind != tSingleVal {
		t.Error("Wrong token type")
	}
	if toks[1].value != "Value" {
		t.Error("Wrong token value")
	}
}

func TestOkArray(t *testing.T) {
	toks, e := Tokens("List: \n Apple\n Orange\n")
	if e != nil {
		t.Errorf("Unexpected Error: %v\n", e)
	}
	if len(toks) != 3 {
		t.Error("Expected exactly 3 tokens")
	}
	if toks[0].kind != tArrayKey {
		t.Error("Expected array key")
	}
	if toks[0].value != "List" {
		t.Error("Wrong token value")
	}
	if toks[1].kind != tArrayVal {
		t.Error("Expected array value")
	}
	if toks[1].value != "Apple" {
		t.Error("Wrong token value")
	}
	if toks[2].kind != tArrayVal {
		t.Error("Expected array value")
	}
	if toks[2].value != "Orange" {
		t.Error("Wrong token value")
	}
}

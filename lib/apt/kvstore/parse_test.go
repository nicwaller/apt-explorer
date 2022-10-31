package kvstore

import (
	"fmt"
	"strings"
	"testing"
)

func Records(input string) []Block {
	sr := strings.NewReader(input)
	recs := make([]Block, 0)
	Parse(sr, func(block Block) {
		recs = append(recs, block)
	})
	return recs
}

func TestOkSingle(t *testing.T) {
	r := Records("Key: Value\n")
	if len(r) != 1 {
		t.Error("Expected 1 record")
		return
	}
	if len(r[0].SingleValues) != 1 {
		t.Error("Expected 1 key")
		return
	}
	for k, v := range r[0].SingleValues {
		if k != "Key" {
			t.Error("Expected key=Key")
			return
		}
		if v != "Value" {
			t.Error("Expected value=Value")
			return
		}
	}
}

func TestOkList(t *testing.T) {
	r := Records("List: \n Apple\n Orange\n")
	if len(r) != 1 {
		t.Error("Expected 1 record")
		return
	}
	if len(r[0].MultiValues) != 1 {
		t.Error("Expected 1 key")
		return
	}
	for k, v := range r[0].MultiValues {
		if k != "List" {
			t.Error("Expected key=List")
			return
		}
		if v[0] != "Apple" {
			fmt.Printf(">>>%v<<<", v[0])
			t.Error("Expected value=Apple")
			return
		}
		if v[1] != "Orange" {
			t.Error("Expected value=Orange")
			return
		}
	}
}

func TestRejectEmptyList(t *testing.T) {
	r := Records("List: \n\n")
	if len(r) != 0 {
		t.Error("Expected 0 records")
		return
	}
}

func TestTwoRecords(t *testing.T) {
	records := Records("Key: Val\n\nKey: Val\n\n")
	if len(records) != 2 {
		fmt.Println(records)
		t.Errorf("Expected 2 records, got %d", len(records))
		return
	}
	for _, r := range records {
		if len(r.SingleValues) != 1 {
			t.Error("Expected 1 key in record")
			return
		}
		for k, v := range r.SingleValues {
			if k != "Key" {
				t.Error("Wrong key")
				return
			}
			if v != "Val" {
				t.Error("Wrong value")
				return
			}
		}
	}
}

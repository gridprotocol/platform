package utils

import (
	"fmt"
	"testing"
)

func TestDur2TS(t *testing.T) {
	month := "1"
	ts, err := DurToTS(month)
	if err != nil {
		t.Error("dur to ts error:", err)
	}
	fmt.Println("ts:", ts)
}

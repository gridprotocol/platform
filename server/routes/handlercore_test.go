package routes

import (
	"testing"

	"github.com/rockiecn/platform/lib/utils"
)

func Test_CalcValue(t *testing.T) {
	o := OrderInfo{
		NumCPU:   utils.Uint64ToString(2),
		PriCPU:   utils.Uint64ToString(1),
		NumGPU:   utils.Uint64ToString(1),
		PriGPU:   utils.Uint64ToString(10),
		NumStore: utils.Uint64ToString(100 * 1024 * 1024 * 1024),
		PriStore: utils.Uint64ToString(1),
		NumMem:   utils.Uint64ToString(1 * 1024 * 1024 * 1024),
		PriMem:   utils.Uint64ToString(10),
		Dur:      utils.Uint64ToString(1 * 30 * 86400),
	}

	res, err := CalcValue(&o)
	if err != nil {
		t.Log("calca value failed:", err)
		t.FailNow()
	}

	t.Log("value:", res)
}

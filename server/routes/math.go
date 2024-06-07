package routes

import (
	"regexp"

	"github.com/rockiecn/platform/lib/utils"
)

// calc credit cost of an order
func CalcCost(o *OrderInfo) (int64, error) {
	nCPU, err := utils.StringToInt64(o.NumCPU)
	if err != nil {
		return 0, err
	}
	pCPU, err := utils.StringToInt64(o.PriCPU)
	if err != nil {
		return 0, err
	}

	nGPU, err := utils.StringToInt64(o.NumGPU)
	if err != nil {
		return 0, err
	}
	pGPU, err := utils.StringToInt64(o.PriGPU)
	if err != nil {
		return 0, err
	}

	// unit: G
	nMem, err := utils.StringToInt64(o.NumMem)
	if err != nil {
		return 0, err
	}

	pMem, err := utils.StringToInt64(o.PriMem)
	if err != nil {
		return 0, err
	}

	// unit: T
	nStor, err := utils.StringToInt64(o.NumDisk)
	if err != nil {
		return 0, err
	}

	pStor, err := utils.StringToInt64(o.PriDisk)
	if err != nil {
		return 0, err
	}

	// dur: month
	dur, err := utils.StringToInt64(o.Dur)
	if err != nil {
		return 0, err
	}

	// month to min
	min := dur * 30 * 24 * 60

	// calc cost
	cost := (nCPU*pCPU + nGPU*pGPU + nMem*pMem + nStor*pStor) * min

	logger.Debug("credit cost of order:", cost)

	// return credit cost
	return cost, nil
}

// check number
func isNumber(s string) bool {
	pattern := `^[0-9]+(\.[0-9]+)?$`
	matched, _ := regexp.MatchString(pattern, s)
	return matched
}

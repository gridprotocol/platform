package routes

import (
	"regexp"

	"github.com/rockiecn/platform/lib/utils"
)

// calc credit cost of an order
func CalcCost(o *OrderInfo) (uint64, error) {
	nCPU, err := utils.StringToUint64(o.NumCPU)
	if err != nil {
		return 0, err
	}
	pCPU, err := utils.StringToUint64(o.PriCPU)
	if err != nil {
		return 0, err
	}

	nGPU, err := utils.StringToUint64(o.NumGPU)
	if err != nil {
		return 0, err
	}
	pGPU, err := utils.StringToUint64(o.PriGPU)
	if err != nil {
		return 0, err
	}

	// T
	nMem, err := utils.StringToUint64(o.NumMem)
	if err != nil {
		return 0, err
	}
	// T to byte
	nMem = nMem * 1024 * 1024 * 1024 * 1024
	pMem, err := utils.StringToUint64(o.PriMem)
	if err != nil {
		return 0, err
	}

	// G
	nStor, err := utils.StringToUint64(o.NumStore)
	if err != nil {
		return 0, err
	}
	// G to byte
	nStor = nStor * 1024 * 1024 * 1024
	pStor, err := utils.StringToUint64(o.PriStore)
	if err != nil {
		return 0, err
	}

	dur, err := utils.StringToUint64(o.Dur)
	if err != nil {
		return 0, err
	}

	// get wei value
	value := (nCPU*pCPU + nGPU*pGPU + nMem*pMem + nStor*pStor) * dur
	logger.Debug("wei of order:", value)

	cost := value / 1000 / 1000 / 1000 / 1000
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

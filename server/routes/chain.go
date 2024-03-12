package routes

// todo: add tx confirming logic
// check if a tx confirmed on chain
func checkTxConfirmed(txHash string) (bool, error) {
	_ = txHash
	return true, nil
}

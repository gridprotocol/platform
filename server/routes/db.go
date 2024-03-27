package routes

import (
	"encoding/json"
	"fmt"

	"github.com/dgraph-io/badger/v2"
	"github.com/rockiecn/platform/lib/utils"
)

// get order id key
func OrderIDKey(userAddr string) []byte {
	key := fmt.Sprintf("ORDER_ID_%s", userAddr)
	return []byte(key)
}

// payinfo id key
func PayInfoIDKey(userAddr string) []byte {
	key := fmt.Sprintf("PI_ID_%s", userAddr)
	return []byte(key)
}

// key for payinfo, with user and id
// PI_user_id
func PayInfoKey(userAddr string, id string) []byte {
	key := fmt.Sprintf("PI_%s_%s", userAddr, id)
	return []byte(key)
}

func TransferIDKey(userAddr string) []byte {
	key := fmt.Sprintf("TRANS_ID_%s", userAddr)
	return []byte(key)
}

// make order key for user with addr and id
// order key: ORDER_user_id
func OrderKey(userAddr string, id string) []byte {
	key := fmt.Sprintf("ORDER_%s_%s", userAddr, id)
	return []byte(key)
}

// make order list key for cp
// order list key: CP_ORDERS_*
func OrderListKey(cpAddr string) []byte {
	key := fmt.Sprintf("CP_ORDERS_%s", cpAddr)
	return []byte(key)
}

// key for cp info
// CP_INFO_*
func CPInfoKey(cpAddr string) []byte {
	key := fmt.Sprintf("CP_INFO_%s", cpAddr)
	return []byte(key)
}

// key for transferinfo, with user and id
// TI_user_id
func TransferInfoKey(userAddr string, id string) []byte {
	key := fmt.Sprintf("TI_%s_%s", userAddr, id)
	return []byte(key)
}

// credit key
// CREDIT_*
func CreditKey(addr string) []byte {
	key := fmt.Sprintf("CREDIT_%s", addr)
	return []byte(key)
}

// order id key: user_*
func (hc *HandlerCore) getOrderID(addr string) (string, error) {
	key := OrderIDKey(addr)

	var id string
	data, err := hc.LocalDB.Get([]byte(key))
	if err != nil {
		// if no id record, init with 0
		if err.Error() == "Key not found" {
			id = "0"
		} else {
			return "", err
		}
	} else {
		logger.Debugf("data:%s", data)
		id = string(data)
	}

	return id, nil
}

// payinfo id key: user_*
func (hc *HandlerCore) getPayInfoID(addr string) (string, error) {
	key := PayInfoIDKey(addr)

	var id string
	data, err := hc.LocalDB.Get([]byte(key))
	if err != nil {
		// if no id record, init with 0
		if err.Error() == "Key not found" {
			id = "0"
		} else {
			return "", err
		}
	} else {
		logger.Debugf("data:%s", data)
		id = string(data)
	}

	return id, nil
}

// transfer id key: user_*
func (hc *HandlerCore) getTransferID(addr string) (string, error) {
	key := TransferIDKey(addr)

	var id string
	data, err := hc.LocalDB.Get([]byte(key))
	if err != nil {
		// if no id record, init with 0
		if err.Error() == "Key not found" {
			id = "0"
		} else {
			return "", err
		}
	} else {
		logger.Debugf("data:%s", data)
		id = string(data)
	}

	return id, nil
}

// query credit for addr
func (hc *HandlerCore) queryCredit(addr string) (string, error) {
	// get key
	creKey := CreditKey(addr)

	var credit string

	// get credit
	data, err := hc.LocalDB.Get([]byte(creKey))
	if err != nil {
		if err.Error() == "Key not found" {
			credit = "0"
			return credit, nil
		} else {
			return "", err
		}
	}
	credit = string(data)

	return credit, nil
}

// get user's all transfer info from db
func (hc *HandlerCore) getUserTransfers(userAddr string) ([]TransferInfo, error) {
	transList := make([]TransferInfo, 0, 100)

	// get transfer id
	transID, err := hc.getTransferID(userAddr)
	if err != nil {
		return nil, err
	}

	logger.Debug("user's transfer id:", transID)

	// number of transfers
	num, err := utils.StringToInt64(transID)
	if err != nil {
		return nil, err
	}
	for i := int64(0); i < num; i++ {
		// make key: trans_user_id
		key := TransferInfoKey(userAddr, utils.Int64ToString(i))
		logger.Debug("transfer key:", key)
		// get transfer
		data, err := hc.LocalDB.Get([]byte(key))
		if err != nil {
			return nil, err
		}
		transfer := TransferInfo{}
		err = json.Unmarshal(data, &transfer)
		if err != nil {
			return nil, err
		}
		transList = append(transList, transfer)
	}

	return transList, nil
}

// set an order's settled state with key
func (hc *HandlerCore) setOrderSettled(key []byte, settled bool) (k []byte, v []byte, err error) {
	// get order info with key
	data, err := hc.LocalDB.Get([]byte(key))
	if err != nil {
		return nil, nil, err
	}

	order := OrderInfo{}
	err = json.Unmarshal(data, &order)
	if err != nil {
		return nil, nil, err
	}

	// set order's settled state
	order.Settled = settled

	// marshal new order
	data, err = json.Marshal(order)
	if err != nil {
		return nil, nil, err
	}

	// update order info
	// err = hc.LocalDB.Put(key, data)
	// if err != nil {
	// 	return nil, nil, err
	// }

	// return k,v for db
	return key, data, nil
}

// set transfer confirmed state
func (hc *HandlerCore) setTransferConfirmed(key []byte, confirmed bool) (k []byte, v []byte, err error) {
	// get transfer info with key
	data, err := hc.LocalDB.Get([]byte(key))
	if err != nil {
		return nil, nil, err
	}

	ti := TransferInfo{}
	err = json.Unmarshal(data, &ti)
	if err != nil {
		return nil, nil, err
	}

	// set order's settled state
	ti.TxConfirmed = confirmed

	// marshal new order
	data, err = json.Marshal(ti)
	if err != nil {
		return nil, nil, err
	}

	// return k,v for db
	return key, data, nil
}

// get user's pay info list from db
func (hc *HandlerCore) getPayInfoList(addr string) ([]PayInfo, error) {
	ciList := make([]PayInfo, 0, 100)

	// get payinfo id
	payID, err := hc.getPayInfoID(addr)
	if err != nil {
		return nil, err
	}
	logger.Debug("account's pay id:", payID)

	// number of order
	num, err := utils.StringToInt64(payID)
	if err != nil {
		return nil, err
	}
	for i := int64(0); i < num; i++ {
		// make payInfo key
		piKey := PayInfoKey(addr, utils.Int64ToString(i))
		logger.Debugf("order key:%s", piKey)
		// get payinfo
		data, err := hc.LocalDB.Get([]byte(piKey))
		if err != nil {
			return nil, err
		}
		ci := &PayInfo{}
		err = json.Unmarshal(data, ci)
		if err != nil {
			return nil, err
		}
		ciList = append(ciList, *ci)
	}

	// payinfo list
	return ciList, nil
}

// get order list for cp
func (hc *HandlerCore) getCpOrders(cpAddr string) ([]OrderInfo, error) {
	// 'cp' _ 'address' as cp key
	orderListKey := OrderListKey(cpAddr)

	// init an empty order list
	orderList := make([]OrderInfo, 0, 100)

	// read db for cp order keys data
	data, err := hc.LocalDB.Get([]byte(orderListKey))
	if err != nil {
		// if no order id, return empty order list
		if err.Error() == "Key not found" {
			return orderList, nil
		} else {
			return nil, err
		}
	}

	var orderKeys []string
	// unmarshal data into order keys if data is not empty
	if len(data) != 0 {
		err = json.Unmarshal([]byte(data), &orderKeys)
		if err != nil {
			return nil, err
		}
	} else { // if no key data, return empty list
		return orderList, nil
	}

	// get order list with order keys
	for i := 0; i < len(orderKeys); i++ {
		// each item is an order key
		key := orderKeys[i]
		// get order
		data, err := hc.LocalDB.Get([]byte(key))
		if err != nil {
			return nil, err
		}
		order := OrderInfo{}
		err = json.Unmarshal(data, &order)
		if err != nil {
			return nil, err
		}
		orderList = append(orderList, order)
	}

	return orderList, nil
}

// get user's order list from db
func (hc *HandlerCore) getUserOrders(userAddr string) ([]OrderInfo, error) {
	orderList := make([]OrderInfo, 0, 100)

	// get order id, key: user_*
	orderID, err := hc.getOrderID(userAddr)
	if err != nil {
		return nil, err
	}

	logger.Debug("user's order id:", orderID)

	// number of order equal to order id
	num, err := utils.StringToInt64(orderID)
	if err != nil {
		return nil, err
	}
	for i := int64(0); i < num; i++ {
		// make key
		key := OrderKey(userAddr, utils.Int64ToString(i))
		logger.Debug("order key:", key)
		// get order
		data, err := hc.LocalDB.Get([]byte(key))
		if err != nil {
			return nil, err
		}
		order := OrderInfo{}
		err = json.Unmarshal(data, &order)
		if err != nil {
			return nil, err
		}
		orderList = append(orderList, order)
	}

	return orderList, nil
}

// append db item into cps
func (hc *HandlerCore) appendResult(cps *[]CPInfo, item *badger.Item) error {
	// append each item
	err := item.Value(func(val []byte) error {
		logger.Debugf("Key:%s Value:%s", string(item.Key()), string(val))
		cp := CPInfo{}
		err := json.Unmarshal(val, &cp)
		if err != nil {
			return err
		}
		// append
		*cps = append(*cps, cp)
		return nil
	})
	if err != nil {
		return fmt.Errorf("error processing result: %w", err)
	}

	return nil
}

// add an account with some credit, return k,v for db write
func (hc *HandlerCore) addCredit(addr string, credit int64) (k []byte, v []byte, err error) {

	// credit key
	creKey := CreditKey(addr)

	// old credit
	var old string

	// get old credit from db, if key not found, init with 0
	data, err := hc.LocalDB.Get([]byte(creKey))
	if err != nil {
		if err.Error() == "Key not found" {
			old = "0"
		} else {
			return nil, nil, err
		}
	} else {
		old = string(data)
	}
	logger.Debug("old credit:", old)

	// accumulate credit
	old64, err := utils.StringToInt64(old)
	if err != nil {
		return nil, nil, err
	}
	new64 := old64 + credit
	new := utils.Int64ToString(new64)

	logger.Debug("new credit:", new)

	// return k v for multiput
	return creKey, []byte(new), nil
}

// append an order key for a cp
func (hc *HandlerCore) appendOrder(cpAddr string, orderKey string) (k []byte, v []byte, err error) {
	// 'cp' _ 'address' as cp key
	cpordersKey := OrderListKey(cpAddr)

	// key of orders
	var orders []string = make([]string, 0)

	// read orders from db
	data, err := hc.LocalDB.Get([]byte(cpordersKey))
	if err != nil {
		// if no order keys, init an empty data
		if err.Error() == "Key not found" {
			data = []byte{}
		} else {
			return nil, nil, err
		}
	}

	// if data not empty, unmarshal it to orders
	if len(data) != 0 {
		err = json.Unmarshal(data, &orders)
		if err != nil {
			return nil, nil, err
		}
	}

	// append into keys
	orders = append(orders, orderKey)

	data, err = json.Marshal(orders)
	if err != nil {
		return nil, nil, err
	}

	// return k,v for multiput
	return cpordersKey, data, nil
}
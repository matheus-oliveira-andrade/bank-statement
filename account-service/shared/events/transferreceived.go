package events

type TransferReceived struct {
	FromNumber string `json:"fromNumber"`
	ToNumber   string `json:"toNumber"`
	Value      int64  `json:"value"`
	Balance    int64  `json:"balance"`
}

func NewTransferReceived(fromNumber, toNumber string, value int64, balance int64) *TransferReceived {
	return &TransferReceived{
		FromNumber: fromNumber,
		ToNumber:   toNumber,
		Value:      value,
		Balance:    balance,
	}
}

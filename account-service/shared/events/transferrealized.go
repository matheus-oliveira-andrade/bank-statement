package events

type TransferRealized struct {
	FromNumber string `json:"fromNumber"`
	ToNumber   string `json:"toNumber"`
	Value      int64  `json:"value"`
	Balance    int64  `json:"balance"`
}

func NewTransferRealized(fromNumber, toNumber string, value int64, balance int64) *TransferRealized {
	return &TransferRealized{
		FromNumber: fromNumber,
		ToNumber:   toNumber,
		Value:      value,
		Balance:    balance,
	}
}

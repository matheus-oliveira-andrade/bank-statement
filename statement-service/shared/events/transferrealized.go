package events

type TransferRealized struct {
	FromNumber string `json:"fromNumber"`
	ToNumber   string `json:"toNumber"`
	Value      int64  `json:"value"`
	Balance    int64  `json:"balance"`
}

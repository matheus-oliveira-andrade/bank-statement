package events

type FundsDeposited struct {
	Number string `json:"number"`
	Value  int64  `json:"value"`
}

func NewFundsDeposited(number string, value int64) *FundsDeposited {
	return &FundsDeposited{
		Number: number,
		Value:  value,
	}
}

package events

type FundsDeposited struct {
	Number string `json:"number"`
	Value  int64  `json:"value"`
}

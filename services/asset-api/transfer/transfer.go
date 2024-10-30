package transfer

type TransferRequest struct {
	From    string  `json:"from" example:"0x123abc456def"`
	To      string  `json:"to" example:"0x987def456def"`
	Network string  `json:"network" example:"Ethereum"`
	Amount  float64 `json:"amount" example:"100.50"`
}

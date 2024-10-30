package transaction

type Request struct {
	From          string  `json:"from" example:"wallet123"`
	To            string  `json:"to" example:"wallet456"`
	Network       string  `json:"network" example:"mainnet"`
	Amount        float64 `json:"amount" example:"100.50"`
	ScheduledTime string  `json:"scheduled_time" example:"2023-12-31T12:00:00Z"`
}

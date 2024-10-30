package publisher

func NewScheduleTransactionEvent(toWalletId string, transactionId int) ScheduleTransactionEvent {
	return ScheduleTransactionEvent{
		Key: toWalletId,
		Value: struct {
			ID int `json:"id"`
		}{ID: transactionId},
	}
}

type ScheduleTransactionEvent struct {
	Key   string `json:"key"`
	Value struct {
		ID int `json:"id"`
	} `json:"value"`
}

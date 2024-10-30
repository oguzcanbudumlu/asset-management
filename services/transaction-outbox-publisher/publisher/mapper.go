package publisher

import "asset-management/internal/schedule"

// MapScheduleTransactionsToEvents maps an array of ScheduleTransaction objects to ScheduleTransactionEvent structs
func MapScheduleTransactionsToEvents(transactions []schedule.ScheduleTransaction) []ScheduleTransactionEvent {
	var events []ScheduleTransactionEvent

	for _, tx := range transactions {
		event := NewScheduleTransactionEvent(tx.ToWallet, tx.ID)
		events = append(events, event)
	}

	return events
}

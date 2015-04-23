package rtb

// The transaction logger interface is designed to assist with implementing a bomb proof transaction logger
type TransactionLogger interface {
	// A way to determine if someone is listening to this log. This is important, as we should not start bidding if nobody is going to consume the transaction
	ConsumerListening() (bool, error)
	// A return of nil (no error) means that a transaction was logged and acknoleged
	// Any error returned should be treated as though the transaction was not logged.
	LogTransaction(transaction *Transaction) error
}

// A transaction is kept light intentionally, as we don't want to "bog down" any accounting system
type Transaction struct {
	CampaignId             int64
	BidResponseId          string
	AmountInMicroCents     int64
	TimestampInNanoSeconds int64
	Ext                    interface{} // Extended data
}

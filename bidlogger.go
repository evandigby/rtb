package rtb

// BidLogItem defines the structure of a bidder log
type BidLogItem struct {
	Domain                            string           `json:"d,omitempty"`
	BidRequest                        *BidRequest      `json:"rq,omitempty"`
	BidResponse                       *BidResponse     `json:"rp,omitempty"`
	RemainingDailyBudgetsInMicroCents map[string]int64 `json:"b,omitempty"`
	StartTimestampInNanoseconds       int64            `json:"sts,omitempty"`
	EndTimestampInNanoseconds         int64            `json:"ets,omitempty"`
}

// BidLogProducer defines a type that can log bid requests
type BidLogProducer interface {
	// Safe inside goroutine
	LogItem(logItem *BidLogItem)
}

// BidLogConsumer defines a type that can consume bid log requests
type BidLogConsumer interface {
	LogChannel() chan *BidLogItem
}

// BidLogger defines a type that can both produce and consume bid log requests
type BidLogger interface {
	BidLogProducer
	BidLogConsumer
}

package rtb

// Bidder defines a type which can bid on bid requests
type Bidder interface {
	Bid() (response *BidResponse, campaignRemainingDailyBudget map[string]int64, err error)
}

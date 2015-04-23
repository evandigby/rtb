package rtb

import (
	"time"
)

// CampaignProvider defines a way to access a campaign data store
type CampaignProvider interface {
	// ReadByTargeting returns any campaigns that have available funds and meet the target criteria passed in.
	// Campaigns are returned in order from highest cpm to lowest
	// Available funds are is measured at the time of the query, and may be spent by the time DebitCampaign is called.
	ReadByTargeting(bidFloorInMicroCents int64, targets []Target) []Campaign

	// DebitCampaign subtracts an amount from the daily budget of the campaign
	// Returns the remaining daily budget after the transaction, and an error if the transaction was unsuccessful
	DebitCampaign(campaignId int64, amountInMicroCents int64, dailyBudgetExpiration time.Time) (remainingDailyBudgetInMicroCents int64, err error)

	// Creates a new persisted campaign
	CreateCampaign(campaignId int64, bidCpmInMicroCents int64, dailyBudgetInMicroCents int64, targets []Target) Campaign

	// Reads a persisted campaign
	ReadCampaign(campaignId int64) Campaign

	// Lists campaigns
	ListCampaigns() []int64
}

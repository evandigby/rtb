package mocks

import (
	"github.com/evandigby/rtb"
	"time"
)

type MockCampaignProvider struct {
	readByTargetingResult []rtb.Campaign
	debitCampaignResults  map[int64]int64
	debitCampaignErrors   map[int64]error

	campaigns map[int64]rtb.Campaign
}

func (cp *MockCampaignProvider) ReadByTargeting(bidFloorInMicroCents int64, targets []rtb.Target) []rtb.Campaign {
	return cp.readByTargetingResult
}

func (cp *MockCampaignProvider) DebitCampaign(campaignId int64, amountInMicroCents int64, dailyBudgetExpiration time.Time) (remainingDailyBudgetInMicroCents int64, err error) {
	return cp.debitCampaignResults[campaignId], cp.debitCampaignErrors[campaignId]
}

func (cp *MockCampaignProvider) CreateCampaign(campaignId int64, bidCpmInMicroCents int64, dailyBudgetInMicroCents int64, targets []rtb.Target) rtb.Campaign {
	return cp.campaigns[campaignId]
}

func (cp *MockCampaignProvider) ReadCampaign(campaignId int64) rtb.Campaign {
	return cp.campaigns[campaignId]
}

func (cp *MockCampaignProvider) ListCampaigns() []int64 {
	keys := make([]int64, 0, len(cp.campaigns))

	for k := range cp.campaigns {
		keys = append(keys, k)
	}
	return keys
}

// NewMockCampaignProvider creates a mock campaign.
// debitCampaignResults returns the result mapped to the campaignId
func NewMockCampaignProvider(readByTargetingResult []rtb.Campaign, debitCampaignResults map[int64]int64, debitCampaignErrors map[int64]error, campaigns map[int64]rtb.Campaign) rtb.CampaignProvider {
	cp := new(MockCampaignProvider)

	cp.readByTargetingResult = readByTargetingResult
	cp.debitCampaignResults = debitCampaignResults
	cp.debitCampaignErrors = debitCampaignErrors
	cp.campaigns = campaigns

	return cp
}

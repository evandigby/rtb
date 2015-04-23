package mocks

import (
	"github.com/evandigby/rtb"
)

type MockCampaign struct {
	campaignId              int64
	bidCpmInMicroCents      int64
	dailyBudgetInMicroCents int64
	targets                 *map[rtb.TargetType]string
}

func (c *MockCampaign) Id() int64 {
	return c.campaignId
}

func (c *MockCampaign) BidCpmInMicroCents() int64 {
	return c.bidCpmInMicroCents
}

func (c *MockCampaign) DailyBudgetInMicroCents() int64 {
	return c.dailyBudgetInMicroCents
}

func (c *MockCampaign) Targets() *map[rtb.TargetType]string {
	return c.targets
}

func NewMockCampaign(id int64, bidCpmInMicroCents int64, dailyBudgetInMicroCents int64, targets *map[rtb.TargetType]string) rtb.Campaign {
	c := new(MockCampaign)

	c.campaignId = id
	c.bidCpmInMicroCents = bidCpmInMicroCents
	c.dailyBudgetInMicroCents = dailyBudgetInMicroCents
	c.targets = targets

	return c
}

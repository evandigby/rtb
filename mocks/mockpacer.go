package mocks

import (
	"github.com/evandigby/rtb"
)

type MockPacer struct {
	canBidResults map[int64]bool
}

func (p *MockPacer) CanBid(campaign rtb.Campaign) bool {
	return p.canBidResults[campaign.Id()]
}

func NewMockPacer(canBidResults map[int64]bool) rtb.Pacer {
	p := new(MockPacer)

	p.canBidResults = canBidResults

	return p
}

package mocks

import (
	"github.com/evandigby/rtb"
)

type MockPacer struct {
	canBidResults map[int64]bool
}

func (p *MockPacer) CanBid(account int64) bool {
	return p.canBidResults[account]
}

func NewMockPacer(canBidResults map[int64]bool) rtb.Pacer {
	p := new(MockPacer)

	p.canBidResults = canBidResults

	return p
}

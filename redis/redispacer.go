package redis

import (
	"github.com/evandigby/rtb"
	"strconv"
	"time"
)

// Implements a simple time segmented pacer. Dividing the daily budget into time segments
// based on the segment duration provided
type RedisPacer struct {
	da     NoDbDataAccess
	banker rtb.Banker
	cp     rtb.CampaignProvider

	segment time.Duration
}

func (p *RedisPacer) paceAccountKey(account int64) string {
	return "pacer:account:" + strconv.FormatInt(account, 16)
}

func (p *RedisPacer) CanBid(campaign rtb.Campaign) bool {
	id := campaign.Id()
	key := p.paceAccountKey(campaign.Id())

	// Update remaining budget every time to compensate for unspent bids last cycle
	cpi := rtb.MicroCentsPerImpression(campaign.BidCpmInMicroCents())
	budget := p.banker.RemainingDailyBudgetInMicroCents(id)

	// Account not yet configured, or they actually have no budget. The pacer is not needed.
	if budget == 0 || cpi == 0 {
		return true // Allow it to pass and configure it
	}

	now := time.Now().UTC()
	midnight := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, time.UTC)

	// TODO: This definitely only needs to be calculated at the end of each segment. Not every time we request a bid.
	durationToMidnight := midnight.Sub(now)
	segmentsToMidnight := int64(durationToMidnight / p.segment)
	remaining := (budget / cpi) / segmentsToMidnight

	remainingBudget, err := p.da.DebitIfNotZero(key, 1, remaining, time.Now().UTC().Add(p.segment))

	//	fmt.Printf("Campign: %v, Bids Per Segment: %v, Segment: %v, Remaining Budget: %v\n", id, remaining, p.segment, remainingBudget)

	if err == nil {
		return remainingBudget > 0
	} else {
		return false
	}
}

func (p *RedisPacer) Segment() time.Duration {
	return p.segment
}

func NewRedisPacer(cp rtb.CampaignProvider, da NoDbDataAccess, banker rtb.Banker, segment time.Duration) rtb.Pacer {
	p := new(RedisPacer)

	p.da = da
	p.banker = banker
	p.cp = cp
	p.segment = segment

	return p
}

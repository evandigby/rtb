package inmemory

import (
	"github.com/evandigby/rtb"
	"strconv"
	"time"
)

type BidRequestBidder struct {
	Request                   *rtb.BidRequest
	CampaignProvider          rtb.CampaignProvider
	Pacer                     rtb.Pacer
	dailyBudgetExpirationTime time.Time
}

// highestBidder Returns the highest bidder with available funds.
// campaigns must be in order from highest CPM to lowest
// Caller is committed to using the campaign returned by this
func (b *BidRequestBidder) highestBidder(campaigns []rtb.Campaign) (campaign rtb.Campaign, remainingDailyBudgetInMicroCents int64) {
	for _, campaign := range campaigns {
		id := campaign.Id()

		if b.Pacer.CanBid(id) {
			bid := campaign.BidCpmInMicroCents()

			if remainingDailyBudgetInMicroCents, err := b.CampaignProvider.DebitCampaign(id, rtb.MicroCentsPerImpression(bid), b.dailyBudgetExpirationTime); err == nil {
				return campaign, remainingDailyBudgetInMicroCents
			}
		}
	}

	return nil, 0
}

// If the bid is nil, the remaining remainingDailyBudgetInMicroCents and campaign id are invalid
func (b *BidRequestBidder) impressionBid(imp *rtb.Imp, userTargets []rtb.Target) (bid *rtb.Bid, remainingDailyBudgetInMicroCents int64) {
	targets := append(userTargets, imp.Targeting()...)

	campaigns := b.CampaignProvider.ReadByTargeting(rtb.CpmToMicroCents(imp.Bidfloor), targets)
	if len(campaigns) == 0 {
		return nil, 0
	}

	// Returns nil if none of the campaigns have available budget at the time of the call
	// We need to double check this, in case the budget is spent by the time we've decided to bid.
	// To be fair, we're committed to using the result of this call.
	campaign, remainingDailyBudgetInMicroCents := b.highestBidder(campaigns)

	// Either the pacer rejected them all, or none of them had available remainingDailyBudgetInMicroCents
	if campaign == nil {
		return nil, 0
	}

	bid = new(rtb.Bid)
	bid.Price = rtb.MicroCentsToCpm(campaign.BidCpmInMicroCents())
	bid.Impid = imp.ID
	bid.Cid = strconv.FormatInt(campaign.Id(), 10)

	// Stub the rest of the fields
	bid.ID = "stub"
	bid.Adid = "stub"
	bid.Adm = "<span>stub</span>"
	bid.Adomain = []string{"stub.go2mobi.com"}
	bid.Crid = "stub"

	return bid, remainingDailyBudgetInMicroCents
}

// BidResponse returns nil if there is no bid
func (b *BidRequestBidder) Bid() (response *rtb.BidResponse, campaignRemainingDailyBudgetsInMicroCents map[string]int64, err error) {
	targets := b.Request.Targeting()

	campaignRemainingDailyBudgetsInMicroCents = make(map[string]int64)
	// Allocate up to the amount of impressions
	bids := make([]rtb.Bid, 0, len(b.Request.Imp))

	for _, imp := range b.Request.Imp {
		ibid, remainingDailyBudgetInMicroCents := b.impressionBid(&imp, targets)

		if ibid != nil {
			bids = append(bids, *ibid)
			campaignRemainingDailyBudgetsInMicroCents[ibid.Cid] = remainingDailyBudgetInMicroCents
		}
	}

	// No bids
	if len(bids) <= 0 {
		return nil, campaignRemainingDailyBudgetsInMicroCents, nil
	}

	// We have bids to submit! Build a response
	response = new(rtb.BidResponse)

	response.Bidid = "stub"
	response.Cur = "USD"
	response.Seatbid = make([]rtb.Seatbid, 1)
	response.Seatbid[0].Bid = bids

	return response, campaignRemainingDailyBudgetsInMicroCents, nil
}

func NewBidRequestBidder(r *rtb.BidRequest, cp rtb.CampaignProvider, pacer rtb.Pacer, now time.Time) rtb.Bidder {
	b := new(BidRequestBidder)

	b.Request = r
	b.CampaignProvider = cp
	b.Pacer = pacer

	// Budgets expire on calendar days (not 24 hour periods)
	b.dailyBudgetExpirationTime = time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, time.UTC)

	return b
}

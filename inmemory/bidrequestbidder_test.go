package inmemory

import (
	"github.com/evandigby/rtb"
	"github.com/evandigby/rtb/mocks"
	"strconv"
	"testing"
	"time"
)

// TestBiddingMatchedTarget tests that the bidder will bid on a campaign when it is returned by the campaign provider
// Expected result is a bid response that is populated by a single bid with a campaign ID matching the campaign returned by the provider
func TestBiddingMatchedTarget(t *testing.T) {
	r := new(rtb.BidRequest)
	r.Imp = make([]rtb.Imp, 1)
	r.Imp[0] = rtb.Imp{}

	pacer := mocks.NewMockPacer(map[int64]bool{100: true})

	campaign := mocks.NewMockCampaign(100, rtb.CpmToMicroCents(0.25), rtb.DollarsToMicroCents(1), nil)

	campaignsReturnedByTargeting := []rtb.Campaign{campaign}

	cp := mocks.NewMockCampaignProvider(campaignsReturnedByTargeting, map[int64]int64{100: rtb.DollarsToMicroCents(1)}, map[int64]error{100: nil}, nil)

	b := NewBidRequestBidder(r, cp, pacer, time.Now().UTC())

	response, _, err := b.Bid()

	if err != nil {
		t.Fail()
	}

	if response == nil {
		t.FailNow()
	}

	if response.Seatbid == nil {
		t.FailNow()
	}

	if len(response.Seatbid) != 1 {
		t.Fail()
	}

	if response.Seatbid[0].Bid == nil {
		t.FailNow()
	}

	if len(response.Seatbid[0].Bid) != 1 {
		t.FailNow()
	}

	if response.Seatbid[0].Bid[0].Cid != strconv.FormatInt(campaignsReturnedByTargeting[0].Id(), 10) {
		t.Fail()
	}
}

// TestBiddingMatchedTargetCorrectAmount tests that the bidder bids the correct amount
// Expected result is a bid response that is populated by a single bid with the correct bid amount
func TestBiddingMatchedTargetCorrectAmount(t *testing.T) {
	r := new(rtb.BidRequest)
	r.Imp = make([]rtb.Imp, 1)
	r.Imp[0] = rtb.Imp{}

	pacer := mocks.NewMockPacer(map[int64]bool{100: true})

	campaign := mocks.NewMockCampaign(100, rtb.CpmToMicroCents(0.32), rtb.DollarsToMicroCents(1), nil)

	campaignsReturnedByTargeting := []rtb.Campaign{campaign}

	cp := mocks.NewMockCampaignProvider(campaignsReturnedByTargeting, map[int64]int64{100: rtb.DollarsToMicroCents(1)}, map[int64]error{100: nil}, nil)

	b := NewBidRequestBidder(r, cp, pacer, time.Now().UTC())

	expectedBidAmount := float64(0.32)

	response, _, err := b.Bid()

	if err != nil {
		t.Fail()
	}

	if response == nil {
		t.FailNow()
	}

	if response.Seatbid == nil {
		t.FailNow()
	}

	if len(response.Seatbid) != 1 {
		t.Fail()
	}

	if response.Seatbid[0].Bid == nil {
		t.FailNow()
	}

	if len(response.Seatbid[0].Bid) != 1 {
		t.FailNow()
	}

	if response.Seatbid[0].Bid[0].Cid != strconv.FormatInt(campaignsReturnedByTargeting[0].Id(), 10) {
		t.Fail()
	}

	if response.Seatbid[0].Bid[0].Price != expectedBidAmount {
		t.Fail()
	}
}

// TestBiddingNoMatchedTarget tests that the bidder will not bid when no campaign is returned by the campaign provider
// Expected result is nil
func TestBiddingNoMatchedTarget(t *testing.T) {
	r := new(rtb.BidRequest)
	r.Imp = make([]rtb.Imp, 1)
	r.Imp[0] = rtb.Imp{}

	pacer := mocks.NewMockPacer(map[int64]bool{100: true})

	campaignsReturnedByTargeting := []rtb.Campaign{}

	cp := mocks.NewMockCampaignProvider(campaignsReturnedByTargeting, map[int64]int64{100: rtb.DollarsToMicroCents(1)}, map[int64]error{100: nil}, nil)

	b := NewBidRequestBidder(r, cp, pacer, time.Now().UTC())

	response, _, err := b.Bid()

	if err != nil {
		t.Fail()
	}

	if response != nil {
		t.FailNow()
	}
}

// TestBiddingMatchedTarget tests that the bidder will not bid if a campaign is out of funds
// Expected result is nil
func TestBiddingMatchedTargetNoFundsAvailable(t *testing.T) {
	r := new(rtb.BidRequest)
	r.Imp = make([]rtb.Imp, 1)
	r.Imp[0] = rtb.Imp{}

	pacer := mocks.NewMockPacer(map[int64]bool{100: true})

	campaign := mocks.NewMockCampaign(100, rtb.CpmToMicroCents(0.25), rtb.DollarsToMicroCents(1), nil)

	campaignsReturnedByTargeting := []rtb.Campaign{campaign}

	cp := mocks.NewMockCampaignProvider(campaignsReturnedByTargeting, map[int64]int64{100: rtb.DollarsToMicroCents(1)}, map[int64]error{100: rtb.NewTransactionError("error", true)}, nil)

	b := NewBidRequestBidder(r, cp, pacer, time.Now().UTC())

	response, _, err := b.Bid()

	if err != nil {
		t.Fail()
	}

	if response != nil {
		t.FailNow()
	}
}

// TestBiddingMultipleMatchedTarget tests that the bidder will bid on the first campaign returned by the bidder
// Expected result is a bid response that is populated by a single bid with a campaign ID matching the first campaigned returned by the provider
func TestBiddingMultipleMatchedTarget(t *testing.T) {
	r := new(rtb.BidRequest)
	r.Imp = make([]rtb.Imp, 1)
	r.Imp[0] = rtb.Imp{}

	pacer := mocks.NewMockPacer(map[int64]bool{100: true})

	campaign1 := mocks.NewMockCampaign(100, rtb.CpmToMicroCents(0.25), rtb.DollarsToMicroCents(1), nil)
	campaign2 := mocks.NewMockCampaign(101, rtb.CpmToMicroCents(0.25), rtb.DollarsToMicroCents(1), nil)

	campaignsReturnedByTargeting := []rtb.Campaign{campaign1, campaign2}

	cp := mocks.NewMockCampaignProvider(campaignsReturnedByTargeting, map[int64]int64{100: rtb.DollarsToMicroCents(1)}, map[int64]error{100: nil}, nil)

	b := NewBidRequestBidder(r, cp, pacer, time.Now().UTC())

	response, _, err := b.Bid()

	if err != nil {
		t.Fail()
	}

	if response == nil {
		t.FailNow()
	}

	if response.Seatbid == nil {
		t.FailNow()
	}

	if len(response.Seatbid) != 1 {
		t.Fail()
	}

	if response.Seatbid[0].Bid == nil {
		t.FailNow()
	}

	if len(response.Seatbid[0].Bid) != 1 {
		t.FailNow()
	}

	if response.Seatbid[0].Bid[0].Cid != strconv.FormatInt(campaignsReturnedByTargeting[0].Id(), 10) {
		t.Fail()
	}
}

// TestBiddingMatchedTarget tests that the bidder will bid on the first campaign returned by the bidder that has a remaining daily budget
// Expected result is a bid response that is populated by a single bid with a campaign ID matching the first campaigned returned by the provider
func TestBiddingFirstCampaignOutOfFunds(t *testing.T) {
	r := new(rtb.BidRequest)
	r.Imp = make([]rtb.Imp, 1)
	r.Imp[0] = rtb.Imp{}

	pacer := mocks.NewMockPacer(map[int64]bool{100: true, 101: true})

	campaign1 := mocks.NewMockCampaign(100, rtb.CpmToMicroCents(0.25), rtb.DollarsToMicroCents(1), nil)
	campaign2 := mocks.NewMockCampaign(101, rtb.CpmToMicroCents(0.25), rtb.DollarsToMicroCents(1), nil)

	campaignsReturnedByTargeting := []rtb.Campaign{campaign1, campaign2}

	cp := mocks.NewMockCampaignProvider(campaignsReturnedByTargeting, map[int64]int64{100: rtb.DollarsToMicroCents(1), 101: rtb.DollarsToMicroCents(1)}, map[int64]error{100: rtb.NewTransactionError("test", true), 101: nil}, nil)

	b := NewBidRequestBidder(r, cp, pacer, time.Now().UTC())

	response, _, err := b.Bid()

	if err != nil {
		t.Fail()
	}

	if response == nil {
		t.FailNow()
	}

	if response.Seatbid == nil {
		t.FailNow()
	}

	if len(response.Seatbid) != 1 {
		t.Fail()
	}

	if response.Seatbid[0].Bid == nil {
		t.FailNow()
	}

	if len(response.Seatbid[0].Bid) != 1 {
		t.FailNow()
	}

	if response.Seatbid[0].Bid[0].Cid != strconv.FormatInt(campaignsReturnedByTargeting[1].Id(), 10) {
		t.Fail()
	}
}

// TestBiddingPacerFalse tests that the bidder will not bid when the pacer replys false to CanBid
// Expected result is nil
func TestBiddingPacerFalse(t *testing.T) {
	r := new(rtb.BidRequest)
	r.Imp = make([]rtb.Imp, 1)
	r.Imp[0] = rtb.Imp{}

	pacer := mocks.NewMockPacer(map[int64]bool{100: false})

	campaign := mocks.NewMockCampaign(100, rtb.CpmToMicroCents(0.25), rtb.DollarsToMicroCents(1), nil)

	campaignsReturnedByTargeting := []rtb.Campaign{campaign}

	cp := mocks.NewMockCampaignProvider(campaignsReturnedByTargeting, map[int64]int64{100: rtb.DollarsToMicroCents(1)}, map[int64]error{100: nil}, nil)

	b := NewBidRequestBidder(r, cp, pacer, time.Now().UTC())

	response, _, err := b.Bid()

	if err != nil {
		t.Fail()
	}

	if response != nil {
		t.FailNow()
	}
}

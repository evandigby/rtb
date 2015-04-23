package redis

import (
	"github.com/evandigby/rtb"
	"testing"
	"time"
)

// Test creating a campaign
// Expected result is a campaign is returned and the member values match the inputs
func TestCreateCampaign(t *testing.T) {
	b := NewRedisBanker(testDataAccess)
	cp := NewRedisCampaignProvider(testDataAccess, b)

	campaignId := int64(300)
	bidCpmInMicroCents := int64(100)
	dailyBudgetInMicroCents := int64(100)
	target := rtb.Target{Type: rtb.Placement, Value: "Words With Friends 2 iPad"}
	targets := []rtb.Target{target}

	c := cp.CreateCampaign(campaignId, bidCpmInMicroCents, dailyBudgetInMicroCents, targets)

	if c.Id() != campaignId {
		t.Fail()
	}

	if c.BidCpmInMicroCents() != bidCpmInMicroCents {
		t.Fail()
	}

	if c.DailyBudgetInMicroCents() != dailyBudgetInMicroCents {
		t.Fail()
	}

	resultTargets := *c.Targets()

	if len(targets) != len(resultTargets) {
		t.FailNow()
	}

	if resultTargets[target.Type] != target.Value {
		t.Fail()
	}
}

// Test creating a campaign and then reading it back
// Expected result is a campaign is returned and the member values match the inputs
func TestReadCampaign(t *testing.T) {
	b := NewRedisBanker(testDataAccess)
	cp := NewRedisCampaignProvider(testDataAccess, b)

	campaignId := int64(301)
	bidCpmInMicroCents := int64(100)
	dailyBudgetInMicroCents := int64(100)
	target := rtb.Target{Type: rtb.Placement, Value: "Words With Friends 2 iPad"}
	targets := []rtb.Target{target}

	c := cp.CreateCampaign(campaignId, bidCpmInMicroCents, dailyBudgetInMicroCents, targets)

	c = cp.ReadCampaign(campaignId)

	if c.Id() != campaignId {
		t.Fail()
	}

	if c.BidCpmInMicroCents() != bidCpmInMicroCents {
		t.Fail()
	}

	if c.DailyBudgetInMicroCents() != dailyBudgetInMicroCents {
		t.Fail()
	}

	resultTargets := *c.Targets()

	if len(targets) != len(resultTargets) {
		t.FailNow()
	}

	if resultTargets[target.Type] != target.Value {
		t.Fail()
	}
}

// Test creating a campaign and then reading it back by targeting
// Expected result is the correct campaign is returned
func TestReadCampaignByTargetingMatch(t *testing.T) {
	b := NewRedisBanker(testDataAccess)
	cp := NewRedisCampaignProvider(testDataAccess, b)

	campaignId := int64(302)
	bidCpmInMicroCents := int64(100)
	dailyBudgetInMicroCents := int64(100)
	target := rtb.Target{Type: rtb.Placement, Value: "Unique Targeting 1"}
	targets := []rtb.Target{target}
	numCampaigns := 1

	c := cp.CreateCampaign(campaignId, bidCpmInMicroCents, dailyBudgetInMicroCents, targets)

	campaigns := cp.ReadByTargeting(0, targets)

	if len(campaigns) != numCampaigns {
		t.Fail()
	}

	c = campaigns[0]

	if c.Id() != campaignId {
		t.Fail()
	}
}

// Test creating a campaign and test it is not returned by reading unmatching targets
// Expected result is a campaign is not returned
func TestReadCampaignByTargetingNoMatch(t *testing.T) {
	b := NewRedisBanker(testDataAccess)
	cp := NewRedisCampaignProvider(testDataAccess, b)

	campaignId := int64(303)
	bidCpmInMicroCents := int64(100)
	dailyBudgetInMicroCents := int64(100)
	target := rtb.Target{Type: rtb.Placement, Value: "Unique Targeting With No Match"}
	targets := []rtb.Target{target}

	unMatchingTarget := rtb.Target{Type: rtb.Placement, Value: "Unique Targeting That Doesn't Match"}
	unMatchingTargetS := []rtb.Target{unMatchingTarget}

	numCampaigns := 0

	cp.CreateCampaign(campaignId, bidCpmInMicroCents, dailyBudgetInMicroCents, targets)

	campaigns := cp.ReadByTargeting(0, unMatchingTargetS)

	if len(campaigns) != numCampaigns {
		t.Fail()
	}
}

// Test creating a campaign and then reading it back by targeting with multiple possible targets
// Expected result is the correct campaign is returned
func TestReadCampaignByTargetingMatchMultipleTargets(t *testing.T) {
	b := NewRedisBanker(testDataAccess)
	cp := NewRedisCampaignProvider(testDataAccess, b)

	campaignId := int64(304)
	bidCpmInMicroCents := int64(100)
	dailyBudgetInMicroCents := int64(100)
	target1 := rtb.Target{Type: rtb.Placement, Value: "Unique Targeting 3"}
	target2 := rtb.Target{Type: rtb.CreativeSize, Value: "123x456"}
	targets := []rtb.Target{target1}
	targetsToMatch := []rtb.Target{target1, target2}
	numCampaigns := 1

	c := cp.CreateCampaign(campaignId, bidCpmInMicroCents, dailyBudgetInMicroCents, targets)

	campaigns := cp.ReadByTargeting(0, targetsToMatch)

	if len(campaigns) != numCampaigns {
		t.FailNow()
	}

	c = campaigns[0]

	if c.Id() != campaignId {
		t.Fail()
	}
}

// Test creating 2 campaigns and then reading both back by targeting with multiple possible targets
// Expected result is the correct campaign is returned
func TestReadCampaignByTargetingMatchMultipleCampaignsSingleTargetEach(t *testing.T) {
	b := NewRedisBanker(testDataAccess)
	cp := NewRedisCampaignProvider(testDataAccess, b)

	campaignId1 := int64(305)
	campaignId2 := int64(306)
	bidCpmInMicroCents := int64(100)
	dailyBudgetInMicroCents := int64(100)
	target1 := rtb.Target{Type: rtb.Placement, Value: "Unique Targeting 4"}
	target2 := rtb.Target{Type: rtb.CreativeSize, Value: "234x567"}
	targets1 := []rtb.Target{target1}
	targets2 := []rtb.Target{target2}

	targetsToMatch := []rtb.Target{target1, target2}
	numCampaigns := 2

	cp.CreateCampaign(campaignId1, bidCpmInMicroCents, dailyBudgetInMicroCents, targets1)
	cp.CreateCampaign(campaignId2, bidCpmInMicroCents, dailyBudgetInMicroCents, targets2)

	campaigns := cp.ReadByTargeting(0, targetsToMatch)

	if len(campaigns) != numCampaigns {
		t.Fail()
	}
}

// Test creating 2 campaigns and then reading one back by targeting
// Expected result is the correct campaign is returned
func TestReadCampaignByTargetingMatchMultipleCampaignsSingleTargetMatch(t *testing.T) {
	b := NewRedisBanker(testDataAccess)
	cp := NewRedisCampaignProvider(testDataAccess, b)

	campaignId1 := int64(307)
	campaignId2 := int64(308)
	bidCpmInMicroCents := int64(100)
	dailyBudgetInMicroCents := int64(100)
	target1 := rtb.Target{Type: rtb.Placement, Value: "Unique Targeting 5"}
	target2 := rtb.Target{Type: rtb.CreativeSize, Value: "345x678"}
	targets1 := []rtb.Target{target1}
	targets2 := []rtb.Target{target2}

	targetsToMatch := []rtb.Target{target1}
	numCampaigns := 1

	cp.CreateCampaign(campaignId1, bidCpmInMicroCents, dailyBudgetInMicroCents, targets1)
	cp.CreateCampaign(campaignId2, bidCpmInMicroCents, dailyBudgetInMicroCents, targets2)

	campaigns := cp.ReadByTargeting(0, targetsToMatch)

	if len(campaigns) != numCampaigns {
		t.Fail()
	}
}

// Test creating 4 campaigns and then reading all back and ensuring they're in the correct order
// Expected result is all campaigns are returned in order from highest bid cpm to lowest
func TestReadCampaignByTargetingOrder(t *testing.T) {
	b := NewRedisBanker(testDataAccess)
	cp := NewRedisCampaignProvider(testDataAccess, b)

	campaignId1 := int64(309)
	campaignId2 := int64(310)
	campaignId3 := int64(311)
	campaignId4 := int64(312)
	bidCpmInMicroCents1 := int64(100)
	bidCpmInMicroCents2 := int64(104)
	bidCpmInMicroCents3 := int64(102)
	bidCpmInMicroCents4 := int64(99)
	dailyBudgetInMicroCents := int64(100)
	target := rtb.Target{Type: rtb.Placement, Value: "Unique Targeting 6"}
	targets := []rtb.Target{target}

	expectedResults := []int64{campaignId2, campaignId3, campaignId1, campaignId4}

	cp.CreateCampaign(campaignId1, bidCpmInMicroCents1, dailyBudgetInMicroCents, targets)
	cp.CreateCampaign(campaignId2, bidCpmInMicroCents2, dailyBudgetInMicroCents, targets)
	cp.CreateCampaign(campaignId3, bidCpmInMicroCents3, dailyBudgetInMicroCents, targets)
	cp.CreateCampaign(campaignId4, bidCpmInMicroCents4, dailyBudgetInMicroCents, targets)

	campaigns := cp.ReadByTargeting(0, targets)

	if len(campaigns) != len(expectedResults) {
		t.FailNow()
	}

	for i, campaign := range campaigns {
		if campaign.Id() != expectedResults[i] {
			t.Fail()
		}
	}
}

// Test a debit when the account has enough
// Expected result is a successful debit with a remainingDailyBudgetInMicroCents of the initial remainingDailyBudgetInMicroCents minus the debit
func TestDebitCampaignAccountWithEnoughRemainingDailyBudgetInMicroCents(t *testing.T) {
	b := NewRedisBanker(testDataAccess)
	cp := NewRedisCampaignProvider(testDataAccess, b)

	campaignId := int64(313)
	bidCpmInMicroCents := int64(100)
	dailyBudgetInMicroCents := int64(100)
	target := rtb.Target{Type: rtb.Placement, Value: "Words With Friends 2 iPad"}
	targets := []rtb.Target{target}

	dailyBudgetExpiration := time.Now().UTC().AddDate(0, 0, 1)

	amount := int64(32)

	cp.CreateCampaign(campaignId, bidCpmInMicroCents, dailyBudgetInMicroCents, targets)

	expectedRemainingDailyBudgetInMicroCents := dailyBudgetInMicroCents - amount

	b.SetRemainingDailyBudgetInMicroCents(campaignId, dailyBudgetInMicroCents, dailyBudgetExpiration)

	result, err := cp.DebitCampaign(campaignId, amount, dailyBudgetExpiration)

	updatedRemainingDailyBudgetInMicroCents := b.RemainingDailyBudgetInMicroCents(campaignId)

	if err != nil {
		t.Fail()
	}

	if result != expectedRemainingDailyBudgetInMicroCents {
		t.Fail()
	}

	if updatedRemainingDailyBudgetInMicroCents != expectedRemainingDailyBudgetInMicroCents {
		t.Fail()
	}
}

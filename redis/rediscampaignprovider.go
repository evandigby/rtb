package redis

import (
	"github.com/evandigby/rtb"
	"strconv"
	"time"
)

type RedisCampaignProvider struct {
	appDomain      string
	network        string
	addr           string
	campaignSetKey string

	banker rtb.Banker
	da     NoDbDataAccess
}

func (cp *RedisCampaignProvider) ReadByTargeting(bidFloorInMicroCents int64, targets []rtb.Target) []rtb.Campaign {
	keys := TargetKeysForTargets("targets:", targets)

	reply := cp.da.SortedSetUnion(keys)

	campaigns := make([]rtb.Campaign, 0, len(reply))

	for _, campaignId := range reply {
		id, err := strconv.ParseInt(campaignId, 10, 64)
		if err != nil {
			panic(err)
		}

		campaign := cp.ReadCampaign(id)
		campaigns = append(campaigns, campaign)
	}

	return campaigns
}

func (cp *RedisCampaignProvider) DebitCampaign(campaignId int64, amountInMicroCents int64, dailyBudgetExpiration time.Time) (remainingDailyBudgetInMicroCents int64, err error) {
	campaign := cp.ReadCampaign(campaignId)

	dailyBudget := campaign.DailyBudgetInMicroCents()

	return cp.banker.DebitAccount(campaignId, amountInMicroCents, dailyBudget, dailyBudgetExpiration)
}

func TargetKeysForTargets(prepend string, targets []rtb.Target) []string {
	values := make([]string, 0, len(targets))

	for _, target := range targets {
		key := prepend + strconv.FormatInt(int64(target.Type), 10) + ":" + target.Value
		values = append(values, key)
	}

	return values
}

func (cp *RedisCampaignProvider) addTargetsToCampaign(targetKey string, targets []rtb.Target) {
	targetKeys := TargetKeysForTargets("", targets)
	members := make([]interface{}, len(targetKeys))
	for i, v := range targetKeys {
		members[i] = v
	}

	cp.da.AddMembersToSet(targetKey, members)
}

func (cp *RedisCampaignProvider) CreateCampaign(campaignId int64, bidCpmInMicroCents int64, dailyBudgetInMicroCents int64, targets []rtb.Target) rtb.Campaign {
	accountKey := "campaign:" + strconv.FormatInt(campaignId, 16)
	targetKey := accountKey + ":targets"

	campaign := NewRedisCampaign(cp.da, campaignId, accountKey, targetKey)

	cp.da.HSetInt64(accountKey, "bidCpmInMicroCents", bidCpmInMicroCents)
	cp.da.HSetInt64(accountKey, "dailyBudgetInMicroCents", dailyBudgetInMicroCents)

	cp.addTargetsToCampaign(targetKey, targets)

	cp.da.AddMembersToSet(cp.campaignSetKey, []interface{}{campaignId})

	members := make(map[string]SortedSetMember, len(targets))

	for _, target := range targets {
		targetKey := "targets:" + strconv.FormatInt(int64(target.Type), 10) + ":" + target.Value
		members[targetKey] = SortedSetMember{Member: campaignId, Score: bidCpmInMicroCents}
	}

	cp.da.AddMembersToSortedSets(members)

	return campaign
}

func (cp *RedisCampaignProvider) ListCampaigns() []int64 {

	keysAsString := cp.da.GetSetMembers(cp.campaignSetKey)

	keys := make([]int64, 0, len(keysAsString))

	for _, val := range keysAsString {
		i, err := strconv.ParseInt(val, 10, 64)
		if err == nil {
			keys = append(keys, i)
		}
	}

	return keys
}

func (cp *RedisCampaignProvider) campaignAccountKey(campaignId int64) string {
	return "campaign:" + strconv.FormatInt(campaignId, 16)
}

func (cp *RedisCampaignProvider) campaignTargetKey(accountKey string) string {
	return accountKey + ":targets"
}

func (cp *RedisCampaignProvider) ReadCampaign(campaignId int64) rtb.Campaign {
	accountKey := cp.campaignAccountKey(campaignId)
	targetKey := cp.campaignTargetKey(accountKey)

	return NewRedisCampaign(cp.da, campaignId, accountKey, targetKey)
}

func NewRedisCampaignProvider(da NoDbDataAccess, banker rtb.Banker) rtb.CampaignProvider {
	cp := new(RedisCampaignProvider)
	cp.da = da

	cp.banker = banker
	cp.campaignSetKey = "campaigns"

	return cp
}

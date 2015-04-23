package redis

import (
	"github.com/evandigby/rtb"
	"strconv"
	"strings"
)

type RedisCampaign struct {
	accountKey   string
	targetSetKey string

	campaignId int64

	bidCpmInMicroCents            int64
	bidCpmInMicroCentsLoaded      bool
	dailyBudgetInMicroCents       int64
	dailyBudgetInMicroCentsLoaded bool
	targets                       *map[rtb.TargetType]string
	targetsLoaded                 bool

	da NoDbDataAccess
}

func (c *RedisCampaign) loadBidCpmInMicroCents() {
	c.bidCpmInMicroCents = c.da.HGetInt64(c.accountKey, "bidCpmInMicroCents")
	c.bidCpmInMicroCentsLoaded = true
}

func (c *RedisCampaign) loadDailyBudgetInMicroCents() {
	c.dailyBudgetInMicroCents = c.da.HGetInt64(c.accountKey, "dailyBudgetInMicroCents")
	c.dailyBudgetInMicroCentsLoaded = true
}

func (c *RedisCampaign) loadTargets() {
	reply := c.da.GetSetMembers(c.targetSetKey)

	targets := make(map[rtb.TargetType]string, len(reply))
	for _, target := range reply {
		kvSplit := strings.SplitN(target, ":", 2)

		if len(kvSplit) < 2 {
			panic("not enough strings")
		}

		key, err := strconv.ParseInt(kvSplit[0], 10, 64)

		if err != nil {
			panic(err)
		}

		targets[rtb.TargetType(key)] = kvSplit[1]
	}

	c.targets = &targets
	c.targetsLoaded = true
}

func (c *RedisCampaign) Id() int64 {
	return c.campaignId
}

func (c *RedisCampaign) BidCpmInMicroCents() int64 {
	if !c.bidCpmInMicroCentsLoaded {
		c.loadBidCpmInMicroCents()
	}
	return c.bidCpmInMicroCents
}

func (c *RedisCampaign) DailyBudgetInMicroCents() int64 {
	if !c.dailyBudgetInMicroCentsLoaded {
		c.loadDailyBudgetInMicroCents()
	}
	return c.dailyBudgetInMicroCents
}

func (c *RedisCampaign) Targets() *map[rtb.TargetType]string {
	if !c.targetsLoaded {
		c.loadTargets()
	}
	return c.targets
}

func NewRedisCampaign(da NoDbDataAccess, id int64, accountKey string, targetSetKey string) rtb.Campaign {
	c := new(RedisCampaign)

	c.da = da
	c.campaignId = id

	c.accountKey = accountKey
	c.targetSetKey = targetSetKey

	return c
}

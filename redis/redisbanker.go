package redis

import (
	"github.com/evandigby/rtb"
	"strconv"
	"time"
)

type RedisBanker struct {
	da NoDbDataAccess
}

func (b *RedisBanker) accountKey(account int64) string {
	return "banker:account:" + strconv.FormatInt(account, 16)
}

func (b *RedisBanker) DeleteAccount(account int64) {
	b.da.DeleteKeys([]string{b.accountKey(account)})
}

func (b *RedisBanker) DebitAccount(account int64, amount int64, dailyBudget int64, dailyBudgetExpiration time.Time) (remainingDailyBudgetInMicroCents int64, err error) {
	accountKey := b.accountKey(account)

	return b.da.DebitIfNotZero(accountKey, amount, dailyBudget, dailyBudgetExpiration)
}

func (b *RedisBanker) RemainingDailyBudgetInMicroCents(account int64) int64 {
	accountKey := b.accountKey(account)

	success, val := b.da.GetInt64(accountKey)

	if success {
		return val
	} else {
		return 0
	}
}

func (b *RedisBanker) SetRemainingDailyBudgetInMicroCents(account int64, amount int64, dailyBudgetExpiration time.Time) {
	accountKey := b.accountKey(account)

	b.da.SetInt64(accountKey, amount)

	b.da.ExpireKey(accountKey, dailyBudgetExpiration)
}

func NewRedisBanker(da NoDbDataAccess) rtb.Banker {
	b := new(RedisBanker)
	b.da = da

	return b
}

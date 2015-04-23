package redis

import (
	"time"
)

type NoDbDataAccess interface {
	DeleteKeys(keys []string)
	HGetInt64(accountKey string, key string) int64
	HSetInt64(accountKey string, key string, val int64)
	GetInt64(accountKey string) (success bool, result int64)
	SetInt64(accountKey string, val int64)
	GetSetMembers(setKey string) []string
	AddMembersToSet(setKey string, members []interface{})
	AddMembersToSortedSets(keyValues map[string]SortedSetMember)
	SortedSetUnion(keys []string) []string
	DebitIfNotZero(accountKey string, amount int64, dailyBudget int64, dailyBudgetExpiration time.Time) (remainingDailyBudgetInMicroCents int64, err error)
	ExpireKey(key string, expirationTime time.Time)
	// Should only be used for testing
	GetKeys() []string
}

type SortedSetMember struct {
	Score  int64
	Member interface{}
}

package mocks

import (
	"github.com/evandigby/rtb"
	"time"
)

type MockBanker struct {
	debitAccountResult                     int64
	debitAccountError                      error
	remainingDailyBudgetInMicroCentsResult int64
}

func (b *MockBanker) DeleteAccount(account int64) {

}

func (b *MockBanker) DebitAccount(account int64, amount int64, dailyBudget int64, dailyBudgetExpiration time.Time) (remainingDailyBudgetInMicroCentsInMicroCents int64, err error) {
	return b.debitAccountResult, b.debitAccountError
}

func (b *MockBanker) RemainingDailyBudgetInMicroCents(account int64) int64 {
	return b.remainingDailyBudgetInMicroCentsResult
}

func (b *MockBanker) SetRemainingDailyBudgetInMicroCents(account int64, amount int64, dailyBudgetExpiration time.Time) {

}

func NewMockBanker(debitAccountResult int64, debitAccountError error, remainingDailyBudgetInMicroCentsResult int64) rtb.Banker {
	b := new(MockBanker)
	b.debitAccountError = debitAccountError
	b.debitAccountResult = debitAccountResult
	b.remainingDailyBudgetInMicroCentsResult = remainingDailyBudgetInMicroCentsResult

	return b
}

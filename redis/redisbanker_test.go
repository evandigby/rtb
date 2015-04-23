package redis

import (
	"testing"
	"time"
)

// Test setting the remainingDailyBudgetInMicroCents
// Expected result is the remainingDailyBudgetInMicroCents should equal the remainingDailyBudgetInMicroCents it's set to
func TestSetRemainingDailyBudgetInMicroCents(t *testing.T) {
	b := NewRedisBanker(testDataAccess)

	account := int64(100)
	expectedRemainingDailyBudgetInMicroCents := int64(22)

	dailyBudgetExpiration := time.Now().UTC().AddDate(0, 0, 1)

	b.SetRemainingDailyBudgetInMicroCents(account, expectedRemainingDailyBudgetInMicroCents, dailyBudgetExpiration)

	remainingDailyBudgetInMicroCents := b.RemainingDailyBudgetInMicroCents(account)

	if expectedRemainingDailyBudgetInMicroCents != remainingDailyBudgetInMicroCents {
		t.Fail()
	}
}

// Ensure a non-existant account returns zero
// Expected result is a successful call to remainingDailyBudgetInMicroCents with a result of zero
func TestZeroRemainingDailyBudgetInMicroCentsWithNoAccount(t *testing.T) {
	b := NewRedisBanker(testDataAccess)

	account := int64(200) // Not used in other tests
	expectedRemainingDailyBudgetInMicroCents := int64(0)

	b.DeleteAccount(account) // Ensure the account isn't in the system

	remainingDailyBudgetInMicroCents := b.RemainingDailyBudgetInMicroCents(account)

	if expectedRemainingDailyBudgetInMicroCents != remainingDailyBudgetInMicroCents {
		t.Fail()
	}
}

// Test a debit when the account has enough
// Expected result is a successful debit with a remainingDailyBudgetInMicroCents of the initial remainingDailyBudgetInMicroCents minus the debit
func TestDebitAccountWithEnoughRemainingDailyBudgetInMicroCents(t *testing.T) {
	b := NewRedisBanker(testDataAccess)

	account := int64(100)

	amount := int64(32)
	dailyBudget := int64(32)
	remainingDailyBudgetInMicroCents := int64(32)

	dailyBudgetExpiration := time.Now().UTC().AddDate(0, 0, 1)

	expectedRemainingDailyBudgetInMicroCents := remainingDailyBudgetInMicroCents - amount

	b.SetRemainingDailyBudgetInMicroCents(account, remainingDailyBudgetInMicroCents, dailyBudgetExpiration)

	result, err := b.DebitAccount(account, amount, dailyBudget, dailyBudgetExpiration)

	if err != nil {
		t.Fail()
	}

	if result != expectedRemainingDailyBudgetInMicroCents {
		t.Fail()
	}
}

// Test a debit when the account does not have enough
// Expected result is a failed debit with an expected remainingDailyBudgetInMicroCents of the initial remainingDailyBudgetInMicroCents
func TestDebitAccountWithNotEnoughRemainingDailyBudgetInMicroCents(t *testing.T) {
	b := NewRedisBanker(testDataAccess)

	account := int64(100)

	amount := int64(32)
	dailyBudget := int64(10)
	remainingDailyBudgetInMicroCents := int64(10)

	dailyBudgetExpiration := time.Now().UTC().AddDate(0, 0, 1)

	expectedRemainingDailyBudgetInMicroCents := remainingDailyBudgetInMicroCents

	b.SetRemainingDailyBudgetInMicroCents(account, remainingDailyBudgetInMicroCents, dailyBudgetExpiration)

	result, err := b.DebitAccount(account, amount, dailyBudget, dailyBudgetExpiration)

	if err == nil {
		t.Fail()
	}

	if result != expectedRemainingDailyBudgetInMicroCents {
		t.Fail()
	}

}

// Test a debit when there is no initial remainingDailyBudgetInMicroCents
// Expected result is a successful debit with a final remainingDailyBudgetInMicroCents of the daily budget minus the debit
func TestDebitAccountDailyBudgetRolloverWithNoRemainingDailyBudgetInMicroCents(t *testing.T) {
	b := NewRedisBanker(testDataAccess)

	account := int64(100)

	amount := int64(32)
	dailyBudget := int64(100)

	dailyBudgetExpiration := time.Now().UTC().AddDate(0, 0, 1)

	expectedRemainingDailyBudgetInMicroCents := dailyBudget - amount

	b.DeleteAccount(account) // Ensure the account isn't in the system

	result, err := b.DebitAccount(account, amount, dailyBudget, dailyBudgetExpiration)

	if err != nil {
		t.Fail()
	}

	if result != expectedRemainingDailyBudgetInMicroCents {
		t.Fail()
	}

}

// Test the daily budget is reinstated once the remainingDailyBudgetInMicroCents exires
// Expected result is a successful debit with a final remainingDailyBudgetInMicroCents of the daily budget minus the debit
func TestDebitAccountDailyBudgetRolloverWithExpiredRemainingDailyBudgetInMicroCents(t *testing.T) {
	b := NewRedisBanker(testDataAccess)

	account := int64(100)

	amount := int64(32)
	dailyBudget := int64(100)
	remainingDailyBudgetInMicroCents := int64(10)

	expectedRemainingDailyBudgetInMicroCents := dailyBudget - amount

	initialExpiration := time.Now().UTC().AddDate(0, 0, -1)
	dailyBudgetExpiration := time.Now().UTC().AddDate(0, 0, 1)

	b.SetRemainingDailyBudgetInMicroCents(account, remainingDailyBudgetInMicroCents, initialExpiration)

	result, err := b.DebitAccount(account, amount, dailyBudget, dailyBudgetExpiration)

	if err != nil {
		t.Fail()
	}

	if result != expectedRemainingDailyBudgetInMicroCents {
		t.Fail()
	}

}

// Test the daily budget is reinstated once the remainingDailyBudgetInMicroCents exires
// Expected result is a failed debit with a final remainingDailyBudgetInMicroCents of the daily budget
func TestDebitAccountWithExpiredRemainingDailyBudgetInMicroCentsAndNotEnoughDaily(t *testing.T) {
	b := NewRedisBanker(testDataAccess)

	account := int64(100)

	amount := int64(32)
	dailyBudget := int64(20)
	remainingDailyBudgetInMicroCents := int64(10)

	expectedRemainingDailyBudgetInMicroCents := dailyBudget

	initialExpiration := time.Now().UTC().AddDate(0, 0, -1)
	dailyBudgetExpiration := time.Now().UTC().AddDate(0, 0, 1)

	b.SetRemainingDailyBudgetInMicroCents(account, remainingDailyBudgetInMicroCents, initialExpiration)

	result, err := b.DebitAccount(account, amount, dailyBudget, dailyBudgetExpiration)

	if err == nil {
		t.Fail()
	}

	if result != expectedRemainingDailyBudgetInMicroCents {
		t.Fail()
	}

}

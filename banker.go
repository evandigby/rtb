package rtb

import (
	"time"
)

// Banker defines an interface for a campaign provider to track daily campaign budgets
//
// Implementations of this interface should be designed with speed in mind.
//
// Although this should be close to 100% accurate, users of this interface
// should not depended on implementations for a true transaction log and accounting purposes.
type Banker interface {
	// DebitAccount subtracts an amount from an account
	// Returns the remaining remainingDailyBudgetInMicroCents after the transaction, and an error if the transaction was unsuccessful
	DebitAccount(account int64, amount int64, dailyBudget int64, dailyBudgetExpiration time.Time) (remainingDailyBudgetInMicroCentsInMicroCents int64, err error)
	// Returns the remainingDailyBudgetInMicroCents for the account, or zero for a non-existant account
	RemainingDailyBudgetInMicroCents(account int64) int64
	// Deletes an account
	DeleteAccount(account int64)
	// Sets the account's remainingDailyBudgetInMicroCents to a specific amount, expiring at a certain time
	SetRemainingDailyBudgetInMicroCents(account int64, amount int64, dailyBudgetExpiration time.Time)
}

// Defines a transaction error used by the banker
type TransactionError struct {
	// Error message
	msg string
	// Whether or not the value contained in any result is valid
	valueValid bool
}

func (e *TransactionError) Error() string { return e.msg }

func NewTransactionError(msg string, valueValid bool) error {
	e := new(TransactionError)
	e.msg = msg
	e.valueValid = valueValid

	return e
}

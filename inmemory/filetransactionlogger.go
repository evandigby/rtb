package inmemory

import (
	"fmt"
	"github.com/evandigby/rtb"
	"os"
)

type FileTransactionLogger struct {
	file *os.File
	csv  bool
}

func (l *FileTransactionLogger) ConsumerListening() (bool, error) {
	return l.file != nil, nil
}

func (l *FileTransactionLogger) LogTransaction(transaction *rtb.Transaction) error {
	var format string

	if l.csv {
		format = "%v,%v,%v,%v\n"
	} else {
		format = "Campaign Id: %v / Bid Response Id: %v / Amount In Micro Cents: %v / Timestamp: %v\n"
	}

	_, err := fmt.Fprintf(l.file, format, transaction.CampaignId, transaction.BidResponseId, transaction.AmountInMicroCents, transaction.TimestampInNanoSeconds)

	return err
}

func NewFileTransactionLogger(file *os.File, commaSeparatedValues bool) rtb.TransactionLogger {
	l := new(FileTransactionLogger)
	l.file = file
	l.csv = commaSeparatedValues
	return l
}

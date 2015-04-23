package inmemory

import (
	"encoding/json"
	"fmt"
	"github.com/evandigby/rtb"
	"os"
)

type FileBidLogger struct {
	file    *os.File
	fullbid bool
}

func (l *FileBidLogger) LogItem(logItem *rtb.BidLogItem) {
	if l.fullbid {
		js, err := json.Marshal(logItem)

		if err == nil {
			fmt.Fprintln(l.file, string(js[:]))
		}
	} else {
		responseTimeInMs := (logItem.EndTimestampInNanoseconds - logItem.StartTimestampInNanoseconds) / 1000000

		fmt.Fprint(l.file, logItem.Domain, " / Request ID: ", logItem.BidRequest.ID, " / ")
		if logItem.BidResponse != nil && logItem.BidResponse.Seatbid != nil && len(logItem.BidResponse.Seatbid) >= 1 && logItem.BidResponse.Seatbid[0].Bid != nil && len(logItem.BidResponse.Seatbid[0].Bid) >= 1 {
			for _, bid := range logItem.BidResponse.Seatbid[0].Bid {
				fmt.Fprint(l.file,
					"Campaign: ", bid.Cid,
					" / Bid: $", bid.Price,
					" / Remaining Daily Budget: $", rtb.MicroCentsToDollarsRounded(logItem.RemainingDailyBudgetsInMicroCents[bid.Cid], 5),
					" / Response Time: ", responseTimeInMs, "ms")
			}
		} else {
			fmt.Fprint(l.file, "No Bid.")
		}

		fmt.Fprintf(l.file, "\n")
	}
}

func NewFileBidLogger(file *os.File, fullbid bool) rtb.BidLogProducer {
	l := new(FileBidLogger)
	l.file = file
	l.fullbid = fullbid
	return l
}

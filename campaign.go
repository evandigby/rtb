package rtb

// TargetType defines the type of targeting
type TargetType int

// Not using IOTA as this maps to historical logs

const (
	// Placement defines a target based on the app or site name
	Placement TargetType = 1
	// CreativeSize defines a target based on the size of the creative the impression is looking for
	CreativeSize = 2
	// Country defines a target based on the country of the user the request came form
	Country = 3
	// OS defines the OS of the device the request came form
	OS = 4
)

// Target defines the type and value of a specific targeting
type Target struct {
	Type  TargetType
	Value string
}

// Campaign defines the structure of a campaign and what requests it is targeting
type Campaign interface {
	// Id defines the unique identifier of this campaign
	Id() int64
	// BidCpmInMicroCents defines the CPM this campaign is willing to bid for a matching impression
	BidCpmInMicroCents() int64
	// DailyBudgetInMicroCents defines the daily budget for this campaign
	DailyBudgetInMicroCents() int64
	// Targets define the type of request this campaign is targeting
	Targets() *map[TargetType]string
}

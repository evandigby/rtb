package rtb

type Bid struct {
	Adid    string                 `json:"adid,omitempty"`
	Adm     string                 `json:"adm,omitempty"`
	Adomain []string               `json:"adomain,omitempty"`
	Attr    []float64              `json:"attr,omitempty"`
	Cid     string                 `json:"cid,omitempty"`
	Crid    string                 `json:"crid,omitempty"`
	Crtype  string                 `json:"crtype,omitempty"`
	Ext     map[string]interface{} `json:"ext,omitempty"`
	ID      string                 `json:"id,omitempty"`
	Impid   string                 `json:"impid,omitempty"`
	Iurl    string                 `json:"iurl,omitempty"`
	Nurl    string                 `json:"nur,omitempty"`
	Price   float64                `json:"price,omitempty"`
}

type Seatbid struct {
	Bid  []Bid  `json:"bid,omitempty"`
	Seat string `json:"seat,omitempty"`
}

// BidRequest defines a bid response object and how to parse it from JSON
type BidResponse struct {
	Bidid   string    `json:"bidid,omitempty"`
	Cur     string    `json:"cur,omitempty"`
	ID      string    `json:"id,omitempty"`
	Seatbid []Seatbid `json:"seatbid,omitempty"`
}

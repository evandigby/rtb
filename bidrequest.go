package rtb

import (
	"fmt"
)

type Publisher struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type App struct {
	Bundle    string     `json:"bundle,omitempty"`
	Cat       []string   `json:"cat,omitempty"`
	ID        string     `json:"id,omitempty"`
	Name      string     `json:"name,omitempty"`
	Publisher *Publisher `json:"publisher,omitempty"`
	Storeurl  string     `json:"storeurl,omitempty"`
	Ver       string     `json:"ver,omitempty"`
}

type Geo struct {
	City    string  `json:"city,omitempty"`
	Country string  `json:"country,omitempty"`
	Lat     float64 `json:"lat,omitempty"`
	Lon     float64 `json:"lon,omitempty"`
	Metro   string  `json:"metro,omitempty"`
	Region  string  `json:"region,omitempty"`
	Zip     string  `json:"zip,omitempty"`
}

type Device struct {
	Carrier        string      `json:"carrier,omitempty"`
	Connectiontype float64     `json:"connectiontype,omitempty"`
	Devicetype     float64     `json:"devicetype,omitempty"`
	Dnt            float64     `json:"dnt,omitempty"`
	Dpidmd5        string      `json:"dpidmd5,omitempty"`
	Dpidsha1       string      `json:"dpidsha1,omitempty"`
	Ext            interface{} `json:"ext,omitempty"`
	Geo            *Geo        `json:"geo,omitempty"`
	Ip             string      `json:"ip,omitempty"`
	Js             float64     `json:"js,omitempty"`
	Language       string      `json:"language,omitempty"`
	Make           string      `json:"make,omitempty"`
	Model          string      `json:"model,omitempty"`
	Os             string      `json:"os,omitempty"`
	Osv            string      `json:"osv,omitempty"`
	Ua             string      `json:"ua,omitempty"`
}

type Deal struct {
	At       float64 `json:"at,omitempty"`
	Bidfloor float64 `json:"bidfloor,omitempty"`
	ID       string  `json:"id,omitempty"`
}

type Pmp struct {
	Deals   []Deal  `json:"deals,omitempty"`
	Private float64 `json:"private,omitempty"`
}

type Mraid struct {
	Functions []string `json:"functions,omitempty"`
	Version   string   `json:"version,omitempty"`
}

type Video struct {
	Linearity   float64  `json:"linearity,omitempty"`
	Maxduration float64  `json:"maxduration,omitempty"`
	Minduration float64  `json:"minduration,omitempty"`
	Type        []string `json:"type,omitempty"`
}

type Banner struct {
	Api   []float64   `json:"api,omitempty"`
	Battr []float64   `json:"battr,omitempty"`
	Btype []float64   `json:"btype,omitempty"`
	Ext   interface{} `json:"ext,omitempty"`
	H     int32       `json:"h,omitempty"`
	Pos   int32       `json:"pos,omitempty"`
	W     int32       `json:"w,omitempty"`
}

type Imp struct {
	Banner            *Banner `json:"banner,omitempty"`
	Bidfloor          float64 `json:"bidfloor,omitempty"`
	Displaymanager    string  `json:"displaymanager,omitempty"`
	Displaymanagerver string  `json:"displaymanagerver,omitempty"`
	ID                string  `json:"id,omitempty"`
	Instl             float64 `json:"instl,omitempty"`
	Tagid             string  `json:"tagid,omitempty"`
}

type Site struct {
	Cat       []string   `json:"cat,omitempty"`
	Domain    string     `json:"domain,omitempty"`
	ID        string     `json:"id,omitempty"`
	Name      string     `json:"name,omitempty"`
	Publisher *Publisher `json:"publisher,omitempty"`
}

type Segment struct {
	ID    string `json:"id,omitempty"`
	Value string `json:"value,omitempty"`
}

type UserData struct {
	ID      string    `json:"id,omitempty"`
	Name    string    `json:"name,omitempty"`
	Segment []Segment `json:"segment,omitempty"`
}

type User struct {
	Data     []UserData `json:"data,omitempty"`
	Gender   string     `json:"gender,omitempty"`
	Keywords string     `json:"keywords,omitempty"`
	Yob      string     `json:"yob,omitempty"`
}

// BidRequest defines a bid request object and how to parse it from JSON
type BidRequest struct {
	App    *App        `json:"app,omitempty"`
	At     float64     `json:"at,omitempty"`
	Badv   []string    `json:"badv,omitempty"`
	Bcat   []string    `json:"bcat,omitempty"`
	Device *Device     `json:"device,omitempty"`
	Ext    interface{} `json:"ext,omitempty"`
	ID     string      `json:"id,omitempty"`
	Imp    []Imp       `json:"imp,omitempty"`
	Site   *Site       `json:"site,omitempty"`
	User   *User       `json:"user,omitempty"`
}

// Targeting creates a list of targets from a bid request
func (r *BidRequest) Targeting() []Target {
	targets := make([]Target, 0, 3)

	if r.Device != nil {
		if r.Device.Geo != nil && r.Device.Geo.Country != "" {
			country := Target{Type: Country, Value: r.Device.Geo.Country}
			targets = append(targets, country)
		}

		if r.Device.Os != "" {
			os := Target{Type: OS, Value: r.Device.Os}
			targets = append(targets, os)
		}
	}

	if r.App != nil && r.App.Name != "" {
		placement := Target{Type: Placement, Value: r.App.Name}
		targets = append(targets, placement)
	} else if r.Site != nil && r.Site.Name != "" {
		placement := Target{Type: Placement, Value: r.Site.Name}
		targets = append(targets, placement)
	}

	return targets
}

// Targeting creates a list of targets from a specific impression
func (i *Imp) Targeting() []Target {
	targets := make([]Target, 0, 1)

	if i.Banner != nil && i.Banner.W != 0 && i.Banner.H != 0 {
		creative := Target{Type: CreativeSize, Value: fmt.Sprintf("%dx%d", i.Banner.W, i.Banner.H)}
		targets = append(targets, creative)
	}

	return targets
}

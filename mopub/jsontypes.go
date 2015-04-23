package mopub

import (
	"encoding/json"
	"github.com/evandigby/rtb"
)

type MoPubRequestDeviceExt struct {
	Idfa string `json:"idfa"`
}
type MoPubRequestExt struct {
	Pmp rtb.Pmp `json:"pmp"`
}

type MoPubRequestBannerExt struct {
	Mraid              []rtb.Mraid `json:"mraid"`
	Nativebrowserclick float64     `json:"nativebrowserclick"`
	Video              rtb.Video   `json:"video"`
}

type MoPubResponseExtDataSegment struct {
	ID string `json:"id"`
}

type MoPubResponseExtData struct {
	ID      string                        `json:"id"`
	Name    string                        `json:"name"`
	Segment []MoPubResponseExtDataSegment `json:"segment"`
}

type MoPubResponseExtVideo struct {
	Duration  float64 `json:"duration"`
	Linearity float64 `json:"linearity"`
	Type      string  `json:"type"`
}

type MoPubResponseExt struct {
	Data   []MoPubResponseExtData `json:"data"`
	DealID string                 `json:"deal_id"`
	Video  MoPubResponseExtVideo  `json:"video"`
}

package entities

import (
	"encoding/json"
	"gorm.io/gorm"
)

type Quotation struct {
	Bid       string `json:"bid"`
	Code      string `json:"code"`
	CodeIn    string `json:"codein"`
	Name      string `json:"name"`
	High      string `json:"high"`
	Low       string `json:"low"`
	VarBid    string `json:"varBid"`
	PctChange string `json:"pctChange"`
	Ask       string `json:"ask"`
	Timestamp string `json:"timestamp"`
	gorm.Model
}

func (q Quotation) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Bid string `json:"bid"`
	}{
		Bid: q.Bid,
	})
}

package models

type ResponseAPI struct {
	Timestamp int64 `json:"timestamp"`
	Asks      []Ask `json:"asks"`
	Bids      []Bid `json:"bids"`
}

type ResponseService struct {
	Timestamp int64
	Ask       Ask
	Bid       Bid
}

type Ask struct {
	Price  string `json:"price"`
	Volume string `json:"volume"`
	Amount string `json:"amount"`
	Factor string `json:"factor"`
	Type   string `json:"type"`
}

type Bid struct {
	Price  string `json:"price"`
	Volume string `json:"volume"`
	Amount string `json:"amount"`
	Factor string `json:"factor"`
	Type   string `json:"type"`
}

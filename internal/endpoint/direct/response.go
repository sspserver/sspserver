package direct

type debugResponse struct {
	ID                uint64 `json:"id,omitempty"`
	ZoneID            uint64 `json:"zone_id,omitempty"`
	AuctionID         string `json:"auction_id,omitempty"`
	ImpressionID      string `json:"impression_id,omitempty"`
	IsAlternativeLink bool   `json:"is_alternative_link,omitempty"`
	Link              string `json:"link,omitempty"`
	Superfailover     string `json:"superfailover,omitempty"`
	Error             error  `json:"error,omitempty"`
	IsEmpty           bool   `json:"is_empty,omitempty"`
}

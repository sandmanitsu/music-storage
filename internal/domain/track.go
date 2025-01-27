package domain

type Track struct {
	ID          int    `json:"id"`
	GroupName   string `json:"group_name"`
	Song        string `json:"song"`
	Text        string `json:"text"`
	RealiseDate string `json:"realise_date"`
	Link        string `json:"link"`
}

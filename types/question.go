package types

type Question struct {
	ID           int     `json:"id"`
	Name         string  `json:"name"`
	URL          string  `json:"url"`
	Difficulty   string  `json:"difficulty"`
	LastReviewed *string `json:"last_reviewed"`
	SRScore      int     `json:"sr_score"`
	Attempted    bool    `json:"attempted"`
}

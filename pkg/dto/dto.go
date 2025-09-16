package dto

// создать слайс dto
type GetRatesRequest struct {
	Currencies []string `json:"currencies"`
	Aggs       []string `json:"aggs, omitempty"`
	Since      string   `json:"since, omitempty"`
	Until      string   `json:"until, omitempty"`
}

type RateItem struct {
	Base   string `json:"base"`
	Price  string `json:"price"`
	Ts     string `json:"ts"`
	Source string `json:"source, omitempty"`
}

type GetRatesResponse struct {
	Rates []RateItem `json:"rates"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

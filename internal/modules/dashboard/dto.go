package dashboard

type OverviewResponse struct {
	Views         MetricCard `json:"views"`
	Messages      MetricCard `json:"messages"`
	SalesAmount   MetricCard `json:"sales_amount"`
	PurchaseCount MetricCard `json:"purchase_count"`
}

type MetricCard struct {
	Label  string  `json:"label"`
	Value  float64 `json:"value"`
	Unit   string  `json:"unit,omitempty"`
	Delta  float64 `json:"delta"`
	Trend  string  `json:"trend"`
	Source string  `json:"source"`
}

type SalesTrendQuery struct {
	Range string `form:"range" binding:"omitempty,oneof=day week month" example:"day"`
}

type SalesTrendResponse struct {
	Range  string       `json:"range"`
	Points []TrendPoint `json:"points"`
	Total  float64      `json:"total"`
	Unit   string       `json:"unit"`
	Source string       `json:"source"`
}

type TrendPoint struct {
	Label string  `json:"label"`
	Value float64 `json:"value"`
}

type VisitSourceResponse struct {
	Items  []NamedValue `json:"items"`
	Total  float64      `json:"total"`
	Source string       `json:"source"`
}

type SalesCategoryResponse struct {
	Items  []NamedValue `json:"items"`
	Total  float64      `json:"total"`
	Unit   string       `json:"unit"`
	Source string       `json:"source"`
}

type NamedValue struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
}

type OnlineUsersResponse struct {
	OnlineUsers int64  `json:"online_users"`
	Source      string `json:"source"`
}

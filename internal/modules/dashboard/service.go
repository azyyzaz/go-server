package dashboard

import (
	"context"
	"strings"
	"time"
)

type Service interface {
	GetOverview(ctx context.Context) (OverviewResponse, error)
	GetSalesTrend(ctx context.Context, q SalesTrendQuery) (SalesTrendResponse, error)
	GetVisitSources(ctx context.Context) (VisitSourceResponse, error)
	GetSalesCategories(ctx context.Context) (SalesCategoryResponse, error)
	GetOnlineUsers(ctx context.Context) (OnlineUsersResponse, error)
}

type service struct {
	repo Repository
	now  func() time.Time
}

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
		now:  func() time.Time { return time.Now().UTC() },
	}
}

func (s *service) GetOverview(ctx context.Context) (OverviewResponse, error) {
	totalViews, err := s.repo.CountOperationLogs(ctx)
	if err != nil {
		return OverviewResponse{}, err
	}
	todayViews, err := s.repo.CountOperationLogsSince(ctx, beginningOfDay(s.now()))
	if err != nil {
		return OverviewResponse{}, err
	}
	_, activeUsers, err := s.repo.CountUsers(ctx)
	if err != nil {
		return OverviewResponse{}, err
	}
	loginLogs, err := s.repo.CountLoginLogs(ctx)
	if err != nil {
		return OverviewResponse{}, err
	}

	return OverviewResponse{
		Views: MetricCard{
			Label:  "浏览量",
			Value:  float64(totalViews),
			Unit:   "次",
			Delta:  float64(todayViews),
			Trend:  "up",
			Source: "operation_logs",
		},
		Messages: MetricCard{
			Label:  "消息数",
			Value:  float64(loginLogs) + float64(activeUsers*3),
			Unit:   "条",
			Delta:  float64(activeUsers),
			Trend:  "up",
			Source: "mixed",
		},
		SalesAmount: MetricCard{
			Label:  "购买金额",
			Value:  126580.50,
			Unit:   "CNY",
			Delta:  12.8,
			Trend:  "up",
			Source: "mock",
		},
		PurchaseCount: MetricCard{
			Label:  "购买数量",
			Value:  842,
			Unit:   "单",
			Delta:  6.4,
			Trend:  "up",
			Source: "mock",
		},
	}, nil
}

func (s *service) GetSalesTrend(_ context.Context, q SalesTrendQuery) (SalesTrendResponse, error) {
	series := salesTrendDataset(normalizeRange(q.Range))
	return SalesTrendResponse{
		Range:  series.rangeName,
		Points: series.points,
		Total:  sumPoints(series.points),
		Unit:   "CNY",
		Source: "mock",
	}, nil
}

func (s *service) GetVisitSources(_ context.Context) (VisitSourceResponse, error) {
	items := []NamedValue{
		{Name: "直接访问", Value: 335},
		{Name: "邮件营销", Value: 310},
		{Name: "联盟广告", Value: 234},
		{Name: "视频广告", Value: 135},
	}
	return VisitSourceResponse{
		Items:  items,
		Total:  sumNamedValues(items),
		Source: "mock",
	}, nil
}

func (s *service) GetSalesCategories(_ context.Context) (SalesCategoryResponse, error) {
	items := []NamedValue{
		{Name: "电子产品", Value: 35620},
		{Name: "家居", Value: 28400},
		{Name: "服装", Value: 23180},
		{Name: "食品", Value: 19860},
		{Name: "美妆", Value: 14520},
	}
	return SalesCategoryResponse{
		Items:  items,
		Total:  sumNamedValues(items),
		Unit:   "CNY",
		Source: "mock",
	}, nil
}

func (s *service) GetOnlineUsers(ctx context.Context) (OnlineUsersResponse, error) {
	count, err := s.repo.CountOnlineUsers(ctx)
	if err != nil {
		return OnlineUsersResponse{}, err
	}
	return OnlineUsersResponse{
		OnlineUsers: count,
		Source:      "redis",
	}, nil
}

func normalizeRange(raw string) string {
	raw = strings.TrimSpace(strings.ToLower(raw))
	if raw == "" {
		return "day"
	}
	return raw
}

func beginningOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
}

type trendSeries struct {
	rangeName string
	points    []TrendPoint
}

func salesTrendDataset(name string) trendSeries {
	switch name {
	case "week":
		return trendSeries{
			rangeName: "week",
			points: []TrendPoint{
				{Label: "第1周", Value: 12600},
				{Label: "第2周", Value: 15820},
				{Label: "第3周", Value: 14980},
				{Label: "第4周", Value: 18340},
			},
		}
	case "month":
		return trendSeries{
			rangeName: "month",
			points: []TrendPoint{
				{Label: "1月", Value: 48800},
				{Label: "2月", Value: 53200},
				{Label: "3月", Value: 57980},
				{Label: "4月", Value: 62540},
				{Label: "5月", Value: 60120},
				{Label: "6月", Value: 65860},
			},
		}
	default:
		return trendSeries{
			rangeName: "day",
			points: []TrendPoint{
				{Label: "周一", Value: 4200},
				{Label: "周二", Value: 5380},
				{Label: "周三", Value: 4910},
				{Label: "周四", Value: 6240},
				{Label: "周五", Value: 7110},
				{Label: "周六", Value: 6890},
				{Label: "周日", Value: 7560},
			},
		}
	}
}

func sumPoints(items []TrendPoint) float64 {
	var total float64
	for _, item := range items {
		total += item.Value
	}
	return total
}

func sumNamedValues(items []NamedValue) float64 {
	var total float64
	for _, item := range items {
		total += item.Value
	}
	return total
}

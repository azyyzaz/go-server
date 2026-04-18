package dashboard

import (
	"context"
	"testing"
	"time"
)

type stubRepository struct {
	totalUsers          int64
	activeUsers         int64
	operationLogs       int64
	todayOperationLogs  int64
	loginLogs           int64
	onlineUsers         int64
}

func (r stubRepository) CountUsers(context.Context) (int64, int64, error) {
	return r.totalUsers, r.activeUsers, nil
}

func (r stubRepository) CountOperationLogs(context.Context) (int64, error) {
	return r.operationLogs, nil
}

func (r stubRepository) CountLoginLogs(context.Context) (int64, error) {
	return r.loginLogs, nil
}

func (r stubRepository) CountOperationLogsSince(context.Context, time.Time) (int64, error) {
	return r.todayOperationLogs, nil
}

func (r stubRepository) CountOnlineUsers(context.Context) (int64, error) {
	return r.onlineUsers, nil
}

func TestGetOverviewBuildsMetricCards(t *testing.T) {
	svc := &service{
		repo: stubRepository{
			totalUsers:         18,
			activeUsers:        7,
			operationLogs:      128,
			todayOperationLogs: 12,
			loginLogs:          20,
		},
		now: func() time.Time {
			return time.Date(2026, 4, 18, 10, 0, 0, 0, time.UTC)
		},
	}

	result, err := svc.GetOverview(context.Background())
	if err != nil {
		t.Fatalf("get overview: %v", err)
	}

	if result.Views.Value != 128 {
		t.Fatalf("expected views 128, got %v", result.Views.Value)
	}
	if result.Views.Delta != 12 {
		t.Fatalf("expected today delta 12, got %v", result.Views.Delta)
	}
	if result.Messages.Value != 41 {
		t.Fatalf("expected messages 41, got %v", result.Messages.Value)
	}
	if result.SalesAmount.Source != "mock" {
		t.Fatalf("expected mock sales amount source, got %q", result.SalesAmount.Source)
	}
	if result.PurchaseCount.Value != 842 {
		t.Fatalf("expected purchase count 842, got %v", result.PurchaseCount.Value)
	}
}

func TestGetSalesTrendSupportsRangeAggregation(t *testing.T) {
	svc := NewService(NewInMemoryRepository())

	day, err := svc.GetSalesTrend(context.Background(), SalesTrendQuery{})
	if err != nil {
		t.Fatalf("get day trend: %v", err)
	}
	if day.Range != "day" {
		t.Fatalf("expected default day range, got %q", day.Range)
	}
	if len(day.Points) != 7 {
		t.Fatalf("expected 7 day points, got %d", len(day.Points))
	}
	if day.Total <= 0 {
		t.Fatalf("expected positive total, got %v", day.Total)
	}

	month, err := svc.GetSalesTrend(context.Background(), SalesTrendQuery{Range: "month"})
	if err != nil {
		t.Fatalf("get month trend: %v", err)
	}
	if month.Range != "month" {
		t.Fatalf("expected month range, got %q", month.Range)
	}
	if len(month.Points) != 6 {
		t.Fatalf("expected 6 month points, got %d", len(month.Points))
	}
}

func TestGetVisitSourcesAndSalesCategoriesProvideChartData(t *testing.T) {
	svc := NewService(NewInMemoryRepository())

	visitSources, err := svc.GetVisitSources(context.Background())
	if err != nil {
		t.Fatalf("get visit sources: %v", err)
	}
	if len(visitSources.Items) != 4 {
		t.Fatalf("expected 4 visit sources, got %d", len(visitSources.Items))
	}
	if visitSources.Total != 1014 {
		t.Fatalf("expected visit source total 1014, got %v", visitSources.Total)
	}

	categories, err := svc.GetSalesCategories(context.Background())
	if err != nil {
		t.Fatalf("get sales categories: %v", err)
	}
	if len(categories.Items) != 5 {
		t.Fatalf("expected 5 categories, got %d", len(categories.Items))
	}
	if categories.Total != 121580 {
		t.Fatalf("expected category total 121580, got %v", categories.Total)
	}
}

func TestGetOnlineUsersReadsRepositoryCount(t *testing.T) {
	svc := NewService(stubRepository{onlineUsers: 9})

	result, err := svc.GetOnlineUsers(context.Background())
	if err != nil {
		t.Fatalf("get online users: %v", err)
	}
	if result.OnlineUsers != 9 {
		t.Fatalf("expected 9 online users, got %d", result.OnlineUsers)
	}
	if result.Source != "redis" {
		t.Fatalf("expected redis source, got %q", result.Source)
	}
}

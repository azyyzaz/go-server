package audit

import (
	"context"
	"testing"
	"time"
)

func TestSummarizeRequestBodyRedactsSensitiveFields(t *testing.T) {
	summary := SummarizeRequestBody([]byte(`{"username":"alice","password":"secret","access_token":"abc"}`))
	if summary == "" {
		t.Fatal("expected non-empty summary")
	}
	if contains(summary, "secret") || contains(summary, "abc") {
		t.Fatalf("expected sensitive fields to be redacted, got %s", summary)
	}
}

func TestListOperationLogsDefaultsPagination(t *testing.T) {
	svc := NewService(NewInMemoryRepository(), "未知")
	if err := svc.RecordOperation(context.Background(), OperationLogEntry{
		RequestID:  "req-1",
		Username:   "alice",
		Method:     "GET",
		Path:       "/api/v1/system/users",
		StatusCode: 200,
	}); err != nil {
		t.Fatalf("record operation: %v", err)
	}

	result, err := svc.ListOperationLogs(context.Background(), ListOperationLogsQuery{})
	if err != nil {
		t.Fatalf("list operation logs: %v", err)
	}
	if result.Page != 1 || result.PageSize != 10 || result.Total != 1 {
		t.Fatalf("unexpected pagination result: %#v", result)
	}
}

func TestCleanupExpiredRemovesOldLogs(t *testing.T) {
	repo := NewInMemoryRepository()
	svc := NewService(repo, "未知")

	memRepo := repo.(*inMemoryRepository)
	oldTime := time.Now().UTC().AddDate(0, 0, -40)
	memRepo.operationLogs = append(memRepo.operationLogs, OperationLog{ID: 1, CreatedAt: oldTime})
	memRepo.loginLogs = append(memRepo.loginLogs, LoginLog{ID: 1, CreatedAt: oldTime})

	if err := svc.CleanupExpired(context.Background(), time.Now().UTC().AddDate(0, 0, -30)); err != nil {
		t.Fatalf("cleanup expired: %v", err)
	}

	ops, err := svc.ListOperationLogs(context.Background(), ListOperationLogsQuery{})
	if err != nil {
		t.Fatalf("list operation logs: %v", err)
	}
	if ops.Total != 0 {
		t.Fatalf("expected 0 operation logs after cleanup, got %d", ops.Total)
	}

	logins, err := svc.ListLoginLogs(context.Background(), ListLoginLogsQuery{})
	if err != nil {
		t.Fatalf("list login logs: %v", err)
	}
	if logins.Total != 0 {
		t.Fatalf("expected 0 login logs after cleanup, got %d", logins.Total)
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(sub) > 0 && (stringIndex(s, sub) >= 0))
}

func stringIndex(s, sub string) int {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}

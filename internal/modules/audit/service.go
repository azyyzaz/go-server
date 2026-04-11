package audit

import (
	"context"
	"strings"
	"time"
)

type Service interface {
	RecordOperation(ctx context.Context, entry OperationLogEntry) error
	RecordLogin(ctx context.Context, entry LoginLogEntry) error
	ListOperationLogs(ctx context.Context, q ListOperationLogsQuery) (OperationLogPageResult, error)
	ListLoginLogs(ctx context.Context, q ListLoginLogsQuery) (LoginLogPageResult, error)
	CleanupExpired(ctx context.Context, before time.Time) error
}

type service struct {
	repo           Repository
	regionFallback string
}

func NewService(repo Repository, regionFallback string) Service {
	return &service{repo: repo, regionFallback: strings.TrimSpace(regionFallback)}
}

func (s *service) RecordOperation(ctx context.Context, entry OperationLogEntry) error {
	module, action := describeOperation(entry.Method, entry.Path)
	item := OperationLog{
		RequestID:    strings.TrimSpace(entry.RequestID),
		UserID:       entry.UserID,
		Username:     strings.TrimSpace(entry.Username),
		Method:       strings.ToUpper(strings.TrimSpace(entry.Method)),
		Path:         strings.TrimSpace(entry.Path),
		Module:       module,
		Action:       action,
		IP:           strings.TrimSpace(entry.IP),
		StatusCode:   entry.StatusCode,
		Result:       resultFromStatus(entry.StatusCode),
		RequestBody:  strings.TrimSpace(entry.RequestBody),
		ErrorMessage: truncate(strings.TrimSpace(entry.ErrorMessage), 255),
		LatencyMS:    entry.LatencyMS,
		CreatedAt:    time.Now().UTC(),
	}
	return s.repo.CreateOperationLog(ctx, item)
}

func (s *service) RecordLogin(ctx context.Context, entry LoginLogEntry) error {
	region := strings.TrimSpace(entry.Region)
	if region == "" {
		region = s.regionFallback
	}
	item := LoginLog{
		RequestID:  strings.TrimSpace(entry.RequestID),
		UserID:     entry.UserID,
		Username:   strings.TrimSpace(entry.Username),
		IP:         strings.TrimSpace(entry.IP),
		Region:     region,
		UserAgent:  truncate(strings.TrimSpace(entry.UserAgent), 512),
		Success:    entry.Success,
		FailReason: truncate(strings.TrimSpace(entry.FailReason), 255),
		CreatedAt:  time.Now().UTC(),
	}
	return s.repo.CreateLoginLog(ctx, item)
}

func (s *service) ListOperationLogs(ctx context.Context, q ListOperationLogsQuery) (OperationLogPageResult, error) {
	normalizeOperationQuery(&q)
	start, end, err := parseTimeRange(q.StartTime, q.EndTime)
	if err != nil {
		return OperationLogPageResult{}, err
	}
	items, total, err := s.repo.ListOperationLogs(ctx, q, start, end)
	if err != nil {
		return OperationLogPageResult{}, err
	}

	result := make([]OperationLogResponse, 0, len(items))
	for _, item := range items {
		result = append(result, OperationLogResponse{
			ID:           item.ID,
			RequestID:    item.RequestID,
			UserID:       item.UserID,
			Username:     item.Username,
			Method:       item.Method,
			Path:         item.Path,
			Module:       item.Module,
			Action:       item.Action,
			IP:           item.IP,
			StatusCode:   item.StatusCode,
			Result:       item.Result,
			RequestBody:  item.RequestBody,
			ErrorMessage: item.ErrorMessage,
			LatencyMS:    item.LatencyMS,
			CreatedAt:    item.CreatedAt,
		})
	}

	return OperationLogPageResult{Items: result, Total: total, Page: q.Page, PageSize: q.PageSize}, nil
}

func (s *service) ListLoginLogs(ctx context.Context, q ListLoginLogsQuery) (LoginLogPageResult, error) {
	normalizeLoginQuery(&q)
	start, end, err := parseTimeRange(q.StartTime, q.EndTime)
	if err != nil {
		return LoginLogPageResult{}, err
	}
	items, total, err := s.repo.ListLoginLogs(ctx, q, start, end)
	if err != nil {
		return LoginLogPageResult{}, err
	}

	result := make([]LoginLogResponse, 0, len(items))
	for _, item := range items {
		result = append(result, LoginLogResponse{
			ID:         item.ID,
			RequestID:  item.RequestID,
			UserID:     item.UserID,
			Username:   item.Username,
			IP:         item.IP,
			Region:     item.Region,
			UserAgent:  item.UserAgent,
			Success:    item.Success,
			FailReason: item.FailReason,
			CreatedAt:  item.CreatedAt,
		})
	}

	return LoginLogPageResult{Items: result, Total: total, Page: q.Page, PageSize: q.PageSize}, nil
}

func (s *service) CleanupExpired(ctx context.Context, before time.Time) error {
	if _, err := s.repo.DeleteOperationLogsBefore(ctx, before); err != nil {
		return err
	}
	if _, err := s.repo.DeleteLoginLogsBefore(ctx, before); err != nil {
		return err
	}
	return nil
}

func normalizeOperationQuery(q *ListOperationLogsQuery) {
	if q.Page <= 0 {
		q.Page = 1
	}
	if q.PageSize <= 0 {
		q.PageSize = 10
	}
}

func normalizeLoginQuery(q *ListLoginLogsQuery) {
	if q.Page <= 0 {
		q.Page = 1
	}
	if q.PageSize <= 0 {
		q.PageSize = 10
	}
}

func parseTimeRange(startRaw, endRaw string) (*time.Time, *time.Time, error) {
	parse := func(raw string) (*time.Time, error) {
		raw = strings.TrimSpace(raw)
		if raw == "" {
			return nil, nil
		}
		t, err := time.Parse(time.RFC3339, raw)
		if err != nil {
			return nil, ErrAuditInvalidTimeRange
		}
		utc := t.UTC()
		return &utc, nil
	}

	start, err := parse(startRaw)
	if err != nil {
		return nil, nil, err
	}
	end, err := parse(endRaw)
	if err != nil {
		return nil, nil, err
	}
	if start != nil && end != nil && start.After(*end) {
		return nil, nil, ErrAuditInvalidTimeRange
	}
	return start, end, nil
}

func describeOperation(method, path string) (string, string) {
	method = strings.ToUpper(strings.TrimSpace(method))
	path = strings.Trim(strings.TrimSpace(path), "/")
	segments := strings.Split(path, "/")
	if len(segments) >= 4 && segments[0] == "api" && segments[1] == "v1" {
		return segments[2] + "/" + segments[3], strings.ToLower(method)
	}
	if len(segments) >= 1 && segments[0] != "" {
		return segments[0], strings.ToLower(method)
	}
	return "unknown", strings.ToLower(method)
}

func resultFromStatus(status int) string {
	if status >= 400 {
		return "failed"
	}
	return "success"
}

func truncate(s string, limit int) string {
	if limit <= 0 || len(s) <= limit {
		return s
	}
	return s[:limit]
}

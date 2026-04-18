package dashboard

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type repositoryMock struct {
	mock.Mock
}

func (m *repositoryMock) CountUsers(ctx context.Context) (int64, int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Get(1).(int64), args.Error(2)
}

func (m *repositoryMock) CountOperationLogs(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

func (m *repositoryMock) CountLoginLogs(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

func (m *repositoryMock) CountOperationLogsSince(ctx context.Context, since time.Time) (int64, error) {
	args := m.Called(ctx, since)
	return args.Get(0).(int64), args.Error(1)
}

func (m *repositoryMock) CountOnlineUsers(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

func TestServiceGetOverviewWithTestifyMock(t *testing.T) {
	repo := &repositoryMock{}
	now := time.Date(2026, 4, 18, 15, 30, 0, 0, time.UTC)
	dayStart := time.Date(2026, 4, 18, 0, 0, 0, 0, time.UTC)

	repo.On("CountOperationLogs", mock.Anything).Return(int64(320), nil).Once()
	repo.On("CountOperationLogsSince", mock.Anything, dayStart).Return(int64(18), nil).Once()
	repo.On("CountUsers", mock.Anything).Return(int64(40), int64(9), nil).Once()
	repo.On("CountLoginLogs", mock.Anything).Return(int64(21), nil).Once()

	svc := &service{repo: repo, now: func() time.Time { return now }}

	result, err := svc.GetOverview(context.Background())
	require.NoError(t, err)
	require.Equal(t, float64(320), result.Views.Value)
	require.Equal(t, float64(18), result.Views.Delta)
	require.Equal(t, float64(48), result.Messages.Value)

	repo.AssertExpectations(t)
}

func TestServiceGetOnlineUsersWithTestifyMock(t *testing.T) {
	repo := &repositoryMock{}
	repo.On("CountOnlineUsers", mock.Anything).Return(int64(11), nil).Once()

	svc := NewService(repo)
	result, err := svc.GetOnlineUsers(context.Background())
	require.NoError(t, err)
	require.Equal(t, int64(11), result.OnlineUsers)
	require.Equal(t, "redis", result.Source)

	repo.AssertExpectations(t)
}

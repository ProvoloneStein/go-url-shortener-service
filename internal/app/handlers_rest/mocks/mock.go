// Code generated by MockGen. DO NOT EDIT.
// Source: handlers.go

// Package mock_handlersrest is a generated GoMock package.
package mock_handlersrest

import (
	context "context"
	reflect "reflect"

	models "github.com/ProvoloneStein/go-url-shortener-service/internal/app/models"
	gomock "github.com/golang/mock/gomock"
)

// MockService is a mock of Service interface.
type MockService struct {
	ctrl     *gomock.Controller
	recorder *MockServiceMockRecorder
}

// MockServiceMockRecorder is the mock recorder for MockService.
type MockServiceMockRecorder struct {
	mock *MockService
}

// NewMockService creates a new mock instance.
func NewMockService(ctrl *gomock.Controller) *MockService {
	mock := &MockService{ctrl: ctrl}
	mock.recorder = &MockServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockService) EXPECT() *MockServiceMockRecorder {
	return m.recorder
}

// BatchCreate mocks base method.
func (m *MockService) BatchCreate(ctx context.Context, userID string, data []models.BatchCreateRequest) ([]models.BatchCreateResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BatchCreate", ctx, userID, data)
	ret0, _ := ret[0].([]models.BatchCreateResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// BatchCreate indicates an expected call of BatchCreate.
func (mr *MockServiceMockRecorder) BatchCreate(ctx, userID, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BatchCreate", reflect.TypeOf((*MockService)(nil).BatchCreate), ctx, userID, data)
}

// CreateShortURL mocks base method.
func (m *MockService) CreateShortURL(ctx context.Context, userID, fullURL string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateShortURL", ctx, userID, fullURL)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateShortURL indicates an expected call of CreateShortURL.
func (mr *MockServiceMockRecorder) CreateShortURL(ctx, userID, fullURL interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateShortURL", reflect.TypeOf((*MockService)(nil).CreateShortURL), ctx, userID, fullURL)
}

// DeleteUserURLsBatch mocks base method.
func (m *MockService) DeleteUserURLsBatch(ctx context.Context, userID string, data []string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteUserURLsBatch", ctx, userID, data)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteUserURLsBatch indicates an expected call of DeleteUserURLsBatch.
func (mr *MockServiceMockRecorder) DeleteUserURLsBatch(ctx, userID, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteUserURLsBatch", reflect.TypeOf((*MockService)(nil).DeleteUserURLsBatch), ctx, userID, data)
}

// GenerateToken mocks base method.
func (m *MockService) GenerateToken(ctx context.Context) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateToken", ctx)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GenerateToken indicates an expected call of GenerateToken.
func (mr *MockServiceMockRecorder) GenerateToken(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateToken", reflect.TypeOf((*MockService)(nil).GenerateToken), ctx)
}

// GetFullByID mocks base method.
func (m *MockService) GetFullByID(ctx context.Context, shortURL string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFullByID", ctx, shortURL)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFullByID indicates an expected call of GetFullByID.
func (mr *MockServiceMockRecorder) GetFullByID(ctx, shortURL interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFullByID", reflect.TypeOf((*MockService)(nil).GetFullByID), ctx, shortURL)
}

// GetListByUser mocks base method.
func (m *MockService) GetListByUser(ctx context.Context, userID string) ([]models.GetURLResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetListByUser", ctx, userID)
	ret0, _ := ret[0].([]models.GetURLResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetListByUser indicates an expected call of GetListByUser.
func (mr *MockServiceMockRecorder) GetListByUser(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetListByUser", reflect.TypeOf((*MockService)(nil).GetListByUser), ctx, userID)
}

// ParseToken mocks base method.
func (m *MockService) ParseToken(accessToken string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ParseToken", accessToken)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ParseToken indicates an expected call of ParseToken.
func (mr *MockServiceMockRecorder) ParseToken(accessToken interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ParseToken", reflect.TypeOf((*MockService)(nil).ParseToken), accessToken)
}

// Ping mocks base method.
func (m *MockService) Ping() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Ping")
	ret0, _ := ret[0].(error)
	return ret0
}

// Ping indicates an expected call of Ping.
func (mr *MockServiceMockRecorder) Ping() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Ping", reflect.TypeOf((*MockService)(nil).Ping))
}

// Stats mocks base method.
func (m *MockService) Stats(ctx context.Context) (models.StatsData, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Stats", ctx)
	ret0, _ := ret[0].(models.StatsData)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Stats indicates an expected call of Stats.
func (mr *MockServiceMockRecorder) Stats(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stats", reflect.TypeOf((*MockService)(nil).Stats), ctx)
}
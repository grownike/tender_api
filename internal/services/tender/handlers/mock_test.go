// Code generated by MockGen. DO NOT EDIT.
// Source: deps.go

// Package handlers is a generated GoMock package.
package handlers

import (
	models "avito_tenders/internal/models"
	reflect "reflect"

	gin "github.com/gin-gonic/gin"
	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
	gorm "gorm.io/gorm"
)

// Mockstorage is a mock of storage interface.
type Mockstorage struct {
	ctrl     *gomock.Controller
	recorder *MockstorageMockRecorder
}

// MockstorageMockRecorder is the mock recorder for Mockstorage.
type MockstorageMockRecorder struct {
	mock *Mockstorage
}

// NewMockstorage creates a new mock instance.
func NewMockstorage(ctrl *gomock.Controller) *Mockstorage {
	mock := &Mockstorage{ctrl: ctrl}
	mock.recorder = &MockstorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Mockstorage) EXPECT() *MockstorageMockRecorder {
	return m.recorder
}

// CreateTender mocks base method.
func (m *Mockstorage) CreateTender(c *gin.Context, tender *models.Tender) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateTender", c, tender)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateTender indicates an expected call of CreateTender.
func (mr *MockstorageMockRecorder) CreateTender(c, tender interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTender", reflect.TypeOf((*Mockstorage)(nil).CreateTender), c, tender)
}

// EditTenderStatus mocks base method.
func (m *Mockstorage) EditTenderStatus(tenderId uuid.UUID, username, newStatus string) (*models.Tender, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EditTenderStatus", tenderId, username, newStatus)
	ret0, _ := ret[0].(*models.Tender)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// EditTenderStatus indicates an expected call of EditTenderStatus.
func (mr *MockstorageMockRecorder) EditTenderStatus(tenderId, username, newStatus interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EditTenderStatus", reflect.TypeOf((*Mockstorage)(nil).EditTenderStatus), tenderId, username, newStatus)
}

// GetMyTenders mocks base method.
func (m *Mockstorage) GetMyTenders(username string, limit, offset int) ([]models.Tender, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMyTenders", username, limit, offset)
	ret0, _ := ret[0].([]models.Tender)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMyTenders indicates an expected call of GetMyTenders.
func (mr *MockstorageMockRecorder) GetMyTenders(username, limit, offset interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMyTenders", reflect.TypeOf((*Mockstorage)(nil).GetMyTenders), username, limit, offset)
}

// GetTenderStatus mocks base method.
func (m *Mockstorage) GetTenderStatus(tenderId uuid.UUID, username string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTenderStatus", tenderId, username)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTenderStatus indicates an expected call of GetTenderStatus.
func (mr *MockstorageMockRecorder) GetTenderStatus(tenderId, username interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTenderStatus", reflect.TypeOf((*Mockstorage)(nil).GetTenderStatus), tenderId, username)
}

// GetTenders mocks base method.
func (m *Mockstorage) GetTenders(limit, offset int, serviceType []string) *gorm.DB {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTenders", limit, offset, serviceType)
	ret0, _ := ret[0].(*gorm.DB)
	return ret0
}

// GetTenders indicates an expected call of GetTenders.
func (mr *MockstorageMockRecorder) GetTenders(limit, offset, serviceType interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTenders", reflect.TypeOf((*Mockstorage)(nil).GetTenders), limit, offset, serviceType)
}

// RollbackTender mocks base method.
func (m *Mockstorage) RollbackTender(tenderId uuid.UUID, version int, username string) (*models.Tender, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RollbackTender", tenderId, version, username)
	ret0, _ := ret[0].(*models.Tender)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RollbackTender indicates an expected call of RollbackTender.
func (mr *MockstorageMockRecorder) RollbackTender(tenderId, version, username interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RollbackTender", reflect.TypeOf((*Mockstorage)(nil).RollbackTender), tenderId, version, username)
}

// UpdateTender mocks base method.
func (m *Mockstorage) UpdateTender(c *gin.Context, tenderId uuid.UUID, updates map[string]interface{}, username string) (*models.Tender, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateTender", c, tenderId, updates, username)
	ret0, _ := ret[0].(*models.Tender)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateTender indicates an expected call of UpdateTender.
func (mr *MockstorageMockRecorder) UpdateTender(c, tenderId, updates, username interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateTender", reflect.TypeOf((*Mockstorage)(nil).UpdateTender), c, tenderId, updates, username)
}

package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func Test_CreateTender_Error_DB(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	storageMock := NewMockstorage(ctrl)

	storageMock.EXPECT().CreateTender(gomock.Any(), gomock.Any()).Return(errors.New("some error"))

	handler := New(storageMock)

	w := httptest.NewRecorder()
	body := strings.NewReader(`{
		"name": "test_name",  
		"description": "test_tender",
		"serviceType": "Delivery",
		"organizationId": "550e8400-e29b-41d4-a716-446655440000",
		"creatorUsername": "admin"
	}`)
	req, _ := http.NewRequest(http.MethodPost, "/tender", body)
	req.Header.Set("Content-Type", "application/json")

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.CreateTender()(c)

	require.Equal(t, http.StatusInternalServerError, w.Code)

}

func Test_CreateTender_Error_Input(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	storageMock := NewMockstorage(ctrl)
	handler := New(storageMock)

	w := httptest.NewRecorder()
	body := strings.NewReader(`{
  "nameeee": "new_tender",  
  "description": "test_tender",
  "serviceType": "Delivery",
  "organizationId": "550e8400-e29b-41d4-a716-446655440000",
  "creatorUsername": "admin"
}`)
	req, _ := http.NewRequest(http.MethodPost, "/tender", body)
	req.Header.Set("Content-Type", "application/json")
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	handler.CreateTender()(c)
	require.Equal(t, http.StatusBadRequest, w.Code)
}

func Test_CreateTender_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	storageMock := NewMockstorage(ctrl)

	storageMock.EXPECT().CreateTender(gomock.Any(), gomock.Any()).Return(nil)

	handler := New(storageMock)

	w := httptest.NewRecorder()
	body := strings.NewReader(`{
  "name": "new_tender",
  "description": "test_tender",
  "serviceType": "Delivery",
  "organizationId": "550e8400-e29b-41d4-a716-446655440000",
  "creatorUsername": "admin"
}`)
	req, _ := http.NewRequest(http.MethodPost, "/tender", body)
	req.Header.Set("Content-Type", "application/json")

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.CreateTender()(c)

	require.Equal(t, http.StatusOK, w.Code)
}

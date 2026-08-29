package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHealthCheckHandler_Connected(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPing()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	handler := HealthCheckHandler(db)
	handler(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response HealthResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "ok", response.Status)
	assert.Equal(t, "healthy", response.Message)
	assert.Equal(t, "connected", response.Database)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestHealthCheckHandler_Disconnected(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPing().WillReturnError(errors.New("db connection refused"))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	handler := HealthCheckHandler(db)
	handler(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response HealthResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "ok", response.Status)
	assert.Equal(t, "healthy", response.Message)
	assert.Equal(t, "disconnected", response.Database)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestHealthCheckHandler_NilDB(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	var db *sql.DB = nil
	handler := HealthCheckHandler(db)
	handler(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response HealthResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "ok", response.Status)
	assert.Equal(t, "healthy", response.Message)
	assert.Equal(t, "disconnected", response.Database)
}

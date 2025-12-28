package tasks

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/avraam311/tasks-service/internal/api/responses"
	"github.com/avraam311/tasks-service/internal/mocks"
	"github.com/avraam311/tasks-service/internal/models"
	"github.com/avraam311/tasks-service/internal/repository/tasks"
)

func TestCreateTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	handler := New(mockService)

	tests := []struct {
		name         string
		method       string
		body         interface{}
		expectedCode int
		expectedErr  string
		serviceMock  func()
	}{
		{
			name:         "MethodNotAllowed",
			method:       http.MethodGet,
			body:         nil,
			expectedCode: http.StatusMethodNotAllowed,
			expectedErr:  responses.ErrMethodNotAllowed,
		},
		{
			name:         "InvalidJSON",
			method:       http.MethodPost,
			body:         "invalid json",
			expectedCode: http.StatusBadRequest,
			expectedErr:  responses.ErrInvalidJSON,
		},
		{
			name:   "ServiceError",
			method: http.MethodPost,
			body: models.TaskDTO{
				Header:      "Test Task",
				Description: "Test Description",
				Finished:    false,
			},
			expectedCode: http.StatusInternalServerError,
			expectedErr:  responses.ErrInternalServer,
			serviceMock: func() {
				mockService.EXPECT().CreateTask(gomock.Any(), gomock.Any()).
					Return(uint(0), assert.AnError)
			},
		},
		{
			name:   "Success",
			method: http.MethodPost,
			body: models.TaskDTO{
				Header:      "Test Task",
				Description: "Test Description",
				Finished:    false,
			},
			expectedCode: http.StatusCreated,
			expectedErr:  "",
			serviceMock: func() {
				mockService.EXPECT().CreateTask(gomock.Any(), gomock.Any()).
					Return(uint(1), nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var req *http.Request
			if tt.body != nil {
				jsonBody, _ := json.Marshal(tt.body)
				req = httptest.NewRequest(tt.method, "/todos", bytes.NewBuffer(jsonBody))
			} else {
				req = httptest.NewRequest(tt.method, "/todos", nil)
			}

			w := httptest.NewRecorder()

			if tt.serviceMock != nil {
				tt.serviceMock()
			}

			handler.CreateTask(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)

			if tt.expectedErr != "" {
				var errorResp responses.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &errorResp)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedErr, errorResp.Error.Code)
			} else {
				var successResp responses.Success
				err := json.Unmarshal(w.Body.Bytes(), &successResp)
				assert.NoError(t, err)
				assert.NotNil(t, successResp.Result)
			}
		})
	}
}

func TestGetAllTasks(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	handler := New(mockService)

	tests := []struct {
		name         string
		method       string
		expectedCode int
		expectedErr  string
		serviceMock  func()
	}{
		{
			name:         "MethodNotAllowed",
			method:       http.MethodPost,
			expectedCode: http.StatusMethodNotAllowed,
			expectedErr:  responses.ErrMethodNotAllowed,
		},
		{
			name:         "ServiceError",
			method:       http.MethodGet,
			expectedCode: http.StatusInternalServerError,
			expectedErr:  responses.ErrInternalServer,
			serviceMock: func() {
				mockService.EXPECT().GetAllTasks(gomock.Any()).
					Return(nil, assert.AnError)
			},
		},
		{
			name:         "Success",
			method:       http.MethodGet,
			expectedCode: http.StatusOK,
			expectedErr:  "",
			serviceMock: func() {
				tasks := []*models.TaskDomain{
					{
						ID:          1,
						Header:      "Test Task",
						Description: "Test Description",
						Finished:    false,
					},
				}
				mockService.EXPECT().GetAllTasks(gomock.Any()).
					Return(tasks, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/todos", nil)
			w := httptest.NewRecorder()

			if tt.serviceMock != nil {
				tt.serviceMock()
			}

			handler.GetAllTasks(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)

			if tt.expectedErr != "" {
				var errorResp responses.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &errorResp)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedErr, errorResp.Error.Code)
			} else {
				var successResp responses.Success
				err := json.Unmarshal(w.Body.Bytes(), &successResp)
				assert.NoError(t, err)
				assert.NotNil(t, successResp.Result)
			}
		})
	}
}

func TestGetTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	handler := New(mockService)

	tests := []struct {
		name         string
		method       string
		path         string
		expectedCode int
		expectedErr  string
		serviceMock  func()
	}{
		{
			name:         "MethodNotAllowed",
			method:       http.MethodPost,
			path:         "/todos/1",
			expectedCode: http.StatusMethodNotAllowed,
			expectedErr:  responses.ErrMethodNotAllowed,
		},
		{
			name:         "InvalidID",
			method:       http.MethodGet,
			path:         "/todos/invalid",
			expectedCode: http.StatusBadRequest,
			expectedErr:  responses.ErrInvalidID,
		},
		{
			name:         "TaskNotFound",
			method:       http.MethodGet,
			path:         "/todos/1",
			expectedCode: http.StatusBadRequest,
			expectedErr:  responses.ErrTaskNotFound,
			serviceMock: func() {
				mockService.EXPECT().GetTask(gomock.Any(), uint(1)).
					Return(nil, tasks.ErrTaskNotFound)
			},
		},
		{
			name:         "ServiceError",
			method:       http.MethodGet,
			path:         "/todos/1",
			expectedCode: http.StatusInternalServerError,
			expectedErr:  responses.ErrInternalServer,
			serviceMock: func() {
				mockService.EXPECT().GetTask(gomock.Any(), uint(1)).
					Return(nil, assert.AnError)
			},
		},
		{
			name:         "Success",
			method:       http.MethodGet,
			path:         "/todos/1",
			expectedCode: http.StatusOK,
			expectedErr:  "",
			serviceMock: func() {
				task := &models.TaskDomain{
					ID:          1,
					Header:      "Test Task",
					Description: "Test Description",
					Finished:    false,
				}
				mockService.EXPECT().GetTask(gomock.Any(), uint(1)).
					Return(task, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()

			if tt.serviceMock != nil {
				tt.serviceMock()
			}

			handler.GetTask(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)

			if tt.expectedErr != "" {
				var errorResp responses.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &errorResp)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedErr, errorResp.Error.Code)
			} else {
				var successResp responses.Success
				err := json.Unmarshal(w.Body.Bytes(), &successResp)
				assert.NoError(t, err)
				assert.NotNil(t, successResp.Result)
			}
		})
	}
}

func TestUpdateTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	handler := New(mockService)

	tests := []struct {
		name         string
		method       string
		path         string
		body         interface{}
		expectedCode int
		expectedErr  string
		serviceMock  func()
	}{
		{
			name:         "MethodNotAllowed",
			method:       http.MethodPost,
			path:         "/todos/1",
			body:         nil,
			expectedCode: http.StatusMethodNotAllowed,
			expectedErr:  responses.ErrMethodNotAllowed,
		},
		{
			name:         "InvalidID",
			method:       http.MethodPut,
			path:         "/todos/invalid",
			body:         nil,
			expectedCode: http.StatusBadRequest,
			expectedErr:  responses.ErrInvalidID,
		},
		{
			name:         "InvalidJSON",
			method:       http.MethodPut,
			path:         "/todos/1",
			body:         "invalid json",
			expectedCode: http.StatusBadRequest,
			expectedErr:  responses.ErrInvalidJSON,
		},
		{
			name:   "TaskNotFound",
			method: http.MethodPut,
			path:   "/todos/1",
			body: models.TaskDTO{
				Header:      "Updated Task",
				Description: "Updated Description",
				Finished:    true,
			},
			expectedCode: http.StatusBadRequest,
			expectedErr:  responses.ErrTaskNotFound,
			serviceMock: func() {
				mockService.EXPECT().UpdateTask(gomock.Any(), uint(1), gomock.Any()).
					Return(tasks.ErrTaskNotFound)
			},
		},
		{
			name:   "ServiceError",
			method: http.MethodPut,
			path:   "/todos/1",
			body: models.TaskDTO{
				Header:      "Updated Task",
				Description: "Updated Description",
				Finished:    true,
			},
			expectedCode: http.StatusInternalServerError,
			expectedErr:  responses.ErrInternalServer,
			serviceMock: func() {
				mockService.EXPECT().UpdateTask(gomock.Any(), uint(1), gomock.Any()).
					Return(assert.AnError)
			},
		},
		{
			name:   "Success",
			method: http.MethodPut,
			path:   "/todos/1",
			body: models.TaskDTO{
				Header:      "Updated Task",
				Description: "Updated Description",
				Finished:    true,
			},
			expectedCode: http.StatusCreated,
			expectedErr:  "",
			serviceMock: func() {
				mockService.EXPECT().UpdateTask(gomock.Any(), uint(1), gomock.Any()).
					Return(nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var req *http.Request
			if tt.body != nil {
				if bodyStr, ok := tt.body.(string); ok {
					req = httptest.NewRequest(tt.method, tt.path, strings.NewReader(bodyStr))
				} else {
					jsonBody, _ := json.Marshal(tt.body)
					req = httptest.NewRequest(tt.method, tt.path, bytes.NewBuffer(jsonBody))
				}
			} else {
				req = httptest.NewRequest(tt.method, tt.path, nil)
			}

			w := httptest.NewRecorder()

			if tt.serviceMock != nil {
				tt.serviceMock()
			}

			handler.UpdateTask(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)

			if tt.expectedErr != "" {
				var errorResp responses.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &errorResp)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedErr, errorResp.Error.Code)
			} else {
				var successResp responses.Success
				err := json.Unmarshal(w.Body.Bytes(), &successResp)
				assert.NoError(t, err)
				assert.Equal(t, responses.SuccessTaskUpdated, successResp.Result)
			}
		})
	}
}

func TestDeleteTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	handler := New(mockService)

	tests := []struct {
		name         string
		method       string
		path         string
		expectedCode int
		expectedErr  string
		serviceMock  func()
	}{
		{
			name:         "MethodNotAllowed",
			method:       http.MethodGet,
			path:         "/todos/1",
			expectedCode: http.StatusMethodNotAllowed,
			expectedErr:  responses.ErrMethodNotAllowed,
		},
		{
			name:         "InvalidID",
			method:       http.MethodDelete,
			path:         "/todos/invalid",
			expectedCode: http.StatusBadRequest,
			expectedErr:  responses.ErrInvalidID,
		},
		{
			name:         "TaskNotFound",
			method:       http.MethodDelete,
			path:         "/todos/1",
			expectedCode: http.StatusBadRequest,
			expectedErr:  responses.ErrTaskNotFound,
			serviceMock: func() {
				mockService.EXPECT().DeleteTask(gomock.Any(), uint(1)).
					Return(tasks.ErrTaskNotFound)
			},
		},
		{
			name:         "ServiceError",
			method:       http.MethodDelete,
			path:         "/todos/1",
			expectedCode: http.StatusInternalServerError,
			expectedErr:  responses.ErrInternalServer,
			serviceMock: func() {
				mockService.EXPECT().DeleteTask(gomock.Any(), uint(1)).
					Return(assert.AnError)
			},
		},
		{
			name:         "Success",
			method:       http.MethodDelete,
			path:         "/todos/1",
			expectedCode: http.StatusOK,
			expectedErr:  "",
			serviceMock: func() {
				mockService.EXPECT().DeleteTask(gomock.Any(), uint(1)).
					Return(nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()

			if tt.serviceMock != nil {
				tt.serviceMock()
			}

			handler.DeleteTask(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)

			if tt.expectedErr != "" {
				var errorResp responses.ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &errorResp)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedErr, errorResp.Error.Code)
			} else {
				var successResp responses.Success
				err := json.Unmarshal(w.Body.Bytes(), &successResp)
				assert.NoError(t, err)
				assert.Equal(t, responses.SuccessTaskDeleted, successResp.Result)
			}
		})
	}
}

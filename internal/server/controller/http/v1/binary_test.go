package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nextlag/keeper/internal/entity"
	"github.com/nextlag/keeper/internal/utils/errs"
)

func TestAddBinary(t *testing.T) {
	c, mockUseCase, ctrl := loadTest(t)
	defer ctrl.Finish()

	userID := uuid.New()
	expectedUser := entity.User{ID: userID}

	tests := []struct {
		name           string
		queryParams    map[string]string
		fileContent    []byte
		expectedStatus int
		expectedBody   string
		expectedError  error
	}{
		{
			name: "successful add binary",
			queryParams: map[string]string{
				"name": "test-binary",
			},
			fileContent:    []byte("file content"),
			expectedStatus: http.StatusCreated,
			expectedBody:   `{"name":"test-binary","file_name":"file.jpg","uuid":"00000000-0000-0000-0000-000000000000","meta":null}`,
			expectedError:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, userBinary, nil)
			req = req.WithContext(context.WithValue(req.Context(), currentUserKey, expectedUser))

			if len(tt.queryParams) > 0 {
				q := req.URL.Query()
				for key, value := range tt.queryParams {
					q.Add(key, value)
				}
				req.URL.RawQuery = q.Encode()
			}

			if tt.fileContent != nil {
				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)
				part, err := writer.CreateFormFile("file", "file.jpg")
				if err != nil {
					t.Fatalf("Failed to create form file: %v", err)
				}
				_, err = part.Write(tt.fileContent)
				if err != nil {
					t.Fatalf("Failed to write file content: %v", err)
				}
				writer.Close()
				req.Body = io.NopCloser(body)
				req.Header.Set("Content-Type", writer.FormDataContentType())
			}

			rr := httptest.NewRecorder()

			if tt.expectedError != nil {
				mockUseCase.EXPECT().AddBinary(gomock.Any(), gomock.Any(), gomock.Any(), expectedUser.ID).Return(tt.expectedError)
			} else {
				mockUseCase.EXPECT().AddBinary(gomock.Any(), gomock.Any(), gomock.Any(), expectedUser.ID).Return(nil)
			}

			http.HandlerFunc(c.AddBinary).ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			assert.JSONEq(t, tt.expectedBody, rr.Body.String())
		})
	}
}

func TestGetBinaries(t *testing.T) {
	c, mockUseCase, ctrl := loadTest(t)
	defer ctrl.Finish()

	expectedUser := entity.User{ID: uuid.New()}

	binaries := []entity.Binary{
		{Name: "binary1", FileName: "file1.jpg", ID: uuid.New(), Meta: nil},
		{Name: "binary2", FileName: "file2.jpg", ID: uuid.New(), Meta: nil},
	}

	tests := []struct {
		name           string
		mockReturn     interface{}
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "successful retrieval of binaries",
			mockReturn:     binaries,
			expectedStatus: http.StatusOK,
			expectedBody:   `[{"file_name":"file1.jpg","name":"binary1"},{"file_name":"file2.jpg","name":"binary2"}]`,
		},
		{
			name:           "no binaries found",
			mockReturn:     []entity.Binary{},
			expectedStatus: http.StatusNoContent,
			expectedBody:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, userBinary, nil)
			req = req.WithContext(context.WithValue(req.Context(), currentUserKey, expectedUser))

			rr := httptest.NewRecorder()

			if tt.mockReturn == errs.ErrUnexpectedError {
				mockUseCase.EXPECT().GetBinaries(gomock.Any(), gomock.Any()).Return(nil, tt.mockReturn.(error))
			} else {
				mockUseCase.EXPECT().GetBinaries(gomock.Any(), gomock.Any()).Return(tt.mockReturn.([]entity.Binary), nil)
			}

			http.HandlerFunc(c.GetBinaries).ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)

			if tt.expectedStatus == http.StatusNoContent {
				assert.Empty(t, rr.Body.String())
			} else {
				var resp []map[string]interface{}
				err := json.Unmarshal(rr.Body.Bytes(), &resp)
				if err != nil {
					t.Fatalf("Failed to unmarshal response body: %v", err)
				}

				var actualBinaries []map[string]interface{}
				for _, binary := range resp {
					filteredBinary := map[string]interface{}{
						"file_name": binary["file_name"],
						"name":      binary["name"],
					}
					actualBinaries = append(actualBinaries, filteredBinary)
				}

				var expectedBinaries []map[string]interface{}
				err = json.Unmarshal([]byte(tt.expectedBody), &expectedBinaries)
				if err != nil {
					t.Fatalf("Failed to unmarshal expected body: %v", err)
				}

				assert.ElementsMatch(t, expectedBinaries, actualBinaries)
			}
		})
	}
}

func TestDownloadBinary(t *testing.T) {
	c, mockUseCase, ctrl := loadTest(t)
	defer ctrl.Finish()

	expectedUser := entity.User{ID: uuid.New()}
	binaryUUID := uuid.New()
	dir := c.cfg.Log.ProjectPath
	file, err := os.CreateTemp(dir, "testfile-*.txt")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer func(name string) {
		if err = os.Remove(name); err != nil {
			t.Errorf("Failed to remove temporary file: %v", err)
		}
	}(file.Name())

	fileContent := "file content"
	if _, err = file.Write([]byte(fileContent)); err != nil {
		t.Fatalf("Failed to write to temporary file: %v", err)
	}
	if err = file.Close(); err != nil {
		t.Fatalf("Failed to close temporary file: %v", err)
	}

	tests := []struct {
		name           string
		getUserError   error
		getBinaryError error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "successful download binary",
			getUserError:   nil,
			getBinaryError: nil,
			expectedStatus: http.StatusOK,
			expectedBody:   fileContent,
		},
		{
			name:           "GetUserBinary error",
			getUserError:   nil,
			getBinaryError: errs.ErrUnexpectedError,
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "unexpected error\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, userBinary+binaryUUID.String(), nil)
			req = req.WithContext(context.WithValue(req.Context(), currentUserKey, expectedUser))
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", binaryUUID.String())
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
			rr := httptest.NewRecorder()
			switch tt.name {
			case "successful download binary":
				mockUseCase.EXPECT().GetUserBinary(gomock.Any(), gomock.Any(), gomock.Any()).Return(file.Name(), nil)
			case "GetUserBinary error":
				mockUseCase.EXPECT().GetUserBinary(gomock.Any(), gomock.Any(), gomock.Any()).Return("", tt.getBinaryError)
			}
			http.HandlerFunc(c.DownloadBinary).ServeHTTP(rr, req)
			assert.Equal(t, tt.expectedStatus, rr.Code)
			if tt.expectedStatus == http.StatusOK {
				assert.Equal(t, fileContent, rr.Body.String())
			}
		})
	}
}

func TestDelBinary(t *testing.T) {
	c, mockUseCase, ctrl := loadTest(t)
	defer ctrl.Finish()

	expectedUser := entity.User{ID: uuid.New()}
	binaryUUID := uuid.New()

	dir := c.cfg.Log.ProjectPath
	file, err := os.CreateTemp(dir, "testfile-*.txt")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	filePath := file.Name()
	fmt.Println(filePath)
	fileContent := "file content"
	if _, err = file.Write([]byte(fileContent)); err != nil {
		t.Fatalf("Failed to write to temporary file: %v", err)
	}
	if err = file.Close(); err != nil {
		t.Fatalf("Failed to close temporary file: %v", err)
	}

	defer func() {
		if err = os.Remove(filePath); err != nil {
			t.Errorf("Failed to remove temporary file: %v", err)
		}
	}()

	tests := []struct {
		name           string
		getUserError   error
		delBinaryError error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "successful delete binary",
			getUserError:   nil,
			delBinaryError: nil,
			expectedStatus: http.StatusAccepted,
			expectedBody:   `{"status":"delete accepted"}`,
		},
		{
			name:           "DelUserBinary error",
			getUserError:   nil,
			delBinaryError: errs.ErrUnexpectedError,
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"some unexpected error"}` + "\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodDelete, userBinary+binaryUUID.String(), nil)
			req = req.WithContext(context.WithValue(req.Context(), currentUserKey, expectedUser))
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", binaryUUID.String())
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			rr := httptest.NewRecorder()

			switch tt.name {
			case "successful delete binary":
				mockUseCase.EXPECT().DelUserBinary(gomock.Any(), gomock.Any(), binaryUUID).Return(nil)
			case "DelUserBinary error":
				mockUseCase.EXPECT().DelUserBinary(gomock.Any(), gomock.Any(), binaryUUID).Return(tt.delBinaryError)
			}

			http.HandlerFunc(c.DelBinary).ServeHTTP(rr, req)
			assert.Equal(t, tt.expectedStatus, rr.Code)
			assert.Equal(t, tt.expectedBody, rr.Body.String())
		})
	}
}

func TestAddBinaryMeta(t *testing.T) {
	c, mockUseCase, ctrl := loadTest(t)
	defer ctrl.Finish()

	r := chi.NewRouter()
	r.Post(userBinaryAddMeta, c.AddBinaryMeta)

	tests := []struct {
		name           string
		binaryUUID     string
		payload        interface{}
		expectedStatus int
		expectedBody   string
	}{
		{
			name:       "Valid Request",
			binaryUUID: "89dacc37-e9cb-4e9a-833b-7b8c0062b449",
			payload: []entity.Meta{
				{ID: uuid.MustParse("8a585ca3-a6a1-484e-b941-62e058ee5efa"), Name: "test", Value: "value"},
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   `[{"uuid":"8a585ca3-a6a1-484e-b941-62e058ee5efa","name":"test","value":"value"}]`,
		},
		{
			name:           "Invalid UUID",
			binaryUUID:     "invalid-uuid",
			payload:        []entity.Meta{{Name: "test", Value: "value"}},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"invalid UUID length: 12"}` + "\n",
		},
		{
			name:           "Error in UseCase",
			binaryUUID:     "89dacc37-e9cb-4e9a-833b-7b8c0062b449",
			payload:        []entity.Meta{{Name: "test", Value: "value"}},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"internal error"}` + "\n",
		},
		{
			name:           "Invalid JSON Payload",
			binaryUUID:     "89dacc37-e9cb-4e9a-833b-7b8c0062b449",
			payload:        "invalid json",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"invalid character 'i' looking for beginning of value"}` + "\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.name {
			case "Valid Request":
				mockUseCase.EXPECT().AddBinaryMeta(
					gomock.Any(), // context.Context
					gomock.Any(), // *entity.User
					uuid.MustParse(tt.binaryUUID),
					gomock.Any(), // []entity.Meta
				).Return(&entity.Binary{Meta: []entity.Meta{
					{ID: uuid.MustParse("8a585ca3-a6a1-484e-b941-62e058ee5efa"), Name: "test", Value: "value"},
				}}, nil).Times(1)

			case "Error in UseCase":
				mockUseCase.EXPECT().AddBinaryMeta(
					gomock.Any(),
					gomock.Any(),
					uuid.MustParse(tt.binaryUUID),
					gomock.Any(),
				).Return(nil, errors.New("internal error")).Times(1)

			case "Invalid UUID", "Invalid JSON Payload":
			default:
				t.Fatalf("Unknown test case: %s", tt.name)
			}

			var payloadBuffer *bytes.Buffer
			if jsonPayload, ok := tt.payload.(string); ok && jsonPayload == "invalid json" {
				payloadBuffer = bytes.NewBufferString(jsonPayload)
			} else {
				payloadData, err := json.Marshal(tt.payload)
				require.NoError(t, err)
				payloadBuffer = bytes.NewBuffer(payloadData)
			}

			req, err := http.NewRequest(http.MethodPost, "/user/binary/"+tt.binaryUUID+"/meta", payloadBuffer)
			require.NoError(t, err)

			ctx := context.WithValue(context.Background(), currentUserKey, entity.User{
				ID:    uuid.MustParse("da7400c9-0438-4489-983c-0ec9d4a77bd3"),
				Email: "example@mail.ru",
			})
			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Logf("Response Body: %s", rr.Body.String())
			}

			assert.Equal(t, tt.expectedStatus, rr.Code)
			if tt.expectedStatus == http.StatusCreated {
				assert.JSONEq(t, tt.expectedBody, rr.Body.String())
			} else {
				assert.Equal(t, tt.expectedBody, rr.Body.String())
			}
		})
	}
}

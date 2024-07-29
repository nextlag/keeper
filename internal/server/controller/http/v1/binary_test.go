package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

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
			req := httptest.NewRequest(http.MethodPost, userBinaryAddPath, nil)
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
			req := httptest.NewRequest(http.MethodGet, userBinaryGetPath, nil)
			req = req.WithContext(context.WithValue(req.Context(), currentUserKey, expectedUser))

			rr := httptest.NewRecorder()

			if tt.mockReturn == errs.ErrUnexpectedError {
				mockUseCase.EXPECT().GetBinaries(gomock.Any(), expectedUser).Return(nil, tt.mockReturn.(error))
			} else {
				mockUseCase.EXPECT().GetBinaries(gomock.Any(), expectedUser).Return(tt.mockReturn.([]entity.Binary), nil)
			}

			http.HandlerFunc(c.GetBinaries).ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)

			if tt.expectedStatus == http.StatusNoContent {
				assert.Empty(t, rr.Body.String())
			} else {
				var response []map[string]interface{}
				err := json.Unmarshal(rr.Body.Bytes(), &response)
				if err != nil {
					t.Fatalf("Failed to unmarshal response body: %v", err)
				}

				var actualBinaries []map[string]interface{}
				for _, binary := range response {
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

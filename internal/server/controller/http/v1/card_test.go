package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nextlag/keeper/internal/entity"
)

func TestAddCard(t *testing.T) {
	c, mockUseCase, ctrl := loadTest(t)
	defer ctrl.Finish()

	expectedUser := entity.User{ID: uuid.New()}

	tests := []struct {
		name           string
		payload        *entity.Card
		contextError   error
		useCaseError   error
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "successful add card",
			payload: &entity.Card{
				ID:              uuid.New(),
				Name:            "test-card",
				CardHolderName:  "John Doe",
				Number:          "1234567812345678",
				Brand:           "Visa",
				ExpirationMonth: "12",
				ExpirationYear:  "2025",
				SecurityCode:    "123",
				Meta:            nil,
			},
			contextError:   nil,
			useCaseError:   nil,
			expectedStatus: http.StatusAccepted,
			expectedBody:   `{"name":"test-card","card_holder_name":"John Doe","number":"1234567812345678","brand":"Visa","expiration_month":"12","expiration_year":"2025","security_code":"123","meta":null}`,
		},
		{
			name:           "bad request error",
			payload:        &entity.Card{},
			contextError:   nil,
			useCaseError:   errors.New("bad request"),
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "bad request\n",
		},
		{
			name:           "context error",
			payload:        &entity.Card{},
			contextError:   errors.New("some unexpected error"),
			useCaseError:   nil,
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `some unexpected error` + "\n" + `{"uuid":"00000000-0000-0000-0000-000000000000","name":"","card_holder_name":"","number":"","brand":"","expiration_month":"","expiration_year":"","security_code":"","meta":null}` + "\n",
		},
		{
			name: "use case error",
			payload: &entity.Card{
				ID:              uuid.New(),
				Name:            "test-card",
				CardHolderName:  "John Doe",
				Number:          "1234567812345678",
				Brand:           "Visa",
				ExpirationMonth: "12",
				ExpirationYear:  "2025",
				SecurityCode:    "123",
				Meta:            nil,
			},
			contextError:   nil,
			useCaseError:   errors.New("use case error"),
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "use case error\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var reqBody *bytes.Buffer
			if tt.payload != nil {
				body, err := json.Marshal(tt.payload)
				require.NoError(t, err)
				reqBody = bytes.NewBuffer(body)
			} else {
				reqBody = bytes.NewBuffer(nil)
			}

			req := httptest.NewRequest(http.MethodPost, userCards, reqBody)
			ctx := context.WithValue(req.Context(), currentUserKey, expectedUser)
			if tt.contextError != nil {
				ctx = context.WithValue(req.Context(), errorKey, tt.contextError)
			}

			req = req.WithContext(ctx)
			rr := httptest.NewRecorder()

			if tt.useCaseError != nil {
				mockUseCase.EXPECT().
					AddCard(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(tt.useCaseError).Times(1)
			} else {
				mockUseCase.EXPECT().
					AddCard(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil).Times(1)
			}

			http.HandlerFunc(c.AddCard).ServeHTTP(rr, req)
			assert.Equal(t, tt.expectedStatus, rr.Code)

			if tt.expectedStatus == http.StatusAccepted {
				var actualResponse map[string]interface{}
				err := json.Unmarshal(rr.Body.Bytes(), &actualResponse)
				require.NoError(t, err)

				actualUUID, ok := actualResponse["uuid"].(string)
				require.True(t, ok, "UUID should be a string")
				_, err = uuid.Parse(actualUUID)
				require.NoError(t, err, "UUID should be a valid UUID")

				delete(actualResponse, "uuid")

				expectedBodyWithoutUUID := `{"name":"test-card","card_holder_name":"John Doe","number":"1234567812345678","brand":"Visa","expiration_month":"12","expiration_year":"2025","security_code":"123","meta":null}`
				var expectedResponse map[string]interface{}
				err = json.Unmarshal([]byte(expectedBodyWithoutUUID), &expectedResponse)
				require.NoError(t, err)

				assert.Equal(t, expectedResponse, actualResponse)
			} else {
				assert.Equal(t, tt.expectedBody, rr.Body.String())
			}
		})
	}
}

func TestGetCards(t *testing.T) {
	c, mockUseCase, ctrl := loadTest(t)
	defer ctrl.Finish()

	expectedUser := entity.User{ID: uuid.New()}

	tests := []struct {
		name           string
		contextError   error
		useCaseError   error
		userCards      []entity.Card
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "successful get cards",
			contextError:   nil,
			useCaseError:   nil,
			userCards:      []entity.Card{{ID: uuid.New(), Name: "test-card", CardHolderName: "", Number: "", Brand: "", ExpirationMonth: "", ExpirationYear: "", SecurityCode: "", Meta: nil}},
			expectedStatus: http.StatusOK,
			expectedBody:   `[{"name":"test-card","card_holder_name":"","number":"","brand":"","expiration_month":"","expiration_year":"","security_code":"","meta":null}]`,
		},
		{
			name:           "context error",
			contextError:   errors.New("context error"),
			useCaseError:   nil,
			userCards:      nil,
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "context error\n",
		},
		{
			name:           "use case error",
			contextError:   nil,
			useCaseError:   errors.New("use case error"),
			userCards:      nil,
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "use case error\n",
		},
		{
			name:           "no cards",
			contextError:   nil,
			useCaseError:   nil,
			userCards:      nil,
			expectedStatus: http.StatusNoContent,
			expectedBody:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, userCards, nil)
			req = req.WithContext(context.WithValue(req.Context(), currentUserKey, expectedUser))

			rr := httptest.NewRecorder()

			if tt.contextError != nil {
				mockUseCase.EXPECT().
					GetCards(gomock.Any(), gomock.Any()).
					Return(nil, tt.contextError).
					Times(1)
			} else {
				if tt.useCaseError != nil {
					mockUseCase.EXPECT().
						GetCards(gomock.Any(), gomock.Any()).
						Return(nil, tt.useCaseError).
						Times(1)
				} else {
					mockUseCase.EXPECT().
						GetCards(gomock.Any(), gomock.Any()).
						Return(tt.userCards, nil).
						Times(1)
				}
			}

			http.HandlerFunc(c.GetCards).ServeHTTP(rr, req)
			assert.Equal(t, tt.expectedStatus, rr.Code)

			if tt.expectedStatus == http.StatusNoContent {
				assert.Empty(t, rr.Body.String())
			} else {
				if tt.expectedStatus == http.StatusInternalServerError {
					assert.Equal(t, tt.expectedBody, rr.Body.String())
				} else {
					var actualResponse []map[string]interface{}
					err := json.Unmarshal(rr.Body.Bytes(), &actualResponse)
					require.NoError(t, err)

					for i := range actualResponse {
						delete(actualResponse[i], "uuid")
					}

					var expectedResponse []map[string]interface{}
					err = json.Unmarshal([]byte(tt.expectedBody), &expectedResponse)
					require.NoError(t, err)

					assert.ElementsMatch(t, expectedResponse, actualResponse)
				}
			}
		})
	}
}

func TestUpdateCard(t *testing.T) {
	c, mockUseCase, ctrl := loadTest(t)
	defer ctrl.Finish()

	validUUID := uuid.New()
	expectedUser := entity.User{ID: validUUID}

	card := entity.Card{
		Name:            "Green",
		CardHolderName:  "DIGIT",
		Number:          "1123123410001234",
		Brand:           "MIR",
		ExpirationMonth: "07",
		ExpirationYear:  "2024",
		SecurityCode:    "003",
	}

	tests := []struct {
		name           string
		mockReturn     error
		expectedStatus int
		expectedBody   string
		cardID         string
		reqBody        interface{}
	}{
		{
			name:           "successful update card",
			mockReturn:     nil,
			expectedStatus: http.StatusAccepted,
			expectedBody:   "",
			cardID:         validUUID.String(),
			reqBody:        card,
		},
		{
			name:           "error from use case",
			mockReturn:     errors.New("update failed"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "update failed\n",
			cardID:         validUUID.String(),
			reqBody:        card,
		},
		{
			name:           "invalid UUID in URL",
			mockReturn:     nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "invalid UUID length: 0\n",
			cardID:         "",
			reqBody:        card,
		},
		{
			name:           "non-UUID string in URL",
			mockReturn:     nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "invalid UUID length: 10\n",
			cardID:         "123a45test",
			reqBody:        card,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expectedStatus == http.StatusAccepted || tt.expectedStatus == http.StatusInternalServerError {
				mockUseCase.EXPECT().
					UpdateCard(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(tt.mockReturn).Times(1)
			}

			reqBody, err := json.Marshal(tt.reqBody)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPatch, userCards+tt.cardID, bytes.NewBuffer(reqBody))
			req = req.WithContext(context.WithValue(req.Context(), currentUserKey, expectedUser))

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.cardID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
			rr := httptest.NewRecorder()

			http.HandlerFunc(c.UpdateCard).ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			assert.Equal(t, tt.expectedBody, rr.Body.String())
		})
	}
}

func TestDelCard(t *testing.T) {
	c, mockUseCase, ctrl := loadTest(t)
	defer ctrl.Finish()

	validUUID := uuid.New()
	expectedUser := entity.User{ID: validUUID}

	tests := []struct {
		name           string
		mockReturn     error
		expectedStatus int
		expectedBody   string
		cardID         string
	}{
		{
			name:           "successful delete card",
			mockReturn:     nil,
			expectedStatus: http.StatusAccepted,
			expectedBody:   "",
			cardID:         validUUID.String(),
		},
		{
			name:           "error from use case",
			mockReturn:     errors.New("delete failed"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "delete failed\n",
			cardID:         validUUID.String(),
		},
		{
			name:           "invalid UUID in URL",
			mockReturn:     nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "invalid UUID length: 0\n",
			cardID:         "",
		},
		{
			name:           "non-UUID string in URL",
			mockReturn:     nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "invalid UUID length: 10\n",
			cardID:         "123a45test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expectedStatus == http.StatusAccepted || tt.expectedStatus == http.StatusInternalServerError {
				mockUseCase.EXPECT().
					DelCard(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(tt.mockReturn).Times(1)
			}

			req := httptest.NewRequest(http.MethodDelete, userCards+tt.cardID, nil)
			req = req.WithContext(context.WithValue(req.Context(), currentUserKey, expectedUser))

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.cardID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
			rr := httptest.NewRecorder()

			http.HandlerFunc(c.DelCard).ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			assert.Equal(t, tt.expectedBody, rr.Body.String())
		})
	}
}

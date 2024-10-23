package api

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid/v5"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/josephlbailey/alert-service/internal/api/models"
	"github.com/josephlbailey/alert-service/internal/db"
	"github.com/josephlbailey/alert-service/internal/db/domain"
	mockdb "github.com/josephlbailey/alert-service/internal/mock"
)

type testCase struct {
	name          string
	externalID    string
	body          gin.H
	buildStubs    func(store *mockdb.MockStore)
	checkResponse func(recorder *httptest.ResponseRecorder)
}

func TestCreateAlert(t *testing.T) {
	alert, message := randomAlert()

	testCases := []testCase{
		{
			name: "create alert with valid body",
			body: gin.H{
				"message": message,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateAlertTX(gomock.Any(), gomock.All(
						gomock.Cond(func(x any) bool { return x.(domain.CreateAlertParams).Message == message })),
					).
					Times(1).
					Return(alert, nil)

			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
				requireBodyMatchAlert(t, alert, recorder.Body)
			},
		},
		{
			name: "create alert with invalid body",
			body: gin.H{
				"message": nil,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateAlertTX(gomock.Any(), gomock.Any()).
					Times(0)

			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			testCase.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(testCase.body)
			require.NoError(t, err)

			url := "/alert"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(data))
			require.NoError(t, err)

			// Add basic auth
			auth := "integrationUser:integrationUserPassword"
			encodedAuth := base64.StdEncoding.EncodeToString([]byte(auth))
			request.Header.Add("Authorization", fmt.Sprintf("Basic %s", encodedAuth))

			server.router.ServeHTTP(recorder, request)
			testCase.checkResponse(recorder)

		})
	}

}

func TestGetAlertByExternalID(t *testing.T) {
	alert, _ := randomAlert()

	testCases := []testCase{
		{
			name:       "get existing alert by valid external ID",
			externalID: alert.ExternalID.String(),
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAlertByExternalID(gomock.Any(), alert.ExternalID).
					Times(1).
					Return(alert, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAlert(t, alert, recorder.Body)
			},
		},
		{
			name:       "get non-existing alert by valid external ID",
			externalID: "f47ac10b-58cc-0372-8567-0e02b2c3d479",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAlertByExternalID(gomock.Any(), gomock.Eq(uuid.Must(uuid.FromString("f47ac10b-58cc-0372-8567-0e02b2c3d479")))).
					Times(1).
					Return(nil, db.ErrAlertNotExists)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
				require.Contains(t, recorder.Body.String(), "alert not found")
			},
		},
		{
			name:       "get alert with invalid external ID format",
			externalID: "invalidUUID",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAlertByExternalID(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				require.Contains(t, recorder.Body.String(), "invalid identifier format")
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			testCase.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/alert/%s", testCase.externalID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			// Add basic auth
			auth := "integrationUser:integrationUserPassword"
			encodedAuth := base64.StdEncoding.EncodeToString([]byte(auth))
			request.Header.Add("Authorization", fmt.Sprintf("Basic %s", encodedAuth))

			server.router.ServeHTTP(recorder, request)
			testCase.checkResponse(recorder)
		})
	}
}

func TestUpdateAlertByExternalID(t *testing.T) {
	alert, message := randomAlert()
	param := randomUpdateAlertParams()
	exID := uuid.Must(uuid.NewV4())
	alert.ExternalID = exID
	upd := time.Now()
	param.UpdatedAt = upd

	testCases := []testCase{
		{
			name:       "update existing alert by valid external ID",
			externalID: alert.ExternalID.String(),
			body: gin.H{
				"message": message,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAlertByExternalID(gomock.Any(), alert.ExternalID).
					Times(1).
					Return(alert, nil)

				store.EXPECT().
					UpdateAlertByIDTX(gomock.Any(), gomock.Cond(func(x any) bool { return x.(domain.UpdateAlertByIDParams).ID == alert.ID })).
					Times(1).
					Return(alert, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAlert(t, alert, recorder.Body)
			},
		},
		{
			name:       "update non-existing alert by valid external ID",
			externalID: "f47ac10b-58cc-0372-8567-0e02b2c3d479",
			body: gin.H{
				"message": message,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAlertByExternalID(gomock.Any(), gomock.Eq(uuid.Must(uuid.FromString("f47ac10b-58cc-0372-8567-0e02b2c3d479")))).
					Times(1).
					Return(nil, db.ErrAlertNotExists)

				store.EXPECT().
					UpdateAlertByIDTX(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
				require.Contains(t, recorder.Body.String(), "alert not found")
			},
		},
		{
			name:       "update alert with invalid external ID format",
			externalID: "invalidUUID",
			body: gin.H{
				"message": message,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAlertByExternalID(gomock.Any(), gomock.Any()).
					Times(0)

				store.EXPECT().
					UpdateAlertByIDTX(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				require.Contains(t, recorder.Body.String(), "invalid identifier format")
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			testCase.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(testCase.body)
			require.NoError(t, err)

			url := fmt.Sprintf("/alert/%s", testCase.externalID)
			request, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(data))
			require.NoError(t, err)

			// Add basic auth
			auth := "integrationUser:integrationUserPassword"
			encodedAuth := base64.StdEncoding.EncodeToString([]byte(auth))
			request.Header.Add("Authorization", fmt.Sprintf("Basic %s", encodedAuth))

			server.router.ServeHTTP(recorder, request)
			testCase.checkResponse(recorder)
		})
	}
}

func TestDeleteAlertByExternalID(t *testing.T) {
	alert, _ := randomAlert()

	testCases := []testCase{
		{
			name:       "delete existing alert by valid external ID",
			externalID: alert.ExternalID.String(),
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAlertByExternalID(gomock.Any(), alert.ExternalID).
					Times(1).
					Return(alert, nil)

				store.EXPECT().
					DeleteAlertByIDTX(gomock.Any(), alert.ID).
					Times(1).
					Return(nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:       "delete non-existing alert by valid external ID",
			externalID: "f47ac10b-58cc-0372-8567-0e02b2c3d479",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAlertByExternalID(gomock.Any(), gomock.Eq(uuid.Must(uuid.FromString("f47ac10b-58cc-0372-8567-0e02b2c3d479")))).
					Times(1).
					Return(nil, db.ErrAlertNotExists)

				store.EXPECT().
					DeleteAlertByIDTX(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
				require.Contains(t, recorder.Body.String(), "alert not found")
			},
		},
		{
			name:       "delete alert with invalid external ID format",
			externalID: "invalidUUID",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAlertByExternalID(gomock.Any(), gomock.Any()).
					Times(0)

				store.EXPECT().
					DeleteAlertByIDTX(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				require.Contains(t, recorder.Body.String(), "invalid identifier format")
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			testCase.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/alert/%s", testCase.externalID)
			request, err := http.NewRequest(http.MethodDelete, url, nil)
			require.NoError(t, err)

			// Add basic auth
			auth := "integrationUser:integrationUserPassword"
			encodedAuth := base64.StdEncoding.EncodeToString([]byte(auth))
			request.Header.Add("Authorization", fmt.Sprintf("Basic %s", encodedAuth))

			server.router.ServeHTTP(recorder, request)
			testCase.checkResponse(recorder)
		})
	}
}

func randomAlert() (alert *domain.Alert, message string) {

	message = "Hello there"

	alert = &domain.Alert{
		ID:        1,
		CreatedAt: time.Now(),
		Message:   message,
	}
	return
}

func randomUpdateAlertParams() (params *domain.UpdateAlertByIDParams) {
	message := "Steven why are you like this"

	params = &domain.UpdateAlertByIDParams{
		ID:      1,
		Message: message,
	}
	return
}

func requireBodyMatchAlert(t *testing.T, alert *domain.Alert, body *bytes.Buffer) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	fmt.Println("Response body:", string(data))

	var gotAlert models.AlertRes
	err = json.Unmarshal(data, &gotAlert)

	require.NoError(t, err)
	require.Equalf(t, alert.ExternalID, gotAlert.ExternalID, "want ExternalID: %v, got ExternalID: %v", alert.ExternalID, gotAlert.ExternalID)
	require.NotNilf(t, gotAlert.CreatedAt, "expected alert to contain CreatedAt")
	require.NotNilf(t, gotAlert.UpdatedAt, "expected alert to contain UpdatedAt")
	require.Equalf(t, alert.Message, gotAlert.Message, "want Message: %v, got Message: %v", alert.Message, gotAlert.Message)
}

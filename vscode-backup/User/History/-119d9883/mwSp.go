package tabapay

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"strings"
	"testing"

	"github.com/Rhymond/go-money"
	"github.com/golang/mock/gomock"
	"github.com/wingocard/braavos/internal/wlog"
	"github.com/wingocard/braavos/pkg/client"
	"github.com/wingocard/braavos/pkg/entities/acquiring"
	"github.com/wingocard/braavos/pkg/testutil"
	"gotest.tools/v3/assert"
)

func TestBuildRequest(t *testing.T) {
	tt := []struct {
		name        string
		method      string
		url         string
		body        interface{}
		expectedErr error
	}{
		{
			name:        "error building request",
			method:      http.MethodPost,
			url:         "/blah",
			body:        math.NaN(),
			expectedErr: errors.New("tabapay has been disabled for program wind down"),
		},
		{
			name:        "success",
			method:      http.MethodPost,
			url:         "/test",
			body:        &acquiring.FundingAccount{},
			expectedErr: errors.New("tabapay has been disabled for program wind down"),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tbp := New(&Config{
				BaseURL:  "base",
				ClientID: "1111",
			})

			req, err := tbp.buildRequest(context.Background(), tc.method, tc.url, tc.body)

			if tc.expectedErr == nil {
				assert.NilError(t, err)
				assert.Equal(t, req.URL.String(), tbp.baseURL+tc.url)
				assert.Equal(t, req.Header.Get("content-type"), "application/json")
				assert.Equal(t, req.Header.Get(HeaderAuthorizationKey), HeaderAuthorizationValue+tbp.bearerToken)

				b, err := ioutil.ReadAll(req.Body)
				assert.NilError(t, err)
				req.Body.Close()
				assert.Equal(t, bytes.Index(b, []byte("\n")), -1)
				return
			}

			assert.Assert(t, req == nil)
			assert.ErrorContains(t, err, tc.expectedErr.Error())
		})
	}
}

func TestCreateFundingAccount(t *testing.T) {
	testRefID := "9999"
	testAccountID := "0000"

	tt := []struct {
		name           string
		fundingAccount *acquiring.FundingAccount
		setupMocks     func(mc *client.MockClient)
		expectedErr    error
	}{
		{
			name:           "account details error",
			fundingAccount: &acquiring.FundingAccount{},
			setupMocks: func(mc *client.MockClient) {
				mc.EXPECT().DoJSON(gomock.Any(), gomock.Any()).Times(0)
			},
			expectedErr: errors.New("tabapayclient: account details are not of type *acquiring.Card"),
		},
		{
			name: "error making request",
			fundingAccount: &acquiring.FundingAccount{
				AccountDetails: &acquiring.Card{},
				AddressLine1:   "123 avenue",
				City:           "city",
				ProvState:      "QC",
				PostalCodeZip:  "A1A1A1",
				Country:        "CA",
			},
			setupMocks: func(mc *client.MockClient) {
				mc.EXPECT().
					DoJSON(gomock.Any(), gomock.Any()).
					Return(errors.New("failure")).
					Times(0)
			},
			expectedErr: errors.New("tabapay has been disabled for program wind down"),
		},
		{
			name: "duplicate found unmarhsal error",
			fundingAccount: &acquiring.FundingAccount{
				AccountDetails: &acquiring.Card{},
				AddressLine1:   "123 avenue",
				City:           "city",
				ProvState:      "CA",
				PostalCodeZip:  "90210",
				Country:        "US",
			},
			setupMocks: func(mc *client.MockClient) {
				mc.EXPECT().
					DoJSON(gomock.Any(), gomock.Any()).
					Return(&client.HTTPError{
						StatusCode: http.StatusConflict,
						Body:       `{`,
					}).
					Times(0)
			},
			expectedErr: errors.New("tabapay has been disabled for program wind down"),
		},
		{
			name: "more than 1 duplicate found error",
			fundingAccount: &acquiring.FundingAccount{
				AccountDetails: &acquiring.Card{},
				AddressLine1:   "123 avenue",
				City:           "city",
				ProvState:      "CA",
				PostalCodeZip:  "90210",
				Country:        "US",
			},
			setupMocks: func(mc *client.MockClient) {
				mc.EXPECT().
					DoJSON(gomock.Any(), gomock.Any()).
					Return(&client.HTTPError{
						StatusCode: http.StatusConflict,
						Body:       `{"duplicateAccountIDs": ["123", "456"]}`,
					}).
					Times(0)
			},
			expectedErr: errors.New("tabapay has been disabled for program wind down"),
		},
		{
			name: "duplicate found not a match error",
			fundingAccount: &acquiring.FundingAccount{
				AccountDetails: &acquiring.Card{
					Last4:          "1234",
					ExpirationDate: "2000-12",
				},
				FirstName:     "billy",
				LastName:      "blaze",
				AddressLine1:  "123 avenue",
				City:          "city",
				ProvState:     "QC",
				PostalCodeZip: "A1A1A1",
				Country:       "CA",
				Phone:         "+15551234567",
			},
			setupMocks: func(mc *client.MockClient) {
				first := mc.EXPECT().
					DoJSON(gomock.Any(), gomock.Any()).
					Return(&client.HTTPError{
						StatusCode: http.StatusConflict,
						Body:       `{"duplicateAccountIDs": ["123"]}`,
					}).
					Times(0)
				mc.EXPECT().
					DoJSON(gomock.Any(), gomock.Any()).
					DoAndReturn(func(req *http.Request, res *getFundingAccountResponse) error {
						*res = getFundingAccountResponse{
							ReferenceID: testRefID,
							Card: cardResponse{
								Last4:          "1234",
								ExpirationDate: "200012",
							},
							Owner: accountOwner{
								Address: &address{
									Line1:   "123 avenue",
									City:    "city",
									State:   "QC",
									ZipCode: "A1A 1A1",
									Country: caCountryCode,
								},
								Name: &name{
									First: "george",
									Last:  "costanza",
								},
								Phone: &phone{
									Number: "5551234567",
								},
							},
						}
						return nil
					}).
					Times(0).
					After(first)
			},
			expectedErr: errors.New("tabapay has been disabled for program wind down"),
		},
		{
			name: "duplicate found success",
			fundingAccount: &acquiring.FundingAccount{
				AccountDetails: &acquiring.Card{
					Last4:          "1234",
					ExpirationDate: "2000-12",
				},
				FirstName:     "billy",
				LastName:      "blaze",
				AddressLine1:  "123 avenue",
				City:          "city",
				ProvState:     "CA",
				PostalCodeZip: "90210",
				Country:       "US",
				Phone:         "+15551234567",
			},
			setupMocks: func(mc *client.MockClient) {
				first := mc.EXPECT().
					DoJSON(gomock.Any(), gomock.Any()).
					Return(&client.HTTPError{
						StatusCode: http.StatusConflict,
						Body:       `{"duplicateAccountIDs": ["` + testAccountID + `"]}`,
					}).
					Times(0)
				mc.EXPECT().
					DoJSON(gomock.Any(), gomock.Any()).
					DoAndReturn(func(req *http.Request, res *getFundingAccountResponse) error {
						*res = getFundingAccountResponse{
							ReferenceID: testRefID,
							Card: cardResponse{
								Last4:          "1234",
								ExpirationDate: "200012",
							},
							Owner: accountOwner{
								Address: &address{
									Line1:   "123 avenue",
									City:    "city",
									State:   "CA",
									ZipCode: "90210",
								},
								Name: &name{
									First: "billy",
									Last:  "blaze",
								},
								Phone: &phone{
									Number: "5551234567",
								},
							},
						}
						return nil
					}).
					Times(0).
					After(first)
			},
			expectedErr: errors.New("tabapay has been disabled for program wind down"),
		},
		{
			name: "success",
			fundingAccount: &acquiring.FundingAccount{
				ID:            "1234",
				UserID:        "5678",
				FirstName:     "Billy",
				LastName:      "Blaze",
				AddressLine1:  "123 avenue",
				City:          "city",
				ProvState:     "QC",
				PostalCodeZip: "A1A1A1",
				Country:       "CA",
				Phone:         "+15551234567",
				AccountDetails: &acquiring.Card{
					Type:           acquiring.Visa,
					ExpirationDate: "1122",
					EncryptedData:  "carddata",
					Last4:          "5555",
				},
			},
			setupMocks: func(mc *client.MockClient) {
				mc.EXPECT().
					DoJSON(gomock.Any(), gomock.Any()).
					DoAndReturn(func(req *http.Request, res *createFundingAccountResponse) error {
						res.AccountID = testAccountID
						return nil
					}).
					Times(0)
			},
			expectedErr: errors.New("tabapay has been disabled for program wind down"),
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			mc := client.NewMockClient(ctrl)
			tc.setupMocks(mc)

			tbp := &TabaPay{
				httpClient: mc,
				baseURL:    "https://example.com",
				refIDGen:   &staticRefID{ID: testRefID},
			}

			err := tbp.CreateFundingAccount(context.Background(), wlog.NewNopLogger(), tc.fundingAccount)

			if tc.expectedErr == nil {
				assert.NilError(t, err)
				assert.Assert(t, tc.fundingAccount.ExternalAccountID == testAccountID)
				return
			}

			assert.ErrorContains(t, err, tc.expectedErr.Error())
		})
	}
}

func TestUpdateFundingAccount(t *testing.T) {
	tt := []struct {
		name           string
		fundingAccount *acquiring.FundingAccount
		setupMocks     func(mc *client.MockClient)
		expectedErr    error
	}{
		{
			name: "account details error",
			fundingAccount: &acquiring.FundingAccount{
				AddressLine1:  "123 avenue",
				City:          "city",
				ProvState:     "QC",
				PostalCodeZip: "A1A1A1",
				Country:       "CA",
				Phone:         "+15551234567",
			},
			setupMocks: func(mc *client.MockClient) {
				mc.EXPECT().DoJSON(gomock.Any(), gomock.Any()).Times(0)
			},
			expectedErr: errors.New("tabapayclient: account details are not of type *acquiring.Card"),
		},
		{
			name: "error making request",
			fundingAccount: &acquiring.FundingAccount{
				AccountDetails: &acquiring.Card{},
				AddressLine1:   "123 avenue",
				City:           "city",
				ProvState:      "QC",
				PostalCodeZip:  "A1A1A1",
				Country:        "CA",
				Phone:          "+15551234567",
			},
			setupMocks: func(mc *client.MockClient) {
				mc.EXPECT().
					DoJSON(gomock.Any(), gomock.Any()).
					Return(errors.New("failure")).
					Times(0)
			},
			expectedErr: errors.New("tabapay has been disabled for program wind down"),
		},
		{
			name: "success",
			fundingAccount: &acquiring.FundingAccount{
				AccountDetails: &acquiring.Card{},
				AddressLine1:   "123 avenue",
				City:           "city",
				ProvState:      "QC",
				PostalCodeZip:  "A1A1A1",
				Country:        "CA",
				Phone:          "+15551234567",
			},
			setupMocks: func(mc *client.MockClient) {
				mc.EXPECT().
					DoJSON(gomock.Any(), gomock.Any()).
					Return(nil).
					Times(0)
			},
			expectedErr: errors.New("tabapay has been disabled for program wind down"),
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			mc := client.NewMockClient(ctrl)
			tc.setupMocks(mc)

			tbp := &TabaPay{
				httpClient: mc,
				baseURL:    "https://example.com",
				refIDGen:   newReferenceID(),
			}

			err := tbp.UpdateFundingAccount(context.Background(), "12345", tc.fundingAccount)

			if tc.expectedErr == nil {
				assert.NilError(t, err)
				return
			}

			assert.ErrorContains(t, err, tc.expectedErr.Error())
		})
	}
}

func TestDeleteFundingAccount(t *testing.T) {
	tt := []struct {
		name        string
		amount      *money.Money
		setupMocks  func(mc *client.MockClient)
		expectedErr error
	}{
		{
			name: "error making request",
			setupMocks: func(mc *client.MockClient) {
				mc.EXPECT().
					DoJSON(gomock.Any(), gomock.Any()).
					Return(errors.New("failure")).
					Times(1)
			},
			expectedErr: errors.New("tabapay has been disabled for program wind down"),
		},
		{
			name: "success",
			setupMocks: func(mc *client.MockClient) {
				mc.EXPECT().
					DoJSON(gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)
			},
			expectedErr: errors.New("tabapay has been disabled for program wind down"),
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			mc := client.NewMockClient(ctrl)
			tc.setupMocks(mc)

			tbp := &TabaPay{
				httpClient: mc,
				baseURL:    "https://example.com",
			}

			err := tbp.DeleteFundingAccount(context.Background(), "12345")

			if tc.expectedErr == nil {
				assert.NilError(t, err)
				return
			}

			assert.ErrorContains(t, err, tc.expectedErr.Error())
		})
	}
}

func TestFormatAmount(t *testing.T) {
	tt := []struct {
		name           string
		amount         *money.Money
		expectedFormat string
	}{
		{
			name:           "1 dollar",
			amount:         money.New(100, "CAD"),
			expectedFormat: "1.00",
		},
		{
			name:           "> 1 dollar",
			amount:         money.New(1000, "CAD"),
			expectedFormat: "10.00",
		},
		{
			name:           "< 1 dollar",
			amount:         money.New(1, "USD"),
			expectedFormat: "0.01",
		},
		{
			name:           "not a whole amount",
			amount:         money.New(135, "USD"),
			expectedFormat: "1.35",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			amountStr := formatAmount(tc.amount.Amount())
			assert.Equal(t, amountStr, tc.expectedFormat)
		})
	}
}

type TransactionBodyMatcher struct {
	expectedBodyJSON string
}

func (t *TransactionBodyMatcher) Matches(x interface{}) bool {
	testReq, ok := x.(*http.Request)
	if !ok {
		return false
	}

	testReqBodyJSON, err := ioutil.ReadAll(testReq.Body)
	if err != nil {
		return false
	}

	res := testutil.CompareResponseJSON(string(testReqBodyJSON), t.expectedBodyJSON)()
	return res.Success()
}

func (t *TransactionBodyMatcher) String() string {
	return fmt.Sprintf("transaction body matches %s", t.expectedBodyJSON)
}

func MatchTransactionBody(expectedBody *createTransactionPayload) gomock.Matcher {
	expectedBodyJSON, err := json.Marshal(expectedBody)
	if err != nil {
		panic(err)
	}
	return &TransactionBodyMatcher{expectedBodyJSON: string(expectedBodyJSON)}
}

func TestFundBankAccount(t *testing.T) {
	testSrcAccountID := "123456"
	testDestAccountID := "000000"
	testAmount := money.New(100, "USD")

	tt := []struct {
		name            string
		expectedReqBody *createTransactionPayload
		expectedErr     error
		setupMocks      func(mc *client.MockClient, expectedBody *createTransactionPayload)
	}{
		{
			name: "error making request",
			expectedReqBody: &createTransactionPayload{
				ReferenceID: "",
				Type:        acquiring.Pull,
				Accounts: &accounts{
					SourceAccountID:      testSrcAccountID,
					DestinationAccountID: testDestAccountID,
				},
				Currency: testAmount.Currency().NumericCode,
				Amount:   formatAmount(testAmount.Amount()),
				Memo:     "",
			},
			expectedErr: errors.New("tabapayclient: error making request"),
			setupMocks: func(mc *client.MockClient, expectedBody *createTransactionPayload) {
				mc.EXPECT().
					DoJSON(MatchTransactionBody(expectedBody), gomock.Any()).
					Return(errors.New("failure")).
					Times(1)
			},
		},
		{
			name: "network response code error",
			expectedReqBody: &createTransactionPayload{
				ReferenceID: "",
				Type:        acquiring.Pull,
				Accounts: &accounts{
					SourceAccountID:      testSrcAccountID,
					DestinationAccountID: testDestAccountID,
				},
				Currency: testAmount.Currency().NumericCode,
				Amount:   formatAmount(testAmount.Amount()),
				Memo:     "",
			},
			expectedErr: errors.New("transaction error (transactionID: 1234), received network response code 51"),
			setupMocks: func(mc *client.MockClient, expectedBody *createTransactionPayload) {
				mc.EXPECT().
					DoJSON(MatchTransactionBody(expectedBody), gomock.Any()).
					DoAndReturn(func(req *http.Request, res *createTransactionResponse) error {
						res.TransactionID = "1234"
						res.NetworkRC = "51"
						return nil
					}).
					Times(1)
			},
		},
		{
			name: "success - nrc approved",
			expectedReqBody: &createTransactionPayload{
				ReferenceID: "",
				Type:        acquiring.Pull,
				Accounts: &accounts{
					SourceAccountID:      testSrcAccountID,
					DestinationAccountID: testDestAccountID,
				},
				Currency: testAmount.Currency().NumericCode,
				Amount:   formatAmount(testAmount.Amount()),
				Memo:     "",
			},
			expectedErr: nil,
			setupMocks: func(mc *client.MockClient, expectedBody *createTransactionPayload) {
				mc.EXPECT().
					DoJSON(MatchTransactionBody(expectedBody), gomock.Any()).
					DoAndReturn(func(req *http.Request, res *createTransactionResponse) error {
						res.TransactionID = "1234"
						res.NetworkRC = nrcApproved
						return nil
					}).
					Times(1)
			},
		},
		{
			name: "success - nrc no reason to decline",
			expectedReqBody: &createTransactionPayload{
				ReferenceID: "",
				Type:        acquiring.Pull,
				Accounts: &accounts{
					SourceAccountID:      testSrcAccountID,
					DestinationAccountID: testDestAccountID,
				},
				Currency: testAmount.Currency().NumericCode,
				Amount:   formatAmount(testAmount.Amount()),
				Memo:     "",
			},
			expectedErr: nil,
			setupMocks: func(mc *client.MockClient, expectedBody *createTransactionPayload) {
				mc.EXPECT().
					DoJSON(MatchTransactionBody(expectedBody), gomock.Any()).
					DoAndReturn(func(req *http.Request, res *createTransactionResponse) error {
						res.TransactionID = "1234"
						res.NetworkRC = nrcNoReasonToDecline
						return nil
					}).
					Times(1)
			},
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			mc := client.NewMockClient(ctrl)
			tc.setupMocks(mc, tc.expectedReqBody)

			tbp := &TabaPay{
				httpClient:          mc,
				baseURL:             "https://example.com",
				settlementAccountID: testDestAccountID,
				refIDGen:            &staticRefID{},
			}

			fundingRes := tbp.FundBankAccount(context.Background(), testSrcAccountID, testAmount)

			if tc.expectedErr == nil {
				assert.NilError(t, fundingRes.Error())
				assert.Assert(t, fundingRes.TransactionID != "")
				return
			}

			assert.ErrorContains(t, fundingRes.Error(), tc.expectedErr.Error())
		})
	}
}

func TestFormatAddress(t *testing.T) {
	tt := []struct {
		name            string
		address         *address
		expectedAddress *address
		expectedErr     error
	}{
		{
			name: "unsupported country",
			address: &address{
				Country: "the moon",
			},
			expectedAddress: nil,
			expectedErr:     errors.New(`unsupported country: "the moon"`),
		},
		{
			name: "invalid canadian postal code",
			address: &address{
				Country: "CA",
				ZipCode: "123",
			},
			expectedAddress: nil,
			expectedErr:     errors.New("postal code is invalid"),
		},
		{
			name: "format US address",
			address: &address{
				Line1:   "123 abc",
				Line2:   "suite 8",
				City:    "city",
				State:   "CA",
				Country: "US",
				ZipCode: "90210",
			},
			expectedAddress: &address{
				Line1:   "123 abc",
				Line2:   "suite 8",
				City:    "city",
				State:   "CA",
				Country: usCountryCode,
				ZipCode: "90210",
			},
			expectedErr: nil,
		},
		{
			name: "format CA address",
			address: &address{
				Line1:   "123 abc",
				Line2:   "suite 8",
				City:    "city",
				State:   "QC",
				Country: "CA",
				ZipCode: "A1A1A1",
			},
			expectedAddress: &address{
				Line1:   "123 abc",
				Line2:   "suite 8",
				City:    "city",
				State:   "QC",
				Country: caCountryCode,
				ZipCode: "A1A 1A1",
			},
			expectedErr: nil,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.address.Format()

			if tc.expectedErr == nil {
				assert.NilError(t, err)
				assert.DeepEqual(t, tc.address, tc.expectedAddress)
				return
			}

			assert.ErrorContains(t, err, tc.expectedErr.Error())
		})
	}
}

func TestGetEncryptionKey(t *testing.T) {
	tt := []struct {
		name      string
		confKeyID string
		confKey   string
	}{
		{
			name:      "success",
			confKeyID: "abcd",
			confKey:   "1234",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tbp := &TabaPay{
				publicKeyID: tc.confKeyID,
				publicKey:   tc.confKey,
			}

			key := tbp.GetEncryptionKey(context.Background())
			assert.Equal(t, key.ID, tc.confKeyID)
			assert.Equal(t, key.Key, tc.confKey)
		})
	}
}

func TestGetFundingAccount(t *testing.T) {
	tt := []struct {
		name        string
		setupMocks  func(mc *client.MockClient)
		expectedErr error
	}{
		{
			name: "error making request",
			setupMocks: func(mc *client.MockClient) {
				mc.EXPECT().
					DoJSON(gomock.Any(), gomock.Any()).
					Return(errors.New("failure")).
					Times(1)
			},
			expectedErr: errors.New("tabapayclient: error making request"),
		},
		{
			name: "success",
			setupMocks: func(mc *client.MockClient) {
				mc.EXPECT().
					DoJSON(gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)
			},
			expectedErr: nil,
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			mc := client.NewMockClient(ctrl)
			tc.setupMocks(mc)

			tbp := &TabaPay{
				httpClient: mc,
				refIDGen:   newReferenceID(),
			}

			_, err := tbp.getFundingAccount(context.Background(), "1234")

			if tc.expectedErr == nil {
				assert.NilError(t, err)
				return
			}

			assert.ErrorContains(t, err, tc.expectedErr.Error())
		})
	}
}

func TestDuplicateMatchCurrentReq(t *testing.T) {
	tt := []struct {
		name        string
		curReq      *createFundingAccountPayload
		setupMocks  func(mc *client.MockClient)
		expectedErr error
	}{
		{
			name:   "error retrieving duplicate account",
			curReq: nil,
			setupMocks: func(mc *client.MockClient) {
				mc.EXPECT().
					DoJSON(gomock.Any(), gomock.Any()).
					Return(errors.New("boom")).
					Times(1)
			},
			expectedErr: errors.New("tabapayclient: error retrieving duplicate account with id"),
		},
		{
			name: "no match bad ref ID",
			curReq: &createFundingAccountPayload{
				ReferenceID: "bad",
				Owner: &accountOwner{
					Name: &name{
						First: "Peter",
						Last:  "Parker",
					},
					Address: &address{
						Line1:   "123 place",
						City:    "asd",
						State:   "QC",
						ZipCode: "A1A 1A1",
						Country: caCountryCode,
					},
					Phone: &phone{
						Number: "5551234567",
					},
				},
			},
			setupMocks: func(mc *client.MockClient) {
				mc.EXPECT().
					DoJSON(gomock.Any(), gomock.Any()).
					DoAndReturn(func(req *http.Request, res *getFundingAccountResponse) error {
						*res = getFundingAccountResponse{
							ReferenceID: "1234567",
							Owner: accountOwner{
								Name: &name{
									First: "Peter",
									Last:  "Parker",
								},
								Address: &address{
									Line1:   "123 place",
									City:    "asd",
									State:   "QC",
									ZipCode: "A1A 1A1",
									Country: caCountryCode,
								},
								Phone: &phone{
									Number: "5551234567",
								},
							}}

						return nil
					}).
					Times(1)
			},
			expectedErr: errors.New("tabapayclient: duplicate accountID: 123 found but does not match current request info"),
		},
		{
			name: "no match bad first name",
			curReq: &createFundingAccountPayload{
				ReferenceID: "1234567",
				Owner: &accountOwner{
					Name: &name{
						First: "Frank",
						Last:  "Parker",
					},
					Address: &address{
						Line1:   "123 place",
						City:    "asd",
						State:   "QC",
						ZipCode: "A1A 1A1",
						Country: caCountryCode,
					},
					Phone: &phone{
						Number: "5551234567",
					},
				},
			},
			setupMocks: func(mc *client.MockClient) {
				mc.EXPECT().
					DoJSON(gomock.Any(), gomock.Any()).
					DoAndReturn(func(req *http.Request, res *getFundingAccountResponse) error {
						*res = getFundingAccountResponse{
							ReferenceID: "1234567",
							Owner: accountOwner{
								Name: &name{
									First: "Peter",
									Last:  "Parker",
								},
								Address: &address{
									Line1:   "123 place",
									City:    "asd",
									State:   "QC",
									ZipCode: "A1A 1A1",
									Country: caCountryCode,
								},
								Phone: &phone{
									Number: "5551234567",
								},
							}}

						return nil
					}).
					Times(1)
			},
			expectedErr: errors.New("tabapayclient: duplicate accountID: 123 found but does not match current request info"),
		},
		{
			name: "no match bad last name",
			curReq: &createFundingAccountPayload{
				ReferenceID: "1234567",
				Owner: &accountOwner{
					Name: &name{
						First: "Peter",
						Last:  "Park",
					},
					Address: &address{
						Line1:   "123 place",
						City:    "asd",
						State:   "QC",
						ZipCode: "A1A 1A1",
						Country: caCountryCode,
					},
					Phone: &phone{
						Number: "5551234567",
					},
				},
			},
			setupMocks: func(mc *client.MockClient) {
				mc.EXPECT().
					DoJSON(gomock.Any(), gomock.Any()).
					DoAndReturn(func(req *http.Request, res *getFundingAccountResponse) error {
						*res = getFundingAccountResponse{
							ReferenceID: "1234567",
							Owner: accountOwner{
								Name: &name{
									First: "Peter",
									Last:  "Parker",
								},
								Address: &address{
									Line1:   "123 place",
									City:    "asd",
									State:   "QC",
									ZipCode: "A1A 1A1",
									Country: caCountryCode,
								},
								Phone: &phone{
									Number: "5551234567",
								},
							}}

						return nil
					}).
					Times(1)
			},
			expectedErr: errors.New("tabapayclient: duplicate accountID: 123 found but does not match current request info"),
		},
		{
			name: "no match bad phone number",
			curReq: &createFundingAccountPayload{
				ReferenceID: "1234567",
				Owner: &accountOwner{
					Name: &name{
						First: "Peter",
						Last:  "Parker",
					},
					Address: &address{
						Line1:   "123 place",
						City:    "asd",
						State:   "QC",
						ZipCode: "A1A 1A1",
						Country: caCountryCode,
					},
					Phone: &phone{
						Number: "88888",
					},
				},
			},
			setupMocks: func(mc *client.MockClient) {
				mc.EXPECT().
					DoJSON(gomock.Any(), gomock.Any()).
					DoAndReturn(func(req *http.Request, res *getFundingAccountResponse) error {
						*res = getFundingAccountResponse{
							ReferenceID: "1234567",
							Owner: accountOwner{
								Name: &name{
									First: "Peter",
									Last:  "Parker",
								},
								Address: &address{
									Line1:   "123 place",
									City:    "asd",
									State:   "QC",
									ZipCode: "A1A 1A1",
									Country: caCountryCode,
								},
								Phone: &phone{
									Number: "5551234567",
								},
							}}

						return nil
					}).
					Times(1)
			},
			expectedErr: errors.New("tabapayclient: duplicate accountID: 123 found but does not match current request info"),
		},
		{
			name: "no match bad address line 1",
			curReq: &createFundingAccountPayload{
				ReferenceID: "1234567",
				Owner: &accountOwner{
					Name: &name{
						First: "Peter",
						Last:  "Parker",
					},
					Address: &address{
						Line1:   "456 street",
						City:    "asd",
						State:   "QC",
						ZipCode: "A1A 1A1",
						Country: caCountryCode,
					},
					Phone: &phone{
						Number: "5551234567",
					},
				},
			},
			setupMocks: func(mc *client.MockClient) {
				mc.EXPECT().
					DoJSON(gomock.Any(), gomock.Any()).
					DoAndReturn(func(req *http.Request, res *getFundingAccountResponse) error {
						*res = getFundingAccountResponse{
							ReferenceID: "1234567",
							Owner: accountOwner{
								Name: &name{
									First: "Peter",
									Last:  "Parker",
								},
								Address: &address{
									Line1:   "123 place",
									City:    "asd",
									State:   "QC",
									ZipCode: "A1A 1A1",
									Country: caCountryCode,
								},
								Phone: &phone{
									Number: "5551234567",
								},
							}}

						return nil
					}).
					Times(1)
			},
			expectedErr: errors.New("tabapayclient: duplicate accountID: 123 found but does not match current request info"),
		},
		{
			name: "no match bad city",
			curReq: &createFundingAccountPayload{
				ReferenceID: "1234567",
				Owner: &accountOwner{
					Name: &name{
						First: "Peter",
						Last:  "Parker",
					},
					Address: &address{
						Line1:   "123 place",
						City:    "nope",
						State:   "QC",
						ZipCode: "A1A 1A1",
						Country: caCountryCode,
					},
					Phone: &phone{
						Number: "5551234567",
					},
				},
			},
			setupMocks: func(mc *client.MockClient) {
				mc.EXPECT().
					DoJSON(gomock.Any(), gomock.Any()).
					DoAndReturn(func(req *http.Request, res *getFundingAccountResponse) error {
						*res = getFundingAccountResponse{
							ReferenceID: "1234567",
							Owner: accountOwner{
								Name: &name{
									First: "Peter",
									Last:  "Parker",
								},
								Address: &address{
									Line1:   "123 place",
									City:    "asd",
									State:   "QC",
									ZipCode: "A1A 1A1",
									Country: caCountryCode,
								},
								Phone: &phone{
									Number: "5551234567",
								},
							}}

						return nil
					}).
					Times(1)
			},
			expectedErr: errors.New("tabapayclient: duplicate accountID: 123 found but does not match current request info"),
		},
		{
			name: "no match bad state",
			curReq: &createFundingAccountPayload{
				ReferenceID: "1234567",
				Owner: &accountOwner{
					Name: &name{
						First: "Peter",
						Last:  "Parker",
					},
					Address: &address{
						Line1:   "123 place",
						City:    "asd",
						State:   "ON",
						ZipCode: "A1A 1A1",
						Country: caCountryCode,
					},
					Phone: &phone{
						Number: "5551234567",
					},
				},
			},
			setupMocks: func(mc *client.MockClient) {
				mc.EXPECT().
					DoJSON(gomock.Any(), gomock.Any()).
					DoAndReturn(func(req *http.Request, res *getFundingAccountResponse) error {
						*res = getFundingAccountResponse{
							ReferenceID: "1234567",
							Owner: accountOwner{
								Name: &name{
									First: "Peter",
									Last:  "Parker",
								},
								Address: &address{
									Line1:   "123 place",
									City:    "asd",
									State:   "QC",
									ZipCode: "A1A 1A1",
									Country: caCountryCode,
								},
								Phone: &phone{
									Number: "5551234567",
								},
							}}

						return nil
					}).
					Times(1)
			},
			expectedErr: errors.New("tabapayclient: duplicate accountID: 123 found but does not match current request info"),
		},
		{
			name: "no match bad zipcode",
			curReq: &createFundingAccountPayload{
				ReferenceID: "1234567",
				Owner: &accountOwner{
					Name: &name{
						First: "Peter",
						Last:  "Parker",
					},
					Address: &address{
						Line1:   "123 place",
						City:    "asd",
						State:   "QC",
						ZipCode: "A1A 2A2",
						Country: caCountryCode,
					},
					Phone: &phone{
						Number: "5551234567",
					},
				},
			},
			setupMocks: func(mc *client.MockClient) {
				mc.EXPECT().
					DoJSON(gomock.Any(), gomock.Any()).
					DoAndReturn(func(req *http.Request, res *getFundingAccountResponse) error {
						*res = getFundingAccountResponse{
							ReferenceID: "1234567",
							Owner: accountOwner{
								Name: &name{
									First: "Peter",
									Last:  "Parker",
								},
								Address: &address{
									Line1:   "123 place",
									City:    "asd",
									State:   "QC",
									ZipCode: "A1A 1A1",
									Country: caCountryCode,
								},
								Phone: &phone{
									Number: "5551234567",
								},
							}}

						return nil
					}).
					Times(1)
			},
			expectedErr: errors.New("tabapayclient: duplicate accountID: 123 found but does not match current request info"),
		},
		{
			name: "no match bad country",
			curReq: &createFundingAccountPayload{
				ReferenceID: "1234567",
				Owner: &accountOwner{
					Name: &name{
						First: "Peter",
						Last:  "Parker",
					},
					Address: &address{
						Line1:   "123 place",
						City:    "asd",
						State:   "QC",
						ZipCode: "A1A 1A1",
						Country: usCountryCode,
					},
					Phone: &phone{
						Number: "5551234567",
					},
				},
			},
			setupMocks: func(mc *client.MockClient) {
				mc.EXPECT().
					DoJSON(gomock.Any(), gomock.Any()).
					DoAndReturn(func(req *http.Request, res *getFundingAccountResponse) error {
						*res = getFundingAccountResponse{
							ReferenceID: "1234567",
							Owner: accountOwner{
								Name: &name{
									First: "Peter",
									Last:  "Parker",
								},
								Address: &address{
									Line1:   "123 place",
									City:    "asd",
									State:   "QC",
									ZipCode: "A1A 1A1",
									Country: caCountryCode,
								},
								Phone: &phone{
									Number: "5551234567",
								},
							}}

						return nil
					}).
					Times(1)
			},
			expectedErr: errors.New("tabapayclient: duplicate accountID: 123 found but does not match current request info"),
		},
		{
			name: "successful match",
			curReq: &createFundingAccountPayload{
				ReferenceID: "1234567",
				Owner: &accountOwner{
					Name: &name{
						First: "Peter",
						Last:  "Parker",
					},
					Address: &address{
						Line1:   "123 place",
						City:    "city",
						State:   "CA",
						ZipCode: "90210",
						Country: usCountryCode,
					},
					Phone: &phone{
						Number: "5551234567",
					},
				},
			},
			setupMocks: func(mc *client.MockClient) {
				mc.EXPECT().
					DoJSON(gomock.Any(), gomock.Any()).
					DoAndReturn(func(req *http.Request, res *getFundingAccountResponse) error {
						*res = getFundingAccountResponse{
							ReferenceID: "1234567",
							Owner: accountOwner{
								Name: &name{
									First: "Peter",
									Last:  "Parker",
								},
								Address: &address{
									Line1:   "123 place",
									City:    "city",
									State:   "CA",
									ZipCode: "90210",
								},
								Phone: &phone{
									Number: "5551234567",
								},
							}}

						return nil
					}).
					Times(1)
			},
			expectedErr: nil,
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			mc := client.NewMockClient(ctrl)
			tc.setupMocks(mc)

			tbp := &TabaPay{
				httpClient: mc,
				refIDGen:   newReferenceID(),
			}

			err := tbp.duplicateMatchCurrentReq(context.Background(), "123", tc.curReq)

			if tc.expectedErr == nil {
				assert.NilError(t, err)
				return
			}

			assert.ErrorContains(t, err, tc.expectedErr.Error())
		})
	}
}

func TestValidateCard(t *testing.T) {
	tt := []struct {
		name        string
		info        *acquiring.AVSCVVInfo
		setupMocks  func(mc *client.MockClient)
		expectedErr error
	}{
		{
			name: "request error",
			info: &acquiring.AVSCVVInfo{
				PostalCodeZip: "90210",
				Country:       "US",
			},
			setupMocks: func(mc *client.MockClient) {
				mc.EXPECT().
					DoJSON(gomock.Any(), gomock.Any()).
					Return(errors.New("network err")).
					Times(1)
			},
			expectedErr: errors.New("tabapayclient: error making request"),
		},
		{
			name: "address format error",
			info: &acquiring.AVSCVVInfo{
				PostalCodeZip: "90210",
				Country:       "CA",
			},
			setupMocks: func(mc *client.MockClient) {
				mc.EXPECT().DoJSON(gomock.Any(), gomock.Any()).Times(0)
			},
			expectedErr: errors.New("tabapayclient: address format error"),
		},
		{
			name: "cvv does not match error",
			info: &acquiring.AVSCVVInfo{
				PostalCodeZip: "90210",
				Country:       "US",
			},
			setupMocks: func(mc *client.MockClient) {
				mc.EXPECT().
					DoJSON(gomock.Any(), gomock.Any()).
					DoAndReturn(func(req *http.Request, res *queryCardResponse) error {
						res.AVS.CodeSecurityCode = cvvNotMatched
						return nil
					}).
					Times(1)
			},
			expectedErr: errors.New("tabapayclient: CVV check failed"),
		},
		{
			name: "avs unavailable error",
			info: &acquiring.AVSCVVInfo{
				PostalCodeZip: "90210",
				Country:       "US",
			},
			setupMocks: func(mc *client.MockClient) {
				mc.EXPECT().
					DoJSON(gomock.Any(), gomock.Any()).
					DoAndReturn(func(req *http.Request, res *queryCardResponse) error {
						res.AVS.CodeSecurityCode = cvvMatched
						res.AVS.CodeAVS = avsUnavailable
						return nil
					}).
					Times(1)
			},
			expectedErr: errors.New("tabapayclient: AVS check failed"),
		},
		{
			name: "avs no info error",
			info: &acquiring.AVSCVVInfo{
				PostalCodeZip: "90210",
				Country:       "US",
			},
			setupMocks: func(mc *client.MockClient) {
				mc.EXPECT().
					DoJSON(gomock.Any(), gomock.Any()).
					DoAndReturn(func(req *http.Request, res *queryCardResponse) error {
						res.AVS.CodeSecurityCode = cvvMatched
						res.AVS.CodeAVS = avsNoInfo
						return nil
					}).
					Times(1)
			},
			expectedErr: errors.New("tabapayclient: AVS check failed"),
		},
		{
			name: "avs zip address not matched error",
			info: &acquiring.AVSCVVInfo{
				PostalCodeZip: "90210",
				Country:       "US",
			},
			setupMocks: func(mc *client.MockClient) {
				mc.EXPECT().
					DoJSON(gomock.Any(), gomock.Any()).
					DoAndReturn(func(req *http.Request, res *queryCardResponse) error {
						res.AVS.CodeSecurityCode = cvvMatched
						res.AVS.CodeAVS = avsZipAddressNotMatched
						return nil
					}).
					Times(1)
			},
			expectedErr: errors.New("tabapayclient: AVS check failed"),
		},
		{
			name: "avs only address match error",
			info: &acquiring.AVSCVVInfo{
				PostalCodeZip: "90210",
				Country:       "US",
			},
			setupMocks: func(mc *client.MockClient) {
				mc.EXPECT().
					DoJSON(gomock.Any(), gomock.Any()).
					DoAndReturn(func(req *http.Request, res *queryCardResponse) error {
						res.AVS.CodeSecurityCode = cvvMatched
						res.AVS.CodeAVS = avsAddressMatch
						return nil
					}).
					Times(1)
			},
			expectedErr: errors.New("tabapayclient: AVS check failed"),
		},
		{
			name: "avs only zip match error",
			info: &acquiring.AVSCVVInfo{
				PostalCodeZip: "90210",
				Country:       "US",
			},
			setupMocks: func(mc *client.MockClient) {
				mc.EXPECT().
					DoJSON(gomock.Any(), gomock.Any()).
					DoAndReturn(func(req *http.Request, res *queryCardResponse) error {
						res.AVS.CodeSecurityCode = cvvMatched
						res.AVS.CodeAVS = avsZipMatch
						return nil
					}).
					Times(1)
			},
			expectedErr: errors.New("tabapayclient: AVS check failed"),
		},
		{
			name: "success - avs zip and address match",
			info: &acquiring.AVSCVVInfo{
				PostalCodeZip: "90210",
				Country:       "US",
			},
			setupMocks: func(mc *client.MockClient) {
				mc.EXPECT().
					DoJSON(gomock.Any(), gomock.Any()).
					DoAndReturn(func(req *http.Request, res *queryCardResponse) error {
						res.AVS.CodeSecurityCode = cvvMatched
						res.AVS.CodeAVS = avsZipAddressMatch
						return nil
					}).
					Times(1)
			},
			expectedErr: nil,
		},
		{
			name: "success - avs zip and address match mastercard and discover",
			info: &acquiring.AVSCVVInfo{
				PostalCodeZip: "90210",
				Country:       "US",
			},
			setupMocks: func(mc *client.MockClient) {
				mc.EXPECT().
					DoJSON(gomock.Any(), gomock.Any()).
					DoAndReturn(func(req *http.Request, res *queryCardResponse) error {
						res.AVS.CodeSecurityCode = cvvMatched
						res.AVS.CodeAVS = avsMCDiscZipAddressMatch
						return nil
					}).
					Times(1)
			},
			expectedErr: nil,
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			mc := client.NewMockClient(ctrl)
			tc.setupMocks(mc)

			tbp := &TabaPay{
				httpClient: mc,
				baseURL:    "https://example.com",
			}

			err := tbp.ValidateCard(context.Background(), tc.info)

			if tc.expectedErr == nil {
				assert.NilError(t, err)
				return
			}

			assert.ErrorContains(t, err, tc.expectedErr.Error())
		})
	}
}

func TestReferenceID(t *testing.T) {
	tt := []struct {
		name          string
		numIDs        int
		expectedCount int
	}{
		{
			name:          "1 id",
			numIDs:        1,
			expectedCount: 1,
		},
		{
			name:          "a few ids",
			numIDs:        5,
			expectedCount: 5,
		},
		{
			name:          "> 100 ids",
			numIDs:        101,
			expectedCount: 1,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			refIDGen := newReferenceID()

			for i := 0; i < tc.numIDs; i++ {
				id := refIDGen.Next()

				assert.Equal(t, refIDGen.count, (i+1)%maxRefID)
				assert.Equal(t, len(id), 15)
				assert.Assert(t, strings.HasSuffix(id, fmt.Sprintf("%02d", i%maxRefID)))
			}

			assert.Equal(t, refIDGen.count, tc.expectedCount)
		})
	}
}

// Package tabapay is the implementation of the TabaPay acquiring processor
package tabapay

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Rhymond/go-money"
	"github.com/wingocard/braavos/internal/legacy/service/acquiringsvc"
	"github.com/wingocard/braavos/internal/wlog"
	"github.com/wingocard/braavos/pkg/client"
	"github.com/wingocard/braavos/pkg/client/middleware"
	"github.com/wingocard/braavos/pkg/entities/acquiring"
)

const (
	// HeaderAuthorizationKey is the header key for the Bearer Token.
	HeaderAuthorizationKey = "Authorization"
	// HeaderAuthorizationValue is the prefix for the Bearer value.
	HeaderAuthorizationValue = "Bearer "
	// ISO3166 numeric country codes
	usCountryCode = "840"
	caCountryCode = "124"

	postalCodeLength = 6

	// ISO8583 network response codes
	nrcApproved          = "00"
	nrcNoReasonToDecline = "85"

	// AVS and CVV checks
	cvvMatched               = "M"
	cvvNotMatched            = "N"
	avsZipAddressMatch       = "Y"
	avsZipMatch              = "Z"
	avsZipAddressNotMatched  = "N"
	avsAddressMatch          = "A"
	avsNoInfo                = "U"
	avsUnavailable           = "R"
	avsMCDiscZipAddressMatch = "X" // full match only for Mastercard and Discover
)

// Config is the struct that represents the configuration needed to setup
// a functioning TabaPay Client.
type Config struct {
	BearerToken         string
	ClientID            string
	SubClientID         string
	SettlementAccountID string
	BaseURL             string
	PublicKeyID         string
	PublicKey           string
	AllowDuplicateCards bool
	DebugRequests       bool
	DebugLogger         wlog.Logger
}

// TabaPay represents a TabaPay Acquiring Processor implementation that adheres
// to the AcquiringAccountProvider interface.
type TabaPay struct {
	httpClient          client.Client
	baseURL             string
	bearerToken         string
	settlementAccountID string
	clientID            string
	subClientID         string
	publicKeyID         string
	publicKey           string
	refIDGen            refIDGen
	allowDuplicateCards bool
}

// New will return a reference to a TabaPay client.
func New(conf *Config) *TabaPay {
	opts := []client.Option{}
	if conf.DebugRequests && conf.DebugLogger != nil {
		opts = append(opts, middleware.WithRequestLogger(conf.DebugLogger, conf.DebugRequests))
	}

	return &TabaPay{
		httpClient:          client.NewHTTPClient(opts...),
		baseURL:             conf.BaseURL,
		bearerToken:         conf.BearerToken,
		settlementAccountID: conf.SettlementAccountID,
		clientID:            conf.ClientID,
		subClientID:         conf.SubClientID,
		publicKeyID:         conf.PublicKeyID,
		publicKey:           conf.PublicKey,
		refIDGen:            newReferenceID(),
		allowDuplicateCards: conf.AllowDuplicateCards,
	}
}

func (t *TabaPay) buildRequest(
	ctx context.Context, method string, url string, body interface{},
) (*http.Request, error) {
	// WIN-1426 Disable TabaPay for program wind down
	return nil, fmt.Errorf("tabapay has been disabled for program wind down")
}

// CreateFundingAccount calls the create account endpoint on Tabapay.
// It updates the AccountID of the funding account on a
// successful call. It returns an error otherwise.
func (t *TabaPay) CreateFundingAccount(ctx context.Context, wl wlog.Logger, fa *acquiring.FundingAccount) error {
	cardDetails, ok := fa.AccountDetails.(*acquiring.Card)
	if !ok {
		return errors.New("tabapayclient: account details are not of type *acquiring.Card")
	}

	addr := &address{
		Line1:   fa.AddressLine1,
		City:    fa.City,
		State:   fa.ProvState,
		ZipCode: fa.PostalCodeZip,
		Country: fa.Country,
	}
	if fa.AddressLine2 != nil {
		addr.Line2 = *fa.AddressLine2
	}
	if err := addr.Format(); err != nil {
		return fmt.Errorf("tabapayclient: address format error: %s", err)
	}

	reqData := &createFundingAccountPayload{
		ReferenceID: t.refIDGen.Next(),
		Card: &encryptedCard{
			KeyID: cardDetails.EncryptionKeyID,
			Data:  cardDetails.EncryptedData,
		},
		Owner: &accountOwner{
			Name: &name{
				First: fa.FirstName,
				Last:  fa.LastName,
			},
			Address: addr,
			Phone: &phone{
				Number: strings.TrimPrefix(fa.Phone, "+1"),
			},
		},
	}
	if fa.AddressLine2 != nil {
		reqData.Owner.Address.Line2 = *fa.AddressLine2
	}

	urlSuffix := "/accounts?RejectDuplicateCard"
	if t.allowDuplicateCards {
		urlSuffix = "/accounts"
	}
	url := "/v1/clients/" + t.clientID + urlSuffix

	req, err := t.buildRequest(ctx, http.MethodPost, url, reqData)
	if err != nil {
		return fmt.Errorf("tabapayclient: error creating request: %w", err)
	}

	res := &createFundingAccountResponse{}
	if err := t.httpClient.DoJSON(req, res); err != nil {
		var httpErr *client.HTTPError
		if !errors.As(err, &httpErr) || httpErr.StatusCode != http.StatusConflict {
			// unexpected error
			return fmt.Errorf("tabapayclient: error making request: %w", err)
		}

		// was a conflict which indicates a duplicate account was found
		wl.Info("tabapayclient: found duplicate account on create")

		// get response body from original request
		if err := json.Unmarshal([]byte(httpErr.Body), res); err != nil {
			return fmt.Errorf("tabapayclient: duplicate found error: %w", err)
		}

		// should have one duplicate found
		if len(res.DuplicateAccountIds) != 1 {
			return fmt.Errorf(
				"tabapayclient: expected 1 duplicate account to be found got %d duplicates",
				len(res.DuplicateAccountIds),
			)
		}
		duplicateID := res.DuplicateAccountIds[0]

		// check if duplicate matches current request
		err := t.duplicateMatchCurrentReq(ctx, duplicateID, reqData)
		if err != nil {
			return err
		}

		// account is duplicate of this request - ok
		res.AccountID = duplicateID
	}

	// update funding account properties
	fa.ExternalAccountID = res.AccountID

	return nil
}

// UpdateFundingAccount calls the update account endpoint on Tabapay.
func (t *TabaPay) UpdateFundingAccount(
	ctx context.Context, accountID string, fa *acquiring.FundingAccount,
) error {
	if accountID == "" {
		return fmt.Errorf("tabapayclient: missing account id")
	}
	cardDetails, ok := fa.AccountDetails.(*acquiring.Card)
	if !ok {
		return errors.New("tabapayclient: account details are not of type *acquiring.Card")
	}

	addr := &address{
		Line1:   fa.AddressLine1,
		City:    fa.City,
		State:   fa.ProvState,
		ZipCode: fa.PostalCodeZip,
		Country: fa.Country,
	}
	if fa.AddressLine2 != nil {
		addr.Line2 = *fa.AddressLine2
	}
	if err := addr.Format(); err != nil {
		return fmt.Errorf("tabapayclient: address format error: %s", err)
	}

	reqData := &updateFundingAccountPayload{
		Card: &encryptedCard{
			KeyID: cardDetails.EncryptionKeyID,
			Data:  cardDetails.EncryptedData,
		},
		Owner: &accountOwner{
			Name: &name{
				First: fa.FirstName,
				Last:  fa.LastName,
			},
			Address: addr,
			Phone: &phone{
				Number: strings.TrimPrefix(fa.Phone, "+1"),
			},
		},
	}
	if fa.AddressLine2 != nil {
		reqData.Owner.Address.Line2 = *fa.AddressLine2
	}

	urlSuffix := "/accounts/" + accountID
	url := "/v1/clients/" + t.clientID + urlSuffix

	req, err := t.buildRequest(ctx, http.MethodPut, url, reqData)
	if err != nil {
		return fmt.Errorf("tabapayclient: error creating request: %w", err)
	}

	res := &updateFundingAccountResponse{}
	if err := t.httpClient.DoJSON(req, res); err != nil {
		return fmt.Errorf("tabapayclient: error making request: %w", err)
	}

	return nil
}

// DeleteFundingAccount calls the delete account endpoint on Tabapay.
func (t *TabaPay) DeleteFundingAccount(ctx context.Context, accountID string) error {
	urlSuffix := "/accounts/" + accountID + "?DeleteDuplicateCard"
	url := "/v1/clients/" + t.clientID + urlSuffix
	req, err := t.buildRequest(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return fmt.Errorf("tabapayclient: error creating request: %w", err)
	}

	res := &deleteFundingAccountResponse{}
	if err := t.httpClient.DoJSON(req, res); err != nil {
		return fmt.Errorf("tabapayclient: error making request: %w", err)
	}

	return nil
}

// Formats the amount into a string that TabaPay is expecting.
// ex: 100.00 .
func formatAmount(amount int64) string {
	fraction := 2   // the amount of decimal places
	decimal := "."  // the character to use for decimal notation
	template := "1" // go-money template
	// the rest of the parameters are intentionally left blank
	return money.NewFormatter(fraction, decimal, "", "", template).Format(amount)
}

// FundBankAccount calls the create transaction endpoint on Tabapay.
func (t *TabaPay) FundBankAccount(
	ctx context.Context, srcAccountID string, amount *money.Money,
) acquiring.FundAccountResponse {
	reqData := &createTransactionPayload{
		ReferenceID: t.refIDGen.Next(),
		Type:        acquiring.Pull,
		Accounts: &accounts{
			SourceAccountID:      srcAccountID,
			DestinationAccountID: t.settlementAccountID,
		},
		Currency: amount.Currency().NumericCode,
		Amount:   formatAmount(amount.Amount()),
		Memo:     "", // TODO: figure out what to put here
	}

	fundingResponse := acquiring.FundAccountResponse{
		Timestamp:        time.Now(),
		FundingAccountID: srcAccountID,
		Amount:           amount,
		ReferenceID:      reqData.ReferenceID,
	}

	urlSuffix := "/transactions"
	url := "/v1/clients/" + t.clientID + "_" + t.subClientID + urlSuffix

	req, err := t.buildRequest(ctx, http.MethodPost, url, reqData)
	if err != nil {
		err = fmt.Errorf("tabapayclient: error creating request: %w", err)
		return fundingResponse.WithErr(err)
	}

	res := &createTransactionResponse{}
	err = t.httpClient.DoJSON(req, res)
	if err != nil {
		var httpErr *client.HTTPError
		if errors.As(err, &httpErr) {
			fundingResponse.RawResponse = httpErr.Body
		}
		err = fmt.Errorf("tabapayclient: error making request: %w", err)
		return fundingResponse.WithErr(err)
	}

	rawResponse, marshallErr := json.Marshal(res)
	fundingResponse.TransactionID = res.TransactionID
	if marshallErr == nil {
		fundingResponse.RawResponse = string(rawResponse)
	}

	// check network response codes
	if res.NetworkRC != nrcApproved && res.NetworkRC != nrcNoReasonToDecline {
		err = &client.FundAccountError{
			NetworkRC: res.NetworkRC,
			Err: fmt.Errorf(
				"tabapayclient: transaction error (transactionID: %s), received network response code %s",
				res.TransactionID,
				res.NetworkRC,
			),
		}
	}
	return fundingResponse.WithErr(err)
}

// ValidateCard calls the query card endpoint on TabaPay
// to perform AVS and CVV checks on a card.
func (t *TabaPay) ValidateCard(ctx context.Context, info *acquiring.AVSCVVInfo) error {
	addr := &address{
		Line1:   info.AddressLine1,
		City:    info.City,
		State:   info.ProvState,
		ZipCode: info.PostalCodeZip,
		Country: info.Country,
	}
	if info.AddressLine2 != nil {
		addr.Line2 = *info.AddressLine2
	}
	if err := addr.Format(); err != nil {
		return fmt.Errorf("tabapayclient: address format error: %s", err)
	}

	reqData := &queryCardRequest{
		Card: &encryptedCard{
			KeyID: info.EncryptionKeyID,
			Data:  info.EncryptedCVVData,
		},
		Owner: &accountOwner{
			Name: &name{
				First: info.FirstName,
				Last:  info.LastName,
			},
			Address: addr,
			Phone: &phone{
				Number: strings.TrimPrefix(info.Phone, "+1"),
			},
		},
	}

	url := "/v1/clients/" + t.clientID + "_" + t.subClientID + "/cards?AVS"
	req, err := t.buildRequest(ctx, http.MethodPost, url, reqData)
	if err != nil {
		return fmt.Errorf("tabapayclient: error creating request: %w", err)
	}

	res := &queryCardResponse{}
	if err := t.httpClient.DoJSON(req, res); err != nil {
		return fmt.Errorf("tabapayclient: error making request: %w", err)
	}

	// check CVV response
	if res.AVS.CodeSecurityCode != cvvMatched {
		return &acquiringsvc.AVSCVVValidationError{
			Err: fmt.Errorf("tabapayclient: CVV check failed, response: %+v", res.AVS),
		}
	}

	// check AVS response
	if res.AVS.CodeAVS != avsZipAddressMatch && res.AVS.CodeAVS != avsMCDiscZipAddressMatch {
		return &acquiringsvc.AVSCVVValidationError{
			Err: fmt.Errorf("tabapayclient: AVS check failed, response: %+v", res.AVS),
		}
	}

	return nil
}

// GetEncryptionKey returns the currently configured TabaPay encryption key
// to use for encrypting credit/debit card details.
func (t *TabaPay) GetEncryptionKey(ctx context.Context) *acquiring.EncryptionKey {
	return &acquiring.EncryptionKey{
		ID:  t.publicKeyID,
		Key: t.publicKey,
	}
}

func (t *TabaPay) getFundingAccount(
	ctx context.Context, accountID string,
) (*getFundingAccountResponse, error) {
	req, err := t.buildRequest(ctx, http.MethodGet, "/accounts/"+accountID, nil)
	if err != nil {
		return nil, fmt.Errorf("tabapayclient: error creating request: %w", err)
	}

	res := &getFundingAccountResponse{}
	if err := t.httpClient.DoJSON(req, res); err != nil {
		return nil, fmt.Errorf("tabapayclient: error making request: %w", err)
	}

	return res, nil
}

func (t *TabaPay) duplicateMatchCurrentReq(
	ctx context.Context, duplicateID string, curReq *createFundingAccountPayload,
) error {
	res, err := t.getFundingAccount(ctx, duplicateID)
	if err != nil {
		return fmt.Errorf(
			"tabapayclient: error retrieving duplicate account with id: %s, %s",
			duplicateID,
			err,
		)
	}

	if res.Owner.Name == nil || res.Owner.Address == nil || res.Owner.Phone == nil {
		return fmt.Errorf(
			"tabapayclient: duplicate accountID: %s found but owner name, address, or phone was nil",
			duplicateID,
		)
	}

	// match response with original request data
	// fix response address country since tabapay doesn't return
	// a country for the US.
	if res.Owner.Address.Country == "" {
		res.Owner.Address.Country = usCountryCode
	}

	switch {
	case
		res.ReferenceID != curReq.ReferenceID,
		res.Owner.Name.First != curReq.Owner.Name.First,
		res.Owner.Name.Last != curReq.Owner.Name.Last,
		res.Owner.Phone.Number != curReq.Owner.Phone.Number,
		res.Owner.Address.Line1 != curReq.Owner.Address.Line1,
		res.Owner.Address.Line2 != curReq.Owner.Address.Line2,
		res.Owner.Address.City != curReq.Owner.Address.City,
		res.Owner.Address.State != curReq.Owner.Address.State,
		res.Owner.Address.ZipCode != curReq.Owner.Address.ZipCode,
		res.Owner.Address.Country != curReq.Owner.Address.Country:
		return fmt.Errorf(
			"tabapayclient: duplicate accountID: %s found but does not match current request info",
			duplicateID,
		)
	}

	return nil
}

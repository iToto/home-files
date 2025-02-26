// Package legacyidologyclient is a wrapper over idology rest api
package legacyidologyclient

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"github.com/wingocard/braavos/internal/wlog"
	"github.com/wingocard/braavos/pkg/client"
	"github.com/wingocard/braavos/pkg/client/middleware"
	"github.com/wingocard/braavos/pkg/client/tracking"
	"github.com/wingocard/braavos/pkg/entities/kyc"
	"github.com/wingocard/braavos/pkg/entities/user"
)

const (
	// baseURL
	baseURL = "localhost" // WIN-1425 Disable KYC for program wind down
	// VerifyIdentityRPC is the suffix of the URL to verify user identity
	verifyIdentityRPC string = "/api/idiq.svc"
	// idScanRPC is the suffix of the URL to get id scan result
	idScanRPC string = "/api/idscan.svc"
	// idScanImageRPC is the suffix of the URL to download documents
	idScanImageRPC string = "/api/idscanimage.svc"
	ssnPlaceHolder string = `<ssn>xx-xxx-xxxx</ssn>`
	ssnLength      int    = 9
)

var ssnRegex = regexp.MustCompile(`<ssn>\d{5}xxxx</ssn>`)

// Config is the struct that represents the configuration needed to setup
// a functioning idology Client
type Config struct {
	Username      string
	Password      string
	DebugRequests bool
	DebugLogger   wlog.Logger
}

// Client is an http client to communicate with Client api
type Client struct {
	username string
	password string
	client   *client.HTTPClient
}

// New will return a reference to a Client client
func New(conf Config) *Client {
	opts := []client.Option{}
	if conf.DebugRequests && conf.DebugLogger != nil {
		opts = append(opts, middleware.WithRequestLogger(conf.DebugLogger, conf.DebugRequests))
	}

	return &Client{
		username: conf.Username,
		password: conf.Password,
		client: client.NewHTTPClient(
			opts...,
		),
	}
}

type expectIDRequest struct {
	Username  string `url:"username"`
	Password  string `url:"password"`
	FirstName string `url:"firstName"`
	LastName  string `url:"lastName"`
	Email     string `url:"email"`
	Address   string `url:"address"`
	State     string `url:"state"`
	City      string `url:"city"`
	Country   string `url:"country"`
	Zip       string `url:"zip"`
	SSN       string `url:"ssn,omitempty"`
	SSN4      string `url:"ssnLast4,omitempty"`
	DobYear   int    `url:"dobYear"`
	DobMonth  int    `url:"dobMonth"`
	DobDay    int    `url:"dobDay"`
}

type expectIDResponse struct {
	XMLName       xml.Name `xml:"response"`
	Text          string   `xml:",chardata"`
	IDNumber      string   `xml:"id-number"`
	SummaryResult struct {
		Text    string `xml:",chardata"`
		Key     string `xml:"key"`
		Message string `xml:"message"`
	} `xml:"summary-result"`
	Results struct {
		Text    string `xml:",chardata"`
		Key     string `xml:"key"`
		Message string `xml:"message"`
	} `xml:"results"`
	Qualifiers struct {
		Text      string `xml:",chardata"`
		Qualifier []struct {
			Text    string `xml:",chardata"`
			Key     string `xml:"key"`
			Message string `xml:"message"`
		} `xml:"qualifier"`
	} `xml:"qualifiers"`
	LocatedRecord struct {
		SSN string `xml:"ssn"`
	} `xml:"located-record"`
	IDScan    string `xml:"id-scan"`
	IDScanURL string `xml:"id-scan-url"`
	Failed    string `xml:"failed"`
}

// VerifyIdentity verifies the identify of a given user
func (c *Client) VerifyIdentity(
	ctx context.Context,
	tracking tracking.Client,
	u *user.User,
	address *user.Address,
	ident *user.Identification,
) (*kyc.IDVerification, error) {
	reqBody := &expectIDRequest{
		Username:  c.username,
		Password:  c.password,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
		Address:   address.AddressLine1,
		State:     address.ProvState,
		City:      address.City,
		Country:   address.Country,
		Zip:       address.PostalCodeZip,
		DobYear:   u.DateOfBirth.Year(),
		DobMonth:  int(u.DateOfBirth.Month()),
		DobDay:    u.DateOfBirth.Day(),
	}

	switch ident.Type {
	case user.SSN4:
		reqBody.SSN4 = ident.ID
	case user.SSN:
		reqBody.SSN = strings.ReplaceAll(ident.ID, "-", "")
	default:
		return nil, fmt.Errorf("idologyclient: unsupported IDType: %d", ident.Type)
	}

	if address.AddressLine2 != nil {
		reqBody.Address = fmt.Sprintf("%s %s", *address.AddressLine2, address.AddressLine1)
	}

	req, err := client.NewURLEncodeRequest(ctx, http.MethodPost, baseURL+verifyIdentityRPC, reqBody)
	if err != nil {
		return nil,
			fmt.Errorf("idologyclient: error building request %s: %s", verifyIdentityRPC, err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("idologyclient: error calling endpoint %s: %s", verifyIdentityRPC, err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if err = c.recordKycResponse(ctx, tracking, u.ID, string(body)); err != nil {
		return nil, err
	}

	// unmarshall response
	var expectID expectIDResponse
	if err := xml.Unmarshal(body, &expectID); err != nil {
		return nil, err
	}
	// check for errors
	if expectID.Failed != "" {
		return nil, fmt.Errorf("idologyclient: response error: %s", expectID.Failed)
	}

	// hard fail if no match was located, in these cases we can not get the rest of the SSN
	if expectID.Results.Key != "result.match" && ident.Type == user.SSN4 {
		kycErrors := append(parseErrors(expectID), kyc.Error{
			Key:     expectID.Results.Key,
			Message: expectID.Results.Message,
		})

		return &kyc.IDVerification{
			UserID:   u.ID,
			VendorID: expectID.IDNumber,
			Status:   kyc.Fail,
			Errors:   kycErrors,
		}, nil
	}

	// if we only have the last 4 pull the rest from the reply
	if ident.Type == user.SSN4 {
		ident.ID = strings.ReplaceAll(expectID.LocatedRecord.SSN, "xxxx", ident.ID)
		ident.Type = user.SSN

		if len(ident.ID) != ssnLength {
			return nil, fmt.Errorf("idologyclient: SSN length missmatch")
		}
	}

	var verification kyc.IDVerification
	switch {
	case expectID.IDScan == "yes":
		verification = kyc.IDVerification{
			UserID:       u.ID,
			VendorID:     expectID.IDNumber,
			Status:       kyc.SoftFail,
			DocUploadURL: expectID.IDScanURL,
			Errors:       parseErrors(expectID),
			// save the full SSN to use after they pass ID check
			SSN: ident.ID,
		}
	case expectID.SummaryResult.Key == "id.failure":
		verification = kyc.IDVerification{
			UserID:   u.ID,
			VendorID: expectID.IDNumber,
			Status:   kyc.Fail,
			Errors:   parseErrors(expectID),
		}
	case expectID.SummaryResult.Key == "id.success":
		verification = kyc.IDVerification{
			UserID:   u.ID,
			VendorID: expectID.IDNumber,
			Status:   kyc.Pass,
		}
	default:
		return nil, errors.New("unknown kyc status")
	}

	return &verification, nil
}

func parseErrors(resp expectIDResponse) []kyc.Error {
	errors := make([]kyc.Error, len(resp.Qualifiers.Qualifier))
	for i, q := range resp.Qualifiers.Qualifier {
		errors[i] = kyc.Error{
			Key:     q.Key,
			Message: q.Message,
		}
	}
	return errors
}

type verifyResultsReq struct {
	Username string `url:"username"`
	Password string `url:"password"`
	QueryID  string `url:"queryId"`
}

type verifyResultResponse struct {
	XMLName  xml.Name `xml:"response"`
	Text     string   `xml:",chardata"`
	IDNumber string   `xml:"id-number"`
	Results  struct {
		Text    string `xml:",chardata"`
		Key     string `xml:"key"`
		Message string `xml:"message"`
	} `xml:"id-scan-result"`
	IDScanSummaryResult struct {
		Text    string `xml:",chardata"`
		Key     string `xml:"key"`
		Message string `xml:"message"`
	} `xml:"id-scan-summary-result"`
}

// GetIDScanResult returns ID scan (IDVerification) results for an existing soft fail IDVerification
func (c *Client) GetIDScanResult(
	ctx context.Context,
	tracking tracking.Client,
	oldID *kyc.IDVerification,
) (*kyc.IDVerification, error) {
	reqBody := &verifyResultsReq{
		Username: c.username,
		Password: c.password,
		QueryID:  oldID.VendorID,
	}
	req, err := client.NewURLEncodeRequest(ctx, http.MethodPost, baseURL+idScanRPC, reqBody)
	if err != nil {
		return nil,
			fmt.Errorf("idologyclient: error building request %s: %s", idScanRPC, err)
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("idologyclient: error calling endpoint %s: %s", idScanRPC, err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var v verifyResultResponse
	if err := xml.Unmarshal(body, &v); err != nil {
		return nil, err
	}

	if v.Results.Message != "Approved" && v.Results.Message != "ID Approved" {
		return oldID, nil
	}

	if err = c.recordKycResponse(ctx, tracking, oldID.UserID, string(body)); err != nil {
		return nil, err
	}

	switch v.IDScanSummaryResult.Message {
	case "PASS":
		return &kyc.IDVerification{
			VendorID: oldID.VendorID,
			UserID:   oldID.UserID,
			Status:   kyc.Pass,
		}, nil
	default:
		return oldID, nil
	}
}

type documentResultsReq struct {
	Username string `url:"username"`
	Password string `url:"password"`
	QueryID  string `url:"queryId"`
	ScanType string `url:"scanType"`
}

// GetDocuments download KYC scanned documents for a given vendorID
func (c *Client) GetDocuments(ctx context.Context, vendorID string) ([]*kyc.Document, error) {
	var docs []*kyc.Document
	for _, scanType := range []string{"first", "firstBack"} {
		reqBody := &documentResultsReq{
			Username: c.username,
			Password: c.password,
			QueryID:  vendorID,
			ScanType: scanType,
		}
		req, err := client.NewURLEncodeRequest(ctx, http.MethodPost, baseURL+idScanImageRPC, reqBody)
		if err != nil {
			return nil,
				fmt.Errorf("idologyclient: error building request %s: %s", verifyIdentityRPC, err)
		}
		resp, err := c.client.Do(req)
		if err != nil {
			// passports don't have back images, skip
			if scanType == "firstBack" &&
				errors.Is(err, &client.HTTPError{StatusCode: 400, Body: "Image Not Available"}) { // nolint
				continue
			}
			return nil, fmt.Errorf("idologyclient: error calling endpoint %s: %s", verifyIdentityRPC, err)
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		// IDology may return either an image or an XML response error in the body,
		if resp.Header.Get("Content-Type") != "image/jpg" {
			return nil, &client.HTTPError{
				URL:        req.URL.String(),
				StatusCode: resp.StatusCode,
				Body:       string(body),
			}
		}

		docs = append(docs, &kyc.Document{
			Key:     fmt.Sprintf("%s.jpg", scanType),
			Content: body,
		})
	}
	return docs, nil
}

func (c *Client) recordKycResponse(
	ctx context.Context,
	tracking tracking.Client,
	userID string,
	response string,
) error {
	// remove the first 5 of their SSN from the response we save
	response = ssnRegex.ReplaceAllString(response, ssnPlaceHolder)
	// send new response to segment for audit purposes
	return tracking.Track(
		ctx,
		userID,
		"KYC Response",
		map[string]interface{}{"content": response},
	)
}

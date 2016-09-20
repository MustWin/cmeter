package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/MustWin/cmeter/context"
	"github.com/MustWin/cmeter/reporting"
	"github.com/MustWin/cmeter/reporting/factory"
)

const (
	CLIENT_USER_AGENT     = "cmeter-http-reporter"
	CLIENT_VERSION_HEADER = "X-CMETER-VERSION"

	DEFAULT_RECEIPT_HEADER = "X-CMETER-RECEIPT"
)

var (
	ErrInvalidEndpoint = errors.New("invalid endpoint url")
	ErrInvalidReceipt  = errors.New("received an invalid or empty receipt")
	ErrInvalidHeaders  = errors.New("error reading additional header configuration")
)

type driverFactory struct{}

func (factory *driverFactory) Create(parameters map[string]interface{}) (reporting.Driver, error) {
	endpointUrl, ok := parameters["url"].(string)
	if !ok || endpointUrl == "" {
		return nil, ErrInvalidEndpoint
	}

	if _, err := url.Parse(endpointUrl); err != nil {
		return nil, fmt.Errorf("invalid endpoint url: %v", err)
	}

	httpMethod, ok := parameters["method"].(string)
	if !ok || httpMethod == "" {
		httpMethod = http.MethodPost
	}

	httpMethod = strings.ToUpper(httpMethod)

	receiptHeader, ok := parameters["receipt_header"].(string)
	if !ok || receiptHeader == "" {
		receiptHeader = DEFAULT_RECEIPT_HEADER
	}

	headers, ok := parameters["headers"].(http.Header)
	if !ok {
		return nil, ErrInvalidHeaders
	}

	return &Driver{
		Endpoint:       endpointUrl,
		Method:         httpMethod,
		ReceiptHeader:  receiptHeader,
		IdentityHeader: identityHeader,
		IdentityLabel:  identityLabel,
		ExtraHeaders:   headers,
	}, nil
}

func init() {
	factory.Register("http", &driverFactory{})
}

type Driver struct {
	Endpoint      string
	Method        string
	ReceiptHeader string
	ExtraHeaders  http.Header
}

func (d *Driver) Report(ctx context.Context, e *reporting.Event) (reporting.Receipt, error) {
	blob, err := json.Marshal(e)
	if err != nil {
		return reporting.EmptyReceipt, fmt.Errorf("error encoding event: %v", err)
	}

	r, err := http.NewRequest(d.Method, d.Endpoint, bytes.NewReader(blob))
	if err != nil {
		return reporting.EmptyReceipt, fmt.Errorf("error creating request: %v", err)
	}

	r.Header.Add("Content-Length", strconv.FormatInt(int64(len(blob)), 10))
	r.Header.Add("Content-Type", "application/json")

	version := context.GetVersion(ctx)
	r.Header.Add("User-Agent", fmt.Sprintf("%s/%s", CLIENT_USER_AGENT, version))
	r.Header.Add(CLIENT_VERSION_HEADER, version)
	for hn, hvs := range d.ExtraHeaders {
		for _, v := range hvs {
			r.Header.Add(hn, v)
		}
	}

	client := http.DefaultClient
	resp, err := client.Do(r)
	if err != nil {
		return reporting.EmptyReceipt, fmt.Errorf("error sending request: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode > 299 || resp.StatusCode < 200 {
		return reporting.EmptyReceipt, fmt.Errorf("unexpected response: %q", resp.Status)
	}

	receiptStr := resp.Header.Get(d.ReceiptHeader)
	if receiptStr == "" {
		return reporting.EmptyReceipt, ErrInvalidReceipt
	}

	return reporting.Receipt(receiptStr), nil
}

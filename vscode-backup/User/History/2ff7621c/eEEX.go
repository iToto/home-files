package okxapi

import "fmt"

// NoDataError is a custom error type for handling cases where no data is returned in a response.
type NoDataError struct {
	Endpoint string
}

// Error method to implement the error interface.
func (e *NoDataError) Error() string {
	return fmt.Sprintf("No data returned from endpoint: %s", e.Endpoint)
}

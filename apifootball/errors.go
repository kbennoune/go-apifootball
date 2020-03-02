package apifootball

import (
	"fmt"
	"strings"
)

// Error represents a Discogs API error
type Error struct {
	Message string
}

func (e *Error) Error() string {
	return fmt.Sprintf("api football error: %s", strings.ToLower(e.Message))
}

// APIErrors
var (
	ErrUnauthorized        = &Error{"authentication required"}
	ErrTimeout             = &Error{"timeout"}
	ErrInternalServerError = &Error{"internal server error"}
)

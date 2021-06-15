package messages

import (
	"fmt"
	"github.com/pkg/errors"
	"time"
)

func nowUTC() time.Time {
	loc, err := time.LoadLocation("UTC")
	if err != nil {
		panic(errors.WithStack(fmt.Errorf("failed to load UTC timezone: %w", err)))
	}
	return time.Now().In(loc).Round(time.Millisecond)
}

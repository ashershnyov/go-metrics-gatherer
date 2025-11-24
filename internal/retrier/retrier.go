package retrier

import (
	"errors"
	"time"
)

const sleep = 1 * time.Second

var ErrUnretriable = errors.New("Unretriable")

// WithRetry executes f until succeeds or maxRetries is reached.
func WithRetry(maxRetires int, f func() error) error {
	var err error
	for i := 0; i < maxRetires; i++ {
		err = f()

		if errors.Is(err, ErrUnretriable) {
			return err
		}

		if err == nil {
			break
		}

		time.Sleep(time.Second + sleep*time.Duration(i))
	}

	return err
}

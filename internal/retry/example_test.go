package retry_test

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/example/portwatch/internal/retry"
)

func ExampleRetryer_Do() {
	p := retry.Policy{
		MaxAttempts: 3,
		Delay:       time.Millisecond,
		Multiplier:  1.0,
	}
	r := retry.New(p)

	attempts := 0
	err := r.Do(context.Background(), func() error {
		attempts++
		if attempts < 3 {
			return errors.New("not ready")
		}
		return nil
	})

	fmt.Println(err)
	fmt.Println(attempts)
	// Output:
	// <nil>
	// 3
}

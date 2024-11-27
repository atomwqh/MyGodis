package timewheel

import (
	"testing"
	"time"
)

func TestDelay(t *testing.T) {
	ch := make(chan time.Time)

	beginTime := time.Now()
	Delay(time.Second, "", func() {
		ch <- time.Now()
	})
	execAt := <-ch
	delayDuration := execAt.Sub(beginTime)

	if delayDuration < time.Second || delayDuration > 3*time.Second {
		t.Error("wrong exec time")
	}
}

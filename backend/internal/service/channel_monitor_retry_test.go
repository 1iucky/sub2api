//go:build unit

package service

import (
	"testing"
	"time"
)

func TestRunMonitorCheckWithRetries_RetriesFailedResults(t *testing.T) {
	attempts := 0
	result := runMonitorCheckWithRetries(2, func() *CheckResult {
		attempts++
		return &CheckResult{
			Model:     "deepseek-v4-flash",
			Status:    MonitorStatusError,
			Message:   "upstream 502",
			CheckedAt: time.Now(),
		}
	})

	if attempts != 3 {
		t.Fatalf("expected first attempt plus 2 retries, got %d attempts", attempts)
	}
	if result.Status != MonitorStatusError {
		t.Fatalf("expected final error result, got %q", result.Status)
	}
}

func TestRunMonitorCheckWithRetries_StopsAfterSuccess(t *testing.T) {
	attempts := 0
	result := runMonitorCheckWithRetries(3, func() *CheckResult {
		attempts++
		if attempts == 1 {
			return &CheckResult{Model: "glm-5.2", Status: MonitorStatusFailed, Message: "challenge mismatch", CheckedAt: time.Now()}
		}
		return &CheckResult{Model: "glm-5.2", Status: MonitorStatusOperational, CheckedAt: time.Now()}
	})

	if attempts != 2 {
		t.Fatalf("expected retry to stop after successful result, got %d attempts", attempts)
	}
	if result.Status != MonitorStatusOperational {
		t.Fatalf("expected operational result, got %q", result.Status)
	}
}

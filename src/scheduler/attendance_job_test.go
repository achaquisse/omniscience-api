package scheduler

import (
	"testing"
)

func TestTriggerManualReport_NoError(t *testing.T) {
	err := TriggerManualReport()
	if err != nil {
		t.Errorf("Expected no error from TriggerManualReport, got %v", err)
	}
}

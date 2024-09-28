package ssv

import (
	"testing"
)

func TestEvent(t *testing.T) {
	events := GetAllSSVEvent()
	for _, event := range events {
		t.Log("event name", event.Name)
	}
	t.Log(len(events))
}

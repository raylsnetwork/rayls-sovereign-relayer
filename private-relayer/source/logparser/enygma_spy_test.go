package logparser_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/raylsnetwork/rayls-sovereign-relayer/private-relayer/source/service"
	"github.com/stretchr/testify/assert"
)

// eventVerifier helps verify that events are pushed with expected resource IDs and event types
// Each unique (resourceID, eventType) combination should be seen exactly once
type eventVerifier struct {
	t              *testing.T
	expectedEvents map[string]expectedEvent // key: "resourceID-eventType"
	seenEvents     map[string]bool          // key: "resourceID-eventType"
	eventCounter   int
}

type expectedEvent struct {
	resourceID      string
	eventType       service.EnygmaEventType
	serializedEvent []byte
}

func newEventVerifier(t *testing.T, eventType service.EnygmaEventType) *eventVerifier {
	return &eventVerifier{
		t:              t,
		expectedEvents: make(map[string]expectedEvent),
		seenEvents:     make(map[string]bool),
	}
}

// expectEvent registers an expected event for a given resource ID and event type
func (v *eventVerifier) expectEvent(resourceID string, eventType service.EnygmaEventType, event interface{}) {
	serialized, err := json.Marshal(event)
	assert.NoError(v.t, err)

	key := fmt.Sprintf("%s-%d", resourceID, eventType)
	v.expectedEvents[key] = expectedEvent{
		resourceID:      resourceID,
		eventType:       eventType,
		serializedEvent: serialized,
	}
}

// verifyBatch is called for each pushed batch and verifies all events match expectations
func (v *eventVerifier) verifyBatch(events []service.EnygmaSerializedEvent) error {
	for _, event := range events {
		v.eventCounter++

		// Check if this event is expected
		key := fmt.Sprintf("%s-%d", event.ResourceID, event.Type)
		expected, ok := v.expectedEvents[key]
		if !ok {
			assert.Fail(v.t, fmt.Sprintf("unexpected event: resourceID=%s, eventType=%d", event.ResourceID, event.Type))
			continue
		}

		// Mark as seen
		v.seenEvents[key] = true

		// Verify the serialized event matches
		// For transfer events, we need to ignore the MessageId field since it's dynamically generated
		if event.Type == service.EnygmaTransferEvent {
			// Unmarshal and compare without MessageId
			var expectedEvent, actualEvent map[string]interface{}
			if err := json.Unmarshal(expected.serializedEvent, &expectedEvent); err != nil {
				assert.Fail(v.t, fmt.Sprintf("failed to unmarshal expected event: %v", err))
				continue
			}
			if err := json.Unmarshal(event.SerializedEvent, &actualEvent); err != nil {
				assert.Fail(v.t, fmt.Sprintf("failed to unmarshal actual event: %v", err))
				continue
			}

			// Remove MessageId from both before comparison
			delete(expectedEvent, "MessageId")
			delete(actualEvent, "MessageId")
			assert.Equal(v.t, expectedEvent, actualEvent, "event mismatch for resourceID=%s, eventType=%d", event.ResourceID, event.Type)
		} else {
			assert.JSONEq(v.t, string(expected.serializedEvent), string(event.SerializedEvent))
		}
	}

	return nil
}

// assertAllSeen verifies that all expected events were pushed
func (v *eventVerifier) assertAllSeen() {
	assert.Equal(v.t, len(v.expectedEvents), v.eventCounter, "expected %d events, got %d", len(v.expectedEvents), v.eventCounter)

	for key, expected := range v.expectedEvents {
		assert.True(v.t, v.seenEvents[key], "expected event not seen: resourceID=%s, eventType=%d", expected.resourceID, expected.eventType)
	}
}

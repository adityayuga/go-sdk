package debug

import "encoding/json"

// TraceEnvelope wraps a message body with trace headers for message brokers
// that have no native header support (e.g. NSQ). The Envelope field is a
// discriminator: consumers check Envelope == true to distinguish wrapped
// messages from legacy plain payloads. The job handler always receives
// OriginalBody — the unwrapped payload.
type TraceEnvelope struct {
	Envelope     bool              `json:"_envelope"`     // always true; discriminator
	DebugID      string            `json:"debug_id"`      // propagated x-debug-id
	Headers      map[string]string `json:"headers"`       // all other x-* trace headers
	OriginalBody json.RawMessage   `json:"original_body"` // the original message payload
}

// WrapEnvelope encodes body inside a TraceEnvelope together with the provided
// trace headers. Returns the serialised envelope bytes ready to publish.
func WrapEnvelope(debugID string, headers map[string]string, body []byte) ([]byte, error) {
	env := TraceEnvelope{
		Envelope:     true,
		DebugID:      debugID,
		Headers:      headers,
		OriginalBody: json.RawMessage(body),
	}
	return json.Marshal(env)
}

// UnwrapEnvelope attempts to decode a TraceEnvelope from data.
// Returns (envelope, true) when data is a valid envelope, or (zero, false)
// when data is a legacy message that was not wrapped by the publisher.
func UnwrapEnvelope(data []byte) (TraceEnvelope, bool) {
	var env TraceEnvelope
	if err := json.Unmarshal(data, &env); err != nil {
		return TraceEnvelope{}, false
	}
	if !env.Envelope || len(env.OriginalBody) == 0 {
		return TraceEnvelope{}, false
	}
	return env, true
}

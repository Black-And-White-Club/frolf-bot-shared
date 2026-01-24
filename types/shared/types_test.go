package sharedtypes

import (
	"encoding/json"
	"testing"

	"github.com/google/uuid"
)

func TestRoundID_MarshalJSON(t *testing.T) {
	id := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
	roundID := RoundID(id)

	data, err := json.Marshal(roundID)
	if err != nil {
		t.Fatalf("MarshalJSON failed: %v", err)
	}

	expected := `"550e8400-e29b-41d4-a716-446655440000"`
	if string(data) != expected {
		t.Errorf("MarshalJSON = %s, want %s", string(data), expected)
	}
}

func TestRoundID_UnmarshalJSON_StringFormat(t *testing.T) {
	jsonData := []byte(`"550e8400-e29b-41d4-a716-446655440000"`)

	var roundID RoundID
	if err := json.Unmarshal(jsonData, &roundID); err != nil {
		t.Fatalf("UnmarshalJSON failed: %v", err)
	}

	expected := "550e8400-e29b-41d4-a716-446655440000"
	if roundID.String() != expected {
		t.Errorf("UnmarshalJSON = %s, want %s", roundID.String(), expected)
	}
}

func TestRoundID_RoundTrip(t *testing.T) {
	original := RoundID(uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"))

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("MarshalJSON failed: %v", err)
	}

	var decoded RoundID
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("UnmarshalJSON failed: %v", err)
	}

	if original.String() != decoded.String() {
		t.Errorf("Round trip failed: got %s, want %s", decoded.String(), original.String())
	}
}

func TestRoundID_InStruct(t *testing.T) {
	type JobArgs struct {
		GuildID string  `json:"guild_id"`
		RoundID RoundID `json:"round_id"`
	}

	job := JobArgs{
		GuildID: "123456789",
		RoundID: RoundID(uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")),
	}

	data, err := json.Marshal(job)
	if err != nil {
		t.Fatalf("Marshal struct failed: %v", err)
	}

	// Verify JSON contains string UUID, not byte array
	expected := `{"guild_id":"123456789","round_id":"550e8400-e29b-41d4-a716-446655440000"}`
	if string(data) != expected {
		t.Errorf("Marshal struct = %s, want %s", string(data), expected)
	}

	var decoded JobArgs
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Unmarshal struct failed: %v", err)
	}

	if decoded.RoundID.String() != job.RoundID.String() {
		t.Errorf("Struct round trip failed: got %s, want %s", decoded.RoundID.String(), job.RoundID.String())
	}
}

func TestEventMessageID_MarshalJSON(t *testing.T) {
	id := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
	eventID := EventMessageID(id)

	data, err := json.Marshal(eventID)
	if err != nil {
		t.Fatalf("MarshalJSON failed: %v", err)
	}

	expected := `"550e8400-e29b-41d4-a716-446655440000"`
	if string(data) != expected {
		t.Errorf("MarshalJSON = %s, want %s", string(data), expected)
	}
}

func TestEventMessageID_UnmarshalJSON(t *testing.T) {
	jsonData := []byte(`"550e8400-e29b-41d4-a716-446655440000"`)

	var eventID EventMessageID
	if err := json.Unmarshal(jsonData, &eventID); err != nil {
		t.Fatalf("UnmarshalJSON failed: %v", err)
	}

	expected := "550e8400-e29b-41d4-a716-446655440000"
	if eventID.String() != expected {
		t.Errorf("UnmarshalJSON = %s, want %s", eventID.String(), expected)
	}
}

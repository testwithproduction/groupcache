package testpb

import (
	"testing"
	"github.com/golang/protobuf/proto"
)

func TestTestMessageProto(t *testing.T) {
	tests := []struct {
		name string
		city string
	}{
		{"hello", "NYC"},
		{"", ""}, // Test empty values
		{"long name with spaces", "city with spaces"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := &TestMessage{
				Name: proto.String(tt.name),
				City: proto.String(tt.city),
			}
			data, err := proto.Marshal(msg)
			if err != nil {
				t.Fatalf("Marshal failed: %v", err)
			}
			var out TestMessage
			err = proto.Unmarshal(data, &out)
			if err != nil {
				t.Fatalf("Unmarshal failed: %v", err)
			}
			if out.GetName() != tt.name || out.GetCity() != tt.city {
				t.Errorf("unexpected values after roundtrip: got %+v, want %+v", out, msg)
			}
		})
	}
}

func TestEmptyProto(t *testing.T) {
	msg := &Empty{}
	data, err := proto.Marshal(msg)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}
	var out Empty
	err = proto.Unmarshal(data, &out)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}
}

func TestTestMessageNilFields(t *testing.T) {
	msg := &TestMessage{} // No fields set
	data, err := proto.Marshal(msg)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}
	var out TestMessage
	err = proto.Unmarshal(data, &out)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}
	if out.GetName() != "" || out.GetCity() != "" {
		t.Errorf("expected empty strings for nil fields, got %+v", out)
	}
}

func TestTestMessageEquality(t *testing.T) {
	msg1 := &TestMessage{Name: proto.String("test"), City: proto.String("city")}
	msg2 := &TestMessage{Name: proto.String("test"), City: proto.String("city")}
	msg3 := &TestMessage{Name: proto.String("different"), City: proto.String("city")}

	if !proto.Equal(msg1, msg2) {
		t.Error("identical messages should be equal")
	}
	if proto.Equal(msg1, msg3) {
		t.Error("different messages should not be equal")
	}
}

func TestTestMessageString(t *testing.T) {
	msg := &TestMessage{Name: proto.String("test"), City: proto.String("city")}
	str := msg.String()
	if str == "" {
		t.Error("String() should return non-empty string")
	}
}
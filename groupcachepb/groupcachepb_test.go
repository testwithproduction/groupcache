package groupcachepb

import (
	"testing"
	"github.com/golang/protobuf/proto"
)

func TestGetRequestProto(t *testing.T) {
	msg := &GetRequest{
		Group: proto.String("test-group"),
		Key:   proto.String("test-key"),
	}
	data, err := proto.Marshal(msg)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}
	var out GetRequest
	err = proto.Unmarshal(data, &out)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}
	if out.GetGroup() != "test-group" || out.GetKey() != "test-key" {
		t.Errorf("unexpected values after roundtrip: got %+v", out)
	}
}

func TestGetResponseProto(t *testing.T) {
	msg := &GetResponse{
		Value: []byte("hello world"),
	}
	data, err := proto.Marshal(msg)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}
	var out GetResponse
	err = proto.Unmarshal(data, &out)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}
	if string(out.Value) != "hello world" {
		t.Errorf("unexpected value after roundtrip: got %+v", out)
	}
}
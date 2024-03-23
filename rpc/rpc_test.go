package rpc_test

import (
	"langserver/rpc"
	"testing"
)

type encodingExample struct {
	Testing bool
}

func TestEncode(t *testing.T) {
	expected := "Content-Length: 16\r\n\r\n{\"Testing\":true}"
	actual := rpc.EncodeMessage(encodingExample{Testing: true})

	if expected != actual {
		t.Fatalf("Expected: %s, Actual: %s", expected, actual)
	}
}

func TestDecode(t *testing.T) {
	incomingMsg := "Content-Length: 15\r\n\r\n{\"Method\":\"hi\"}"
	expectedContentLen := 15
	expectedMethod := "hi"

	actualMethod, content, err := rpc.DecodeMessage([]byte(incomingMsg))
	if err != nil {
		t.Fatal(err)
	}

	if expectedContentLen != len(content) {
		t.Fatalf("Expected: %d, Actual: %d", expectedContentLen, len(content))
	}
	if expectedMethod != actualMethod {
		t.Fatalf("Expected: %s, Actual: %s", expectedMethod, actualMethod)
	}
}

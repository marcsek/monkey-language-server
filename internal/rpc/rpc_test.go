package rpc_test

import (
	"fmt"
	"testing"

	"github.com/marcsek/monkey-language-server/internal/rpc"
)

type EncodingExample struct {
	Method string
}

func TestEncode(t *testing.T) {
	expected := fmt.Sprintf("Content-Length: 15\r\n\r\n{\"Method\":\"ok\"}")
	result := rpc.EncodeMessage(EncodingExample{Method: "ok"})

	if expected != result {
		t.Fatalf("encoding failed, want=%s; got=%s", expected, result)
	}
}

func TestDecode(t *testing.T) {
	expectedMethod, expectedLen := "ok", 15
	method, content, err := rpc.DecodeMessage(
		[]byte("Content-Length: 15\r\n\r\n{\"Method\":\"ok\"}"),
	)
	contentLength := len(content)
	if err != nil {
		t.Fatal(err)
	}

	if expectedLen != contentLength {
		t.Fatalf("decoding failed, want=%d; got=%d", expectedLen, contentLength)
	}

	if expectedMethod != method {
		t.Fatalf("decoding failed, want=%s; got=%s", expectedMethod, method)
	}
}

func TestSplitFunc(t *testing.T) {
	type ExpectedOutput struct {
		Advance int
		Token   []byte
		Err     error
	}

	tests := []struct {
		Input  string
		Result ExpectedOutput
	}{
		{
			Input: "Content-Length: 15\r\n\r\n{\"Method\":\"ok\"}",
			Result: ExpectedOutput{
				Advance: len("Content-Length: 15\r\n\r\n{\"Method\":\"ok\"}"),
				Token:   []byte("Content-Length: 15\r\n\r\n{\"Method\":\"ok\"}"),
				Err:     nil,
			},
		},
		{
			Input: "Content-Length: 1d\r\n\r\n{\"Method\":\"ok\"}",
			Result: ExpectedOutput{
				Advance: 0,
				Token:   nil,
				Err:     fmt.Errorf(""),
			},
		},
		{
			Input: "",
			Result: ExpectedOutput{
				Advance: 0,
				Token:   nil,
				Err:     nil,
			},
		},
		{
			Input: "Content-Length: 1d\r\n",
			Result: ExpectedOutput{
				Advance: 0,
				Token:   nil,
				Err:     nil,
			},
		},
	}

	for _, test := range tests {
		advance, token, err := rpc.SplicFunc([]byte(test.Input), false)

		if advance != test.Result.Advance {
			t.Fatalf("In advance, want=%d; got=%d", test.Result.Advance, advance)
		}

		if len(token) != len(test.Result.Token) {
			t.Fatalf("Wrong token len, want=%d; got=%d", len(test.Result.Token), len(token))
		}

		if err != nil && test.Result.Err == nil {
			t.Fatalf("Unexpected error, got=%s", err)
		}

		if err == nil && test.Result.Err != nil {
			t.Fatalf("Expected error, got=%s", err)
		}
	}
}

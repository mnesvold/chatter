package main

import (
	"bytes"
	"reflect"
	"testing"
	"time"
)

type fakeRWC struct {
	input, output *bytes.Buffer
	isOpen        bool
}

func (rwc *fakeRWC) Close() error {
	rwc.isOpen = false
	return nil
}

func (rwc *fakeRWC) Read(p []byte) (n int, err error) {
	n, err = rwc.input.Read(p)
	return
}

func (rwc *fakeRWC) Write(p []byte) (n int, err error) {
	n, err = rwc.output.Write(p)
	return
}

func newFakeRWC() *fakeRWC {
	var input, output bytes.Buffer
	return &fakeRWC{&input, &output, true}
}

type testContext struct {
	client *client
	input  *bytes.Buffer
	output *bytes.Buffer
	send   chan []byte
	recv   chan map[string]interface{}
}

func newTestContext() *testContext {
	rwc := newFakeRWC()
	send := make(chan []byte)
	recv := make(chan map[string]interface{})
	client := newClient(rwc, send, recv)
	return &testContext{
		client: client,
		input:  rwc.input,
		output: rwc.output,
		send:   send,
		recv:   recv,
	}
}

func TestClientRead(t *testing.T) {
	expected := map[string]interface{}{
		"hello": "world",
	}
	ctx := newTestContext()
	_, err := ctx.input.Write([]byte(`{"hello":"world"}`))
	if err != nil {
		t.Fatal(err)
	}
	select {
	case actual := <-ctx.recv:
		if !reflect.DeepEqual(actual, expected) {
			t.Fatalf("Expected %q; recv'd %q instead", expected, actual)
		}
	case <-time.After(1 * time.Second):
		t.Fatal("timeout")
	}
}

func TestClientWrite(t *testing.T) {
	expected := []byte(`{"hello":"world"}`)
	ctx := newTestContext()
	select {
	case ctx.send <- expected:
		// continue the test
	case <-time.After(1 * time.Second):
		t.Fatal("timeout")
	}
}

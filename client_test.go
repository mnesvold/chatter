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
	rwc    *fakeRWC
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
		rwc:    rwc,
		input:  rwc.input,
		output: rwc.output,
		send:   send,
		recv:   recv,
	}
}

func (ctx *testContext) Close(t *testing.T) (err error) {
	err = ctx.client.Close()
	if err != nil {
		t.Error(err)
	}
	close(ctx.send)
	close(ctx.recv)
	return
}

func TestClientRead(t *testing.T) {
	expected := map[string]interface{}{
		"hello": "world",
	}
	ctx := newTestContext()
	defer ctx.Close(t)
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
	defer ctx.Close(t)
	ctx.client.Send(expected)

	bufferReady := make(chan bool)
	go func() {
		for ctx.output.Len() < len(expected) {
		}
		bufferReady <- true
	}()
	select {
	case <-bufferReady:
		actual := ctx.output.Bytes()
		if !bytes.Equal(actual, expected) {
			t.Fatalf("Expected %q; client sent %q instead", expected, actual)
		}
	case <-time.After(1 * time.Second):
		t.Fatal("timeout")
	}
}

func TestClientClose(t *testing.T) {
	ctx := newTestContext()
	if !ctx.rwc.isOpen {
		t.Fatal("fake RWC was not open to begin with -- bad test setup")
	}
	ctx.Close(t)
	if ctx.rwc.isOpen {
		t.Error("client did not close its RWC on close")
	}
}

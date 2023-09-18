package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func FuzzEnvHandler(f *testing.F) {

	// given
	testcases := []string{"test", "hello!", "$%#iwo"}
	for _, tc := range testcases {
		f.Add(tc)
	}

	f.Fuzz(func(t *testing.T, a string) {
		f := func() *string {
			return &a
		}
		// when
		srvr := httptest.NewServer(http.HandlerFunc(showMsg(f)))
		defer srvr.Close()

		resp, err := http.Get(srvr.URL)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
		var buf = new(bytes.Buffer)
		if _, err := io.Copy(buf, resp.Body); err != nil {
			t.Fatal(err)
		}

		// then
		got := buf.String()
		if got != a {
			t.Fatalf("wrong output. Expected: '%v' [%d], got: '%v'[%d]", a, len(a), got, len(got))
		}
	})
}

func TestHandler(t *testing.T) {

	// given
	v := "some test value from ENV"
	os.Setenv("MSG", v)

	srvr := httptest.NewServer(http.HandlerFunc(showMsg(fromEnv("MSG"))))
	defer srvr.Close()

	// when
	resp, err := http.Get(srvr.URL)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	// then
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("wrong status code. expected: %d, got: %d", http.StatusOK, resp.StatusCode)
	}

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, resp.Body); err != nil {
		t.Fatalf("could not read response body: %v", err)
	}

	if buf.String() != v {
		t.Fatalf("wrong response\nexpected: %s\ngot: %s", v, buf.String())
	}
}

func TestHandlerMissingENV(t *testing.T) {

	// given
	srvr := httptest.NewServer(http.HandlerFunc(showMsg(func() *string { return nil })))
	defer srvr.Close()

	// when
	resp, err := http.Get(srvr.URL)
	if err != nil {
		t.Fatalf("could not request HTTP: %v", err)
	}
	defer resp.Body.Close()

	// then
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("wrong status code. expected: %d, got: %d", http.StatusNotFound, resp.StatusCode)
	}
}

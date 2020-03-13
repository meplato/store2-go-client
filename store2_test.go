package store2_test

import (
	"bufio"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"strings"
	"testing"

	store2 "github.com/meplato/store2-go-client"
)

func getService(responseFile string) (*store2.Service, *httptest.Server, error) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slurp, err := ioutil.ReadFile(path.Join("testdata", responseFile))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		res, err := http.ReadResponse(bufio.NewReader(strings.NewReader(string(slurp))), r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer res.Body.Close()
		bs, err := ioutil.ReadAll(res.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(res.StatusCode)
		fmt.Fprint(w, string(bs))
	}))

	service, err := store2.New(http.DefaultClient)
	if err != nil {
		return service, nil, err
	}
	service.BaseURL = ts.URL // "http://store2.go/api/v2"
	service.User = os.Getenv("STORE2_USER")
	service.Password = os.Getenv("STORE2_PASSWORD")
	return service, ts, nil
}

func TestPing(t *testing.T) {
	service, ts, err := getService("ping.success")
	if err != nil {
		t.Fatal(err)
	}
	if service == nil {
		t.Fatal("expected service; got: nil")
	}
	defer ts.Close()

	err = service.Ping().Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
}

func TestPingUnauthorized(t *testing.T) {
	service, ts, err := getService("ping.unauthorized")
	if err != nil {
		t.Fatal(err)
	}
	if service == nil {
		t.Fatal("expected service; got: nil")
	}
	defer ts.Close()

	err = service.Ping().Do(context.Background())
	if err == nil {
		t.Fatalf("expected error; got: %v", err)
	}
}

func TestMe(t *testing.T) {
	service, ts, err := getService("me.success")
	if err != nil {
		t.Fatal(err)
	}
	if service == nil {
		t.Fatal("expected service; got: nil")
	}
	defer ts.Close()

	info, err := service.Me().Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if info == nil {
		t.Fatal("expected response; got: nil")
	}
	if info.Kind != "store#me" {
		t.Errorf("expected kind %q; got: %q", "store#me", info.Kind)
	}
	if strings.HasSuffix("/api/v2", info.SelfLink) {
		t.Errorf("expected selfLink suffix %s; got: %s", "/api/v2", info.SelfLink)
	}
}

func TestMeUnauthorized(t *testing.T) {
	service, ts, err := getService("me.unauthorized")
	service.User = ""
	service.Password = ""
	if err != nil {
		t.Fatal(err)
	}
	if service == nil {
		t.Fatal("expected service; got: nil")
	}
	defer ts.Close()

	info, err := service.Me().Do(context.Background())
	if err == nil {
		t.Fatal("expected error; got: nil")
	}
	if info != nil {
		t.Fatalf("expected no response; got: %v", info)
	}
	if err.Error() != "meplatoapi: Error 401: Unauthorized" {
		t.Errorf("expected error %q; got: %q", "meplatoapi: Error 401: Unauthorized", err.Error())
	}
}

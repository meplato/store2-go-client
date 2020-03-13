package jobs_test

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

	"github.com/meplato/store2-go-client/jobs"
)

func getService(responseFile string) (*jobs.Service, *httptest.Server, error) {
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

	service, err := jobs.New(http.DefaultClient)
	if err != nil {
		return service, nil, err
	}
	service.BaseURL = ts.URL
	service.User = os.Getenv("STORE2_USER")
	service.Password = os.Getenv("STORE2_PASSWORD")
	return service, ts, nil
}

func TestJobsSearch(t *testing.T) {
	service, ts, err := getService("jobs.search.success")
	if err != nil {
		t.Fatal(err)
	}
	if service == nil {
		t.Fatal("expected service; got: nil")
	}
	defer ts.Close()

	res, err := service.Search().Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if res == nil {
		t.Fatal("expected response; got: nil")
	}
}

func TestJobGet(t *testing.T) {
	service, ts, err := getService("jobs.get.success")
	if err != nil {
		t.Fatal(err)
	}
	if service == nil {
		t.Fatal("expected service; got: nil")
	}
	defer ts.Close()

	job, err := service.Get().ID("58097dc3-b279-49b5-a5da-23eb1c77d840").Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if job == nil {
		t.Fatal("expected response; got: nil")
	}
	if job.Kind != "store#job" {
		t.Errorf("expected %q; got: %q", "store#job", job.Kind)
	}
	if job.ID != "58097dc3-b279-49b5-a5da-23eb1c77d840" {
		t.Errorf("expected %q; got: %q", "58097dc3-b279-49b5-a5da-23eb1c77d840", job.ID)
	}
}

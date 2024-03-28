package availabilities_test

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/meplato/store2-go-client/v2/availabilities"
)

func getService(responseFile string) (*availabilities.Service, *httptest.Server, error) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slurp, err := os.ReadFile(path.Join("testdata", responseFile))
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
		bs, err := io.ReadAll(res.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(res.StatusCode)
		fmt.Fprint(w, string(bs))
	}))

	service, err := availabilities.New(http.DefaultClient)
	if err != nil {
		return service, nil, err
	}
	service.BaseURL = ts.URL
	service.User = os.Getenv("STORE2_USER")
	service.Password = os.Getenv("STORE2_PASSWORD")
	return service, ts, nil
}

func TestAvailabilitiesGet(t *testing.T) {
	service, ts, err := getService("availabilities.get.success")
	if err != nil {
		t.Fatal(err)
	}
	if service == nil {
		t.Fatal("expected service; got: nil")
	}
	defer ts.Close()

	res, err := service.Get().Spn("1234").Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if res == nil {
		t.Fatal("expected response; got: nil")
	}

	if len(res.Items) != 3 {
		t.Fatalf("expected items to be of length 3; got: %d", len(res.Items))
	}

	if res.Kind != "store#availabilities/getResponse" {
		t.Fatalf("expected kind %q; got: %v", "store#availabilities/getResponse", res.Kind)
	}
	firstItem := res.Items[0]

	if firstItem == nil {
		t.Fatal("expected first entry; got: nil")
	} else {
		if firstItem.Spn != "1234" {
			t.Fatalf("expected availability Spn %v; got: %v", "1234", firstItem.Spn)
		}
	}
}
func TestAvailabilitiesGetNotFound(t *testing.T) {
	service, ts, err := getService("availabilities.get.not_found")
	if err != nil {
		t.Fatal(err)
	}
	if service == nil {
		t.Fatal("expected service; got: nil")
	}
	defer ts.Close()

	res, err := service.Get().Spn("1235").Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if res == nil {
		t.Fatal("expected response; got: nil")
	}

	if len(res.Items) != 0 {
		t.Fatalf("expected items to be of length 0; got: %d", len(res.Items))
	}

	if res.Kind != "store#availabilities/getResponse" {
		t.Fatalf("expected kind %q; got: %v", "store#availabilities/getResponse", res.Kind)
	}
}
func TestAvailabilitiesUpsert(t *testing.T) {
	service, ts, err := getService("availabilities.upsert.success")
	if err != nil {
		t.Fatal(err)
	}
	if service == nil {
		t.Fatal("expected service; got: nil")
	}
	defer ts.Close()

	var Quantity = 0.0

	res, err := service.Upsert().Spn("1234").Availability(&availabilities.UpsertRequest{
		Message:  "not in stock",
		Quantity: &Quantity,
		Region:   "AQ",
		Updated:  "Q1/2024",
		ZipCode:  "1234",
	}).Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if res == nil {
		t.Fatal("expected response; got: nil")
	}

	if res.Kind != "store#availabilities/upsertResponse" {
		t.Fatalf("expected kind %q; got: %v", "store#availabilities/upsertResponse", res.Kind)
	}
}

func TestAvailabilitiesDelete(t *testing.T) {
	service, ts, err := getService("availabilities.delete.success")
	if err != nil {
		t.Fatal(err)
	}
	if service == nil {
		t.Fatal("expected service; got: nil")
	}
	defer ts.Close()

	res, err := service.Delete().Spn("1234").Region("DE").ZipCode("12345").Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if res == nil {
		t.Fatal("expected response; got: nil")
	}

	if res.Kind != "store#availabilities/deleteResponse" {
		t.Fatalf("expected kind %q; got: %v", "store#availabilities/deleteResponse", res.Kind)
	}
}

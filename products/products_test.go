package products_test

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/meplato/store2-go-client/products"
)

func getService(responseFile string) (*products.Service, *httptest.Server, error) {
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
		log.Print(string(bs))
		fmt.Fprint(w, string(bs))
	}))

	service, err := products.New(http.DefaultClient)
	if err != nil {
		return service, nil, err
	}
	service.BaseURL = ts.URL // "http://store2.go/api/v2"
	service.User = os.Getenv("STORE2_USER")
	service.Password = os.Getenv("STORE2_PASSWORD")
	return service, ts, nil
}

func TestProductSearch(t *testing.T) {
	service, ts, err := getService("products.search.success")
	if err != nil {
		t.Fatal(err)
	}
	if service == nil {
		t.Fatal("expected service; got: nil")
	}
	defer ts.Close()

	res, err := service.Search().PIN("PIN").Area("work").Skip(0).Take(30).Do()
	if err != nil {
		t.Fatal(err)
	}
	if res == nil {
		t.Fatal("expected response; got: nil")
	}
}

func TestProductGet(t *testing.T) {
	service, ts, err := getService("products.get.success")
	if err != nil {
		t.Fatal(err)
	}
	if service == nil {
		t.Fatal("expected service; got: nil")
	}
	defer ts.Close()

	res, err := service.Get().PIN("AD8CCDD5F9").Area("work").ID("50763599@12").Do()
	if err != nil {
		t.Fatal(err)
	}
	if res == nil {
		t.Fatal("expected response; got: nil")
	}
}

func TestProductCreate(t *testing.T) {
	service, ts, err := getService("products.create.success")
	if err != nil {
		t.Fatal(err)
	}
	if service == nil {
		t.Fatal("expected service; got: nil")
	}
	defer ts.Close()

	create := &products.CreateProduct{
		Spn:       "1000",
		Name:      "Produkt 1000",
		Price:     4.99,
		OrderUnit: "PCE",
	}

	cres, err := service.Create().PIN("AD8CCDD5F9").Area("work").Product(create).Do()
	if err != nil {
		t.Fatal(err)
	}
	if cres == nil {
		t.Fatal("expected response; got: nil")
	}
	if cres.ID == "" {
		t.Fatalf("expected ID; got: %v", cres.ID)
	}
	if cres.Link == "" {
		t.Fatalf("expected link to product; got: %v", cres.Link)
	}
}

func TestProductDelete(t *testing.T) {
	service, ts, err := getService("products.create.success")
	if err != nil {
		t.Fatal(err)
	}
	if service == nil {
		t.Fatal("expected service; got: nil")
	}
	defer ts.Close()

	err = service.Delete().PIN("AD8CCDD5F9").Area("work").ID("1000@14").Do()
	if err != nil {
		t.Fatal(err)
	}
}

func TestProductUpdate(t *testing.T) {
	service, ts, err := getService("products.create.success")
	if err != nil {
		t.Fatal(err)
	}
	if service == nil {
		t.Fatal("expected service; got: nil")
	}
	defer ts.Close()

	newName := "Produkt 1000 (geändert)"
	newPrice := 3.99
	update := &products.UpdateProduct{
		Name:  &newName,
		Price: &newPrice,
	}

	ures, err := service.Update().PIN("AD8CCDD5F9").Area("work").ID("MBA11@12").Product(update).Do()
	if err != nil {
		t.Fatal(err)
	}
	if ures == nil {
		t.Fatal("expected response; got: nil")
	}
	if ures.ID == "" {
		t.Fatalf("expected ID; got: %v", ures.ID)
	}
	if ures.Link == "" {
		t.Fatalf("expected link to product; got: %v", ures.Link)
	}
}

func TestProductReplace(t *testing.T) {
	service, ts, err := getService("products.replace.success")
	if err != nil {
		t.Fatal(err)
	}
	if service == nil {
		t.Fatal("expected service; got: nil")
	}
	defer ts.Close()

	replace := &products.ReplaceProduct{
		Name:      "Produkt 1000 (NEU!)",
		Price:     2.50,
		OrderUnit: "PK",
	}

	rres, err := service.Replace().PIN("AD8CCDD5F9").Area("work").ID("MBA11@12").Product(replace).Do()
	if err != nil {
		t.Fatal(err)
	}
	if rres == nil {
		t.Fatal("expected response; got: nil")
	}
	if rres.ID == "" {
		t.Fatalf("expected ID; got: %v", rres.ID)
	}
	if rres.Link == "" {
		t.Fatalf("expected link to product; got: %v", rres.Link)
	}
}

func TestProductCreateParameterMissing(t *testing.T) {
	service, ts, err := getService("products.create.parameter_missing")
	if err != nil {
		t.Fatal(err)
	}
	if service == nil {
		t.Fatal("expected service; got: nil")
	}
	defer ts.Close()

	create := &products.CreateProduct{
		Spn:       "", // we don't provide a SPN
		Name:      "Produkt 1000",
		Price:     4.99,
		OrderUnit: "PCE",
	}

	cres, err := service.Create().PIN("AD8CCDD5F9").Area("work").Product(create).Do()
	if err == nil {
		t.Fatal(err)
	}
	// Error is: meplatoapi: Error 400: SPN must not be blank
	if cres != nil {
		t.Fatalf("expected no create response; got: %v", cres)
	}
}

func TestProductScroll(t *testing.T) {
	service, ts, err := getService("products.scroll.success.1")
	if err != nil {
		t.Fatal(err)
	}
	if service == nil {
		t.Fatal("expected service; got: nil")
	}
	defer ts.Close()

	// Get first result set
	res, err := service.Scroll().PIN("AD8CCDD5F9").Area("work").Do()
	if err != nil {
		t.Fatal(err)
	}
	if res == nil {
		t.Fatal("expected response; got: nil")
	}
	if res.PageToken == "" {
		t.Fatalf("expected page token; got: %v", res.PageToken)
	}
	/*
		pageToken := res.PageToken
		for {
			res, err := service.Scroll().PIN(pin).Area("work").PageToken(pageToken).Do()
			if err != nil {
				t.Fatal(err)
			}
			if res == nil {
				t.Fatal("expected response; got: nil")
			}
			if res.PageToken == "" {
				break
			}
			if len(res.Items) == 0 {
				t.Fatalf("expected some results; got: %v", res.Items)
			}
			if res.Kind != "store#products" {
				t.Fatalf("expected store#products; got: %s", res.Kind)
			}

			pageToken = res.PageToken
		}
	*/
}

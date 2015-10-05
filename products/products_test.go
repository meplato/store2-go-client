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
	service.BaseURL = ts.URL
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

	res, err := service.Get().PIN("AD8CCDD5F9").Area("work").Spn("50763599").Do()
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
	if cres.Kind != "store#productsCreateResponse" {
		t.Fatalf("expected kind %q; got: %v", "store#productsCreateResponse", cres.Kind)
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

	err = service.Delete().PIN("AD8CCDD5F9").Area("work").Spn("1000").Do()
	if err != nil {
		t.Fatal(err)
	}
}

func TestProductUpdate(t *testing.T) {
	service, ts, err := getService("products.update.success")
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

	ures, err := service.Update().PIN("AD8CCDD5F9").Area("work").Spn("MBA11").Product(update).Do()
	if err != nil {
		t.Fatal(err)
	}
	if ures == nil {
		t.Fatal("expected response; got: nil")
	}
	if ures.Kind != "store#productsUpdateResponse" {
		t.Fatalf("expected kind %q; got: %v", "store#productsUpdateResponse", ures.Kind)
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

	rres, err := service.Replace().PIN("AD8CCDD5F9").Area("work").Spn("MBA11").Product(replace).Do()
	if err != nil {
		t.Fatal(err)
	}
	if rres == nil {
		t.Fatal("expected response; got: nil")
	}
	if rres.Kind != "store#productsReplaceResponse" {
		t.Fatalf("expected kind %q; got: %v", "store#productsReplaceResponse", rres.Kind)
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

func TestProductUpsert(t *testing.T) {
	service, ts, err := getService("products.upsert.success")
	if err != nil {
		t.Fatal(err)
	}
	if service == nil {
		t.Fatal("expected service; got: nil")
	}
	defer ts.Close()

	up := &products.UpsertProduct{
		Spn:       "1000",
		Name:      "Produkt 1000",
		Price:     4.99,
		OrderUnit: "PCE",
	}

	res, err := service.Upsert().PIN("AD8CCDD5F9").Area("work").Product(up).Do()
	if err != nil {
		t.Fatal(err)
	}
	if res == nil {
		t.Fatal("expected response; got: nil")
	}
	if res.Kind != "store#productsUpsertResponse" {
		t.Fatalf("expected kind %q; got: %v", "store#productsUpsertResponse", res.Kind)
	}
	if res.Link == "" {
		t.Fatalf("expected link to product; got: %v", res.Link)
	}
}

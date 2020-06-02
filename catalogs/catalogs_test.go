package catalogs_test

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

	"github.com/meplato/store2-go-client/v2/catalogs"
)

func getService(responseFile string) (*catalogs.Service, *httptest.Server, error) {
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

	service, err := catalogs.New(http.DefaultClient)
	if err != nil {
		return service, nil, err
	}
	service.BaseURL = ts.URL
	service.User = os.Getenv("STORE2_USER")
	service.Password = os.Getenv("STORE2_PASSWORD")
	return service, ts, nil
}

func TestCatalogSearch(t *testing.T) {
	service, ts, err := getService("catalogs.search.success")
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

func TestCatalogCreate(t *testing.T) {
	service, ts, err := getService("catalogs.create.success")
	if err != nil {
		t.Fatal(err)
	}
	if service == nil {
		t.Fatal("expected service; got: nil")
	}
	defer ts.Close()

	create := &catalogs.CreateCatalog{
		MerchantID: 1,
		Name:       "test2",
		// Description: "",
		// ProjectID:    0,
		ProjectMpcc:  "meplato",
		ValidFrom:    nil,
		ValidUntil:   nil,
		Country:      "DE",
		Currency:     "EUR",
		Language:     "de",
		Target:       "mall",
		Type:         "CC",
		SageNumber:   "",
		SageContract: "",
	}

	cres, err := service.Create().Catalog(create).Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if cres == nil {
		t.Fatal("expected response; got: nil")
	}
	if cres.Kind != "store#catalog" {
		t.Fatalf("expected kind %q; got: %v", "store#catalog", cres.Kind)
	}
	if cres.ID != 81 {
		t.Fatalf("expected catalog id of %d; got: %d", 81, cres.ID)
	}
	if cres.PIN != "48F31F33AD" {
		t.Fatalf("expected catalog pin of %q; got: %q", "48F31F33AD", cres.PIN)
	}
	if cres.Type != "CC" {
		t.Fatalf("expected catalog type of %q; got: %q", "CC", cres.Type)
	}
}

func TestCatalogGet(t *testing.T) {
	service, ts, err := getService("catalogs.get.success")
	if err != nil {
		t.Fatal(err)
	}
	if service == nil {
		t.Fatal("expected service; got: nil")
	}
	defer ts.Close()

	c, err := service.Get().PIN("5094310527").Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if c == nil {
		t.Fatal("expected response; got: nil")
	}
	if c.Kind != "store#catalog" {
		t.Errorf("expected %q; got: %q", "store#catalog", c.Kind)
	}
	if c.ID != 14 {
		t.Errorf("expected %d; got: %d", 14, c.ID)
	}
}

func TestCatalogPublish(t *testing.T) {
	service, ts, err := getService("catalogs.publish.success")
	if err != nil {
		t.Fatal(err)
	}
	if service == nil {
		t.Fatal("expected service; got: nil")
	}
	defer ts.Close()

	// Publish
	pub, err := service.Publish().PIN("5094310527").Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if pub == nil {
		t.Fatal("expected response; got: nil")
	}
	if pub.Kind != "store#catalogPublish" {
		t.Errorf("expected %q; got: %q", "store#catalogPublish", pub.Kind)
	}
	if pub.SelfLink == "" {
		t.Errorf("expected self link; got: %q", pub.SelfLink)
	}
	if pub.StatusLink == "" {
		t.Errorf("expected status link; got: %q", pub.StatusLink)
	}
}

/*
	// Watch status for max. 10 seconds
	var i int
	const N = 10
	for {
		time.Sleep(5 * time.Second)

		status, err := service.PublishStatus().PIN("AD8CCDD5F9").Do()
		if err != nil {
			t.Fatal(err)
		}
		if status == nil {
			t.Fatal("expected response; got: nil")
		}
		if status.Kind != "store#catalogPublishStatus" {
			t.Errorf("expected %q; got: %q", "store#catalogPublishStatus", status.Kind)
		}
		if status.Done {
			break
		}
		i += 1
		if i > N {
			t.Fatal("expected publish to complete after a while")
		}
	}
*/

func TestCatalogPurge(t *testing.T) {
	service, ts, err := getService("catalogs.purge.success")
	if err != nil {
		t.Fatal(err)
	}
	if service == nil {
		t.Fatal("expected service; got: nil")
	}
	defer ts.Close()

	c, err := service.Purge().PIN("5094310527").Area("work").Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if c == nil {
		t.Fatal("expected response; got: nil")
	}
	if c.Kind != "store#catalogPurge" {
		t.Errorf("expected %q; got: %q", "store#catalogPurge", c.Kind)
	}
}

// Copyright (c) 2013-present Meplato GmbH.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except
// in compliance with the License. You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software distributed under the License
// is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
// or implied. See the License for the specific language governing permissions and limitations under
// the License.

// Package catalogs implements the Meplato Store API.
//
// See https://developer.meplato.com/store2/.
package catalogs

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/meplato/store2-go-client/v2/internal/meplatoapi"
)

// Always reference these packages, just in case.
var (
	_ = bytes.NewBuffer
	_ = http.Get
	_ = fmt.Print
	_ = bytes.NewBuffer
	_ = json.NewDecoder
	_ = errors.New
	_ = fmt.Print
	_ = io.Copy
	_ = http.Get
	_ = url.Parse
	_ = strconv.Itoa
	_ = strings.HasPrefix
	_ = time.Parse
	_ = meplatoapi.CheckResponse
)

const (
	title   = "Meplato Store API"
	version = "2.1.8"
	baseURL = "https://store.meplato.com/api/v2"
)

type Service struct {
	client   *http.Client
	BaseURL  string
	User     string
	Password string
}

func New(client *http.Client) (*Service, error) {
	if client == nil {
		return nil, errors.New("client is nil")
	}
	return &Service{client: client, BaseURL: baseURL}, nil
}

func (s *Service) Create() *CreateService {
	return NewCreateService(s)
}

func (s *Service) Get() *GetService {
	return NewGetService(s)
}

func (s *Service) Publish() *PublishService {
	return NewPublishService(s)
}

func (s *Service) PublishStatus() *PublishStatusService {
	return NewPublishStatusService(s)
}

func (s *Service) Purge() *PurgeService {
	return NewPurgeService(s)
}

func (s *Service) Search() *SearchService {
	return NewSearchService(s)
}

// Catalog is a container for products, to be used in a certain project.
type Catalog struct {
	// Country is the ISO-3166 alpha-2 code for the country that the catalog
	// is destined for (e.g. DE or US).
	Country string `json:"country,omitempty"`
	// Created is the creation date and time of the catalog.
	Created *time.Time `json:"created,omitempty"`
	// Currency is the ISO-4217 currency code that is used for all products in
	// the catalog (e.g. EUR or USD).
	Currency string `json:"currency,omitempty"`
	// CustFields is an array of generic name/value pairs for
	// customer-specific attributes.
	CustFields []*CustField `json:"custFields,omitempty"`
	// Description of the catalog.
	Description string `json:"description,omitempty"`
	// DownloadChecksum represents the checksum of the catalog last
	// downloaded.
	DownloadChecksum string `json:"downloadChecksum,omitempty"`
	// DownloadInterval represents the interval to use for checking new
	// versions of a catalog at the DownloadURL.
	DownloadInterval string `json:"downloadInterval,omitempty"`
	// DownloadURL represents a URL which is periodically downloaded and
	// imported as a new catalog.
	DownloadURL string `json:"downloadUrl,omitempty"`
	// ErpNumberBuyer: ERPNumberBuyer is the number of the merchant of this
	// catalog in the SAP/ERP system of the buyer.
	ErpNumberBuyer string `json:"erpNumberBuyer,omitempty"`
	// Expired indicates whether the catalog is expired as of now.
	Expired bool `json:"expired,omitempty"`
	// HubURL represents the Meplato Hub URL for this catalog, e.g.
	// https://hub.meplato.de/forward/12345/shop
	HubURL string `json:"hubUrl,omitempty"`
	// ID is a unique (internal) identifier of the catalog.
	ID int64 `json:"id,omitempty"`
	// KeepOriginalBlobs indicates whether the URLs in a blob will be passed
	// through and not cached by Store.
	KeepOriginalBlobs bool `json:"keepOriginalBlobs,omitempty"`
	// Kind is store#catalog for a catalog entity.
	Kind string `json:"kind,omitempty"`
	// KpiSummary: KPISummary returns the outcome of analyzing the contents
	// for key performance indicators.
	KpiSummary *KPISummary `json:"kpiSummary,omitempty"`
	// Language is the IETF language tag of the language of all products in
	// the catalog (e.g. de or pt-BR).
	Language string `json:"language,omitempty"`
	// LastImported is the date and time the catalog was last imported.
	LastImported *time.Time `json:"lastImported,omitempty"`
	// LastPublished is the date and time the catalog was last published.
	LastPublished *time.Time `json:"lastPublished,omitempty"`
	// LockedForDownload indicates whether a catalog is locked and cannot be
	// downloaded.
	LockedForDownload bool `json:"lockedForDownload,omitempty"`
	// MerchantID: ID of the merchant.
	MerchantID int64 `json:"merchantId,omitempty"`
	// MerchantMpcc: MPCC of the merchant.
	MerchantMpcc string `json:"merchantMpcc,omitempty"`
	// MerchantMpsc: MPSC of the merchant.
	MerchantMpsc string `json:"merchantMpsc,omitempty"`
	// MerchantName: Name of the merchant.
	MerchantName string `json:"merchantName,omitempty"`
	// Name of the catalog.
	Name string `json:"name,omitempty"`
	// NumProductsLive: Number of products currently in the live area (only
	// returned when getting the details of a catalog).
	NumProductsLive *int64 `json:"numProductsLive,omitempty"`
	// NumProductsWork: Number of products currently in the work area (only
	// returned when getting the details of a catalog).
	NumProductsWork *int64 `json:"numProductsWork,omitempty"`
	// OciURL represents the OCI punchout URL that the supplier specified for
	// this catalog, e.g. https://my-shop.com/oci?param1=a
	OciURL string `json:"ociUrl,omitempty"`
	// PIN of the catalog.
	PIN string `json:"pin,omitempty"`
	// Project references the project that this catalog belongs to.
	Project *Project `json:"project,omitempty"`
	// ProjectID: ID of the project.
	ProjectID int64 `json:"projectId,omitempty"`
	// ProjectMpbc: MPBC of the project.
	ProjectMpbc string `json:"projectMpbc,omitempty"`
	// ProjectMpcc: MPCC of the project.
	ProjectMpcc string `json:"projectMpcc,omitempty"`
	// ProjectName: Name of the project.
	ProjectName string `json:"projectName,omitempty"`
	// PublishedVersion is the version number of the published catalog. It is
	// incremented when the publish task publishes the catalog.
	PublishedVersion *int64 `json:"publishedVersion,omitempty"`
	// SageContract represents the internal identifier at Meplato for the
	// contract of this catalog.
	SageContract string `json:"sageContract,omitempty"`
	// SageNumber represents the internal identifier at Meplato for the
	// merchant of this catalog.
	SageNumber string `json:"sageNumber,omitempty"`
	// SelfLink: URL to this page.
	SelfLink string `json:"selfLink,omitempty"`
	// Slug of the catalog.
	Slug string `json:"slug,omitempty"`
	// State describes the current state of the catalog, e.g. idle.
	State string `json:"state,omitempty"`
	// SupportsOciBackgroundsearch indicates whether a catalog supports the
	// OCI BACKGROUNDSEARCH transaction.
	SupportsOciBackgroundsearch bool `json:"supportsOciBackgroundsearch,omitempty"`
	// SupportsOciDetail indicates whether a catalog supports the OCI DETAIL
	// transaction.
	SupportsOciDetail bool `json:"supportsOciDetail,omitempty"`
	// SupportsOciDetailadd indicates whether a catalog supports the OCI
	// DETAILADD transaction.
	SupportsOciDetailadd bool `json:"supportsOciDetailadd,omitempty"`
	// SupportsOciDownloadjson indicates whether a catalog supports the OCI
	// DOWNLOADJSON transaction.
	SupportsOciDownloadjson bool `json:"supportsOciDownloadjson,omitempty"`
	// SupportsOciQuantitycheck indicates whether a catalog supports the OCI
	// QUANTITYCHECK transaction.
	SupportsOciQuantitycheck bool `json:"supportsOciQuantitycheck,omitempty"`
	// SupportsOciSourcing indicates whether a catalog supports the OCI
	// SOURCING transaction.
	SupportsOciSourcing bool `json:"supportsOciSourcing,omitempty"`
	// SupportsOciValidate indicates whether a catalog supports the OCI
	// VALIDATE transaction.
	SupportsOciValidate bool `json:"supportsOciValidate,omitempty"`
	// Target represents the target system which can be either an empty
	// string, "catscout" or "mall".
	Target string `json:"target,omitempty"`
	// Type of catalog, e.g. corporate or basic.
	Type string `json:"type,omitempty"`
	// Updated is the last modification date and time of the catalog.
	Updated *time.Time `json:"updated,omitempty"`
	// ValidFrom is the date the catalog becomes effective (YYYY-MM-DD).
	ValidFrom *string `json:"validFrom,omitempty"`
	// ValidUntil is the date the catalog expires (YYYY-MM-DD).
	ValidUntil *string `json:"validUntil,omitempty"`
}

// CreateCatalog holds the properties of a new catalog.
type CreateCatalog struct {
	// Country is the ISO-3166 alpha-2 code for the country that the catalog
	// is destined for (e.g. DE or US).
	Country string `json:"country,omitempty"`
	// Currency is the ISO-4217 currency code that is used for all products in
	// the catalog (e.g. EUR or USD).
	Currency string `json:"currency,omitempty"`
	// Description of the catalog.
	Description string `json:"description,omitempty"`
	// Language is the IETF language tag of the language of all products in
	// the catalog (e.g. de or pt-BR).
	Language string `json:"language,omitempty"`
	// MerchantID: ID of the merchant.
	MerchantID int64 `json:"merchantId,omitempty"`
	// Name of the catalog.
	Name string `json:"name,omitempty"`
	// ProjectID: ID of the project.
	ProjectID int64 `json:"projectId,omitempty"`
	// ProjectMpcc: MPCC of the project.
	ProjectMpcc string `json:"projectMpcc,omitempty"`
	// SageContract represents the internal identifier at Meplato for the
	// contract of this catalog.
	SageContract string `json:"sageContract,omitempty"`
	// SageNumber represents the internal identifier at Meplato for the
	// merchant of this catalog.
	SageNumber string `json:"sageNumber,omitempty"`
	// Target represents the target system which can be either an empty
	// string, "catscout" or "mall".
	Target string `json:"target,omitempty"`
	// ValidFrom is the date the catalog becomes effective (YYYY-MM-DD).
	ValidFrom *string `json:"validFrom,omitempty"`
	// ValidUntil is the date the catalog expires (YYYY-MM-DD).
	ValidUntil *string `json:"validUntil,omitempty"`
}

// CustField describes a generic name/value pair. Its purpose is to
// provide a mechanism for customer-specific fields.
type CustField struct {
	// Name is the name of the customer-specific field, e.g. TaxRate.
	Name string `json:"name,omitempty"`
	// Value is the value of the customer-specific field, e.g. 19%%.
	Value string `json:"value,omitempty"`
}

// KPISummary represents the outcome of analyzing the contents for key
// performance indicators.
type KPISummary struct {
	// Coefficients represents the weight that is used to calculate the
	// weighted coefficients for a criteria. It relies on the medal stored in
	// DegreesOfFulfillment.
	Coefficients map[string]float64 `json:"coefficients,omitempty"`
	// CreatedAt is the date/time when the KPI summary has been created.
	CreatedAt time.Time `json:"createdAt,omitempty"`
	// DegreesOfFulfillment represents the medal for all KPI criteria: 3 for
	// gold, 2 for silver, 1 for bronze, 0 for no medal.
	DegreesOfFulfillment map[string]int `json:"degreesOfFulfillment,omitempty"`
	// FinalResult returns a value between 0.0 and 1.0 that describes the
	// weighted sum of all content-related test criteria.
	FinalResult float64 `json:"finalResult,omitempty"`
	// OverallResult returns 3 for a gold medal, 2 for a silver medal, 1 for a
	// bronze medal, and 0 for no medal.
	OverallResult int `json:"overallResult,omitempty"`
	// TestResults represents the unweighted outcome for a specific KPI
	// criteria, i.e. the percentage of products that fulfill the criteria.
	TestResults map[string]float64 `json:"testResults,omitempty"`
	// WeightedCoefficients is a value between 0.0 and 1.0 that represents the
	// weighted outcome of a KPI criteria, as calculated by the coefficient
	// and the test result.
	WeightedCoefficients map[string]float64 `json:"weightedCoefficients,omitempty"`
}

// Project describes customer-specific settings, typically encompassing a
// set of catalogs.
type Project struct {
	// Country specifies the country code where catalogs for this project are
	// located.
	Country string `json:"country,omitempty"`
	// Created is the creation date and time of the project.
	Created *time.Time `json:"created,omitempty"`
	// ID is a unique (internal) identifier of the project.
	ID int64 `json:"id,omitempty"`
	// Kind is store#project for a project entity.
	Kind string `json:"kind,omitempty"`
	// Language specifies the language code of the catalogs of this project.
	Language string `json:"language,omitempty"`
	// Mpbc: MPBC is the Meplato Buyer Code that identifies a set of buy-side
	// companies that belong together.
	Mpbc string `json:"mpbc,omitempty"`
	// Mpcc: MPCC is the Meplato Company Code that uniquely identifies the
	// buy-side.
	Mpcc string `json:"mpcc,omitempty"`
	// Name is a short description of the project.
	Name string `json:"name,omitempty"`
	// SelfLink: URL to this page.
	SelfLink string `json:"selfLink,omitempty"`
	// Type describes the type of project which can be either corporate or
	// basic.
	Type string `json:"type,omitempty"`
	// Updated is the last modification date and time of the project.
	Updated *time.Time `json:"updated,omitempty"`
	// Visible indicates whether this project is visible to merchants.
	Visible bool `json:"visible,omitempty"`
}

// PublishResponse is the response of the request to publish a catalog.
type PublishResponse struct {
	// Kind is store#catalogPublish for this kind of response.
	Kind string `json:"kind,omitempty"`
	// SelfLink returns the URL to this page.
	SelfLink string `json:"selfLink,omitempty"`
	// StatusLink returns the URL that returns the current status of the
	// request.
	StatusLink string `json:"statusLink,omitempty"`
}

// PublishStatusResponse returns current information about the status of a
// publish request.
type PublishStatusResponse struct {
	// Busy indicates whether the catalog is still busy.
	Busy bool `json:"busy,omitempty"`
	// Canceled indicates whether the publishing process has been canceled.
	Canceled bool `json:"canceled,omitempty"`
	// CurrentStep is an indicator of the current step in the total list of
	// steps. Use in combination with TotalSteps to retrieve the progress in
	// percent.
	CurrentStep int64 `json:"currentStep,omitempty"`
	// Done indicates whether publishing is finished.
	Done bool `json:"done,omitempty"`
	// Kind is store#catalogPublishStatus for this kind of response.
	Kind string `json:"kind,omitempty"`
	// Percent indicates the progress of the publish request.
	Percent int `json:"percent,omitempty"`
	// SelfLink returns the URL to this page.
	SelfLink string `json:"selfLink,omitempty"`
	// Status describes the general status of the publish request.
	Status string `json:"status,omitempty"`
	// TotalSteps is an indicator of the total number steps required to
	// complete the publish request. Use in combination with CurrentStep.
	TotalSteps int64 `json:"totalSteps,omitempty"`
}

// PurgeResponse is the response of the request to purge an area of a
// catalog.
type PurgeResponse struct {
	// Kind is store#catalogPurge for this kind of response.
	Kind string `json:"kind,omitempty"`
}

// SearchResponse is a partial listing of catalogs.
type SearchResponse struct {
	// Items is the slice of catalogs of this result.
	Items []*Catalog `json:"items,omitempty"`
	// Kind is store#catalogs for this kind of response.
	Kind string `json:"kind,omitempty"`
	// NextLink returns the URL to the next slice of catalogs (if any).
	NextLink string `json:"nextLink,omitempty"`
	// PreviousLink returns the URL of the previous slice of catalogs (if
	// any).
	PreviousLink string `json:"previousLink,omitempty"`
	// SelfLink returns the URL to this page.
	SelfLink string `json:"selfLink,omitempty"`
	// TotalItems describes the total number of catalogs found.
	TotalItems int64 `json:"totalItems,omitempty"`
}

// Create a new catalog (admin only).
type CreateService struct {
	s       *Service
	opt_    map[string]interface{}
	hdr_    map[string]interface{}
	catalog *CreateCatalog
}

// NewCreateService creates a new instance of CreateService.
func NewCreateService(s *Service) *CreateService {
	rs := &CreateService{s: s, opt_: make(map[string]interface{}), hdr_: make(map[string]interface{})}
	return rs
}

// Catalog properties of the new catalog.
func (s *CreateService) Catalog(catalog *CreateCatalog) *CreateService {
	s.catalog = catalog
	return s
}

// Do executes the operation.
func (s *CreateService) Do(ctx context.Context) (*Catalog, error) {
	var body io.Reader
	body, err := meplatoapi.ReadJSON(s.catalog)
	if err != nil {
		return nil, err
	}
	params := make(map[string]interface{})
	path := "/catalogs"
	if len(params) > 0 {
		query := url.Values{}
		for k, v := range params {
			query.Add(k, fmt.Sprintf("%v", v))
		}
		path += "?" + query.Encode()
	}
	req, err := http.NewRequest("POST", s.s.BaseURL+path, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Charset", "utf-8")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", meplatoapi.UserAgent)
	if s.s.User != "" || s.s.Password != "" {
		req.Header.Set("Authorization", meplatoapi.HTTPBasicAuthorizationHeader(s.s.User, s.s.Password))
	}
	res, err := s.s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer meplatoapi.CloseBody(res)
	if err := meplatoapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := new(Catalog)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
		return nil, err
	}
	return ret, nil
}

// Get a single catalog.
type GetService struct {
	s    *Service
	opt_ map[string]interface{}
	hdr_ map[string]interface{}
	pin  string
}

// NewGetService creates a new instance of GetService.
func NewGetService(s *Service) *GetService {
	rs := &GetService{s: s, opt_: make(map[string]interface{}), hdr_: make(map[string]interface{})}
	return rs
}

// PIN of the catalog.
func (s *GetService) PIN(pin string) *GetService {
	s.pin = pin
	return s
}

// Do executes the operation.
func (s *GetService) Do(ctx context.Context) (*Catalog, error) {
	var body io.Reader
	params := make(map[string]interface{})
	params["pin"] = s.pin
	path, err := meplatoapi.Expand("/catalogs/{pin}", params)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("GET", s.s.BaseURL+path, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Charset", "utf-8")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", meplatoapi.UserAgent)
	if s.s.User != "" || s.s.Password != "" {
		req.Header.Set("Authorization", meplatoapi.HTTPBasicAuthorizationHeader(s.s.User, s.s.Password))
	}
	res, err := s.s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer meplatoapi.CloseBody(res)
	if err := meplatoapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := new(Catalog)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
		return nil, err
	}
	return ret, nil
}

// Publishes a catalog.
type PublishService struct {
	s    *Service
	opt_ map[string]interface{}
	hdr_ map[string]interface{}
	pin  string
}

// NewPublishService creates a new instance of PublishService.
func NewPublishService(s *Service) *PublishService {
	rs := &PublishService{s: s, opt_: make(map[string]interface{}), hdr_: make(map[string]interface{})}
	return rs
}

// PIN of the catalog to publish.
func (s *PublishService) PIN(pin string) *PublishService {
	s.pin = pin
	return s
}

// Do executes the operation.
func (s *PublishService) Do(ctx context.Context) (*PublishResponse, error) {
	var body io.Reader
	params := make(map[string]interface{})
	params["pin"] = s.pin
	path, err := meplatoapi.Expand("/catalogs/{pin}/publish", params)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", s.s.BaseURL+path, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Charset", "utf-8")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", meplatoapi.UserAgent)
	if s.s.User != "" || s.s.Password != "" {
		req.Header.Set("Authorization", meplatoapi.HTTPBasicAuthorizationHeader(s.s.User, s.s.Password))
	}
	res, err := s.s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer meplatoapi.CloseBody(res)
	if err := meplatoapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := new(PublishResponse)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
		return nil, err
	}
	return ret, nil
}

// Status of a publish process.
type PublishStatusService struct {
	s    *Service
	opt_ map[string]interface{}
	hdr_ map[string]interface{}
	pin  string
}

// NewPublishStatusService creates a new instance of PublishStatusService.
func NewPublishStatusService(s *Service) *PublishStatusService {
	rs := &PublishStatusService{s: s, opt_: make(map[string]interface{}), hdr_: make(map[string]interface{})}
	return rs
}

// PIN of the catalog to get the publish status from.
func (s *PublishStatusService) PIN(pin string) *PublishStatusService {
	s.pin = pin
	return s
}

// Do executes the operation.
func (s *PublishStatusService) Do(ctx context.Context) (*PublishStatusResponse, error) {
	var body io.Reader
	params := make(map[string]interface{})
	params["pin"] = s.pin
	path, err := meplatoapi.Expand("/catalogs/{pin}/publish/status", params)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("GET", s.s.BaseURL+path, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Charset", "utf-8")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", meplatoapi.UserAgent)
	if s.s.User != "" || s.s.Password != "" {
		req.Header.Set("Authorization", meplatoapi.HTTPBasicAuthorizationHeader(s.s.User, s.s.Password))
	}
	res, err := s.s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer meplatoapi.CloseBody(res)
	if err := meplatoapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := new(PublishStatusResponse)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
		return nil, err
	}
	return ret, nil
}

// Purge the work or live area of a catalog, i.e. remove all products in
// the given area, but do not delete the catalog itself.
type PurgeService struct {
	s    *Service
	opt_ map[string]interface{}
	hdr_ map[string]interface{}
	pin  string
	area string
}

// NewPurgeService creates a new instance of PurgeService.
func NewPurgeService(s *Service) *PurgeService {
	rs := &PurgeService{s: s, opt_: make(map[string]interface{}), hdr_: make(map[string]interface{})}
	return rs
}

// Area of the catalog to purge, i.e. work or live.
func (s *PurgeService) Area(area string) *PurgeService {
	s.area = area
	return s
}

// PIN of the catalog to purge.
func (s *PurgeService) PIN(pin string) *PurgeService {
	s.pin = pin
	return s
}

// Do executes the operation.
func (s *PurgeService) Do(ctx context.Context) (*PurgeResponse, error) {
	var body io.Reader
	params := make(map[string]interface{})
	params["area"] = s.area
	params["pin"] = s.pin
	path, err := meplatoapi.Expand("/catalogs/{pin}/{area}", params)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("DELETE", s.s.BaseURL+path, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Charset", "utf-8")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", meplatoapi.UserAgent)
	if s.s.User != "" || s.s.Password != "" {
		req.Header.Set("Authorization", meplatoapi.HTTPBasicAuthorizationHeader(s.s.User, s.s.Password))
	}
	res, err := s.s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer meplatoapi.CloseBody(res)
	if err := meplatoapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := new(PurgeResponse)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
		return nil, err
	}
	return ret, nil
}

// Search for catalogs.
type SearchService struct {
	s    *Service
	opt_ map[string]interface{}
	hdr_ map[string]interface{}
}

// NewSearchService creates a new instance of SearchService.
func NewSearchService(s *Service) *SearchService {
	rs := &SearchService{s: s, opt_: make(map[string]interface{}), hdr_: make(map[string]interface{})}
	return rs
}

// Q defines are full text query.
func (s *SearchService) Q(q string) *SearchService {
	s.opt_["q"] = q
	return s
}

// Skip specifies how many catalogs to skip (default 0).
func (s *SearchService) Skip(skip int64) *SearchService {
	s.opt_["skip"] = skip
	return s
}

// Sort order, e.g. name or id or -created (default: score).
func (s *SearchService) Sort(sort string) *SearchService {
	s.opt_["sort"] = sort
	return s
}

// Take defines how many catalogs to return (max 100, default 20).
func (s *SearchService) Take(take int64) *SearchService {
	s.opt_["take"] = take
	return s
}

// Do executes the operation.
func (s *SearchService) Do(ctx context.Context) (*SearchResponse, error) {
	var body io.Reader
	params := make(map[string]interface{})
	if v, ok := s.opt_["q"]; ok {
		params["q"] = v
	}
	if v, ok := s.opt_["skip"]; ok {
		params["skip"] = v
	}
	if v, ok := s.opt_["sort"]; ok {
		params["sort"] = v
	}
	if v, ok := s.opt_["take"]; ok {
		params["take"] = v
	}
	path, err := meplatoapi.Expand("/catalogs{?q,skip,take,sort}", params)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("GET", s.s.BaseURL+path, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Charset", "utf-8")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", meplatoapi.UserAgent)
	if s.s.User != "" || s.s.Password != "" {
		req.Header.Set("Authorization", meplatoapi.HTTPBasicAuthorizationHeader(s.s.User, s.s.Password))
	}
	res, err := s.s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer meplatoapi.CloseBody(res)
	if err := meplatoapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := new(SearchResponse)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
		return nil, err
	}
	return ret, nil
}

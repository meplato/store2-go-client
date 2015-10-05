// Copyright (c) 2015 Meplato GmbH, Switzerland.
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

// Package catalogs implements the Meplato Store 2 API.
//
// See https://developer.meplato.com/store2/.
package catalogs

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/meplato/store2-go-client/internal/meplatoapi"
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
	title   = "Meplato Store 2 API"
	version = "2.0.0.beta3"
	baseURL = "https://store2.meplato.com/api/v2"
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
	// Created is the creation date and time of the catalog.
	Created *time.Time `json:"created,omitempty"`
	// Currency is the ISO-4217 currency code that is used for all products in
	// the catalog.
	Currency string `json:"currency,omitempty"`
	// Description of the catalog.
	Description string `json:"description,omitempty"`
	// ErpNumberBuyer: ERPNumberBuyer is the number of the merchant of this
	// catalog in the SAP/ERP system of the buyer.
	ErpNumberBuyer string `json:"erpNumberBuyer,omitempty"`
	// ID is a unique (internal) identifier of the catalog.
	ID int64 `json:"id,omitempty"`
	// Kind is store#catalog for a catalog entity.
	Kind string `json:"kind,omitempty"`
	// Language is the IETF language tag of the language of all products in
	// the catalog.
	Language string `json:"language,omitempty"`
	// LastImported is the date and time the catalog was last imported.
	LastImported *time.Time `json:"lastImported,omitempty"`
	// LastPublished is the date and time the catalog was last published.
	LastPublished *time.Time `json:"lastPublished,omitempty"`
	// MerchantID: ID of the merchant.
	MerchantID int64 `json:"merchantId,omitempty"`
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
	// PIN of the catalog.
	PIN string `json:"pin,omitempty"`
	// ProjectID: ID of the project.
	ProjectID int64 `json:"projectId,omitempty"`
	// PublishedVersion is the version number of the published catalog. It is
	// incremented when the publish task publishes the catalog.
	PublishedVersion *int64 `json:"publishedVersion,omitempty"`
	// SelfLink: URL to this page.
	SelfLink string `json:"selfLink,omitempty"`
	// Slug of the catalog.
	Slug string `json:"slug,omitempty"`
	// State describes the current state of the catalog, e.g. idle.
	State string `json:"state,omitempty"`
	// Updated is the last modification date and time of the catalog.
	Updated *time.Time `json:"updated,omitempty"`
	// ValidFrom is the date the catalog becomes effective.
	ValidFrom *string `json:"validFrom,omitempty"`
	// ValidUntil is the date the catalog expires.
	ValidUntil *string `json:"validUntil,omitempty"`
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
func (s *GetService) Do() (*Catalog, error) {
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
func (s *PublishService) Do() (*PublishResponse, error) {
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
func (s *PublishStatusService) Do() (*PublishStatusResponse, error) {
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
func (s *PurgeService) Do() (*PurgeResponse, error) {
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

// Skip specifies how many catalogs to skip (default 0).
func (s *SearchService) Skip(skip int64) *SearchService {
	s.opt_["skip"] = skip
	return s
}

// Sort order, e.g. name or id or -created (default: name).
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
func (s *SearchService) Do() (*SearchResponse, error) {
	var body io.Reader
	params := make(map[string]interface{})
	if v, ok := s.opt_["skip"]; ok {
		params["skip"] = v
	}
	if v, ok := s.opt_["sort"]; ok {
		params["sort"] = v
	}
	if v, ok := s.opt_["take"]; ok {
		params["take"] = v
	}
	path, err := meplatoapi.Expand("/catalogs{?skip,take,sort}", params)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("GET", s.s.BaseURL+path, body)
	if err != nil {
		return nil, err
	}
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

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

// Package availabilities implements the Meplato Store API.
//
// See https://developer.meplato.com/store2/.
package availabilities

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
	version = "2.2.0"
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

func (s *Service) Delete() *DeleteService {
	return NewDeleteService(s)
}

func (s *Service) Get() *GetService {
	return NewGetService(s)
}

func (s *Service) Upsert() *UpsertService {
	return NewUpsertService(s)
}

// Availability information of a product in a location
type Availability struct {
	// Message: Contains the stock state description; i.e. in stock; out of
	// stock; limited availability; on display to order
	Message string `json:"message,omitempty"`
	// Mpcc: Unique internal identifier of the merchant
	Mpcc string `json:"mpcc,omitempty"`
	// Quantity: Reflects the amount of items available
	Quantity *float64 `json:"quantity,omitempty"`
	// Region: 2-letter ISO code of the country/region where the product is
	// stored
	Region string `json:"region,omitempty"`
	// Spn: Merchant's unique identifier of a product
	Spn string `json:"spn,omitempty"`
	// Updated: Update date given by the merchant i.e. Q4/2022, 2022/10/12
	Updated string `json:"updated,omitempty"`
	// ZipCode: Zip code where the product is stored
	ZipCode string `json:"zipCode,omitempty"`
}

// DeleteResponse is the outcome of a successful request to delete an
// availability.
type DeleteResponse struct {
	// Kind describes this entity, it will be
	// store#availability/deleteResponse.
	Kind string `json:"kind,omitempty"`
}

// GetResponse is the collection of availability information for an SPN.
type GetResponse struct {
	// Items: Collection of availability information associated with an SPN
	// for a merchant.
	Items []*Availability `json:"items,omitempty"`
	// Kind is store#availability/getResponse for this kind of response.
	Kind string `json:"kind,omitempty"`
}

// UpsertRequest holds the properties of the availability information to
// create or update.
type UpsertRequest struct {
	// Message: Contains the stock state description; i.e. in stock; out of
	// stock; limited availability; on display to order
	Message string `json:"message,omitempty"`
	// Mpcc: Unique internal identifier of the merchant (optional)
	Mpcc string `json:"mpcc,omitempty"`
	// Quantity: Reflects the amount of items available
	Quantity *float64 `json:"quantity,omitempty"`
	// Region: 2-letter ISO code of the country/region where the product is
	// stored
	Region string `json:"region,omitempty"`
	// Updated: Update date given by the merchant i.e. Q4/2022, 2022/10/12
	Updated string `json:"updated,omitempty"`
	// ZipCode: Zip code where the product is stored
	ZipCode string `json:"zipCode,omitempty"`
}

// UpsertResponse is the outcome of a successful request to upsert an
// availability.
type UpsertResponse struct {
	// Kind describes this entity, it will be
	// store#availability/upsertResponse.
	Kind string `json:"kind,omitempty"`
	// Link includes the URL where this resource will be available
	Link string `json:"link,omitempty"`
}

// Delete availability information of a product. It is an asynchronous
// operation.
type DeleteService struct {
	s    *Service
	opt_ map[string]interface{}
	hdr_ map[string]interface{}
	spn  string
}

// NewDeleteService creates a new instance of DeleteService.
func NewDeleteService(s *Service) *DeleteService {
	rs := &DeleteService{s: s, opt_: make(map[string]interface{}), hdr_: make(map[string]interface{})}
	return rs
}

// 2-letter ISO code of the country/region where the product is stored
func (s *DeleteService) Region(region string) *DeleteService {
	s.opt_["region"] = region
	return s
}

// SPN is the unique identifier of a product within a merchant.
func (s *DeleteService) Spn(spn string) *DeleteService {
	s.spn = spn
	return s
}

// Zip code where the product is stored
func (s *DeleteService) ZipCode(zipCode string) *DeleteService {
	s.opt_["zipCode"] = zipCode
	return s
}

// Do executes the operation.
func (s *DeleteService) Do(ctx context.Context) (*DeleteResponse, error) {
	var body io.Reader
	params := make(map[string]interface{})
	if v, ok := s.opt_["region"]; ok {
		params["region"] = v
	}
	params["spn"] = s.spn
	if v, ok := s.opt_["zipCode"]; ok {
		params["zipCode"] = v
	}
	path, err := meplatoapi.Expand("/api/v2/products/{spn}/availabilities{?region,zipCode}", params)
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
	ret := new(DeleteResponse)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
		return nil, err
	}
	return ret, nil
}

// Read availability information of a product
type GetService struct {
	s    *Service
	opt_ map[string]interface{}
	hdr_ map[string]interface{}
	spn  string
}

// NewGetService creates a new instance of GetService.
func NewGetService(s *Service) *GetService {
	rs := &GetService{s: s, opt_: make(map[string]interface{}), hdr_: make(map[string]interface{})}
	return rs
}

// 2-letter ISO code of the country/region where the product is stored
func (s *GetService) Region(region string) *GetService {
	s.opt_["region"] = region
	return s
}

// SPN is the unique identifier of a product within a merchant.
func (s *GetService) Spn(spn string) *GetService {
	s.spn = spn
	return s
}

// Zip code where the product is stored
func (s *GetService) ZipCode(zipCode string) *GetService {
	s.opt_["zipCode"] = zipCode
	return s
}

// Do executes the operation.
func (s *GetService) Do(ctx context.Context) (*GetResponse, error) {
	var body io.Reader
	params := make(map[string]interface{})
	if v, ok := s.opt_["region"]; ok {
		params["region"] = v
	}
	params["spn"] = s.spn
	if v, ok := s.opt_["zipCode"]; ok {
		params["zipCode"] = v
	}
	path, err := meplatoapi.Expand("/api/v2/products/{spn}/availabilities{?region,zipCode}", params)
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
	ret := new(GetResponse)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
		return nil, err
	}
	return ret, nil
}

// Update or create availability information of a product. It is an
// asynchronous operation.
type UpsertService struct {
	s            *Service
	opt_         map[string]interface{}
	hdr_         map[string]interface{}
	spn          string
	availability *UpsertRequest
}

// NewUpsertService creates a new instance of UpsertService.
func NewUpsertService(s *Service) *UpsertService {
	rs := &UpsertService{s: s, opt_: make(map[string]interface{}), hdr_: make(map[string]interface{})}
	return rs
}

// Availability properties of the product.
func (s *UpsertService) Availability(availability *UpsertRequest) *UpsertService {
	s.availability = availability
	return s
}

// SPN is the unique identifier of a product within a merchant.
func (s *UpsertService) Spn(spn string) *UpsertService {
	s.spn = spn
	return s
}

// Do executes the operation.
func (s *UpsertService) Do(ctx context.Context) (*UpsertResponse, error) {
	var body io.Reader
	body, err := meplatoapi.ReadJSON(s.availability)
	if err != nil {
		return nil, err
	}
	params := make(map[string]interface{})
	params["spn"] = s.spn
	path, err := meplatoapi.Expand("/api/v2/products/{spn}/availabilities", params)
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
	ret := new(UpsertResponse)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
		return nil, err
	}
	return ret, nil
}

// Copyright (c) 2015-2016 Meplato GmbH, Switzerland.
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

// Package jobs implements the Meplato Store API.
//
// See https://developer.meplato.com/store2/.
package jobs

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
	title   = "Meplato Store API"
	version = "2.1.3"
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

func (s *Service) Get() *GetService {
	return NewGetService(s)
}

func (s *Service) Search() *SearchService {
	return NewSearchService(s)
}

// Job that processes a task in the background, e.g. publishing a catalog.
type Job struct {
	// CatalogID: ID of the catalog.
	CatalogID int64 `json:"catalogId,omitempty"`
	// CatalogName: Name of the catalog.
	CatalogName string `json:"catalogName,omitempty"`
	// Completed is the date and time when the job has been completed, either
	// successfully or failed.
	Completed *time.Time `json:"completed,omitempty"`
	// Created is the creation date and time of the job.
	Created *time.Time `json:"created,omitempty"`
	// Email of the user that initiated the job.
	Email string `json:"email,omitempty"`
	// ID is a unique (internal) identifier of the job.
	ID string `json:"id,omitempty"`
	// Kind is store#job for a job entity.
	Kind string `json:"kind,omitempty"`
	// MerchantID: ID of the merchant.
	MerchantID int64 `json:"merchantId,omitempty"`
	// MerchantMpcc: MPCC of the merchant.
	MerchantMpcc string `json:"merchantMpcc,omitempty"`
	// MerchantName: Name of the merchant.
	MerchantName string `json:"merchantName,omitempty"`
	// SelfLink: URL to this page.
	SelfLink string `json:"selfLink,omitempty"`
	// Started is the date and time when the job has been started.
	Started *time.Time `json:"started,omitempty"`
	// State describes the current state of the job, i.e. one of
	// waiting,working,succeeded, or failed.
	State string `json:"state,omitempty"`
	// Topic of the job, e.g. if it was an import or a validation task.
	Topic string `json:"topic,omitempty"`
}

// SearchResponse is a partial listing of jobs.
type SearchResponse struct {
	// Items is the slice of jobs of this result.
	Items []*Job `json:"items,omitempty"`
	// Kind is store#jobs for this kind of response.
	Kind string `json:"kind,omitempty"`
	// NextLink returns the URL to the next slice of jobs (if any).
	NextLink string `json:"nextLink,omitempty"`
	// PreviousLink returns the URL of the previous slice of jobs (if any).
	PreviousLink string `json:"previousLink,omitempty"`
	// SelfLink returns the URL to this page.
	SelfLink string `json:"selfLink,omitempty"`
	// TotalItems describes the total number of jobs found.
	TotalItems int64 `json:"totalItems,omitempty"`
}

// Get a single job.
type GetService struct {
	s    *Service
	opt_ map[string]interface{}
	hdr_ map[string]interface{}
	id   string
}

// NewGetService creates a new instance of GetService.
func NewGetService(s *Service) *GetService {
	rs := &GetService{s: s, opt_: make(map[string]interface{}), hdr_: make(map[string]interface{})}
	return rs
}

// ID of the job.
func (s *GetService) ID(id string) *GetService {
	s.id = id
	return s
}

// Do executes the operation.
func (s *GetService) Do() (*Job, error) {
	var body io.Reader
	params := make(map[string]interface{})
	params["id"] = s.id
	path, err := meplatoapi.Expand("/jobs/{id}", params)
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
	ret := new(Job)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
		return nil, err
	}
	return ret, nil
}

// Search for jobs.
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

// State filter, e.g. waiting,working,succeeded,failed.
func (s *SearchService) State(state string) *SearchService {
	s.opt_["state"] = state
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
	if v, ok := s.opt_["state"]; ok {
		params["state"] = v
	}
	if v, ok := s.opt_["take"]; ok {
		params["take"] = v
	}
	path, err := meplatoapi.Expand("/jobs{?merchantId,skip,take,state}", params)
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

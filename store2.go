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

// Package store2 implements the Meplato Store 2 API.
//
// See https://developer.meplato.com/store2/.
package store2

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

func (s *Service) Me() *MeService {
	return NewMeService(s)
}

func (s *Service) Ping() *PingService {
	return NewPingService(s)
}

// MeResponse returns various information about the user and endpoints.
type MeResponse struct {
	// CatalogsLink is the URL for retrieving the list of catalogs.
	CatalogsLink string `json:"catalogsLink,omitempty"`
	// Kind is store#me for this entity.
	Kind string `json:"kind,omitempty"`
	// Merchant returns information about your merchant account.
	Merchant *Merchant `json:"merchant,omitempty"`
	// SelfLink is the URL of this request.
	SelfLink string `json:"selfLink,omitempty"`
	// User returns information about your user account.
	User *User `json:"user,omitempty"`
}

// Merchant holds account data for the merchant/supplier in Meplato Store.
type Merchant struct {
	// Country is the ISO code for the country of the merchant, e.g. DE or CH.
	Country string `json:"country,omitempty"`
	// Created is the date/time when the merchant was created, e.g.
	// 2015-03-19T12:09:45Z
	Created *time.Time `json:"created,omitempty"`
	// Currency is the default ISO code for new catalogs, e.g. EUR or GBP.
	Currency string `json:"currency,omitempty"`
	// ID is a unique (internal) identifier of the merchant.
	ID int64 `json:"id,omitempty"`
	// Kind is store#merchant for this entity.
	Kind string `json:"kind,omitempty"`
	// Language is the code for the language of the merchant, e.g. de or en.
	Language string `json:"language,omitempty"`
	// Locale is the regional code in the format de_AT.
	Locale string `json:"locale,omitempty"`
	// Mpcc: MPCC is the Meplato Company Code, a unique identifier.
	Mpcc string `json:"mpcc,omitempty"`
	// Mpsc: MPSC is the Meplato Supplier Code.
	Mpsc string `json:"mpsc,omitempty"`
	// Name is the name of the merchant.
	Name string `json:"name,omitempty"`
	// Ou: OU is the default ISO code of the order unit, e.g. PCE or EA.
	Ou string `json:"ou,omitempty"`
	// SelfLink is the URL for this merchant.
	SelfLink string `json:"selfLink,omitempty"`
	// TimeZone is the time zone in the format Europe/Berlin.
	TimeZone string `json:"timeZone,omitempty"`
	// Token is a shared token for this merchant.
	Token string `json:"token,omitempty"`
	// Updated is the date/time when the merchant was last modified, e.g.
	// 2015-03-19T12:09:45Z
	Updated *time.Time `json:"updated,omitempty"`
}

// User holds account data for the user in Meplato Store.
type User struct {
	// Country is the ISO code for the country, e.g. DE or CH.
	Country string `json:"country,omitempty"`
	// Created is the date/time when the user was created, e.g.
	// 2015-03-19T12:09:45Z
	Created *time.Time `json:"created,omitempty"`
	// Currency is the default ISO code for currencies, e.g. EUR or GBP.
	Currency string `json:"currency,omitempty"`
	// Email is the email address.
	Email string `json:"email,omitempty"`
	// ID is a unique (internal) identifier of the user.
	ID int64 `json:"id,omitempty"`
	// Kind is store#user for this entity.
	Kind string `json:"kind,omitempty"`
	// Language is the code for the language, e.g. de or en.
	Language string `json:"language,omitempty"`
	// Locale is the regional code in the format de_AT.
	Locale string `json:"locale,omitempty"`
	// Name is the user, including first and last name.
	Name string `json:"name,omitempty"`
	// Provider is used internally.
	Provider string `json:"provider,omitempty"`
	// TimeZone is the time zone in the format Europe/Berlin.
	TimeZone string `json:"timeZone,omitempty"`
	// Uid: UID is used internally.
	Uid string `json:"uid,omitempty"`
	// Updated is the date/time when the user was last modified, e.g.
	// 2015-03-19T12:09:45Z
	Updated *time.Time `json:"updated,omitempty"`
}

// Me returns information about your user profile and the API endpoints of
// the Meplato Store 2.0 API.
type MeService struct {
	s    *Service
	opt_ map[string]interface{}
	hdr_ map[string]interface{}
}

// NewMeService creates a new instance of MeService.
func NewMeService(s *Service) *MeService {
	rs := &MeService{s: s, opt_: make(map[string]interface{}), hdr_: make(map[string]interface{})}
	return rs
}

// Do executes the operation.
func (s *MeService) Do() (*MeResponse, error) {
	var body io.Reader
	path := "/"
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
	ret := new(MeResponse)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
		return nil, err
	}
	return ret, nil
}

// Ping allows you to test if the Meplato Store 2.0 API is currently
// operational.
type PingService struct {
	s    *Service
	opt_ map[string]interface{}
	hdr_ map[string]interface{}
}

// NewPingService creates a new instance of PingService.
func NewPingService(s *Service) *PingService {
	rs := &PingService{s: s, opt_: make(map[string]interface{}), hdr_: make(map[string]interface{})}
	return rs
}

// Do executes the operation.
func (s *PingService) Do() error {
	var body io.Reader
	path := "/"
	req, err := http.NewRequest("HEAD", s.s.BaseURL+path, body)
	if err != nil {
		return err
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
		return err
	}
	defer meplatoapi.CloseBody(res)
	if err := meplatoapi.CheckResponse(res); err != nil {
		return err
	}
	return nil
}

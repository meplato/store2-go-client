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

// Package store2 implements the Meplato Store API.
//
// See https://developer.meplato.com/store2/.
package store2

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"runtime"
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
	version = "2.1.9"
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
		client = &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
				DialContext: (&net.Dialer{
					Timeout:   30 * time.Second,
					KeepAlive: 30 * time.Second,
					DualStack: true,
				}).DialContext,
				MaxIdleConns:          100,
				IdleConnTimeout:       90 * time.Second,
				TLSHandshakeTimeout:   10 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,
				MaxIdleConnsPerHost:   runtime.GOMAXPROCS(0) + 1,
				TLSClientConfig: &tls.Config{
					MinVersion: tls.VersionTLS12,
				},
			},
		}
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
	// Country/Region is the ISO code for the country/region of the merchant,
	// e.g. DE or CH.
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
	// SelfService indicates whether this merchant is on self-service or
	// managed by Meplato.
	SelfService bool `json:"selfService,omitempty"`
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
	// Country/Region is the ISO code for the country/region, e.g. DE or CH.
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
	// MerchantID: merchantId is a unique (internal) identifier of the
	// merchant of the user.
	MerchantID int64 `json:"merchantId,omitempty"`
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
func (s *MeService) Do(ctx context.Context) (*MeResponse, error) {
	var body io.Reader
	path := "/"
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
func (s *PingService) Do(ctx context.Context) error {
	var body io.Reader
	path := "/"
	req, err := http.NewRequest("HEAD", s.s.BaseURL+path, body)
	if err != nil {
		return err
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
		return err
	}
	defer meplatoapi.CloseBody(res)
	if err := meplatoapi.CheckResponse(res); err != nil {
		return err
	}
	return nil
}

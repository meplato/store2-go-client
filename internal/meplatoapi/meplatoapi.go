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
package meplatoapi

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"runtime"
)

const (
	Version   = "2.0"
	UserAgent = "meplato-store-go-client/" + Version + " (" + runtime.GOOS + "/" + runtime.GOARCH + ")"
)

// Error contains an error response from the server.
type Error struct {
	// Code is the HTTP response status code and will always be populated.
	Code int `json:"code"`
	// Message is the server response message.
	Message string `json:"message,omitempty"`
	// Details contains error details.
	Details []string `json:"details,omitempty"`
	// Body is the response body.
	Body string
}

func (e *Error) Error() string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "meplatoapi: Error %d: ", e.Code)
	if e.Message != "" {
		fmt.Fprintf(&buf, "%s", e.Message)
	}
	return buf.String()
}

type errorReply struct {
	Error *Error `json:"error"`
}

// CheckResponse returns an error (of type *Error) if the response status
// code is not 2xx.
func CheckResponse(res *http.Response) error {
	if res.StatusCode >= 200 && res.StatusCode <= 299 {
		return nil
	}
	slurp, err := ioutil.ReadAll(res.Body)
	if err == nil {
		jerr := new(errorReply)
		err = json.Unmarshal(slurp, jerr)
		if err == nil && jerr.Error != nil {
			if jerr.Error.Code == 0 {
				jerr.Error.Code = res.StatusCode
			}
			jerr.Error.Body = string(slurp)
			return jerr.Error
		}
	}
	return &Error{
		Code: res.StatusCode,
		Body: string(slurp),
	}
}

func ReadJSON(v interface{}) (io.Reader, error) {
	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(v)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func Expand(path string, expansions map[string]interface{}) (string, error) {
	template, err := Parse(path)
	if err != nil {
		return "", err
	}
	return template.Expand(expansions)
}

func CloseBody(res *http.Response) {
	if res == nil || res.Body == nil {
		return
	}
	res.Body.Close()
}

func HTTPBasicAuthorizationHeader(user, pass string) string {
	s := user + ":" + pass
	return fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(s)))
}

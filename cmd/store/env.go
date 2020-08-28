package main

import (
	"crypto/tls"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"path"
	"runtime"
	"time"

	"github.com/bgentry/go-netrc/netrc"

	"github.com/meplato/store2-go-client/v2/catalogs"
	"github.com/meplato/store2-go-client/v2/products"
)

func GetBaseURL() string {
	if url := os.Getenv("STORE_URL"); url != "" {
		return url
	}
	if url := os.Getenv("STORE2_URL"); url != "" {
		return url
	}
	return "https://store.meplato.com/api/v2"
}

func getUsername() string {
	if s := os.Getenv("STORE_USER"); s != "" {
		return s
	}
	if s := os.Getenv("STORE2_USER"); s != "" {
		return s
	}

	// Retrieve from .netrc
	username, _ := getLoginAndPasswordFromNetrc(GetBaseURL())
	return username
}

func getPassword() string {
	if s := os.Getenv("STORE_PASSWORD"); s != "" {
		return s
	}
	if s := os.Getenv("STORE2_PASSWORD"); s != "" {
		return s
	}

	// Retrieve from .netrc
	_, password := getLoginAndPasswordFromNetrc(GetBaseURL())
	return password
}

func getLoginAndPasswordFromNetrc(serviceEndpoint string) (username, password string) {
	username = ""
	password = ""

	// Get hostname from BaseURL
	u, err := url.Parse(serviceEndpoint)
	if err != nil {
		return
	}
	machine := u.Host

	// Get user's home directory to find ~/.netrc
	user, err := user.Current()
	if err != nil {
		return
	}
	netrcfile := path.Join(user.HomeDir, ".netrc")

	// Find entry in .netrc for the given machine
	m, err := netrc.FindMachine(netrcfile, machine)
	if err != nil {
		return
	}
	if m != nil {
		username = m.Login
		password = m.Password
	}
	return
}

func GetHttpClient() (*http.Client, error) {
	client := &http.Client{
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
				MinVersion:         tls.VersionTLS12,
				InsecureSkipVerify: true,
			},
		},
	}
	return client, nil
}

func GetCatalogsService() (*catalogs.Service, error) {
	client, err := GetHttpClient()
	if err != nil {
		return nil, err
	}
	service, err := catalogs.New(client)
	if err != nil {
		return nil, err
	}
	if url := GetBaseURL(); url != "" {
		service.BaseURL = url
	}
	service.User = getUsername()
	service.Password = getPassword()
	return service, nil
}

func GetProductsService() (*products.Service, error) {
	client, err := GetHttpClient()
	if err != nil {
		return nil, err
	}
	service, err := products.New(client)
	if err != nil {
		return nil, err
	}
	if url := GetBaseURL(); url != "" {
		service.BaseURL = url
	}
	service.User = getUsername()
	service.Password = getPassword()
	return service, nil
}

package main

import (
	"net/http"
	"net/url"
	"os"
	"os/user"
	"path"

	"github.com/bgentry/go-netrc/netrc"

	"github.com/meplato/store2-go-client/catalogs"
	"github.com/meplato/store2-go-client/products"
)

func GetBaseURL() string {
	if url := os.Getenv("STORE2_URL"); url != "" {
		return url
	}
	return "https://store2.meplato.com/api/v2"
}

func getUsername() string {
	if s := os.Getenv("STORE2_USER"); s != "" {
		return s
	}

	// Retrieve from .netrc
	username, _ := getLoginAndPasswordFromNetrc(GetBaseURL())
	return username
}

func getPassword() string {
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
	return http.DefaultClient, nil
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

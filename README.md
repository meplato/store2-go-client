# Meplato Store 2 API Go Client

This is the Go client for the Meplato Store 2 API. It consists of a library
to integrate your infrastructure with Meplato suite for suppliers. It also
has a command line tool you can use to immediately interact with Meplato Store.
The command line tool also illustrates how to write your own clients.

## Prerequisites

You need at two things to use the Meplato Store 2 API.

1. A login to Meplato Store 2.
2. An API token.

Get your login by contacting Meplato Supplier Network Services. The API token
is required to securely communicate with the Meplato Store 2 API. You can
find it in the personalization section when logged into Meplato Store.

## Getting started

We use the command line tool to illustrate how easy it is to get started
with Meplato Store. First, we need to clone the repository and build the
command line client:

1. Get the library via: `go get github.com/meplato/store2-go-client`
2. Install the dependencies: `make deps`
3. Set the API token, e.g. via an environment variable:
   `export STORE2_USER=<your-api-token>`.
   If you're an old Unix guy, you can also add a line to your `~/.netrc`
   file, e.g. `machine store2.meplato.com login <your-api-token>`; the
   command line client will happily pick it up without any environment
   variables.
4. Build the command line client: `make`

You should now have an executable binary at `./store`. If you run it, you
can see a help page regarding the various commands.

If you successfully set your API token, either via environment variable or
via netrc, you should now be able to run `./store catalogs` and it should
list all the catalogs in your Meplato Store.

```bash
$ ./store catalogs
4 catalogs found.
 ID  Name                                               Created
==============================================================================
  4. Excel-Test                                         2015-07-17
  3. Demokatalog                                        2015-06-18
  2. Büromaterial                                       2015-06-18
  1. Büromaterial                                       2015-06-18
```

## Using the library

Using the library is actually quite simple. All functionality is separated
into services. So you e.g. have a service to work with catalogs, another
service to work with products in a catalog etc. All services need to be
initialized with your API token.

```go
import (
	"net/http"
	"log"

	catalogs "github.com/meplato/store2-go-client/catalogs"
)

...

// Create and initialize your service with your API token
service, err := catalogs.New(http.DefaultClient)
if err != nil {
  log.Fatal(err)
}
service.User = "<your-api-token>"
```

Now that you have access to your service, you can set up parameters and
execute the service call. For example, the following snippet will print
the first 10 catalogs in your Meplato Store, sorted by catalog name.

```go
res, err := catalogs.Search().Skip(0).Take(10).Sort("name").Do()
if err != nil {
	log.Fatal(err)
}

fmt.Printf("You have a total of %d catalog(s).\n", res.TotalItems)
for _, c := range res.Items {
	fmt.Printf("Catalog with ID=%d has name %q and was created at %v.\n",
		c.ID, c.Name, c.Created.Format("2006-01-02"))
}
```

Feel free to read the unit tests for the various usage scenarios of the
library.

## Running tests

To run all tests use `go test ./...`

## Documentation

Complete documentation for the Meplato Store 2 API can be found at
[https://developer.meplato.com/store2](https://developer.meplato.com/store2).

# License

This software is licensed under the Apache 2 license.

    Copyright (c) 2015 Meplato GmbH, Switzerland <http://www.meplato.com>

		Licensed under the Apache License, Version 2.0 (the "License");
		you may not use this file except in compliance with the License.
		You may obtain a copy of the License at

		    http://www.apache.org/licenses/LICENSE-2.0

		Unless required by applicable law or agreed to in writing, software
		distributed under the License is distributed on an "AS IS" BASIS,
		WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
		See the License for the specific language governing permissions and
		limitations under the License.

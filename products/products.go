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

// Package products implements the Meplato Store 2 API.
//
// See https://developer.meplato.com/store2/.
package products

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

func (s *Service) Create() *CreateService {
	return NewCreateService(s)
}

func (s *Service) Delete() *DeleteService {
	return NewDeleteService(s)
}

func (s *Service) Get() *GetService {
	return NewGetService(s)
}

func (s *Service) Replace() *ReplaceService {
	return NewReplaceService(s)
}

func (s *Service) Scroll() *ScrollService {
	return NewScrollService(s)
}

func (s *Service) Search() *SearchService {
	return NewSearchService(s)
}

func (s *Service) Update() *UpdateService {
	return NewUpdateService(s)
}

func (s *Service) Upsert() *UpsertService {
	return NewUpsertService(s)
}

// Availability details product availability.
type Availability struct {
	// Message gives a textual description of the availability message, e.g.
	// "in stock", "out of stock", "limited availability", or "on display to
	// order".
	Message string `json:"message,omitempty"`
	// Qty describes the quantity for certain kinds of availability messages.
	// E.g. you can indicate the number of items in stock.
	Qty *float64 `json:"qty,omitempty"`
	// Updated indicates when the availability message has been last updated.
	Updated string `json:"updated,omitempty"`
}

// Blob describes external product data, e.g. an image or a datasheet.
type Blob struct {
	// Kind describes the type of blob, e.g. image, detail, thumbnail,
	// datasheet, or safetysheet.
	Kind string `json:"kind,omitempty"`
	// Language indicates the language of the blob, e.g. the language of a
	// datasheet.
	Language string `json:"lang,omitempty"`
	// Source is either a (relative) file name or a URL.
	Source string `json:"source,omitempty"`
	// Text gives a textual description the blob.
	Text string `json:"text,omitempty"`
}

// Condition describes a product status, e.g. refurbished or used.
type Condition struct {
	// Kind describes the condition, e.g. bargain, new_product, old_product,
	// new, used, refurbished, or core_product.
	Kind string `json:"kind,omitempty"`
	// Text gives a textual description of the condition.
	Text string `json:"text,omitempty"`
}

// CreateProduct holds the properties of a new product.
type CreateProduct struct {
	// Availability allows the update of product availability data, e.g. the
	// number of items in stock or the date when the product will be available
	// again.
	Availability *Availability `json:"availability,omitempty"`
	// Blobs specifies external data, e.g. images or datasheets.
	Blobs []*Blob `json:"blobs,omitempty"`
	// Bpn: BPN is the buyer part number of the product.
	Bpn string `json:"bpn,omitempty"`
	// Categories is a list of (supplier-specific) category names the product
	// belongs to.
	Categories []string `json:"categories,omitempty"`
	// Conditions describes the product conditions, e.g. refurbished or used.
	Conditions []*Condition `json:"conditions,omitempty"`
	// ContentUnit is the content unit of the product, a 3-character ISO code
	// (usually project-specific).
	ContentUnit string `json:"cu,omitempty"`
	// CuPerOu describes the number of content units per order unit, e.g. the
	// 12 in '1 case contains 12 bottles'.
	CuPerOu *float64 `json:"cuPerOu,omitempty"`
	// CustField1 is the CUST_FIELD1 of the SAP OCI specification. It has a
	// maximum length of 10 characters.
	CustField1 string `json:"custField1,omitempty"`
	// CustField2 is the CUST_FIELD2 of the SAP OCI specification. It has a
	// maximum length of 10 characters.
	CustField2 string `json:"custField2,omitempty"`
	// CustField3 is the CUST_FIELD3 of the SAP OCI specification. It has a
	// maximum length of 10 characters.
	CustField3 string `json:"custField3,omitempty"`
	// CustField4 is the CUST_FIELD4 of the SAP OCI specification. It has a
	// maximum length of 20 characters.
	CustField4 string `json:"custField4,omitempty"`
	// CustField5 is the CUST_FIELD5 of the SAP OCI specification. It has a
	// maximum length of 50 characters.
	CustField5 string `json:"custField5,omitempty"`
	// CustFields is an array of generic name/value pairs for
	// customer-specific attributes.
	CustFields []*CustField `json:"custFields,omitempty"`
	// Datasheet is the name of an datasheet file (in the media files) or a
	// URL to the datasheet on the internet.
	Datasheet string `json:"datasheet,omitempty"`
	// Description of the product.
	Description string `json:"description,omitempty"`
	// Eclasses is a list of eCl@ss categories the product belongs to.
	Eclasses []*Eclass `json:"eclasses,omitempty"`
	// ErpGroupSupplier: erpGroupSupplier is the material group of the product
	// on the merchant-/supplier-side.
	ErpGroupSupplier string `json:"erpGroupSupplier,omitempty"`
	// Excluded is a flag that indicates whether to exclude this product from
	// the catalog. If true, this product will not be published into the live
	// area.
	Excluded bool `json:"excluded,omitempty"`
	// ExtCategory is the EXT_CATEGORY field of the SAP OCI specification.
	ExtCategory string `json:"extCategory,omitempty"`
	// ExtCategoryID is the EXT_CATEGORY_ID field of the SAP OCI
	// specification.
	ExtCategoryID string `json:"extCategoryId,omitempty"`
	// ExtSchemaType is the EXT_SCHEMA_TYPE field of the SAP OCI
	// specification.
	ExtSchemaType string `json:"extSchemaType,omitempty"`
	// Features defines product features, i.e. additional properties of the
	// product.
	Features []*Feature `json:"features,omitempty"`
	// Gtin: GTIN is the global trade item number of the product (used to be
	// EAN).
	Gtin string `json:"gtin,omitempty"`
	// Hazmats classifies hazardous/dangerous goods.
	Hazmats []*Hazmat `json:"hazmats,omitempty"`
	// Image is the name of an image file (in the media files) or a URL to the
	// image on the internet.
	Image string `json:"image,omitempty"`
	// Keywords is a list of aliases for the product.
	Keywords []string `json:"keywords,omitempty"`
	// Leadtime is the number of days for delivery.
	Leadtime *float64 `json:"leadtime,omitempty"`
	// ListPrice is the net list price of the product.
	ListPrice *float64 `json:"listPrice,omitempty"`
	// Manufactcode is the manufacturer code as used in the SAP OCI
	// specification.
	Manufactcode string `json:"manufactcode,omitempty"`
	// Manufacturer is the name of the manufacturer.
	Manufacturer string `json:"manufacturer,omitempty"`
	// Matgroup is the material group of the product on the buy-side.
	Matgroup string `json:"matgroup,omitempty"`
	// Mpn: MPN is the manufacturer part number.
	Mpn string `json:"mpn,omitempty"`
	// Name of the product.
	Name string `json:"name,omitempty"`
	// OrderUnit is the order unit of the product, a 3-character ISO code
	// (usually project-specific).
	OrderUnit string `json:"ou,omitempty"`
	// Price is the net price (per order unit) of the product for the
	// end-user.
	Price float64 `json:"price,omitempty"`
	// PriceQty is the quantity for which the price is specified (default:
	// 1.0).
	PriceQty *float64 `json:"priceQty,omitempty"`
	// QuantityInterval is the interval in which this product can be ordered.
	// E.g. if the quantity interval is 5, the end-user can only order in
	// quantities of 5,10,15 etc.
	QuantityInterval *float64 `json:"quantityInterval,omitempty"`
	// QuantityMax is the maximum order quantity for this product.
	QuantityMax *float64 `json:"quantityMax,omitempty"`
	// QuantityMin is the minimum order quantity for this product.
	QuantityMin *float64 `json:"quantityMin,omitempty"`
	// References defines cross-product references, e.g. alternatives or
	// follow-up products.
	References []*Reference `json:"references,omitempty"`
	// Safetysheet is the name of an safetysheet file (in the media files) or
	// a URL to the safetysheet on the internet.
	Safetysheet string `json:"safetysheet,omitempty"`
	// ScalePrices can be used when the price of the product is dependent on
	// the ordered quantity.
	ScalePrices []*ScalePrice `json:"scalePrices,omitempty"`
	// Service indicates if the is a good (false) or a service (true). The
	// default value is false.
	Service bool `json:"service,omitempty"`
	// Spn: SPN is the supplier part number.
	Spn string `json:"spn,omitempty"`
	// TaxCode to use for this product. This is typically project-specific.
	TaxCode string `json:"taxCode,omitempty"`
	// Thumbnail is the name of an thumbnail image file (in the media files)
	// or a URL to the image on the internet.
	Thumbnail string `json:"thumbnail,omitempty"`
	// Unspscs is a list of UNSPSC categories the product belongs to.
	Unspscs []*Unspsc `json:"unspscs,omitempty"`
}

// CreateProductResponse is the outcome of a successful request to create
// a product.
type CreateProductResponse struct {
	// Kind describes this entity.
	Kind string `json:"kind,omitempty"`
	// Link returns a URL to the representation of the newly created product.
	Link string `json:"link,omitempty"`
}

// CustField describes a generic name/value pair. Its purpose is to
// provide a mechanism for customer-specific fields.
type CustField struct {
	// Name is the name of the customer-specific field, e.g. TaxRate.
	Name string `json:"name,omitempty"`
	// Value is the value of the customer-specific field, e.g. 19%%.
	Value string `json:"value,omitempty"`
}

// Eclass is used to tie a product to an eCl@ss schema.
type Eclass struct {
	// Code is the eCl@ss code. Only use digits for encoding, e.g. 19010203.
	Code string `json:"code,omitempty"`
	// Version is the eCl@ss version in the major.minor format, e.g. 5.1 or
	// 7.0.
	Version string `json:"version,omitempty"`
}

// Feature describes additional features of a product.
type Feature struct {
	// Kind describes the type of feature, e.g. ECLASS-5.1 to describe a
	// feature of eCl@ss 5.1.
	Kind string `json:"kind,omitempty"`
	// Name specifies the name/code of the feature. If you are refering to a
	// standard classification like eCl@ss, you should use the code of the
	// feature. Otherwise, you are free to use a textual description like
	// "Weight" or "Diameter" or "Voltage".
	Name string `json:"name,omitempty"`
	// Unit is a unit of measurement to describe the feature value. E.g. if
	// you specify the weight, you should probably use KGM as unit to describe
	// that the weight is given in kilograms.
	Unit string `json:"unit,omitempty"`
	// Values describes the feature value(s).
	Values []string `json:"values,omitempty"`
}

// Hazmat describes a hazardous/dangerous good.
type Hazmat struct {
	// Kind describes the classification system, e.g. GGVS.
	Kind string `json:"kind,omitempty"`
	// Text gives a textual description or a code of the harm that the product
	// can do to people, property, or the environment.
	Text string `json:"text,omitempty"`
}

// Product is a good or service in a catalog.
type Product struct {
	// Availability allows the update of product availability data, e.g. the
	// number of items in stock or the date when the product will be available
	// again.
	Availability *Availability `json:"availability,omitempty"`
	// Blobs specifies external data, e.g. images or datasheets.
	Blobs []*Blob `json:"blobs,omitempty"`
	// Bpn: BPN is the buyer part number of the product.
	Bpn string `json:"bpn,omitempty"`
	// CatalogID: ID of the catalog this products belongs to.
	CatalogID int64 `json:"catalogId,omitempty"`
	// Categories is a list of (supplier-specific) category names the product
	// belongs to.
	Categories []string `json:"categories,omitempty"`
	// Conditions describes the product conditions, e.g. refurbished or used.
	Conditions []*Condition `json:"conditions,omitempty"`
	// Created is the creation date and time of the product.
	Created *time.Time `json:"created,omitempty"`
	// ContentUnit is the content unit of the product, a 3-character ISO code
	// (usually project-specific).
	ContentUnit string `json:"cu,omitempty"`
	// CuPerOu describes the number of content units per order unit, e.g. the
	// 12 in '1 case contains 12 bottles'.
	CuPerOu float64 `json:"cuPerOu,omitempty"`
	// Currency is the ISO currency code for the prices, e.g. EUR or GBP.
	Currency string `json:"currency,omitempty"`
	// CustField1 is the CUST_FIELD1 of the SAP OCI specification. It has a
	// maximum length of 10 characters.
	CustField1 string `json:"custField1,omitempty"`
	// CustField2 is the CUST_FIELD2 of the SAP OCI specification. It has a
	// maximum length of 10 characters.
	CustField2 string `json:"custField2,omitempty"`
	// CustField3 is the CUST_FIELD3 of the SAP OCI specification. It has a
	// maximum length of 10 characters.
	CustField3 string `json:"custField3,omitempty"`
	// CustField4 is the CUST_FIELD4 of the SAP OCI specification. It has a
	// maximum length of 20 characters.
	CustField4 string `json:"custField4,omitempty"`
	// CustField5 is the CUST_FIELD5 of the SAP OCI specification. It has a
	// maximum length of 50 characters.
	CustField5 string `json:"custField5,omitempty"`
	// CustFields is an array of generic name/value pairs for
	// customer-specific attributes.
	CustFields []*CustField `json:"custFields,omitempty"`
	// Datasheet is the name of an datasheet file (in the media files) or a
	// URL to the datasheet on the internet.
	Datasheet string `json:"datasheet,omitempty"`
	// DatasheetURL is the URL to the data sheet (if available).
	DatasheetURL string `json:"datasheetURL,omitempty"`
	// Description of the product.
	Description string `json:"description,omitempty"`
	// Eclasses is a list of eCl@ss categories the product belongs to.
	Eclasses []*Eclass `json:"eclasses,omitempty"`
	// ErpGroupSupplier: erpGroupSupplier is the material group of the product
	// on the merchant-/supplier-side.
	ErpGroupSupplier string `json:"erpGroupSupplier,omitempty"`
	// Excluded is a flag that indicates whether to exclude this product from
	// the catalog. If true, this product will not be published into the live
	// area.
	Excluded bool `json:"excluded,omitempty"`
	// ExtCategory is the EXT_CATEGORY field of the SAP OCI specification.
	ExtCategory string `json:"extCategory,omitempty"`
	// ExtCategoryID is the EXT_CATEGORY_ID field of the SAP OCI
	// specification.
	ExtCategoryID string `json:"extCategoryId,omitempty"`
	// ExtProductID is the EXT_PRODUCT_ID field of the SAP OCI specification.
	ExtProductID string `json:"extProductId,omitempty"`
	// ExtSchemaType is the EXT_SCHEMA_TYPE field of the SAP OCI
	// specification.
	ExtSchemaType string `json:"extSchemaType,omitempty"`
	// Features defines product features, i.e. additional properties of the
	// product.
	Features []*Feature `json:"features,omitempty"`
	// Gtin: GTIN is the global trade item number of the product (used to be
	// EAN).
	Gtin string `json:"gtin,omitempty"`
	// Hazmats classifies hazardous/dangerous goods.
	Hazmats []*Hazmat `json:"hazmats,omitempty"`
	// ID is a unique (internal) identifier of the product.
	ID string `json:"id,omitempty"`
	// Image is the name of an image file (in the media files) or a URL to the
	// image on the internet.
	Image string `json:"image,omitempty"`
	// ImageURL is the URL to the image.
	ImageURL string `json:"imageURL,omitempty"`
	// Keywords is a list of aliases for the product.
	Keywords []string `json:"keywords,omitempty"`
	// Kind is store#product for a product entity.
	Kind string `json:"kind,omitempty"`
	// Leadtime is the number of days for delivery.
	Leadtime *float64 `json:"leadtime,omitempty"`
	// ListPrice is the net list price of the product.
	ListPrice float64 `json:"listPrice,omitempty"`
	// Manufactcode is the manufacturer code as used in the SAP OCI
	// specification.
	Manufactcode string `json:"manufactcode,omitempty"`
	// Manufacturer is the name of the manufacturer.
	Manufacturer string `json:"manufacturer,omitempty"`
	// Matgroup is the material group of the product on the buy-side.
	Matgroup string `json:"matgroup,omitempty"`
	// MeplatoPrice is the Meplato price of the product.
	MeplatoPrice float64 `json:"meplatoPrice,omitempty"`
	// MerchantID: ID of the merchant.
	MerchantID int64 `json:"merchantId,omitempty"`
	// Mpn: MPN is the manufacturer part number.
	Mpn string `json:"mpn,omitempty"`
	// Name of the product.
	Name string `json:"name,omitempty"`
	// OrderUnit is the order unit of the product, a 3-character ISO code
	// (usually project-specific).
	OrderUnit string `json:"ou,omitempty"`
	// Price is the net price (per order unit) of the product for the
	// end-user.
	Price float64 `json:"price,omitempty"`
	// PriceQty is the quantity for which the price is specified (default:
	// 1.0).
	PriceQty float64 `json:"priceQty,omitempty"`
	// ProjectID: ID of the project.
	ProjectID int64 `json:"projectId,omitempty"`
	// QuantityInterval is the interval in which this product can be ordered.
	// E.g. if the quantity interval is 5, the end-user can only order in
	// quantities of 5,10,15 etc.
	QuantityInterval *float64 `json:"quantityInterval,omitempty"`
	// QuantityMax is the maximum order quantity for this product.
	QuantityMax *float64 `json:"quantityMax,omitempty"`
	// QuantityMin is the minimum order quantity for this product.
	QuantityMin *float64 `json:"quantityMin,omitempty"`
	// References defines cross-product references, e.g. alternatives or
	// follow-up products.
	References []*Reference `json:"references,omitempty"`
	// Safetysheet is the name of an safetysheet file (in the media files) or
	// a URL to the safetysheet on the internet.
	Safetysheet string `json:"safetysheet,omitempty"`
	// SafetysheetURL is the URL to the safety data sheet (if available).
	SafetysheetURL string `json:"safetysheetURL,omitempty"`
	// ScalePrices can be used when the price of the product is dependent on
	// the ordered quantity.
	ScalePrices []*ScalePrice `json:"scalePrices,omitempty"`
	// SelfLink: URL to this page.
	SelfLink string `json:"selfLink,omitempty"`
	// Service indicates if the is a good (false) or a service (true). The
	// default value is false.
	Service bool `json:"service,omitempty"`
	// Spn: SPN is the supplier part number.
	Spn string `json:"spn,omitempty"`
	// TaxCode to use for this product. This is typically project-specific.
	TaxCode string `json:"taxCode,omitempty"`
	// Thumbnail is the name of an thumbnail image file (in the media files)
	// or a URL to the image on the internet.
	Thumbnail string `json:"thumbnail,omitempty"`
	// ThumbnailURL: ThubmnailURL is the URL to the thumbnail image.
	ThumbnailURL string `json:"thumbnailURL,omitempty"`
	// Unspscs is a list of UNSPSC categories the product belongs to.
	Unspscs []*Unspsc `json:"unspscs,omitempty"`
	// Updated is the last modification date and time of the product.
	Updated *time.Time `json:"updated,omitempty"`
}

// Reference describes a reference from one product to another product.
type Reference struct {
	// Kind describes the type of reference.
	Kind string `json:"kind,omitempty"`
	// Qty describes the quantity for certain kinds of references. E.g. the
	// consists_of kind must use the quantity field to be useful for the
	// end-user.
	Qty *float64 `json:"qty,omitempty"`
	// Spn: SPN specifies the supplier product number of the other product.
	Spn string `json:"spn,omitempty"`
}

// ReplaceProduct replace all properties of an existing product.
type ReplaceProduct struct {
	// Availability allows the update of product availability data, e.g. the
	// number of items in stock or the date when the product will be available
	// again.
	Availability *Availability `json:"availability,omitempty"`
	// Blobs contains information about external data, e.g. attachments like
	// images or datasheets.
	Blobs []*Blob `json:"blobs,omitempty"`
	// Bpn: BPN is the buyer part number of the product.
	Bpn string `json:"bpn,omitempty"`
	// Categories is a list of (supplier-specific) category names the product
	// belongs to.
	Categories []string `json:"categories,omitempty"`
	// Conditions describes the product conditions, e.g. refurbished or used.
	Conditions []*Condition `json:"conditions,omitempty"`
	// ContentUnit is the content unit of the product, a 3-character ISO code
	// (usually project-specific).
	ContentUnit string `json:"cu,omitempty"`
	// CuPerOu describes the number of content units per order unit, e.g. the
	// 12 in '1 case contains 12 bottles'.
	CuPerOu *float64 `json:"cuPerOu,omitempty"`
	// CustField1 is the CUST_FIELD1 of the SAP OCI specification. It has a
	// maximum length of 10 characters.
	CustField1 string `json:"custField1,omitempty"`
	// CustField2 is the CUST_FIELD2 of the SAP OCI specification. It has a
	// maximum length of 10 characters.
	CustField2 string `json:"custField2,omitempty"`
	// CustField3 is the CUST_FIELD3 of the SAP OCI specification. It has a
	// maximum length of 10 characters.
	CustField3 string `json:"custField3,omitempty"`
	// CustField4 is the CUST_FIELD4 of the SAP OCI specification. It has a
	// maximum length of 20 characters.
	CustField4 string `json:"custField4,omitempty"`
	// CustField5 is the CUST_FIELD5 of the SAP OCI specification. It has a
	// maximum length of 50 characters.
	CustField5 string `json:"custField5,omitempty"`
	// CustFields is an array of generic name/value pairs for
	// customer-specific attributes.
	CustFields []*CustField `json:"custFields,omitempty"`
	// Datasheet is the name of an datasheet file (in the media files) or a
	// URL to the datasheet on the internet.
	Datasheet string `json:"datasheet,omitempty"`
	// Description of the product.
	Description string `json:"description,omitempty"`
	// Eclasses is a list of eCl@ss categories the product belongs to.
	Eclasses []*Eclass `json:"eclasses,omitempty"`
	// ErpGroupSupplier: erpGroupSupplier is the material group of the product
	// on the merchant-/supplier-side.
	ErpGroupSupplier string `json:"erpGroupSupplier,omitempty"`
	// Excluded is a flag that indicates whether to exclude this product from
	// the catalog. If true, this product will not be published into the live
	// area.
	Excluded bool `json:"excluded,omitempty"`
	// ExtCategory is the EXT_CATEGORY field of the SAP OCI specification.
	ExtCategory string `json:"extCategory,omitempty"`
	// ExtCategoryID is the EXT_CATEGORY_ID field of the SAP OCI
	// specification.
	ExtCategoryID string `json:"extCategoryId,omitempty"`
	// ExtSchemaType is the EXT_SCHEMA_TYPE field of the SAP OCI
	// specification.
	ExtSchemaType string `json:"extSchemaType,omitempty"`
	// Features defines product features, i.e. additional properties of the
	// product.
	Features []*Feature `json:"features,omitempty"`
	// Gtin: GTIN is the global trade item number of the product (used to be
	// EAN).
	Gtin string `json:"gtin,omitempty"`
	// Hazmats classifies hazardous/dangerous goods.
	Hazmats []*Hazmat `json:"hazmats,omitempty"`
	// Image is the name of an image file (in the media files) or a URL to the
	// image on the internet.
	Image string `json:"image,omitempty"`
	// Keywords is a list of aliases for the product.
	Keywords []string `json:"keywords,omitempty"`
	// Leadtime is the number of days for delivery.
	Leadtime *float64 `json:"leadtime,omitempty"`
	// ListPrice is the net list price of the product.
	ListPrice *float64 `json:"listPrice,omitempty"`
	// Manufactcode is the manufacturer code as used in the SAP OCI
	// specification.
	Manufactcode string `json:"manufactcode,omitempty"`
	// Manufacturer is the name of the manufacturer.
	Manufacturer string `json:"manufacturer,omitempty"`
	// Matgroup is the material group of the product on the buy-side.
	Matgroup string `json:"matgroup,omitempty"`
	// Mpn: MPN is the manufacturer part number.
	Mpn string `json:"mpn,omitempty"`
	// Name of the product.
	Name string `json:"name,omitempty"`
	// OrderUnit is the order unit of the product, a 3-character ISO code
	// (usually project-specific).
	OrderUnit string `json:"ou,omitempty"`
	// Price is the net price (per order unit) of the product for the
	// end-user.
	Price float64 `json:"price,omitempty"`
	// PriceQty is the quantity for which the price is specified (default:
	// 1.0).
	PriceQty float64 `json:"priceQty,omitempty"`
	// QuantityInterval is the interval in which this product can be ordered.
	// E.g. if the quantity interval is 5, the end-user can only order in
	// quantities of 5,10,15 etc.
	QuantityInterval *float64 `json:"quantityInterval,omitempty"`
	// QuantityMax is the maximum order quantity for this product.
	QuantityMax *float64 `json:"quantityMax,omitempty"`
	// QuantityMin is the minimum order quantity for this product.
	QuantityMin *float64 `json:"quantityMin,omitempty"`
	// References defines cross-product references, e.g. alternatives or
	// follow-up products.
	References []*Reference `json:"references,omitempty"`
	// Safetysheet is the name of an safetysheet file (in the media files) or
	// a URL to the safetysheet on the internet.
	Safetysheet string `json:"safetysheet,omitempty"`
	// ScalePrices can be used when the price of the product is dependent on
	// the ordered quantity.
	ScalePrices []*ScalePrice `json:"scalePrices,omitempty"`
	// Service indicates if the is a good (false) or a service (true). The
	// default value is false.
	Service bool `json:"service,omitempty"`
	// TaxCode to use for this product. This is typically project-specific.
	TaxCode string `json:"taxCode,omitempty"`
	// Thumbnail is the name of an thumbnail image file (in the media files)
	// or a URL to the image on the internet.
	Thumbnail string `json:"thumbnail,omitempty"`
	// Unspscs is a list of UNSPSC categories the product belongs to.
	Unspscs []*Unspsc `json:"unspscs,omitempty"`
}

// ReplaceProductResponse is the outcome of a successful replacement of a
// product.
type ReplaceProductResponse struct {
	// Kind describes this entity.
	Kind string `json:"kind,omitempty"`
	// Link returns a URL to the representation of the replaced product.
	Link string `json:"link,omitempty"`
}

// ScalePrice describes a price that is dependent on the ordered quantity.
type ScalePrice struct {
	// Lbound: LBound is the lower bound when this price will become
	// effective.
	Lbound float64 `json:"lbound,omitempty"`
	// Price is the net price for the given lower bound.
	Price float64 `json:"price,omitempty"`
}

// ScrollResponse is a slice of products from a catalog.
type ScrollResponse struct {
	// Items is the slice of products of this result.
	Items []*Product `json:"items,omitempty"`
	// Kind is store#products/scroll for this kind of response.
	Kind string `json:"kind,omitempty"`
	// NextLink returns the URL to the next slice of products (if any).
	NextLink string `json:"nextLink,omitempty"`
	// PageToken needs to be passed to get the next slice of products. It is
	// blank if there are no more products. Instead of using pageToken for
	// your next request, you can also follow nextLink.
	PageToken string `json:"pageToken,omitempty"`
	// PreviousLink returns the URL of the previous slice of products (if
	// any).
	PreviousLink string `json:"previousLink,omitempty"`
	// SelfLink returns the URL to this page.
	SelfLink string `json:"selfLink,omitempty"`
	// TotalItems describes the total number of products found.
	TotalItems int64 `json:"totalItems,omitempty"`
}

// SearchResponse is a partial listing of products.
type SearchResponse struct {
	// Items is the slice of products of this result.
	Items []*Product `json:"items,omitempty"`
	// Kind is store#products/search for this kind of response.
	Kind string `json:"kind,omitempty"`
	// NextLink returns the URL to the next slice of products (if any).
	NextLink string `json:"nextLink,omitempty"`
	// PreviousLink returns the URL of the previous slice of products (if
	// any).
	PreviousLink string `json:"previousLink,omitempty"`
	// SelfLink returns the URL to this page.
	SelfLink string `json:"selfLink,omitempty"`
	// TotalItems describes the total number of products found.
	TotalItems int64 `json:"totalItems,omitempty"`
}

// Unspsc is used to tie a product to a UNSPSC schema.
type Unspsc struct {
	// Code is the UNSPSC code. Only use digits for encoding, e.g. 43211503.
	Code string `json:"code,omitempty"`
	// Version is the UNSPSC version in the major.minor format, e.g. 16.0901.
	Version string `json:"version,omitempty"`
}

// UpdateProduct holds the properties of a product that need to be
// updated.
type UpdateProduct struct {
	// Availability allows the update of product availability data, e.g. the
	// number of items in stock or the date when the product will be available
	// again.
	Availability *Availability `json:"availability,omitempty"`
	// Blobs specifies external data, e.g. images or datasheets.
	Blobs []*Blob `json:"blobs,omitempty"`
	// Bpn: BPN is the buyer part number of the product.
	Bpn *string `json:"bpn,omitempty"`
	// Categories is a list of (supplier-specific) category names the product
	// belongs to.
	Categories []string `json:"categories,omitempty"`
	// Conditions describes the product conditions, e.g. refurbished or used.
	Conditions []*Condition `json:"conditions,omitempty"`
	// ContentUnit is the content unit of the product, a 3-character ISO code
	// (usually project-specific).
	ContentUnit *string `json:"cu,omitempty"`
	// CuPerOu describes the number of content units per order unit, e.g. the
	// 12 in '1 case contains 12 bottles'.
	CuPerOu *float64 `json:"cuPerOu,omitempty"`
	// CustField1 is the CUST_FIELD1 of the SAP OCI specification. It has a
	// maximum length of 10 characters.
	CustField1 *string `json:"custField1,omitempty"`
	// CustField2 is the CUST_FIELD2 of the SAP OCI specification. It has a
	// maximum length of 10 characters.
	CustField2 *string `json:"custField2,omitempty"`
	// CustField3 is the CUST_FIELD3 of the SAP OCI specification. It has a
	// maximum length of 10 characters.
	CustField3 *string `json:"custField3,omitempty"`
	// CustField4 is the CUST_FIELD4 of the SAP OCI specification. It has a
	// maximum length of 20 characters.
	CustField4 *string `json:"custField4,omitempty"`
	// CustField5 is the CUST_FIELD5 of the SAP OCI specification. It has a
	// maximum length of 50 characters.
	CustField5 *string `json:"custField5,omitempty"`
	// CustFields is an array of generic name/value pairs for
	// customer-specific attributes.
	CustFields []*CustField `json:"custFields,omitempty"`
	// Datasheet is the name of an datasheet file (in the media files) or a
	// URL to the datasheet on the internet.
	Datasheet *string `json:"datasheet,omitempty"`
	// Description of the product.
	Description *string `json:"description,omitempty"`
	// Eclasses is a list of eCl@ss categories the product belongs to.
	Eclasses []*Eclass `json:"eclasses,omitempty"`
	// ErpGroupSupplier: erpGroupSupplier is the material group of the product
	// on the merchant-/supplier-side.
	ErpGroupSupplier *string `json:"erpGroupSupplier,omitempty"`
	// Excluded is a flag that indicates whether to exclude this product from
	// the catalog. If true, this product will not be published into the live
	// area.
	Excluded *bool `json:"excluded,omitempty"`
	// ExtCategory is the EXT_CATEGORY field of the SAP OCI specification.
	ExtCategory *string `json:"extCategory,omitempty"`
	// ExtCategoryID is the EXT_CATEGORY_ID field of the SAP OCI
	// specification.
	ExtCategoryID *string `json:"extCategoryId,omitempty"`
	// ExtSchemaType is the EXT_SCHEMA_TYPE field of the SAP OCI
	// specification.
	ExtSchemaType *string `json:"extSchemaType,omitempty"`
	// Features defines product features, i.e. additional properties of the
	// product.
	Features []*Feature `json:"features,omitempty"`
	// Gtin: GTIN is the global trade item number of the product (used to be
	// EAN).
	Gtin *string `json:"gtin,omitempty"`
	// Hazmats classifies hazardous/dangerous goods.
	Hazmats []*Hazmat `json:"hazmats,omitempty"`
	// Image is the name of an image file (in the media files) or a URL to the
	// image on the internet.
	Image *string `json:"image,omitempty"`
	// Keywords is a list of aliases for the product.
	Keywords []string `json:"keywords,omitempty"`
	// Leadtime is the number of days for delivery.
	Leadtime *float64 `json:"leadtime,omitempty"`
	// ListPrice is the net list price of the product.
	ListPrice *float64 `json:"listPrice,omitempty"`
	// Manufactcode is the manufacturer code as used in the SAP OCI
	// specification.
	Manufactcode *string `json:"manufactcode,omitempty"`
	// Manufacturer is the name of the manufacturer.
	Manufacturer *string `json:"manufacturer,omitempty"`
	// Matgroup is the material group of the product on the buy-side.
	Matgroup *string `json:"matgroup,omitempty"`
	// Mpn: MPN is the manufacturer part number.
	Mpn *string `json:"mpn,omitempty"`
	// Name of the product.
	Name *string `json:"name,omitempty"`
	// OrderUnit is the order unit of the product, a 3-character ISO code
	// (usually project-specific).
	OrderUnit *string `json:"ou,omitempty"`
	// Price is the net price (per order unit) of the product for the
	// end-user.
	Price *float64 `json:"price,omitempty"`
	// PriceQty is the quantity for which the price is specified (default:
	// 1.0).
	PriceQty *float64 `json:"priceQty,omitempty"`
	// QuantityInterval is the interval in which this product can be ordered.
	// E.g. if the quantity interval is 5, the end-user can only order in
	// quantities of 5,10,15 etc.
	QuantityInterval *float64 `json:"quantityInterval,omitempty"`
	// QuantityMax is the maximum order quantity for this product.
	QuantityMax *float64 `json:"quantityMax,omitempty"`
	// QuantityMin is the minimum order quantity for this product.
	QuantityMin *float64 `json:"quantityMin,omitempty"`
	// References defines cross-product references, e.g. alternatives or
	// follow-up products.
	References []*Reference `json:"references,omitempty"`
	// Safetysheet is the name of an safetysheet file (in the media files) or
	// a URL to the safetysheet on the internet.
	Safetysheet *string `json:"safetysheet,omitempty"`
	// ScalePrices can be used when the price of the product is dependent on
	// the ordered quantity.
	ScalePrices []*ScalePrice `json:"scalePrices,omitempty"`
	// Service indicates if the is a good (false) or a service (true). The
	// default value is false.
	Service *bool `json:"service,omitempty"`
	// TaxCode to use for this product. This is typically project-specific.
	TaxCode *string `json:"taxCode,omitempty"`
	// Thumbnail is the name of an thumbnail image file (in the media files)
	// or a URL to the image on the internet.
	Thumbnail *string `json:"thumbnail,omitempty"`
	// Unspscs is a list of UNSPSC categories the product belongs to.
	Unspscs []*Unspsc `json:"unspscs,omitempty"`
}

// UpdateProductResponse is the outcome of a successful request to update
// a product.
type UpdateProductResponse struct {
	// Kind describes this entity.
	Kind string `json:"kind,omitempty"`
	// Link returns a URL to the representation of the updated product.
	Link string `json:"link,omitempty"`
}

// UpsertProduct holds the properties of the product to create or update.
type UpsertProduct struct {
	// Availability allows the update of product availability data, e.g. the
	// number of items in stock or the date when the product will be available
	// again.
	Availability *Availability `json:"availability,omitempty"`
	// Blobs specifies external data, e.g. images or datasheets.
	Blobs []*Blob `json:"blobs,omitempty"`
	// Bpn: BPN is the buyer part number of the product.
	Bpn string `json:"bpn,omitempty"`
	// Categories is a list of (supplier-specific) category names the product
	// belongs to.
	Categories []string `json:"categories,omitempty"`
	// Conditions describes the product conditions, e.g. refurbished or used.
	Conditions []*Condition `json:"conditions,omitempty"`
	// ContentUnit is the content unit of the product, a 3-character ISO code
	// (usually project-specific).
	ContentUnit string `json:"cu,omitempty"`
	// CuPerOu describes the number of content units per order unit, e.g. the
	// 12 in '1 case contains 12 bottles'.
	CuPerOu *float64 `json:"cuPerOu,omitempty"`
	// CustField1 is the CUST_FIELD1 of the SAP OCI specification. It has a
	// maximum length of 10 characters.
	CustField1 string `json:"custField1,omitempty"`
	// CustField2 is the CUST_FIELD2 of the SAP OCI specification. It has a
	// maximum length of 10 characters.
	CustField2 string `json:"custField2,omitempty"`
	// CustField3 is the CUST_FIELD3 of the SAP OCI specification. It has a
	// maximum length of 10 characters.
	CustField3 string `json:"custField3,omitempty"`
	// CustField4 is the CUST_FIELD4 of the SAP OCI specification. It has a
	// maximum length of 20 characters.
	CustField4 string `json:"custField4,omitempty"`
	// CustField5 is the CUST_FIELD5 of the SAP OCI specification. It has a
	// maximum length of 50 characters.
	CustField5 string `json:"custField5,omitempty"`
	// CustFields is an array of generic name/value pairs for
	// customer-specific attributes.
	CustFields []*CustField `json:"custFields,omitempty"`
	// Datasheet is the name of an datasheet file (in the media files) or a
	// URL to the datasheet on the internet.
	Datasheet string `json:"datasheet,omitempty"`
	// Description of the product.
	Description string `json:"description,omitempty"`
	// Eclasses is a list of eCl@ss categories the product belongs to.
	Eclasses []*Eclass `json:"eclasses,omitempty"`
	// ErpGroupSupplier: erpGroupSupplier is the material group of the product
	// on the merchant-/supplier-side.
	ErpGroupSupplier string `json:"erpGroupSupplier,omitempty"`
	// Excluded is a flag that indicates whether to exclude this product from
	// the catalog. If true, this product will not be published into the live
	// area.
	Excluded bool `json:"excluded,omitempty"`
	// ExtCategory is the EXT_CATEGORY field of the SAP OCI specification.
	ExtCategory string `json:"extCategory,omitempty"`
	// ExtCategoryID is the EXT_CATEGORY_ID field of the SAP OCI
	// specification.
	ExtCategoryID string `json:"extCategoryId,omitempty"`
	// ExtSchemaType is the EXT_SCHEMA_TYPE field of the SAP OCI
	// specification.
	ExtSchemaType string `json:"extSchemaType,omitempty"`
	// Features defines product features, i.e. additional properties of the
	// product.
	Features []*Feature `json:"features,omitempty"`
	// Gtin: GTIN is the global trade item number of the product (used to be
	// EAN).
	Gtin string `json:"gtin,omitempty"`
	// Hazmats classifies hazardous/dangerous goods.
	Hazmats []*Hazmat `json:"hazmats,omitempty"`
	// Image is the name of an image file (in the media files) or a URL to the
	// image on the internet.
	Image string `json:"image,omitempty"`
	// Keywords is a list of aliases for the product.
	Keywords []string `json:"keywords,omitempty"`
	// Leadtime is the number of days for delivery.
	Leadtime *float64 `json:"leadtime,omitempty"`
	// ListPrice is the net list price of the product.
	ListPrice *float64 `json:"listPrice,omitempty"`
	// Manufactcode is the manufacturer code as used in the SAP OCI
	// specification.
	Manufactcode string `json:"manufactcode,omitempty"`
	// Manufacturer is the name of the manufacturer.
	Manufacturer string `json:"manufacturer,omitempty"`
	// Matgroup is the material group of the product on the buy-side.
	Matgroup string `json:"matgroup,omitempty"`
	// Mpn: MPN is the manufacturer part number.
	Mpn string `json:"mpn,omitempty"`
	// Name of the product. The product name is a required field
	Name string `json:"name,omitempty"`
	// OrderUnit is the order unit of the product, a 3-character ISO code
	// (usually project-specific). OrderUnit is a required field.
	OrderUnit string `json:"ou,omitempty"`
	// Price is the net price (per order unit) of the product for the
	// end-user. Price is a required field.
	Price float64 `json:"price,omitempty"`
	// PriceQty is the quantity for which the price is specified (default:
	// 1.0).
	PriceQty *float64 `json:"priceQty,omitempty"`
	// QuantityInterval is the interval in which this product can be ordered.
	// E.g. if the quantity interval is 5, the end-user can only order in
	// quantities of 5,10,15 etc.
	QuantityInterval *float64 `json:"quantityInterval,omitempty"`
	// QuantityMax is the maximum order quantity for this product.
	QuantityMax *float64 `json:"quantityMax,omitempty"`
	// QuantityMin is the minimum order quantity for this product.
	QuantityMin *float64 `json:"quantityMin,omitempty"`
	// References defines cross-product references, e.g. alternatives or
	// follow-up products.
	References []*Reference `json:"references,omitempty"`
	// Safetysheet is the name of an safetysheet file (in the media files) or
	// a URL to the safetysheet on the internet.
	Safetysheet string `json:"safetysheet,omitempty"`
	// ScalePrices can be used when the price of the product is dependent on
	// the ordered quantity.
	ScalePrices []*ScalePrice `json:"scalePrices,omitempty"`
	// Service indicates if the is a good (false) or a service (true). The
	// default value is false.
	Service bool `json:"service,omitempty"`
	// Spn: SPN is the supplier part number. SPN is a required field.
	Spn string `json:"spn,omitempty"`
	// TaxCode to use for this product. This is typically project-specific.
	TaxCode string `json:"taxCode,omitempty"`
	// Thumbnail is the name of an thumbnail image file (in the media files)
	// or a URL to the image on the internet.
	Thumbnail string `json:"thumbnail,omitempty"`
	// Unspscs is a list of UNSPSC categories the product belongs to.
	Unspscs []*Unspsc `json:"unspscs,omitempty"`
}

// UpsertProductResponse is the outcome of a successful request to upsert
// a product.
type UpsertProductResponse struct {
	// Kind describes this entity.
	Kind string `json:"kind,omitempty"`
	// Link returns a URL to the representation of the created or updated
	// product.
	Link string `json:"link,omitempty"`
}

// Create a new product in the given catalog and area.
type CreateService struct {
	s       *Service
	opt_    map[string]interface{}
	hdr_    map[string]interface{}
	pin     string
	area    string
	product *CreateProduct
}

// NewCreateService creates a new instance of CreateService.
func NewCreateService(s *Service) *CreateService {
	rs := &CreateService{s: s, opt_: make(map[string]interface{}), hdr_: make(map[string]interface{})}
	return rs
}

// Area of the catalog, e.g. work or live.
func (s *CreateService) Area(area string) *CreateService {
	s.area = area
	return s
}

// PIN of the catalog.
func (s *CreateService) PIN(pin string) *CreateService {
	s.pin = pin
	return s
}

// Product properties of the new product.
func (s *CreateService) Product(product *CreateProduct) *CreateService {
	s.product = product
	return s
}

// Do executes the operation.
func (s *CreateService) Do() (*CreateProductResponse, error) {
	var body io.Reader
	body, err := meplatoapi.ReadJSON(s.product)
	if err != nil {
		return nil, err
	}
	params := make(map[string]interface{})
	params["area"] = s.area
	params["pin"] = s.pin
	path, err := meplatoapi.Expand("/catalogs/{pin}/{area}/products", params)
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
	ret := new(CreateProductResponse)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
		return nil, err
	}
	return ret, nil
}

// Delete a product.
type DeleteService struct {
	s    *Service
	opt_ map[string]interface{}
	hdr_ map[string]interface{}
	pin  string
	area string
	spn  string
}

// NewDeleteService creates a new instance of DeleteService.
func NewDeleteService(s *Service) *DeleteService {
	rs := &DeleteService{s: s, opt_: make(map[string]interface{}), hdr_: make(map[string]interface{})}
	return rs
}

// Area of the catalog, e.g. work or live.
func (s *DeleteService) Area(area string) *DeleteService {
	s.area = area
	return s
}

// PIN of the catalog.
func (s *DeleteService) PIN(pin string) *DeleteService {
	s.pin = pin
	return s
}

// SPN is the supplier part number of the product to delete.
func (s *DeleteService) Spn(spn string) *DeleteService {
	s.spn = spn
	return s
}

// Do executes the operation.
func (s *DeleteService) Do() error {
	var body io.Reader
	params := make(map[string]interface{})
	params["area"] = s.area
	params["pin"] = s.pin
	params["spn"] = s.spn
	path, err := meplatoapi.Expand("/catalogs/{pin}/{area}/products/{spn}", params)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("DELETE", s.s.BaseURL+path, body)
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

// Get returns a single product by its Supplier Part Number (SPN).
type GetService struct {
	s    *Service
	opt_ map[string]interface{}
	hdr_ map[string]interface{}
	pin  string
	area string
	spn  string
}

// NewGetService creates a new instance of GetService.
func NewGetService(s *Service) *GetService {
	rs := &GetService{s: s, opt_: make(map[string]interface{}), hdr_: make(map[string]interface{})}
	return rs
}

// Area of the catalog, e.g. work or live.
func (s *GetService) Area(area string) *GetService {
	s.area = area
	return s
}

// PIN of the catalog.
func (s *GetService) PIN(pin string) *GetService {
	s.pin = pin
	return s
}

// SPN is the supplier part number of the product to get.
func (s *GetService) Spn(spn string) *GetService {
	s.spn = spn
	return s
}

// Do executes the operation.
func (s *GetService) Do() (*Product, error) {
	var body io.Reader
	params := make(map[string]interface{})
	params["area"] = s.area
	params["pin"] = s.pin
	params["spn"] = s.spn
	path, err := meplatoapi.Expand("/catalogs/{pin}/{area}/products/{spn}", params)
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
	ret := new(Product)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
		return nil, err
	}
	return ret, nil
}

// Replace all fields of a product. Use Update to update only certain
// fields of a product.
type ReplaceService struct {
	s       *Service
	opt_    map[string]interface{}
	hdr_    map[string]interface{}
	pin     string
	area    string
	spn     string
	product *ReplaceProduct
}

// NewReplaceService creates a new instance of ReplaceService.
func NewReplaceService(s *Service) *ReplaceService {
	rs := &ReplaceService{s: s, opt_: make(map[string]interface{}), hdr_: make(map[string]interface{})}
	return rs
}

// Area of the catalog, e.g. work or live.
func (s *ReplaceService) Area(area string) *ReplaceService {
	s.area = area
	return s
}

// PIN of the catalog.
func (s *ReplaceService) PIN(pin string) *ReplaceService {
	s.pin = pin
	return s
}

// New properties of the product.
func (s *ReplaceService) Product(product *ReplaceProduct) *ReplaceService {
	s.product = product
	return s
}

// SPN is the supplier part number of the product to replace.
func (s *ReplaceService) Spn(spn string) *ReplaceService {
	s.spn = spn
	return s
}

// Do executes the operation.
func (s *ReplaceService) Do() (*ReplaceProductResponse, error) {
	var body io.Reader
	body, err := meplatoapi.ReadJSON(s.product)
	if err != nil {
		return nil, err
	}
	params := make(map[string]interface{})
	params["area"] = s.area
	params["pin"] = s.pin
	params["spn"] = s.spn
	path, err := meplatoapi.Expand("/catalogs/{pin}/{area}/products/{spn}", params)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("PUT", s.s.BaseURL+path, body)
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
	ret := new(ReplaceProductResponse)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
		return nil, err
	}
	return ret, nil
}

// Scroll through products of a catalog (area). If you need to iterate
// through all products in a catalog, this is the most effective way to do
// so. If you want to search for products, use the Search endpoint.
type ScrollService struct {
	s    *Service
	opt_ map[string]interface{}
	hdr_ map[string]interface{}
	pin  string
	area string
}

// NewScrollService creates a new instance of ScrollService.
func NewScrollService(s *Service) *ScrollService {
	rs := &ScrollService{s: s, opt_: make(map[string]interface{}), hdr_: make(map[string]interface{})}
	return rs
}

// Area of the catalog, e.g. work or live.
func (s *ScrollService) Area(area string) *ScrollService {
	s.area = area
	return s
}

// PageToken must be passed in the 2nd and all consective requests to get
// the next page of results. You do not need to pass the page token
// manually. You should just follow the nextUrl link in the metadata to
// get the next slice of data. If there is no nextUrl returned, you have
// reached the last page of products for the given catalog. A scroll
// request is kept alive for 2 minutes. If you fail to ask for the next
// slice of products within this period, a new scroll request will be
// created and you start over at the first page of results.
func (s *ScrollService) PageToken(pageToken string) *ScrollService {
	s.opt_["pageToken"] = pageToken
	return s
}

// PIN of the catalog.
func (s *ScrollService) PIN(pin string) *ScrollService {
	s.pin = pin
	return s
}

// Do executes the operation.
func (s *ScrollService) Do() (*ScrollResponse, error) {
	var body io.Reader
	params := make(map[string]interface{})
	params["area"] = s.area
	if v, ok := s.opt_["pageToken"]; ok {
		params["pageToken"] = v
	}
	params["pin"] = s.pin
	path, err := meplatoapi.Expand("/catalogs/{pin}/{area}/products/scroll{?pageToken}", params)
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
	ret := new(ScrollResponse)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
		return nil, err
	}
	return ret, nil
}

// Search for products. Do not use this method for iterating through all
// of the products in a catalog; use the Scroll endpoint instead. It is
// much more efficient.
type SearchService struct {
	s    *Service
	opt_ map[string]interface{}
	hdr_ map[string]interface{}
	pin  string
	area string
}

// NewSearchService creates a new instance of SearchService.
func NewSearchService(s *Service) *SearchService {
	rs := &SearchService{s: s, opt_: make(map[string]interface{}), hdr_: make(map[string]interface{})}
	return rs
}

// Area of the catalog, e.g. work or live.
func (s *SearchService) Area(area string) *SearchService {
	s.area = area
	return s
}

// PIN of the catalog.
func (s *SearchService) PIN(pin string) *SearchService {
	s.pin = pin
	return s
}

// Q defines are full text query.
func (s *SearchService) Q(q string) *SearchService {
	s.opt_["q"] = q
	return s
}

// Skip specifies how many products to skip (default 0).
func (s *SearchService) Skip(skip int64) *SearchService {
	s.opt_["skip"] = skip
	return s
}

// Take defines how many products to return (max 100, default 20).
func (s *SearchService) Take(take int64) *SearchService {
	s.opt_["take"] = take
	return s
}

// Do executes the operation.
func (s *SearchService) Do() (*SearchResponse, error) {
	var body io.Reader
	params := make(map[string]interface{})
	params["area"] = s.area
	params["pin"] = s.pin
	if v, ok := s.opt_["q"]; ok {
		params["q"] = v
	}
	if v, ok := s.opt_["skip"]; ok {
		params["skip"] = v
	}
	if v, ok := s.opt_["take"]; ok {
		params["take"] = v
	}
	path, err := meplatoapi.Expand("/catalogs/{pin}/{area}/products{?q,skip,take}", params)
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

// Update the fields of a product selectively. Use Replace to replace the
// product as a whole.
type UpdateService struct {
	s       *Service
	opt_    map[string]interface{}
	hdr_    map[string]interface{}
	pin     string
	area    string
	spn     string
	product *UpdateProduct
}

// NewUpdateService creates a new instance of UpdateService.
func NewUpdateService(s *Service) *UpdateService {
	rs := &UpdateService{s: s, opt_: make(map[string]interface{}), hdr_: make(map[string]interface{})}
	return rs
}

// Area of the catalog, e.g. work or live.
func (s *UpdateService) Area(area string) *UpdateService {
	s.area = area
	return s
}

// PIN of the catalog.
func (s *UpdateService) PIN(pin string) *UpdateService {
	s.pin = pin
	return s
}

// Products properties of the updated product.
func (s *UpdateService) Product(product *UpdateProduct) *UpdateService {
	s.product = product
	return s
}

// SPN is the supplier part number of the product to update.
func (s *UpdateService) Spn(spn string) *UpdateService {
	s.spn = spn
	return s
}

// Do executes the operation.
func (s *UpdateService) Do() (*UpdateProductResponse, error) {
	var body io.Reader
	body, err := meplatoapi.ReadJSON(s.product)
	if err != nil {
		return nil, err
	}
	params := make(map[string]interface{})
	params["area"] = s.area
	params["pin"] = s.pin
	params["spn"] = s.spn
	path, err := meplatoapi.Expand("/catalogs/{pin}/{area}/products/{spn}", params)
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
	ret := new(UpdateProductResponse)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
		return nil, err
	}
	return ret, nil
}

// Upsert a product in the given catalog and area. Upsert will create if
// the product does not exist yet, otherwise it will update.
type UpsertService struct {
	s       *Service
	opt_    map[string]interface{}
	hdr_    map[string]interface{}
	pin     string
	area    string
	product *UpsertProduct
}

// NewUpsertService creates a new instance of UpsertService.
func NewUpsertService(s *Service) *UpsertService {
	rs := &UpsertService{s: s, opt_: make(map[string]interface{}), hdr_: make(map[string]interface{})}
	return rs
}

// Area of the catalog, e.g. work or live.
func (s *UpsertService) Area(area string) *UpsertService {
	s.area = area
	return s
}

// PIN of the catalog.
func (s *UpsertService) PIN(pin string) *UpsertService {
	s.pin = pin
	return s
}

// Product properties of the new product.
func (s *UpsertService) Product(product *UpsertProduct) *UpsertService {
	s.product = product
	return s
}

// Do executes the operation.
func (s *UpsertService) Do() (*UpsertProductResponse, error) {
	var body io.Reader
	body, err := meplatoapi.ReadJSON(s.product)
	if err != nil {
		return nil, err
	}
	params := make(map[string]interface{})
	params["area"] = s.area
	params["pin"] = s.pin
	path, err := meplatoapi.Expand("/catalogs/{pin}/{area}/products/upsert", params)
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
	ret := new(UpsertProductResponse)
	if err := json.NewDecoder(res.Body).Decode(ret); err != nil {
		return nil, err
	}
	return ret, nil
}

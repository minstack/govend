// Package vend handles interactions with the Vend API.
package vend

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

// Product is a basic product object
type Product struct {
	ID                      *string          `json:"id"`
	SourceID                *string          `json:"source_id"`
	VariantSourceID         *string          `json:"variant_source_id"`
	Handle                  *string          `json:"handle"`
	HasVariants             bool             `json:"has_variants"`
	VariantParentID         *string          `json:"variant_parent_id"`
	VariantOptionOneName    *string          `json:"variant_option_one_name"`
	VariantOptionOneValue   *string          `json:"variant_option_one_value"`
	VariantOptionTwoName    *string          `json:"variant_option_two_name"`
	VariantOptionTwoValue   *string          `json:"variant_option_two_value"`
	VariantOptionThreeName  *string          `json:"variant_option_three_name"`
	VariantOptionThreeValue *string          `json:"variant_option_three_value"`
	VariantName             *string          `json:"variant_name,omitempty"`
	Active                  bool             `json:"active"`
	Name                    *string          `json:"name"`
	Description             *string          `json:"description"`
	Image                   *string          `json:"image"`
	ImageURL                *string          `json:"image_url"`
	ImageLarge              *string          `json:"image_large"`
	Images                  []Image          `json:"images"`
	SKU                     *string          `json:"sku"`
	Tags                    *string          `json:"tags"`
	BrandID                 *string          `json:"brand_id"`
	BrandName               *string          `json:"brand_name"`
	SupplierName            *string          `json:"supplier_name"`
	SupplierCode            *string          `json:"supplier_code"`
	SupplyPrice             *float64         `json:"supply_price"`
	AccountCodePurchase     *string          `json:"account_code_purchase"`
	AccountCodeSales        *string          `json:"account_code_sales"`
	TrackInventory          bool             `json:"track_inventory"`
	Inventory               []Inventory      `json:"inventory"`
	PriceBookEntries        []PriceBookEntry `json:"price_book_entries"`
	Price                   *float64         `json:"price"`
	Tax                     *float64         `json:"tax"`
	TaxID                   *string          `json:"tax_id"`
	TaxRate                 *float64         `json:"tax_rate"`
	TaxName                 *string          `json:"tax_name"`
	Taxes                   []Tax            `json:"taxes"`
	UpdatedAt               *string          `json:"updated_at"`
	DeletedAt               *string          `json:"deleted_at"`
}

// Inventory houses product inventory object
type Inventory struct {
	OutletID     string `json:"outlet_id"`
	OutletName   string `json:"outlet_name"`
	Count        string `json:"count"`
	ReorderPoint string `json:"reorder_point"`
	RestockLevel string `json:"restock_level"`
}

// PriceBookEntry houses product pricing object
type PriceBookEntry struct {
	ID                             string  `json:"id"`
	ProductID                      string  `json:"product_id"`
	PriceBookID                    string  `json:"price_book_id"`
	PriceBookName                  string  `json:"price_book_name"`
	Type                           string  `json:"type"`
	OutletName                     string  `json:"outlet_name"`
	OutletID                       string  `json:"outlet_id"`
	CustomerGroupName              string  `json:"customer_group_name"`
	CustomerGroupID                string  `json:"customer_group_id"`
	Price                          float64 `json:"price"`
	LoyaltyValue                   int64   `json:"loyalty_value"`
	Tax                            float64 `json:"tax"`
	TaxID                          string  `json:"tax_id"`
	TaxRate                        float64 `json:"tax_rate"`
	TaxName                        string  `json:"tax_name"`
	DisplayRetailPriceTaxInclusive int64   `json:"display_retail_price_tax_inclusive"`
	MinUnits                       string  `json:"min_units"`
	MaxUnits                       string  `json:"max_units"`
	ValidFrom                      string  `json:"valid_from"`
	ValidTo                        string  `json:"valid_to"`
}

// Tax houses product tax object
type Tax struct {
	OutletID string `json:"outlet_id"`
	TaxID    string `json:"tax_id"`
}

// Image is the info contained in each Vend image object.
type Image struct {
	ID      *string `json:"id,omitempty"`
	URL     *string `json:"url,omitempty"`
	Version *int64  `json:"version"`
}

// ImageUpload holds data for Images
type ImageUpload struct {
	Data Data `json:"data,omitempty"`
}

// Data is the information for each image contained in the response.
type Data struct {
	ID        *string `json:"id,omitempty"`
	ProductID *string `json:"product_id,omitempty"`
	Position  *int64  `json:"position,omitempty"`
	Status    *string `json:"status,omitempty"`
	Version   *int64  `json:"version,omitempty"`
}

// ProductUpload contains the fields needed to post an image to a product in Vend.
type ProductUpload struct {
	ID       string `json:"id,omitempty"`
	Handle   string `json:"handle,omitempty"`
	SKU      string `json:"sku,omitempty"`
	ImageURL string `json:"image_url,omitempty"`
}

// Products grabs and collates all products in pages of 10,000.
func (c Client) Products() ([]Product, map[string]Product, error) {

	productMap := make(map[string]Product)
	var products, page []Product
	var data []byte
	var v int64

	// v is a version that is used to get products by page.
	// Here we get the first page.
	data, v, err := ResourcePage(0, c.DomainPrefix, c.Token, "products")

	// Unmarshal payload into sales object.
	err = json.Unmarshal(data, &page)

	products = append(products, page...)

	for len(page) > 0 {
		page = []Product{}

		// Continue grabbing pages until we receive an empty one.
		data, v, err = ResourcePage(v, c.DomainPrefix, c.Token, "products")
		if err != nil {
			return nil, nil, err
		}

		// Unmarshal payload into product object.
		err = json.Unmarshal(data, &page)

		// Append page to list.
		products = append(products, page...)
	}

	productMap = buildProductMap(products)

	return products, productMap, err
}

func buildProductMap(products []Product) map[string]Product {
	productMap := make(map[string]Product)

	for _, product := range products {
		productMap[*product.ID] = product
	}

	return productMap
}

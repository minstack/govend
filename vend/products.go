// Package vend handles interactions with the Vend API.
package vend

import (
	"encoding/json"
	"time"
)

// Product is the basic Vend product structure.
type Product struct {
	ID              *string    `json:"id,omitempty"`
	SourceID        *string    `json:"source_id,omitempty"`
	SourceVariantID *string    `json:"source_variant_id,omitempty"`
	VariantParentID *string    `json:"variant_parent_id,omitempty"`
	Name            *string    `json:"name,omitempty"`
	VariantName     *string    `json:"variant_name,omitempty"`
	Handle          *string    `json:"handle,omitempty"`
	SKU             *string    `json:"sku,omitempty"`
	Active          *bool      `json:"active,omitempty"`
	Source          *string    `json:"source"`
	DeletedAt       *time.Time `json:"deleted_at,omitempty"`
	Version         *int64     `json:"version"`
	ImageURL        *string    `json:"image_url,omitempty"`
	Images          []Image    `json:"images,omitempty"`
}

// Image is the info contained in each Vend image object.
type Image struct {
	ID      *string `json:"id,omitempty"`
	URL     *string `json:"url,omitempty"`
	Version *int64  `json:"version"`
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

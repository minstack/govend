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
	products := []Product{}
	page := []Product{}
	data := []byte{}
	var v int64

	// v is a version that is used to get products by page.
	// Here we get the first page.
	data, v, err := c.ResourcePage(0, "GET", "products")
	if err != nil {
		fmt.Println(err)
	}

	// Unmarshal payload into product object.
	err = json.Unmarshal(data, &page)
	if err != nil {
		fmt.Println(err)
	}

	products = append(products, page...)

	for len(page) > 0 {
		page = []Product{}

		// Continue grabbing pages until we receive an empty one.
		data, v, err = c.ResourcePage(v, "GET", "products")
		if err != nil {
			fmt.Println(err)
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

// UploadImage uploads a single product image to Vend.
func (c Client) UploadImage(imagePath string, product ProductUpload) error {
	var err error

	// This checks we actually have an image to post.
	if len(product.ImageURL) > 0 {

		// First grab and save the image from the URL.
		imageURL := fmt.Sprintf("%s", product.ImageURL)

		var body bytes.Buffer
		// Start multipart writer.
		writer := multipart.NewWriter(&body)

		// Key "image" value is the image binary.
		var part io.Writer
		part, err = writer.CreateFormFile("image", imageURL)
		if err != nil {
			fmt.Printf("Error creating multipart form file")
			return err
		}

		// Open image file.
		var file *os.File
		file, err = os.Open(imagePath)
		if err != nil {
			fmt.Printf("Error opening image file")
			return err
		}

		// Make sure file is closed and then removed at end.
		defer file.Close()
		defer os.Remove(imageURL)

		// Copying image binary to form file.
		_, err = io.Copy(part, file)
		if err != nil {
			log.Fatalf("Error copying file for requst body: %s", err)
			return err
		}

		err = writer.Close()
		if err != nil {
			fmt.Printf("Error closing writer")
			return err
		}

		// Create the Vend URL to send our image to.
		url := c.ImageUploadURLFactory(product.ID)

		fmt.Printf("Uploading image to %v, ", product.ID)

		req, _ := http.NewRequest("POST", url, &body)

		// Headers
		req.Header.Set("User-agent", "vend-image-upload")
		req.Header.Set("Content-Type", writer.FormDataContentType())
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token))

		client := &http.Client{}

		// Make the request.
		var attempt int
		var res *http.Response
		for {
			time.Sleep(time.Second)
			res, err = client.Do(req)
			if err != nil || !ResponseCheck(res.StatusCode) {
				log.Fatalf("Couldnt source image: %s", res.Status)
				// Delays between attempts will be exponentially longer each time.
				attempt++
				delay := BackoffDuration(attempt)
				time.Sleep(delay)
			} else {
				// Ensure that image file is removed after it's uploaded.
				os.Remove(imagePath)
				break
			}
		}

		// Make sure response body is closed at end.
		defer res.Body.Close()

		var resBody []byte
		resBody, err = ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Printf("Error reading response body")
			return err
		}

		// Unmarshal JSON response into our respone struct.
		// from this we can find info about the image status.
		response := ImageUpload{}
		err = json.Unmarshal(resBody, &response)
		if err != nil {
			fmt.Printf("Error unmarhsalling response body")
			return err
		}

		payload := response.Data

		fmt.Printf("image created at Position: %v\n\n", *payload.Position)

	}
	return err
}

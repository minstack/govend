// Package vend handles interactions with the Vend API.
package vend

import (
	"encoding/json"
	"log"
)

// Supplier contains supplier data.
type Supplier struct {
	ID          string    `json:"id"`
	RetailerID  string    `json:"retailer_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Source      string    `json:"source"`
	Contact     []Contact `json:"contact"`
}

// Contact is a supplier object
type Contact struct {
	FirstName         string `json:"first_name"`
	LastName          string `json:"last_name"`
	CompanyName       string `json:"company_name"`
	Phone             string `json:"phone"`
	Mobile            string `json:"mobile"`
	Fax               string `json:"fax"`
	Email             string `json:"email"`
	Twitter           string `json:"twitter"`
	Website           string `json:"website"`
	PhysicalAddress1  string `json:"physical_address1"`
	PhysicalAddress2  string `json:"physical_address2"`
	PhysicalSuburb    string `json:"physical_suburb"`
	PhysicalCity      string `json:"physical_city"`
	PhysicalPostcode  string `json:"physical_postcode"`
	PhysicalState     string `json:"physical_state"`
	PhysicalCountryID string `json:"physical_country_id"`
	PostalAddress1    string `json:"postal_address1"`
	PostalAddress2    string `json:"postal_address2"`
	PostalSuburb      string `json:"postal_suburb"`
	PostalCity        string `json:"postal_city"`
	PostalPostcode    string `json:"postal_postcode"`
	PostalState       string `json:"postal_state"`
	PostalCountryID   string `json:"postal_country_id"`
}

// TO DO: Supplier uses API 0.9 for details.
// Suppliers grabs and collates all suppliers in pages of 10,000.
func (c Client) Suppliers() ([]Supplier, error) {

	suppliers := []Supplier{}
	page := []Supplier{}

	// v is a version that is used to get supplier by page.
	// Here we get the first page.
	data, v, err := ResourcePage(0, c.DomainPrefix, c.Token, "suppliers")

	// Unmarshal payload into sales object.
	err = json.Unmarshal(data, &page)
	if err != nil {
		log.Printf("error while unmarshalling: %s", err)
	}

	suppliers = append(suppliers, page...)

	for len(page) > 0 {
		page = []Supplier{}

		// Continue grabbing pages until we receive an empty one.
		data, v, err = ResourcePage(v, c.DomainPrefix, c.Token, "suppliers")
		if err != nil {
			return nil, err
		}

		// Unmarshal payload into customer object.
		err = json.Unmarshal(data, &page)

		// Append customer page to list of customers.
		suppliers = append(suppliers, page...)
	}

	return suppliers, err
}

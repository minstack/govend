// Package vend handles interactions with the Vend API.
package vend

// Supplier contains supplier data.
type SupplierBase struct {
	Name        string   `json:"name"`
	Description *string  `json:"description"`
	Source      *string  `json:"source"`
	Contact     *Contact `json:"contact"`
}

// Supplier contains supplier data.
type Supplier struct {
	ID *string `json:"id"`
	SupplierBase
}

// Contact is a supplier object
type Contact struct {
	FirstName         *string `json:"first_name"`
	LastName          *string `json:"last_name"`
	CompanyName       *string `json:"company_name"`
	Phone             *string `json:"phone"`
	Mobile            *string `json:"mobile"`
	Fax               *string `json:"fax"`
	Email             *string `json:"email"`
	Twitter           *string `json:"twitter"`
	Website           *string `json:"website"`
	PhysicalAddress1  *string `json:"physical_address1"`
	PhysicalAddress2  *string `json:"physical_address2"`
	PhysicalSuburb    *string `json:"physical_suburb"`
	PhysicalCity      *string `json:"physical_city"`
	PhysicalPostcode  *string `json:"physical_postcode"`
	PhysicalState     *string `json:"physical_state"`
	PhysicalCountryID *string `json:"physical_country_id"`
	PostalAddress1    *string `json:"postal_address1"`
	PostalAddress2    *string `json:"postal_address2"`
	PostalSuburb      *string `json:"postal_suburb"`
	PostalCity        *string `json:"postal_city"`
	PostalPostcode    *string `json:"postal_postcode"`
	PostalState       *string `json:"postal_state"`
	PostalCountryID   *string `json:"postal_country_id"`
}

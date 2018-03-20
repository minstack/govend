// Package vend handles interactions with the Vend API.
package vend

import (
	"encoding/json"
	"time"
)

// ConsignmentPayload contains register data and versioning info.
type ConsignmentPayload struct {
	Data    []Consignment    `json:"data,omitempty"`
	Version map[string]int64 `json:"version,omitempty"`
}

// Consignment is a register object.
type Consignment struct {
	ID              *string    `json:"id,omitempty"`
	OutletID        *string    `json:"outlet_id,omitempty"`
	Name            *string    `json:"name,omitempty"`
	Type            *string    `json:"type,omitempty"`
	Status          *string    `json:"status,omitempty"`
	ConsignmentDate *string    `json:"consignment_date,omitempty"` // NOTE: Using string for ParseVendDT.
	DeletedAt       *time.Time `json:"deleted_at,omitempty"`
}

// Consignments gets all stock consignments and transfers from a store.
func (c Client) Consignments() ([]Consignment, error) {

	var consignments, page []Consignment
	var v int64

	// v is a version that is used to objects by page.
	// Here we get the first page.
	data, v, err := ResourcePage(0, c.DomainPrefix, c.Token, "consignments")

	// Unmarshal payload into sales object.
	err = json.Unmarshal(data, &page)

	// Append page to list.
	consignments = append(consignments, page...)

	// NOTE: Turns out empty response is 2bytes.
	for len(data) > 2 {
		page = []Consignment{}

		// Continue grabbing pages until we receive an empty one.
		data, v, err = ResourcePage(v, c.DomainPrefix, c.Token, "consignments")
		if err != nil {
			return nil, err
		}

		// Unmarshal payload into sales object.
		err = json.Unmarshal(data, &page)

		// Append page to list.
		consignments = append(consignments, page...)
	}

	return consignments, err
}

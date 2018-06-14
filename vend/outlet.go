// Package vend handles interactions with the Vend API.
package vend

import (
	"encoding/json"
	"log"
	"time"
)

// Vend API Docs: https://docs.vendhq.com/v0.9/reference#outlets-2

// OutletPayload contains outlet data and versioning info.
type OutletPayload struct {
	Data    []Outlet         `json:"data,omitempty"`
	Version map[string]int64 `json:"version,omitempty"`
}

// Outlet is usually a physical store location.
type Outlet struct {
	ID        *string    `json:"id,omitempty"`
	Name      *string    `json:"name,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

// Outlets gets all outlets from a store.
func (c Client) Outlets() ([]Outlet, map[string][]Outlet, error) {

	outlets := []Outlet{}
	page := []Outlet{}

	// v is a version that is used to get outlets by page.
	data, v, err := c.ResourcePage(0, "GET", "outlets")

	// Unmarshal payload into outlet object.
	err = json.Unmarshal(data, &page)
	if err != nil {
		log.Printf("error while unmarshalling: %s", err)
	}

	outlets = append(outlets, page...)

	for len(page) > 0 {
		page = []Outlet{}

		// Continue grabbing pages until we receive an empty one.
		data, v, err = c.ResourcePage(v, "GET", "outlets")
		if err != nil {
			return nil, nil, err
		}

		// Unmarshal payload into outlet object.
		err = json.Unmarshal(data, &page)

		// Append outlet page to list of outlets.
		outlets = append(outlets, page...)
	}

	outletMap := make(map[string][]Outlet)
	for _, outlet := range outlets {
		outletMap[*outlet.ID] = append(outletMap[*outlet.ID], outlet)
	}

	return outlets, outletMap, err
}

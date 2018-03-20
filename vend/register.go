// Package vend handles interactions with the Vend API.
package vend

import (
	"encoding/json"
	"log"
	"time"
)

// RegisterPayload contains register data and versioning info.
type RegisterPayload struct {
	Data    []Register       `json:"data,omitempty"`
	Version map[string]int64 `json:"version,omitempty"`
}

// Register is a register object.
type Register struct {
	ID        *string    `json:"id,omitempty"`
	Name      *string    `json:"name,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

// Registers gets all registers from a store.
func (c Client) Registers() ([]Register, error) {

	registers := []Register{}
	page := []Register{}

	// v is a version that is used to get registers by page.
	// Here we get the first page.
	data, v, err := ResourcePage(0, c.DomainPrefix, c.Token, "registers")

	// Unmarshal payload into sales object.
	err = json.Unmarshal(data, &page)
	if err != nil {
		log.Printf("error while unmarshalling: %s", err)
	}

	registers = append(registers, page...)

	for len(page) > 0 {
		page = []Register{}

		// Continue grabbing pages until we receive an empty one.
		data, v, err = ResourcePage(v, c.DomainPrefix, c.Token, "registers")
		if err != nil {
			return nil, err
		}

		// Unmarshal payload into register object.
		err = json.Unmarshal(data, &page)

		// Append register page to list of registers.
		registers = append(registers, page...)
	}

	return registers, err
}

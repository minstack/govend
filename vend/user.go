// Package vend handles interactions with the Vend API.
package vend

import (
	"encoding/json"
	"log"
	"time"
)

// UserPayload contains sales data and versioning info.
type UserPayload struct {
	Data    []User           `json:"data,omitempty"`
	Version map[string]int64 `json:"version,omitempty"`
}

// User is a basic user object.
type User struct {
	ID          *string    `json:"id,omitempty"`
	Username    *string    `json:"username,omitempty"`
	DisplayName *string    `json:"display_name,omitempty"`
	Email       *string    `json:"email,omitempty"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}

// Users gets all users from a store.
func (c Client) Users() ([]User, error) {

	users := []User{}
	page := []User{}

	// v is a version that is used to get registers by page.
	// Here we get the first page.
	data, v, err := ResourcePage(0, c.DomainPrefix, c.Token, "users")

	// Unmarshal payload into sales object.
	err = json.Unmarshal(data, &page)
	if err != nil {
		log.Printf("error while unmarshalling: %s", err)
	}

	users = append(users, page...)

	for len(page) > 0 {
		page = []User{}

		// Continue grabbing pages until we receive an empty one.
		data, v, err = ResourcePage(v, c.DomainPrefix, c.Token, "users")
		if err != nil {
			return nil, err
		}

		// Unmarshal payload into register object.
		err = json.Unmarshal(data, &page)

		// Append register page to list of registers.
		users = append(users, page...)
	}

	return users, err
}

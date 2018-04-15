// Package vend handles interactions with the Vend API.
package vend

import (
	"encoding/json"
	"fmt"
)

// StoreCreditPayload hold Gift Card data
type StoreCreditPayload struct {
	Data []StoreCredit `json:"data"`
}

// StoreCredit contains Store Credit data
type StoreCredit struct {
	ID                      *string                  `json:"id"`
	CustomerID              *string                  `json:"customer_id"`
	CreatedAt               *string                  `json:"created_at"`
	Customer                *string                  `json:"customer"`
	Balance                 *float64                 `json:"balance"`
	TotalCreditIssued       *float64                 `json:"total_credit_issued"`
	TotalCreditRedeemed     *float64                 `json:"total_credit_redeemed"`
	StoreCreditTransactions []StoreCreditTransaction `json:"store_credit_transactions"`
}

// StoreCreditTransaction is a Store Credit object.
type StoreCreditTransaction struct {
	ID           *string `json:"id,omitempty"`
	CustomerCode string  `json:"-"`
	CustomerID   *string `json:"-"`
	Amount       float64 `json:"amount"`
	Type         string  `json:"type"`
	Notes        *string `json:"notes"`
	UserID       *string `json:"user_id"`
	SaleID       *string `json:"sale_id,omitempty"`
	ClientID     *string `json:"client_id,omitempty"`
	CreatedAt    *string `json:"created_at,omitempty"`
}

// StoreCredits gets all Store Credit data from a store.
func (c Client) StoreCredits() ([]StoreCredit, error) {

	storecredits := []StoreCredit{}

	// Here we get the first page.
	data, lastID, err := c.ResourcePageFlake("", "GET", "store_credits")
	if err != nil {
		return []StoreCredit{}, fmt.Errorf("Failed to retrieve a page of data %v", err)
	}

	payload := StoreCreditPayload{}

	// Unmarshal payload into Store Credit object.
	err = json.Unmarshal(data, &payload)
	if err != nil {
		return []StoreCredit{}, err
	}

	fmt.Printf(string(data))

	// Append page to list.
	storecredits = append(storecredits, payload.Data...)

	// NOTE: Turns out empty response is 2bytes.
	for len(payload.Data) > 1 {
		payload = StoreCreditPayload{}

		// Continue grabbing pages until we receive an empty one.
		data, lastID, err = c.ResourcePageFlake(lastID, "GET", "store_credits")
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(data, &payload)
		if err != nil {
			return []StoreCredit{}, err
		}

		// Last page will always return a Store Credit from the previous payload, removes the last Store Credit.
		if len(payload.Data) > 1 {
			storecredits = append(storecredits, payload.Data...)
		}
	}

	return storecredits, err
}

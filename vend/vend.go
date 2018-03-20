// Package vend handles interactions with the Vend API.
package vend

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"net/url"
	"os"
	"time"
)

// Client contains API authentication details.
type Client struct {
	Token        string
	DomainPrefix string
	TimeZone     string
}

// NewClient is called to pass authentication details to the manager.
func NewClient(token, domainPrefix, tz string) Client {
	return Client{token, domainPrefix, tz}
}

// ResourcePage gets a single page of data from a 2.0 API resource using a version attribute.
func ResourcePage(version int64, domainPrefix, key, resource string) ([]byte, int64, error) {

	// Build the URL for the resource page.
	url := urlFactory(version, domainPrefix, "", resource)

	body, err := GetDataFromURL(key, url)
	if err != nil {
		fmt.Printf("Error getting resource: %s", err)
	}

	// Decode the raw JSON.
	response := Payload{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Printf("\nError unmarshalling payload: %s", err)
		return nil, 0, err
	}

	// Data is the resource body.
	data := response.Data

	// Version contains the maximum version number of the resources.
	version = response.Version["max"]

	return data, version, err
}

// ResourcePageFlake gets a single page of data from a 2.0 API resource using a Flake ID attribute.
func ResourcePageFlake(id, domainPrefix, key, resource string) ([]byte, string, error) {

	// Build the URL for the resource page.
	url := urlFactoryFlake(id, domainPrefix, resource)
	body, err := GetDataFromURL(key, url)
	if err != nil {
		fmt.Printf("Error getting resource: %s", err)
	}
	// Decode the raw JSON.
	payload := map[string][]interface{}{}
	err = json.Unmarshal(body, &payload)
	if err != nil {
		fmt.Printf("\nError unmarshalling payload: %s", err)
		return nil, "", err
	}

	items := payload["data"]

	// Retrieve the last id from the payload to be used to request subsequent page
	// **TODO** Last ID will be stripped as its included in the previous payload, need a better way to handle this
	i := items[(len(items) - 1)]
	m := i.(map[string]interface{})
	lastID := m["id"].(string)

	return body, lastID, err
}

// GetDataFromURL performs a get request on a Vend API endpoint.
func GetDataFromURL(key, url string) ([]byte, error) {

	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("\nError creating http request: %s", err)
		return nil, err
	}

	// Using personal token authentication.
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", key))
	req.Header.Set("User-Agent", "GO-Vend")

	// Doing the request.
	var attempt int
	var res *http.Response
	for {
		res, err = client.Do(req)
		if err != nil {
			fmt.Printf("\nError performing request: %s", err)
			// Delays between attempts will be exponentially longer each time.
			attempt++
			delay := BackoffDuration(attempt)
			time.Sleep(delay)
		} else {
			break
		}
	}
	// Make sure response body is closed at end.
	defer res.Body.Close()

	// Check for invalid status codes.
	ResponseCheck(res.StatusCode)

	// Read what we got back.
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("\nError while reading response body: %s\n", err)
		return nil, err
	}

	return body, err
}

// ResponseCheck checks the HTTP status codes of responses.
func ResponseCheck(statusCode int) bool {
	switch statusCode {
	case 200, 201:
		return true
	case 401:
		fmt.Printf("\nAccess denied - check personal API token. Status: %d", statusCode)
		os.Exit(0)
	case 404:
		fmt.Printf("\nURL not found - check domain prefix. Status: %d", statusCode)
		os.Exit(0)
	case 429:
		fmt.Printf("\nRate limited by the Vend API :S Status: %d", statusCode)
	case 500:
		fmt.Printf("\nServer error. Status: %d", statusCode)
	case 502:
		fmt.Printf("\nServer received an invalid response :S Status: %d", statusCode)
		os.Exit(0)
	default:
		fmt.Printf("\nGot an unknown status code - Google it. Status: %d", statusCode)
	}
	return false
}

// BackoffDuration ...
func BackoffDuration(attempt int) time.Duration {
	if attempt <= 0 {
		attempt = 1
	}
	seconds := math.Pow(float64(attempt), 3.5) + 5
	return time.Second * time.Duration(seconds)
}

// urlFactory creates a Vend API 2.0 URL based on a resource.
func urlFactory(version int64, domainPrefix, objectID, resource string) string {
	// Page size is capped at ten thousand for all endpoints except sales which it is capped at five hundred.
	const (
		pageSize = 10000
		deleted  = true
	)

	// Using 2.x Endpoint.
	address := fmt.Sprintf("https://%s.vendhq.com/api/2.0/", domainPrefix)
	query := url.Values{}
	query.Add("after", fmt.Sprintf("%d", version))

	if objectID != "" {
		address += fmt.Sprintf("%s/%s/products?%s", resource, objectID, query.Encode())
	} else {
		address += fmt.Sprintf("%s?%s", resource, query.Encode())
	}

	return address
}

// urlFactoryFlake creates a Vend API 2.0 URL based on a resource.
func urlFactoryFlake(id, domainPrefix, resource string) string {
	// Page size is capped at ten thousand for all endpoints except sales which it is capped at five hundred.
	const (
		pageSize = 10000
		deleted  = true
	)

	// Using 2.x Endpoint.
	address := fmt.Sprintf("https://%s.vendhq.com/api/2.0/%s", domainPrefix, resource)

	// Iterate through pages using the ?before= FLAKE ID attribute.
	if id != "" {
		query := url.Values{}
		query.Add("before", fmt.Sprintf("%s", id))
		address += fmt.Sprintf("?%s", query.Encode())
	}

	return address
}

// ImageUploadURLFactory creates the Vend URL for uploading an image.
func ImageUploadURLFactory(domainPrefix, productID string) string {
	url := fmt.Sprintf("https://%s.vendhq.com/api/2.0/products/%s/actions/image_upload",
		domainPrefix, productID)
	return url
}

// ParseVendDT converts the default Vend timestamp string into a go Time.time value.
func ParseVendDT(dt, tz string) time.Time {

	// Load store's timezone as location.
	loc, err := time.LoadLocation(tz)
	if err != nil {
		fmt.Printf("Error loading timezone as location: %s", err)
	}

	// Default Vend timedate layout.
	const longForm = "2006-01-02T15:04:05Z07:00"
	t, err := time.Parse(longForm, dt)
	if err != nil {
		log.Fatalf("Error parsing time into deafult timestamp: %s", err)
	}

	// Time in retailer's timezone.
	dtWithTimezone := t.In(loc)

	return dtWithTimezone

	// Time string with timezone removed.
	// timeStr := timeLoc.String()[0:19]
}

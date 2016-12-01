package realtime

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

// BASERIPTAAPIURL is the base url for the realtime api
const BASERIPTAAPIURL = "http://realtime.ripta.com:81"

// RIPTAAPITripUpdates is the url for trip updates
const RIPTAAPITripUpdates = "/api/tripupdates?format=json"

// RIPTAAPIVehiclePositions is the url for vehicle positions
const RIPTAAPIVehiclePositions = "/api/vehiclepositions?format=json"

// RIPTAAPIServiceAlerts is the url for service alerts
const RIPTAAPIServiceAlerts = "/api/servicealerts?format=json"

type Client struct {
	BaseURL *url.URL
	client  *http.Client

	common service

	TripUpdates      *TripUpdatesService
	VehiclePositions *VehiclePositionsService
	ServiceAlerts    *ServiceAlertsService
}

type service struct {
	client *Client
}

func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	baseURL, _ := url.Parse(BASERIPTAAPIURL)

	c := &Client{
		client:  httpClient,
		BaseURL: baseURL,
	}

	c.common.client = c

	c.TripUpdates = (*TripUpdatesService)(&c.common)
	c.VehiclePositions = (*VehiclePositionsService)(&c.common)
	c.ServiceAlerts = (*ServiceAlertsService)(&c.common)

	return c
}

func (c *Client) NewRequest(method, urlStr string, body interface{}) (*http.Request, error) {
	rel, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	u := c.BaseURL.ResolveReference(rel)

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return req, nil
}

func (c *Client) Do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	if v != nil {
		err = json.NewDecoder(resp.Body).Decode(v)
		if err == io.EOF {
			err = nil // ignore EOF errors caused by empty response body
		}
	}

	return resp, err
}

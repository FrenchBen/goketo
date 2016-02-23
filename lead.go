package goketo

import (
	"encoding/json"
	"net/url"
	"strconv"
	"strings"

	"github.com/Sirupsen/logrus"
)

// Leads response from list request
type Leads struct {
	client    *Client
	RequestID string          `json:"requestId"`
	Result    json.RawMessage `json:"result"`
	Success   bool            `json:"success"`
	Next      string          `json:"nextPageToken,omitempty"`
}

// LeadResult default result struct as part of the lead - can be customized to allow greater fields
type LeadResult struct {
	ID        int    `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Created   string `json:"createdAt"`
	Updated   string `json:"updatedAt"`
}

// LeadRequest builds a request for data retrieval
type LeadRequest struct {
	ID     int    // List ID
	Next   string // Next page Token
	Fields string
}

// LeadUpdate builds the data for an update
type LeadUpdate struct {
	Action string          `json:"action"` // createOnly - updateOnly - createOrUpdate(default request) - createDuplicate
	Lookup string          `json:"lookupField"`
	Input  json.RawMessage `json:"input"`
}

// LeadUpdateResponse data format for update response
type LeadUpdateResponse struct {
	ID      string             `json:"requestId"`
	Success bool               `json:"success"`
	Result  []LeadUpdateResult `json:"result,omitempty"`
	Error   []LeadError        `json:"errors,omitempty"`
}

// LeadUpdateResult holds result for all updates
type LeadUpdateResult struct {
	ID     int    `json:"id"`
	Status string `json:"status"`
}

// LeadError shows the error code and message for response
type LeadError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Leads Get leads by list Id
func (c *Client) Leads(list *LeadRequest) (leads *Leads, err error) {
	var nextPage string
	if list.Next != "" {
		nextPage = "?nextPageToken=" + list.Next
	} else {
		nextPage = ""
	}
	body, err := c.Get("/list/" + strconv.Itoa(list.ID) + "/leads.json" + nextPage)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &leads)
	leads.client = c
	return leads, err
}

// Lead Get lead by Id - aka member by ID
func (c *Client) Lead(leadReq *LeadRequest) (lead *Leads, err error) {
	fields := url.Values{}
	if len(leadReq.Fields) > 0 {
		fields.Set("fields", strings.Join(strings.Fields(leadReq.Fields), ""))
	}
	logrus.Info("Fields: ", fields.Encode())
	body, err := c.Get("/lead/" + strconv.Itoa(leadReq.ID) + ".json" + "?" + fields.Encode())
	if err != nil {
		return
	}
	err = json.Unmarshal(body, &lead)
	lead.client = c
	return
}

// UpdateLeads post update of data for a lead
func (c *Client) UpdateLeads(update *LeadUpdate) ([]byte, error) {
	data, err := json.Marshal(update)
	if err != nil {
		return nil, err
	}
	body, err := c.Post("/leads.json", data)
	return body, err
}

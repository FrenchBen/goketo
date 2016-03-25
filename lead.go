package goketo

import (
	"encoding/json"
	"net/url"
	"strconv"
	"strings"

	"github.com/Sirupsen/logrus"
)

// Convert all to Interfaces for re-usability
// Add fmt.Sprintf("%v:%v", host, port) to build strings

// LeadResponse response from list request
type LeadResponse struct {
	client    *Client
	RequestID string          `json:"requestId"`
	Result    json.RawMessage `json:"result"`
	Success   bool            `json:"success"`
	Next      string          `json:"nextPageToken,omitempty"`
	More      bool            `json:"moreResult,omitempty"`
	Errors    []struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"errors,omitempty"`
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

// LeadFieldResponse response for all fields
type LeadFieldResponse struct {
	client    *Client
	RequestID string      `json:"requestId"`
	Result    []LeadField `json:"result"`
	Success   bool        `json:"success"`
}

// LeadField describes all possible fields for Leads
type LeadField struct {
	ID     int    `json:"id"`
	Name   string `json:"displayName"`
	Type   string `json:"dataType"`
	Length int    `json:"length"`
	Rest   struct {
		Name     string `json:"name"`
		ReadOnly bool   `json:"readOnly"`
	} `json:"rest"`
	Soap struct {
		Name     string `json:"name"`
		ReadOnly bool   `json:"readOnly"`
	} `json:"soap"`
}

// DeletedLeadResponse response of Deleted lead request
type DeletedLeadResponse struct {
	*LeadResponse
	Result []DeletedLead `json:"result"`
}

// DeletedLead result
type DeletedLead struct {
	ID         int      `json:"id"`
	LeadID     int      `json:"leadId"`
	Date       string   `json:"activityDate"`
	TypeID     int      `json:"activityTypeId"`
	PrimaryID  int      `json:"primaryAttributeValueId"`
	PrimaryVal string   `json:"primaryAttributeValue"`
	Attributes []string `json:"attributes"`
}

// Leads Get leads by list Id
func (c *Client) Leads(leadReq *LeadRequest) (leads *LeadResponse, err error) {
	nextPage := url.Values{}
	if leadReq.Next != "" {
		nextPage.Set("&nextPageToken", leadReq.Next)
	}
	fields := url.Values{}
	if len(leadReq.Fields) > 0 {
		fields.Set("fields", strings.Join(strings.Fields(leadReq.Fields), ""))
	}
	logrus.Debug("Fields: ", fields.Encode())
	body, err := c.Get("/list/" + strconv.Itoa(leadReq.ID) + "/leads.json" + "?" + fields.Encode() + nextPage.Encode())
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &leads)
	leads.client = c
	return leads, err
}

// Lead Get lead by Id - aka member by ID
func (c *Client) Lead(leadReq *LeadRequest) (lead *LeadResponse, err error) {
	fields := url.Values{}
	if len(leadReq.Fields) > 0 {
		fields.Set("fields", strings.Join(strings.Fields(leadReq.Fields), ""))
	}
	logrus.Debug("Fields: ", fields.Encode())
	body, err := c.Get("/lead/" + strconv.Itoa(leadReq.ID) + ".json" + "?" + fields.Encode())
	if err != nil {
		return
	}
	err = json.Unmarshal(body, &lead)
	lead.client = c
	return lead, err
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

// LeadFields return all fields and the data type of a lead object
func (c *Client) LeadFields() (fields *LeadFieldResponse, err error) {
	body, err := c.Get("/leads/describe.json")
	err = json.Unmarshal(body, &fields)
	return fields, err
}

// DeletedLeads returns a list of leads that were deleted
func (c *Client) DeletedLeads(leadReq *LeadRequest) (deletedLeads *DeletedLeadResponse, err error) {
	nextPage := url.Values{}
	if leadReq.Next != "" {
		nextPage.Set("nextPageToken", leadReq.Next)
	}
	body, err := c.Get("/activities/deletedleads.json?" + nextPage.Encode())
	err = json.Unmarshal(body, &deletedLeads)
	return deletedLeads, err
}

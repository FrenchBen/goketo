package goketo

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

// Convert all to Interfaces for re-usability
// Add fmt.Sprintf("%v:%v", host, port) to build strings
// apiResponse common api response structure
type apiResponse struct {
	RequestID string `json:"requestId"`
	Success   bool   `json:"success"`
	Next      string `json:"nextPageToken,omitempty"`
	More      bool   `json:"moreResult,omitempty"`
	Errors    []struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"errors,omitempty"`
}

// LeadResponse response from list request
type LeadResponse struct {
	apiResponse
	Result json.RawMessage `json:"result,omitempty"`
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
	Errors  []struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"errors,omitempty"`
}

// LeadUpdateResult holds result for all updates
type LeadUpdateResult struct {
	ID     int    `json:"id"`
	Status string `json:"status"`
}

// LeadFieldResponse response for all fields
type LeadFieldResponse struct {
	client *Client
	apiResponse
	Result []LeadField `json:"result"`
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
func Leads(req Requester, leadReq *LeadRequest) (leads *LeadResponse, err error) {
	urlQuery := url.Values{}
	if leadReq.Next != "" {
		urlQuery.Set("nextPageToken", leadReq.Next)
	}
	if len(leadReq.Fields) > 0 {
		// join fields that are separate by space ' ' instead of ','
		urlQuery.Set("fields", strings.Join(strings.Fields(leadReq.Fields), ","))
	}
	url := fmt.Sprintf("list/%s/leads.json?%s", strconv.Itoa(leadReq.ID), urlQuery.Encode())
	logrus.Debug("Get: ", url)
	body, err := req.Get(url)
	if err != nil {
		return nil, err
	}
	logrus.Debug("Body: ", string(body))
	err = json.Unmarshal(body, &leads)
	return leads, err
}

// Lead Get lead by Id - aka member by ID
func Lead(req Requester, leadReq *LeadRequest) (lead *LeadResponse, err error) {
	urlQuery := url.Values{}
	if len(leadReq.Fields) > 0 {
		// join fields that are separate by space ' ' instead of ','
		urlQuery.Set("fields", strings.Join(strings.Fields(leadReq.Fields), ","))
	}
	url := fmt.Sprintf("lead/%s.json?%s", strconv.Itoa(leadReq.ID), urlQuery.Encode())
	logrus.Debug("Get: ", url)
	body, err := req.Get(url)
	if err != nil {
		return
	}
	logrus.Debug("Body: ", string(body))
	err = json.Unmarshal(body, &lead)
	return lead, err
}

// LeadsFilter Get leads by filter Type
// Common filter types:
//  - id
//  - cookies
//  - email
//  - twitterId
//  - facebookId
//  - linkedInId
//  - sfdcAccountId
//  - sfdcContactId
//  - sfdcLeadId
//  - sfdcLeadOwnerId
//  - sfdcOpptyId
func LeadsFilter(req Requester, leadReq *LeadRequest, filterType string, filterValues []string) (leads *LeadResponse, err error) {
	urlQuery := url.Values{}
	urlQuery.Set("filterType", filterType)
	if leadReq.Next != "" {
		urlQuery.Set("nextPageToken", leadReq.Next)
	}
	if len(leadReq.Fields) > 0 {
		// join fields that are separate by space ' ' instead of ','
		urlQuery.Set("fields", strings.Join(strings.Fields(leadReq.Fields), ","))
	}
	if len(filterValues) > 0 {
		urlQuery.Set("filterValues", strings.Join(filterValues, ","))
	}
	url := fmt.Sprintf("leads.json?%s", urlQuery.Encode())
	logrus.Debug("Get: ", url)
	body, err := req.Get(url)
	if err != nil {
		return nil, err
	}
	logrus.Debug("Body: ", string(body))
	err = json.Unmarshal(body, &leads)
	return leads, err
}

// UpdateLeads post update of data for a lead
func UpdateLeads(req Requester, update *LeadUpdate) ([]byte, error) {
	data, err := json.Marshal(update)
	if err != nil {
		return nil, err
	}
	body, err := req.Post("leads.json", data)
	return body, err
}

// LeadFields return all fields and the data type of a lead object
func LeadFields(req Requester) (fields *LeadFieldResponse, err error) {
	body, err := req.Get("leads/describe.json")
	err = json.Unmarshal(body, &fields)
	return fields, err
}

// DeletedLeads returns a list of leads that were deleted
func DeletedLeads(req Requester, leadReq *LeadRequest) (deletedLeads *DeletedLeadResponse, err error) {
	urlQuery := url.Values{}
	if leadReq.Next != "" {
		urlQuery.Set("&nextPageToken", leadReq.Next)
	}
	url := fmt.Sprintf("activities/deletedleads.json?%s", urlQuery.Encode())
	logrus.Debug("Get: ", url)
	body, err := req.Get(url)
	err = json.Unmarshal(body, &deletedLeads)
	return deletedLeads, err
}

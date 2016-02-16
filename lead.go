package goketo

import "encoding/json"

// Lead response from lead request
type Lead struct {
	client    *Client
	RequestID string       `json:"requestId"`
	Result    []LeadResult `json:"result"`
	Success   bool         `json:"success"`
}

// Leads response from list request
type Leads struct {
	client    *Client
	RequestID string       `json:"requestId"`
	Result    []LeadResult `json:"result"`
	Success   bool         `json:"success"`
	Next      string       `json:"nextPageToken"`
}

// LeadResult result struct as part of the lead
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
	ID   string // List ID
	Next string // Next page Token
}

// LeadUpdate builds the data for an update
type LeadUpdate struct {
	Action string        `json:"action"`
	Lookup string        `json:"lookupField"`
	Data   []interface{} `json:"input"`
}

// Leads Get leads by list Id
func (c *Client) Leads(list *LeadRequest) (leads *Leads, err error) {
	var nextPage string
	if list.Next != "" {
		nextPage = "?nextPageToken=" + list.Next
	} else {
		nextPage = ""
	}
	body, err := c.Get("/list/" + list.ID + "/leads.json" + nextPage)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &leads)
	leads.client = c
	return
}

// UpdateLeads post update of data for a lead
func (c *Client) UpdateLeads(update LeadUpdate) ([]byte, error) {
	data, err := json.Marshal(update)
	if err != nil {
		return nil, err
	}
	body, err := c.Post("/leads.json", data)
	return body, err
}

// Lead Get lead by Id - aka member by ID
func (c *Client) Lead(leadID string) (lead *Lead, err error) {
	body, err := c.Get("/lead/" + leadID + ".json")
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &lead)
	lead.client = c
	return
}

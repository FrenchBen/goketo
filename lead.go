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

// Leads Get leads by list Id
func (c *Client) Leads(listID string) (leads *Leads, err error) {
	body, err := c.Get("/list/" + listID + "/leads.json")
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &leads)
	leads.client = c
	return
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

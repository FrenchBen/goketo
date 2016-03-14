package goketo

import "encoding/json"

// ErrorResponse response from list request
type ErrorResponse struct {
	client    *Client
	RequestID string `json:"requestId"`
	Result    struct {
		Date   string  `json:"date"`
		Total  int     `json:"total"`
		Errors []Error `json:"errors"`
	} `json:"result"`
	Success bool `json:"success"`
}

// Error contains code and count
type Error struct {
	Code  string `json:"errorCode"`
	Count int    `json:"count"`
}

// DailyError returns error codes and their count for the day
func (c *Client) DailyError() (errors *ErrorResponse, err error) {
	body, err := c.Get("/stats/errors.json")
	err = json.Unmarshal(body, &errors)
	return errors, err
}

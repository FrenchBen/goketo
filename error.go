package goketo

import "encoding/json"

// ErrorResponse response from list request
type ErrorResponse struct {
	RequestID string `json:"requestId"`
	Result    Result `json:"result"`
	Success   bool   `json:"success"`
}

// Error contains code and count
type Error struct {
	Code  string `json:"errorCode"`
	Count int    `json:"count"`
}

// Result contains result stack
type Result struct {
	Date   string  `json:"date"`
	Total  int     `json:"total"`
	Errors []Error `json:"errors"`
}

// DailyError returns error codes and their count for the day
func DailyError(req Requester) (errors *ErrorResponse, err error) {
	body, err := req.Get("/stats/errors.json")
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(body, &errors); err != nil {
		return nil, err
	}
	return errors, nil
}

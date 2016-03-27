package goketo

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"

	"github.com/Sirupsen/logrus"
	"github.com/jinzhu/now"
)

type pagingToken struct {
	apiResponse
}

// ActivityType is the response for a list of actity types
type ActivityType struct {
	apiResponse
	Result []struct {
		ID               int    `json:"id"`
		Name             string `json:"name"`
		Description      string `json:"description"`
		PrimaryAttribute struct {
			Name     string `json:"name"`
			DataType string `json:"dataType"`
		} `json:"primaryAttribute"`
		Attributes []struct {
			Name     string `json:"name"`
			DataType string `json:"dataType"`
		} `json:"attributes"`
	} `json:"result,omitempty"`
}

// ActivityRequest is the building block for an activity request
type ActivityRequest struct {
	ActivityTypeID string
	DateTime       string
	ListID         string
	LeadIDs        []string
}

// Activity is the response from a get activity request
type Activity struct {
	apiResponse
	Result []struct {
		ID                      int    `json:"id"`
		LeadID                  int    `json:"leadId"`
		ActivityDate            string `json:"activityDate"`
		ActivityTypeID          int    `json:"activityTypeId"`
		PrimaryAttributeValueID int    `json:"primaryAttributeValueId"`
		PrimaryAttributeValue   string `json:"primaryAttributeValue"`
		Attributes              []struct {
			Name  string `json:"name"`
			Value string `json:"value"`
		} `json:"attributes"`
	} `json:"result,omitempty"`
}

// LeadChange response to a Lead Changes request
type LeadChange struct {
	apiResponse
	Result []struct {
		ID             int    `json:"id"`
		LeadID         int    `json:"leadId"`
		ActivityDate   string `json:"activityDate"`
		ActivityTypeID int    `json:"activityTypeId"`
		Fields         []struct {
			ID       int    `json:"id"`
			Name     string `json:"name"`
			NewValue string `json:"newValue"`
			OldValue string `json:"oldValue"`
		} `json:"fields"`
		Attributes []struct {
			Name  string `json:"name"`
			Value string `json:"value"`
		} `json:"attributes"`
	} `json:"result,omitempty"`
}

// getPagingToken sends a request for a paging token to be used with activity
func getPagingToken(req Requester, dateTime string) (string, error) {

	url := fmt.Sprintf("activities/pagingtoken.json?sinceDatetime=%s", dateTime)
	logrus.Debug("Get: ", url)
	body, err := req.Get(url)
	if err != nil {
		return "", err
	}
	logrus.Debug("Body: ", string(body))
	tokenResponse := pagingToken{}
	err = json.Unmarshal(body, &tokenResponse)
	return tokenResponse.Next, err
}

// GetActivityTypes returns a list of activities accepted by the GetActivity request
func GetActivityTypes(req Requester) (activityType *ActivityType) {
	body, err := req.Get("activities/types.json")
	if err != nil {
		return nil
	}
	logrus.Debug("Body: ", string(body))
	err = json.Unmarshal(body, &activityType)
	if err != nil {
		logrus.Error("Error with JSON: ", err)
	}
	return
}

// GetActivity get a series of activities based on a data/time string and optional list/leads ID
func GetActivity(req Requester, activityReq ActivityRequest) (activities *Activity, err error) {
	if activityReq.ActivityTypeID == "" {
		logrus.Error("Missing activity ID")
		return nil, errors.New("missing activity ID")
	}
	tNow, err := now.Parse(activityReq.DateTime)
	if err != nil {
		logrus.Error("couldn't parse date-time: ", err)
		return nil, err
	}
	activityReq.DateTime = tNow.Format("2006-01-02T15:04-0800")
	token, err := getPagingToken(req, activityReq.DateTime)
	if err != nil {
		return nil, err
	}
	urlQuery := url.Values{}
	urlQuery.Set("nextPageToken", token)
	urlQuery.Set("activityTypeIds", activityReq.ActivityTypeID)
	if activityReq.ListID != "" {
		urlQuery.Set("listId", activityReq.ListID)
	}
	if len(activityReq.LeadIDs) > 0 {
		for _, leadID := range activityReq.LeadIDs {
			urlQuery.Add("leadIds", leadID)
		}
	}

	url := fmt.Sprintf("activities.json?%s", urlQuery.Encode())
	logrus.Debug("url: ", url)
	body, err := req.Get(url)
	if err != nil {
		return nil, err
	}
	logrus.Debug("Body: ", string(body))
	err = json.Unmarshal(body, &activities)
	if err != nil {
		logrus.Error("Error with JSON: ", err)
	}
	return activities, err
}

// GetLeadChanges get a series of changes based on a data/time string and list ID
func GetLeadChanges(req Requester, dateTime string, listID string, fields string) (leadChanges *LeadChange, err error) {
	tNow, err := now.Parse(dateTime)
	if err != nil {
		logrus.Error("couldn't parse date-time: ", err)
		return nil, err
	}
	dateTime = tNow.Format("2006-01-02T15:04-0800")
	token, err := getPagingToken(req, dateTime)
	if err != nil {
		return nil, err
	}
	urlQuery := url.Values{}
	urlQuery.Set("nextPageToken", token)
	urlQuery.Set("fields", fields)
	urlQuery.Set("listId", listID)

	url := fmt.Sprintf("activities/leadchanges.json?%s", urlQuery.Encode())
	logrus.Debug("url: ", url)
	body, err := req.Get(url)
	if err != nil {
		return nil, err
	}
	logrus.Debug("Body: ", string(body))
	err = json.Unmarshal(body, &leadChanges)
	if err != nil {
		logrus.Error("Error with JSON: ", err)
	}
	return leadChanges, err
}

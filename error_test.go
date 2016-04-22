package goketo

import (
	"net/http"
	"reflect"
	"testing"
)

//@TODO Implement Unit Test and Integration tests (BATS)

type FakeRequester struct {
	Body []byte
	Err  error
}

func (f FakeRequester) do(req *http.Request) ([]byte, error) {
	if f.Err != nil {
		return nil, f.Err
	}

	return f.Body, nil
}
func (f FakeRequester) Get(resource string) ([]byte, error) {
	if f.Err != nil {
		return nil, f.Err
	}

	return f.Body, nil
}
func (f FakeRequester) Post(resource string) ([]byte, error) {
	if f.Err != nil {
		return nil, f.Err
	}

	return f.Body, nil
}

func TestDailyError(t *testing.T) {
	f := []struct {
		f            FakeRequester
		url          string
		expectedResp ErrorResponse
		expectedErr  error
	}{
		{
			f: FakeRequester{
				Body: []byte(`{
              "requestId": "123ab#456c789de10",
              "result": [
                {
                  "date": "2016-04-07",
                  "total": 10,
                  "errors": [
                    {
                      "errorCode": "1003",
                      "count": 10
                    }
                  ]
                }
              ],
              "success": true
            }`),
				Err: nil,
			},
			url: "/stats/errors.json",
			expectedResp: ErrorResponse{
				RequestID: "123ab#456c789de10",
				Result: {
					Date:  "2016-04-07",
					Total: 10,
					Errors: []Error{
						Code:  "1003",
						Count: 10,
					},
				},
				Success: true,
			},
			expectedErr: nil,
		},
	}
	curr := f[0]
	resp, err := DailyError(curr.f)
	if !reflect.DeepEqual(err, curr.expectedErr) {
		t.Errorf("Expected err to be %q but it was %q", curr.expectedErr, err)
	}
	if !reflect.DeepEqual(resp, curr.expectedResp) {
		t.Fatalf("Expected %v but got %v", curr.expectedResp, resp)
	}

}

package goketo

import (
	"net/http"
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

}

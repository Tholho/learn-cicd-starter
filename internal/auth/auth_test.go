package auth

import (
	"errors"
	"fmt"
	"net/http/httptest"
	"testing"
)

func isSpecificError(expected error) func(error) bool {
	return func(err error) bool {
		return errors.Is(err, expected)
	}
}

func hasErrorMessage(expectedMsg string) func(error) bool {
	return func(err error) bool {
		return err != nil && err.Error() == expectedMsg
	}
}

func TestGetAPIKey(t *testing.T) {
	tests := map[string]struct {
		inputAPIKey string
		want        string
		checkErr    func(error) bool
	}{
		"empty":     {inputAPIKey: "", want: "", checkErr: isSpecificError(ErrNoAuthHeaderIncluded)},
		"malformed": {inputAPIKey: "apikey someincorrectAPIKEY", want: "", checkErr: hasErrorMessage("malformed authorization header")},
		"nominal":   {inputAPIKey: "ApiKey somevalidkey", want: "somevalidkey", checkErr: func(err error) bool { return err == nil }},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			fakeGet := httptest.NewRequest("GET", "/", nil)
			fakeGet.Header.Set("Authorization", tc.inputAPIKey)
			got, err := GetAPIKey(fakeGet.Header)
			fmt.Println(got, err)
			if !tc.checkErr(err) {
				t.Errorf("expected error message: %v", err.Error())
			}
		})
	}
}

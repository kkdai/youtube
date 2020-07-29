package decipher

import (
	"fmt"
	"net/http"
)

type errUnexpectedStatusCode int

func (err errUnexpectedStatusCode) Error() string {
	return fmt.Sprintf("unexpected status code: %d", err)
}

func (d Decipher) httpGet(url string) (resp *http.Response, err error) {
	resp, err = d.client.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, errUnexpectedStatusCode(resp.StatusCode)
	}
	return
}

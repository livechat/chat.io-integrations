package chat_io_capi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type RESTAPI struct {
	httpClient *http.Client
}

func NewRESTAPI() *RESTAPI {
	return &RESTAPI{
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

func (r *RESTAPI) Send(api, version, action, accessToken string, q *url.URL, payload interface{}) error {

	raw, err := json.Marshal(payload)

	req, err := http.NewRequest("POST", api+"/v"+version+"/action/"+action+"?"+q.RawQuery, bytes.NewBuffer(raw)) // <-- URL-encoded payload
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	resp, err := r.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Println(string(body), api+"/v"+version+"/action/"+action+"?"+q.Query().Encode())

	if resp.StatusCode != http.StatusOK {
		return errors.New("internal_error")
	}

	return nil
}

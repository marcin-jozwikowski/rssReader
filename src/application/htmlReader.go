package application

import (
	"io/ioutil"
	"net/http"
	"crypto/tls"
)

func GetHtmlContent(url string) ([]byte, error) {
	transport := &http.Transport{
        TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
    }

    client := &http.Client{Transport: transport}

	resp, err := client.Get(url)
	// handle the error if there is one
	if err != nil {
		return nil, err
	}
	// do this now so it won't be forgotten
	defer resp.Body.Close()
	// reads html as a slice of bytes
	html, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return html, nil
}

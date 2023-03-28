package deliveroo

import (
	"bytes"
	"net/http"
)

type ApolloOperation string

const (
	CreatePaymentPlan ApolloOperation = "CreatePaymentPlan"
	None              ApolloOperation = ""
)

func (d *Deliveroo) sendGET(url string) (*http.Response, error) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	d.setHeaders(req)
	return client.Do(req)
}

func (d *Deliveroo) sendPOST(url string, data []byte, op ApolloOperation) (*http.Response, error) {
	client := &http.Client{}
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(data))
	d.setPostHeaders(req, op)
	return client.Do(req)
}

func (d *Deliveroo) setHeaders(req *http.Request) {
	req.Header.Set("Accept-Language", "en-UK")
	req.Header.Set("User-Agent", "Deliveroo/3.98.0 (samsung SM-G935F;Android 8.0.0;it-IT;releaseEnv release)")
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("X-Roo-Guid", d.uID)
	req.Header.Set("X-Roo-Sticky-Guid", d.uID)
	req.Header.Set("Authorization", d.auth)
}

func setApolloHeaders(req *http.Request, op ApolloOperation) {
	req.Header.Set("X-APOLLO-OPERATION-NAME", string(op))
	req.Header.Set("X-APOLLO-CACHE-FETCH-STRATEGY", "NETWORK_ONLY")
	req.Header.Set("X-APOLLO-EXPIRE-TIMEOUT", "0")
	req.Header.Set("X-APOLLO-EXPIRE-AFTER-READ", "false")
	req.Header.Set("X-APOLLO-PREFETCH", "false")
	req.Header.Set("X-APOLLO-CACHE-DO-NOT-STORE", "false")
}

func (d *Deliveroo) setPostHeaders(req *http.Request, op ApolloOperation) {
	if op != None {
		setApolloHeaders(req, op)
	}

	req.Header.Set("Accept-Language", "en-UK")
	req.Header.Set("User-Agent", "Deliveroo/3.98.0 (samsung SM-G935F;Android 8.0.0;it-IT;releaseEnv release)")
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("X-Roo-Guid", d.uID)
	req.Header.Set("X-Roo-Sticky-Guid", d.uID)
	req.Header.Set("Authorization", d.auth)
}

package deliveroo

import (
	"bytes"
	"github.com/gofrs/uuid"
	"net/http"
)

func (d *Deliveroo) sendGET(url string) (*http.Response, error) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	d.setHeaders(req)
	return client.Do(req)
}

func (d *Deliveroo) sendPOST(url string, data []byte) (*http.Response, error) {
	client := &http.Client{}
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(data))
	d.setPostHeaders(req)
	return client.Do(req)
}

func (d *Deliveroo) setHeaders(req *http.Request) {
	authKeyValue, _ := uuid.DefaultGenerator.NewV1()

	req.Header.Set("Accept-Language", "en-UK")
	req.Header.Set("User-Agent", "Deliveroo/3.98.0 (samsung SM-G935F;Android 8.0.0;it-IT;releaseEnv release)")
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("X-Roo-Guid", authKeyValue.String())
	req.Header.Set("X-Roo-Sticky-Guid", authKeyValue.String())
	req.Header.Set("Authorization", d.auth)
}

func (d *Deliveroo) setPostHeaders(req *http.Request) {
	authKeyValue, _ := uuid.DefaultGenerator.NewV1()

	req.Header.Set("Accept-Language", "en-UK")
	req.Header.Set("User-Agent", "Deliveroo/3.98.0 (samsung SM-G935F;Android 8.0.0;it-IT;releaseEnv release)")
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("X-Roo-Guid", authKeyValue.String())
	req.Header.Set("X-Roo-Sticky-Guid", authKeyValue.String())
	req.Header.Set("Authorization", d.auth)
}

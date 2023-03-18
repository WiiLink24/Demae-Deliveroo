package deliveroo

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func (d *Deliveroo) DownloadAndReturnImage(id string) []byte {
	// If the file exists, serve it
	file, err := os.ReadFile(fmt.Sprintf("./images/%s.jpg", id))
	if err != nil {
		if !os.IsNotExist(err) {
			return nil
		}
	} else {
		return file
	}

	response, err := d.sendGET(fmt.Sprintf("https://api.deliveroo.com/orderapp/v1/restaurants/%s?track=1&lat=%.2f&lng=%.2f&include_unavailable=true&restaurant_fulfillments_supported=true&fulfillment_method=DELIVERY", id, 44.49616860666642, 11.341424435377123))
	if err != nil {
		return nil
	}

	defer response.Body.Close()
	respBytes, _ := io.ReadAll(response.Body)

	jsonData := map[string]any{}
	err = json.Unmarshal(respBytes, &jsonData)
	if err != nil {
		return nil
	}

	imageUrl := jsonData["image_url"].(string)
	imageUrl = strings.Replace(imageUrl, "{w}", "160", -1)
	imageUrl = strings.Replace(imageUrl, "{h}", "160", -1)
	imageUrl = strings.Replace(imageUrl, "{&quality}", "100", -1)

	response, err = http.Get(imageUrl)
	if err != nil {
		return nil
	}

	defer response.Body.Close()
	respBytes, _ = io.ReadAll(response.Body)

	os.WriteFile(fmt.Sprintf("./images/%s.jpg", id), respBytes, 0666)
	return respBytes
}

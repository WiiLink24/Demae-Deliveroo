package deliveroo

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mitchellh/go-wordwrap"
	"golang.org/x/exp/slices"
	"golang.org/x/image/draw"
	"image"
	"image/jpeg"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
)

// GetUserAddress sets the coordinates for delivery. The address must be labeled `demae` or it will return an error
func (d *Deliveroo) GetUserAddress() error {
	response, err := d.sendGET(fmt.Sprintf("https://api.deliveroo.com/orderapp/v1/users/%s/addresses", d.userId))
	if err != nil {
		return err
	}

	defer response.Body.Close()
	respBytes, _ := io.ReadAll(response.Body)

	d.response = string(respBytes)
	d.responseCode = response.StatusCode

	jsonData := map[string]any{}
	err = json.Unmarshal(respBytes, &jsonData)
	if err != nil {
		return err
	}

	for _, address := range jsonData["addresses"].([]any) {
		if address.(map[string]any)["label"].(string) != "demae" {
			continue
		}

		d.longitude = address.(map[string]any)["coordinates"].([]any)[0].(float64)
		d.latitude = address.(map[string]any)["coordinates"].([]any)[1].(float64)
		break
	}

	return nil
}

// GetBareShops returns the absolute bare minimum required for the main menu to display food types
func (d *Deliveroo) GetBareShops() (HasFoodType, error) {
	graphQLQuery, err := GetShopsQuery(d.longitude, d.latitude)
	if err != nil {
		return HasFoodType{}, err
	}

	response, err := d.sendPOST("https://api.deliveroo.com/consumer/graphql/", graphQLQuery, None)
	if err != nil {
		return HasFoodType{}, err
	}

	defer response.Body.Close()
	respBytes, _ := io.ReadAll(response.Body)

	d.response = string(respBytes)
	d.responseCode = response.StatusCode

	jsonData := map[string]any{}
	err = json.Unmarshal(respBytes, &jsonData)
	if err != nil {
		return HasFoodType{}, err
	}

	var has HasFoodType
	filters := jsonData["data"].(map[string]any)["results"].(map[string]any)["ui_control_groups"].(map[string]any)["filters"].([]any)

	for _, filter := range filters {
		if filter.(map[string]any)["id"].(string) == "category" {
			for _, category := range filter.(map[string]any)["options"].([]any) {
				if category.(map[string]any)["count"].(float64) > 0 {
					switch category.(map[string]any)["header"].(string) {
					case "American":
						has.Western = true
						break
					case "Burgers":
						has.Western = true
						has.FastFood = true
						break
					case "Burritos":
						has.Western = true
						break
					case "Chicken":
						has.FastFood = true
						break
					case "Chinese":
						has.Chinese = true
						break
					case "Dessert":
						has.DessertAndDrinks = true
					case "Dumplings":
						has.Chinese = true
						break
					case "Ice cream":
						has.DessertAndDrinks = true
						break
					case "Italian":
						has.Pizza = true
						break
					case "Indian":
						has.Curry = true
						break
					case "Japanese":
						has.Japanese = true
						break
					case "Mexican":
						has.Western = true
						break
					case "Pizza":
						has.Pizza = true
						break
					case "Sandwiches":
						has.PartyFood = true
						break
					case "Sushi":
						has.Sushi = true
						break
					case "Wine":
						has.DessertAndDrinks = true
						break
					}
				}
			}
		}
	}

	return has, nil
}

var foodTypes = map[CategoryCode][]string{
	Pizza:            {"Pizza", "Italian"},
	Western:          {"American", "Burritos", "Mexican", "Burgers"},
	FastFood:         {"American", "Burgers", "Chicken"},
	Chinese:          {"Chinese", "Dumplings"},
	DrinksAndDessert: {"Dessert", "Ice cream", "Wine"},
	Curry:            {"Indian"},
	Japanese:         {"Japanese"},
	PartyFood:        {"Sandwiches"},
	Sushi:            {"Sushi"},
}

func (d *Deliveroo) GetShops(code CategoryCode) ([]Store, error) {
	graphQLQuery, err := GetShopsQuery(d.longitude, d.latitude)
	if err != nil {
		return nil, err
	}

	response, err := d.sendPOST("https://api.deliveroo.com/consumer/graphql/", graphQLQuery, None)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	respBytes, _ := io.ReadAll(response.Body)

	d.response = string(respBytes)
	d.responseCode = response.StatusCode

	jsonData := map[string]any{}
	err = json.Unmarshal(respBytes, &jsonData)
	if err != nil {
		return nil, err
	}

	restaurants := jsonData["data"].(map[string]any)["results"].(map[string]any)["ui_layout_groups"].([]any)[0].(map[string]any)["ui_layouts"].([]any)

	var shopIDs []string
	for _, restaurant := range restaurants {
		object := restaurant.(map[string]any)["ui_blocks"].([]any)[0].(map[string]any)

		if object["restaurant"] == nil {
			continue
		}

		shopID := object["restaurant"].(map[string]any)["id"].(string)
		if slices.Contains(shopIDs, shopID) {
			continue
		}

		shopIDs = append(shopIDs, shopID)
	}

	wg := sync.WaitGroup{}
	wg.Add(len(shopIDs))

	var stores []Store
	for _, iD := range shopIDs {
		iD := iD
		go func() {
			defer wg.Done()
			_ = iD
			response, err := d.sendGET(fmt.Sprintf("https://api.deliveroo.com/orderapp/v1/restaurants/%s?track=1&lat=%.2f&lng=%.2f&include_unavailable=true&restaurant_fulfillments_supported=true&fulfillment_method=DELIVERY", iD, d.latitude, d.longitude))
			if err != nil {
				return
			}

			defer response.Body.Close()
			respBytes, _ := io.ReadAll(response.Body)

			jsonData := map[string]any{}
			err = json.Unmarshal(respBytes, &jsonData)
			if err != nil {
				return
			}

			if jsonData["category"] == nil {
				return
			}

			if !slices.Contains(foodTypes[code], jsonData["category"].(string)) {
				return
			}

			description := "No Description"
			if jsonData["description"] != nil {
				if jsonData["description"].(string) != "" {
					description = wordwrap.WrapString(jsonData["description"].(string), 38)
					for i3, s := range strings.Split(description, "\n") {
						switch i3 {
						case 0:
							description = s
							break
						case 1:
							description += "\n"
							description += s
							break
						default:
							break
						}
					}
				}
			}

			d.mutex.Lock()
			stores = append(stores, Store{
				Name:         jsonData["name"].(string),
				StoreID:      iD,
				Address:      "",
				WaitTime:     jsonData["total_time"].(float64),
				MinPrice:     jsonData["min_order"].(string),
				IsOpen:       jsonData["open"].(bool),
				DetailedWait: "",
				Phone:        "",
				ServiceHours: ServiceHours{},
				Information:  description,
			})
			d.mutex.Unlock()
		}()
	}

	wg.Wait()

	return stores, nil
}

func (d *Deliveroo) GetStore(id string) (Store, error) {
	response, err := d.sendGET(fmt.Sprintf("https://api.deliveroo.com/orderapp/v1/restaurants/%s?track=1&lat=%.2f&lng=%.2f&include_unavailable=true&restaurant_fulfillments_supported=true&fulfillment_method=DELIVERY", id, d.latitude, d.longitude))
	if err != nil {
		return Store{}, err
	}

	defer response.Body.Close()
	respBytes, _ := io.ReadAll(response.Body)

	d.response = string(respBytes)
	d.responseCode = response.StatusCode

	jsonData := map[string]any{}
	err = json.Unmarshal(respBytes, &jsonData)
	if err != nil {
		return Store{}, err
	}

	description := "No Description"
	if jsonData["description"] != nil {
		description = wordwrap.WrapString(jsonData["description"].(string), 28)
		for i3, s := range strings.Split(description, "\n") {
			switch i3 {
			case 0:
				description = s
				break
			default:
				description += "\n"
				description += s
				break
			}
		}
	}

	amenity := "None"
	if jsonData["promotion_incentive"] != nil {
		// TODO: Figure out other incentives
		if jsonData["promotion_incentive"].(map[string]any)["type"].(string) == "free_delivery" {
			amenity = fmt.Sprintf("Free Delivery at £%.f", jsonData["promotion_incentive"].(map[string]any)["threshold"].(float64))
		}
	}

	return Store{
		Name:     jsonData["name"].(string),
		StoreID:  id,
		Address:  jsonData["address"].(map[string]any)["address1"].(string),
		WaitTime: jsonData["total_time"].(float64),
		MinPrice: jsonData["min_order"].(string),
		// jsonData["open"].(bool)
		IsOpen:       true,
		DetailedWait: "a",
		Phone:        jsonData["phone_number"].(string),
		ServiceHours: ServiceHours{},
		Information:  description,
		Amenity:      amenity,
	}, nil
}

func (d *Deliveroo) GetMenuCategories(id string) ([]Category, error) {
	response, err := d.sendGET(fmt.Sprintf("https://api.deliveroo.com/orderapp/v1/restaurants/%s?track=1&lat=%.2f&lng=%.2f&include_unavailable=true&restaurant_fulfillments_supported=true&fulfillment_method=DELIVERY", id, d.latitude, d.longitude))
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	respBytes, _ := io.ReadAll(response.Body)

	d.response = string(respBytes)
	d.responseCode = response.StatusCode

	jsonData := map[string]any{}
	err = json.Unmarshal(respBytes, &jsonData)
	if err != nil {
		return nil, err
	}

	categoriesJson := jsonData["menu"].(map[string]any)["menu_categories"].([]any)
	var categories []Category
	for _, category := range categoriesJson {
		if !category.(map[string]any)["top_level"].(bool) {
			continue
		}

		categories = append(categories, Category{
			Name: category.(map[string]any)["name"].(string),
			Code: int(category.(map[string]any)["id"].(float64)),
		})
	}

	return categories, nil
}

func (d *Deliveroo) GetItems(shopCode string, categoryID string) ([]Item, error) {
	response, err := d.sendGET(fmt.Sprintf("https://api.deliveroo.com/orderapp/v1/restaurants/%s?track=1&lat=%.2f&lng=%.2f&include_unavailable=true&restaurant_fulfillments_supported=true&fulfillment_method=DELIVERY", shopCode, d.latitude, d.longitude))
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	respBytes, _ := io.ReadAll(response.Body)

	d.response = string(respBytes)
	d.responseCode = response.StatusCode

	jsonData := map[string]any{}
	err = json.Unmarshal(respBytes, &jsonData)
	if err != nil {
		return nil, err
	}

	menuItems := jsonData["menu"].(map[string]any)["menu_items"].([]any)
	var items []Item
	for _, item := range menuItems {
		itemCategoryID := fmt.Sprintf("%.f", item.(map[string]any)["category_id"].(float64))
		if itemCategoryID != categoryID {
			continue
		}

		imageID := fmt.Sprintf("%.f", item.(map[string]any)["id"].(float64))
		if item.(map[string]any)["image_url"] != nil {
			_, err := os.ReadFile(fmt.Sprintf("./images/%s/%s.jpg", shopCode, imageID))
			if err != nil {
				// This is typically unsafe, but we should be able to get away with it for our purposes
				if _, err := os.Stat(fmt.Sprintf("./images/%s", shopCode)); os.IsNotExist(err) {
					os.Mkdir(fmt.Sprintf("./images/%s", shopCode), 0777)
				}

				imageURL := item.(map[string]any)["image_url"].(string)
				imageURL = strings.Replace(imageURL, "{w}", "160", -1)
				imageURL = strings.Replace(imageURL, "{h}", "160", -1)
				imageURL = strings.Replace(imageURL, "{&quality}", "100", -1)
				response, _ = http.Get(imageURL)
				newImage := ConvertImage(response.Body)
				if newImage != nil {
					os.WriteFile(fmt.Sprintf("./images/%s/%s.jpg", shopCode, imageID), newImage, 0666)
					response.Body.Close()
				}
			}
		}

		description := item.(map[string]any)["name"].(string)
		if item.(map[string]any)["description"] != nil {
			description = item.(map[string]any)["description"].(string)
		}

		items = append(items, Item{
			Name:        item.(map[string]any)["name"].(string),
			Description: description,
			ImgID:       imageID,
			SoldOut:     item.(map[string]any)["available"].(bool),
			Price:       item.(map[string]any)["price"].(string),
		})
	}

	return items, nil
}

func (d *Deliveroo) GetItem(shopCode string, categoryCode string, itemCode string) (Item, error) {
	response, err := d.sendGET(fmt.Sprintf("https://api.deliveroo.com/orderapp/v1/restaurants/%s?track=1&lat=%.2f&lng=%.2f&include_unavailable=true&restaurant_fulfillments_supported=true&fulfillment_method=DELIVERY", shopCode, d.latitude, d.longitude))
	if err != nil {
		return Item{}, err
	}

	defer response.Body.Close()
	respBytes, _ := io.ReadAll(response.Body)

	d.response = string(respBytes)
	d.responseCode = response.StatusCode

	jsonData := map[string]any{}
	err = json.Unmarshal(respBytes, &jsonData)
	if err != nil {
		return Item{}, err
	}

	menuItems := jsonData["menu"].(map[string]any)["menu_items"].([]any)
	var item Item
	var modifierGroups []ModifierGroup
	for _, _item := range menuItems {
		itemCategoryID := fmt.Sprintf("%.f", _item.(map[string]any)["category_id"].(float64))
		itemID := fmt.Sprintf("%.f", _item.(map[string]any)["id"].(float64))
		if itemCategoryID != categoryCode || itemCode != itemID {
			continue
		}

		// First we deal with modifiers if possible
		var itemModifierIDs []string
		for _, modifier := range _item.(map[string]any)["modifier_group_ids"].([]any) {
			itemModifierIDs = append(itemModifierIDs, fmt.Sprintf("%.f", modifier.(float64)))
		}

		for _, modifier := range jsonData["menu"].(map[string]any)["menu_modifier_groups"].([]any) {
			if !slices.Contains(itemModifierIDs, fmt.Sprintf("%.f", modifier.(map[string]any)["id"].(float64))) {
				continue
			}

			var modifiers []Modifier
			for _, modifierId := range modifier.(map[string]any)["modifier_item_ids"].([]any) {
				modifiers = append(modifiers, Modifier{ID: fmt.Sprintf("%.f", modifierId.(float64))})
			}

			modifierGroups = append(modifierGroups, ModifierGroup{
				ID:           fmt.Sprintf("%.f", modifier.(map[string]any)["id"].(float64)),
				Name:         modifier.(map[string]any)["name"].(string),
				MinSelection: modifier.(map[string]any)["min_selection_points"].(float64),
				MaxSelection: modifier.(map[string]any)["max_selection_points"].(float64),
				Modifiers:    modifiers,
			})
		}

		item = Item{
			Price:          _item.(map[string]any)["price"].(string),
			ModifierGroups: modifierGroups,
		}
	}

	for i, modifierGroup := range item.ModifierGroups {
		for i2, modifier := range modifierGroup.Modifiers {
			for _, _item := range menuItems {
				if modifier.ID != fmt.Sprintf("%.f", _item.(map[string]any)["id"].(float64)) {
					continue
				}

				imageID := fmt.Sprintf("%.f", _item.(map[string]any)["id"].(float64))
				if _item.(map[string]any)["image_url"] != nil {
					_, err := os.ReadFile(fmt.Sprintf("./images/%s/%s.jpg", shopCode, imageID))
					if err != nil {
						// This is typically unsafe, but we should be able to get away with it for our purposes
						if _, err := os.Stat(fmt.Sprintf("./images/%s", shopCode)); os.IsNotExist(err) {
							os.Mkdir(fmt.Sprintf("./images/%s", shopCode), 0777)
						}

						imageURL := _item.(map[string]any)["image_url"].(string)
						imageURL = strings.Replace(imageURL, "{w}", "160", -1)
						imageURL = strings.Replace(imageURL, "{h}", "160", -1)
						imageURL = strings.Replace(imageURL, "{&quality}", "100", -1)
						response, _ = http.Get(imageURL)
						newImage := ConvertImage(response.Body)
						if newImage != nil {
							os.WriteFile(fmt.Sprintf("./images/%s/%s.jpg", shopCode, imageID), newImage, 0666)
							response.Body.Close()
						}
					}
				} else {
					imageID = "non"
				}

				item.ModifierGroups[i].Modifiers[i2] = Modifier{
					ID:          modifier.ID,
					Name:        _item.(map[string]any)["name"].(string),
					Description: _item.(map[string]any)["name"].(string),
					Price:       _item.(map[string]any)["price"].(string),
					ImageID:     imageID,
				}
			}

		}
	}

	return item, nil
}

// GetItemsWithModifiers returns an array of Item with the selected modifiers
func (d *Deliveroo) GetItemsWithModifiers(shopCode string, itemCodes []string, modifierGroups [][]ModifierGroup) ([]Item, error) {
	response, err := d.sendGET(fmt.Sprintf("https://api.deliveroo.com/orderapp/v1/restaurants/%s?track=1&lat=%.2f&lng=%.2f&include_unavailable=true&restaurant_fulfillments_supported=true&fulfillment_method=DELIVERY", shopCode, d.latitude, d.longitude))
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	respBytes, _ := io.ReadAll(response.Body)

	d.response = string(respBytes)
	d.responseCode = response.StatusCode

	jsonData := map[string]any{}
	err = json.Unmarshal(respBytes, &jsonData)
	if err != nil {
		return nil, err
	}

	menuItems := jsonData["menu"].(map[string]any)["menu_items"].([]any)
	var items []Item

	for i, itemCode := range itemCodes {
		for _, _item := range menuItems {
			itemID := fmt.Sprintf("%.f", _item.(map[string]any)["id"].(float64))
			if itemCode == itemID {

				item := Item{
					Name:           _item.(map[string]any)["name"].(string),
					Price:          _item.(map[string]any)["price"].(string),
					ModifierGroups: modifierGroups[i],
				}

				for i, group := range item.ModifierGroups {
					for _, modifier := range jsonData["menu"].(map[string]any)["menu_modifier_groups"].([]any) {
						if group.ID == fmt.Sprintf("%.f", modifier.(map[string]any)["id"].(float64)) {
							item.ModifierGroups[i].Name = modifier.(map[string]any)["name"].(string)
							break
						}
					}
				}

				for i2, modifierGroup := range modifierGroups[i] {
					for i3, modifier := range modifierGroup.Modifiers {
						for _, __item := range menuItems {
							if modifier.ID != fmt.Sprintf("%.f", __item.(map[string]any)["id"].(float64)) {
								continue
							}

							item.ModifierGroups[i2].Modifiers[i3] = Modifier{
								ID:          modifier.ID,
								Name:        __item.(map[string]any)["name"].(string),
								Description: __item.(map[string]any)["name"].(string),
								Price:       __item.(map[string]any)["price"].(string),
							}
						}

					}
				}

				items = append(items, item)
				break
			}
		}
	}

	return items, nil
}

func (d *Deliveroo) SendBasket(shopCode string, basket []map[string]any) (string, string, float64, error) {
	payload := map[string]any{
		"basket": map[string]any{
			"allergy_note":               "",
			"driver_tip":                 0,
			"fulfillment_method":         "DELIVERY",
			"items":                      basket,
			"order_modifiers_collection": map[string]any{},
			"restaurant_id":              shopCode,
			"scheduled_delivery_day":     "today",
			"scheduled_delivery_time":    "ASAP",
		},
		"basket_upsell_version": 0,
		"corporate":             false,
		"deliver_to": []float64{
			d.longitude,
			d.latitude,
		},
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return "", "", 0, err
	}

	response, err := d.sendPOST("https://api.deliveroo.com/orderapp/v1/basket", data, None)
	if err != nil {
		return "", "", 0, err
	}

	defer response.Body.Close()
	respBytes, _ := io.ReadAll(response.Body)

	d.response = string(respBytes)
	d.responseCode = response.StatusCode

	jsonData := map[string]any{}
	err = json.Unmarshal(respBytes, &jsonData)
	if err != nil {
		return "", "", 0, err
	}

	fee, _ := strconv.ParseFloat(jsonData["basket"].(map[string]any)["fee"].(string), 64)
	surcharge, _ := strconv.ParseFloat(jsonData["basket"].(map[string]any)["surcharge"].(string), 64)

	return jsonData["basket"].(map[string]any)["total"].(string), jsonData["basket"].(map[string]any)["subtotal"].(string), surcharge + fee, nil
}

func (d *Deliveroo) CreatePaymentPlan() (*PaymentMethod, error) {
	query, err := GetCreatePaymentQuery()
	if err != nil {
		return nil, err
	}

	response, err := d.sendPOST("https://api.deliveroo.com/checkout-api/graphql-query", query, CreatePaymentPlan)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	respBytes, _ := io.ReadAll(response.Body)

	d.response = string(respBytes)
	d.responseCode = response.StatusCode

	if response.StatusCode != http.StatusOK {
		return nil, errors.New("failed to create payment plan")
	}

	jsonData := map[string]any{}
	err = json.Unmarshal(respBytes, &jsonData)
	if err != nil {
		return nil, err
	}

	creditCard := "In App Credit"
	if jsonData["data"].(map[string]any)["payment_plan"].(map[string]any)["payment_options"].(map[string]any)["selected_completing"] != nil {
		creditCard = jsonData["data"].(map[string]any)["payment_plan"].(map[string]any)["payment_options"].(map[string]any)["selected_completing"].(map[string]any)["description"].([]any)[0].(string)
	}

	return &PaymentMethod{
		RestaurantName:  jsonData["data"].(map[string]any)["payment_plan"].(map[string]any)["fulfillment_details"].(map[string]any)["restaurant"].(string),
		ID:              jsonData["data"].(map[string]any)["payment_plan"].(map[string]any)["id"].(string),
		DeliveryAddress: jsonData["data"].(map[string]any)["payment_plan"].(map[string]any)["delivery_addresses"].(map[string]any)["selected"].(map[string]any)["short_description"].([]any)[0].(string),
		CreditCard:      creditCard,
	}, nil
}

func ConvertImage(data io.Reader) []byte {
	origImage, err := jpeg.Decode(data)
	if err != nil {
		return nil
	}

	newImage := image.NewRGBA(image.Rect(0, 0, 160, 160))
	draw.BiLinear.Scale(newImage, newImage.Bounds(), origImage, origImage.Bounds(), draw.Over, nil)

	var outputImgWriter bytes.Buffer
	err = jpeg.Encode(bufio.NewWriter(&outputImgWriter), newImage, nil)
	if err != nil {
		return nil
	}

	return outputImgWriter.Bytes()
}

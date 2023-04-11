package main

import (
	"DemaeDeliveroo/deliveroo"
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/mitchellh/go-wordwrap"
	"io"
	"net/http"
	"strings"
	"time"
)

const InsertAuthkey = `UPDATE "user" SET auth_key = $1 WHERE wii_id = $2`
const InsertPaymentID = `UPDATE "user" SET payment_id = $1 WHERE wii_id = $2`
const QueryUserBasket = `SELECT "user".basket, "user".auth_key FROM "user" WHERE "user".wii_id = $1 LIMIT 1`
const QueryUser = `SELECT "user".basket, "user".auth_key, "user".discord_id FROM "user" WHERE "user".wii_id = $1 LIMIT 1`
const ClearBasket = `UPDATE "user" SET order_id = $1, price = $2, basket = $3 WHERE wii_id = $4`

func authKey(r *Response) {
	authKeyValue, err := uuid.DefaultGenerator.NewV1()
	if err != nil {
		r.ReportError(err, http.StatusUnauthorized)
		return
	}

	// First we query to determine if the user already has an auth key. If they do, reset the basket.
	var _authKey string
	row := pool.QueryRow(context.Background(), QueryUserBasket, r.request.Header.Get("X-WiiID"))
	err = row.Scan(nil, &_authKey)
	if err != nil {
		r.ReportError(err, http.StatusInternalServerError)
		return
	}

	if _authKey != "" {
		_, err = pool.Exec(context.Background(), ClearBasket, "", "", "[]", r.request.Header.Get("X-WiiID"))
		if err != nil {
			r.ReportError(err, http.StatusInternalServerError)
			return
		}
	}

	_, err = pool.Exec(context.Background(), InsertAuthkey, authKeyValue.String(), r.request.Header.Get("X-WiiID"))
	if err != nil {
		r.ReportError(err, http.StatusInternalServerError)
		return
	}

	r.ResponseFields = []any{
		KVField{
			XMLName: xml.Name{Local: "authKey"},
			Value:   authKeyValue.String(),
		},
	}
}

func basketAdd(r *Response) {
	itemCode := r.request.PostForm.Get("itemCode")
	quantity := r.request.PostForm.Get("quantity")

	modifiers := parseOptions(r.request.PostForm)
	basket := BasketJSON{
		ItemCode:  itemCode,
		Quantity:  quantity,
		Modifiers: modifiers,
	}

	if len(modifiers) == 0 {
		basket.Modifiers = []ModifierJSON{}
	}

	var lastBasket string
	row := pool.QueryRow(context.Background(), QueryUserBasket, r.request.Header.Get("X-WiiID"))
	err := row.Scan(&lastBasket, nil)
	if err != nil {
		r.ReportError(err, http.StatusInternalServerError)
		return
	}

	var actualBasket []map[string]any
	err = json.Unmarshal([]byte(lastBasket), &actualBasket)
	if err != nil {
		r.ReportError(err, http.StatusInternalServerError)
		return
	}

	// This is literal insanity but go with it
	data, err := json.Marshal(basket)
	if err != nil {
		r.ReportError(err, http.StatusInternalServerError)
		return
	}

	var newBasket map[string]any
	err = json.Unmarshal(data, &newBasket)
	if err != nil {
		r.ReportError(err, http.StatusInternalServerError)
		return
	}

	actualBasket = append(actualBasket, newBasket)

	// Convert basket to JSON then insert to database
	jsonStr, err := json.Marshal(actualBasket)
	if err != nil {
		r.ReportError(err, http.StatusInternalServerError)
		return
	}

	_, err = pool.Exec(context.Background(), `UPDATE "user" SET basket = $1 WHERE wii_id = $2`, jsonStr, r.request.Header.Get("X-WiiID"))
	if err != nil {
		r.ReportError(err, http.StatusInternalServerError)
		return
	}
}

func basketList(r *Response) {
	var basketStr string
	row := pool.QueryRow(context.Background(), QueryUserBasket, r.request.Header.Get("X-WiiID"))
	err := row.Scan(&basketStr, nil)
	if err != nil {
		r.ReportError(err, http.StatusInternalServerError)
		return
	}

	var basket []BasketJSON
	err = json.Unmarshal([]byte(basketStr), &basket)
	if err != nil {
		r.ReportError(err, http.StatusInternalServerError)
		return
	}

	var mapBasket []map[string]any
	err = json.Unmarshal([]byte(basketStr), &mapBasket)
	if err != nil {
		r.ReportError(err, http.StatusInternalServerError)
		return
	}

	d, err := deliveroo.NewDeliveroo(pool, r.request)
	if err != nil {
		r.ReportError(err, http.StatusInternalServerError)
		return
	}

	total, subtotal, charge, err, resp := d.SendBasket(r.request.URL.Query().Get("shopCode"), mapBasket)
	if err != nil {
		r.ReportError(fmt.Errorf("%v\nResponse: %s", err, resp), http.StatusInternalServerError)
		return
	}

	basketPrice := KVField{
		XMLName: xml.Name{Local: "basketPrice"},
		Value:   charge,
	}

	chargePrice := KVField{
		XMLName: xml.Name{Local: "chargePrice"},
		Value:   subtotal,
	}

	totalPrice := KVField{
		XMLName: xml.Name{Local: "totalPrice"},
		Value:   total,
	}

	status := KVFieldWChildren{
		XMLName: xml.Name{Local: "Status"},
		Value: []any{
			KVField{
				XMLName: xml.Name{Local: "isOrder"},
				Value:   BoolToInt(true),
			},
			KVFieldWChildren{
				XMLName: xml.Name{Local: "messages"},
				Value: []any{KVField{
					XMLName: xml.Name{Local: "hey"},
					Value:   "how are you?",
				}},
			},
		},
	}

	var itemCodes []string
	var modifierGroups [][]deliveroo.ModifierGroup

	for _, basketJSON := range basket {
		itemCodes = append(itemCodes, basketJSON.ItemCode)

		var groups []deliveroo.ModifierGroup
		for _, modifier := range basketJSON.Modifiers {
			group := deliveroo.ModifierGroup{
				ID:           modifier.ModifierGroupID,
				Name:         "",
				MinSelection: 0,
				MaxSelection: 0,
				Modifiers:    nil,
			}

			var modifiers []deliveroo.Modifier
			for _, s := range modifier.ModifierID {
				modifiers = append(modifiers, deliveroo.Modifier{
					ID:          s,
					Name:        "",
					Description: "",
					Price:       "",
				})
			}

			group.Modifiers = modifiers
			groups = append(groups, group)
		}
		modifierGroups = append(modifierGroups, groups)
	}

	items, err, resp := d.GetItemsWithModifiers(r.request.URL.Query().Get("shopCode"), itemCodes, modifierGroups)
	if err != nil {
		r.ReportError(fmt.Errorf("%v\nResponse: %s", err, resp), http.StatusInternalServerError)
		return
	}

	var basketItems []BasketItem
	for i, item := range items {
		var modifierGroups []ItemOne
		for i2, group := range item.ModifierGroups {
			options := ItemOne{
				XMLName: xml.Name{Local: fmt.Sprintf("container%d", i2)},
				Info:    CDATA{""},
				Code:    CDATA{0},
				Type:    CDATA{0},
				Name:    CDATA{group.Name},
				List:    KVFieldWChildren{},
			}

			for _, option := range group.Modifiers {
				options.List.Value = append(options.List.Value, Item{
					MenuCode:   CDATA{0},
					ItemCode:   CDATA{0},
					Name:       CDATA{option.Name},
					Price:      CDATA{option.Price},
					Info:       CDATA{option.Description},
					IsSelected: &CDATA{BoolToInt(true)},
					Image:      CDATA{"non"},
					IsSoldout:  CDATA{BoolToInt(false)},
				})
			}

			modifierGroups = append(modifierGroups, options)
		}

		name := wordwrap.WrapString(item.Name, 29)
		for i3, s := range strings.Split(name, "\n") {
			switch i3 {
			case 0:
				name = s
				break
			case 1:
				name += "\n"
				name += s
				break
			default:
				// If it is too long it becomes ... so we are fine
				name += " " + s
				break
			}
		}

		basketItems = append(basketItems, BasketItem{
			XMLName:       xml.Name{Local: fmt.Sprintf("container%d", i)},
			BasketNo:      CDATA{i + 1},
			MenuCode:      CDATA{1},
			ItemCode:      CDATA{0},
			Name:          CDATA{name},
			Price:         CDATA{item.Price},
			Size:          CDATA{""},
			IsSoldout:     CDATA{BoolToInt(false)},
			Quantity:      CDATA{1},
			SubTotalPrice: CDATA{item.Price},
			Menu: KVFieldWChildren{
				XMLName: xml.Name{Local: "Menu"},
				Value: []any{
					KVField{
						XMLName: xml.Name{Local: "name"},
						Value:   "Menu",
					},
					KVFieldWChildren{
						XMLName: xml.Name{Local: "lunchMenuList"},
						Value: []any{
							KVField{
								XMLName: xml.Name{Local: "isLunchTimeMenu"},
								Value:   BoolToInt(false),
							},
							KVField{
								XMLName: xml.Name{Local: "isOpen"},
								Value:   BoolToInt(true),
							},
						},
					},
				},
			},
			OptionList: KVFieldWChildren{
				XMLName: xml.Name{Local: ""},
				Value: []any{
					modifierGroups,
				},
			},
		})
	}

	cart := KVFieldWChildren{
		XMLName: xml.Name{Local: "List"},
		Value:   []any{basketItems[:]},
	}

	r.ResponseFields = []any{
		basketPrice,
		chargePrice,
		totalPrice,
		status,
		cart,
	}
}

func orderDone(r *Response) {
	// We don't handle the actual order here, we dispatch a verification message to the current user.
	// The Discord bot then places the order if all is right.
	var basketStr string
	var _authKey string
	var discordID string
	row := pool.QueryRow(context.Background(), QueryUser, r.request.Header.Get("X-WiiID"))
	err := row.Scan(&basketStr, &_authKey, &discordID)
	if err != nil {
		r.ReportError(err, http.StatusInternalServerError)
		return
	}

	payload := map[string]string{
		"recipient_id": discordID,
	}

	data, _ := json.Marshal(payload)

	client := &http.Client{}
	req, _ := http.NewRequest("POST", "https://discord.com/api/v9/users/@me/channels", bytes.NewBuffer(data))
	req.Header.Add("Authorization", "Bot MTA4NDk1Mjk1NzQ1MzM1NzEwOA.GrxKnX.2gzsdeDLdInbhtB1UWY-zPLiHfVPQf5wlpsb_U")
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		r.ReportError(err, http.StatusInternalServerError)
		return
	}

	respBytes, _ := io.ReadAll(resp.Body)
	var dm struct {
		ID string `json:"id"`
	}

	err = json.Unmarshal(respBytes, &dm)
	if err != nil {
		r.ReportError(err, http.StatusInternalServerError)
		return
	}

	d, err := deliveroo.NewDeliveroo(pool, r.request)
	if err != nil {
		r.ReportError(err, http.StatusInternalServerError)
		return
	}

	var basket []BasketJSON
	err = json.Unmarshal([]byte(basketStr), &basket)
	if err != nil {
		r.ReportError(err, http.StatusInternalServerError)
		return
	}

	var mapBasket []map[string]any
	err = json.Unmarshal([]byte(basketStr), &mapBasket)
	if err != nil {
		r.ReportError(err, http.StatusInternalServerError)
		return
	}

	var itemCodes []string
	var modifierGroups [][]deliveroo.ModifierGroup

	for _, basketJSON := range basket {
		itemCodes = append(itemCodes, basketJSON.ItemCode)

		var groups []deliveroo.ModifierGroup
		for _, modifier := range basketJSON.Modifiers {
			group := deliveroo.ModifierGroup{
				ID:           modifier.ModifierGroupID,
				Name:         "",
				MinSelection: 0,
				MaxSelection: 0,
				Modifiers:    nil,
			}

			var modifiers []deliveroo.Modifier
			for _, s := range modifier.ModifierID {
				modifiers = append(modifiers, deliveroo.Modifier{
					ID:          s,
					Name:        "",
					Description: "",
					Price:       "",
				})
			}

			group.Modifiers = modifiers
			groups = append(groups, group)
		}
		modifierGroups = append(modifierGroups, groups)
	}

	items, err, res := d.GetItemsWithModifiers(r.request.PostForm.Get("shop[ShopCode]"), itemCodes, modifierGroups)
	if err != nil {
		r.ReportError(fmt.Errorf("%v\nResponse: %s", err, res), http.StatusInternalServerError)
		return
	}

	itemsStr := ""
	for _, item := range items {
		itemsStr += fmt.Sprintf("%s - %s\n", item.Name, item.Price)
	}

	total, _, _, err, res := d.SendBasket(r.request.PostForm.Get("shop[ShopCode]"), mapBasket)
	if err != nil {
		r.ReportError(fmt.Errorf("%v\nResponse: %s", err, res), http.StatusInternalServerError)
		return
	}

	payment, err, res := d.CreatePaymentPlan()
	if err != nil {
		r.ReportError(fmt.Errorf("%v\nResponse: %s", err, resp), http.StatusInternalServerError)
		return
	}

	_, err = pool.Exec(context.Background(), InsertPaymentID, payment.ID, r.request.Header.Get("X-WiiID"))
	if err != nil {
		r.ReportError(err, http.StatusInternalServerError)
		return
	}

	message := map[string]any{
		"content": nil,
		"embeds": []map[string]any{
			{
				"title":       fmt.Sprintf("Deliveroo Order Request - %s", payment.RestaurantName),
				"description": fmt.Sprintf("Demae Deliveroo has recieved a request for a %s order with the following basket:.\n\n%s", total, itemsStr),
				"color":       3666886,
				"fields": []map[string]any{
					{
						"name":  "Delivery Address",
						"value": payment.DeliveryAddress,
					},
					{
						"name":  "Selected Credit Card",
						"value": payment.CreditCard,
					},
					{
						"name":  "Accept Request",
						"value": fmt.Sprintf("To accept the request, enter the following:\n```I agree to placing this order. I acknowledge that once this goes through, there is no cancelling the order. Auth Key: %s```", _authKey),
					},
				},
			},
		},
	}

	data, _ = json.Marshal(message)
	req, _ = http.NewRequest("POST", fmt.Sprintf("https://discord.com/api/v9/channels/%s/messages", dm.ID), bytes.NewBuffer(data))
	req.Header.Add("Authorization", "Bot MTA4NDk1Mjk1NzQ1MzM1NzEwOA.GrxKnX.2gzsdeDLdInbhtB1UWY-zPLiHfVPQf5wlpsb_U")
	req.Header.Add("Content-Type", "application/json")

	resp, err = client.Do(req)
	if err != nil {
		r.ReportError(err, http.StatusInternalServerError)
		return
	}

	PostDiscordWebhook(
		"An order has been queued. Awaiting confirmation",
		fmt.Sprintf("The order was placed by user id %s", r.request.Header.Get("X-WiiID")),
		config.OrderWebhook,
		65311,
	)

	currentTime := time.Now().Format("200602011504")
	r.AddKVWChildNode("Message", KVField{
		XMLName: xml.Name{Local: "contents"},
		Value:   "Thank you! Your order has been placed!",
	})
	r.AddKVNode("order_id", "1")
	r.AddKVNode("orderDay", currentTime)
	r.AddKVNode("hashKey", "Testing: 1, 2, 3")
	r.AddKVNode("hour", currentTime)
}

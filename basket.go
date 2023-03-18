package main

import (
	"DemaeDeliveroo/deliveroo"
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/mitchellh/go-wordwrap"
	"net/http"
	"strings"
)

const InsertAuthkey = `UPDATE "user" SET auth_key = $1 WHERE wii_id = $2`
const QueryUserBasket = `SELECT "user".basket, "user".auth_key FROM "user" WHERE "user".wii_id = $1 LIMIT 1`
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

	total, subtotal, charge, err := d.SendBasket(r.request.URL.Query().Get("shopCode"), mapBasket)
	if err != nil {
		r.ReportError(err, http.StatusInternalServerError)
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

	items, err := d.GetItemsWithModifiers(r.request.URL.Query().Get("shopCode"), itemCodes, modifierGroups)
	if err != nil {
		r.ReportError(err, http.StatusInternalServerError)
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

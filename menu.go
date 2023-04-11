package main

import (
	"DemaeDeliveroo/deliveroo"
	"encoding/xml"
	"fmt"
	"github.com/mitchellh/go-wordwrap"
	"net/http"
	"strings"
)

func menuList(r *Response) {
	d, err := deliveroo.NewDeliveroo(pool, r.request)
	if err != nil {
		r.ReportError(err, http.StatusInternalServerError)
		return
	}

	menuData, err, resp := d.GetMenuCategories(r.request.URL.Query().Get("shopCode"))
	if err != nil {
		r.ReportError(fmt.Errorf("%v\nResponse: %s", err, resp), http.StatusInternalServerError)
		return
	}

	var menus []Menu
	for _, menu := range menuData {
		menus = append(menus, Menu{
			XMLName:     xml.Name{Local: fmt.Sprintf("container_%s", menu.Code)},
			MenuCode:    CDATA{menu.Code},
			LinkTitle:   CDATA{menu.Name},
			EnabledLink: CDATA{1},
			Name:        CDATA{menu.Name},
			Info:        CDATA{menu.Name},
			SetNum:      CDATA{1},
			LunchMenuList: struct {
				IsLunchTimeMenu CDATA `xml:"isLunchTimeMenu"`
				Hour            KVFieldWChildren
				IsOpen          CDATA `xml:"isOpen"`
				Message         CDATA `xml:"message"`
			}{
				IsLunchTimeMenu: CDATA{BoolToInt(false)},
				Hour: KVFieldWChildren{
					XMLName: xml.Name{Local: "hour"},
					Value: []any{
						KVField{
							XMLName: xml.Name{Local: "start"},
							Value:   "00:00:00",
						},
						KVField{
							XMLName: xml.Name{Local: "end"},
							Value:   "24:59:59",
						},
					},
				},
				IsOpen:  CDATA{BoolToInt(true)},
				Message: CDATA{"Where does this show up?"},
			},
		})
	}

	// Append 1 more as a placeholder
	placeholder := menus[0]
	placeholder.XMLName = xml.Name{Local: "placeholder"}
	menus = append(menus, placeholder)
	r.AddCustomType(menus)
}

func itemList(r *Response) {
	var items []NestedItem
	d, err := deliveroo.NewDeliveroo(pool, r.request)
	if err != nil {
		r.ReportError(err, http.StatusInternalServerError)
		return
	}

	itemData, err, resp := d.GetItems(r.request.URL.Query().Get("shopCode"), r.request.URL.Query().Get("menuCode"))
	if err != nil {
		r.ReportError(fmt.Errorf("%v\nResponse: %s", err, resp), http.StatusInternalServerError)
		return
	}

	for i, item := range itemData {
		name := wordwrap.WrapString(item.Name, 20)
		for i, s := range strings.Split(name, "\n") {
			switch i {
			case 0:
				name = s
				break
			default:
				name += "\n"
				name += s
				break
			}
		}

		description := wordwrap.WrapString(item.Description, 36)
		for i, s := range strings.Split(description, "\n") {
			switch i {
			case 0:
				description = s
				break
			case 1:
				description += "\n"
				description += s
				break
			case 2:
				description += "\n"
				description += strings.Split(wordwrap.WrapString(s, 22), "\n")[0]
				break
			default:
				break
			}
		}

		nestedItem := NestedItem{
			XMLName: xml.Name{Local: fmt.Sprintf("container%d", i)},
			Name:    CDATA{name},
			Item: Item{
				XMLName:   xml.Name{Local: "item"},
				MenuCode:  CDATA{r.request.URL.Query().Get("menuCode")},
				ItemCode:  CDATA{item.ImgID},
				Price:     CDATA{"vee"},
				Info:      CDATA{description},
				Size:      &CDATA{"something"},
				Image:     CDATA{item.ImgID},
				IsSoldout: CDATA{BoolToInt(item.SoldOut)},
				SizeList: &KVFieldWChildren{
					XMLName: xml.Name{Local: "sizeList"},
					Value: []any{
						ItemSize{
							XMLName:   xml.Name{Local: "item0"},
							ItemCode:  CDATA{item.ImgID},
							Size:      CDATA{name},
							Price:     CDATA{item.Price},
							IsSoldout: CDATA{BoolToInt(false)},
						},
					},
				},
			},
		}

		items = append(items, nestedItem)
	}

	r.ResponseFields = []any{
		KVField{
			XMLName: xml.Name{Local: "Count"},
			Value:   len(items),
		},
		KVFieldWChildren{
			XMLName: xml.Name{Local: "List"},
			Value:   []any{items[:]},
		},
	}
}

func itemOne(r *Response) {
	d, err := deliveroo.NewDeliveroo(pool, r.request)
	if err != nil {
		r.ReportError(err, http.StatusInternalServerError)
		return
	}

	item, err, resp := d.GetItem(r.request.URL.Query().Get("shopCode"), r.request.URL.Query().Get("menuCode"), r.request.URL.Query().Get("itemCode"))
	if err != nil {
		r.ReportError(fmt.Errorf("%v\nResponse: %s", err, resp), http.StatusInternalServerError)
		return
	}

	var modifiers []ItemOne
	for i, modifier := range item.ModifierGroups {
		buttonType := "box"
		if modifier.MaxSelection == modifier.MinSelection {
			buttonType = "radio"
		}

		parent := ItemOne{
			XMLName: xml.Name{Local: fmt.Sprintf("container%d", i)},
			Info:    CDATA{fmt.Sprintf("Max item selection is %.f, minimum %.f", modifier.MaxSelection, modifier.MinSelection)},
			Code:    CDATA{modifier.ID},
			Type:    CDATA{buttonType},
			Name:    CDATA{modifier.Name},
			List: KVFieldWChildren{
				XMLName: xml.Name{Local: "list"},
			},
		}

		for _, _item := range modifier.Modifiers {
			name := wordwrap.WrapString(_item.Name, 32)
			for i3, s := range strings.Split(name, "\n") {
				switch i3 {
				case 0:
					name = s
					break
				default:
					name += "\n"
					name += s
					break
				}
			}

			fmt.Println(_item.ImageID)
			parent.List.Value = append(parent.List.Value, Item{
				MenuCode:  CDATA{modifier.ID},
				ItemCode:  CDATA{_item.ID},
				Name:      CDATA{name},
				Price:     CDATA{_item.Price},
				Info:      CDATA{""},
				Size:      nil,
				Image:     CDATA{_item.ImageID},
				IsSoldout: CDATA{BoolToInt(false)},
				SizeList:  nil,
			})
		}

		modifiers = append(modifiers, parent)
	}

	r.ResponseFields = []any{
		KVField{
			XMLName: xml.Name{Local: "price"},
			Value:   item.Price,
		},
		KVFieldWChildren{
			XMLName: xml.Name{Local: "optionList"},
			Value:   []any{modifiers[:]},
		},
	}
}

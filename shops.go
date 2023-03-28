package main

import (
	"DemaeDeliveroo/deliveroo"
	"encoding/xml"
	"net/http"
)

func shopList(r *Response) {
	categoryCode := r.request.URL.Query().Get("categoryCode")

	d, err := deliveroo.NewDeliveroo(pool, r.request)
	if err != nil {
		r.ReportError(err, http.StatusInternalServerError)
		return
	}

	stores, err := d.GetShops(deliveroo.CategoryCode(categoryCode))
	if err != nil {
		r.ReportError(err, http.StatusInternalServerError)
		return
	}

	var storesXML []BasicShop
	for _, storeData := range stores {
		store := BasicShop{
			ShopCode:    CDATA{storeData.StoreID},
			HomeCode:    CDATA{storeData.StoreID},
			Name:        CDATA{storeData.Name},
			Catchphrase: CDATA{storeData.Information},
			MinPrice:    CDATA{storeData.MinPrice},
			Yoyaku:      CDATA{1},
			Activate:    CDATA{"on"},
			WaitTime:    CDATA{storeData.WaitTime},
			PaymentList: KVFieldWChildren{
				XMLName: xml.Name{Local: "paymentList"},
				Value: []any{
					KVField{
						XMLName: xml.Name{Local: "athing"},
						Value:   "Fox Card",
					},
				},
			},
			ShopStatus: KVFieldWChildren{
				XMLName: xml.Name{Local: "shopStatus"},
				Value: []any{
					KVFieldWChildren{
						XMLName: xml.Name{Local: "status"},
						Value: []any{
							KVField{
								XMLName: xml.Name{Local: "isOpen"},
								Value:   1,
							},
						},
					},
				},
			},
		}

		storesXML = append(storesXML, store)
	}

	category := KVFieldWChildren{
		XMLName: xml.Name{Local: "Pizza"},
		Value: []any{
			KVField{
				XMLName: xml.Name{Local: "LargeCategoryName"},
				Value:   "Meal",
			},
			KVFieldWChildren{
				XMLName: xml.Name{Local: "CategoryList"},
				Value: []any{
					KVFieldWChildren{
						XMLName: xml.Name{Local: "TestingCategory"},
						Value: []any{
							KVField{
								XMLName: xml.Name{Local: "CategoryCode"},
								Value:   categoryCode,
							},
							KVFieldWChildren{
								XMLName: xml.Name{Local: "ShopList"},
								Value: []any{
									storesXML,
								},
							},
						},
					},
				},
			},
		},
	}

	r.AddCustomType(category)
}

func shopOne(r *Response) {
	shopCode := r.request.URL.Query().Get("shopCode")

	d, err := deliveroo.NewDeliveroo(pool, r.request)
	if err != nil {
		r.ReportError(err, http.StatusInternalServerError)
		return
	}

	shopData, _ := d.GetStore(shopCode)

	shop := ShopOne{
		CategoryCode:  CDATA{"01"},
		Address:       CDATA{shopData.Address},
		Information:   CDATA{shopData.Information},
		Attention:     CDATA{"why"},
		Amenity:       CDATA{shopData.Amenity},
		MenuListCode:  CDATA{1},
		Activate:      CDATA{"on"},
		WaitTime:      CDATA{shopData.WaitTime},
		TimeOrder:     CDATA{"y"},
		Tel:           CDATA{shopData.Phone},
		YoyakuMinDate: CDATA{1},
		YoyakuMaxDate: CDATA{30},
		PaymentList: KVFieldWChildren{
			XMLName: xml.Name{Local: "paymentList"},
			Value: []any{
				KVField{
					XMLName: xml.Name{Local: "athing"},
					Value:   "Fox Card",
				},
			},
		},
		ShopStatus: ShopStatus{
			Hours: KVFieldWChildren{
				XMLName: xml.Name{Local: "hours"},
				Value: []any{
					KVFieldWChildren{
						XMLName: xml.Name{Local: "all"},
						Value: []any{
							KVField{
								XMLName: xml.Name{Local: "message"},
								Value:   shopData.DetailedWait,
							},
						},
					},
					KVFieldWChildren{
						XMLName: xml.Name{Local: "today"},
						Value: []any{
							KVFieldWChildren{
								XMLName: xml.Name{Local: "values"},
								Value: []any{
									KVField{
										XMLName: xml.Name{Local: "start"},
										Value:   "01:00:00",
									},
									KVField{
										XMLName: xml.Name{Local: "end"},
										Value:   "23:45:00",
									},
									KVField{
										XMLName: xml.Name{Local: "holiday"},
										Value:   "n",
									},
								},
							},
							KVFieldWChildren{
								XMLName: xml.Name{Local: "values1"},
								Value: []any{
									KVField{
										XMLName: xml.Name{Local: "start"},
										Value:   "01:00:00",
									},
									KVField{
										XMLName: xml.Name{Local: "end"},
										Value:   "23:45:00",
									},
									KVField{
										XMLName: xml.Name{Local: "holiday"},
										Value:   "n",
									},
								},
							},
						},
					},
					KVFieldWChildren{
						XMLName: xml.Name{Local: "delivery"},
						Value: []any{
							KVFieldWChildren{
								XMLName: xml.Name{Local: "values"},
								Value: []any{
									KVField{
										XMLName: xml.Name{Local: "start"},
										Value:   "01:00:00",
									},
									KVField{
										XMLName: xml.Name{Local: "end"},
										Value:   "23:45:00",
									},
									KVField{
										XMLName: xml.Name{Local: "holiday"},
										Value:   "n",
									},
								},
							},
						},
					},
					KVFieldWChildren{
						XMLName: xml.Name{Local: "selList"},
						Value: []any{
							KVFieldWChildren{
								XMLName: xml.Name{Local: "values"},
								Value: []any{
									KVFieldWChildren{
										XMLName: xml.Name{Local: "one"},
										Value: []any{
											KVField{
												XMLName: xml.Name{Local: "id"},
												Value:   "1",
											},
											KVField{
												XMLName: xml.Name{Local: "name"},
												Value:   "n",
											},
										},
									},
									KVFieldWChildren{
										XMLName: xml.Name{Local: "two"},
										Value: []any{
											KVField{
												XMLName: xml.Name{Local: "id"},
												Value:   "2",
											},
											KVField{
												XMLName: xml.Name{Local: "name"},
												Value:   "n",
											},
										},
									},
								},
							},
						},
					},
					KVFieldWChildren{
						XMLName: xml.Name{Local: "status"},
						Value: []any{
							KVField{
								XMLName: xml.Name{Local: "isOpen"},
								Value:   BoolToInt(shopData.IsOpen),
							},
						},
					},
				},
			},
			Interval: CDATA{5},
			Holiday:  CDATA{"No ordering on Canada Day"},
		},
		RecommendedItemList: KVFieldWChildren{
			Value: []any{
				Item{
					XMLName:   xml.Name{Local: "container1"},
					MenuCode:  CDATA{10},
					ItemCode:  CDATA{1},
					Name:      CDATA{"Pizza"},
					Price:     CDATA{10},
					Info:      CDATA{"Fresh"},
					Size:      &CDATA{1},
					Image:     CDATA{"PIZZA"},
					IsSoldout: CDATA{0},
					SizeList: &KVFieldWChildren{
						XMLName: xml.Name{Local: "sizeList"},
						Value: []any{
							ItemSize{
								XMLName:   xml.Name{Local: "item1"},
								ItemCode:  CDATA{1},
								Size:      CDATA{1},
								Price:     CDATA{10},
								IsSoldout: CDATA{0},
							},
						},
					},
				},
			},
		},
	}

	// Strip the parent response tag
	r.ResponseFields = shop
}

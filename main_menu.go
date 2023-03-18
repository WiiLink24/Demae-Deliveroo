package main

import (
	"DemaeDeliveroo/deliveroo"
	"encoding/xml"
	"net/http"
)

func documentTemplate(r *Response) {
	r.AddKVWChildNode("container0", KVField{
		XMLName: xml.Name{Local: "contents"},
		Value:   "By clicking agree, you verify you have read and agree to https://demae.wiilink24.com/privacypolicy and https://demae.wiilink24.com/tos",
	})
	r.AddKVWChildNode("container1", KVField{
		XMLName: xml.Name{Local: "contents"},
		Value:   "Among Us",
	})
	r.AddKVWChildNode("container2", KVField{
		XMLName: xml.Name{Local: "contents"},
		Value:   "Among Us",
	})
}

func categoryList(r *Response) {
	d, err := deliveroo.NewDeliveroo(pool, r.request)
	if err != nil {
		r.ReportError(err, http.StatusInternalServerError)
		return
	}

	has, err := d.GetBareShops()
	if err != nil {
		r.ReportError(err, http.StatusInternalServerError)
		return
	}

	if has.Pizza {
		r.MakeCategoryXML(deliveroo.Pizza)
	}
	if has.Bento {
		r.MakeCategoryXML(deliveroo.BentoBox)
	}
	if has.Sushi {
		r.MakeCategoryXML(deliveroo.Sushi)
	}
	if has.Japanese {
		r.MakeCategoryXML(deliveroo.Japanese)
	}
	if has.Chinese {
		r.MakeCategoryXML(deliveroo.Chinese)
	}
	if has.Western {
		r.MakeCategoryXML(deliveroo.Western)
	}
	if has.FastFood {
		r.MakeCategoryXML(deliveroo.FastFood)
	}
	if has.PartyFood {
		r.MakeCategoryXML(deliveroo.PartyFood)
	}
	if has.DessertAndDrinks {
		r.MakeCategoryXML(deliveroo.DrinksAndDessert)
	}

	placeholder := KVFieldWChildren{
		XMLName: xml.Name{Local: "Placeholder"},
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
								Value:   "11",
							},
							KVFieldWChildren{
								XMLName: xml.Name{Local: "ShopList"},
								Value: []any{
									BasicShop{
										ShopCode:    CDATA{0},
										HomeCode:    CDATA{1},
										Name:        CDATA{"Test"},
										Catchphrase: CDATA{"A"},
										MinPrice:    CDATA{1},
										Yoyaku:      CDATA{1},
										Activate:    CDATA{"on"},
										WaitTime:    CDATA{10},
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
									},
								},
							},
						},
					},
				},
			},
		},
	}

	r.AddCustomType(placeholder)

	/*// It there is no nearby stores, we do not add the placeholder. This will tell the user there are no stores.
	if storesXML != nil && r.request.URL.Query().Get("action") != "webApi_shop_list" {
		r.AddCustomType(placeholder)
	}*/
}

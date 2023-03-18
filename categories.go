package main

import (
	"DemaeDeliveroo/deliveroo"
	"encoding/xml"
)

func (r *Response) MakeCategoryXML(code deliveroo.CategoryCode) {
	var categoryName string
	switch code {
	case deliveroo.Pizza:
		categoryName = "Pizza"
		break
	case deliveroo.BentoBox:
		categoryName = "BentoBox"
		break
	case deliveroo.Sushi:
		categoryName = "Sushi"
		break
	case deliveroo.Japanese:
		categoryName = "Japanese"
		break
	case deliveroo.Chinese:
		categoryName = "Chinese"
		break
	case deliveroo.Western:
		categoryName = "Western"
		break
	case deliveroo.FastFood:
		categoryName = "FastFood"
		break
	case deliveroo.Curry:
		categoryName = "Curry"
		break
	case deliveroo.PartyFood:
		categoryName = "PartyFood"
		break
	case deliveroo.DrinksAndDessert:
		categoryName = "DrinksAndDessert"
		break
	}

	r.AddCustomType(KVFieldWChildren{
		XMLName: xml.Name{Local: categoryName},
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
								Value:   code,
							},
							KVFieldWChildren{
								XMLName: xml.Name{Local: "ShopList"},
								Value: []any{
									BasicShop{
										ShopCode:    CDATA{code},
										HomeCode:    CDATA{code},
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
	})
}

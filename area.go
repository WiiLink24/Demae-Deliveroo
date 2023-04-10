package main

import (
	"context"
	"encoding/xml"
	"math/rand"
	"net/http"
)

const numberBytes = "0123456789"
const InsertUser = `UPDATE "user" SET area_code = $1 WHERE wii_id = $2`

func GenerateAreaCode() string {
	b := make([]byte, 11)
	for i := range b {
		b[i] = numberBytes[rand.Intn(len(numberBytes))]
	}

	return string(b)
}

func areaList(r *Response) {
	areaCode := r.request.URL.Query().Get("areaCode")

	// Nintendo, for whatever reason, require a separate "selectedArea" element
	// as a root node within output.
	// This violates about every XML specification in existence.
	// I am reasonably certain there was a mistake as their function to
	// interpret nodes at levels accepts a parent node, to which they seem to
	// have passed NULL instead of response.
	//
	// We are not going to bother spending time to deal with this.
	if r.request.URL.Query().Get("zipCode") != "" {
		version, apiStatus := GenerateVersionAndAPIStatus()
		r.ResponseFields = []any{
			KVFieldWChildren{
				XMLName: xml.Name{Local: "response"},
				Value: []any{
					KVFieldWChildren{
						XMLName: xml.Name{Local: "areaList"},
						Value: []any{
							KVField{
								XMLName: xml.Name{Local: "segment"},
								Value:   "United States",
							},
							KVFieldWChildren{
								XMLName: xml.Name{Local: "list"},
								Value: []any{
									KVFieldWChildren{
										XMLName: xml.Name{Local: "areaPlace"},
										Value: []any{AreaNames{
											AreaName: CDATA{"place name"},
											AreaCode: CDATA{2},
										}},
									},
								},
							},
						},
					},
					KVField{
						XMLName: xml.Name{Local: "areaCount"},
						Value:   "1",
					},
					version,
					apiStatus,
				},
			},
			KVFieldWChildren{
				XMLName: xml.Name{Local: "selectedArea"},
				Value: []any{
					KVField{
						XMLName: xml.Name{Local: "areaCode"},
						Value:   1,
					},
				},
			},
		}
		return
	}

	if areaCode == "0" {
		r.AddKVWChildNode("areaList", []any{
			KVFieldWChildren{
				XMLName: xml.Name{Local: "place"},
				Value: []any{
					KVField{
						XMLName: xml.Name{Local: "segment"},
						Value:   "Deliveroo",
					},
					KVFieldWChildren{
						XMLName: xml.Name{Local: "list"},
						Value: []any{
							AreaNames{
								AreaName: CDATA{"Deliveroo"},
								AreaCode: CDATA{1},
							},
						},
					},
				},
			},
		})
		r.AddKVNode("areaCount", "2")
		return
	}

	newAreaCode := GenerateAreaCode()
	_, err := pool.Exec(context.Background(), InsertUser, newAreaCode, r.request.Header.Get("X-WiiID"))
	if err != nil {
		r.ReportError(err, http.StatusInternalServerError)
		return
	}
	
	r.AddKVWChildNode("areaList", KVFieldWChildren{
		XMLName: xml.Name{Local: "place"},
		Value: []any{
			KVField{
				XMLName: xml.Name{Local: "container0"},
				Value:   "aaaaa",
			},
			KVField{
				XMLName: xml.Name{Local: "segment"},
				Value:   "Something idk",
			},
			KVFieldWChildren{
				XMLName: xml.Name{Local: "list"},
				Value: []any{
					Area{
						AreaName:   CDATA{"Deliveroo"},
						AreaCode:   CDATA{newAreaCode},
						IsNextArea: CDATA{0},
						Display:    CDATA{1},
						Kanji1:     CDATA{"Deliveroo"},
						Kanji2:     CDATA{"Deliveroo"},
						Kanji3:     CDATA{""},
						Kanji4:     CDATA{""},
					},
				},
			},
		},
	})
	r.AddKVNode("areaCount", "1")
}

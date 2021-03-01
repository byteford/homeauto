package main

// LightItem -
type LightItem struct {
	EntityID string     `json:"entity_id"`
	State    string     `json:"state"`
	Attr     Attributes `json:"attributes"`
}

//Attributes -
type Attributes struct {
	Brightness        int       `json:"brightness"`
	HsColor           []float64 `json:"hs_color"`
	RgbColor          []int     `json:"rgb_color"`
	XyColor           []float64 `json:"xy_color"`
	WhiteValue        int       `json:"white_value"`
	Name              string    `json:"friendly_name"`
	ColorMode         string    `json:"color_mode"`
	SupportedFeatures int       `json:"supported_features"`
}

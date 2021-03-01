package main

// LightItem -
type LightItem struct {
	EntityID          string  `json:"entity_id"`
	State             string  `json:"state"`
	Brightness        int     `json:"brightness"`
	HsColor           []float `json:"hs_color"`
	RgbColor          []int   `json:"rgb_color"`
	XyColor           []float `json:"xy_color"`
	WhiteValue        int     `json:"white_value"`
	Name              string  `json:"friendly_name"`
	ColorMode         string  `json:"color_mode"`
	SupportedFeatures int     `json:"supported_features"`
}

package homeauto

// LightItem is used to hold the data from terraform in a way go can use
// defines the json name for the api for easy conversion
type LightItem struct {
	EntityID string     `json:"entity_id"`
	State    string     `json:"state"`
	Attr     Attributes `json:"attributes"`
}

//Attributes holds the attributes for the light,
//its split from the LightItem to more easily  create the json object
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

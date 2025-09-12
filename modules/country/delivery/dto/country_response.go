package dto

type CountryResponse struct {
	Name         string `json:"name"`
	Code         string `json:"alpha_code"`
	Flags        Flag   `json:"flags"`
	CallingCodes string `json:"calling_code"`
}

type Flag struct {
	Png string `json:"png"`
	Svg string `json:"svg"`
}

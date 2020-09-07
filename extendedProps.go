package msgraph

type SingleValueExtendedProp struct {
	ID    string `json:"id"`
	Value string `json:"value"`
}

type MultiValueExtendedProp struct {
	ID    string   `json:"id"`
	Value []string `json:"value"`
}

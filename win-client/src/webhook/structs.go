package webhook

type AvailPaths struct {
	Name   string `json:"Name,omitempty"`
	Letter string `json:"Letter,omitempty"`
	RW     bool   `json:"RW,omitempty"`
}

type Response struct {
	ID         string       `json:"id,omitempty"`
	Fullname   string       `json:"Fullname,omitempty"`
	AccessKey  string       `json:"AccessKey,omitempty"`
	AvailPaths []AvailPaths `json:"AvailablePaths,omitempty"`
}

var ActivePaths []AvailPaths

package structs

type RcloneInput struct {
	Name string `json:"name,omitempty"`
}

type RcloneResponse struct {
	Error  string      `json:"error,omitempty"`
	Input  RcloneInput `json:"input,omitempty"`
	Path   string      `json:"path,omitempty"`
	Status int         `json:"status,omitempty"`
}

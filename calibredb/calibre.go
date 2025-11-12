package calibredb

type Calibre struct {
	LibraryPath string `json:"library-path,omitempty"`
	Username    string `json:"username,omitempty"`
	Password    string `json:"password,omitempty"`
	Timeout     string `json:"timeout,omitempty"`
}

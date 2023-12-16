package server

type cat struct {
	Name  string `json:"name" xml:"name"`
	Color string `json:"color" xml:"color"`
}

type methodResponse map[string]response

type response struct {
	status int
	body   any
}

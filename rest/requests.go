package rest

type pathRequestBody struct {
	Path      string `json:"path"`
	Namespace string `json:"namespace"`
}

type userRequestBody struct {
	Username  string `json:"username"`
	Namespace string `json:"namespace"`
}

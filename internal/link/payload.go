package link

type CreateLinkRequest struct {
	Url string `json:"url"`
}

type CreateLinkResponse struct {
	Success bool `json:"success"`
}

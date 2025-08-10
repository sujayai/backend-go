package models

type Post struct {
	Slug     string   `json:"slug"`
	Title    string   `json:"title"`
	Date     string   `json:"date"`
	Tags     []string `json:"tags"`
	Summary  string   `json:"summary"`
	Cover    *string  `json:"cover"`
	CoverAlt string   `json:"coverAlt"`
}

type PostWithContent struct {
	Post
	HTML string `json:"html"`
}

type CreatePostRequest struct {
	Title    string `form:"title" binding:"required,min=3,max=140"`
	Date     string `form:"date"`
	Tags     string `form:"tags"`
	Summary  string `form:"summary"`
	Content  string `form:"content"`
	CoverAlt string `form:"coverAlt"`
}

type APIResponse struct {
	OK    bool   `json:"ok,omitempty"`
	Slug  string `json:"slug,omitempty"`
	Error string `json:"error,omitempty"`
}

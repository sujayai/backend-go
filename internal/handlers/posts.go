package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sujayai/backend-go/internal/models"
	"github.com/sujayai/backend-go/internal/storage"
	"github.com/sujayai/backend-go/internal/utils"
)

type PostsHandler struct {
	storage *storage.Storage
}

func NewPostsHandler(s *storage.Storage) *PostsHandler {
	return &PostsHandler{storage: s}
}

func (h *PostsHandler) GetPosts(c *gin.Context) {
	posts, err := h.storage.GetPosts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get posts"})
		return
	}

	c.JSON(http.StatusOK, posts)
}

func (h *PostsHandler) GetPost(c *gin.Context) {
	slug := c.Param("slug")

	post, err := h.storage.GetPost(slug)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not_found"})
		return
	}

	content, err := h.storage.GetPostContent(slug)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get content"})
		return
	}

	html := utils.MarkdownToHTML(content)

	postWithContent := models.PostWithContent{
		Post: *post,
		HTML: html,
	}

	c.JSON(http.StatusOK, postWithContent)
}

func (h *PostsHandler) CreatePost(c *gin.Context) {
	var req models.CreatePostRequest

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title_invalid"})
		return
	}

	// Get existing posts for slug uniqueness
	posts, _ := h.storage.GetPosts()

	// Generate slug
	baseSlug := utils.Slugify(req.Title)
	slug := utils.UniqueSlug(baseSlug, posts)

	// Parse tags
	tags := utils.ParseTags(req.Tags)

	// Format date
	date := utils.FormatDate(req.Date)

	// Get content from form or uploaded file
	content := req.Content
	if mdFile, err := c.FormFile("mdFile"); err == nil {
		file, err := mdFile.Open()
		if err == nil {
			defer file.Close()
			data := make([]byte, mdFile.Size)
			file.Read(data)
			content = string(data)
		}
	}

	if strings.TrimSpace(content) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "content_missing"})
		return
	}

	// Handle cover image upload
	var cover *string
	if coverFile, err := c.FormFile("cover"); err == nil {
		file, err := coverFile.Open()
		if err == nil {
			defer file.Close()
			data := make([]byte, coverFile.Size)
			file.Read(data)

			coverURL, err := h.storage.SaveUpload(slug, coverFile.Filename, data)
			if err == nil {
				cover = &coverURL
			}
		}
	}

	// Limit summary length
	summary := req.Summary
	if len(summary) > 240 {
		summary = summary[:240]
	}

	// Limit coverAlt length
	coverAlt := req.CoverAlt
	if len(coverAlt) > 120 {
		coverAlt = coverAlt[:120]
	}

	// Create post
	post := models.Post{
		Slug:     slug,
		Title:    req.Title,
		Date:     date,
		Tags:     tags,
		Summary:  summary,
		Cover:    cover,
		CoverAlt: coverAlt,
	}

	// Save post
	if err := h.storage.SavePost(post, content); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save post"})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		OK:   true,
		Slug: slug,
	})
}

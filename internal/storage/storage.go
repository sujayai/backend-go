package storage

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/sujayai/backend-go/internal/models"
	"github.com/sujayai/backend-go/internal/utils"
	"github.com/sujayai/backend-go/pkg/config"
)

type Storage struct {
	config *config.Config
}

func New(cfg *config.Config) (*Storage, error) {
	s := &Storage{config: cfg}

	// Create directories
	dirs := []string{cfg.PostsDir, cfg.UploadsDir, cfg.DataDir}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// Create empty posts.json if it doesn't exist
	if _, err := os.Stat(cfg.DBFile); os.IsNotExist(err) {
		if err := ioutil.WriteFile(cfg.DBFile, []byte("[]"), 0644); err != nil {
			return nil, fmt.Errorf("failed to create posts.json: %w", err)
		}
	}

	return s, nil
}

func (s *Storage) GetPosts() ([]models.Post, error) {
	data, err := ioutil.ReadFile(s.config.DBFile)
	if err != nil {
		return []models.Post{}, nil
	}

	var posts []models.Post
	if err := json.Unmarshal(data, &posts); err != nil {
		return []models.Post{}, nil
	}

	utils.SortPostsByDate(posts)
	return posts, nil
}

func (s *Storage) GetPost(slug string) (*models.Post, error) {
	posts, err := s.GetPosts()
	if err != nil {
		return nil, err
	}

	for _, post := range posts {
		if post.Slug == slug {
			return &post, nil
		}
	}

	return nil, fmt.Errorf("post not found")
}

func (s *Storage) GetPostContent(slug string) (string, error) {
	mdPath := filepath.Join(s.config.PostsDir, slug+".md")

	if _, err := os.Stat(mdPath); os.IsNotExist(err) {
		return "", nil
	}

	content, err := ioutil.ReadFile(mdPath)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

func (s *Storage) SavePost(post models.Post, content string) error {
	// Save markdown content
	mdPath := filepath.Join(s.config.PostsDir, post.Slug+".md")
	if err := ioutil.WriteFile(mdPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to save markdown: %w", err)
	}

	// Update posts.json
	posts, err := s.GetPosts()
	if err != nil {
		return err
	}

	// Remove existing post with same slug
	filtered := make([]models.Post, 0, len(posts))
	for _, p := range posts {
		if p.Slug != post.Slug {
			filtered = append(filtered, p)
		}
	}

	// Add new post
	filtered = append(filtered, post)
	utils.SortPostsByDate(filtered)

	// Save to file
	data, err := json.MarshalIndent(filtered, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(s.config.DBFile, data, 0644)
}

func (s *Storage) SaveUpload(slug, originalName string, fileData []byte) (string, error) {
	ext := filepath.Ext(originalName)
	if ext == "" {
		ext = ".jpg"
	}

	filename := slug + ext
	filePath := filepath.Join(s.config.UploadsDir, filename)

	if err := ioutil.WriteFile(filePath, fileData, 0644); err != nil {
		return "", fmt.Errorf("failed to save upload: %w", err)
	}

	return "/uploads/" + filename, nil
}

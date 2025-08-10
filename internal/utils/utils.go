package utils

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/microcosm-cc/bluemonday"
	"github.com/sujayai/backend-go/internal/models"
)

func Slugify(s string) string {
	// Convert to lowercase and replace non-alphanumeric with spaces
	reg := regexp.MustCompile(`[^a-z0-9\s-]`)
	s = strings.ToLower(s)
	s = reg.ReplaceAllString(s, "")
	s = strings.TrimSpace(s)

	// Replace multiple spaces/dashes with single dash
	reg2 := regexp.MustCompile(`\s+`)
	s = reg2.ReplaceAllString(s, "-")
	reg3 := regexp.MustCompile(`-+`)
	s = reg3.ReplaceAllString(s, "-")

	// Limit length
	if len(s) > 80 {
		s = s[:80]
	}

	return strings.Trim(s, "-")
}

func UniqueSlug(base string, posts []models.Post) string {
	if base == "" {
		base = fmt.Sprintf("post-%d", time.Now().Unix())
	}

	slug := base
	n := 2

	for {
		exists := false
		for _, p := range posts {
			if p.Slug == slug {
				exists = true
				break
			}
		}
		if !exists {
			break
		}
		slug = fmt.Sprintf("%s-%d", base, n)
		n++
	}

	return slug
}

func ParseTags(tags string) []string {
	if tags == "" {
		return []string{}
	}

	tagList := strings.Split(tags, ",")
	result := make([]string, 0, len(tagList))

	for _, tag := range tagList {
		tag = strings.TrimSpace(tag)
		if tag != "" && len(tag) <= 24 {
			result = append(result, tag)
		}
		if len(result) >= 8 {
			break
		}
	}

	return result
}

func MarkdownToHTML(md string) string {
	// Create parser with extensions
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)

	// Create HTML renderer
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	// Parse markdown to HTML
	htmlBytes := markdown.ToHTML([]byte(md), p, renderer)

	// Sanitize HTML
	policy := bluemonday.UGCPolicy()
	policy.AllowElements("h1", "h2", "h3", "figure", "figcaption", "pre", "code", "img")
	policy.AllowAttrs("src", "alt", "title", "loading").OnElements("img")

	return policy.Sanitize(string(htmlBytes))
}

func FormatDate(dateStr string) string {
	if dateStr == "" {
		return time.Now().Format("2006-01-02")
	}

	if t, err := time.Parse(time.RFC3339, dateStr); err == nil {
		return t.Format("2006-01-02")
	}

	if t, err := time.Parse("2006-01-02", dateStr); err == nil {
		return t.Format("2006-01-02")
	}

	return time.Now().Format("2006-01-02")
}

func SortPostsByDate(posts []models.Post) {
	sort.Slice(posts, func(i, j int) bool {
		dateI, _ := time.Parse("2006-01-02", posts[i].Date)
		dateJ, _ := time.Parse("2006-01-02", posts[j].Date)
		return dateI.After(dateJ)
	})
}

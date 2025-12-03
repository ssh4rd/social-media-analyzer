package controller

import (
	"fmt"
	"html/template"
	"net/http"

	"social-media-analyzer/internal/models"
	"social-media-analyzer/internal/transport/http/router"

	"gorm.io/gorm"
)

type MainController struct {
	db *gorm.DB
}

type PageData struct {
	Groups []GroupStats
}

type GroupStats struct {
	ID                 uint
	Domain             string
	TotalPosts         int
	TotalLikes         int
	AvgLikesPerPost    float64
	MaxLikesPerPost    int
	AvgCommentsPerPost float64
	PostsLastWeek      int
}

func NewMainController(db *gorm.DB) *MainController {
	return &MainController{db: db}
}

// GetMainPage handles GET / requests
func (mc *MainController) GetMainPage(w http.ResponseWriter, r *http.Request, params router.Params) {
	// Fetch groups
	var groups []models.Group
	if err := mc.db.Find(&groups).Error; err != nil {
		http.Error(w, "Failed to fetch groups", http.StatusInternalServerError)
		return
	}

	// Calculate stats for each group
	pageData := PageData{}
	for _, group := range groups {
		var posts []models.Post

		if err := mc.db.Where("group_id = ?", group.ID).Find(&posts).Error; err != nil {
			continue
		}

		stats := GroupStats{
			ID:         group.ID,
			Domain:     group.Domain,
			TotalPosts: len(posts),
		}

		fmt.Println(posts)

		if len(posts) > 0 {
			totalLikes := 0
			maxLikes := 0
			totalComments := 0

			for _, post := range posts {
				totalLikes += post.Likes
				totalComments += post.Comments
				if post.Likes > maxLikes {
					maxLikes = post.Likes
				}
			}

			stats.TotalLikes = totalLikes
			stats.MaxLikesPerPost = maxLikes
			stats.AvgLikesPerPost = float64(totalLikes) / float64(len(posts))
			stats.AvgCommentsPerPost = float64(totalComments) / float64(len(posts))
		}

		pageData.Groups = append(pageData.Groups, stats)
	}

	tpl := template.Must(template.ParseFiles("web/templates/main.html"))
	tpl.Execute(w, pageData)
}

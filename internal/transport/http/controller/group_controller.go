package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"social-media-analyzer/internal/models"
	"social-media-analyzer/internal/service"
	"social-media-analyzer/internal/transport/http/router"

	"gorm.io/gorm"
)

type GroupController struct {
	db        *gorm.DB
	vkService *service.VKService
}

type AddGroupRequest struct {
	Link string `json:"link"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

type SuccessResponse struct {
	Message string `json:"message"`
	GroupID uint   `json:"group_id"`
}

func NewGroupController(db *gorm.DB, vkService *service.VKService) *GroupController {
	return &GroupController{db: db, vkService: vkService}
}

// AddGroup handles POST /api/groups requests
func (gc *GroupController) AddGroup(w http.ResponseWriter, r *http.Request, params router.Params) {

	fmt.Println("dfsdfsdgs")

	w.Header().Set("Content-Type", "application/json")

	var req AddGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "Invalid request body"})
		return
	}

	fmt.Println(req)

	if req.Link == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "Link is required"})
		return
	}

	// Parse group from link using VK API
	group, err := gc.vkService.ParseGroupFromLink(req.Link)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Message: err.Error()})
		return
	}

	// Check if group already exists
	// var existingGroup models.Group
	// if result := gc.db.Where("domain = ?", group.Domain).First(&existingGroup); result.Error == nil {
	// 	w.WriteHeader(http.StatusConflict)
	// 	json.NewEncoder(w).Encode(ErrorResponse{Message: "Group already exists"})
	// 	return
	// }

	// Create group in database
	if result := gc.db.Create(&group); result.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "Failed to add group"})
		return
	}

	// Fetch wall posts asynchronously
	go func() {
		if err := gc.fetchAndSaveWallPosts(group); err != nil {
			log.Printf("Failed to fetch wall posts for group %s: %v\n", group.Domain, err)
		}
	}()

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(SuccessResponse{
		Message: "Group added successfully",
		GroupID: group.ID,
	})
}

// fetchAndSaveWallPosts fetches wall posts from VK API and saves them to database
func (gc *GroupController) fetchAndSaveWallPosts(group *models.Group) error {
	// Fetch wall posts from VK API
	wallPosts, err := gc.vkService.GetWallPosts(group.Domain, 100)
	if err != nil {
		return err
	}

	// Convert VK posts to model posts
	for _, vkPost := range wallPosts {
		post := models.Post{
			GroupID:   group.ID,
			Date:      time.Unix(int64(vkPost.Date), 0).Format("2006-01-02"),
			Text:      vkPost.Text,
			Views:     vkPost.Views.Count,
			Reactions: vkPost.Likes.Count, // Using likes as reactions for now
			Likes:     vkPost.Likes.Count,
			Comments:  vkPost.Comments.Count,
		}

		// Check if post already exists
		var existingPost models.Post
		if result := gc.db.Where("group_id = ? AND date = ? AND text = ?", group.ID, post.Date, post.Text).First(&existingPost); result.Error == nil {
			continue // Post already exists
		}

		// Save post to database
		if result := gc.db.Create(&post); result.Error != nil {
			log.Printf("Failed to save post for group %s: %v\n", group.Domain, result.Error)
			continue
		}
	}

	return nil
}

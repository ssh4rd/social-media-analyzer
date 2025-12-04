package models

import (
	"testing"
	"time"
)

// TestGroupModelCreation tests Group model creation and fields
func TestGroupModelCreation(t *testing.T) {
	now := time.Now()
	group := Group{
		ID:          1,
		Domain:      "testgroup",
		Subscribers: 1000,
		ParsedAt:    now,
		Posts:       []Post{},
	}

	tests := []struct {
		name     string
		field    string
		expected interface{}
		actual   interface{}
	}{
		{"ID", "ID", uint(1), group.ID},
		{"Domain", "Domain", "testgroup", group.Domain},
		{"Subscribers", "Subscribers", 1000, group.Subscribers},
		{"ParsedAt", "ParsedAt", now, group.ParsedAt},
		{"Posts", "Posts", 0, len(group.Posts)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expected != tt.actual {
				t.Errorf("%s: expected %v, got %v", tt.field, tt.expected, tt.actual)
			}
		})
	}
}

// TestGroupModelWithPosts tests Group model with related posts
func TestGroupModelWithPosts(t *testing.T) {
	posts := []Post{
		{ID: 1, GroupID: 1, Likes: 100, Comments: 10},
		{ID: 2, GroupID: 1, Likes: 200, Comments: 20},
		{ID: 3, GroupID: 1, Likes: 150, Comments: 15},
	}

	group := Group{
		ID:          1,
		Domain:      "groupwithposts",
		Subscribers: 5000,
		ParsedAt:    time.Now(),
		Posts:       posts,
	}

	if len(group.Posts) != 3 {
		t.Errorf("Expected 3 posts, got %d", len(group.Posts))
	}

	for i, post := range group.Posts {
		if post.GroupID != group.ID {
			t.Errorf("Post %d: expected GroupID %d, got %d", i, group.ID, post.GroupID)
		}
	}
}

// TestGroupModelEmptyDomain tests Group with empty domain
func TestGroupModelEmptyDomain(t *testing.T) {
	group := Group{
		ID:       1,
		Domain:   "",
		Posts:    []Post{},
	}

	if group.Domain != "" {
		t.Errorf("Expected empty domain, got '%s'", group.Domain)
	}
}

// TestGroupModelZeroSubscribers tests Group with zero subscribers
func TestGroupModelZeroSubscribers(t *testing.T) {
	group := Group{
		ID:          1,
		Domain:      "newgroup",
		Subscribers: 0,
		Posts:       []Post{},
	}

	if group.Subscribers != 0 {
		t.Errorf("Expected 0 subscribers, got %d", group.Subscribers)
	}
}

// TestPostModelCreation tests Post model creation and fields
func TestPostModelCreation(t *testing.T) {
	post := Post{
		ID:        1,
		GroupID:   1,
		Date:      "2025-12-04",
		Group:     Group{ID: 1, Domain: "testgroup"},
		Views:     1000,
		Reactions: 100,
		Likes:     50,
		Text:      "Test post content",
		Comments:  10,
	}

	tests := []struct {
		name     string
		field    string
		expected interface{}
		actual   interface{}
	}{
		{"ID", "ID", uint(1), post.ID},
		{"GroupID", "GroupID", uint(1), post.GroupID},
		{"Date", "Date", "2025-12-04", post.Date},
		{"Views", "Views", 1000, post.Views},
		{"Reactions", "Reactions", 100, post.Reactions},
		{"Likes", "Likes", 50, post.Likes},
		{"Text", "Text", "Test post content", post.Text},
		{"Comments", "Comments", 10, post.Comments},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expected != tt.actual {
				t.Errorf("%s: expected %v, got %v", tt.field, tt.expected, tt.actual)
			}
		})
	}
}

// TestPostModelGroupRelationship tests Post-Group relationship
func TestPostModelGroupRelationship(t *testing.T) {
	group := Group{
		ID:          5,
		Domain:      "relatedgroup",
		Subscribers: 2000,
	}

	post := Post{
		ID:      1,
		GroupID: 5,
		Group:   group,
		Likes:   100,
	}

	if post.GroupID != group.ID {
		t.Errorf("Expected GroupID %d, got %d", group.ID, post.GroupID)
	}

	if post.Group.Domain != "relatedgroup" {
		t.Errorf("Expected group domain 'relatedgroup', got '%s'", post.Group.Domain)
	}
}

// TestPostModelEngagementMetrics tests post engagement metrics
func TestPostModelEngagementMetrics(t *testing.T) {
	tests := []struct {
		name      string
		likes     int
		comments  int
		views     int
		reactions int
	}{
		{
			name:      "High engagement post",
			likes:     500,
			comments:  100,
			views:     5000,
			reactions: 500,
		},
		{
			name:      "Low engagement post",
			likes:     10,
			comments:  2,
			views:     100,
			reactions: 10,
		},
		{
			name:      "Zero engagement post",
			likes:     0,
			comments:  0,
			views:     0,
			reactions: 0,
		},
		{
			name:      "High comments low likes",
			likes:     50,
			comments:  200,
			views:     1000,
			reactions: 50,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			post := Post{
				ID:        1,
				GroupID:   1,
				Likes:     tt.likes,
				Comments:  tt.comments,
				Views:     tt.views,
				Reactions: tt.reactions,
			}

			if post.Likes != tt.likes {
				t.Errorf("Expected %d likes, got %d", tt.likes, post.Likes)
			}
			if post.Comments != tt.comments {
				t.Errorf("Expected %d comments, got %d", tt.comments, post.Comments)
			}
		})
	}
}

// TestPostModelWithLongText tests Post with long text content
func TestPostModelWithLongText(t *testing.T) {
	longText := "This is a very long post content that might contain multiple paragraphs and various information about the group and its activities. " +
		"It tests whether the Post model can properly handle extensive text data without any issues. " +
		"The text field should be able to store any amount of text without truncation or errors."

	post := Post{
		ID:      1,
		GroupID: 1,
		Text:    longText,
	}

	if post.Text != longText {
		t.Errorf("Text mismatch: expected length %d, got %d", len(longText), len(post.Text))
	}
}

// TestPostModelDateFormats tests Post with various date formats
func TestPostModelDateFormats(t *testing.T) {
	tests := []struct {
		name     string
		dateStr  string
	}{
		{"ISO date", "2025-12-04"},
		{"Full timestamp", "2025-12-04 15:30:45"},
		{"Alternative format", "04/12/2025"},
		{"Empty date", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			post := Post{
				ID:      1,
				GroupID: 1,
				Date:    tt.dateStr,
			}

			if post.Date != tt.dateStr {
				t.Errorf("Expected date '%s', got '%s'", tt.dateStr, post.Date)
			}
		})
	}
}

// TestGroupModelBulkCreation tests creating multiple groups
func TestGroupModelBulkCreation(t *testing.T) {
	groupCount := 10
	groups := make([]Group, groupCount)

	for i := 0; i < groupCount; i++ {
		groups[i] = Group{
			ID:          uint(i + 1),
			Domain:      "group" + string(rune(i+1)),
			Subscribers: (i + 1) * 100,
			ParsedAt:    time.Now(),
		}
	}

	if len(groups) != groupCount {
		t.Errorf("Expected %d groups, got %d", groupCount, len(groups))
	}

	for i, group := range groups {
		expectedID := uint(i + 1)
		if group.ID != expectedID {
			t.Errorf("Group %d: expected ID %d, got %d", i, expectedID, group.ID)
		}
	}
}

// TestPostModelBulkCreation tests creating multiple posts
func TestPostModelBulkCreation(t *testing.T) {
	postCount := 20
	posts := make([]Post, postCount)

	for i := 0; i < postCount; i++ {
		posts[i] = Post{
			ID:       uint(i + 1),
			GroupID:  1,
			Date:     "2025-12-04",
			Likes:    (i + 1) * 10,
			Comments: (i + 1) * 2,
			Views:    (i + 1) * 100,
		}
	}

	if len(posts) != postCount {
		t.Errorf("Expected %d posts, got %d", postCount, len(posts))
	}

	totalLikes := 0
	for _, post := range posts {
		totalLikes += post.Likes
	}

	expectedTotalLikes := (postCount * (postCount + 1) / 2) * 10
	if totalLikes != expectedTotalLikes {
		t.Errorf("Expected total likes %d, got %d", expectedTotalLikes, totalLikes)
	}
}

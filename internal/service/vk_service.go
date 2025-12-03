package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"social-media-analyzer/internal/config"
	"social-media-analyzer/internal/models"
)

type VKService struct {
	accessToken string
	apiVersion  string
	httpClient  *http.Client
}

type VKGroupInfo struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Domain string `json:"screen_name"`
}

type VKGroupResponse struct {
	Response []VKGroupInfo `json:"response"`
	Error    struct {
		ErrorCode int    `json:"error_code"`
		ErrorMsg  string `json:"error_msg"`
	} `json:"error"`
}

type VKWallPost struct {
	ID       int    `json:"id"`
	OwnerID  int    `json:"owner_id"`
	Date     int    `json:"date"`
	Text     string `json:"text"`
	Likes    struct {
		Count int `json:"count"`
	} `json:"likes"`
	Comments struct {
		Count int `json:"count"`
	} `json:"comments"`
	Reposts struct {
		Count int `json:"count"`
	} `json:"reposts"`
	Views struct {
		Count int `json:"count"`
	} `json:"views"`
}

type VKWallResponse struct {
	Response struct {
		Count int          `json:"count"`
		Items []VKWallPost `json:"items"`
	} `json:"response"`
	Error struct {
		ErrorCode int    `json:"error_code"`
		ErrorMsg  string `json:"error_msg"`
	} `json:"error"`
}

func NewVKService(cfg *config.VKConfig) *VKService {
	return &VKService{
		accessToken: cfg.AccessToken,
		apiVersion:  cfg.APIVersion,
		httpClient:  &http.Client{},
	}
}

// ExtractGroupScreenName extracts group screen_name from various VK link formats
// Supports formats: https://vk.com/groupname, vk.com/groupname, https://vk.com/club123456
func (s *VKService) ExtractGroupScreenName(link string) (string, error) {
	// Remove protocol and www if present
	link = strings.TrimPrefix(link, "https://")
	link = strings.TrimPrefix(link, "http://")
	link = strings.TrimPrefix(link, "www.")

	// Extract path part (after vk.com/)
	parts := strings.Split(link, "/")
	var screenName string

	for _, part := range parts {
		if strings.Contains(part, "vk.com") {
			continue
		}
		if part != "" {
			screenName = part
			break
		}
	}

	if screenName == "" {
		return "", fmt.Errorf("could not extract group name from link")
	}

	// If it's a club ID (starts with "club" or "-"), keep it as is
	// Otherwise, it's a screen name
	return screenName, nil
}

// GetGroupInfo fetches group information from VK API
func (s *VKService) GetGroupInfo(screenName string) (*models.Group, error) {
	if s.accessToken == "" {
		return nil, fmt.Errorf("VK access token not configured")
	}

	// Build API request URL
	url := fmt.Sprintf(
		"https://api.vk.com/method/groups.getById?group_ids=%s&fields=members,description&v=%s&access_token=%s",
		screenName, s.apiVersion, s.accessToken,
	)

	resp, err := s.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch group info: %w", err)
	}
	defer resp.Body.Close()

	var vkResp VKGroupResponse
	if err := json.NewDecoder(resp.Body).Decode(&vkResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if vkResp.Error.ErrorCode != 0 {
		return nil, fmt.Errorf("VK API error: %s", vkResp.Error.ErrorMsg)
	}

	if len(vkResp.Response) == 0 {
		return nil, fmt.Errorf("group not found")
	}

	groupInfo := vkResp.Response[0]

	group := &models.Group{
		Domain: groupInfo.Domain,
	}

	return group, nil
}

// GetWallPosts fetches posts from group wall
func (s *VKService) GetWallPosts(screenName string, count int) ([]VKWallPost, error) {
	if s.accessToken == "" {
		return nil, fmt.Errorf("VK access token not configured")
	}

	if count <= 0 || count > 100 {
		count = 100
	}

	// Build API request URL
	// wall.get method with owner_id (negative for groups)
	url := fmt.Sprintf(
		"https://api.vk.com/method/wall.get?domain=%s&count=%d&v=%s&access_token=%s",
		screenName, count, s.apiVersion, s.accessToken,
	)

	resp, err := s.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch wall posts: %w", err)
	}
	defer resp.Body.Close()

	var vkResp VKWallResponse
	if err := json.NewDecoder(resp.Body).Decode(&vkResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if vkResp.Error.ErrorCode != 0 {
		return nil, fmt.Errorf("VK API error: %s", vkResp.Error.ErrorMsg)
	}

	return vkResp.Response.Items, nil
}

// ParseGroupFromLink parses group link and fetches info from VK API
func (s *VKService) ParseGroupFromLink(link string) (*models.Group, error) {
	screenName, err := s.ExtractGroupScreenName(link)
	if err != nil {
		return nil, fmt.Errorf("invalid group link: %w", err)
	}

	fmt.Println(screenName)

	// Validate screen name (alphanumeric, dash, underscore, or starts with club/public)
	validPattern := regexp.MustCompile(`^([a-zA-Z0-9_-]+|club\d+|public\d+|-\d+)$`)
	if !validPattern.MatchString(screenName) {
		return nil, fmt.Errorf("invalid screen name format")
	}

	group, err := s.GetGroupInfo(screenName)
	if err != nil {
		return nil, fmt.Errorf("failed to get group info: %w", err)
	}

	return group, nil
}

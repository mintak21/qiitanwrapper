package model

import "time"

// QiitaItem is Response of Qiita API v2 Item
type QiitaItem struct {
	RenderedBody  string    `json:"rendered_body"`
	Body          string    `json:"body"`
	Coediting     bool      `json:"coediting"`
	CommentsCount int       `json:"comments_count"`
	CreatedAt     time.Time `json:"created_at"`
	Group         struct {
		CreatedAt time.Time `json:"created_at"`
		ID        int       `json:"id"`
		Name      string    `json:"name"`
		Private   bool      `json:"private"`
		UpdatedAt time.Time `json:"updated_at"`
		URLName   string    `json:"url_name"`
	} `json:"group"`
	ID             string `json:"id"`
	LikesCount     int    `json:"likes_count"`
	Private        bool   `json:"private"`
	ReactionsCount int    `json:"reactions_count"`
	Tags           []struct {
		Name     string   `json:"name"`
		Versions []string `json:"versions"`
	} `json:"tags"`
	Title     string    `json:"title"`
	UpdatedAt time.Time `json:"updated_at"`
	URL       string    `json:"url"`
	User      struct {
		Description       string `json:"description"`
		FacebookID        string `json:"facebook_id"`
		FolloweesCount    int    `json:"followees_count"`
		FollowersCount    int    `json:"followers_count"`
		GithubLoginName   string `json:"github_login_name"`
		ID                string `json:"id"`
		ItemsCount        int    `json:"items_count"`
		LinkedinID        string `json:"linkedin_id"`
		Location          string `json:"location"`
		Name              string `json:"name"`
		Organization      string `json:"organization"`
		PermanentID       int    `json:"permanent_id"`
		ProfileImageURL   string `json:"profile_image_url"`
		TeamOnly          bool   `json:"team_only"`
		TwitterScreenName string `json:"twitter_screen_name"`
		WebsiteURL        string `json:"website_url"`
	} `json:"user"`
	PageViewsCount int  `json:"page_views_count"`
	HasNext        bool `json:"-"`
}

package model

type SocialAction struct {
	PostID    string `json:"post_id"`
	Topic     string `json:"topic"`
	Num       int    `json:"num"`
	ThumbType string `json:"thumb_type"`
}

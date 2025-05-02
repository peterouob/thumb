package model

type SocialAction struct {
	PostID  string `json:"post_id"`
	Topic   string `json:"topic"`
	Like    int    `json:"like"`
	Forward int    `json:"forward"`
	Dislike int    `json:"dislike"`
}

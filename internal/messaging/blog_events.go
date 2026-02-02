package messaging

import "time"

// BlogPublishedEvent for notifying followers when blog is published
type BlogPublishedEvent struct {
	BlogID      string    `json:"blog_id"`
	AuthorID    string    `json:"author_id"`
	AuthorName  string    `json:"author_name"`
	Title       string    `json:"title"`
	Excerpt     string    `json:"excerpt"`
	PublishedAt time.Time `json:"published_at"`
}

// CommentAddedEvent for notifying blog owner when comment is added
type CommentAddedEvent struct {
	CommentID   string    `json:"comment_id"`
	BlogID      string    `json:"blog_id"`
	BlogTitle   string    `json:"blog_title"`
	CommenterID string    `json:"commenter_id"`
	Commenter   string    `json:"commenter_name"`
	Comment     string    `json:"comment"`
	BlogOwnerID string    `json:"blog_owner_id"`
	CreatedAt   time.Time `json:"created_at"`
}

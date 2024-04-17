package poll

type Poll struct {
	Id          int
	Title       string
	Description string
	Slug        string
	CreatedBy   string
	PollId      int `json:"poll_id"`
}

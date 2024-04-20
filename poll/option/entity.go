package option

type Option struct {
	Id        int    `json:"option_id,omitempty"`
	Title     string `json:"title"`
	PollId    int    `json:"poll_id"`
	CreatedAt string `json:"created_at"`
}

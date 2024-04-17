package poll

type Poll struct {
	Id       int
	Email    string
    PollId   int `json:"poll_id"`
}

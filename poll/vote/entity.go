package vote

type Vote struct {
	Id       int
	Email    string
    PollId   int `json:"poll_id"`
}

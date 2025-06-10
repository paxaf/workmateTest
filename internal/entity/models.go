package entity

const (
	statusNew  = "new"
	statusErr  = "err"
	statusDone = "done"
)

type Task struct {
	Id      string `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Status  string `json:"status"`
}

type TaskResponse struct {
	Tasks []Task `json:"tasks"`
}

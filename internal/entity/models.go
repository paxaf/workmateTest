package entity

const (
	StatusNew        = "new"
	StatusInProgress = "in progress"
	StatusDone       = "done"
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

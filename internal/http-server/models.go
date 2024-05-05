package httpserver

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

type Task struct {
	Title       string `json:"title"`
	Is_finished bool   `json:"is_finished"`
}

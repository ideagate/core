package websocket

type Event struct {
	Id        string `json:"id"`
	ProjectId string `json:"project_id"`
	Type      string `json:"type"`
	Data      any    `json:"data"`
}

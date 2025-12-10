package log

type LogRequest struct {
	Log    string `json:"log" validate:"required"`
	Labels Labels `json:"labels" validate:"required"`
	Token  string `json:"token" validate:"required,uuid"`
}

type Labels struct {
	Service       string `json:"service" validate:"required"`
	Node          string `json:"node"`
	NodeId        string `json:"node_id"`
	ContainerName string `json:"container_name"`
	ContainerId   string `json:"container_hostname"`
}

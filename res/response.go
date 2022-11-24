package res

type Response struct {
	Success bool                   `json:"success"`
	Message string                 `json:"message,omitempty" extensions:"x-omitempty" example:"Error message"`
	Data    map[string]interface{} `json:"body,omitempty" extensions:"x-omitempty"`
}

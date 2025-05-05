package dto

type HealthResponse struct {
	Status string `json:"status"`
}

func NewHealthResponse(status string) HealthResponse {
	return HealthResponse{
		Status: status,
	}
}

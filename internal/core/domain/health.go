package domain

type HealthStatus struct {
	Environment string `json:"environment"`
	Version     string `json:"version"`
	Status      string `json:"status"`
}

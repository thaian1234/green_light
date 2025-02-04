package domain

import "math"

type Metadata struct {
	CurrentPage  int `json:"current_page,omitempty"`
	CurrentSize  int `json:"current_size,omitempty"`
	FirstPage    int `json:"first_page,omitempty"`
	LastPage     int `json:"last_page,omitempty"`
	TotalRecords int `json:"total_records,omitempty"`
	TotalPages   int `json:"total_pages,omitempty"`
}

func CalculateMetadata(totalRecords, page, size int) Metadata {
	if totalRecords == 0 {
		return Metadata{}
	}
	return Metadata{
		CurrentPage:  page,
		CurrentSize:  size,
		FirstPage:    1,
		LastPage:     int(math.Ceil(float64(totalRecords) / float64(size))),
		TotalRecords: totalRecords,
	}
}

package domain

import "strings"

type Filter struct {
	Page         int    `form:"page" binding:"omitempty,page"`
	Size         int    `form:"size" binding:"omitempty,size"`
	Sort         string `form:"sort" binding:"omitempty,sort"`
	SortSafeList []string
}

func (f *Filter) SortColumn() string {
	for _, safeValue := range f.SortSafeList {
		if f.Sort == safeValue {
			return strings.TrimPrefix(f.Sort, "-")
		}
	}
	return "created_at"
}

func (f *Filter) SortDirection() string {
	if strings.HasPrefix(f.Sort, "-") {
		return "DESC"
	}
	return "ASC"
}

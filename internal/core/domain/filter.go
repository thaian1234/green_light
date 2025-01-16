package domain

type Filter struct {
	Page         int    `form:"page" binding:"omitempty,page"`
	Size         int    `form:"size" binding:"omitempty,size"`
	Sort         string `form:"sort" binding:"omitempty,sort"`
	SortSafeList []string
}

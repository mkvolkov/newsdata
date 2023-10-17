package model

//go:generate reform

// reform:News
type Article struct {
	ID         int32  `json:"Id" reform:"id,pk"`
	Title      string `json:"Title" reform:"title"`
	Content    string `json:"Content" reform:"content"`
	Categories []int  `json:"Categories" reform:"-"`
}

// reform:NewsCategories
type Category struct {
	ID    int32 `reform:"NewsId"`
	CatID int32 `reform:"CategoryId"`
}

type ListNews struct {
	Success bool       `json:"Success"`
	AllNews []*Article `json:"News"`
}

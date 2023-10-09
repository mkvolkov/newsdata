package model

//go:generate reform

// reform:News
type Article struct {
	ID      int32  `reform:"id,pk"`
	Title   string `reform:"title"`
	Content string `reform:"content"`
}

// reform:NewsCategories
type Category struct {
	ID    int32 `reform:"NewsId"`
	CatID int32 `reform:"CategoryId"`
}

type ListNews struct {
	Success bool   `json:"Success"`
	AllNews []News `json:"News"`
}

type News struct {
	ID         int32  `json:"Id"`
	Title      string `json:"Title"`
	Content    string `json:"Content"`
	Categories []int  `json:"Categories"`
}

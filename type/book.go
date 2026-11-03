package types

type Book struct {
	ID            string   `json:"id"`
	Title         string   `json:"title"`
	Authors       []string `json:"authors"`
	CoverURL      *string  `json:"coverUrl"`
	Description   *string  `json:"description"`
	PageCount     *int     `json:"pageCount"`
	PublishedDate *string  `json:"publishedDate"`
	Publisher     *string  `json:"publisher"`
	ISBN          *string  `json:"isbn"`
	AddMode       string   `json:"addMode"`
	VolumeID      *string  `json:"volumeId"`
	BookshelfID   int      `json:"bookshelfId"`
}

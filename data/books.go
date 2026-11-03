package data

import types "erastos-backend/type"

func stringPtr(value string) *string {
	return &value
}

var Books = []types.Book{
	{
		ID:          "1",
		Title:       "The Great Gatsby",
		Authors:     []string{"F. Scott Fitzgerald"},
		ISBN:        stringPtr("9780743273565"),
		AddMode:     "MANUAL",
		BookshelfID: 1,
	},
	{
		ID:          "2",
		Title:       "To Kill a Mockingbird",
		Authors:     []string{"Harper Lee"},
		ISBN:        stringPtr("9780061120084"),
		AddMode:     "MANUAL",
		BookshelfID: 1,
	},
	{
		ID:          "3",
		Title:       "1984",
		Authors:     []string{"George Orwell"},
		ISBN:        stringPtr("9780451524935"),
		AddMode:     "MANUAL",
		BookshelfID: 1,
	},
	{
		ID:          "4",
		Title:       "Pride and Prejudice",
		Authors:     []string{"Jane Austen"},
		ISBN:        stringPtr("9780141439518"),
		AddMode:     "MANUAL",
		BookshelfID: 1,
	},
	{
		ID:          "5",
		Title:       "The Hobbit",
		Authors:     []string{"J. R. R. Tolkien"},
		ISBN:        stringPtr("9780547928227"),
		AddMode:     "MANUAL",
		BookshelfID: 1,
	},
	{
		ID:          "6",
		Title:       "The Catcher in the Rye",
		Authors:     []string{"J. D. Salinger"},
		ISBN:        stringPtr("9780316769488"),
		AddMode:     "MANUAL",
		BookshelfID: 1,
	},
}

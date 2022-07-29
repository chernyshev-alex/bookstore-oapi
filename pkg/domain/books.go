package domain

type Book struct {
	Author      string
	AuthorId    int
	BookId      int
	Publisher   string
	PublisherId int
	Title       string
	Year        int
	Descr       string
	Isbn        string
}

type Books = []Book

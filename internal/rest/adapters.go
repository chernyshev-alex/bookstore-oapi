package rest

import (
	"github.com/chernyshev-alex/go-bookstore-oapi/internal/models"
	"github.com/volatiletech/null/v8"
)

func addBookRequestToDomain(rq AddBookRequest) models.Book {
	return models.Book{
		ID:          null.Int64From(int64(rq.Id)),
		Authorid:    int64(rq.AuthorId),
		Bookid:      int64(rq.Id),
		Publisher:   null.StringFrom(rq.Publisher),
		Publisherid: int64(rq.PublisherId),
		Title:       null.StringFrom(rq.Title),
		Year:        int64(rq.Year),
		Descr:       null.StringFrom(rq.Descr),
		Isbn:        rq.Isbn,
	}
}

func toJsonBook(b models.Book) BookJson {
	return BookJson{
		AuthorId:    int(b.Authorid),
		Descr:       b.Descr.String,
		Id:          int(b.Bookid),
		Isbn:        b.Isbn,
		Publisher:   b.Publisher.String,
		PublisherId: int(b.Publisherid),
		Title:       b.Title.String,
		Year:        int(b.Year),
	}
}

func sliceToJsonBooks(books []*models.Book) []BookJson {
	var bj = make([]BookJson, 0)
	for _, b := range books {
		bj = append(bj, toJsonBook(*b))
	}
	return bj
}

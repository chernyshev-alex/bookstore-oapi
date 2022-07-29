package rest

import "github.com/chernyshev-alex/go-bookstore-oapi/pkg/domain"

func addBookRequestToDomain(rq *AddBookRequest) domain.Book {
	return domain.Book{
		Author:      rq.Author,
		AuthorId:    rq.AuthorId,
		BookId:      rq.Id,
		Publisher:   rq.Publisher,
		PublisherId: rq.PublisherId,
		Title:       rq.Title,
		Year:        rq.Year,
	}
}

func domainToJsonBook(b *domain.Book) BookJson {
	return BookJson{
		Author:      b.Author,
		AuthorId:    b.AuthorId,
		Descr:       b.Descr,
		Id:          b.BookId,
		Isbn:        b.Isbn,
		Publisher:   b.Publisher,
		PublisherId: b.PublisherId,
		Title:       b.Title,
		Year:        b.Year,
	}
}

func domainSliceToJsonBooks(books []domain.Book) []BookJson {
	var bj = make([]BookJson, 0)
	for _, b := range books {
		bj = append(bj, domainToJsonBook(&b))
	}
	return bj
}

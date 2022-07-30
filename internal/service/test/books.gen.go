// Code generated by counterfeiter. DO NOT EDIT.
package test

import (
	"context"
	"sync"

	"github.com/chernyshev-alex/go-bookstore-oapi/internal/models"
	"github.com/chernyshev-alex/go-bookstore-oapi/internal/service"
)

type FakeBooksService struct {
	AddBookStub        func(context.Context, models.Book) (*models.Book, error)
	addBookMutex       sync.RWMutex
	addBookArgsForCall []struct {
		arg1 context.Context
		arg2 models.Book
	}
	addBookReturns struct {
		result1 *models.Book
		result2 error
	}
	addBookReturnsOnCall map[int]struct {
		result1 *models.Book
		result2 error
	}
	DeleteBookStub        func(context.Context, string) error
	deleteBookMutex       sync.RWMutex
	deleteBookArgsForCall []struct {
		arg1 context.Context
		arg2 string
	}
	deleteBookReturns struct {
		result1 error
	}
	deleteBookReturnsOnCall map[int]struct {
		result1 error
	}
	FindBooksByAuthorStub        func(context.Context, string) ([]*models.Book, error)
	findBooksByAuthorMutex       sync.RWMutex
	findBooksByAuthorArgsForCall []struct {
		arg1 context.Context
		arg2 string
	}
	findBooksByAuthorReturns struct {
		result1 []*models.Book
		result2 error
	}
	findBooksByAuthorReturnsOnCall map[int]struct {
		result1 []*models.Book
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeBooksService) AddBook(arg1 context.Context, arg2 models.Book) (*models.Book, error) {
	fake.addBookMutex.Lock()
	ret, specificReturn := fake.addBookReturnsOnCall[len(fake.addBookArgsForCall)]
	fake.addBookArgsForCall = append(fake.addBookArgsForCall, struct {
		arg1 context.Context
		arg2 models.Book
	}{arg1, arg2})
	stub := fake.AddBookStub
	fakeReturns := fake.addBookReturns
	fake.recordInvocation("AddBook", []interface{}{arg1, arg2})
	fake.addBookMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeBooksService) AddBookCallCount() int {
	fake.addBookMutex.RLock()
	defer fake.addBookMutex.RUnlock()
	return len(fake.addBookArgsForCall)
}

func (fake *FakeBooksService) AddBookCalls(stub func(context.Context, models.Book) (*models.Book, error)) {
	fake.addBookMutex.Lock()
	defer fake.addBookMutex.Unlock()
	fake.AddBookStub = stub
}

func (fake *FakeBooksService) AddBookArgsForCall(i int) (context.Context, models.Book) {
	fake.addBookMutex.RLock()
	defer fake.addBookMutex.RUnlock()
	argsForCall := fake.addBookArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeBooksService) AddBookReturns(result1 *models.Book, result2 error) {
	fake.addBookMutex.Lock()
	defer fake.addBookMutex.Unlock()
	fake.AddBookStub = nil
	fake.addBookReturns = struct {
		result1 *models.Book
		result2 error
	}{result1, result2}
}

func (fake *FakeBooksService) AddBookReturnsOnCall(i int, result1 *models.Book, result2 error) {
	fake.addBookMutex.Lock()
	defer fake.addBookMutex.Unlock()
	fake.AddBookStub = nil
	if fake.addBookReturnsOnCall == nil {
		fake.addBookReturnsOnCall = make(map[int]struct {
			result1 *models.Book
			result2 error
		})
	}
	fake.addBookReturnsOnCall[i] = struct {
		result1 *models.Book
		result2 error
	}{result1, result2}
}

func (fake *FakeBooksService) DeleteBook(arg1 context.Context, arg2 string) error {
	fake.deleteBookMutex.Lock()
	ret, specificReturn := fake.deleteBookReturnsOnCall[len(fake.deleteBookArgsForCall)]
	fake.deleteBookArgsForCall = append(fake.deleteBookArgsForCall, struct {
		arg1 context.Context
		arg2 string
	}{arg1, arg2})
	stub := fake.DeleteBookStub
	fakeReturns := fake.deleteBookReturns
	fake.recordInvocation("DeleteBook", []interface{}{arg1, arg2})
	fake.deleteBookMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeBooksService) DeleteBookCallCount() int {
	fake.deleteBookMutex.RLock()
	defer fake.deleteBookMutex.RUnlock()
	return len(fake.deleteBookArgsForCall)
}

func (fake *FakeBooksService) DeleteBookCalls(stub func(context.Context, string) error) {
	fake.deleteBookMutex.Lock()
	defer fake.deleteBookMutex.Unlock()
	fake.DeleteBookStub = stub
}

func (fake *FakeBooksService) DeleteBookArgsForCall(i int) (context.Context, string) {
	fake.deleteBookMutex.RLock()
	defer fake.deleteBookMutex.RUnlock()
	argsForCall := fake.deleteBookArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeBooksService) DeleteBookReturns(result1 error) {
	fake.deleteBookMutex.Lock()
	defer fake.deleteBookMutex.Unlock()
	fake.DeleteBookStub = nil
	fake.deleteBookReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeBooksService) DeleteBookReturnsOnCall(i int, result1 error) {
	fake.deleteBookMutex.Lock()
	defer fake.deleteBookMutex.Unlock()
	fake.DeleteBookStub = nil
	if fake.deleteBookReturnsOnCall == nil {
		fake.deleteBookReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.deleteBookReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeBooksService) FindBooksByAuthor(arg1 context.Context, arg2 string) ([]*models.Book, error) {
	fake.findBooksByAuthorMutex.Lock()
	ret, specificReturn := fake.findBooksByAuthorReturnsOnCall[len(fake.findBooksByAuthorArgsForCall)]
	fake.findBooksByAuthorArgsForCall = append(fake.findBooksByAuthorArgsForCall, struct {
		arg1 context.Context
		arg2 string
	}{arg1, arg2})
	stub := fake.FindBooksByAuthorStub
	fakeReturns := fake.findBooksByAuthorReturns
	fake.recordInvocation("FindBooksByAuthor", []interface{}{arg1, arg2})
	fake.findBooksByAuthorMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeBooksService) FindBooksByAuthorCallCount() int {
	fake.findBooksByAuthorMutex.RLock()
	defer fake.findBooksByAuthorMutex.RUnlock()
	return len(fake.findBooksByAuthorArgsForCall)
}

func (fake *FakeBooksService) FindBooksByAuthorCalls(stub func(context.Context, string) ([]*models.Book, error)) {
	fake.findBooksByAuthorMutex.Lock()
	defer fake.findBooksByAuthorMutex.Unlock()
	fake.FindBooksByAuthorStub = stub
}

func (fake *FakeBooksService) FindBooksByAuthorArgsForCall(i int) (context.Context, string) {
	fake.findBooksByAuthorMutex.RLock()
	defer fake.findBooksByAuthorMutex.RUnlock()
	argsForCall := fake.findBooksByAuthorArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeBooksService) FindBooksByAuthorReturns(result1 []*models.Book, result2 error) {
	fake.findBooksByAuthorMutex.Lock()
	defer fake.findBooksByAuthorMutex.Unlock()
	fake.FindBooksByAuthorStub = nil
	fake.findBooksByAuthorReturns = struct {
		result1 []*models.Book
		result2 error
	}{result1, result2}
}

func (fake *FakeBooksService) FindBooksByAuthorReturnsOnCall(i int, result1 []*models.Book, result2 error) {
	fake.findBooksByAuthorMutex.Lock()
	defer fake.findBooksByAuthorMutex.Unlock()
	fake.FindBooksByAuthorStub = nil
	if fake.findBooksByAuthorReturnsOnCall == nil {
		fake.findBooksByAuthorReturnsOnCall = make(map[int]struct {
			result1 []*models.Book
			result2 error
		})
	}
	fake.findBooksByAuthorReturnsOnCall[i] = struct {
		result1 []*models.Book
		result2 error
	}{result1, result2}
}

func (fake *FakeBooksService) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.addBookMutex.RLock()
	defer fake.addBookMutex.RUnlock()
	fake.deleteBookMutex.RLock()
	defer fake.deleteBookMutex.RUnlock()
	fake.findBooksByAuthorMutex.RLock()
	defer fake.findBooksByAuthorMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeBooksService) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ service.BooksService = new(FakeBooksService)

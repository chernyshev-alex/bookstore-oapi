// Code generated by counterfeiter. DO NOT EDIT.
package test

import (
	"context"
	"sync"

	"github.com/chernyshev-alex/go-bookstore-oapi/internal/repo"
	"github.com/chernyshev-alex/go-bookstore-oapi/pkg/domain"
)

type FakeBooksCrudRepository struct {
	AddBookStub        func(context.Context, domain.Book) (domain.Book, error)
	addBookMutex       sync.RWMutex
	addBookArgsForCall []struct {
		arg1 context.Context
		arg2 domain.Book
	}
	addBookReturns struct {
		result1 domain.Book
		result2 error
	}
	addBookReturnsOnCall map[int]struct {
		result1 domain.Book
		result2 error
	}
	DeleteBookStub        func(context.Context, int) error
	deleteBookMutex       sync.RWMutex
	deleteBookArgsForCall []struct {
		arg1 context.Context
		arg2 int
	}
	deleteBookReturns struct {
		result1 error
	}
	deleteBookReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeBooksCrudRepository) AddBook(arg1 context.Context, arg2 domain.Book) (domain.Book, error) {
	fake.addBookMutex.Lock()
	ret, specificReturn := fake.addBookReturnsOnCall[len(fake.addBookArgsForCall)]
	fake.addBookArgsForCall = append(fake.addBookArgsForCall, struct {
		arg1 context.Context
		arg2 domain.Book
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

func (fake *FakeBooksCrudRepository) AddBookCallCount() int {
	fake.addBookMutex.RLock()
	defer fake.addBookMutex.RUnlock()
	return len(fake.addBookArgsForCall)
}

func (fake *FakeBooksCrudRepository) AddBookCalls(stub func(context.Context, domain.Book) (domain.Book, error)) {
	fake.addBookMutex.Lock()
	defer fake.addBookMutex.Unlock()
	fake.AddBookStub = stub
}

func (fake *FakeBooksCrudRepository) AddBookArgsForCall(i int) (context.Context, domain.Book) {
	fake.addBookMutex.RLock()
	defer fake.addBookMutex.RUnlock()
	argsForCall := fake.addBookArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeBooksCrudRepository) AddBookReturns(result1 domain.Book, result2 error) {
	fake.addBookMutex.Lock()
	defer fake.addBookMutex.Unlock()
	fake.AddBookStub = nil
	fake.addBookReturns = struct {
		result1 domain.Book
		result2 error
	}{result1, result2}
}

func (fake *FakeBooksCrudRepository) AddBookReturnsOnCall(i int, result1 domain.Book, result2 error) {
	fake.addBookMutex.Lock()
	defer fake.addBookMutex.Unlock()
	fake.AddBookStub = nil
	if fake.addBookReturnsOnCall == nil {
		fake.addBookReturnsOnCall = make(map[int]struct {
			result1 domain.Book
			result2 error
		})
	}
	fake.addBookReturnsOnCall[i] = struct {
		result1 domain.Book
		result2 error
	}{result1, result2}
}

func (fake *FakeBooksCrudRepository) DeleteBook(arg1 context.Context, arg2 int) error {
	fake.deleteBookMutex.Lock()
	ret, specificReturn := fake.deleteBookReturnsOnCall[len(fake.deleteBookArgsForCall)]
	fake.deleteBookArgsForCall = append(fake.deleteBookArgsForCall, struct {
		arg1 context.Context
		arg2 int
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

func (fake *FakeBooksCrudRepository) DeleteBookCallCount() int {
	fake.deleteBookMutex.RLock()
	defer fake.deleteBookMutex.RUnlock()
	return len(fake.deleteBookArgsForCall)
}

func (fake *FakeBooksCrudRepository) DeleteBookCalls(stub func(context.Context, int) error) {
	fake.deleteBookMutex.Lock()
	defer fake.deleteBookMutex.Unlock()
	fake.DeleteBookStub = stub
}

func (fake *FakeBooksCrudRepository) DeleteBookArgsForCall(i int) (context.Context, int) {
	fake.deleteBookMutex.RLock()
	defer fake.deleteBookMutex.RUnlock()
	argsForCall := fake.deleteBookArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeBooksCrudRepository) DeleteBookReturns(result1 error) {
	fake.deleteBookMutex.Lock()
	defer fake.deleteBookMutex.Unlock()
	fake.DeleteBookStub = nil
	fake.deleteBookReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeBooksCrudRepository) DeleteBookReturnsOnCall(i int, result1 error) {
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

func (fake *FakeBooksCrudRepository) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.addBookMutex.RLock()
	defer fake.addBookMutex.RUnlock()
	fake.deleteBookMutex.RLock()
	defer fake.deleteBookMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeBooksCrudRepository) recordInvocation(key string, args []interface{}) {
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

var _ repo.BooksCrudRepository = new(FakeBooksCrudRepository)

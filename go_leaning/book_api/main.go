package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
)

type Book struct {
	ID     int
	Title  string
	Author string
}

type Store struct {
	sync.RWMutex
	books  []Book
	nextID int
}

const dataFile = "book.json"

func main() {
	fmt.Println("book_api started")
	books, err := loadBook()
	if err != nil {
		fmt.Println("读取数据失败", err)
		return
	}
	store := newStore(books)

	// testConcurrency(store)
	mux := http.NewServeMux()
	mux.HandleFunc("/books", booksHandler(store))
	mux.HandleFunc("/books/", bookByIDHandler(store))

	log.Println("book_api started on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}

}

func newStore(books []Book) *Store {
	return &Store{
		books:  books,
		nextID: nextID(books),
	}
}
func loadBook() ([]Book, error) {
	data, err := os.ReadFile(dataFile)
	if err != nil {
		if os.IsNotExist(err) {
			return []Book{}, nil
		}
		return nil, err
	}
	if len(data) == 0 {
		return []Book{}, nil
	}

	var books []Book
	if err := json.Unmarshal(data, &books); err != nil {
		return nil, err
	}
	return books, nil
}

func saveBook(books []Book) error {
	data, err := json.MarshalIndent(books, "", " ")
	if err != nil {
		return err
	}
	return os.WriteFile(dataFile, data, 0644)
}

func nextID(books []Book) int {
	maxID := 0
	for _, book := range books {
		if book.ID > maxID {
			maxID = book.ID
		}
	}
	return maxID + 1
}

func (s *Store) Add(title, author string) Book {
	s.Lock()
	defer s.Unlock()

	book := Book{
		ID:     s.nextID,
		Title:  title,
		Author: author,
	}
	s.books = append(s.books, book)
	s.nextID++
	saveBook(s.books)
	return book
}

func (s *Store) List() []Book {
	s.RLock()
	defer s.RUnlock()
	list := make([]Book, len(s.books))
	copy(list, s.books)

	return list
}

func (s *Store) Get(id int) (Book, bool) {
	s.Lock()
	defer s.Unlock()
	for i := range s.books {
		if s.books[i].ID == id {
			return s.books[i], true
		}
	}
	return Book{}, false
}

func (s *Store) Delete(id int) bool {
	s.Lock()
	defer s.Unlock()
	for i := range s.books {
		if s.books[i].ID == id {
			s.books = append(s.books[:i], s.books[i+1:]...)
			saveBook(s.books)
			return true
		}
	}

	return false
}

func (s *Store) Update(id int, title string, author string) (Book, bool) {
	s.Lock()
	defer s.Unlock()
	for i := range s.books {
		if s.books[i].ID == id {
			s.books[i].Title = title
			s.books[i].Author = author
			saveBook(s.books)
			return s.books[i], true
		}
	}

	return Book{}, false
}

func booksHandler(store *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			if err := json.NewEncoder(w).Encode(store.List()); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		case http.MethodPost:
			var req struct {
				Title  string `json:"title"`
				Author string `json:"author"`
			}
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, "json不对", http.StatusBadRequest)
				return
			}
			req.Title = strings.TrimSpace(req.Title)
			if req.Title == "" {
				http.Error(w, "标题为空", http.StatusBadRequest)
				return
			}
			book := store.Add(req.Title, req.Author)
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			json.NewEncoder(w).Encode(book)

		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

	}
}

func bookByIDHandler(store *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idstr := strings.TrimPrefix(r.URL.Path, "/books/")
		if idstr == "" {
			http.Error(w, "错误id", http.StatusBadRequest)
			return
		}
		id, err := strconv.Atoi(idstr)
		if err != nil || id <= 0 {
			http.Error(w, "错误id", http.StatusBadRequest)
			return
		}
		switch r.Method {
		case http.MethodGet:
			book, ok := store.Get(id)
			if !ok {
				http.Error(w, "没找到书", http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "application/json;charset=utf-8")
			json.NewEncoder(w).Encode(book)
		case http.MethodPut:
			var req struct {
				Title  string `json:"title"`
				Author string `json:"author"`
			}
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, "jaon不对", http.StatusBadRequest)
				return
			}
			req.Title = strings.TrimSpace(req.Title)
			if req.Title == "" {
				http.Error(w, "标题为空", http.StatusBadRequest)
				return
			}
			updated, ok := store.Update(id, req.Title, req.Author)
			if !ok {
				http.Error(w, "title required", http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			json.NewEncoder(w).Encode(updated)
		case http.MethodDelete:
			if !store.Delete(id) {
				http.Error(w, "没找到书", http.StatusNotFound)
				return
			}
			w.WriteHeader(http.StatusNoContent)

		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}

	}
}

func testConcurrency(store *Store) {
	var wg sync.WaitGroup
	const count = 100 // 模拟 100 个并发操作

	// 1. 测试并发写入 (Add)
	for i := 0; i < count; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			store.Add(fmt.Sprintf("Book %d", n), "Author")
		}(i)
	}

	// 2. 测试并发读取 (List)
	for i := 0; i < count; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = store.List()
		}()
	}

	wg.Wait()
	fmt.Printf("并发测试完成，当前总数: %d\n", len(store.List()))
}

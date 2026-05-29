package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

type Book struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

type Store struct {
	db *sql.DB
}

const dataFile = "book.json"

func main() {
	fmt.Println("book_api started")
	db, err := openDB()
	if err != nil {
		fmt.Println("打开数据库失败", err)
		return
	}
	defer db.Close()
	store := newStore(db)

	// testConcurrency(store)
	mux := http.NewServeMux()
	mux.HandleFunc("/books", booksHandler(store))
	mux.HandleFunc("/books/", bookByIDHandler(store))

	log.Println("book_api started on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}

}

func newStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) Add(title, author string) (Book, error) {
	result, err := s.db.Exec(`INSERT INTO books (title,author) VALUES(?,?)`, title, author)
	if err != nil {
		return Book{}, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return Book{}, err
	}
	return Book{ID: int(id), Title: title, Author: author}, nil
}

func (s *Store) List() ([]Book, error) {
	rows, err := s.db.Query(`SELECT id, title, author FROM books ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	books := make([]Book, 0)
	for rows.Next() {
		var book Book
		if err := rows.Scan(&book.ID, &book.Title, &book.Author); err != nil {
			return nil, err
		}
		books = append(books, book)
	}

	return books, rows.Err()
}

func (s *Store) Get(id int) (Book, bool, error) {
	var book Book
	err := s.db.QueryRow(`SELECT id ,title,author FROM books WHERE id = ?`, id).Scan(&book.ID, &book.Title, &book.Author)
	if err != nil {
		if err == sql.ErrNoRows {
			return Book{}, false, nil
		}
		return Book{}, false, err
	}

	return book, true, nil
}

func (s *Store) Delete(id int) (bool, error) {
	result, err := s.db.Exec(`DELETE FROM books WHERE id = ?`, id)
	if err != nil {
		return false, err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (s *Store) Update(id int, title string, author string) (Book, bool, error) {
	result, err := s.db.Exec(`UPDATE books SET title = ?, author = ? WHERE id = ?`, title, author, id)
	if err != nil {
		return Book{}, false, err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return Book{}, false, err
	}
	if count == 0 {
		return Book{}, false, nil
	}
	return Book{ID: id, Title: title, Author: author}, true, nil
}

func booksHandler(store *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			books, err := store.List()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			if err := json.NewEncoder(w).Encode(books); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		case http.MethodPost:
			var req struct {
				Title  string `json:"title"`
				Author string `json:"author"`
			}
			defer r.Body.Close()
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, "json不对", http.StatusBadRequest)
				return
			}
			req.Title = strings.TrimSpace(req.Title)
			if req.Title == "" {
				http.Error(w, "标题为空", http.StatusBadRequest)
				return
			}
			books, err := store.Add(req.Title, req.Author)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			if err := json.NewEncoder(w).Encode(books); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

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
			book, ok, err := store.Get(id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if !ok {
				http.Error(w, "没找到书", http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			if err := json.NewEncoder(w).Encode(book); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		case http.MethodPut:
			var req struct {
				Title  string `json:"title"`
				Author string `json:"author"`
			}
			defer r.Body.Close()
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, "json不对", http.StatusBadRequest)
				return
			}
			req.Title = strings.TrimSpace(req.Title)
			if req.Title == "" {
				http.Error(w, "标题为空", http.StatusBadRequest)
				return
			}
			updated, ok, err := store.Update(id, req.Title, req.Author)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if !ok {
				http.Error(w, "没找到书", http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			if err := json.NewEncoder(w).Encode(updated); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		case http.MethodDelete:
			isDelete, err := store.Delete(id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if !isDelete {
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
			_, _ = store.List()
		}()
	}

	wg.Wait()
	books, _ := store.List()
	fmt.Printf("并发测试完成，当前总数: %d\n", len(books))
}

package data

import (
	"database/sql"
	"errors"
	"github.com/Bug-daulet/FinalSPA/internal/validator"
	"github.com/lib/pq"
	"time"
)

type BookModel struct {
	DB *sql.DB
}

func (b BookModel) Insert(book *Book) error {
	query := `INSERT INTO books (title, year, pages, genres)
				VALUES ($1, $2, $3, $4)
				RETURNING id, created_at, version`

	args := []interface{}{book.Title, book.Year, book.Pages, pq.Array(book.Genres)}

	return b.DB.QueryRow(query, args...).Scan(&book.ID, &book.CreatedAt, &book.Version)

}

func (b BookModel) Get(id int64) (*Book, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	query := `SELECT id, created_at, title, year, pages, genres, version
				FROM books
				WHERE id = $1`

	var book Book

	err := b.DB.QueryRow(query, id).Scan(
		&book.ID,
		&book.CreatedAt,
		&book.Title,
		&book.Year,
		&book.Pages,
		pq.Array(&book.Genres),
		&book.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &book, nil

}
func (b BookModel) Update(book *Book) error {
	query := `UPDATE books
				SET title = $1, year = $2, pages = $3, genres = $4, version = version + 1
				WHERE id = $5
				RETURNING version`

	args := []interface{}{
		book.Title,
		book.Year,
		book.Pages,
		pq.Array(book.Genres),
		book.ID,
	}

	return b.DB.QueryRow(query, args...).Scan(&book.Version)
}

func (b BookModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `DELETE FROM books
				WHERE id = $1`

	result, err := b.DB.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil

}

type Book struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"`
	Title     string    `json:"title"`
	Year      int32     `json:"year,omitempty"`
	Pages     Pages     `json:"pages,omitempty"`
	Genres    []string  `json:"genres,omitempty"`
	Version   int32     `json:"version"`
}

func ValidateBooks(v *validator.Validator, input *Book) {
	v.Check(input.Title != "", "title", "must be provided")
	v.Check(len(input.Title) <= 500, "title", "must not be more than 500 bytes long")

	v.Check(input.Year != 0, "year", "must be provided")
	v.Check(input.Year >= 1888, "year", "must be greater than 1888")
	v.Check(input.Year <= int32(time.Now().Year()), "year", "must not be in the future")

	v.Check(input.Pages != 0, "pages", "must be provided")
	v.Check(input.Pages > 0, "pages", "must be a positive integer")

	v.Check(input.Genres != nil, "genres", "must be provided")
	v.Check(len(input.Genres) >= 1, "genres", "must contain at least 1 genre")
	v.Check(len(input.Genres) <= 5, "genres", "must not contain more than 5 genres")

	v.Check(validator.Unique(input.Genres), "genres", "must not contain duplicate values")
}

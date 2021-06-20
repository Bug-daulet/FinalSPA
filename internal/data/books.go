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
	// Create an args slice containing the values for the placeholder parameters.
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

//
//type ComicsModel struct {
//	DB *sql.DB
//}
//
//// Add a placeholder method for inserting a new record in the movies table.
//func (m ComicsModel) Insert(comics *Comics) error {
//	query := `INSERT INTO comics (title, year, pages)
//			VALUES ($1, $2, $3)
//			RETURNING id, created_at, version`
//
//	args := []interface{}{comics.Title, comics.Year, comics.Pages}
//
//	// Create a context with a 3-second timeout.
//	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
//	defer cancel()
//
//	return m.DB.QueryRowContext(ctx, query, args...).Scan(&comics.ID, &comics.CreatedAt, &comics.Version)
//}
//
//// Add a placeholder method for fetching a specific record from the movies table.
//func (m ComicsModel) Get(id int64) (*Comics, error) {
//	if id < 1 {
//		return nil, ErrRecordNotFound
//	}
//	// Define the SQL query for retrieving the comics data.
//	query := `SELECT id, created_at, title, year, pages, version
//			FROM comics
//			WHERE id = $1`
//	// Declare a Movie struct to hold the data returned by the query.
//	var comics Comics
//
//	// Use the context.WithTimeout() function to create a context.Context which carries a
//	// 3-second timeout deadline. Note that we're using the empty context.Background()
//	// as the 'parent' context.
//	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
//
//	// Importantly, use defer to make sure that we cancel the context before the Get()
//	// method returns.
//	defer cancel()
//
//	err := m.DB.QueryRowContext(ctx, query, id).Scan(
//		&comics.ID,
//		&comics.CreatedAt,
//		&comics.Title,
//		&comics.Year,
//		&comics.Pages,
//		&comics.Version,
//	)
//	// Handle any errors. If there was no matching comics found, Scan() will return
//	// a sql.ErrNoRows error. We check for this and return our custom ErrRecordNotFound
//	// error instead.
//	if err != nil {
//		switch {
//		case errors.Is(err, sql.ErrNoRows):
//			return nil, ErrRecordNotFound
//		default:
//			return nil, err
//		}
//	}
//	// Otherwise, return a pointer to the Movie struct.
//	return &comics, nil
//}
//
//// Add a placeholder method for updating a specific record in the movies table.
//func (m ComicsModel) Update(comics *Comics) error {
//	query := `UPDATE comics
//			SET title = $1, year = $2, pages = $3, version = version + 1
//			WHERE id = $4 AND version = $5
//			RETURNING version`
//	// Create an args slice containing the values for the placeholder parameters.
//	args := []interface{}{
//		comics.Title,
//		comics.Year,
//		comics.Pages,
//		comics.ID,
//		comics.Version,
//	}
//
//	// Create a context with a 3-second timeout.
//	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
//	defer cancel()
//
//	// Execute the SQL query. If no matching row could be found, we know the movie
//	// version has changed (or the record has been deleted) and we return our custom
//	// ErrEditConflict error.
//	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&comics.Version)
//	if err != nil {
//		switch {
//		case errors.Is(err, sql.ErrNoRows):
//			return ErrEditConflict
//		default:
//			return err
//		}
//	}
//	return nil
//}
//
//// Add a placeholder method for deleting a specific record from the movies table.
//func (m ComicsModel) Delete(id int64) error {
//	// Return an ErrRecordNotFound error if the movie ID is less than 1.
//	if id < 1 {
//		return ErrRecordNotFound
//	}
//	// Construct the SQL query to delete the record.
//	query := `DELETE FROM comics
//			WHERE id = $1`
//
//	// Create a context with a 3-second timeout.
//	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
//	defer cancel()
//
//	// Execute the SQL query using the Exec() method, passing in the id variable as
//	// the value for the placeholder parameter. The Exec() method returns a sql.Result
//	// object.
//	result, err := m.DB.ExecContext(ctx, query, id)
//	if err != nil {
//		return err
//	}
//	// Call the RowsAffected() method on the sql.Result object to get the number of rows
//	// affected by the query.
//	rowsAffected, err := result.RowsAffected()
//	if err != nil {
//		return err
//	}
//	// If no rows were affected, we know that the movies table didn't contain a record
//	// with the provided ID at the moment we tried to delete it. In that case we
//	// return an ErrRecordNotFound error.
//	if rowsAffected == 0 {
//		return ErrRecordNotFound
//	}
//	return nil
//}
//
//func (m ComicsModel) GetAll() (*[]Comics, error) {
//	// Define the SQL query for retrieving the comics data.
//	query := `SELECT id, created_at, title, year, pages, version
//			FROM comics`
//	// Declare a Movie struct to hold the data returned by the query.
//	var comics Comics
//	var comicses []Comics
//	rows, err := m.DB.Query(query)
//	defer rows.Close()
//	// Handle any errors. If there was no matching comics found, Scan() will return
//	// a sql.ErrNoRows error. We check for this and return our custom ErrRecordNotFound
//	// error instead.
//	if err != nil {
//		switch {
//		case errors.Is(err, sql.ErrNoRows):
//			return nil, ErrRecordNotFound
//		default:
//			return nil, err
//		}
//	}
//	for rows.Next() {
//		err := rows.Scan(
//			&comics.ID,
//			&comics.CreatedAt,
//			&comics.Title,
//			&comics.Year,
//			&comics.Pages,
//			&comics.Version,
//		)
//		if err != nil {
//			log.Fatal(err)
//		}
//
//		comicses = append(comicses, Comics{
//			ID:        comics.ID,
//			CreatedAt: comics.CreatedAt,
//			Title:     comics.Title,
//			Year:      comics.Year,
//			Pages:     comics.Pages,
//			Version:   comics.Version,
//		})
//	}
//	if err := rows.Err(); err != nil {
//		log.Fatal(err)
//	}
//	return &comicses, nil
//}

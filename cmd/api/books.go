package main

import (
	"errors"
	"fmt"
	"github.com/Bug-daulet/FinalSPA/internal/data"
	"github.com/Bug-daulet/FinalSPA/internal/validator"
	"net/http"
)

func (app *application) createBookHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Title  string     `json:"title"`
		Year   int32      `json:"year"`
		Pages  data.Pages `json:"pages"`
		Genres []string   `json:"genres"`
	}

	err := app.readJSON(w, r, &input)

	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	book := &data.Book{
		Title:  input.Title,
		Year:   input.Year,
		Pages:  input.Pages,
		Genres: input.Genres,
	}

	v := validator.New()

	if data.ValidateBooks(v, book); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Books.Insert(book)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprint("/v1/books/%d", book.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"book": book}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	fmt.Fprintf(w, "%+v\n", input)
}

func (app *application) showBookHandler(w http.ResponseWriter, r *http.Request) {

	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	books, err := app.models.Books.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"books": books}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateBookHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	book, err := app.models.Books.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	var input struct {
		Title  string     `json:"title"`
		Year   int32      `json:"year"`
		Pages  data.Pages `json:"pages"`
		Genres []string   `json:"genres"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	book.Title = input.Title
	book.Year = input.Year
	book.Pages = input.Pages
	book.Genres = input.Genres

	v := validator.New()
	if data.ValidateBooks(v, book); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Books.Update(book)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"book": book}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

//func (app *application) deleteComicsHandler(w http.ResponseWriter, r *http.Request) {
//	// Extract the movie ID from the URL.
//	id, err := app.readIDParam(r)
//	if err != nil {
//		app.notFoundResponse(w, r)
//		return
//	}
//	// Delete the movie from the database, sending a 404 Not Found response to the
//	// client if there isn't a matching record.
//	err = app.models.Comics.Delete(id)
//	if err != nil {
//		switch {
//		case errors.Is(err, data.ErrRecordNotFound):
//			app.notFoundResponse(w, r)
//		default:
//			app.serverErrorResponse(w, r, err)
//		}
//		return
//	}
//	// Return a 200 OK status code along with a success message.
//	err = app.writeJSON(w, http.StatusOK, envelope{"message": "comics successfully deleted"}, nil)
//	if err != nil {
//		app.serverErrorResponse(w, r, err)
//	}
//}
//
//func (app *application) showAllComicsHandler(w http.ResponseWriter, r *http.Request) {
//
//	comics, err := app.models.Comics.GetAll()
//	if err != nil {
//		switch {
//		case errors.Is(err, data.ErrRecordNotFound):
//			app.notFoundResponse(w, r)
//		default:
//			app.serverErrorResponse(w, r, err)
//		}
//		return
//	}
//
//	err = app.writeJSON(w, http.StatusOK, envelope{"comics": comics}, nil)
//	if err != nil {
//		app.serverErrorResponse(w, r, err)
//	}
//}

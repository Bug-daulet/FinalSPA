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
		Title  *string     `json:"title"`
		Year   *int32      `json:"year"`
		Pages  *data.Pages `json:"pages"`
		Genres []string    `json:"genres"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Title != nil {
		book.Title = *input.Title
	}

	if input.Year != nil {
		book.Year = *input.Year
	}
	if input.Pages != nil {
		book.Pages = *input.Pages
	}
	if input.Genres != nil {
		book.Genres = input.Genres
	}

	v := validator.New()

	if data.ValidateBooks(v, book); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	err = app.models.Books.Update(book)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"book": book}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) deleteBookHandler(w http.ResponseWriter, r *http.Request) {

	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Books.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Return a 200 OK status code along with a success message.
	err = app.writeJSON(w, http.StatusOK, envelope{"message": "book successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

/*func (app *application) showAllComicsHandler(w http.ResponseWriter, r *http.Request) {

	comics, err := app.models.Books.GetAll()
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"comics": comics}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}*/

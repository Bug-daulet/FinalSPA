package main

import (
	"fmt"
	"github.com/Bug-daulet/FinalSPA/internal/data"
	"net/http"
	"time"
)


func (app *application) createBookHandler(w http.ResponseWriter, r *http.Request) {

	//var input struct {
	//	Title string     `json:"title"`
	//	Year  int32      `json:"year"`
	//	Pages data.Pages `json:"pages"`
	//}
	//
	//err := app.readJSON(w, r, &input)
	//if err != nil {
	//	app.badRequestResponse(w, r, err)
	//	return
	//}
	//
	//comics := &data.Comics{
	//	Title: input.Title,
	//	Year:  input.Year,
	//	Pages: input.Pages,
	//}
	//
	//v := validator.New()
	//
	//if data.ValidateComics(v, comics); !v.Valid() {
	//	app.failedValidationResponse(w, r, v.Errors)
	//	return
	//}
	//
	//err = app.models.Comics.Insert(comics)
	//if err != nil {
	//	app.serverErrorResponse(w, r, err)
	//	return
	//}
	//
	//headers := make(http.Header)
	//headers.Set("Location", fmt.Sprintf("/v1/comics/%d", comics.ID))
	//
	//err = app.writeJSON(w, http.StatusCreated, envelope{"comics": comics}, headers)
	//if err != nil {
	//	app.serverErrorResponse(w, r, err)
	//}
	//
	//fmt.Fprintf(w, "%+v\n", input)

	fmt.Fprintln(w, "create a new book")

}

func (app *application) showBookHandler(w http.ResponseWriter, r *http.Request) {

	//id, err := app.readIDParam(r)
	//if err != nil {
	//	app.notFoundResponse(w, r)
	//	return
	//}
	//
	//comics, err := app.models.Comics.Get(id)
	//if err != nil {
	//	switch {
	//	case errors.Is(err, data.ErrRecordNotFound):
	//		app.notFoundResponse(w, r)
	//	default:
	//		app.serverErrorResponse(w, r, err)
	//	}
	//	return
	//}
	//
	//err = app.writeJSON(w, http.StatusOK, envelope{"comics": comics}, nil)
	//if err != nil {
	//	app.serverErrorResponse(w, r, err)
	//}


	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	book := data.Book{
		ID: id,
		CreatedAt: time.Now(),
		Title: "Casablanca",
		Pages: 102,
		Genres: []string{"drama", "romance", "war"},
		Version: 1,
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"book": book}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

//func (app *application) updateComicsHandler(w http.ResponseWriter, r *http.Request) {
//	// Extract the comics ID from the URL.
//	id, err := app.readIDParam(r)
//	if err != nil {
//		app.notFoundResponse(w, r)
//		return
//	}
//	// Fetch the existing comics record from the database, sending a 404 Not Found
//	// response to the client if we couldn't find a matching record.
//	comics, err := app.models.Comics.Get(id)
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
//	// If the request contains a X-Expected-Version header, verify that the movie
//	// version in the database matches the expected version specified in the header.
//	if r.Header.Get("X-Expected-Version") != "" {
//		if strconv.FormatInt(int64(comics.Version), 32) != r.Header.Get("X-Expected-Version") {
//			app.editConflictResponse(w, r)
//			return
//		}
//	}
//
//	// Declare an input struct to hold the expected data from the client.
//	var input struct {
//		Title *string     `json:"title"`
//		Year  *int32      `json:"year"`
//		Pages *data.Pages `json:"pages"`
//	}
//	// Read the JSON request body data into the input struct.
//	err = app.readJSON(w, r, &input)
//	if err != nil {
//		app.badRequestResponse(w, r, err)
//		return
//	}
//	// Copy the values from the request body to the appropriate fields of the movie
//	// record.
//	if input.Title != nil {
//		comics.Title = *input.Title
//	}
//	if input.Year != nil {
//		comics.Year = *input.Year
//	}
//	if input.Pages != nil {
//		comics.Pages = *input.Pages
//	}
//	// Validate the updated movie record, sending the client a 422 Unprocessable Entity
//	// response if any checks fail.
//	v := validator.New()
//	if data.ValidateComics(v, comics); !v.Valid() {
//		app.failedValidationResponse(w, r, v.Errors)
//		return
//	}
//	// Pass the updated movie record to our new Update() method.
//	// Intercept any ErrEditConflict error and call the new editConflictResponse()
//	// helper.
//	err = app.models.Comics.Update(comics)
//	if err != nil {
//		switch {
//		case errors.Is(err, data.ErrEditConflict):
//			app.editConflictResponse(w, r)
//		default:
//			app.serverErrorResponse(w, r, err)
//		}
//		return
//
//	}
//	// Write the updated movie record in a JSON response.
//	err = app.writeJSON(w, http.StatusOK, envelope{"comics": comics}, nil)
//	if err != nil {
//		app.serverErrorResponse(w, r, err)
//	}
//}
//
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

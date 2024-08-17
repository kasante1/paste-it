 
package main

import (
"fmt"
"net/http"
"strconv"
// "html/template" 
// "log"
"strings"
"unicode/utf8" 
"errors"
"github.com/kasante1/paste_it_backend/internal/models"
"github.com/julienschmidt/httprouter" 
)


func (app *application) home(w http.ResponseWriter, r *http.Request) {
	

	// files := []string{
	// 	"./ui/html/base.tmpl",
	// 	"./ui/html/partials/nav.tmpl",
	// 	"./ui/html/pages/home.tmpl",
	// 	}

	// ts, err := template.ParseFiles(files...)
	
	// if err != nil {
	// 	log.Println(err.Error())
	// 	app.serverError(w, err)
	// 	return
	// }

	// err = ts.ExecuteTemplate(w, "base", nil)

	// if err != nil {
	// 	log.Println(err.Error())
	// 	app.serverError(w, err)
	// }
	
	// w.Write([]byte("Hello from Snippetbox"))
	snippets, err := app.snippets.Latest()
	if err != nil {
	app.serverError(w, err)
	return
	}
	for _, snippet := range snippets {
	fmt.Fprintf(w, "%+v\n", snippet)
	}
}


func (app *application) snippetView(w http.ResponseWriter, r *http.Request)  {
	params := httprouter.ParamsFromContext(r.Context())


	id, err := strconv.Atoi(params.ByName("id"))
		if err != nil || id < 1 {
		app.notFound(w)
		return
		}
	
	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
	return
	}
	
	fmt.Fprintf(w, "%+v", snippet)
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request){
	fieldErrors := make(map[string]string)

	err := r.ParseForm()

	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")

	expires, err := strconv.Atoi(r.PostForm.Get("expires"))
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Check that the title value is not blank and is not more than 100
	// characters long. If it fails either of those checks, add a message to the
	// errors map using the field name as the key.
	if strings.TrimSpace(title) == "" {
	fieldErrors["title"] = "This field cannot be blank"
	} else if utf8.RuneCountInString(title) > 100 {
	fieldErrors["title"] = "This field cannot be more than 100 characters long"
	}
	// Check that the Content value isn't blank.
	if strings.TrimSpace(content) == "" {
	fieldErrors["content"] = "This field cannot be blank"
	}
	// Check the expires value matches one of the permitted values (1, 7 or
	// 365).
	if expires != 1 && expires != 7 && expires != 365 {
	fieldErrors["expires"] = "This field must equal 1, 7 or 365"
	}
	// If there are any errors, dump them in a plain text HTTP response and
	// return from the handler.
	if len(fieldErrors) > 0 {
	fmt.Fprint(w, fieldErrors)
	return
	}
	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	
	}

}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display the form for creating a new snippet..."))
	}
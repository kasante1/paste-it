 
package main

import (
"fmt"
"net/http"
"strconv"
"errors"
"github.com/kasante1/paste_it_backend/internal/models"
"github.com/julienschmidt/httprouter" 
"github.com/kasante1/paste_it_backend/internal/validator"
)

type snippetCreateForm struct {
	Title string
	Content string
	Expires int
	validator.Validator
}

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

	form := snippetCreateForm{
		Title: r.PostForm.Get("title"),
		Content: r.PostForm.Get("content"),
		Expires: expires,
	}

	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be 100 characters long")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
	form.CheckField(validator.PermittedInt(form.Expires, 1, 7, 365), "expires", "This field must equal 1 or 7 or 365")

	if !form.Valid(){
		return
	}
	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)


}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display the form for creating a new snippet..."))
	}
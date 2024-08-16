 
package main

import (
"fmt"
"net/http"
"strconv"
// "html/template" 
// "log"
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
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	
	}

	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	expires := 7

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
	app.serverError(w, err)
	return
	}


	http.Redirect(w, r, fmt.Sprintf("/snippet/view?id=%d", id), http.StatusSeeOther)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display the form for creating a new snippet..."))
	}
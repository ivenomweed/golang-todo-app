package main

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

type templateData struct {
	StringMap       map[string]string
	IntMap          map[string]int
	FloatMap        map[string]float32
	Data            map[string]interface{}
	CSRFToken       string
	Flash           string
	Warning         string
	Error           string
	IsAuthenticated int
	API             string
	CSSVersion      string
}

var functions = template.FuncMap{}

//go:embed templates
var templateFS embed.FS

func (app *application) addDefaultData(
	td *templateData,
	r *http.Request,
) *templateData {
	return td
}

func (app *application) renderTemplate(
	w http.ResponseWriter,
	r *http.Request,
	page string,
	td *templateData,
	partials ...string,
) error {
	var t *template.Template
	var err error
	templateToRender := fmt.Sprintf("templates/%s.page.html", page)
	template, templateInMap := app.templateCache[templateToRender]
	if app.config.env == "prod" && templateInMap {
		t = template
	} else {
		t, err = app.parseTemplate(partials, page, templateToRender)
		if err != nil {
			app.errorLog.Println(err)
			return err
		}
	}
	if td == nil {
		td = &templateData{}
	}
	td = app.addDefaultData(td, r)
	err = t.Execute(w, td)
	if err != nil {
		app.errorLog.Println(err)
		return err
	}
	return nil
}

func (app *application) parseTemplate(
	partials []string,
	page string,
	templateToRender string,
) (*template.Template, error) {
	var t *template.Template
	var err error
	// TODO: Build Partials
	if len(partials) > 0 {
		for i, x := range partials {
			partials[i] = fmt.Sprintf("templates/%s.partial.html", x)
		}
	}
	if len(partials) > 0 {
		t, err = template.New(fmt.Sprintf("%s.page.html", page)).
			Funcs(functions).
			ParseFS(
				templateFS,
				"templates/base.layout.html",
				strings.Join(partials, ","),
				templateToRender,
			)
	} else {
		t, err = template.New(fmt.Sprintf("%s.page.html", page)).
			Funcs(functions).
			ParseFS(
				templateFS,
				"templates/base.layout.html",
				templateToRender,
			)
	}
	if err != nil {
		app.errorLog.Println(err)
		return nil, err
	}

	app.templateCache[templateToRender] = t
	return t, nil
}

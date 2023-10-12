package main

//------------------------------------------------------------------------------

import (
	"embed"
	"errors"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/mattn/go-sqlite3"
)

//------------------------------------------------------------------------------

//go:embed templates/*.tmpl
var templateFS embed.FS

//go:embed static/*
var staticFS embed.FS

var Templates *template.Template
var randomColor color

// ------------------------------------------------------------------------------

type color struct {
	R int8
	G int8
	B int8
}

// Next returns the color slightly shifted up
func (c *color) Next() {
	c.R = c.R + int8(rand.Intn(100))
	c.G = c.G + int8(rand.Intn(100))
	c.B = c.B + int8(rand.Intn(100))

}

type poem struct {
	EnglishLines []string
	ChineseLines []string
	Color        color
	Code         []int
	CodeString   string
}

// ------------------------------------------------------------------------------
// New creates a poem from agruments or in absence a new random poem
func (p *poem) New(c ...int) error {
	var code []int

	if len(c) > 0 {
		code = c
	} else {
		code = FourRandomNumbers()
	}

	for _, num := range code {
		p.EnglishLines = append(p.EnglishLines, embeddedTable[num - 1].English)
		p.ChineseLines = append(p.ChineseLines, embeddedTable[num - 1].Chinese)
		p.Code = append(p.Code, embeddedTable[num - 1].ID)

	}

	p.CodeString = fmt.Sprintf("%v-%v-%v-%v", code[0], code[1], code[2], code[3])
	p.Color.Next()

	return nil
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	err := Templates.ExecuteTemplate(w, "base", nil)
	if err != nil {
		http.Error(w, "Error Generating Index", http.StatusInternalServerError)
	}
}

func about(w http.ResponseWriter, r *http.Request) {
	err := Templates.ExecuteTemplate(w, "about", nil)
	if err != nil {
		http.Error(w, "Error Generating Index", http.StatusInternalServerError)

	}
}

func listAll(w http.ResponseWriter, r *http.Request) {
	err := Templates.ExecuteTemplate(w, "list_all_lines", embeddedTable)
	if err != nil {
		http.Error(w, "Error Generating Index", http.StatusInternalServerError)

	}
}

func FourRandomNumbers() []int {
	var nums []int
	for i := 0; i < 4; i++ {
		nums = append(nums, rand.Intn(999))
	}
	return nums
}

func onePoem(w http.ResponseWriter, r *http.Request) {
	c := chi.URLParam(r, "code")
	parts := strings.Split(c, "-")
	var ints []int
	for _, num := range parts {
		i, err := strconv.Atoi(num)
		if err != nil {
			fmt.Println("Error converting string to int:", err)
			return
		}
		ints = append(ints, i)
	}
	var p poem
	if err := p.New(ints...); err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	err := Templates.ExecuteTemplate(w, "one_poem", p)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("Error executing template: %v", err)
	}
}

func tenPoems(w http.ResponseWriter, r *http.Request) {
	var pl []poem
	for i := 0; i < 10; i++ {
		var p poem
		if err := p.New(); err != nil {
			log.Fatal(err)
		}
		pl = append(pl, p)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	err := Templates.ExecuteTemplate(w, "ten_poems", pl)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("Error executing template: %v", err)
	}
}
//------------------------------------------------------------------------------
func main() {
	var err error
	funcMap := template.FuncMap{
		"sub": func(a, b int) int {
			return a - b
		},
	}

	Templates, err = template.New("").Funcs(funcMap).ParseFS(templateFS, "templates/*.tmpl")

	if err != nil {
		fmt.Printf("Error parsing templates: %v", err)
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	fileServer := http.FileServer(http.FS(staticFS))

	r.Handle("/static/*", http.StripPrefix("", fileServer))
	r.Get("/poem/{code}", onePoem)
	r.Get("/about", about)
	r.Get("/list", listAll)
	r.Get("/mc", tenPoems)
	r.Get("/", indexHandler)

	err = http.ListenAndServe(":8080", r)
	if err != nil {
		fmt.Println("Farts")
	}
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}

}

package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"errors"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

const one_poem_sql = "SELECT eng, rowid FROM new_table ORDER BY random() LIMIT (?)"
const decode_poem_sql = "SELECT eng, rowid FROM new_table WHERE rowid IN (?,?,?,?)"
var Templates *template.Template





type line struct {

	Id int
	Eng string
	Order int
	
}

type poem struct {
	Lines []line
	Color []int
	Code []int
	CodeString string
}

func (p *poem) Encode() (string, error) {
	if len(p.Code) == 0 {
		return "", errors.New("You are trying to encode an uninitialized poem. You can't do that.")
	}
	return fmt.Sprintf("%v-%v-%v-%v", p.Code[0], p.Code[1], p.Code[2], p.Code[3]), nil
}

func (p *poem) New(c ...int) error {
	var a []int
	var query string
	if len(c) > 0 {
		query = decode_poem_sql
		a = c
	} else {
		query = one_poem_sql
		a = []int{4}
	}

	// Convert []int to []interface{}
	// db.exec needs an []interface{} because that's golang's any type.
	interfaceSlice := make([]interface{}, len(a))
	for i, v := range a {
		interfaceSlice[i] = v
	}
	rows, err := db.Query(query, interfaceSlice...)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	
	// TODO - get better colors generated
	p.Color = []int{34, 34, 34}

	// Fill Lines with rows from db
	for rows.Next() {
		var l line
	
		if err := rows.Scan(&l.Eng, &l.Id); err != nil {
			log.Fatal(err)
		}
		p.Code = append(p.Code, l.Id)
		p.Lines = append(p.Lines, l)
	}

	// Encode index of lines as string
	var cs string
	if cs, err = p.Encode(); err != nil {
		log.Fatal(err)
	}
	p.CodeString = cs

	// Return an error is there is one
	if err := rows.Err(); err != nil {
		return err
	}
	return nil
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	err := Templates.ExecuteTemplate(w, "base", nil)
	if err != nil {
		http.Error(w, "Error Generating Index", http.StatusInternalServerError)
	}
}

type DBLine struct {
	id int
	english string
	chinese string
}

func GimmieData() {
	rows, err := db.Query("Select * from new_table")
	if err != nil {
		log.Fatal(err)
	}
	//	var full []DBLine
	var lines []string
	lines = append(lines, "var embeddedTable = []DBLine{")
	for rows.Next() {
		var l DBLine
		if err := rows.Scan(&l.id, &l.english, &l.chinese); err != nil {
			log.Fatal(err)
		}
		lines = append(lines, fmt.Sprintf("{ID: %d, English: %q, Chinese: %q},", l.id, l.english, l.chinese))
	}

	lines = append(lines, "}")

	output := strings.Join(lines, "\n")

	// Write the generated Go code to a file
	err = os.WriteFile("embedded_data.go", []byte(output), 0644)
	if err != nil {
		log.Fatal(err)
	}

}	
// 		fmt.Println("Hey There")
// 		full = append(full, l)
		
// 	}
// 	defer rows.Close()
	
// 	fmt.Printf("%v", full)
// }

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

	// err := Templates.ExecuteTemplate(w, "ten_poems", pl)
	// if err != nil {
	// 	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	// 	log.Printf("Error executing template: %v", err)
	// }
	
	
 	fmt.Fprintln(w, p)
	
}

func tenPoems(w http.ResponseWriter, r *http.Request) {
	var pl []poem
	for i := 0; i < 10; i++  {
		var p poem
		if err := p.New(); err != nil {
			log.Fatal(err)
		}
		pl = append(pl, p)
	}
	fmt.Println(pl)
	
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	err := Templates.ExecuteTemplate(w, "ten_poems", pl)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("Error executing template: %v", err)
	}
	//Templates.ExecuteTemplate(w, "ten_poems", pl)
}

func main() {
	var err error
	funcMap := template.FuncMap{
		"sub": func(a, b int) int {
			return a - b
		},
	}

	Templates, err = template.New("").Funcs(funcMap).ParseGlob("templates/*.tmpl")

	//	Templates, err = template.ParseGlob("templates.tmpl")
	if err != nil {
		fmt.Printf("Error parsing templates: %v", err)
	}
	
	
	db, err = sql.Open("sqlite3", "database.db")
	if err != nil {
	log.Fatalf("Farts: %v", err)
	}

	GimmieData()

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	fileServer := http.FileServer(http.Dir("static"))
	
	r.Handle("/static/*", http.StripPrefix("/static/", fileServer))
	r.Get("/poem/{code}", onePoem)
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

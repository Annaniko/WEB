package main
import (
    "fmt"
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
    "net/http"
    "html/template"
    "log"
    "github.com/gorilla/mux"
)
type Dogs struct{
    Id int
    Breed string
    Description string
}
var database *sql.DB
var templates = template.Must(template.ParseFiles("edit.html", "note.html", "index.html", ))
 
func DeleteHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id := vars["id"]
 
    _, err := database.Exec("delete from coursework.Dogs where id = ?", id)
    if err != nil{
        log.Println(err)
    }
     
    http.Redirect(w, r, "/note", 301)
}
 
func EditPage(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id := vars["id"]
 
    row := database.QueryRow("select * from coursework.Dogs where id = ?", id)
    prod := Dogs{}
    err := row.Scan(&prod.Id, &prod.Breed, &prod.Description)
    if err != nil{
        log.Println(err)
        http.Error(w, http.StatusText(404), http.StatusNotFound)
    }else{
        tmpl, _ := template.ParseFiles("edit.html")
        tmpl.Execute(w, prod)
    }
}
 
func EditHandler(w http.ResponseWriter, r *http.Request) {
    err := r.ParseForm()
    if err != nil {
        log.Println(err)
    }
    id := r.FormValue("id")
    breed := r.FormValue("breed")
	description := r.FormValue("description")
 
    _, err = database.Exec("update coursework.Dogs set  breed=?, description=? where id = ?", 
         breed, description, id)
 
    if err != nil {
        log.Println(err)
    }
    http.Redirect(w, r, "/note", 301)
}
 
func CreateHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == "POST" {
 
        err := r.ParseForm()
        if err != nil {
            log.Println(err)
        }
        breed := r.FormValue("breed")
		description := r.FormValue("description")
 
        _, err = database.Exec("insert into coursework.Dogs (breed, description) values (?, ?)", 
           breed, description)
 
        if err != nil {
            log.Println(err)
        }
        http.Redirect(w, r, "/note", 301)
    }else{
        http.ServeFile(w,r, "create.html")
    }
}
 
func NoteHandler(w http.ResponseWriter, r *http.Request) {
 
    rows, err := database.Query("select * from coursework.Dogs")
    if err != nil {
        log.Println(err)
    }
    defer rows.Close()
    products := []Dogs{}
     
    for rows.Next(){
        p := Dogs{}
        err := rows.Scan(&p.Id, &p.Breed, &p.Description)
        if err != nil{
            fmt.Println(err)
            continue
        }
        products = append(products, p)
    }
 
    tmpl, _ := template.ParseFiles("note.html")
    tmpl.Execute(w, products)
}
 
func IndexHandler(w http.ResponseWriter, r *http.Request) {
    products := []Dogs{}
    tmpl, _ := template.ParseFiles("index.html")
    tmpl.Execute(w, products)
}


func main() {
      
    db, err := sql.Open("mysql", "root:anyaniko@/coursework")
     
    if err != nil {
        log.Println(err)
    }
    database = db
    defer db.Close()
     
    router := mux.NewRouter()
    router.HandleFunc("/index", IndexHandler)
	router.HandleFunc("/note", NoteHandler)
    router.HandleFunc("/create", CreateHandler)
    router.HandleFunc("/edit/{id:[0-9]+}", EditPage).Methods("GET")
    router.HandleFunc("/edit/{id:[0-9]+}", EditHandler).Methods("POST")
    router.HandleFunc("/delete/{id:[0-9]+}", DeleteHandler)
     
    http.Handle("/",router)
 
    fmt.Println("Server is listening...")
    http.ListenAndServe(":8081", nil)
}
package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"personal-web/connection"
	"personal-web/middleware"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

type Blog struct {
	ID           int
	Name         string
	Startdate    time.Time
	Enddate      time.Time
	Format_date  string
	Test         string
	Description  string
	Technologies []string
	Image        string
	Author       string
}

// untuk mendapatkan form register
var Data = map[string]interface{}{
	"Title":   "Personal Web",
	"IsLogin": true,
}

// untuk medapatkan durasi blog
func (tgl Blog) Duration() string {
	// parsing string menjadi time.Time (sama seeperti fungsi new Date() pada javascript)
	start := tgl.Startdate
	// parsing string menjadi time.Time (sama seeperti fungsi new Date() pada javascript)
	end := tgl.Enddate

	// menghitung durasi menggunakan method .Sub milik object time.Time
	duration := end.Sub(start).Hours()

	// konversi ke hari
	var day int = 0
	for duration >= 24 {
		day += 1
		duration -= 24
	}

	return strconv.Itoa(day) + " Day"
}

// data untuk user
type User struct {
	ID       int
	Username string
	Email    string
	Password string
}

var Blogs = []Blog{}

func main() {
	route := mux.NewRouter()

	// connection to databse
	connection.DatabaseConnect()

	//route untuk menginisialisai folder public
	route.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))
	//route untuk menginisialisai folder uploads
	route.PathPrefix("/uploads/").Handler(http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads/"))))

	//route untuk menginisialisai folder views
	route.HandleFunc("/", home).Methods("GET")
	route.HandleFunc("/contact", contact).Methods("GET")
	route.HandleFunc("/formblog", formblog).Methods("GET")
	route.HandleFunc("/blog-detail/{id}", blogDetail).Methods("GET")
	route.HandleFunc("/addblog", middleware.UploadFile(addblog)).Methods("POST")
	route.HandleFunc("/delete-blog/{id}", deleteBlog).Methods("GET")
	route.HandleFunc("/update-blog/{id}", middleware.UploadFile(updateBlog)).Methods("POST")
	route.HandleFunc("/get-update-blog/{id}", getUpdateBlog).Methods("GET")
	route.HandleFunc("/register", formRegister).Methods("GET")
	route.HandleFunc("/register", Register).Methods("POST")
	route.HandleFunc("/login", formLogin).Methods("GET")
	route.HandleFunc("/login", Login).Methods("POST")
	route.HandleFunc("/logout", Logout).Methods("GET")

	//untuk menjalakan server
	fmt.Println("Server berjalan pada port 5000")
	http.ListenAndServe("localhost:5000", route)
}

// untuk mendapatkan form home
func home(w http.ResponseWriter, r *http.Request) {
	// untuk menseting type content yang akan di tampilkan
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	// parsing file html
	tmpt, err := template.ParseFiles("views/index.html")

	//jika terjadi error
	if err != nil {
		w.Write([]byte("Message : " + err.Error()))
		return
	}

	// untuk membuat session
	var store = sessions.NewCookieStore([]byte("SESSION_ID"))
	session, _ := store.Get(r, "SESSION_ID")

	// untuk mengecek apakah user sudah login atau belum
	if session.Values["IsLogin"] != true {
		Data["IsLogin"] = false
	} else {
		Data["IsLogin"] = session.Values["IsLogin"].(bool)
		Data["Username"] = session.Values["Username"].(string)
		Data["ID"] = session.Values["ID"].(int)
	}

	// untuk mendapatkan data dari database menggunakan relasi
	rows, _ := connection.Conn.Query(context.Background(), "SELECT tbl_blog.id, name, start_date, end_date, description, technologies, image, tbl_login.username as author FROM tbl_blog LEFT JOIN tbl_login ON tbl_blog.author_id = tbl_login.id  ORDER BY id ASC")

	// untuk menghapus data yang ada di slice Blogs
	var result []Blog
	// untuk menambahkan data yang baru ke slice Blogs
	for rows.Next() {
		// untuk menampung data dari database
		var each = Blog{}
		// untuk menampung data dari database
		var err = rows.Scan(&each.ID, &each.Name, &each.Startdate, &each.Enddate, &each.Description, &each.Technologies, &each.Image, &each.Author)
		// jika terjadi error
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		// untuk menformat tanggal
		each.Format_date = each.Startdate.Format("2 January 2006")
		each.Test = each.Enddate.Format("2 January 2006")

		// untuk menambahkan data ke slice Blogs
		result = append(result, each)
	}

	// untuk mengirimkan data ke html
	resData := map[string]interface{}{
		"Data":  Data,
		"Blogs": result,
	}

	// memberikan status code 200
	w.WriteHeader(http.StatusOK)
	// untuk menampilkan data di html
	tmpt.Execute(w, resData)
}

// untuk mendapatkan form contact
func contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	tmpt, err := template.ParseFiles("views/contact.html")

	if err != nil {
		w.Write([]byte("Message : " + err.Error()))
		return
	}

	var store = sessions.NewCookieStore([]byte("SESSION_ID"))
	session, _ := store.Get(r, "SESSION_ID")

	if session.Values["IsLogin"] != true {
		Data["IsLogin"] = false
	} else {
		Data["IsLogin"] = session.Values["IsLogin"].(bool)
		Data["UserName"] = session.Values["Username"].(string)
	}
	resData := map[string]interface{}{
		"Data": Data,
	}

	tmpt.Execute(w, resData)
}

// untuk mendapatkan form addblog
func formblog(w http.ResponseWriter, r *http.Request) {
	//set header
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	//parsing file html
	tmpt, err := template.ParseFiles("views/addblog.html")
	//jika error
	if err != nil {
		w.Write([]byte("Message : " + err.Error()))
		return
	}
	//jika tidak error
	var store = sessions.NewCookieStore([]byte("SESSION_ID"))
	session, _ := store.Get(r, "SESSION_ID")

	//jika session tidak ada dan jika session ada
	if session.Values["IsLogin"] != true {
		Data["IsLogin"] = false
	} else {
		Data["IsLogin"] = session.Values["IsLogin"].(bool)
		Data["UserName"] = session.Values["Username"].(string)
	}

	//mengirim data ke html
	resData := map[string]interface{}{
		"Data": Data,
	}
	//menampilkan data di html
	tmpt.Execute(w, resData)
}

// fungsi untuk menambahkan data ke database
func addblog(w http.ResponseWriter, r *http.Request) {
	// ambil data dari form
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	//mengambil data dari form
	name := r.PostForm.Get("name")
	startdate, _ := time.Parse("2006-01-02", r.PostForm.Get("std"))
	enddate, _ := time.Parse("2006-01-02", r.PostForm.Get("etd"))
	description := r.PostForm.Get("description")
	technologies := r.PostForm["technologies"]

	//mengambil data file image
	dataContex := r.Context().Value("dataFile")
	image := dataContex.(string)

	//mengambil data dari session
	var store = sessions.NewCookieStore([]byte("SESSION_ID"))
	session, _ := store.Get(r, "SESSION_ID")
	author := session.Values["ID"].(int)

	//memasukan data ke database
	_, err = connection.Conn.Exec(context.Background(), "INSERT INTO tbl_blog(name, start_date, end_date, description, technologies, image, author_id) VALUES ($1,$2,$3,$4,$5,$6,$7)",
		name, startdate, enddate, description, technologies, image, author)

	//jika error
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	//jika tidak error
	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

// get update blog[index]
func getUpdateBlog(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	tmpt, err := template.ParseFiles("views/update-blog.html")

	if err != nil {
		w.Write([]byte("Message : " + err.Error()))
		return
	}

	// id, _ := strconv.Atoi(mux.Vars(r)["index"])

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	var result []Blog
	var each = Blog{}
	err = connection.Conn.QueryRow(context.Background(), "SELECT id, name, start_date, end_date, description, technologies, image FROM tbl_blog WHERE id=$1", id).Scan(
		&each.ID, &each.Name, &each.Startdate, &each.Enddate, &each.Description, &each.Technologies, &each.Image)
	if err != nil {
		// mengirim status code 500
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}
	each.Format_date = each.Startdate.Format("2006-01-02")
	each.Test = each.Enddate.Format("2006-01-02")

	result = append(result, each)

	resData := map[string]interface{}{
		"Blogs": result,
	}

	tmpt.Execute(w, resData)
}

// fungsi update blog berdasarkan id
func updateBlog(w http.ResponseWriter, r *http.Request) {
	// ambil data dari form
	err := r.ParseMultipartForm(1024)
	//jika error
	if err != nil {
		log.Fatal(err)
	}

	//mengambil data dari form
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	name := r.PostForm.Get("name")
	startdate, _ := time.Parse("2006-01-02", r.PostForm.Get("std"))
	enddate, _ := time.Parse("2006-01-02", r.PostForm.Get("etd"))
	description := r.PostForm.Get("description")
	technologies := r.PostForm["technologies"]

	//mengambil data file image
	dataContex := r.Context().Value("dataFile")
	image := dataContex.(string)

	//mengambil data dari session
	var store = sessions.NewCookieStore([]byte("SESSION_ID"))
	session, _ := store.Get(r, "SESSION_ID")
	author := session.Values["ID"].(int)

	//mengupdate data ke database
	_, err = connection.Conn.Exec(context.Background(), "UPDATE tbl_blog SET name=$1, start_date=$2, end_date=$3, description=$4, technologies=$5, image=$6, author_id=$7 WHERE id=$8",
		name, startdate, enddate, description, technologies, image, author, id)

	//jika error
	if err != nil {
		//mengirim error ke html
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	//jika tidak error
	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

// untuk mendapatkan blog detail
func blogDetail(w http.ResponseWriter, r *http.Request) {
	//untuk menset content type
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	//parsing file html
	tmpt, err := template.ParseFiles("views/blog-detail.html")
	//jika error
	if err != nil {
		w.Write([]byte("Message : " + err.Error()))
		return
	}

	//untuk session
	var store = sessions.NewCookieStore([]byte("SESSION_ID"))
	session, _ := store.Get(r, "SESSION_ID")
	//untuk mengecek session
	if session.Values["IsLogin"] != true {
		Data["IsLogin"] = false
	} else {
		Data["IsLogin"] = session.Values["IsLogin"].(bool)
		Data["Username"] = session.Values["Username"].(string)
		Data["ID"] = session.Values["ID"].(int)
	}

	//convert string ke int
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	// untuk menghapus data yang ada di slice Blogs
	var result []Blog
	//untuk menampung data dari database
	var each = Blog{}
	//untuk mengambil data dari database
	err = connection.Conn.QueryRow(context.Background(), "SELECT id, name, start_date, end_date, description, technologies, image FROM tbl_blog WHERE id=$1", id).Scan(
		&each.ID, &each.Name, &each.Startdate, &each.Enddate, &each.Description, &each.Technologies, &each.Image)
	if err != nil {
		fmt.Println("Message : " + err.Error())
		return
	}

	//untuk mengubah format tanggal
	each.Format_date = each.Startdate.Format("2 January 2006")
	each.Test = each.Enddate.Format("2 January 2006")

	//untuk menambahkan data ke slice Blogs
	result = append(result, each)

	//untuk mengirimkan data ke html
	resData := map[string]interface{}{
		"Data":  Data,
		"Blogs": result,
	}

	//untuk mengeksekusi file html
	tmpt.Execute(w, resData)
}

// untuk menghapus blog
func deleteBlog(w http.ResponseWriter, r *http.Request) {
	//untuk koncersi string ke int
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	//untuk menghapus data dari database
	_, err := connection.Conn.Exec(context.Background(), "DELETE FROM tbl_blog WHERE id=$1", id)
	//jika error
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}
	//jika tidak error
	http.Redirect(w, r, "/", http.StatusFound)
}

// untuk mendapatkan form register
func formRegister(w http.ResponseWriter, r *http.Request) {
	//untuk menset content type
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	//untuk parsing file html
	var tmpl, err = template.ParseFiles("views/formregis.html")
	//jika error
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}
	//jika tidak error
	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, Data)
}

// fungsi untuk register
func Register(w http.ResponseWriter, r *http.Request) {
	//parse form
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	//mengambil data dari form
	name := r.PostForm.Get("username")
	email := r.PostForm.Get("email")
	password := r.PostForm.Get("password")
	//encrypt password
	passwordcrypt, _ := bcrypt.GenerateFromPassword([]byte(password), 10)

	//memasukan data ke database
	_, err = connection.Conn.Exec(context.Background(), "INSERT INTO tbl_login (username, email, password) VALUES ($1, $2, $3)",
		name, email, passwordcrypt)
	//jika error
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	//jika tidak error
	http.Redirect(w, r, "/", http.StatusFound)
}

// untuk mendapatkan form login
func formLogin(w http.ResponseWriter, r *http.Request) {
	//untuk menset content type
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	//untuk parsing file html
	var tmpl, err = template.ParseFiles("views/formlogin.html")
	//	jika error
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}
	//jika tidak error
	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, Data)
}

// fungsi untuk login
func Login(w http.ResponseWriter, r *http.Request) {
	//parse form
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	//mengambil data dari form
	email := r.PostForm.Get("email")
	password := r.PostForm.Get("password")

	//untuk mengecek email dan password
	user := User{}
	//untuk mengambil data dari database
	_ = connection.Conn.QueryRow(context.Background(), "SELECT * FROM tbl_login WHERE email=$1", email).Scan(
		&user.ID, &user.Username, &user.Email, &user.Password)

	//untuk mengecek password
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("message : " + err.Error()))
		return
	}
	//untuk mengecek password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	//untuk membuat session
	var store = sessions.NewCookieStore([]byte("SESSION_ID"))
	session, _ := store.Get(r, "SESSION_ID")

	//untuk mengecek session
	session.Values["IsLogin"] = true
	session.Values["Username"] = user.Username
	session.Values["ID"] = user.ID
	session.Options.MaxAge = 10800

	//untuk menyimpan session
	session.AddFlash("Login success", "message")
	session.Save(r, w)

	//jika tidak error
	http.Redirect(w, r, "/", http.StatusFound)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	//untuk membuat session
	var store = sessions.NewCookieStore([]byte("SESSION_ID"))
	session, _ := store.Get(r, "SESSION_ID")

	//untuk menghapus session
	session.Options.MaxAge = -1
	session.Save(r, w)

	//jika tidak error
	http.Redirect(w, r, "/", http.StatusFound)
}

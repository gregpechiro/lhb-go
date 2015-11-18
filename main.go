package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/cagnosolutions/dbdb"
	"github.com/cagnosolutions/web"
)

const (
	USERNAME = "admin"
	PASSWORD = "admin"
)

var mux = web.NewMux()
var tmpl = web.NewTmplCache()
var db = dbdb.NewDataStore()

func main() {
	//web.Funcs["snake"] = Snake
	db.AddStore("image")
	db.AddStore("listing")
	mux.AddRoutes(home, gallery, about, contact, services, listings, floorPlans, login, logout)
	mux.AddSecureRoutes(ADMIN, webmaster, allListings, oneImage, uploadImage, saveImage, deleteImage, oneListing, addListing, saveListing, deleteListing)
	http.ListenAndServe(":8080", mux)

}

var home = web.Route{"GET", "/", func(w http.ResponseWriter, r *http.Request) {
	tmpl.Render(w, r, "home.tmpl", web.Model{})
}}

var gallery = web.Route{"GET", "/gallery", func(w http.ResponseWriter, r *http.Request) {
	images := db.GetAll("image")
	tmpl.Render(w, r, "gallery.tmpl", web.Model{
		"images": images,
		"cats":   getCategories(images),
	})
}}

var about = web.Route{"GET", "/about", func(w http.ResponseWriter, r *http.Request) {
	tmpl.Render(w, r, "about.tmpl", web.Model{})
}}

var contact = web.Route{"GET", "/contact", func(w http.ResponseWriter, r *http.Request) {
	tmpl.Render(w, r, "contact.tmpl", web.Model{})
}}

var services = web.Route{"GET", "/services", func(w http.ResponseWriter, r *http.Request) {
	tmpl.Render(w, r, "services.tmpl", web.Model{})
}}

var listings = web.Route{"GET", "/listings", func(w http.ResponseWriter, r *http.Request) {
	tmpl.Render(w, r, "listings.tmpl", web.Model{
		"listings": db.GetAll("listing"),
	})
}}

var floorPlans = web.Route{"GET", "/floor-plans", func(w http.ResponseWriter, r *http.Request) {
	tmpl.Render(w, r, "floor-plans.tmpl", web.Model{})
}}

var login = web.Route{"POST", "/login", func(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("username") == USERNAME && r.FormValue("password") == PASSWORD {
		web.Login(w, r, "admin")
		web.SetSuccessRedirect(w, r, "/webmaster", "You are now logged in")
		return
	}
	http.Redirect(w, r, "/", 303)
	return
}}

var logout = web.Route{"GET", "/logout", func(w http.ResponseWriter, r *http.Request) {
	web.Logout(w, r)
	web.SetSuccessRedirect(w, r, "/", "You are now logged out")
}}

var webmaster = web.Route{"GET", "/webmaster", func(w http.ResponseWriter, r *http.Request) {
	images := db.GetAll("image")
	tmpl.Render(w, r, "webmaster.tmpl", web.Model{
		"images": images,
		"cats":   getCategories(images),
		"page":   "webmaster",
	})
}}

var uploadImage = web.Route{"POST", "/upload-image", func(w http.ResponseWriter, r *http.Request) {
	path := "static/img/upload/"
	if err := os.MkdirAll(path, 0755); err != nil {
		fmt.Printf("uploadImage >> MkdirAll: %v\n", err)
		web.SetErrorRedirect(w, r, "/webmaster", "Error uploading file")
		return
	}
	r.ParseMultipartForm(32 << 20) // 32 MB
	file, handler, err := r.FormFile("picture")
	if err != nil || len(handler.Header["Content-Type"]) < 1 {
		fmt.Printf("uploadImage >> Header len < 1: %v\n", err)
		web.SetErrorRedirect(w, r, "/webmaster", "Error uploading file")
		return
	}
	defer file.Close()
	if handler.Header["Content-Type"][0] != "image/png" && handler.Header["Content-Type"][0] != "image/jpeg" {
		fmt.Printf("uploadImage >> Header != png || jpeg: %v\n", err)
		web.SetErrorRedirect(w, r, "/webmaster", "Error uploading file")
		return
	}
	f, err := os.OpenFile(path+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Printf("uploadImage >> OpenFile: %v\n", err)
		web.SetErrorRedirect(w, r, "/webmaster", "Error uploading file")
		return
	}
	defer f.Close()
	io.Copy(f, file)
	doc := map[string]interface{}{
		"category":    r.FormValue("category"),
		"description": r.FormValue("description"),
		"source":      handler.Filename,
	}
	db.Add("image", doc)
	web.SetSuccessRedirect(w, r, "/webmaster", "Successfully uploaded image")
	return

}}

var saveImage = web.Route{"POST", "/save-image/:id", func(w http.ResponseWriter, r *http.Request) {
	id := getId(r.FormValue(":id"))
	img := db.Get("image", id).Data
	img["category"] = r.FormValue("category")
	img["description"] = r.FormValue("description")
	db.Set("image", id, img)
	web.SetSuccessRedirect(w, r, "/webmaster", "Successfully saved image")
}}

var oneImage = web.Route{"GET", "/webmaster/:id", func(w http.ResponseWriter, r *http.Request) {
	tmpl.Render(w, r, "webmaster.tmpl", web.Model{
		"images": db.GetAll("image"),
		"image":  db.Get("image", getId(r.FormValue(":id"))),
		"page":   "webmaster",
	})
}}

var deleteImage = web.Route{"POST", "/webmaster/:id", func(w http.ResponseWriter, r *http.Request) {
	img := db.Get("image", getId(r.FormValue(":id")))
	if img == nil {
		log.Printf("webmasterDeleteImage >> db.Get: image could not be found in database\n")
		web.SetErrorRedirect(w, r, "/webmaster", "Error deleting image.")
		return
	}
	if err := os.Remove("static/img/upload/" + img.Data["source"].(string)); err != nil {
		log.Printf("webmasterDeleteImage >> os.Remove: %v\n", err)
		web.SetErrorRedirect(w, r, "/webmaster", "Error deleting image.")
		return
	}
	db.Del("image", img.Id)
	web.SetSuccessRedirect(w, r, "/webmaster", "Successfully deleted image")
	return
}}

var allListings = web.Route{"GET", "/all-listings", func(w http.ResponseWriter, r *http.Request) {
	tmpl.Render(w, r, "all-listings.tmpl", web.Model{
		"listings": db.GetAll("listing"),
		"page":     "listings",
	})
	return
}}

var oneListing = web.Route{"GET", "/all-listings/:id", func(w http.ResponseWriter, r *http.Request) {
	tmpl.Render(w, r, "all-listings.tmpl", web.Model{
		"listings": db.GetAll("listing"),
		"listing":  db.Get("listing", getId(r.FormValue(":id"))),
		"page":     "listings",
	})
	return
}}

var addListing = web.Route{"POST", "/save-listing", func(w http.ResponseWriter, r *http.Request) {
	listing := map[string]interface{}{
		"street": r.FormValue("street"),
		"city":   r.FormValue("city"),
		"state":  r.FormValue("state"),
		"zip":    r.FormValue("zip"),
		"mls":    r.FormValue("mls"),
		"agent":  r.FormValue("agent"),
		"phone":  r.FormValue("phone"),
	}
	db.Add("listing", listing)
	web.SetSuccessRedirect(w, r, "/all-listings", "Successfuly added listing")
}}

var saveListing = web.Route{"POST", "/save-listing/:id", func(w http.ResponseWriter, r *http.Request) {
	listing := map[string]interface{}{
		"street": r.FormValue("street"),
		"city":   r.FormValue("city"),
		"state":  r.FormValue("state"),
		"zip":    r.FormValue("zip"),
		"mls":    r.FormValue("mls"),
		"agent":  r.FormValue("agent"),
		"phone":  r.FormValue("phone"),
	}
	db.Set("listing", getId(r.FormValue(":id")), listing)
	web.SetSuccessRedirect(w, r, "/all-listings", "Successfuly saved listing")
}}

var deleteListing = web.Route{"POST", "/all-listing/:id", func(w http.ResponseWriter, r *http.Request) {
	db.Del("listing", getId(r.FormValue(":id")))
	web.SetSuccessRedirect(w, r, "/all-listings", "Successfully deleted listing")
	return
}}

func getId(sid string) uint64 {
	id, _ := strconv.ParseUint(sid, 10, 64)
	return id
}

func Snake(s string) string {
	return strings.ToLower(strings.Replace(s, " ", "-", -1))
}

func getCategories(images []*dbdb.Doc) []string {
	m := make(map[string]string)
	ss := make([]string, 0)
	for _, v := range images {
		m[v.Data["category"].(string)] = "v"
	}
	for k := range m {
		ss = append(ss, k)
	}
	sort.Strings(ss)
	return ss
}

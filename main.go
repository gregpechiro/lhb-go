package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/cagnosolutions/dbdb"
	"github.com/cagnosolutions/web"
)

var mux = web.NewMux().CSRF()
var tmpl = web.NewTmplCache()
var db = dbdb.NewDataStore()

func main() {

	db.AddStore("images")
	mux.AddRoutes(home, gallery, about, contact, services, listings, floorPlans, login, logout, msg)
	mux.AddSecureRoutes(ADMIN, webmaster, allListings, uploadImage, webmasterImage, webmasterDeleteImage)
	http.ListenAndServe(":8080", mux)

}

var home = web.Route{"GET", "/", func(w http.ResponseWriter, r *http.Request) {
	tmpl.Render(w, r, "home.tmpl", web.Model{})
}}

var msg = web.Route{"GET", "/msg", func(w http.ResponseWriter, r *http.Request) {
	web.SetSuccessRedirect(w, r, "/", "Message")
	return
}}

var gallery = web.Route{"GET", "/gallery", func(w http.ResponseWriter, r *http.Request) {
	tmpl.Render(w, r, "gallery.tmpl", web.Model{})
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
	tmpl.Render(w, r, "listings.tmpl", web.Model{})
}}

var floorPlans = web.Route{"GET", "/floor-plans", func(w http.ResponseWriter, r *http.Request) {
	tmpl.Render(w, r, "floor-plans.tmpl", web.Model{})
}}

var login = web.Route{"GET", "/login", func(w http.ResponseWriter, r *http.Request) {
	web.Login(w, r, "admin")
	web.SetSuccessRedirect(w, r, "/webmaster", "You are now logged in")
}}

var logout = web.Route{"GET", "/logout", func(w http.ResponseWriter, r *http.Request) {
	web.Logout(w, r)
	web.SetSuccessRedirect(w, r, "/", "You are now logged out")
}}

var webmaster = web.Route{"GET", "/webmaster", func(w http.ResponseWriter, r *http.Request) {

	images := db.GetAll("images")
	tmpl.Render(w, r, "webmaster.tmpl", web.Model{
		"images": images,
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
		"src":         handler.Filename,
	}
	db.Add("images", doc)
	web.SetSuccessRedirect(w, r, "/webmaster", "Successfully uploaded image")
	return

}}

var webmasterImage = web.Route{"GET", "/webmaster/:id", func(w http.ResponseWriter, r *http.Request) {
	tmpl.Render(w, r, "webmaster.tmpl", web.Model{
		"images": db.GetAll("images"),
		"image":  db.Get("images", getId(r.FormValue(":id"))),
	})
}}

var webmasterDeleteImage = web.Route{"POST", "/webmaster/:id", func(w http.ResponseWriter, r *http.Request) {
	img := db.Get("images", getId(r.FormValue(":id")))
	if img == nil {
		log.Printf("webmasterDeleteImage >> db.Get: image could not be found in database\n")
		web.SetErrorRedirect(w, r, "/webmaster", "Error deleting image.")
		return
	}
	if err := os.Remove("static/img/uplaod/" + img.Data["src"].(string)); err != nil {
		log.Printf("webmasterDeleteImage >> os.Remove: %v\n", err)
		web.SetErrorRedirect(w, r, "/webmaster", "Error deleting image.")
		return
	}
	db.Del("images", img.Id)
	web.SetSuccessRedirect(w, r, "/webmaster", "Successfully deleted image")
}}

var allListings = web.Route{"GET", "/all-listings", func(w http.ResponseWriter, r *http.Request) {
	tmpl.Render(w, r, "all-listings.tmpl", web.Model{})
}}

func getId(sid string) uint64 {
	id, _ := strconv.ParseUint(sid, 10, 64)
	return id
}

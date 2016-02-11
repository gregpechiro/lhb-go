package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
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
var tmpl *web.TmplCache
var db = dbdb.NewDataStore()

func init() {
	web.Funcs["split"] = func(s, sep string, i int) string {
		return strings.Split(s, sep)[i]
	}
	web.Funcs["title"] = strings.Title
	tmpl = web.NewTmplCache()
	db.AddStore("image")
	db.AddStore("listing")
}

func main() {

	mux.AddRoutes(home, gallery, about, contact, services, listings, floorPlans, login, logout)

	mux.AddSecureRoutes(WEBMASTER, allFloorplans, uploadFloorplan, renameFloorplan, deleteFloorplan)

	mux.AddSecureRoutes(WEBMASTER, allListings, oneListing, addListing, saveListing, deleteListing)

	mux.AddSecureRoutes(WEBMASTER, webmaster, oneImage, uploadImage, saveImage, deleteImage)

	log.Fatal(http.ListenAndServe(":8080", mux))

}

var home = web.Route{"GET", "/", func(w http.ResponseWriter, r *http.Request) {
	tmpl.Render(w, r, "home.tmpl", web.Model{})
	return
}}

var gallery = web.Route{"GET", "/gallery", func(w http.ResponseWriter, r *http.Request) {
	images := db.GetAll("image")
	tmpl.Render(w, r, "gallery.tmpl", web.Model{
		"images": images,
		"cats":   getCategories(images),
	})
	return
}}

var about = web.Route{"GET", "/about", func(w http.ResponseWriter, r *http.Request) {
	tmpl.Render(w, r, "about.tmpl", web.Model{})
	return
}}

var contact = web.Route{"GET", "/contact", func(w http.ResponseWriter, r *http.Request) {
	tmpl.Render(w, r, "contact.tmpl", web.Model{})
	return
}}

var services = web.Route{"GET", "/services", func(w http.ResponseWriter, r *http.Request) {
	tmpl.Render(w, r, "services.tmpl", web.Model{})
	return
}}

var listings = web.Route{"GET", "/listings", func(w http.ResponseWriter, r *http.Request) {
	tmpl.Render(w, r, "listings.tmpl", web.Model{
		"listings": db.GetAll("listing"),
	})
	return
}}

var floorPlans = web.Route{"GET", "/floor-plans", func(w http.ResponseWriter, r *http.Request) {
	tmpl.Render(w, r, "floor-plans.tmpl", web.Model{
		"floorplans": GetFloorPlans(),
	})
	return
}}

var login = web.Route{"POST", "/login", func(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("username") == USERNAME && r.FormValue("password") == PASSWORD {
		web.Login(w, r, "webmaster")
		web.SetSuccessRedirect(w, r, "/webmaster", "You are now logged in")
		return
	}
	http.Redirect(w, r, "/", 303)
	return
}}

var logout = web.Route{"GET", "/logout", func(w http.ResponseWriter, r *http.Request) {
	web.Logout(w)
	web.SetSuccessRedirect(w, r, "/", "You are now logged out")
	return
}}

func ParseId(v interface{}) float64 {
	var id float64
	var err error
	switch v.(type) {
	case string:
		id, err = strconv.ParseFloat(v.(string), 64)
		if err != nil {
			log.Printf("ParseId() >> strconv.ParseFloat(): ", err)
		}
	case uint64:
		id = float64(v.(uint64))
	case float64:
		id = v.(float64)
	}
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

func GetFloorPlans() []string {
	var fp []string
	filepath.Walk("static/floorplans", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.Name() == "floorplans" {
			return nil
		}
		if info.IsDir() {
			return filepath.SkipDir
		}
		fp = append(fp, info.Name())
		return nil
	})
	sort.Strings(fp)
	return fp
}

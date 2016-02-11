package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/cagnosolutions/web"
)

var webmaster = web.Route{"GET", "/webmaster", func(w http.ResponseWriter, r *http.Request) {
	images := db.GetAll("image")
	tmpl.Render(w, r, "webmaster.tmpl", web.Model{
		"images": images,
		"cats":   getCategories(images),
		"page":   "webmaster",
	})
	return
}}

var uploadImage = web.Route{"POST", "/webmaster/upload-image", func(w http.ResponseWriter, r *http.Request) {
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

var saveImage = web.Route{"POST", "/webmaster/save-image/:id", func(w http.ResponseWriter, r *http.Request) {
	id := ParseId(r.FormValue(":id"))
	img := db.Get("image", id).Data
	img["category"] = r.FormValue("category")
	img["description"] = r.FormValue("description")
	db.Set("image", id, img)
	web.SetSuccessRedirect(w, r, "/webmaster", "Successfully saved image")
	return
}}

var oneImage = web.Route{"GET", "/webmaster/:id", func(w http.ResponseWriter, r *http.Request) {
	tmpl.Render(w, r, "webmaster.tmpl", web.Model{
		"images": db.GetAll("image"),
		"image":  db.Get("image", ParseId(r.FormValue(":id"))),
		"page":   "webmaster",
	})
	return
}}

var deleteImage = web.Route{"POST", "/webmaster/:id", func(w http.ResponseWriter, r *http.Request) {
	img := db.Get("image", ParseId(r.FormValue(":id")))
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

var allListings = web.Route{"GET", "/webmaster/all-listings", func(w http.ResponseWriter, r *http.Request) {
	tmpl.Render(w, r, "webmaster-all-listings.tmpl", web.Model{
		"listings": db.GetAll("listing"),
		"page":     "listings",
	})
	return
}}

var oneListing = web.Route{"GET", "/webmaster/all-listings/:id", func(w http.ResponseWriter, r *http.Request) {
	tmpl.Render(w, r, "webmaster-all-listings.tmpl", web.Model{
		"listings": db.GetAll("listing"),
		"listing":  db.Get("listing", ParseId(r.FormValue(":id"))),
		"page":     "listings",
	})
	return
}}

var addListing = web.Route{"POST", "/webmaster/save-listing", func(w http.ResponseWriter, r *http.Request) {
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
	web.SetSuccessRedirect(w, r, "/webmaster/all-listings", "Successfuly added listing")
	return
}}

var saveListing = web.Route{"POST", "/webmaster/save-listing/:id", func(w http.ResponseWriter, r *http.Request) {
	listing := map[string]interface{}{
		"street": r.FormValue("street"),
		"city":   r.FormValue("city"),
		"state":  r.FormValue("state"),
		"zip":    r.FormValue("zip"),
		"mls":    r.FormValue("mls"),
		"agent":  r.FormValue("agent"),
		"phone":  r.FormValue("phone"),
	}
	db.Set("listing", ParseId(r.FormValue(":id")), listing)
	web.SetSuccessRedirect(w, r, "/webaster/all-listings", "Successfuly saved listing")
	return
}}

var deleteListing = web.Route{"POST", "/webmaster/all-listing/:id", func(w http.ResponseWriter, r *http.Request) {
	db.Del("listing", ParseId(r.FormValue(":id")))
	web.SetSuccessRedirect(w, r, "/all-listings", "Successfully deleted listing")
	return
}}

var allFloorplans = web.Route{"GET", "/webmaster/floorplans", func(w http.ResponseWriter, r *http.Request) {
	tmpl.Render(w, r, "webmaster-floorplans.tmpl", web.Model{
		"floorplans": GetFloorPlans(),
		"page":       "floorplans",
	})
	return
}}

var uploadFloorplan = web.Route{"POST", "/webmaster/upload-floorplan", func(w http.ResponseWriter, r *http.Request) {
	path := "static/floorplans/"
	if err := os.MkdirAll(path, 0755); err != nil {
		fmt.Printf("uploadFloorplan >> MkdirAll: %v\n", err)
		web.SetErrorRedirect(w, r, "/webmaster/floorplans", "Error uploading file")
		return
	}
	r.ParseMultipartForm(32 << 20) // 32 MB
	file, handler, err := r.FormFile("floorplan")
	if err != nil || len(handler.Header["Content-Type"]) < 1 {
		fmt.Printf("uploadFloorplan >> Header len < 1: %v\n", err)
		web.SetErrorRedirect(w, r, "/webmaster/floorplans", "Error uploading file")
		return
	}
	defer file.Close()
	contentType := handler.Header["Content-Type"][0]
	fileName := r.FormValue("name")
	if contentType == "image/png" {
		fileName += ".png"
	} else if contentType == "image/jpeg" {
		fileName += ".jpg"
	} else if contentType == "application/pdf" {
		fileName += ".pdf"
	} else {
		fmt.Printf("uploadFloorplan >> Header[\"content-type\"] != png || jpeg || pdf: content-type is: %s\n", contentType)
		web.SetErrorRedirect(w, r, "/webmaster/floorplans", "Error uploading file")
		return
	}
	f, err := os.OpenFile(path+fileName, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Printf("uploadFloorplan >> OpenFile: %v\n", err)
		web.SetErrorRedirect(w, r, "/webmaster/floorplans", "Error uploading file")
		return
	}
	defer f.Close()
	io.Copy(f, file)
	web.SetSuccessRedirect(w, r, "/webmaster/floorplans", "Successfully uploaded image")
	return
}}

var deleteFloorplan = web.Route{"POST", "/webmaster/floorplan/:name", func(w http.ResponseWriter, r *http.Request) {
	err := os.Remove("static/floorplans/" + r.FormValue(":name"))
	if err != nil {
		web.SetErrorRedirect(w, r, "/webmaster/floorplans", "Error deleting floorplan")
		return
	}
	web.SetSuccessRedirect(w, r, "/webmaster/floorplans", "Successfully deleted floorplan")
	return
}}

var renameFloorplan = web.Route{"POST", "/webmaster/floorplan/rename", func(w http.ResponseWriter, r *http.Request) {
	ext := "." + strings.Split(r.FormValue("oldName"), ".")[1]
	err := os.Rename("static/floorplans/"+r.FormValue("oldName"), "static/floorplans/"+r.FormValue("name")+ext)
	if err != nil {
		web.SetErrorRedirect(w, r, "/webmaster/floorplans", "Error renaming floorplan")
		return
	}
	web.SetSuccessRedirect(w, r, "/webmaster/floorplans", "Successfully renamed floorplan")
	return

}}

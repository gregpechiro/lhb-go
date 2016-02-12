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

var webmaster = web.Route{"GET", "/webmaster/gallery", func(w http.ResponseWriter, r *http.Request) {
	images := db.GetAll("image")
	tmpl.Render(w, r, "webmaster-gallery.tmpl", web.Model{
		"images": images,
		"cats":   getCategories(images),
		"page":   "gallery",
	})
	return
}}

var uploadImage = web.Route{"POST", "/webmaster/gallery/upload", func(w http.ResponseWriter, r *http.Request) {
	path := "upload/gallery/"
	if err := os.MkdirAll(path, 0755); err != nil {
		fmt.Printf("uploadImage >> MkdirAll: %v\n", err)
		web.SetErrorRedirect(w, r, "/webmaster/gallery", "Error uploading file")
		return
	}
	r.ParseMultipartForm(32 << 20) // 32 MB
	file, handler, err := r.FormFile("picture")
	if err != nil || len(handler.Header["Content-Type"]) < 1 {
		fmt.Printf("uploadImage >> Header len < 1: %v\n", err)
		web.SetErrorRedirect(w, r, "/webmaster/gallery", "Error uploading file")
		return
	}
	defer file.Close()
	if handler.Header["Content-Type"][0] != "image/png" && handler.Header["Content-Type"][0] != "image/jpeg" {
		fmt.Printf("uploadImage >> Header != png || jpeg: %v\n", err)
		web.SetErrorRedirect(w, r, "/webmaster/gallery", "Error uploading file")
		return
	}
	f, err := os.OpenFile(path+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Printf("uploadImage >> OpenFile: %v\n", err)
		web.SetErrorRedirect(w, r, "/webmaster/gallery", "Error uploading file")
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
	web.SetSuccessRedirect(w, r, "/webmaster/gallery", "Successfully uploaded image")
	return
}}

var oneImage = web.Route{"GET", "/webmaster/gallery/:id", func(w http.ResponseWriter, r *http.Request) {
	images := db.GetAll("image")
	tmpl.Render(w, r, "webmaster-gallery.tmpl", web.Model{
		"images": db.GetAll("image"),
		"image":  db.Get("image", ParseId(r.FormValue(":id"))),
		"cats":   getCategories(images),
		"page":   "gallery",
	})
	return
}}

var saveImage = web.Route{"POST", "/webmaster/gallery/save/:id", func(w http.ResponseWriter, r *http.Request) {
	id := ParseId(r.FormValue(":id"))
	img := db.Get("image", id).Data
	img["category"] = r.FormValue("category")
	img["description"] = r.FormValue("description")
	db.Set("image", id, img)
	web.SetSuccessRedirect(w, r, "/webmaster/gallery", "Successfully saved image")
	return
}}

var deleteImage = web.Route{"POST", "/webmaster/gallery/:id", func(w http.ResponseWriter, r *http.Request) {
	img := db.Get("image", ParseId(r.FormValue(":id")))
	if img == nil {
		log.Printf("webmasterDeleteImage >> db.Get: image could not be found in database\n")
		web.SetErrorRedirect(w, r, "/webmaster/gallery", "Error deleting image.")
		return
	}
	if err := os.Remove("static/img/upload/" + img.Data["source"].(string)); err != nil {
		log.Printf("webmasterDeleteImage >> os.Remove: %v\n", err)
		web.SetErrorRedirect(w, r, "/webmaster/gallery", "Error deleting image.")
		return
	}
	db.Del("image", img.Id)
	web.SetSuccessRedirect(w, r, "/webmaster/gallery", "Successfully deleted image")
	return
}}

var allListings = web.Route{"GET", "/webmaster/listing", func(w http.ResponseWriter, r *http.Request) {
	tmpl.Render(w, r, "webmaster-listing.tmpl", web.Model{
		"listings": db.GetAll("listing"),
		"page":     "listings",
	})
	return
}}

var addListing = web.Route{"POST", "/webmaster/listing/add", func(w http.ResponseWriter, r *http.Request) {
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
	web.SetSuccessRedirect(w, r, "/webmaster/listing", "Successfuly added listing")
	return
}}

var oneListing = web.Route{"GET", "/webmaster/listing/:id", func(w http.ResponseWriter, r *http.Request) {
	tmpl.Render(w, r, "webmaster-listing.tmpl", web.Model{
		"listings": db.GetAll("listing"),
		"listing":  db.Get("listing", ParseId(r.FormValue(":id"))),
		"page":     "listings",
	})
	return
}}

var saveListing = web.Route{"POST", "/webmaster/listing/save/:id", func(w http.ResponseWriter, r *http.Request) {
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
	web.SetSuccessRedirect(w, r, "/webmaster/listing", "Successfuly saved listing")
	return
}}

var deleteListing = web.Route{"POST", "/webmaster/listing/:id", func(w http.ResponseWriter, r *http.Request) {
	db.Del("listing", ParseId(r.FormValue(":id")))
	web.SetSuccessRedirect(w, r, "/webmaster/listing", "Successfully deleted listing")
	return
}}

var allFloorplans = web.Route{"GET", "/webmaster/floorplan", func(w http.ResponseWriter, r *http.Request) {
	tmpl.Render(w, r, "webmaster-floorplan.tmpl", web.Model{
		"floorplans": GetFloorPlans(),
		"page":       "floorplans",
	})
	return
}}

var uploadFloorplan = web.Route{"POST", "/webmaster/floorplan/upload", func(w http.ResponseWriter, r *http.Request) {
	path := "static/floorplans/"
	if err := os.MkdirAll(path, 0755); err != nil {
		fmt.Printf("uploadFloorplan >> MkdirAll: %v\n", err)
		web.SetErrorRedirect(w, r, "/webmaster/floorplan", "Error uploading file")
		return
	}
	r.ParseMultipartForm(32 << 20) // 32 MB
	file, handler, err := r.FormFile("floorplan")
	if err != nil || len(handler.Header["Content-Type"]) < 1 {
		fmt.Printf("uploadFloorplan >> Header len < 1: %v\n", err)
		web.SetErrorRedirect(w, r, "/webmaster/floorplan", "Error uploading file")
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
		web.SetErrorRedirect(w, r, "/webmaster/floorplan", "Error uploading file")
		return
	}
	f, err := os.OpenFile(path+fileName, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Printf("uploadFloorplan >> OpenFile: %v\n", err)
		web.SetErrorRedirect(w, r, "/webmaster/floorplan", "Error uploading file")
		return
	}
	defer f.Close()
	io.Copy(f, file)
	web.SetSuccessRedirect(w, r, "/webmaster/floorplan", "Successfully uploaded image")
	return
}}

var renameFloorplan = web.Route{"POST", "/webmaster/floorplan/rename", func(w http.ResponseWriter, r *http.Request) {
	ext := "." + strings.Split(r.FormValue("oldName"), ".")[1]
	err := os.Rename("static/floorplans/"+r.FormValue("oldName"), "static/floorplans/"+r.FormValue("name")+ext)
	if err != nil {
		web.SetErrorRedirect(w, r, "/webmaster/floorplan", "Error renaming floorplan")
		return
	}
	web.SetSuccessRedirect(w, r, "/webmaster/floorplan", "Successfully renamed floorplan")
	return

}}

var deleteFloorplan = web.Route{"POST", "/webmaster/floorplan/:name", func(w http.ResponseWriter, r *http.Request) {
	err := os.Remove("static/floorplans/" + r.FormValue(":name"))
	if err != nil {
		web.SetErrorRedirect(w, r, "/webmaster/floorplan", "Error deleting floorplan")
		return
	}
	web.SetSuccessRedirect(w, r, "/webmaster/floorplan", "Successfully deleted floorplan")
	return
}}

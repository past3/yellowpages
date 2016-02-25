package main

import (
	"net/http"
)

//FrontAdminHandler for serving admin page
func FrontAdminHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "admin/index.html")
}

//AddListingViewHandler to render veiw page
func AddListingViewHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "admin/partials/addlisting.html")
}

//UnapprovedViewHandler to render veiw page
func UnapprovedViewHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "admin/partials/view_listings.html")
}

func addcatViewHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "admin/partials/addcat.html")
}

//ClientViewHandler app view for client
func ClientViewHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "client/index.html")
}

//CategoryHandler for view
func CategoryHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "client/partials/category.html")
}

//ResultHandler for result view
func ResultHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "client/partials/result.html")
}

//ListingHandler for listing view
func ListingHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "client/partials/list.html")
}

package models

import (
	"encoding/base64"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/goamz/goamz/aws"
	"github.com/goamz/goamz/s3"
	uuid "github.com/satori/go.uuid"
	"github.com/smilecs/yellowpages/config"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//Listing struct holds all submitted form data for listings
type Listing struct {
	ID             bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	CompanyName    string        `bson:"companyname"`
	Address        string        `bson:"address"`
	Hotline        string        `bson:"hotline"`
	HotlinesList   []string
	Specialisation string `bson:"specialisation"`
	Category       string `bson:"category"`
	// Advert         string    `bson:"advert"`
	// Size           string    `bson:"size"`
	Image    string    `bson:"image"`
	Images   []string  `bson:"images"`
	Slug     string    `bson:"slug"`
	About    string    `bson:"about"`
	Username string    `bson:"username"`
	Password string    `bson:"password"`
	RC       string    `bson:"rc"` //
	Branch   string    `bson:"branch"`
	Product  string    `bson:"product"`
	Email    string    `bson:"email"`
	Website  string    `bson:"website"`
	DHr      string    `bson:"dhr"` //Working days and hours
	Verified string    `bson:"verified"`
	Approved bool      `bson:"approved"`
	Plus     string    `bson:"plus"`
	Expiry   time.Time `bson:"expiry"`
	Duration string    `bson:"duration"`
	Date     time.Time `bson:"date"`
	Pg       Page
	Reviews  string `bson:"reviews"`
}

type Listings struct {
	Data  []Listing
	Page  Page
	Query string
}

//Addlisting function adding listings data to db
func (r Listing) Add(config *config.Conf) error {

	auth, err := aws.EnvAuth()
	if err != nil {
		log.Fatal(err)
	}
	client := s3.New(auth, aws.USWest2)
	bucket := client.Bucket("yellowpagesng")

	p := rand.New(rand.NewSource(time.Now().UnixNano()))
	str := strconv.Itoa(p.Intn(10))
	r.Date = time.Now()
	if r.Image != "" {
		tm, _ := strconv.Atoi(r.Duration)
		t := r.Date.AddDate(tm, 1, 0)
		r.Expiry = t
		byt, er := base64.StdEncoding.DecodeString(strings.Split(r.Image, "base64,")[1])
		if er != nil {
			log.Println(er)
		}

		meta := strings.Split(r.Image, "base64,")[0]
		newmeta := strings.Replace(strings.Replace(meta, "data:", "", -1), ";", "", -1)
		imagename := "listings/" + uuid.NewV1().String()

		err = bucket.Put(imagename, byt, newmeta, s3.PublicReadWrite, s3.Options{})
		if err != nil {
			log.Println(err)
		}

		log.Println(bucket.URL(imagename))
		r.Image = bucket.URL(imagename)
	}

	var images []string
	for _, v := range r.Images {
		var byt []byte
		byt, err = base64.StdEncoding.DecodeString(strings.Split(v, "base64,")[1])
		if err != nil {
			log.Println(err)
		}

		meta := strings.Split(v, "base64,")[0]
		newmeta := strings.Replace(strings.Replace(meta, "data:", "", -1), ";", "", -1)
		imagename := "listings/" + uuid.NewV1().String()

		err = bucket.Put(imagename, byt, newmeta, s3.PublicReadWrite, s3.Options{})
		if err != nil {
			log.Println(err)
		}
		images = append(images, bucket.URL(imagename))
	}

	r.Images = images

	r.Slug = strings.Replace(r.CompanyName, " ", "-", -1) + str
	r.Slug = strings.Replace(r.Slug, "&", "-", -1) + str
	r.Slug = strings.Replace(r.Slug, "/", "-", -1) + str
	r.Slug = strings.Replace(r.Slug, ",", "-", -1) + str
	index := mgo.Index{
		Key:        []string{"$text:specialisation", "$text:companyname"},
		Unique:     true,
		DropDups:   true,
		Background: true,
	}
	mgoSession := config.Database.Session.Copy()
	defer mgoSession.Close()

	collection := config.Database.C("Listings").With(mgoSession)
	collection.EnsureIndex(index)
	collection.Insert(r)
	return err
}

func (r Listing) Edit(config *config.Conf) error {
	log.Println(r)
	auth, err := aws.EnvAuth()
	if err != nil {
		log.Fatal(err)
	}
	client := s3.New(auth, aws.USWest2)
	bucket := client.Bucket("yellowpagesng")

	if r.Image != "" {
		if len(strings.Split(r.Image, "base64,")) > 1 {
			meta := strings.Split(r.Image, "base64,")[0]
			tm, _ := strconv.Atoi(r.Duration)
			t := r.Date.AddDate(tm, 1, 0)
			r.Expiry = t
			byt, er := base64.StdEncoding.DecodeString(strings.Split(r.Image, "base64,")[1])
			if er != nil {
				log.Println(er)
			}

			newmeta := strings.Replace(strings.Replace(meta, "data:", "", -1), ";", "", -1)
			imagename := "listings/" + uuid.NewV1().String()

			err = bucket.Put(imagename, byt, newmeta, s3.PublicReadWrite, s3.Options{})
			if err != nil {
				log.Println(err)
			}

			log.Println(bucket.URL(imagename))
			r.Image = bucket.URL(imagename)
		}
	}

	var images []string
	for _, v := range r.Images {

		if len(strings.Split(v, "base64,")) > 1 {

			var byt []byte
			byt, err = base64.StdEncoding.DecodeString(strings.Split(v, "base64,")[1])
			if err != nil {
				log.Println(err)
			}

			meta := strings.Split(v, "base64,")[0]
			newmeta := strings.Replace(strings.Replace(meta, "data:", "", -1), ";", "", -1)
			imagename := "listings/" + uuid.NewV1().String()

			err = bucket.Put(imagename, byt, newmeta, s3.PublicReadWrite, s3.Options{})
			if err != nil {
				log.Println(err)
			}
			images = append(images, bucket.URL(imagename))
		} else {
			images = append(images, v)
		}
	}

	r.Images = images

	index := mgo.Index{
		Key:        []string{"$text:specialisation", "$text:companyname"},
		Unique:     true,
		DropDups:   true,
		Background: true,
	}
	mgoSession := config.Database.Session.Copy()
	defer mgoSession.Close()

	collection := config.Database.C("Listings").With(mgoSession)
	collection.EnsureIndex(index)
	collection.Update(bson.M{"slug": r.Slug}, bson.M{
		"$set": bson.M{
			"companyname":    r.CompanyName,
			"address":        r.Address,
			"hotline":        r.Hotline,
			"specialisation": r.Specialisation,
			"category":       r.Category,
			"image":          r.Image,
			"images":         r.Images,
			"slug":           r.Slug,
			"about":          r.About,
			"rc":             r.RC,
			"branch":         r.Branch,
			"product":        r.Product,
			"email":          r.Email,
			"website":        r.Website,
			"dhr":            r.DHr,
		},
	})

	err = IndexSingleListingWithBleve(r)
	if err != nil {
		log.Println(err)
	}
	return err

}

func (r Listing) GetOne(config *config.Conf, id string) (Listing, error) {
	result := Listing{}
	mgoSession := config.Database.Session.Copy()
	defer mgoSession.Close()
	collection := config.Database.C("Listings").With(mgoSession)

	err := collection.Find(bson.M{"slug": id}).One(&result)
	if err != nil {
		log.Println(err)
		return result, err
	}
	return result, nil
}

func StringToPhoneNumbers(s string) []string {
	var response []string
	noCommaList := strings.Split(s, ",")
	for _, j := range noCommaList {
		noSpaceList := strings.Split(j, " ")
		for _, jj := range noSpaceList {
			if jj != " " && jj != "," && jj != "" {
				response = append(response, strings.TrimSpace(jj))
			}
		}
	}
	return response
}

//GetListings return listings
func (r Listing) GetAllApproved(config *config.Conf) (Listings, error) {
	result := []Listing{}
	listings := Listings{}
	mgoSession := config.Database.Session.Copy()
	defer mgoSession.Close()

	collection := config.Database.C("Listings").With(mgoSession)

	err := collection.Find(bson.M{"approved": true}).All(&result)
	if err != nil {
		return listings, err
	}
	for i := range result {
		result[i].HotlinesList = StringToPhoneNumbers(result[i].Hotline)
	}
	listings.Data = result
	return listings, nil
}

//Getunapproved function for Getunapproved handler
func (r Listing) GetAllUnapproved(config *config.Conf) (Listings, error) {
	result := []Listing{}
	listings := Listings{}
	mgoSession := config.Database.Session.Copy()
	defer mgoSession.Close()
	collection := config.Database.C("Listings").With(mgoSession)

	err := collection.Find(bson.M{"approved": false}).All(&result)
	if err != nil {
		return listings, err
	}
	for i := range result {
		result[i].HotlinesList = StringToPhoneNumbers(result[i].Hotline)
	}
	listings.Data = result
	return listings, nil
}

//GetcatListing function
func (r Listing) GetAllInCategory(config *config.Conf, id string, page int) (Listings, error) {
	perPage := 10
	//result := []Listing{}
	listings := Listings{}
	mgoSession := config.Database.Session.Copy()
	defer mgoSession.Close()

	collection := config.Database.C("Listings").With(mgoSession)

	q := collection.Find(bson.M{"category": id, "approved": true}).Sort("-plus")

	count, err := q.Count()
	if err != nil {
		return listings, err
	}

	pg := SearchPagination(count, page, perPage)
	err = q.Limit(perPage).Skip(pg.Skip).All(&listings.Data)

	for i := range listings.Data {
		listings.Data[i].HotlinesList = StringToPhoneNumbers(listings.Data[i].Hotline)
	}

	listings.Page = pg

	if err != nil {
		return listings, err
	}
	return listings, nil
}

//GetcatListing function
func (r Listing) GetAllPlusListings(config *config.Conf) (Listings, error) {
	result := []Listing{}
	listings := Listings{}

	mgoSession := config.Database.Session.Copy()
	defer mgoSession.Close()

	collection := config.Database.C("Listings").With(mgoSession)

	err := collection.Find(bson.M{"plus": "true"}).All(&result)
	if err != nil {
		return listings, err
	}
	for i := range result {
		result[i].HotlinesList = StringToPhoneNumbers(result[i].Hotline)
	}

	listings.Data = result
	return listings, nil
}

//Update to approve of a listing to show in the client side
func (r Listing) Approve(config *config.Conf, slug string) error {
	mgoSession := config.Database.Session.Copy()
	defer mgoSession.Close()
	collection := config.Database.C("Listings").With(mgoSession)

	query := bson.M{"slug": slug}
	change := bson.M{"$set": bson.M{"approved": true}}

	err := collection.Update(query, change)
	if err != nil {
		log.Println(err)
		return err
	}

	listing := Listing{}
	err = collection.Find(bson.M{"slug": slug}).One(&listing)
	if err != nil {
		log.Println(err)
	}
	err = IndexSingleListingWithBleve(listing)
	if err != nil {
		log.Println(err)
	}
	return nil
}

//Update to approve of a listing to show in the client side
func (r Listing) Delete(config *config.Conf, slug string) error {
	mgoSession := config.Database.Session.Copy()
	defer mgoSession.Close()
	collection := config.Database.C("Listings").With(mgoSession)

	query := bson.M{"slug": slug}

	err := collection.Remove(query)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

//TimeUpdate to change time
func (r Listing) TimeUpdate(config *config.Conf, id string, expiry string) error {

	x, _ := strconv.Atoi(expiry)

	mgoSession := config.Database.Session.Copy()
	defer mgoSession.Close()

	collection := config.Database.C("Listings").With(mgoSession)

	query := bson.M{"slug": id}
	change := bson.M{"$set": bson.M{"expiry": time.Now().AddDate(x, 1, 0)}}

	err := collection.Update(query, change)

	if err != nil {
		return err
	}

	return nil
}

//Search searches
func (r Listing) Search(config *config.Conf, query string, page int) (Listings, error) {
	perPage := 10
	Results := Listings{}

	index := mgo.Index{
		Key: []string{"$text:specialisation", "$text:companyname"},
	}

	mgoSession := config.Database.Session.Copy()
	defer mgoSession.Close()

	collection := config.Database.C("Listings").With(mgoSession)

	err := collection.EnsureIndex(index)
	if err != nil {
		return Results, err
	}

	searchResult, err := SearchWithIndex(query)
	if err != nil {
		log.Println(err)
	}
	log.Printf("%+v", searchResult.Total)

	documentIdArray := []string{}
	for _, v := range searchResult.Hits {
		documentIdArray = append(documentIdArray, v.ID)

	}
	log.Println(documentIdArray)
	q := collection.Find(
		bson.M{
			"slug": bson.M{
				"$in": documentIdArray,
			},
			"approved": true,
		}).Sort("-plus")

	count, err := q.Count()
	if err != nil {
		return Results, err
	}
	pg := SearchPagination(count, page, perPage)
	err = q.Limit(perPage).Skip(pg.Skip).All(&Results.Data)
	Results.Page = pg
	if err != nil {
		return Results, err
	}
	return Results, nil
}

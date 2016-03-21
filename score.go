package main
// REST API to create, retrieve, update and delete scores
import (
//  "bufio"
//  "bytes"
	"fmt"
	"encoding/json"  //implements encoding and decoding of JSON objects
//  "encoding/hex"
	"net/http"
	"reflect"
	// implements run-time reflection, allowing a program to manipulate 
	// objects with arbitrary types. The typical use is to take a value 
	// with static type interface{} and extract its dynamic type information 
	// by calling TypeOf, which returns a Type.  
//  "image/jpeg"
  "strings"
//  "text/template"
	"html/template"
	"time"
//  "path"
  "log"
//  "io/ioutil"
//  "io"  
//  "os"
//  "encoding/base64"
	"github.com/gorilla/context" //request and response mapping
	"github.com/julienschmidt/httprouter"  //to handle r&r from gorrilla
	"github.com/justinas/alice" //used to chain handlers
	"gopkg.in/mgo.v2"  //driver for mongodb
	"gopkg.in/mgo.v2/bson"  //implementation of the BSON specification for Go
)
// Nicolas Merouze
// Repo BSON spec, mongodb driver

var listEvent = template.Must(template.ParseFiles("templates/base.html", "templates/events/list/list.html"))
var createnewEvent = template.Must(template.ParseFiles("templates/base.html", "templates/events/new/new.html", "templates/events/form.html"))
var updateEvent = template.Must(template.ParseFiles("templates/base.html", "templates/events/form.html", "templates/events/update/update.html"))
var showEvent = template.Must(template.ParseFiles("templates/base.html", "templates/events/show/show.html"))

var getEntrantsNew = template.Must(template.ParseFiles("templates/base.html", "templates/getEntrants/new.html"))
var getEntrantsUpdate = template.Must(template.ParseFiles("templates/base.html", "templates/getEntrants/edit.html"))

var listEntrant = template.Must(template.ParseFiles("templates/base.html", "templates/entrants/list/list.html"))
var createnewEntrant = template.Must(template.ParseFiles("templates/base.html", "templates/entrants/new/new.html", "templates/entrants/form.html"))
var updateEntrant = template.Must(template.ParseFiles("templates/base.html", "templates/entrants/update/update.html", "templates/entrants/form.html"))
var showEntrant = template.Must(template.ParseFiles("templates/base.html", "templates/entrants/show/show.html"))

var listScorecard = template.Must(template.ParseFiles("templates/base.html", "templates/scorecards/list/list.html"))
var createnewScorecard = template.Must(template.ParseFiles("templates/base.html", "templates/scorecards/new/new.html", "templates/scorecards/form.html"))
var updateScorecard = template.Must(template.ParseFiles("templates/base.html", "templates/scorecards/update/update.html", "templates/scorecards/form.html"))
var showScorecard = template.Must(template.ParseFiles("templates/base.html", "templates/scorecards/show/show.html"))

//Event collection////////////////////////////////////////////////////////////////////////////////////


type Event struct {
	Id                  bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	Name                string        `json:"name"`
	Location            string        `json:"location"`
  Division            string        `json:"division"`
  Int_search_areas    string        `json:"int_search_areas"`  
  Ext_search_areas    string        `json:"ext_search_areas"`    
  Cont_search_areas   string        `json:"cont_search_areas"`    
  Veh_search_areas    string        `json:"veh_search_areas"`    
  Elite_search_areas  string        `json:"elite_search_areas"`  
  Int_hides           string        `json:"int_hides"` 
  Ext_hides           string        `json:"ext_hides"`   
  Cont_hides          string        `json:"cont_hides"`   
  Veh_hides           string        `json:"veh_hides"`   
  Elite_hides         string        `json:"elite_hides"`  
  Event_Id            string        `json:"event_id"`
  EntrantAll_Id       []string      `json:"entrantall_id"`
  EntrantSelected_Id  []string      `json:"entrantselected_id"`
}

type EventsCollection struct {
	Data []Event `json:"data"`
}

type EventResource struct {
	Data Event `json:"evdata"`
}

type EventRepo struct {
	coll *mgo.Collection
}

func (r *EventRepo) All() (EventsCollection, error) {
	fmt.Println("In Event All")
	result := EventsCollection{[]Event{}}
	err := r.coll.Find(nil).All(&result.Data)
	if err != nil {
    fmt.Println("out of Event All")
		return result, err
	}
  fmt.Println("out of Event All")
	return result, nil
}

func (r *EventRepo) Find(id string) (EventResource, error) {
	fmt.Println("In Event Find")
  prefix := "ObjectIdHex(\""
  suffix := "\")"
  id = strings.TrimPrefix(id, prefix)
  id = strings.TrimSuffix(id, suffix)
  idb := bson.ObjectIdHex(id)
  fmt.Println("printing idb")
  fmt.Println(idb)
	result := EventResource{}  
  err := r.coll.Find(bson.M{"_id": idb}).One(&result.Data)
  fmt.Println("printing result.Data")
  fmt.Println(result.Data)
	if err != nil {
    fmt.Println("out of Event Find")
		return result, err
	}
	fmt.Println("out of Event Find")
	return result, nil
}

func (r *EventRepo) Create(event *Event) (error, bson.ObjectId) {
	fmt.Println("In Event Create")
	id := bson.NewObjectId()
	_, err := r.coll.UpsertId(id, event)
	if err != nil {
		fmt.Println("out of Event Create")
    return err, id
	}
	event.Id = id
  fmt.Println("out of Event Create")
	return err, id
}

func (r *EventRepo) Update(event *Event) error {
	fmt.Println("In Event Update")
	result := EventResource{}  
  err := r.coll.Find(bson.M{"_id": event.Id}).One(&result.Data)
	if err != nil {
    fmt.Println("Find error")
    fmt.Println("out of Event Update")
		return err
	}  
  err = r.coll.Update(result.Data, event)
	if err != nil {
    fmt.Println("Update error")
    fmt.Println("out of Event Update")
		return err
	}
	fmt.Println("out of Event Update")
	return nil
}

func (r *EventRepo) Delete(id string) error {
  fmt.Println("In Event Delete")
  prefix := "ObjectIdHex(\""
  suffix := "\")"
  id = strings.TrimPrefix(id, prefix)
  id = strings.TrimSuffix(id, suffix)
  idb := bson.ObjectIdHex(id) 
	result := EventResource{}  
  err := r.coll.Find(bson.M{"_id": idb}).One(&result.Data)  
	if err != nil {
    fmt.Println("error in find")
    fmt.Println("out of Event Delete")
		return err
	}
  err = r.coll.Remove(result.Data)
	if err != nil {
    fmt.Println("error in Delete")
    fmt.Println("out of Event Delete")
		return err
	}  
  fmt.Println("out of Event Delete")
	return nil
}

// Entrant collection  //////////////////////////////////////////////////////////////////////////////

type Entrant struct {
	Id bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	Name       string        `json:"name"`
	Team_Id    string        `json:"entrant_id"`
  Event_Id   []string      `json:"event_id"`
}

type EntrantsCollection struct {
	Data []Entrant `json:"data"`
}

type EntrantResource struct {
	Data Entrant `json:"data"`
}

type EntrantRepo struct {
	coll *mgo.Collection
}

func (r *EntrantRepo) All() (EntrantsCollection, error) {
	fmt.Println("In Entrant All")
	result := EntrantsCollection{[]Entrant{}}
	err := r.coll.Find(nil).All(&result.Data)
	if err != nil {
    fmt.Println("out of Entrant All")
		return result, err
	}
  fmt.Println("out of Entrant All")
	return result, nil
}

func (r *EntrantRepo) Find(id string) (EntrantResource, error) {
	fmt.Println("In Entrant Find")
  prefix := "ObjectIdHex(\""
  suffix := "\")"
  id = strings.TrimPrefix(id, prefix)
  id = strings.TrimSuffix(id, suffix)
  idb := bson.ObjectIdHex(id) 
	result := EntrantResource{}  
  err := r.coll.Find(bson.M{"_id": idb}).One(&result.Data)
	if err != nil {
    fmt.Println("out of Entrant Find")
		return result, err
	}
	fmt.Println("out of Entrant Find")
	return result, nil
}

func (r *EntrantRepo) Create(entrant *Entrant) (error, bson.ObjectId) {
	fmt.Println("In Entrant Create")
	id := bson.NewObjectId()
	_, err := r.coll.UpsertId(id, entrant)
	if err != nil {
		fmt.Println("out of Entrant Create")
    return err, id
	}
	entrant.Id = id
  fmt.Println("out of Entrant Create")
	return err, id
}

func (r *EntrantRepo) Update(entrant *Entrant) error {
	fmt.Println("In Entrant Update")
	result := EntrantResource{}  
  err := r.coll.Find(bson.M{"_id": entrant.Id}).One(&result.Data)
	if err != nil {
    fmt.Println("Find error")
    fmt.Println("out of Entrant Update")
		return err
	}  
  err = r.coll.Update(result.Data, entrant)
	if err != nil {
    fmt.Println("Update error")
    fmt.Println("out of Entrant Update")
		return err
	}
	fmt.Println("out of Entrant Update")
	return nil
}

func (r *EntrantRepo) Delete(id string) error {
  fmt.Println("In Entrant Delete")
  prefix := "ObjectIdHex(\""
  suffix := "\")"
  id = strings.TrimPrefix(id, prefix)
  id = strings.TrimSuffix(id, suffix)
  idb := bson.ObjectIdHex(id) 
	result := EntrantResource{}  
  err := r.coll.Find(bson.M{"_id": idb}).One(&result.Data)  
	if err != nil {
    fmt.Println("error in find")
    fmt.Println("out of Entrant Delete")
		return err
	}
  err = r.coll.Remove(result.Data)
	if err != nil {
    fmt.Println("error in Delete")
    fmt.Println("out of Entrant Delete")
		return err
	}  
  fmt.Println("out of Entrant Delete")
	return nil
}

//Scorecard collection////////////////////////////////////////////////////////////////////////////////////


type Scorecard struct {
	Id bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	Element                   string        `json:"element"`
  Maxtime_m                 string        `json:"maxtime_m"`
  Maxtime_s                 string        `json:"maxtime_s"`
  Maxtime_ms                string        `json:"maxtime_ms"`
  Finish_call               string        `json:"finish_call"`
  False_alert_fringe        string        `json:"false_alert_fringe"`
  Timed_out                 string        `json:"timed_out"`
  Dismissed                 string        `json:"dismissed"`
  Excused                   string        `json:"excused"`
  Absent                    string        `json:"absent"`
  Eliminated_during_search  string        `json:"eliminated_during_search"`
  Other_faults_descr        string        `json:"other_faults_descr"`
  Other_faults_count        string        `json:"other_faults_count"`
  Comments                  string        `json:"comments"`
  Time_elapsed_m            string        `json:"time_elapsed_m"`  
  Time_elapsed_s            string        `json:"time_elapsed_s"`
  Time_elapsed_ms           string        `json:"time_elapsed_ms"`
  Pronounced                string        `json:"pronounced"`
  Judge_signature           string        `json:"judge_signature"`
  Scorecard_Id              string        `json:"scorecard_id"`
  Event_Id                  string        `json:"event_id"`
  Entrant_Id                string        `json:"entrant_id"`
  Search_area               string        `json:"search_area"`
  Hides_max                 string        `json:"hides_max"`
  Hides_found               string        `json:"hides_found"`
  Hides_missed              string        `json:"hides_missed"`
  Total_faults              string        `json:"total_faults"`
  Maxpoint                  string        `json:"maxpoint"`
  Total_points              string        `json:"total_points"`  
}

type ScorecardsCollection struct {
	Data []Scorecard `json:"data"`
}

type ScorecardResource struct {
	Data Scorecard `json:"evdata"`
}

type ScorecardRepo struct {
	coll *mgo.Collection
}

func (r *ScorecardRepo) All() (ScorecardsCollection, error) {
	fmt.Println("In Scorecard All")
	result := ScorecardsCollection{[]Scorecard{}}
	err := r.coll.Find(nil).All(&result.Data)
	if err != nil {
    fmt.Println("out of Scorecard All")
		return result, err
	}
  fmt.Println("out of Scorecard All")
	return result, nil
}

func (r *ScorecardRepo) Find(id string) (ScorecardResource, error) {
	fmt.Println("In Scorecard Find")
  prefix := "ObjectIdHex(\""
  suffix := "\")"
  id = strings.TrimPrefix(id, prefix)
  id = strings.TrimSuffix(id, suffix)
  idb := bson.ObjectIdHex(id)
  fmt.Println("printing idb")
  fmt.Println(idb)
	result := ScorecardResource{}  
  err := r.coll.Find(bson.M{"_id": idb}).One(&result.Data)
  fmt.Println("printing result.Data")
  fmt.Println(result.Data)
	if err != nil {
    fmt.Println("out of Scorecard Find")
		return result, err
	}
	fmt.Println("out of Scorecard Find")
	return result, nil
}

func (r *ScorecardRepo) Create(scorecard *Scorecard) (error, bson.ObjectId) {
	fmt.Println("In Scorecard Create")
	id := bson.NewObjectId()
	_, err := r.coll.UpsertId(id, scorecard)
	if err != nil {
		fmt.Println("out of Scorecard Create")
    return err, id
	}
	scorecard.Id = id
  fmt.Println("out of Scorecard Create")
	return err, id
}

func (r *ScorecardRepo) Update(scorecard *Scorecard) error {
	fmt.Println("In Scorecard Update")
	result := ScorecardResource{}  
  err := r.coll.Find(bson.M{"_id": scorecard.Id}).One(&result.Data)
	if err != nil {
    fmt.Println("Find error")
    fmt.Println("out of Scorecard Update")
		return err
	}  
  err = r.coll.Update(result.Data, scorecard)
	if err != nil {
    fmt.Println("Update error")
    fmt.Println("out of Scorecard Update")
		return err
	}
	fmt.Println("out of Scorecard Update")
	return nil
}

func (r *ScorecardRepo) Delete(id string) error {
  fmt.Println("In Scorecard Delete")
  prefix := "ObjectIdHex(\""
  suffix := "\")"
  id = strings.TrimPrefix(id, prefix)
  id = strings.TrimSuffix(id, suffix)
  idb := bson.ObjectIdHex(id) 
	result := ScorecardResource{}  
  err := r.coll.Find(bson.M{"_id": idb}).One(&result.Data)  
	if err != nil {
    fmt.Println("error in find")
    fmt.Println("out of Scorecard Delete")
		return err
	}
  err = r.coll.Remove(result.Data)
	if err != nil {
    fmt.Println("error in Delete")
    fmt.Println("out of Scorecard Delete")
		return err
	}  
  fmt.Println("out of Scorecard Delete")
	return nil
}

// Errors////////////////////////////////////////////////////////////////////////////////////////////

type Errors struct {
	Errors []*Error `json:"errors"`
}

type Error struct {
	Id     string `json:"id"`
	Status int    `json:"status"`
	Title  string `json:"title"`
	Detail string `json:"detail"`
}

func WriteError(w http.ResponseWriter, err *Error) {
  fmt.Println("In WriteError")
	w.Header().Set("Content-Type", "application/vnd.api+json")
	w.WriteHeader(err.Status)
	json.NewEncoder(w).Encode(Errors{[]*Error{err}})
  fmt.Println("Out of WriteError")
}

var (
	ErrBadRequest           = &Error{"bad_request", 400, "Bad request", "Request body is not well-formed. It must be JSON."}
	ErrNotAcceptable        = &Error{"not_acceptable", 406, "Not Acceptable", "Accept header must be set to 'application/vnd.api+json'."}
	ErrUnsupportedMediaType = &Error{"unsupported_media_type", 415, "Unsupported Media Type", "Content-Type header must be set to: 'application/vnd.api+json'."}
	ErrInternalServer       = &Error{"internal_server_error", 500, "Internal Server Error", "Something went wrong."}
)

// Middlewares//////////////////////////////////////////////////////////////////////////////////////////
// go net/http

func recoverHandler(next http.Handler) http.Handler {
	fmt.Println("In recoverHandler")
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
        fmt.Println("out of recoverHandler")
				log.Printf("panic: %+v", err)
				WriteError(w, ErrInternalServer)
				fmt.Println("in defer")
			}
		}()

		next.ServeHTTP(w, r)
	}
	fmt.Println("out of recoverHandler")
	return http.HandlerFunc(fn)
}
//  go net/http
func loggingHandler(next http.Handler) http.Handler {
	fmt.Println("In logging handler")
	fn := func(w http.ResponseWriter, r *http.Request) {
		t1 := time.Now()
		next.ServeHTTP(w, r)
		t2 := time.Now()
		log.Printf("[%s] %q %v\n", r.Method, r.URL.String(), t2.Sub(t1))
	}
	fmt.Println("out of logging handler")
	return http.HandlerFunc(fn)
}
//  go net/http
func acceptHandler(next http.Handler) http.Handler {
	fmt.Println("In acceptHandler")
	
	fn := func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Accept") != "application/vnd.api+json" {
      fmt.Println("out of acceptHandler")
      WriteError(w, ErrNotAcceptable)
			return
		}
		next.ServeHTTP(w, r)
	}
	fmt.Println("out of acceptHandler")
	return http.HandlerFunc(fn)
}
//  go net/http
func contentTypeHandler(next http.Handler) http.Handler {
	fmt.Println("In contentTypeHandler")
	fn := func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/vnd.api+json" {
      fmt.Println("out of contentTypeHandler")
      WriteError(w, ErrUnsupportedMediaType)
			return
		}
		next.ServeHTTP(w, r)
	}
	fmt.Println("out of contentTypeHandler")
	return http.HandlerFunc(fn)
}
//  go net/http, reflect, gorilla context, mongodb driver
func bodyHandler(v interface{}) func(next http.Handler) http.Handler {
	fmt.Println("In bodyHandler")	
	t := reflect.TypeOf(v)                        //type interface{} which may be empty
	m := func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {    
      val := reflect.New(t).Interface()   //val is type interface{}      
//      err := json.NewDecoder(r.Body).Decode(val)  //r.Body is the request body and is type interface io.ReadCloser
//      err := json.NewDecoder(strings.NewReader(evj)).Decode(val)      
//      val = evj
//      fmt.Println("printing err")
//			fmt.Println(err)
//      if err != nil {
//				WriteError(w, ErrBadRequest)
//				return
//			}
			if next != nil {
				context.Set(r, "body", val)     //gorilla context, key "body": val, val is type interface{}  "body" will now retrieve val
        next.ServeHTTP(w, r)
				fmt.Println("Key set")
        fmt.Println("Printing body.value")
        fmt.Println(val)
			}
		}
    fmt.Println("out of bodyHandler")
    return http.HandlerFunc(fn)
  }
	fmt.Println("out of bodyHandler")
	return m
}

// Main handlers /////////////////////////////////////////////////////////////////////////////////////
// gorilla/context bound to mongo db

// MGO Database Type //////////////////////////////////////////////////////////////////////////////////

type appContext struct {
	db *mgo.Database
}


// Event Handlers /////////////////////////////////////////////////////////////////////////////////////

func (c *appContext) eventsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("In eventsHandler")
	repo := EventRepo{c.db.C("events")}
	events, err := repo.All()
	if err != nil {
    fmt.Println("Out of eventsHandler")
		panic(err)
	}
  //	w.Header().Set("Content-Type", "application/vnd.api+json")
  //	json.NewEncoder(w).Encode(events)
	// read BSON into JSON
  if err = listEvent.Execute(w, events.Data); err != nil {
      http.Error(w, err.Error(), http.StatusInternalServerError)
      fmt.Println("Out of eventsHandler")
      return
  }
  fmt.Println("Out of eventsHandler")
}

func (c *appContext) eventHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("In eventHandler")
  params := context.Get(r, "params").(httprouter.Params)  //gorrila context, key "params"
  repo := EventRepo{c.db.C("events")}
	event, err := repo.Find(params.ByName("id")) //getting data from named param :id
  if err != nil {
    fmt.Println("out of eventHandler") 
		panic(err)
	}
  //  w.Header().Set("Content-Type", "application/vnd.api+json")
  //	json.NewEncoder(w).Encode(event)  
  //	if err = show.Execute(w, json.NewEncoder(w).Encode(event)); err != nil {
  //      http.Error(w, err.Error(), http.StatusInternalServerError)
  //      return
  //  }
	// read JSON into BSON 
	if err = showEvent.Execute(w, event.Data); err != nil {
      fmt.Println("out of eventHandler") 
      http.Error(w, err.Error(), http.StatusInternalServerError)
      return
  }  
  fmt.Println("out of eventHandler")  
}

func (c *appContext) newEventHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("In newEventHandler")
  v := r.URL.Query()
  fmt.Println("printing v") 
  fmt.Println(v)   
  body := context.Get(r, "body").(*EventResource)
  fmt.Println("printing &body.Data") 
  fmt.Println(&body.Data)
  repo := EventRepo{c.db.C("events")}
  err, id := repo.Create(&body.Data)
  fmt.Println(err)
  event, err := repo.Find(id.Hex())
  for i:=0; i<len(v["Team_Id"]); i++{
    event.Data.EntrantAll_Id = append(event.Data.EntrantAll_Id, v["Team_Id"][i])
  }

  fmt.Println("printing event.Data")
  fmt.Println(event.Data)
	if err := createnewEvent.Execute(w, event.Data); err != nil {
      fmt.Println("out of newEventHandler") 
      http.Error(w, err.Error(), http.StatusInternalServerError)
      return
  }
  fmt.Println("out of newEventHandler") 
}

//func (c *appContext) createEventHandler(w http.ResponseWriter, r *http.Request) {
//	fmt.Println("In createEventHandler")  
//  body := context.Get(r, "body").(*EventResource)    //gorilla context, key "body" that returns val
//  body.Data.Name = r.FormValue("Name")
//  body.Data.Location = r.FormValue("Location")  
//  body.Data.Division = r.FormValue("Division")
//  body.Data.Event_Id = r.FormValue("Event_Id")
//  err := r.ParseForm()
//  fmt.Println(err)
//  body.Data.EntrantSelected_Id = r.Form["Team_Id"]
//  repo := EventRepo{c.db.C("events")}	
//  err, id := repo.Create(&body.Data) 
//	if err != nil {
//    fmt.Println("out of createEventHandler")
//		panic(err)
//	}
  //	w.Header().Set("Content-Type", "application/vnd.api+json")
  //	w.WriteHeader(201)
	//json.NewEncoder(w).Encode(body)
	// read JSON into BSON
  //  if err := createnew.Execute(w, body.Data); err != nil {
  //      fmt.Println("out of createEventHandler")
  //      http.Error(w, err.Error(), http.StatusInternalServerError)
  //      return
  //  }
  //  _, err = http.Get("/events/show/:" + id.Hex())
//	fmt.Println("out of createEventHandler")  
//  http.Redirect(w, r, "/events/show/" + id.Hex(), 302)
//}

func (c *appContext) editEventHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("In editEventHandler")
  v := r.URL.Query()
  fmt.Println("printing v") 
  fmt.Println(v)  
  params := context.Get(r, "params").(httprouter.Params)    //gorilla context, key "params"
  evRepo := EventRepo{c.db.C("events")}
 	event, err := evRepo.Find(params.ByName("id")) //getting data from named param :id  
//  fmt.Println(event.Data)
  for i:=0; i<len(v["Team_Id"]); i++{
    fmt.Println("printing v[\"Team_Id\"][i]")
    fmt.Println(v["Team_Id"][i])
    event.Data.EntrantAll_Id = append(event.Data.EntrantAll_Id, v["Team_Id"][i])
  }  
  fmt.Println("printing event.Data")
  fmt.Println(event.Data)
  if err = updateEvent.Execute(w, event.Data); err != nil {
      fmt.Println("out of editEventHandler")
      http.Error(w, err.Error(), http.StatusInternalServerError)
      return
  }
	fmt.Println("out of editEventHandler")
}

func (c *appContext) updateEventHandler(w http.ResponseWriter, r *http.Request) {
  fmt.Println("In updateEventHandler") 
	params := context.Get(r, "params").(httprouter.Params)  
 	evRepo := EventRepo{c.db.C("events")}
  event, err := evRepo.Find(params.ByName("id")) //getting data from named param :id
  body := context.Get(r, "body").(*EventResource)    
  body.Data.Id = event.Data.Id
  body.Data.Name = r.FormValue("Name")
  body.Data.Location = r.FormValue("Location")  
  body.Data.Division = r.FormValue("Division")
  body.Data.Event_Id = r.FormValue("Event_Id")
  body.Data.Int_search_areas = r.FormValue("Int_search_areas")
  body.Data.Ext_search_areas = r.FormValue("Ext_search_areas")
  body.Data.Cont_search_areas = r.FormValue("Cont_search_areas")
  body.Data.Veh_search_areas = r.FormValue("Veh_search_areas")
  body.Data.Elite_search_areas = r.FormValue("Elite_search_areas")
  body.Data.Int_hides = r.FormValue("Int_hides")
  body.Data.Ext_hides = r.FormValue("Ext_hides")
  body.Data.Cont_hides = r.FormValue("Cont_hides")
  body.Data.Veh_hides = r.FormValue("Veh_hides")
  body.Data.Elite_hides = r.FormValue("Elite_hides")
  body.Data.EntrantSelected_Id = r.Form["EntrantSelected_Id"]
	body.Data.EntrantAll_Id = r.Form["EntrantAll_Id"]
  
  err = evRepo.Update(&body.Data)
  
  if err != nil {
    fmt.Println("out of updateEventHandler")
		panic(err)
	}
  //  w.WriteHeader(204)
  //	w.Write([]byte("\n"))	
	fmt.Println("out of updateEventHandler")
  c.updateEntrantsEventHandler(w, r)
  if body.Data.Event_Id == ""{
    http.Redirect(w, r, "/events/delete/" + body.Data.Id.Hex(), 302)
  }else{
    http.Redirect(w, r, "/events/show/" + body.Data.Id.Hex(), 302)   
  }
}

func (c *appContext) deleteEventHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("In deleteEventHandler")
	params := context.Get(r, "params").(httprouter.Params)    //gorilla context, key "params"
	repo := EventRepo{c.db.C("events")}
	err := repo.Delete(params.ByName("id"))
	if err != nil {
    fmt.Println("out of deleteEventHandler")
		panic(err)
	}
  //	w.WriteHeader(204)
  //	w.Write([]byte("\n"))
	fmt.Println("out of deleteEventHandler")
  //  _, err = http.Get("/events")
  http.Redirect(w, r, "/events", 302)
}


// EventEntrants Handlers /////////////////////////////////////////////////////////////////////////////////////


func (c *appContext) newEventEntrantsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("In newEventEntrantsHandler")
	repo := EntrantRepo{c.db.C("entrants")}
	entrants, err := repo.All()
  fmt.Println("Printing entrants.Data")
  fmt.Println(entrants.Data)
  body := "/?"
  for i:=0; i<len(entrants.Data); i++{
    if i == 0{
      body = body + "Id=" + entrants.Data[i].Id.Hex() + "&Name=" + entrants.Data[i].Name + "&Team_Id=" + entrants.Data[i].Team_Id
    }else{
       body = body + "&Id=" + entrants.Data[i].Id.Hex() + "&Name=" + entrants.Data[i].Name + "&Team_Id=" + entrants.Data[i].Team_Id
    }
  }
  fmt.Println("Printing body")
  fmt.Println(body)   
  if err != nil {
    fmt.Println("Out of newEventEntrantsHandler")
		panic(err)
	}
//  if err = getEntrantsNew.Execute(w, entrants.Data); err != nil {
//      http.Error(w, err.Error(), http.StatusInternalServerError)
//      fmt.Println("Out of newEventEntrantsHandler")
//      return
//  }
  fmt.Println("Out of newEventEntrantsHandler")     
  http.Redirect(w, r, "/events/new" + body, 302)
}

func (c *appContext) updateEventEntrantsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("In updateEventEntrantsHandler")
	repo := EntrantRepo{c.db.C("entrants")}
	entrants, err := repo.All() 
  fmt.Println("Printing entrants.Data")
  fmt.Println(entrants.Data)
  params := context.Get(r, "params").(httprouter.Params)  
 	evRepo := EventRepo{c.db.C("events")}
  event, err := evRepo.Find(params.ByName("id")) //getting data from named param :id 
  body := event.Data.Id.Hex() + "/?"
  for i:=0; i<len(entrants.Data); i++{
    if i == 0{
      body = body + "Id=" + entrants.Data[i].Id.Hex() + "&Name=" + entrants.Data[i].Name + "&Team_Id=" + entrants.Data[i].Team_Id
    }else{
       body = body + "&Id=" + entrants.Data[i].Id.Hex() + "&Name=" + entrants.Data[i].Name + "&Team_Id=" + entrants.Data[i].Team_Id
    }
  }
  fmt.Println("Printing body")
  fmt.Println(body)  
  if err != nil {
    fmt.Println("Out of updateEventEntrantsHandler")
		panic(err)
	}
  fmt.Println("Out of updateEventEntrantsHandler")    
  http.Redirect(w, r, "/events/edit/" + body, 302)
}


// EntrantsEvent Handlers /////////////////////////////////////////////////////////////////////////////////////

func (c *appContext) updateEntrantsEventHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("In updateEntrantsEventHandler")
	repo := EntrantRepo{c.db.C("entrants")}
	entrants, err := repo.All() 
  fmt.Println("Printing entrants.Data")
  fmt.Println(entrants.Data)
  params := context.Get(r, "params").(httprouter.Params)  
 	evRepo := EventRepo{c.db.C("events")}
  event, err := evRepo.Find(params.ByName("id")) //getting data from named param :id 
  for i:=0; i<len(event.Data.EntrantSelected_Id); i++{
    for j:=0; j<len(entrants.Data); j++{     
      if entrants.Data[j].Team_Id == event.Data.EntrantSelected_Id[i]{
        entrants.Data[j].Event_Id = append(entrants.Data[j].Event_Id, event.Data.Event_Id)
        fmt.Println("printing entrants.Data[j].Event_Id")
        fmt.Println(entrants.Data[j].Event_Id)
        fmt.Println("printing &entrants.Data[j]")
        fmt.Println(&entrants.Data[j])        
        err = repo.Update(&entrants.Data[j])
      }
    }
  }
  fmt.Println("printing entrants.Data")
  fmt.Println(entrants.Data)
  if err != nil {
    fmt.Println("Out of updateEntrantsEventHandler")
		panic(err)
	}
  fmt.Println("Out of updateEntrantsEventHandler")    
  return
}



// EventScorecards Handlers /////////////////////////////////////////////////////////////////////////////////////




// Entrant Handlers /////////////////////////////////////////////////////////////////////////////////////

func (c *appContext) entrantsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("In entrantsHandler")
	repo := EntrantRepo{c.db.C("entrants")}
	entrants, err := repo.All()
	if err != nil {
    fmt.Println("Out of entrantsHandler")
		panic(err)
	}
  if err = listEntrant.Execute(w, entrants.Data); err != nil {
      http.Error(w, err.Error(), http.StatusInternalServerError)
      fmt.Println("Out of entrantsHandler")
      return
  }
  fmt.Println("Out of entrantsHandler")
}

func (c *appContext) entrantHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("In entrantHandler")
  params := context.Get(r, "params").(httprouter.Params)  //gorrila context, key "params"
  repo := EntrantRepo{c.db.C("entrants")}
	entrant, err := repo.Find(params.ByName("id")) //getting data from named param :id
  if err != nil {
    fmt.Println("out of entrantHandler") 
		panic(err)
	}
  fmt.Println("printing entrant.Data")
  fmt.Println(entrant.Data)
	if err = showEntrant.Execute(w, entrant.Data); err != nil {
      fmt.Println("out of entrantHandler") 
      http.Error(w, err.Error(), http.StatusInternalServerError)
      return
  }  
  fmt.Println("out of entrantHandler")  
}

func newEntrantHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("In newEntrantHandler")

	if err := createnewEntrant.Execute(w, nil); err != nil {
      fmt.Println("out of newEntrantHandler") 
      http.Error(w, err.Error(), http.StatusInternalServerError)
      return
  }
  fmt.Println("out of newEntrantHandler") 
}

func (c *appContext) createEntrantHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("In createEntrantHandler")
  body := context.Get(r, "body").(*EntrantResource)    //gorilla context, key "body" that returns val
  body.Data.Name = r.FormValue("Name")
  body.Data.Team_Id = r.FormValue("Team_Id")
  repo := EntrantRepo{c.db.C("entrants")}	
  err, id := repo.Create(&body.Data)
	if err != nil {
    fmt.Println("out of createEntrantHandler")
		panic(err)
	}
	fmt.Println("out of createEntrantHandler")
  http.Redirect(w, r, "/entrants/show/" + id.Hex(), 302)
}

func (c *appContext) editEntrantHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("In editEntrantHandler")
	params := context.Get(r, "params").(httprouter.Params)    //gorilla context, key "params"
	repo := EntrantRepo{c.db.C("entrants")}
 	entrant, err := repo.Find(params.ByName("id")) //getting data from named param :id 
  if err = updateEntrant.Execute(w, entrant.Data); err != nil {
      fmt.Println("out of editEntrantHandler")
      http.Error(w, err.Error(), http.StatusInternalServerError)
      return
  }  
	fmt.Println("out of editEntrantHandler")
}

func (c *appContext) updateEntrantHandler(w http.ResponseWriter, r *http.Request) {
  fmt.Println("In updateEntrantHandler")
	params := context.Get(r, "params").(httprouter.Params)  
 	repo := EntrantRepo{c.db.C("entrants")}
  entrant, err := repo.Find(params.ByName("id")) //getting data from named param :id     
  body := context.Get(r, "body").(*EntrantResource)
  fmt.Println("Printing &body.Data")
  fmt.Println(&body.Data)
  body.Data.Id = entrant.Data.Id
  body.Data.Name = r.FormValue("Name")
  body.Data.Team_Id = r.FormValue("Team_Id")
  body.Data.Event_Id = r.Form["Event_Id"]
  fmt.Println("Printing &body.Data")
  fmt.Println(&body.Data)  
	err = repo.Update(&body.Data)
	if err != nil {
    fmt.Println("out of updateEntrantHandler")
		panic(err)
	}
	fmt.Println("out of updateEntrantHandler")
  http.Redirect(w, r, "/entrants/show/" + body.Data.Id.Hex(), 302)
}

func (c *appContext) deleteEntrantHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("In deleteEntrantHandler")
	params := context.Get(r, "params").(httprouter.Params)    //gorilla context, key "params"
	repo := EntrantRepo{c.db.C("entrants")}
	err := repo.Delete(params.ByName("id"))
	if err != nil {
    fmt.Println("out of deleteEntrantHandler")
		panic(err)
	}
	fmt.Println("out of deleteEntrantHandler")
  http.Redirect(w, r, "/entrants", 302)
}


// Scorecard Handlers /////////////////////////////////////////////////////////////////////////////////////

func (c *appContext) scorecardsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("In scorecardsHandler")
	repo := ScorecardRepo{c.db.C("scorecards")}
	scorecards, err := repo.All()
	if err != nil {
    fmt.Println("Out of scorecardsHandler")
		panic(err)
	}
  //	w.Header().Set("Content-Type", "application/vnd.api+json")
  //	json.NewEncoder(w).Encode(scorecards)
	// read BSON into JSON
  if err = listScorecard.Execute(w, scorecards.Data); err != nil {
      http.Error(w, err.Error(), http.StatusInternalServerError)
      fmt.Println("Out of scorecardsHandler")
      return
  }
  fmt.Println("Out of scorecardsHandler")
}

func (c *appContext) scorecardHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("In scorecardHandler")
  params := context.Get(r, "params").(httprouter.Params)  //gorrila context, key "params"
  repo := ScorecardRepo{c.db.C("scorecards")}
	scorecard, err := repo.Find(params.ByName("id")) //getting data from named param :id
  if err != nil {
    fmt.Println("out of scorecardHandler") 
		panic(err)
	}
  //  w.Header().Set("Content-Type", "application/vnd.api+json")
  //	json.NewEncoder(w).Encode(scorecard)  
  //	if err = show.Execute(w, json.NewEncoder(w).Encode(scorecard)); err != nil {
  //      http.Error(w, err.Error(), http.StatusInternalServerError)
  //      return
  //  }
	// read JSON into BSON 
	if err = showScorecard.Execute(w, scorecard.Data); err != nil {
      fmt.Println("out of scorecardHandler") 
      http.Error(w, err.Error(), http.StatusInternalServerError)
      return
  }  
  fmt.Println("out of scorecardHandler")  
}

func (c *appContext) newScorecardHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("In newScorecardHandler")
//  v := r.URL.Query()
//  fmt.Println("printing v") 
//  fmt.Println(v)   
  body := context.Get(r, "body").(*ScorecardResource)
  fmt.Println("printing &body.Data") 
  fmt.Println(&body.Data)
  repo := ScorecardRepo{c.db.C("scorecards")}
  err, id := repo.Create(&body.Data)
  fmt.Println(err)
  scorecard, err := repo.Find(id.Hex())
//  for i:=0; i<len(v["Team_Id"]); i++{
//    scorecard.Data.EntrantAll_Id = append(scorecard.Data.EntrantAll_Id, v["Team_Id"][i])
//  }

  fmt.Println("printing scorecard.Data")
  fmt.Println(scorecard.Data)
	if err := createnewScorecard.Execute(w, scorecard.Data); err != nil {
      fmt.Println("out of newScorecardHandler") 
      http.Error(w, err.Error(), http.StatusInternalServerError)
      return
  }
  fmt.Println("out of newScorecardHandler") 
}

func (c *appContext) editScorecardHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("In editScorecardHandler")
//  v := r.URL.Query()
//  fmt.Println("printing v") 
//  fmt.Println(v)  
  params := context.Get(r, "params").(httprouter.Params)    //gorilla context, key "params"
  Repo := ScorecardRepo{c.db.C("scorecards")}
 	scorecard, err := Repo.Find(params.ByName("id")) //getting data from named param :id  
//  fmt.Println(scorecard.Data)
//  for i:=0; i<len(v["Team_Id"]); i++{
//    fmt.Println("printing v[\"Team_Id\"][i]")
//    fmt.Println(v["Team_Id"][i])
//    scorecard.Data.EntrantAll_Id = append(scorecard.Data.EntrantAll_Id, v["Team_Id"][i])
//  }  
//  fmt.Println("printing scorecard.Data")
//  fmt.Println(scorecard.Data)
  if err = updateScorecard.Execute(w, scorecard.Data); err != nil {
      fmt.Println("out of editScorecardHandler")
      http.Error(w, err.Error(), http.StatusInternalServerError)
      return
  }
	fmt.Println("out of editScorecardHandler")
}

func (c *appContext) updateScorecardHandler(w http.ResponseWriter, r *http.Request) {
  fmt.Println("In updateScorecardHandler") 
	params := context.Get(r, "params").(httprouter.Params)  
 	Repo := ScorecardRepo{c.db.C("scorecards")}
  scorecard, err := Repo.Find(params.ByName("id")) //getting data from named param :id
  body := context.Get(r, "body").(*ScorecardResource)
  body.Data.Id = scorecard.Data.Id  
  body.Data.Element = r.FormValue("Element")
  body.Data.Maxtime_m = r.FormValue("Maxtime_m")  
  body.Data.Maxtime_s = r.FormValue("Maxtime_s")
  body.Data.Maxtime_ms = r.FormValue("Maxtime_ms")
  body.Data.Finish_call = r.FormValue("Finish_call")
  body.Data.False_alert_fringe = r.FormValue("False_alert_fringe")
  body.Data.Timed_out = r.FormValue("Timed_out")
  body.Data.Dismissed = r.FormValue("Dismissed")
  body.Data.Excused = r.FormValue("Excused")
  body.Data.Absent = r.FormValue("Absent")
  body.Data.Eliminated_during_search = r.FormValue("Eliminated_during_search")
  body.Data.Other_faults_descr = r.FormValue("Other_faults_descr")
  body.Data.Other_faults_count = r.FormValue("Other_faults_count")
  body.Data.Comments = r.FormValue("Comments")
  body.Data.Time_elapsed_m = r.FormValue("Time_elapsed_m")
  body.Data.Time_elapsed_s = r.FormValue("Time_elapsed_s")
  body.Data.Time_elapsed_ms = r.FormValue("Time_elapsed_ms")
  body.Data.Pronounced = r.FormValue("Pronounced")
  body.Data.Judge_signature = r.FormValue("Judge_signature")
  body.Data.Scorecard_Id = r.FormValue("Scorecard_Id")
  body.Data.Event_Id = r.FormValue("Event_Id")
  body.Data.Entrant_Id = r.FormValue("Entrant_Id")
  body.Data.Search_area = r.FormValue("Search_area")
  body.Data.Hides_max = r.FormValue("Hides_max")
  body.Data.Hides_found = r.FormValue("Hides_found")
  body.Data.Hides_missed = r.FormValue("Hides_missed")
  body.Data.Total_faults = r.FormValue("Total_faults")
  body.Data.Maxpoint = r.FormValue("Maxpoint")
  body.Data.Total_points = r.FormValue("Total_points")  
  
  err = Repo.Update(&body.Data)
  
  if err != nil {
    fmt.Println("out of updateScorecardHandler")
		panic(err)
	}
  //  w.WriteHeader(204)
  //	w.Write([]byte("\n"))	
	fmt.Println("out of updateScorecardHandler")
  if body.Data.Scorecard_Id == ""{
    http.Redirect(w, r, "/scorecards/delete/" + body.Data.Id.Hex(), 302)
  }else{
    http.Redirect(w, r, "/scorecards/show/" + body.Data.Id.Hex(), 302)
  }
}

func (c *appContext) get_elmSearchAreas(id string) string {
  screpo := ScorecardRepo{c.db.C("scorecards")}
  scorecard, err := screpo.Find(id)
  evrepo := EventRepo{c.db.C("events")}
	result := EventResource{} 
  err = evrepo.coll.Find(bson.M{"event_id": scorecard.Data.Event_Id}).One(&result.Data)  
  fmt.Println("printing result.Data")
  fmt.Println(result.Data)
  fmt.Println("printing err")
  fmt.Println(err)  
  elm := ""
  switch scorecard.Data.Element {
    case "Container":
      elm = result.Data.Cont_search_areas
    case "Interior":
      elm = result.Data.Int_search_areas
    case "Exterior":
      elm = result.Data.Ext_search_areas
    case "Vehicle":
      elm = result.Data.Veh_search_areas
    case "Elite":
      elm = result.Data.Elite_search_areas
  }
  fmt.Println("printing elm")
  fmt.Println(elm)
  return elm
}

func (c *appContext) get_elmHides(id string) string {
  screpo := ScorecardRepo{c.db.C("scorecards")}
  scorecard, err := screpo.Find(id)
  evrepo := EventRepo{c.db.C("events")}
	result := EventResource{} 
  err = evrepo.coll.Find(bson.M{"event_id": scorecard.Data.Event_Id}).One(&result.Data)  
  fmt.Println("printing result.Data")
  fmt.Println(result.Data)
  fmt.Println("printing err")
  fmt.Println(err)  
  elm := ""
  switch scorecard.Data.Element {
    case "Container":
      elm = result.Data.Cont_hides
    case "Interior":
      elm = result.Data.Int_hides
    case "Exterior":
      elm = result.Data.Ext_hides
    case "Vehicle":
      elm = result.Data.Veh_hides
    case "Elite":
      elm = result.Data.Elite_hides
  }
  fmt.Println("printing elm")
  fmt.Println(elm)
  return elm
}

func (c *appContext) deleteScorecardHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("In deleteScorecardHandler")
	params := context.Get(r, "params").(httprouter.Params)    //gorilla context, key "params"
	repo := ScorecardRepo{c.db.C("scorecards")}
	err := repo.Delete(params.ByName("id"))
	if err != nil {
    fmt.Println("out of deleteScorecardHandler")
		panic(err)
	}
  //	w.WriteHeader(204)
  //	w.Write([]byte("\n"))
	fmt.Println("out of deleteScorecardHandler")
  //  _, err = http.Get("/scorecards")
  http.Redirect(w, r, "/scorecards", 302)
}


// Router //////////////////////////////////////////////////////////////////////////////////////////

type router struct {
	*httprouter.Router
}

func (r *router) Get(path string, handler http.Handler) {
	r.GET(path, wrapHandler(handler))
	fmt.Println("Router getting")
}

func (r *router) Post(path string, handler http.Handler) {
	r.POST(path, wrapHandler(handler))
	fmt.Println("Router posting")
}

//func (r *router) Put(path string, handler http.Handler) {
//	r.PUT(path, wrapHandler(handler))
//	fmt.Println("Router putting")
//}

func (r *router) Delete(path string, handler http.Handler) {
	r.DELETE(path, wrapHandler(handler))
	fmt.Println("Router deleting")
}

// Integrating httprouter to our frameworks where it is incompatible with
// go http.Handler
// We wrap our middleware stack - implementing http.Handler into a 
// httprouter.Handler function

func NewRouter() *router {
  fmt.Println("NewRouter")
	return &router{httprouter.New()}
}

func wrapHandler(h http.Handler) httprouter.Handle {
	fmt.Println("In and out of wrapHandler")
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    context.Set(r, "params", ps)    //gorilla context, key "params"
		fmt.Println("printing value of \"params\"")
    fmt.Println(ps)
    h.ServeHTTP(w, r)
	}
}

// MAIN ////////////////////////////////////////////////////////////////////////////////////////////

func main() {
  
  session, err := mgo.Dial("localhost:27017")
	fmt.Println("Dialed for session")
	
  if err != nil {
		panic(err)
	}
	defer session.Close()
  
	session.SetMode(mgo.Monotonic, true)	
	
  appC := appContext{session.DB("test")}
    
  // commonHandlers := alice.New(context.ClearHandler, loggingHandler, recoverHandler, acceptHandler)
  commonHandlers := alice.New(context.ClearHandler, loggingHandler, recoverHandler)
	// alice is used to chain handlers
	// context from gorrila mapping
  fmt.Println("Chained handlers set up")
  
  router := NewRouter()

	fmt.Println("Called to NewRouter")
  
  router.ServeFiles("/static/*filepath", http.Dir("static"))
  
  //  Event routing  /////////////////
  
  router.Get("/events", commonHandlers.ThenFunc(appC.eventsHandler))
  router.Get("/events/show/:id", commonHandlers.ThenFunc(appC.eventHandler))  
  router.Get("/events/new", commonHandlers.Append(bodyHandler(EventResource{})).ThenFunc(appC.newEventHandler)) 
  //router.Post("/events/create", commonHandlers.Append(contentTypeHandler, bodyHandler(EventResource{})).ThenFunc(appC.createEventHandler))
  // router.Post("/events/create", commonHandlers.Append(bodyHandler(EventResource{})).ThenFunc(appC.createEventHandler))
  router.Get("/events/edit/:id/", commonHandlers.ThenFunc(appC.editEventHandler))
  router.Post("/events/update/:id/", commonHandlers.Append(bodyHandler(EventResource{})).ThenFunc(appC.updateEventHandler))  
  //  router.Put("/events/update/:id", commonHandlers.Append(contentTypeHandler, bodyHandler(EventResource{})).ThenFunc(appC.updateEventHandler))
	router.Get("/events/delete/:id", commonHandlers.ThenFunc(appC.deleteEventHandler))
	
  //  EventEntrants routing  /////////////////
  
  router.Get("/getEntrants/new", commonHandlers.ThenFunc(appC.newEventEntrantsHandler))
  router.Get("/getEntrants/edit/:id", commonHandlers.ThenFunc(appC.updateEventEntrantsHandler)) 
  
  //  Entrant routing  //////////////////
  
  router.Get("/entrants", commonHandlers.ThenFunc(appC.entrantsHandler))
//  router.Post("/entrants/", commonHandlers.ThenFunc(appC.updateEntrantsEventHandler))
  router.Get("/entrants/show/:id", commonHandlers.ThenFunc(appC.entrantHandler))
  router.Get("/entrants/new", commonHandlers.ThenFunc(newEntrantHandler))
  router.Post("/entrants/create", commonHandlers.Append(bodyHandler(EntrantResource{})).ThenFunc(appC.createEntrantHandler))  
  router.Get("/entrants/edit/:id", commonHandlers.ThenFunc(appC.editEntrantHandler))
  router.Post("/entrants/update/:id/", commonHandlers.Append(bodyHandler(EntrantResource{})).ThenFunc(appC.updateEntrantHandler))  
	router.Get("/entrants/delete/:id", commonHandlers.ThenFunc(appC.deleteEntrantHandler))  
 
  //  Scorecard routing  //////////////////
  
  router.Get("/scorecards", commonHandlers.ThenFunc(appC.scorecardsHandler))
  router.Get("/scorecards/show/:id", commonHandlers.ThenFunc(appC.scorecardHandler))
  router.Get("/scorecards/new", commonHandlers.Append(bodyHandler(ScorecardResource{})).ThenFunc(appC.newScorecardHandler))
  router.Get("/scorecards/edit/:id", commonHandlers.ThenFunc(appC.editScorecardHandler))
  router.Post("/scorecards/update/:id/", commonHandlers.Append(bodyHandler(ScorecardResource{})).ThenFunc(appC.updateScorecardHandler))  
	router.Get("/scorecards/delete/:id", commonHandlers.ThenFunc(appC.deleteScorecardHandler))  
 
  //  listening
  
  http.ListenAndServe(":8080", router)
}
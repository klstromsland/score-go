package main
// REST API to create, retrieve, update and delete scores
import (
//  "bufio"
//  "bytes"
	"encoding/json"  //implements encoding and decoding of JSON objects
  "fmt"
  "html/template"
//  "io"
//  "io/ioutil"  
  "log"
  "math"
	"net/http"
//	"net/url"  
//  "os"  
//  "path"  
	"reflect"
	// implements run-time reflection, allowing a program to manipulate 
	// objects with arbitrary types. The typical use is to take a value 
	// with static type interface{} and extract its dynamic type information 
	// by calling TypeOf, which returns a Type.    
  "regexp"
  "strconv"  
  "strings"
	"time" 
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

var listEntrant = template.Must(template.ParseFiles("templates/base.html", "templates/entrants/list/list.html"))
var createnewEntrant = template.Must(template.ParseFiles("templates/base.html", "templates/entrants/new/new.html", "templates/entrants/form.html"))
var updateEntrant = template.Must(template.ParseFiles("templates/base.html", "templates/entrants/update/update.html", "templates/entrants/form.html"))
var showEntrant = template.Must(template.ParseFiles("templates/base.html", "templates/entrants/show/show.html"))

var listScorecard = template.Must(template.ParseFiles("templates/base.html", "templates/scorecards/list/list.html"))
var createnewScorecard = template.Must(template.ParseFiles("templates/base.html", "templates/scorecards/new/new.html", "templates/scorecards/form.html"))
var updateScorecard = template.Must(template.ParseFiles("templates/base.html", "templates/scorecards/update/update.html", "templates/scorecards/form.html"))
var showScorecard = template.Must(template.ParseFiles("templates/base.html", "templates/scorecards/show/show.html"))

var ELEMENTS = []string{"Container", "Interior", "Exterior", "Vehicle", "Elite"}

var dataStart bool = false
var dataReset bool = false
var dataStop bool = false
var timeStart time.Time
var lastTime  time.Duration = 0
var timelimit  time.Duration
var timedata string = ""
var data string = ""
var count int = 0
const second = time.Second
const seconds = 10*second
const minute = time.Minute
const minutes = 0*time.Minute + 15*second
const millisecond = time.Millisecond
const milliseconds = 40*millisecond


//Event collection////////////////////////////////////////////////////////////////////////////////////

type Selected struct{
  Value string
  Selected bool
}

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
  EntrantAll_Id       []Selected    `json:"entrantall_id"`
  EntrantSelected_Id  []string      `json:"entrantselected_id"`
}

type EventsCollection struct {
	Data []Event `json:"data"`
}

type EventResource struct {
	Data Event `json:"data"`
}

type EventScorecardResource struct {
	 SCData ScorecardResource
   EVData EventResource   
}

type EventEntrantsResource struct {
	 ENcoll EntrantsCollection
   EVData EventResource   
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
  fmt.Println("result.Data")
  fmt.Println(result.Data)
  fmt.Println("event")
  fmt.Println(event)  
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
  fmt.Println("printing err")
  fmt.Println(err)
  fmt.Println("printing entrant")
  fmt.Println(entrant)
  fmt.Println("printing &result.Data")
  fmt.Println(&result.Data)
	if err != nil {
    fmt.Println("Find error")
    fmt.Println("out of Entrant Update")
		return err
	}  
  err = r.coll.Update(result.Data, entrant)
	if err != nil {
    fmt.Println("Update error")
    fmt.Println(err)
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
  timelimit = 0
	fmt.Println("In logging handler")
	fn := func(w http.ResponseWriter, r *http.Request) {
//		t1 := time.Now()
		next.ServeHTTP(w, r)
//		t2 := time.Now()
//		log.Printf("[%s] %q %v\n", r.Method, r.URL.String(), t2.Sub(t1))
//    log.Printf("[%s] %q \n", r.Method, r.URL.String())
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
//  v := r.URL.Query()
//  fmt.Println("printing v") 
//  fmt.Println(v)
  enRepo := EntrantRepo{c.db.C("entrants")}
  entrants, err := enRepo.All()
  body := context.Get(r, "body").(*EventResource)
  fmt.Println("printing &body.Data") 
  fmt.Println(&body.Data)
  evRepo := EventRepo{c.db.C("events")}
  err, id := evRepo.Create(&body.Data)
  fmt.Println(err)
  event, err := evRepo.Find(id.Hex())
//  for i:=0; i<len(v["Team_Id"]); i++{
//    newEntrant := Selected{Value: v["Team_Id"][i], Selected: false}
//    event.Data.EntrantAll_Id = append(event.Data.EntrantAll_Id, newEntrant)
// }
  for i:=0; i<len(entrants.Data); i++{
    newEntrant := Selected{Value: entrants.Data[i].Team_Id, Selected: false}
    event.Data.EntrantAll_Id = append(event.Data.EntrantAll_Id, newEntrant)
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
  params := context.Get(r, "params").(httprouter.Params)    //gorilla context, key "params"
  evRepo := EventRepo{c.db.C("events")}
 	event, err := evRepo.Find(params.ByName("id")) //getting data from named param :id  
  enRepo := EntrantRepo{c.db.C("entrants")}
  entrants, err := enRepo.All()  
  for i:=0; i<len(entrants.Data); i++{
    newEntrant := Selected{Value: entrants.Data[i].Team_Id, Selected: false}
    event.Data.EntrantAll_Id = append(event.Data.EntrantAll_Id, newEntrant)
    for i:=0; i<len(event.Data.EntrantAll_Id); i++{
      for j:=0; j<len(event.Data.EntrantSelected_Id); j++{
        if event.Data.EntrantAll_Id[i].Value == event.Data.EntrantSelected_Id[j]{
          event.Data.EntrantAll_Id[i].Selected = true
        }
      }
    }    
  }  
  fmt.Println("printing event.Data")
  fmt.Println(event.Data)
  // posts to /events/update/{{ .Id }} and the update event handler
  if err = updateEvent.Execute(w, event.Data); err != nil {
//  if err = updateEvent.Execute(os.Stdout, event.Data); err != nil {  
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
  enRepo := EntrantRepo{c.db.C("entrants")}
  entrants, err := enRepo.All()   
  for i:=0; i<len(entrants.Data); i++{
    newEntrant := Selected{Value: entrants.Data[i].Team_Id, Selected: false}
    event.Data.EntrantAll_Id = append(event.Data.EntrantAll_Id, newEntrant)
    for i:=0; i<len(event.Data.EntrantAll_Id); i++{
      for j:=0; j<len(event.Data.EntrantSelected_Id); j++{
        if event.Data.EntrantAll_Id[i].Value == event.Data.EntrantSelected_Id[j]{
          event.Data.EntrantAll_Id[i].Selected = true
        }
      }
    }    
  }  
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
  body.Data.EntrantAll_Id = event.Data.EntrantAll_Id
  body.Data.EntrantSelected_Id = r.Form["EntrantSelected_Id"]
  for i:=0; i<len(event.Data.EntrantAll_Id); i++{
    for j:=0; j<len(body.Data.EntrantSelected_Id); j++{
      if event.Data.EntrantAll_Id[i].Value == body.Data.EntrantSelected_Id[j]{
        body.Data.EntrantAll_Id[i].Selected = true
      }
    }
  }
  fmt.Println("&body.Data")
  fmt.Println(&body.Data)
  err = evRepo.Update(&body.Data)
  // Create or update scorecards for each entrant selected and add Event_Id to entrant
  for i:=0; i<len(body.Data.EntrantSelected_Id); i++{
  // func create scorecard for each event element with total hides = to event element hides for each entrant
    
  }
  
  if err != nil {
    fmt.Println("out of updateEventHandler")
		panic(err)
	}
  //  w.WriteHeader(204)
  //	w.Write([]byte("\n"))	
	fmt.Println("out of updateEventHandler")
  c.updateEntrantsEventHandler(w, r)
//  if body.Data.Event_Id == ""{
//    http.Redirect(w, r, "/events/delete/" + body.Data.Id.Hex(), 302)
//  }else{
//    http.Redirect(w, r, "/events/show/" + body.Data.Id.Hex(), 302)   
//  }
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
  //func delete event scorecards
  //func update entrant Event_Id
  //	w.WriteHeader(204)
  //	w.Write([]byte("\n"))
	fmt.Println("out of deleteEventHandler")
  //  _, err = http.Get("/events")
  http.Redirect(w, r, "/events", 302)
}


// EventEntrants Handlers /////////////////////////////////////////////////////////////////////////////////////


// gathers all of the entrants to send to new event handler (not needed?)
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

// gathers all of the entrants to send to edit event handler
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

// adds event_id to entrant and calls event entrant scorecard handler
func (c *appContext) updateEntrantsEventHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("In updateEntrantsEventHandler")
	enRepo := EntrantRepo{c.db.C("entrants")}
	entrants, err := enRepo.All() 
  fmt.Println("Printing entrants.Data")
  fmt.Println(entrants.Data)
  params := context.Get(r, "params").(httprouter.Params)  
 	evRepo := EventRepo{c.db.C("events")}
  event, err := evRepo.Find(params.ByName("id")) //getting data from named param :id 
  found := true
  for i:=0; i<len(event.Data.EntrantSelected_Id); i++{
    found = true
    for j:=0; j<len(entrants.Data); j++{          
      if entrants.Data[j].Team_Id == event.Data.EntrantSelected_Id[i]{
        if len(entrants.Data[j].Event_Id) == 0{
          found = false
        }else{
          for k:=0; k<len(entrants.Data[j].Event_Id); k++{
            if entrants.Data[j].Event_Id[k] != event.Data.Event_Id{
              found = false
            }
          }
        }
        if found == false{
          body := EntrantResource{}
          body.Data.Id = entrants.Data[j].Id
          body.Data.Event_Id = entrants.Data[j].Event_Id
          body.Data.Event_Id = append(body.Data.Event_Id, event.Data.Event_Id)
          body.Data.Team_Id = entrants.Data[j].Team_Id
          body.Data.Name = entrants.Data[j].Name
//          fmt.Println("printing body.Data")
//          fmt.Println(&body.Data)
          err = enRepo.Update(&body.Data)
        }
      }
    }
  }
  if err != nil {
    fmt.Println("Out of updateEntrantsEventHandler")
		panic(err)
	}
  fmt.Println("Out of updateEntrantsEventHandler")
  c.updateEventEntrantScorecardsHandler(w, r)
//  if event.Data.Event_Id == ""{
//    http.Redirect(w, r, "/events/delete/" + event.Data.Id.Hex(), 302)
//  }else{
//    http.Redirect(w, r, "/events/show/" + event.Data.Id.Hex(), 302)   
//  }
}



// EventEntrantScorecards Handlers /////////////////////////////////////////////////////////////////////////////////////

// creates scorecards
func (c *appContext) updateEventEntrantScorecardsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("In updateEventEntrantScorecards")
	enRepo := EntrantRepo{c.db.C("entrants")}
	entrants, err := enRepo.All() 
  fmt.Println("Printing entrants.Data")
  fmt.Println(entrants.Data)
  params := context.Get(r, "params").(httprouter.Params)  
 	evRepo := EventRepo{c.db.C("events")}
  event, err := evRepo.Find(params.ByName("id")) //getting data from named param :id 
//  scRepo := ScorecardRepo{c.db.C("scorecards")}
//  scorecard, err := scRepo.All() //getting data from named param :id    
  search_areas := 0
  scRepo := ScorecardRepo{c.db.C("scorecards")}
  for i:=0; i<len(event.Data.EntrantSelected_Id); i++{
    //create a scorecard for each event element
    
    for j:=0; j<len(ELEMENTS); j++{
      body := ScorecardResource{}    
      body.Data.Event_Id = event.Data.Event_Id
      body.Data.Entrant_Id = event.Data.EntrantSelected_Id[i]

      switch body.Data.Element = ELEMENTS[j]; ELEMENTS[j]{
        case "Container":
          if event.Data.Cont_search_areas != ""{
            search_areas, err = strconv.Atoi(event.Data.Cont_search_areas)
          }else{
            search_areas = 0
          }
        case "Interior":
          if (event.Data.Int_search_areas != ""){
            search_areas, err = strconv.Atoi(event.Data.Int_search_areas)
          }else{
            search_areas = 0
          }
        case "Exterior":
          if event.Data.Ext_search_areas != ""{        
            search_areas, err = strconv.Atoi(event.Data.Ext_search_areas)
          }else{
            search_areas = 0
          }
        case "Vehicle":
          if event.Data.Veh_search_areas != ""{                
            search_areas, err = strconv.Atoi(event.Data.Veh_search_areas)
          }else{
            search_areas = 0
          }
        case "Elite":
          if event.Data.Elite_search_areas != ""{                
            search_areas, err = strconv.Atoi(event.Data.Elite_search_areas)
          }else{
            search_areas = 0
          }
      }
      if search_areas > 0{
        for k:=1; k<=search_areas; k++{
          body.Data.Search_area = strconv.Itoa(k)
          body.Data.Scorecard_Id = event.Data.Event_Id + event.Data.EntrantSelected_Id[i] + body.Data.Search_area
          cbody := body
          err, id := scRepo.Create(&cbody.Data)
          fmt.Println(err)
          fmt.Println(id)
        }
      }
    }
  }
  if err != nil {
    fmt.Println("Out of updateEventEntrantScorecards")
		panic(err)
	}
  fmt.Println("Out of updateEventEntrantScorecards")    
  if event.Data.Event_Id == ""{
    http.Redirect(w, r, "/events/delete/" + event.Data.Id.Hex(), 302)
  }else{
    http.Redirect(w, r, "/events/show/" + event.Data.Id.Hex(), 302)   
  }
}


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
//	if err := createnewScorecard.Execute(w, scorecard.Data); err != nil {
//      fmt.Println("out of newScorecardHandler") 
//      http.Error(w, err.Error(), http.StatusInternalServerError)
//      return
//  }
  fmt.Println("out of newScorecardHandler")
  return
}

func (c *appContext) editScorecardHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("In editScorecardHandler")
  
 
  params := context.Get(r, "params").(httprouter.Params)    //gorilla context, key "params"
  repo := ScorecardRepo{c.db.C("scorecards")}
 	scorecard, err := repo.Find(params.ByName("id")) //getting data from named param :id  
  evRepo := EventRepo{c.db.C("events")}
  event := EventResource{}
  err = evRepo.coll.Find(bson.M{"event_id": scorecard.Data.Event_Id}).One(&event.Data)
  eventscorecard := EventScorecardResource{}
  eventscorecard.EVData = event
  eventscorecard.SCData = scorecard
 
  if err = updateScorecard.Execute(w, eventscorecard); err != nil {
      fmt.Println("out of editScorecardHandler")
      http.Error(w, err.Error(), http.StatusInternalServerError)
      return
  }
  
  http.Redirect(w, r, "/scorecards/edit/" + scorecard.Data.Id.Hex(), 302)
	fmt.Println("out of editScorecardHandler")
  if err != nil{
    fmt.Println("printing err")
    fmt.Println(err)
  }
}


// handler to cater AJAX requests
func handlerGetTime(w http.ResponseWriter, r *http.Request) {
  var idata int
// gets time since an instance of time was declared
  var newTime time.Duration
// gives the difference between the "last" newTime and the current newTime  
  var diff time.Duration
// if we just started to get time go for it
// do this unless JQuery shuts us down through a click on stop or the timelimit
// has been superceded
  if dataStart && !dataStop{
    newTime = time.Since(timeStart)
    // the elapsed time since last call from JQuery
    diff = newTime - lastTime
    if diff < 0{
      diff = -diff
    }
    // if elapsed time (diff) is less than the approximate delay from JQuery, then continue
    // and send processed data to handlerPostTime
    // otherwise we may have been stopped in some way and should respond
    // if the time since timeStart (newTime) is less than the timelimit, we should keep
    // going
    if diff <= milliseconds{
      // just milliseconds to deal with
      if newTime < second{
        data = newTime.String()
        re := regexp.MustCompile("ms")
        data = re.ReplaceAllString(data, "")
        fdata, err := strconv.ParseFloat(data, 64)
        fdata = fdata/10     
        fldata := math.Floor(fdata)
        idata = int(fldata)
        if idata < 10{
          data = "00:00:" + "0" + strconv.Itoa(idata)
        }else{
          data = "00:00:" + strconv.Itoa(idata)
        }
        if err != nil{
          fmt.Println(err)
        }       
        http.Redirect(w, r, "/savetime/" + data, 302)
      }
      // process seconds
      if newTime >= second && newTime < minute{
        data = newTime.String()       
        re := regexp.MustCompile("s")
        data = re.ReplaceAllString(data, "")    
        sdata := strings.SplitN(data, ".", 2)
        sdata1 := sdata[0]
        sdata2 := sdata[1][0:3]
        fdata2, err := strconv.ParseFloat(sdata2, 64)     
        fdata2 = fdata2/10
        fldata := math.Floor(fdata2)
        idata = int(fldata)
        isdata, err := strconv.Atoi(sdata1)      
        if isdata < 10{
          if idata < 10{
            data = "00:" + "0" + sdata1 + ":0" + strconv.Itoa(idata)
          }else{
            data = "00:" + "0" + sdata1 + ":" + strconv.Itoa(idata)
          }
        }else{
          if idata < 10{
            data = "00:" + sdata1 + ":0" + strconv.Itoa(idata)
          }else{
            data = "00:" + sdata1 + ":" + strconv.Itoa(idata)
          }          
        }
        if err != nil{
          fmt.Println(err)    
        }       
        http.Redirect(w, r, "/savetime/" + data, 302)
      }
      // process minutes and seconds
      if newTime >= minute{
        data = newTime.String()      
        sre := regexp.MustCompile("s")
        data = sre.ReplaceAllString(data, "")
        mdata := strings.SplitN(data, "m", 2)
        sdata := strings.SplitN(mdata[1], ".", 2)
        sdata1 := sdata[0]
        sdata2 := sdata[1][0:3]
        fdata2, err := strconv.ParseFloat(sdata2, 64)      
        fdata2 = fdata2/10
        fldata := math.Floor(fdata2)
        idata = int(fldata)
        imdata, err := strconv.Atoi(mdata[0])
        isdata, err := strconv.Atoi(sdata1)    
        if imdata < 10{
          if isdata < 10{
            if idata < 10{
              data = "0" + mdata[0] + ":" + "0" + sdata1 + ":0" + strconv.Itoa(idata)
            }else{
              data = "0" + mdata[0] + ":" + "0" + sdata1 + ":" + strconv.Itoa(idata)
            }            
          }else{
            if idata < 10{
              data = "0" + mdata[0] + ":" + sdata1 + ":0" + strconv.Itoa(idata)
            }else{
              data = "0" + mdata[0] + ":" + sdata1 + ":" + strconv.Itoa(idata)
            }
          }
        }else{
          if isdata < 10{
            if idata < 10{
              data = mdata[0] + ":" + "0" + sdata1 + ":0" + strconv.Itoa(idata)
            }else{
              data = mdata[0] + ":" + "0" + sdata1 + ":" + strconv.Itoa(idata)
            }            
          }else{
            if idata < 10{
              data = mdata[0] + ":" + sdata1 + ":0" + strconv.Itoa(idata)
            }else{
              data = mdata[0] + ":" + sdata1 + ":" + strconv.Itoa(idata)
            }
          }
        }
        if err != nil{
          fmt.Println(err)
        }      
        http.Redirect(w, r, "/savetime/" + data, 302)
      }
    }
    // if newTime is greater than timelimit, we need to go to different logic so that output repeats
    // until user figures it out, otherwise will continue by providing lastTime for next diff check
    if newTime >= timelimit{
      dataReset = true
    }
    lastTime = newTime   
  // if we just started to get time go for it, unless a stop was called
  // a new instant of time.Now is instantiated to serve as a reference for elapsed time
  // in order to process duration data
  }
  if !dataStart && !dataStop{
    timeStart = time.Now()
    dataStart = true
    lastTime = 0
  // we have been stopped by JQuery (diff is too large) or have received a timelimit signal
  // dataStop is true
  }
  if dataStop{
    http.Redirect(w, r, "/savetime/" + data, 302)
    newTime = time.Since(timeStart)
    // the elapsed time since last call from JQuery
    diff = newTime - lastTime
    if diff < 0{
      diff = -diff
    }
    lastTime = newTime    
  }
  if diff > milliseconds{
    http.Redirect(w, r, "/savetime/" + data, 302)
    dataStop = false
    dataStart = false
  }
  if dataReset{
    dataStop = true
    dataReset = false
  }
}

// handler to cater AJAX requests
func handlerPostTime(w http.ResponseWriter, r *http.Request) {
  params := context.Get(r, "params").(httprouter.Params)
  fmt.Fprint(w, params.ByName("data"))
  timedata = params.ByName("data")       
}

func (c *appContext) updateScorecardHandler(w http.ResponseWriter, r *http.Request) {
  fmt.Println("In updateScorecardHandler") 
	params := context.Get(r, "params").(httprouter.Params)  
 	repo := ScorecardRepo{c.db.C("scorecards")}
  scorecard, err := repo.Find(params.ByName("id")) //getting data from named param :id
  body := context.Get(r, "body").(*ScorecardResource)
  body.Data.Id = scorecard.Data.Id  
  body.Data.Element = scorecard.Data.Element
  body.Data.Maxtime_m = r.FormValue("Maxtime_m")  
  body.Data.Maxtime_s = r.FormValue("Maxtime_s")
  body.Data.Maxtime_ms = "00"
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
  body.Data.Time_elapsed_m = timedata
//  body.Data.Time_elapsed_s = r.FormValue("Time_elapsed_s")
//  body.Data.Time_elapsed_ms = r.FormValue("Time_elapsed_ms")
  body.Data.Pronounced = r.FormValue("Pronounced")
  body.Data.Judge_signature = r.FormValue("Judge_signature")
  body.Data.Event_Id = scorecard.Data.Event_Id
  body.Data.Entrant_Id = scorecard.Data.Entrant_Id
  body.Data.Search_area = scorecard.Data.Search_area
  body.Data.Scorecard_Id = scorecard.Data.Scorecard_Id
  body.Data.Hides_max = r.FormValue("Hides_max")
  body.Data.Hides_found = r.FormValue("Hides_found")
  body.Data.Hides_missed = r.FormValue("Hides_missed")
  body.Data.Total_faults = r.FormValue("Total_faults")
  body.Data.Maxpoint = r.FormValue("Maxpoint")
  body.Data.Total_points = r.FormValue("Total_points")  

  
  tmp_time_m, err := strconv.Atoi(body.Data.Maxtime_m)
  tmp_time_m = tmp_time_m*60
  tmp_time_s, err := strconv.Atoi(body.Data.Maxtime_s)
  tmp_time := tmp_time_s + tmp_time_m
  tmp_timeD := time.Duration(tmp_time)*time.Second
  fmt.Println("printing tmp seconds")
  timelimit = tmp_timeD
  fmt.Println(tmp_timeD)
  
  err = repo.Update(&body.Data)
	
  fmt.Println("in hideCountCheck")
  
 	evRepo := EventRepo{c.db.C("events")}

  repo = ScorecardRepo{c.db.C("scorecards")}
  scorecards, err := repo.All()

  fmt.Println("printing body.Data.Event_Id")
  fmt.Println(body.Data.Event_Id)
	
  event := EventResource{}  
  err = evRepo.coll.Find(bson.M{"event_id": body.Data.Event_Id}).One(&event.Data)   

  fmt.Println("printing event.Data")
  fmt.Println(event.Data)
  
  elm_hides := 0
  hideCountCheck := 0  
  hides_max := 0
    
  for j:=0; j<len(ELEMENTS); j++{
    switch body.Data.Element = ELEMENTS[j]; ELEMENTS[j]{
      case "Container":
        if event.Data.Cont_search_areas != ""{
          hideCountCheck, err = strconv.Atoi(event.Data.Cont_hides)
        }else{
          hideCountCheck = 0
        }
      case "Interior":
        if (event.Data.Int_search_areas != ""){
          hideCountCheck, err = strconv.Atoi(event.Data.Int_hides)
        }else{
          hideCountCheck = 0
        }
      case "Exterior":
        if event.Data.Ext_search_areas != ""{        
          hideCountCheck, err = strconv.Atoi(event.Data.Ext_hides)
        }else{
          hideCountCheck = 0
        }
      case "Vehicle":
        if event.Data.Veh_search_areas != ""{                
          hideCountCheck, err = strconv.Atoi(event.Data.Veh_hides)
        }else{
          hideCountCheck = 0
        }
      case "Elite":
        if event.Data.Elite_search_areas != ""{                
          hideCountCheck, err = strconv.Atoi(event.Data.Elite_hides)
        }else{
          hideCountCheck = 0
        }
    }
    elm_hides = hideCountCheck
    if body.Data.Hides_max != ""{
      for i:=0; i<len(scorecards.Data);i++{
        if body.Data.Entrant_Id == scorecards.Data[i].Entrant_Id{
          if body.Data.Event_Id == scorecards.Data[i].Event_Id{
            if body.Data.Element == scorecards.Data[i].Element && (body.Data.Hides_max != "" || body.Data.Hides_max != "0"){
              hides_max, err = strconv.Atoi(body.Data.Hides_max)
              hideCountCheck -= hides_max
            }
          }
        }
      }
      hides_max, err = strconv.Atoi(body.Data.Hides_max)
      if (hideCountCheck > elm_hides || hideCountCheck < 0 ) && (event.Data.Division != "NW1"){
        fmt.Println("Incorrect Hide Count...")
      }else if event.Data.Division == "NW1"{
        if hides_max != 1{
          fmt.Println("Incorrect Hide Count...")
        }
      }else if event.Data.Division == "NW2"{
        if hides_max == 0{
          fmt.Println("Incorrect Hide Count...")
        }
      }
    }    
  }

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
    http.Redirect(w, r, "/scorecards/edit/" + body.Data.Id.Hex(), 302)
  }
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
//		fmt.Println("printing value of \"params\"")
//    fmt.Println(ps)
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
  router.Get("/scorecards/edit/:id", commonHandlers.Append(bodyHandler(ScorecardResource{})).ThenFunc(appC.editScorecardHandler))
//  router.Post("/checkhidecount", commonHandlers.ThenFunc(checkHideCountHandler))
//  router.Get("/savecheckhidecount/:data", commonHandlers.ThenFunc(postHideCountCheckHandler))
  router.Post("/gettime", commonHandlers.ThenFunc(handlerGetTime))
  router.Get("/savetime/:data", commonHandlers.ThenFunc(handlerPostTime))
  router.Post("/scorecards/update/:id/", commonHandlers.Append(bodyHandler(ScorecardResource{})).ThenFunc(appC.updateScorecardHandler))  
	router.Get("/scorecards/delete/:id", commonHandlers.ThenFunc(appC.deleteScorecardHandler))  
 
  //  listening
  
  http.ListenAndServe(":8080", router)
}
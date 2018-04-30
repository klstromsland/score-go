package main
// REST API to create, retrieve, update and delete scores
import (
  "encoding/json"  //implements encoding and decoding of JSON objects
  "fmt"
  "html/template"
  "log"
  "math"
  "math/rand"
  "net/http"
  // "os"
  "reflect"
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
// Based on blog by Nicolas Merouze
// Repo BSON spec, mongodb driver

var listEvent = template.Must(template.ParseFiles("templates/base.html", "templates/events/list/list.html"))
var createnewEvent = template.Must(template.ParseFiles("templates/base.html", "templates/events/new/new.html", "templates/events/form.html"))
var updateEvent = template.Must(template.ParseFiles("templates/base.html", "templates/events/form.html", "templates/events/update/update.html"))
var showEvent = template.Must(template.ParseFiles("templates/base.html", "templates/events/show/show.html"))

var showInfo = template.Must(template.ParseFiles("templates/base.html", "templates/info/info.html"))

var listEntrant = template.Must(template.ParseFiles("templates/base.html", "templates/entrants/list/list.html"))
var createnewEntrant = template.Must(template.ParseFiles("templates/base.html", "templates/entrants/new/new.html", "templates/entrants/form.html"))
var updateEntrant = template.Must(template.ParseFiles("templates/base.html", "templates/entrants/update/update.html", "templates/entrants/form.html"))
var showEntrant = template.Must(template.ParseFiles("templates/base.html", "templates/entrants/show/show.html"))

var listUser = template.Must(template.ParseFiles("templates/base.html", "templates/users/list/list.html"))
var createnewUser = template.Must(template.ParseFiles("templates/base.html", "templates/users/new/new.html", "templates/users/form.html"))
var updateUser = template.Must(template.ParseFiles("templates/base.html", "templates/users/update/update.html", "templates/users/form.html"))
var showUser = template.Must(template.ParseFiles("templates/base.html", "templates/users/show/show.html"))

var listScorecard = template.Must(template.ParseFiles("templates/base.html", "templates/scorecards/list/list.html"))
var createnewScorecard = template.Must(template.ParseFiles("templates/base.html", "templates/scorecards/new/new.html", "templates/scorecards/form.html"))
var updateScorecard = template.Must(template.ParseFiles("templates/base.html", "templates/scorecards/update/update.html", "templates/scorecards/form.html"))
var showScorecard = template.Must(template.ParseFiles("templates/base.html", "templates/scorecards/show/show.html"))

var listTally = template.Must(template.ParseFiles("templates/base.html", "templates/tallies/list/list.html"))
var updateTally = template.Must(template.ParseFiles("templates/base.html", "templates/tallies/update/update.html", "templates/tallies/form.html"))
var showTally = template.Must(template.ParseFiles("templates/base.html", "templates/tallies/show/show.html"))

var createnewSession = template.Must(template.ParseFiles("templates/base.html",  "templates/sessions/form.html", "templates/sessions/new/new.html"))

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
var repeat int
const second = time.Second
const minute = time.Minute
const millisecond = time.Millisecond
const jqDelay = 10000*millisecond


//Global structs//////////////////////////////////////////////////////////////////////////////////////


type Selected struct{
  Value string
  Selected bool
}

type Session struct{
  Current_user string
  Current_email string
  Current_status bool
}


type SessionResource struct{
  SData Session
}

 var current_session = Session{Current_user: "", Current_email: "", Current_status: false}

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
  w.Header().Set("Content-Type", "application/vnd.api+json")
  w.WriteHeader(err.Status)
  json.NewEncoder(w).Encode(Errors{[]*Error{err}})
}

var (
  ErrBadRequest           = &Error{"bad_request", 400, "Bad request", "Request body is not well-formed. It must be JSON."}
  ErrNotAcceptable        = &Error{"not_acceptable", 406, "Not Acceptable", "Accept header must be set to 'application/vnd.api+json'."}
  ErrUnsupportedMediaType = &Error{"unsupported_media_type", 415, "Unsupported Media Type", "Content-Type header must be set to: 'application/vnd.api+json'."}
  ErrInternalServer       = &Error{"internal_server_error", 500, "Internal Server Error", "Something went wrong."}
)


//Event collection////////////////////////////////////////////////////////////////////////////////////

type Event struct {
  Id                  bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
  Name                string        `json:"name"`
  Location            string        `json:"location"`
  Date                string        `json:"data"`
  Host                string        `json:"host"`
  Status              string        `json:"status"`
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
  UserAll_Id          []Selected    `json:"userall_id"`
  UserSelected_Id     []string      `json:"userselected_id"`
}

type EventsCollection struct {
  Data []Event `json:"data"`
}

type EventResource struct {
  Data Event `json:"data"`
  SData Session
}

type EventsResource struct {
  Data []EventResource `json:"data"`
  SData Session
}

type EventShowResource struct {
  EVData EventResource
  ENSData EntrantsResource
  USRSData UsersResource
  TLYSData TalliesResource
  SCSData ScorecardsResource
  SCScompleted  []string
  TLYcompleted  []string
  Rank          []string
  SData Session
}

type EventEditResource struct {
  EVData EventResource
  ENSData EntrantsResource
  USRSData UsersResource
  SData Session
}

type EventRepo struct {
  coll *mgo.Collection
}

func (r *EventRepo) All() (EventsCollection, error) {
  result := EventsCollection{[]Event{}}
  err := r.coll.Find(nil).All(&result.Data)
  if err != nil {
	return result, err
  }
  return result, nil
}

func (r *EventRepo) Find(id string) (EventResource, error) {
  prefix := "ObjectIdHex(\""
  suffix := "\")"
  id = strings.TrimPrefix(id, prefix)
  id = strings.TrimSuffix(id, suffix)
  idb := bson.ObjectIdHex(id)
  result := EventResource{}
  err := r.coll.Find(bson.M{"_id": idb}).One(&result.Data)
  if err != nil {
	return result, err
  }
  return result, nil
}

func (r *EventRepo) Create(event *Event) (error, bson.ObjectId) {
  id := bson.NewObjectId()
  _, err := r.coll.UpsertId(id, event)
  if err != nil {
    return err, id
  }
  event.Id = id
  return err, id
}

func (r *EventRepo) Update(event *Event) error {
  result := EventResource{}
  err := r.coll.Find(bson.M{"_id": event.Id}).One(&result.Data)
  if err != nil {
	return err
  }
  err = r.coll.Update(result.Data, event)
  if err != nil {
	return err
  }
  return nil
}

func (r *EventRepo) Delete(id string) error {
  prefix := "ObjectIdHex(\""
  suffix := "\")"
  id = strings.TrimPrefix(id, prefix)
  id = strings.TrimSuffix(id, suffix)
  idb := bson.ObjectIdHex(id)
  result := EventResource{}
  err := r.coll.Find(bson.M{"_id": idb}).One(&result.Data)
  if err != nil {
	return err
  }
  err = r.coll.Remove(result.Data)
  if err != nil {
	return err
  }
  return nil
}

// Entrant collection  //////////////////////////////////////////////////////////////////////////////

type Entrant struct {
  Id bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
  First_name    string        `json:"first_name"`
  Last_name     string        `json:"last_name"`
  Id_number     string        `json:"id_number"`
  Dog_name      string        `json:"dog_name"`
  Dog_id_number string        `json:"dog_id_number"`
  Breed         string        `json:"breed"`
  Team_Id       string        `json:"entrant_id"`
  Event_Id      []string      `json:"event_id"`
}

type EntrantsCollection struct {
  Data []Entrant `json:"data"`
}

type EntrantResource struct {
  Data Entrant `json:"data"`
  SData Session
}

type EntrantsResource struct {
  Data []EntrantResource `json:"data"`
  SData Session
}

type EntrantRepo struct {
  coll *mgo.Collection
}

func (r *EntrantRepo) All() (EntrantsCollection, error) {
  result := EntrantsCollection{[]Entrant{}}
  err := r.coll.Find(nil).All(&result.Data)
  if err != nil {
	return result, err
  }
  return result, nil
}

func (r *EntrantRepo) Find(id string) (EntrantResource, error) {
  prefix := "ObjectIdHex(\""
  suffix := "\")"
  id = strings.TrimPrefix(id, prefix)
  id = strings.TrimSuffix(id, suffix)
  idb := bson.ObjectIdHex(id)
  result := EntrantResource{}
  err := r.coll.Find(bson.M{"_id": idb}).One(&result.Data)
  if err != nil {
	return result, err
  }
  return result, nil
}

func (r *EntrantRepo) Create(entrant *Entrant) (error, bson.ObjectId) {
  id := bson.NewObjectId()
  _, err := r.coll.UpsertId(id, entrant)
  if err != nil {
      return err, id
  }
  entrant.Id = id
  return err, id
}

func (r *EntrantRepo) Update(entrant *Entrant) error {
  result := EntrantResource{}
  err := r.coll.Find(bson.M{"_id": entrant.Id}).One(&result.Data)
  if err != nil {
	return err
  }
  err = r.coll.Update(result.Data, entrant)
  if err != nil {
	return err
  }
  return nil
}

func (r *EntrantRepo) Delete(id string) error {
  prefix := "ObjectIdHex(\""
  suffix := "\")"
  id = strings.TrimPrefix(id, prefix)
  id = strings.TrimSuffix(id, suffix)
  idb := bson.ObjectIdHex(id)
  result := EntrantResource{}
  err := r.coll.Find(bson.M{"_id": idb}).One(&result.Data)
  if err != nil {
	return err
  }
  err = r.coll.Remove(result.Data)
  if err != nil {
	return err
  }
  return nil
}

//User collection////////////////////////////////////////////////////////////////////////////////////

type User struct {
  Id              bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
  First_name      string        `json:"first_name_name"`
  Last_name       string        `json:"last_name_name"`
  Role            string        `json:"role"`
  Approved        string        `json:"approved"`
  Status          string        `json:"status"`
  Email           string        `json:"email"`
  Event_Id        []string      `json:"event_id"`
  User_Id         string        `json:"user_id"`
  Password        string        `json:"password"`
}

type UsersCollection struct {
  Data []User `json:"data"`
}

type UserResource struct {
  Data User `json:"data"`
  SData Session
}

type UsersResource struct {
  Data []UserResource `json:"data"`
  SData Session
}

type UserRepo struct {
  coll *mgo.Collection
}

func (r *UserRepo) All() (UsersCollection, error) {
  result := UsersCollection{[]User{}}
  err := r.coll.Find(nil).All(&result.Data)
  if err != nil {
	return result, err
  }
  return result, nil
}

func (r *UserRepo) Find(id string) (UserResource, error) {
  prefix := "ObjectIdHex(\""
  suffix := "\")"
  id = strings.TrimPrefix(id, prefix)
  id = strings.TrimSuffix(id, suffix)
  idb := bson.ObjectIdHex(id)
  result := UserResource{}
  err := r.coll.Find(bson.M{"_id": idb}).One(&result.Data)
  if err != nil {
	return result, err
  }
  return result, nil
}

func (r *UserRepo) Create(user *User) (error, bson.ObjectId) {
  id := bson.NewObjectId()
  _, err := r.coll.UpsertId(id, user)
  if err != nil {
    return err, id
  }
  user.Id = id
  return err, id
}

func (r *UserRepo) Update(user *User) error {
  result := UserResource{}
  err := r.coll.Find(bson.M{"_id": user.Id}).One(&result.Data)
  if err != nil {
	return err
  }
  err = r.coll.Update(result.Data, user)
  if err != nil {
	return err
  }
  return nil
}

func (r *UserRepo) Delete(id string) error {
  prefix := "ObjectIdHex(\""
  suffix := "\")"
  id = strings.TrimPrefix(id, prefix)
  id = strings.TrimSuffix(id, suffix)
  idb := bson.ObjectIdHex(id)
  result := UserResource{}
  err := r.coll.Find(bson.M{"_id": idb}).One(&result.Data)
  if err != nil {
	return err
  }
  err = r.coll.Remove(result.Data)
  if err != nil {
    return err
  }
  return nil
}


//Scorecard collection////////////////////////////////////////////////////////////////////////////////////


type Scorecard struct {
  Id bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
  Element                   string        `json:"element"`
  Maxtime_m                 string        `json:"maxtime_m"`
  Maxtime_s                 string        `json:"maxtime_s"`
  Finish_call               Selected      `json:"finish_call"`
  False_alert_fringe        string        `json:"false_alert_fringe"`
  Timed_out                 Selected      `json:"timed_out"`
  Dismissed                 Selected      `json:"dismissed"`
  Excused                   Selected      `json:"excused"`
  Absent                    Selected      `json:"absent"`
  Eliminated_during_search  Selected      `json:"eliminated_during_search"`
  Other_faults_descr        string        `json:"other_faults_descr"`
  Other_faults_count        string        `json:"other_faults_count"`
  Comments                  string        `json:"comments"`
  Total_time                string        `json:"total_time"`
  Pronounced                Selected      `json:"pronounced"`
  Judge_signature           Selected      `json:"judge_signature"`
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
  Data Scorecard `json:"data"`
  SData Session
}

type ScorecardsResource struct {
  Data []ScorecardResource `json:"data"`
  SData Session
}

type ScorecardFormResource struct {
   SCData ScorecardResource
   EVData EventResource
   ENData EntrantResource
   CheckCount string
   SData Session
}

type ScorecardRepo struct {
	coll *mgo.Collection
}

func (r *ScorecardRepo) All() (ScorecardsCollection, error) {
  result := ScorecardsCollection{[]Scorecard{}}
  err := r.coll.Find(nil).All(&result.Data)
  if err != nil {
    return result, err
  }
  return result, nil
}

func (r *ScorecardRepo) Find(id string) (ScorecardResource, error) {
  prefix := "ObjectIdHex(\""
  suffix := "\")"
  id = strings.TrimPrefix(id, prefix)
  id = strings.TrimSuffix(id, suffix)
  idb := bson.ObjectIdHex(id)
  result := ScorecardResource{}
  err := r.coll.Find(bson.M{"_id": idb}).One(&result.Data)
  if err != nil {
	return result, err
  }
  return result, nil
}

func (r *ScorecardRepo) Create(scorecard *Scorecard) (error, bson.ObjectId) {
  id := bson.NewObjectId()
  _, err := r.coll.UpsertId(id, scorecard)
  if err != nil {
    return err, id
  }
  scorecard.Id = id
  return err, id
}

func (r *ScorecardRepo) Update(scorecard *Scorecard) error {
  result := ScorecardResource{}
  err := r.coll.Find(bson.M{"_id": scorecard.Id}).One(&result.Data)
  if err != nil {
	return err
  }
  err = r.coll.Update(result.Data, scorecard)
  if err != nil {
	return err
  }
  return nil
}

func (r *ScorecardRepo) Delete(id string) error {
  prefix := "ObjectIdHex(\""
  suffix := "\")"
  id = strings.TrimPrefix(id, prefix)
  id = strings.TrimSuffix(id, suffix)
  idb := bson.ObjectIdHex(id)
  result := ScorecardResource{}
  err := r.coll.Find(bson.M{"_id": idb}).One(&result.Data)
  if err != nil {
    return err
  }
  err = r.coll.Remove(result.Data)
  if err != nil {
	return err
  }
  return nil
}


//Tally collection////////////////////////////////////////////////////////////////////////////////////


type Tally struct {
  Id                        bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
  Event_Id                  string        `json:"event_id"`
  Entrant_Id                string        `json:"entrant_id"`
  Tally_Id                  string        `json:"tally_id"`
  Total_time                string        `json:"total_time"`
  Total_faults              string        `json:"total_faults"`
  Total_points              string        `json:"total_points"`
  Title                     string        `json:"title"`
  Qualifying_score          string        `json:"qualifying_score"`
  Qualifying_scores         string        `json:"qualifying_scores"`
}

type TalliesCollection struct {
   Data []Tally `json:"data"`
}

type TallyResource struct {
  Data Tally `json:"data"`
  SData Session
}

type TalliesResource struct {
  Data []TallyResource `json:"data"`
  SData Session
}

type TallyFormResource struct {
   TLYData TallyResource
   EVData EventResource
   ENData EntrantResource
   SCSData ScorecardsResource
   SData Session
}

type TallyRepo struct {
  coll *mgo.Collection
}

func (r *TallyRepo) All() (TalliesCollection, error) {
	result := TalliesCollection{[]Tally{}}
	err := r.coll.Find(nil).All(&result.Data)
	if err != nil {
	  return result, err
	}
    return result, nil
}

func (r *TallyRepo) Find(id string) (TallyResource, error) {
  prefix := "ObjectIdHex(\""
  suffix := "\")"
  id = strings.TrimPrefix(id, prefix)
  id = strings.TrimSuffix(id, suffix)
  idb := bson.ObjectIdHex(id)
  result := TallyResource{}
  err := r.coll.Find(bson.M{"_id": idb}).One(&result.Data)
  if err != nil {
	return result, err
  }
  return result, nil
}

func (r *TallyRepo) Create(tally *Tally) (error, bson.ObjectId) {
	id := bson.NewObjectId()
	_, err := r.coll.UpsertId(id, tally)
	if err != nil {
      return err, id
	}
	tally.Id = id
	return err, id
}

func (r *TallyRepo) Update(tally *Tally) error {
	result := TallyResource{}
    err := r.coll.Find(bson.M{"_id": tally.Id}).One(&result.Data)
	if err != nil {
		return err
	}
    err = r.coll.Update(result.Data, tally)
	if err != nil {
		return err
	}
	return nil
}

func (r *TallyRepo) Delete(id string) error {
  prefix := "ObjectIdHex(\""
  suffix := "\")"
  id = strings.TrimPrefix(id, prefix)
  id = strings.TrimSuffix(id, suffix)
  idb := bson.ObjectIdHex(id)
  result := TallyResource{}
  err := r.coll.Find(bson.M{"_id": idb}).One(&result.Data)
  if err != nil {
    return err
  }
  err = r.coll.Remove(result.Data)
  if err != nil {
    return err
  }
  return nil
}


// Middlewares//////////////////////////////////////////////////////////////////////////////////////////
// go net/http

func recoverHandler(next http.Handler) http.Handler {
  fn := func(w http.ResponseWriter, r *http.Request) {
    defer func() {
      if err := recover(); err != nil {
        log.Printf("panic: %+v", err)
		WriteError(w, ErrInternalServer)
      }
     }()
	 next.ServeHTTP(w, r)
   }
   return http.HandlerFunc(fn)
}

//  go net/http
func loggingHandler(next http.Handler) http.Handler {
  timelimit = 0
  fn := func(w http.ResponseWriter, r *http.Request) {
//		t1 := time.Now()
    next.ServeHTTP(w, r)
//		t2 := time.Now()
   }
   return http.HandlerFunc(fn)
}

//  go net/http
func acceptHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
        if r.Header.Get("Accept") != "application/vnd.api+json" {
            WriteError(w, ErrNotAcceptable)
    		return
    	}
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

//  go net/http
func contentTypeHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
        if r.Header.Get("Content-Type") != "application/vnd.api+json" {
            WriteError(w, ErrUnsupportedMediaType)
    		return
        }
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

//  go net/http, reflect, gorilla context, mongodb driver
func bodyHandler(v interface{}) func(next http.Handler) http.Handler {
	t := reflect.TypeOf(v)                        //type interface{} which may be empty
	m := func(next http.Handler) http.Handler {
          fn := func(w http.ResponseWriter, r *http.Request) {
               val := reflect.New(t).Interface()   //val is type interface{}
               // err := json.NewDecoder(r.Body).Decode(val)  //r.Body is the request body and is type interface io.ReadCloser
               // // err := json.NewDecoder(strings.NewReader(evj)).Decode(val)
               // // val = evj
               // if err != nil {
			// 	WriteError(w, ErrBadRequest)
			// 	return
			// }
          	if next != nil {
          		context.Set(r, "body", val)     //gorilla context, key "body": val, val is type interface{}  "body" will now retrieve val
                 next.ServeHTTP(w, r)
          	}
         }
         return http.HandlerFunc(fn)
    }
    return m
}

// Main handlers /////////////////////////////////////////////////////////////////////////////////////
// gorilla/context bound to mongo db

// MGO Database Type //////////////////////////////////////////////////////////////////////////////////

type appContext struct {
  db *mgo.Database
}

func infoHandler(w http.ResponseWriter, r *http.Request) {
  sessionresrc := SessionResource{}
  sessionresrc.SData = current_session

  //  w.Header().Set("Content-Type", "application/vnd.api+json")
  //	json.NewEncoder(w).Encode(event)
  //	if err = show.Execute(w, json.NewEncoder(w).Encode(event)); err != nil {
  //      http.Error(w, err.Error(), http.StatusInternalServerError)
  //      return
  //  }

  // read JSON into BSON
  if err := showInfo.Execute(w, sessionresrc); err != nil {
      http.Error(w, err.Error(), http.StatusInternalServerError)
      return
  }
}


// Event Handlers /////////////////////////////////////////////////////////////////////////////////////

func (c *appContext) eventsHandler(w http.ResponseWriter, r *http.Request) {
  if current_session.Current_status == false{
    http.Redirect(w, r, "/login", 302)
  }else{
    eventsresrc := EventsResource{}
    eventsresrc.SData = current_session
    repo := EventRepo{c.db.C("events")}
    events, err := repo.All()
    for i:=0; i<len(events.Data); i++{
      body := EventResource{}
      body.Data.Id = events.Data[i].Id
      body.Data.Name = events.Data[i].Name
      body.Data.Location = events.Data[i].Location
      body.Data.Status = events.Data[i].Status
      body.Data.Host = events.Data[i].Host
      body.Data.Division = events.Data[i].Division
      body.Data.Event_Id = events.Data[i].Event_Id
      body.Data.Date = events.Data[i].Date
      eventsresrc.Data = append(eventsresrc.Data, body)
    }
    if err != nil {
      panic(err)
    }
    //	w.Header().Set("Content-Type", "application/vnd.api+json")
    //	json.NewEncoder(w).Encode(events)

    // read BSON into JSON
    if err = listEvent.Execute(w, eventsresrc); err != nil {
      http.Error(w, err.Error(), http.StatusInternalServerError)
      return
    }
  }
}

func (c *appContext) eventHandler(w http.ResponseWriter, r *http.Request) {
  if current_session.Current_status == false{
    http.Redirect(w, r, "/login", 302)
  }else{
    params := context.Get(r, "params").(httprouter.Params)  //gorrila context, key "params"
    evRepo := EventRepo{c.db.C("events")}
    event, err := evRepo.Find(params.ByName("id")) //getting data from named param :id
    enRepo := EntrantRepo{c.db.C("entrants")}
    entrants, err := enRepo.All()
    entrants_resrc := EntrantsResource{}
    for i:=0; i<len(event.Data.EntrantSelected_Id); i++{
      for j:=0; j<len(entrants.Data); j++{
        if entrants.Data[j].Team_Id == event.Data.EntrantSelected_Id[i]{
          for k:=0; k<len(entrants.Data[j].Event_Id); k++{
            if entrants.Data[j].Event_Id[k] == event.Data.Event_Id{
              body := EntrantResource{}
              body.Data.Id = entrants.Data[j].Id
              body.Data.Event_Id = entrants.Data[j].Event_Id
              body.Data.First_name = entrants.Data[j].First_name
              body.Data.Last_name = entrants.Data[j].Last_name
              body.Data.Team_Id = entrants.Data[j].Team_Id
              body.Data.Dog_name = entrants.Data[j].Dog_name
              body.Data.Dog_id_number = entrants.Data[j].Dog_id_number
              body.Data.Id_number = entrants.Data[j].Id_number
              entrants_resrc.Data = append(entrants_resrc.Data, body)
            }
          }
        }
      }
    }

    usrRepo := UserRepo{c.db.C("users")}
    users, err := usrRepo.All()
    users_resrc := UsersResource{}
    for i:=0; i<len(event.Data.UserSelected_Id); i++{
      for j:=0; j<len(users.Data); j++{
        if users.Data[j].User_Id == event.Data.UserSelected_Id[i]{
          for k:=0; k<len(users.Data[j].Event_Id); k++{
            if users.Data[j].Event_Id[k] == event.Data.Event_Id{
              body := UserResource{}
              body.Data.Id = users.Data[j].Id
              body.Data.Event_Id = users.Data[j].Event_Id
              body.Data.First_name = users.Data[j].First_name
              body.Data.Last_name = users.Data[j].Last_name
              body.Data.User_Id = users.Data[j].User_Id
              body.Data.Email = users.Data[j].Email
              body.Data.Password = users.Data[j].Password
              body.Data.Role = users.Data[j].Role
              users_resrc.Data = append(users_resrc.Data, body)
            }
          }
        }
      }
    }

    tlyRepo := TallyRepo{c.db.C("tallies")}
    tallies, err := tlyRepo.All()
    tallies_resrc := TalliesResource{}
    for i:=0; i<len(event.Data.EntrantSelected_Id); i++{
      for j:=0; j<len(tallies.Data); j++{
        if tallies.Data[j].Entrant_Id == event.Data.EntrantSelected_Id[i] && tallies.Data[j].Event_Id == event.Data.Event_Id{
          body := TallyResource{}
          body.Data.Id = tallies.Data[j].Id
          body.Data.Tally_Id = tallies.Data[j].Tally_Id
          body.Data.Event_Id = tallies.Data[j].Event_Id
          body.Data.Entrant_Id = tallies.Data[j].Entrant_Id
          body.Data.Total_points = tallies.Data[j].Total_points
          body.Data.Total_faults = tallies.Data[j].Total_faults
          body.Data.Total_time = tallies.Data[j].Total_time
          body.Data.Title = tallies.Data[j].Title
          body.Data.Qualifying_score = tallies.Data[j].Qualifying_score
          body.Data.Qualifying_scores = tallies.Data[j].Qualifying_scores
          tallies_resrc.Data = append(tallies_resrc.Data, body)
        }
      }
    }

    scRepo := ScorecardRepo{c.db.C("scorecards")}
    scorecards, err := scRepo.All()
    scorecards_resrc := ScorecardsResource{}
    for i:=0; i<len(event.Data.EntrantSelected_Id); i++{
      for j:=0; j<len(scorecards.Data); j++{
        if scorecards.Data[j].Entrant_Id == event.Data.EntrantSelected_Id[i] && scorecards.Data[j].Event_Id == event.Data.Event_Id{
          body := ScorecardResource{}
          body.Data.Id = scorecards.Data[j].Id
          body.Data.Event_Id = scorecards.Data[j].Event_Id
          body.Data.Entrant_Id = scorecards.Data[j].Entrant_Id
          body.Data.Search_area = scorecards.Data[j].Search_area
          body.Data.Element = scorecards.Data[j].Element
          body.Data.Total_points = scorecards.Data[j].Total_points
          body.Data.Total_faults = scorecards.Data[j].Total_faults
          body.Data.Total_time = scorecards.Data[j].Total_time
          scorecards_resrc.Data = append(scorecards_resrc.Data, body)
        }
      }
    }

    if err != nil {
      panic(err)
    }
    eventshow := EventShowResource{}
    eventshow.EVData = event
    eventshow.SData = current_session
    eventshow.ENSData = entrants_resrc
    eventshow.USRSData = users_resrc
    eventshow.TLYSData = tallies_resrc
    eventshow.SCSData = scorecards_resrc
    eventshow.Rank = c.place_order(event.Data.Id.Hex())
    eventshow.SCScompleted = c.scorecard_completion(event.Data.Id.Hex())
    eventshow.TLYcompleted = c.tally_completion(event.Data.Id.Hex())

    //  w.Header().Set("Content-Type", "application/vnd.api+json")
    //	json.NewEncoder(w).Encode(event)
    //	if err = show.Execute(w, json.NewEncoder(w).Encode(event)); err != nil {
    //      http.Error(w, err.Error(), http.StatusInternalServerError)
    //      return
    //  }

    // read JSON into BSON
    if err = showEvent.Execute(w, eventshow); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
  }
}

func (c *appContext) newEventHandler(w http.ResponseWriter, r *http.Request) {
  if current_session.Current_status == false{
    http.Redirect(w, r, "/login", 302)
  }else{
    eventresrc := EventEditResource{}
    eventresrc.SData = current_session

    enRepo := EntrantRepo{c.db.C("entrants")}
    entrants, err := enRepo.All()
    entrants_resrc := EntrantsResource{}
    for j:=0; j<len(entrants.Data); j++{
      body := EntrantResource{}
      body.Data.Id = entrants.Data[j].Id
      body.Data.First_name = entrants.Data[j].First_name
      body.Data.Last_name = entrants.Data[j].Last_name
      body.Data.Team_Id = entrants.Data[j].Team_Id
      body.Data.Dog_name = entrants.Data[j].Dog_name
      body.Data.Dog_id_number = entrants.Data[j].Dog_id_number
      body.Data.Id_number = entrants.Data[j].Id_number
      entrants_resrc.Data = append(entrants_resrc.Data, body)
    }

    usrRepo := UserRepo{c.db.C("users")}
    users, err := usrRepo.All()
    users_resrc := UsersResource{}
    for j:=0; j<len(users.Data); j++{
      body := UserResource{}
      body.Data.Id = users.Data[j].Id
      body.Data.First_name = users.Data[j].First_name
      body.Data.Last_name = users.Data[j].Last_name
      body.Data.User_Id = users.Data[j].User_Id
      body.Data.Email = users.Data[j].Email
      body.Data.Password = users.Data[j].Password
      body.Data.Role = users.Data[j].Role
      users_resrc.Data = append(users_resrc.Data, body)
    }

    eventresrc.ENSData = entrants_resrc
    eventresrc.USRSData = users_resrc

    body := EventResource{}
    body.Data.Int_search_areas = "0"
    body.Data.Ext_search_areas = "0"
    body.Data.Cont_search_areas = "0"
    body.Data.Veh_search_areas = "0"
    body.Data.Elite_search_areas = "0"
    body.Data.Int_hides = "0"
    body.Data.Ext_hides = "0"
    body.Data.Cont_hides = "0"
    body.Data.Veh_hides = "0"
    body.Data.Elite_hides = "0"

    for i:=0; i<len(entrants.Data); i++{
      newEntrant := Selected{Value: entrants.Data[i].Team_Id, Selected: false}
      body.Data.EntrantAll_Id = append(body.Data.EntrantAll_Id, newEntrant)
    }
    for i:=0; i<len(users.Data); i++{
      newUser := Selected{Value: users.Data[i].User_Id, Selected: false}
      body.Data.UserAll_Id = append(body.Data.UserAll_Id, newUser)
    }

    eventresrc.EVData = body

    if err := createnewEvent.Execute(w, eventresrc); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    fmt.Println(err)
  }
}

func (c *appContext) createEventHandler(w http.ResponseWriter, r *http.Request) {
  if current_session.Current_status == false{
    http.Redirect(w, r, "/login", 302)
  }else{
    evRepo := EventRepo{c.db.C("events")}
    events, err := evRepo.All()
    rrcount := 0
    evbody := context.Get(r, "body").(*EventResource)    //gorilla context, key "body" that returns val

    enRepo := EntrantRepo{c.db.C("entrants")}
    entrants, err := enRepo.All()

    usRepo := UserRepo{c.db.C("users")}
    users, err := usRepo.All()

    for i:=0; i<len(entrants.Data); i++{
      newEntrant := Selected{Value: entrants.Data[i].Team_Id, Selected: false}
      evbody.Data.EntrantAll_Id = append(evbody.Data.EntrantAll_Id, newEntrant)
    }
    for i:=0; i<len(users.Data); i++{
      newUser := Selected{Value: users.Data[i].User_Id, Selected: false}
      evbody.Data.UserAll_Id = append(evbody.Data.UserAll_Id, newUser)
    }
    evbody.Data.Name = r.FormValue("Name")
    evbody.Data.Location = r.FormValue("Location")
    evbody.Data.Date = r.FormValue("Date")
    evbody.Data.Host = r.FormValue("Host")
    evbody.Data.Status = r.FormValue("Status")
    evbody.Data.Division = r.FormValue("Division")
    evbody.Data.Event_Id = "EV_" + strconv.Itoa(rand.Int())
    for r:=0; r<len(events.Data); r++{
      if evbody.Data.Event_Id == events.Data[r].Event_Id{
        // Event_Id duplicate found - re-naming and re-checking loop 1 to loop 2
        evbody.Data.Event_Id = "EV_" + strconv.Itoa(rand.Int())
      }
      for rr:=0;rr<len(events.Data); rr++{
        if evbody.Data.Event_Id == events.Data[rr].Event_Id{
          // Event_Id duplicate found - re-naming and re-checking loop 2 to outer loop 1
          evbody.Data.Event_Id = "EV_" + strconv.Itoa(rand.Int())
          break
        }else{
          rrcount = rr
        }
      }
      if rrcount == len(events.Data)-1{
        // No duplicates both loops
        break
      }
    }
    evbody.Data.Int_search_areas = r.FormValue("Int_search_areas")
    if evbody.Data.Int_search_areas == ""{
      evbody.Data.Int_search_areas = "0"
    }
    evbody.Data.Ext_search_areas = r.FormValue("Ext_search_areas")
    if evbody.Data.Ext_search_areas == ""{
      evbody.Data.Ext_search_areas = "0"
    }
    evbody.Data.Cont_search_areas = r.FormValue("Cont_search_areas")
    if evbody.Data.Cont_search_areas == ""{
      evbody.Data.Cont_search_areas = "0"
    }
    evbody.Data.Veh_search_areas = r.FormValue("Veh_search_areas")
    if evbody.Data.Veh_search_areas == ""{
      evbody.Data.Veh_search_areas = "0"
    }
    evbody.Data.Elite_search_areas = r.FormValue("Elite_search_areas")
    if evbody.Data.Elite_search_areas == ""{
      evbody.Data.Elite_search_areas = "0"
    }
    evbody.Data.Int_hides = r.FormValue("Int_hides")
    if evbody.Data.Int_hides == ""{
      evbody.Data.Int_hides = "0"
    }
    evbody.Data.Ext_hides = r.FormValue("Ext_hides")
    if evbody.Data.Ext_hides == ""{
      evbody.Data.Ext_hides = "0"
    }
    evbody.Data.Cont_hides = r.FormValue("Cont_hides")
    if evbody.Data.Cont_hides == ""{
      evbody.Data.Cont_hides = "0"
    }
    evbody.Data.Veh_hides = r.FormValue("Veh_hides")
    if evbody.Data.Veh_hides == ""{
      evbody.Data.Veh_hides = "0"
    }
    evbody.Data.Elite_hides = r.FormValue("Elite_hides")
    if evbody.Data.Elite_hides == ""{
      evbody.Data.Elite_hides = "0"
    }
    evbody.Data.EntrantSelected_Id = r.Form["EntrantSelected_Id"]
    evbody.Data.UserSelected_Id = r.Form["UserSelected_Id"]
    for i:=0; i<len(evbody.Data.EntrantAll_Id); i++{
      evbody.Data.EntrantAll_Id[i].Selected = false
      if len(evbody.Data.EntrantSelected_Id)>0{
        for j:=0; j<len(evbody.Data.EntrantSelected_Id); j++{
          if evbody.Data.EntrantAll_Id[i].Value == evbody.Data.EntrantSelected_Id[j]{
            evbody.Data.EntrantAll_Id[i].Selected = true
          }
        }
      }
    }
    for i:=0; i<len(evbody.Data.UserAll_Id); i++{
      evbody.Data.UserAll_Id[i].Selected = false
      if len(evbody.Data.UserSelected_Id)>0{
        for j:=0; j<len(evbody.Data.UserSelected_Id); j++{
          if evbody.Data.UserAll_Id[i].Value == evbody.Data.UserSelected_Id[j]{
            evbody.Data.UserAll_Id[i].Selected = true
          }
        }
      }
    }
    err, id := evRepo.Create(&evbody.Data)

    // Add Event_Id to selected entrants and add scorecards and tallies
    scRepo := ScorecardRepo{c.db.C("scorecards")}
    scorecards, err := scRepo.All()

    taRepo := TallyRepo{c.db.C("tallies")}
    tallies, err := taRepo.All()

    evRepo = EventRepo{c.db.C("events")}
    event, err := evRepo.Find(id.Hex()) //getting data from named param :id

    // Search through all of the entrants registered in event EntrantAll_Id
    for i:=0; i<len(event.Data.EntrantAll_Id); i++{
      // Search through all entrants
      for j:=0; j<len(entrants.Data); j++{
        // Search through entrants registered in event EntrantSelected_Id
        // If there is at least one entrant selected
        if len(event.Data.EntrantSelected_Id) > 0{
          if entrants.Data[j].Team_Id == event.Data.EntrantAll_Id[i].Value{
            if event.Data.EntrantAll_Id[i].Selected == true{
              // Register event in entrant and create scorecards and tally
              entrants.Data[j].Event_Id = append(entrants.Data[j].Event_Id, event.Data.Event_Id)
              enbody := EntrantResource{}
              enbody.Data.Id = entrants.Data[j].Id
              enbody.Data.First_name = entrants.Data[j].First_name
              enbody.Data.Last_name = entrants.Data[j].Last_name
              enbody.Data.Id_number = entrants.Data[j].Id_number
              enbody.Data.Dog_name = entrants.Data[j].Dog_name
              enbody.Data.Dog_id_number = entrants.Data[j].Dog_id_number
              enbody.Data.Breed = entrants.Data[j].Breed
              enbody.Data.Team_Id = entrants.Data[j].Team_Id
              enbody.Data.Event_Id = entrants.Data[j].Event_Id
              err = enRepo.Update(&enbody.Data)

              // create scorecards for this event and entrant
              search_areas := 0
              element := ""
              for elm:=0; elm<len(ELEMENTS); elm++{
                switch element = ELEMENTS[elm]; ELEMENTS[elm]{
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
                  for s:=1; s<=search_areas; s++{
                    scbody := ScorecardResource{}
                    scbody.Data.Element = element
                    scbody.Data.Event_Id = event.Data.Event_Id
                    scbody.Data.Entrant_Id = entrants.Data[j].Team_Id
                    scbody.Data.Search_area = strconv.Itoa(s)
                    scbody.Data.Scorecard_Id = "SC_" + strconv.Itoa(rand.Int())

                    // check for duplicates
                    rrcount := 0
                    for r:=0; r<len(scorecards.Data); r++{
                      if scbody.Data.Scorecard_Id == scorecards.Data[r].Scorecard_Id{
                        // Scorecard_Id duplicate found - re-naming and re-checking loop 1 to loop 2
                        scbody.Data.Scorecard_Id = "SC_" + strconv.Itoa(rand.Int())
                      }
                      for rr:=0;rr<len(scorecards.Data); rr++{
                        if scbody.Data.Scorecard_Id == scorecards.Data[rr].Scorecard_Id{
                          // Scorecard_Id duplicate found - re-naming and re-checking loop 2 to outer loop 1
                          scbody.Data.Scorecard_Id = "SC_" + strconv.Itoa(rand.Int())
                          break
                        }else{
                          rrcount = rr
                        }
                      }
                      if rrcount == len(scorecards.Data)-1{
                        // No duplicates both loops
                        break
                      }
                    }
                    scbody.Data.Hides_max = "0"
                    scbody.Data.Hides_found = "0"
                    scbody.Data.Hides_missed = "0"
                    scbody.Data.Maxpoint = "0"
                    scbody.Data.Total_time = "00:00:00"
                    scbody.Data.False_alert_fringe = "0"
                    scbody.Data.Finish_call = Selected{Value: "yes", Selected: true}
                    scbody.Data.Timed_out = Selected{Value: "no", Selected: false}
                    scbody.Data.Dismissed = Selected{Value: "no", Selected: false}
                    scbody.Data.Excused = Selected{Value: "no", Selected: false}
                    scbody.Data.Absent = Selected{Value: "no", Selected: false}
                    scbody.Data.Eliminated_during_search = Selected{Value: "no", Selected: false}
                    scbody.Data.Pronounced = Selected{Value: "no", Selected: false}
                    scbody.Data.Judge_signature = Selected{Value: "no", Selected: false}
                    err, id := scRepo.Create(&scbody.Data)
                    fmt.Println(err)
                    fmt.Println(id)
                  }
                }
              }
              tabody := TallyResource{}
              tabody.Data.Entrant_Id = entrants.Data[j].Team_Id
              tabody.Data.Event_Id = event.Data.Event_Id
              tabody.Data.Tally_Id = "TLY_" + strconv.Itoa(rand.Int())

              // check for duplicates
              rrcount := 0
              for r:=0; r<len(tallies.Data); r++{
                if tabody.Data.Tally_Id == tallies.Data[r].Tally_Id{
                // Tally_Id duplicate found - re-naming and re-checking loop 1 to loop 2
                  tabody.Data.Tally_Id = "TLY_" + strconv.Itoa(rand.Int())
                }
                for rr:=0;rr<len(tallies.Data); rr++{
                  if tabody.Data.Tally_Id == tallies.Data[rr].Tally_Id{
                    // Tally_Id duplicate found - re-naming and re-checking loop 2 to loop 1
                    tabody.Data.Tally_Id = "TLY_" + strconv.Itoa(rand.Int())
                    break
                  }else{
                    rrcount = rr
                  }
                }
                if rrcount == len(tallies.Data)-1{
                  // No duplicates both loops
                  break
                }
              }
              tabody.Data.Total_time = "0"
              tabody.Data.Total_faults = "0"
              tabody.Data.Title = "not this time"
              tabody.Data.Total_points = "0"
              tabody.Data.Qualifying_score = "0"
              tabody.Data.Qualifying_scores = "0"
              err, id := taRepo.Create(&tabody.Data)
              fmt.Println(err)
              fmt.Println(id)
            }
          }
        }
      }
    }
    for i:=0; i<len(event.Data.UserAll_Id); i++{
      for j:=0; j<len(users.Data); j++{
        if len(event.Data.UserSelected_Id) > 0{
          if users.Data[j].User_Id == event.Data.UserAll_Id[i].Value{
            if event.Data.UserAll_Id[i].Selected == true{
              users.Data[j].Event_Id = append(users.Data[j].Event_Id, event.Data.Event_Id)
              usbody := UserResource{}
              usbody.Data.Id = users.Data[j].Id
              usbody.Data.First_name = users.Data[j].First_name
              usbody.Data.Last_name = users.Data[j].Last_name
              usbody.Data.Status = users.Data[j].Status
              usbody.Data.Approved = users.Data[j].Approved
              usbody.Data.Email = users.Data[j].Email
              usbody.Data.Password = users.Data[j].Password
              usbody.Data.User_Id = users.Data[j].User_Id
              usbody.Data.Role = users.Data[j].Role
              usbody.Data.Event_Id = users.Data[j].Event_Id
              err = usRepo.Update(&usbody.Data)
            }
          }
        }
      }
    }
    if err != nil {
      panic(err)
    }
    if event.Data.Event_Id == ""{
      http.Redirect(w, r, "/events/delete/" + event.Data.Id.Hex(), 302)
    }else{
      http.Redirect(w, r, "/events/show/" + event.Data.Id.Hex(), 302)
    }
  }
}

func (c *appContext) editEventHandler(w http.ResponseWriter, r *http.Request) {
  if current_session.Current_status == false{
    http.Redirect(w, r, "/login", 302)
  }else{
    params := context.Get(r, "params").(httprouter.Params)    //gorilla context, key "params"
    evRepo := EventRepo{c.db.C("events")}
    event, err := evRepo.Find(params.ByName("id")) //getting data from named param :id

    eventresrc := EventEditResource{}

    enRepo := EntrantRepo{c.db.C("entrants")}
    entrants, err := enRepo.All()
    entrants_resrc := EntrantsResource{}
    for j:=0; j<len(entrants.Data); j++{
      body := EntrantResource{}
      body.Data.Id = entrants.Data[j].Id
      body.Data.Event_Id = entrants.Data[j].Event_Id
      body.Data.First_name = entrants.Data[j].First_name
      body.Data.Last_name = entrants.Data[j].Last_name
      body.Data.Team_Id = entrants.Data[j].Team_Id
      body.Data.Dog_name = entrants.Data[j].Dog_name
      body.Data.Dog_id_number = entrants.Data[j].Dog_id_number
      body.Data.Id_number = entrants.Data[j].Id_number
      entrants_resrc.Data = append(entrants_resrc.Data, body)
    }

    usrRepo := UserRepo{c.db.C("users")}
    users, err := usrRepo.All()
    users_resrc := UsersResource{}
    for j:=0; j<len(users.Data); j++{
      body := UserResource{}
      body.Data.Id = users.Data[j].Id
      body.Data.Event_Id = users.Data[j].Event_Id
      body.Data.First_name = users.Data[j].First_name
      body.Data.Last_name = users.Data[j].Last_name
      body.Data.User_Id = users.Data[j].User_Id
      body.Data.Email = users.Data[j].Email
      body.Data.Password = users.Data[j].Password
      body.Data.Role = users.Data[j].Role
      users_resrc.Data = append(users_resrc.Data, body)
    }

    eventresrc.ENSData = entrants_resrc
    eventresrc.USRSData = users_resrc

    if len(event.Data.EntrantAll_Id) == 0{
      for i:=0; i<len(entrants.Data); i++{
        newEntrant := Selected{Value: entrants.Data[i].Team_Id, Selected: false}
        event.Data.EntrantAll_Id = append(event.Data.EntrantAll_Id, newEntrant)
      }
    }
    for j:=0; j<len(event.Data.EntrantAll_Id); j++{
      for k:=0; k<len(event.Data.EntrantSelected_Id); k++{
        if event.Data.EntrantAll_Id[j].Value == event.Data.EntrantSelected_Id[k]{
          event.Data.EntrantAll_Id[j].Selected = true
        }
      }
    }
    if len(event.Data.UserAll_Id) == 0{
      for i:=0; i<len(users.Data); i++{
        newUser := Selected{Value: users.Data[i].User_Id, Selected: false}
        event.Data.UserAll_Id = append(event.Data.UserAll_Id, newUser)
      }
    }
    for j:=0; j<len(event.Data.UserAll_Id); j++{
      for k:=0; k<len(event.Data.UserSelected_Id); k++{
        if event.Data.UserAll_Id[j].Value == event.Data.UserSelected_Id[k]{
          event.Data.UserAll_Id[j].Selected = true
        }
      }
    }
    eventresrc.EVData = event
    eventresrc.SData = current_session

    if err = updateEvent.Execute(w, eventresrc); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
  }
}

func (c *appContext) updateEventHandler(w http.ResponseWriter, r *http.Request) {
  if current_session.Current_status == false{
    http.Redirect(w, r, "/login", 302)
  }else{
    params := context.Get(r, "params").(httprouter.Params)

    evRepo := EventRepo{c.db.C("events")}
    event, err := evRepo.Find(params.ByName("id")) //getting data from named param :id

    enRepo := EntrantRepo{c.db.C("entrants")}
    entrants, err := enRepo.All()

    usRepo := UserRepo{c.db.C("users")}
    users, err := usRepo.All()
    if len(event.Data.EntrantAll_Id) == 0{
      for i:=0; i<len(entrants.Data); i++{
          newEntrant := Selected{Value: entrants.Data[i].Team_Id, Selected: false}
          event.Data.EntrantAll_Id = append(event.Data.EntrantAll_Id, newEntrant)
      }
    }
    if len(event.Data.UserAll_Id) == 0{
      for i:=0; i<len(users.Data); i++{
        newUser := Selected{Value: users.Data[i].User_Id, Selected: false}
        event.Data.UserAll_Id = append(event.Data.UserAll_Id, newUser)
      }
    }
    evbody := context.Get(r, "body").(*EventResource)
    evbody.Data.Id = event.Data.Id
    evbody.Data.Name = r.FormValue("Name")
    evbody.Data.Location = r.FormValue("Location")
    evbody.Data.Date = r.FormValue("Date")
    evbody.Data.Host = r.FormValue("Host")
    evbody.Data.Status = r.FormValue("Status")
    evbody.Data.Division = r.FormValue("Division")
    evbody.Data.Event_Id = event.Data.Event_Id
    evbody.Data.Int_search_areas = r.FormValue("Int_search_areas")
    if evbody.Data.Int_search_areas == ""{
      evbody.Data.Int_search_areas = "0"
    }
    evbody.Data.Ext_search_areas = r.FormValue("Ext_search_areas")
    if evbody.Data.Ext_search_areas == ""{
      evbody.Data.Ext_search_areas = "0"
    }
    evbody.Data.Cont_search_areas = r.FormValue("Cont_search_areas")
    if evbody.Data.Cont_search_areas == ""{
      evbody.Data.Cont_search_areas = "0"
    }
    evbody.Data.Veh_search_areas = r.FormValue("Veh_search_areas")
    if evbody.Data.Veh_search_areas == ""{
      evbody.Data.Veh_search_areas = "0"
    }
    evbody.Data.Elite_search_areas = r.FormValue("Elite_search_areas")
    if evbody.Data.Elite_search_areas == ""{
      evbody.Data.Elite_search_areas = "0"
    }
    evbody.Data.Int_hides = r.FormValue("Int_hides")
    if evbody.Data.Int_hides == ""{
      evbody.Data.Int_hides = "0"
    }
    evbody.Data.Ext_hides = r.FormValue("Ext_hides")
    if evbody.Data.Ext_hides == ""{
      evbody.Data.Ext_hides = "0"
    }
    evbody.Data.Cont_hides = r.FormValue("Cont_hides")
    if evbody.Data.Cont_hides == ""{
      evbody.Data.Cont_hides = "0"
    }
    evbody.Data.Veh_hides = r.FormValue("Veh_hides")
    if evbody.Data.Veh_hides == ""{
      evbody.Data.Veh_hides = "0"
    }
    evbody.Data.Elite_hides = r.FormValue("Elite_hides")
    if evbody.Data.Elite_hides == ""{
      evbody.Data.Elite_hides = "0"
    }
    evbody.Data.EntrantAll_Id = event.Data.EntrantAll_Id
    evbody.Data.EntrantSelected_Id = r.Form["EntrantSelected_Id"]
    evbody.Data.UserAll_Id = event.Data.UserAll_Id
    evbody.Data.UserSelected_Id = r.Form["UserSelected_Id"]
    for i:=0; i<len(event.Data.EntrantAll_Id); i++{
      evbody.Data.EntrantAll_Id[i].Selected = false
      if len(evbody.Data.EntrantSelected_Id)>0{
        for j:=0; j<len(evbody.Data.EntrantSelected_Id); j++{
          if event.Data.EntrantAll_Id[i].Value == evbody.Data.EntrantSelected_Id[j]{
            evbody.Data.EntrantAll_Id[i].Selected = true
          }
        }
      }
    }
    for i:=0; i<len(event.Data.UserAll_Id); i++{
      evbody.Data.UserAll_Id[i].Selected = false
      if len(evbody.Data.UserSelected_Id)>0{
        for j:=0; j<len(evbody.Data.UserSelected_Id); j++{
          if event.Data.UserAll_Id[i].Value == evbody.Data.UserSelected_Id[j]{
            evbody.Data.UserAll_Id[i].Selected = true
          }
        }
      }
    }
    err = evRepo.Update(&evbody.Data)

    // Add/remove Event_Id to selected entrants and add/remove scorecards and tallies
    scRepo := ScorecardRepo{c.db.C("scorecards")}
    scorecards, err := scRepo.All()

    taRepo := TallyRepo{c.db.C("tallies")}
    tallies, err := taRepo.All()

    evRepo = EventRepo{c.db.C("events")}
    event, err = evRepo.Find(params.ByName("id")) //getting data from named param :id

    found := false
    // Search through all of the entrants registered in event EntrantAll_Id
    for i:=0; i<len(event.Data.EntrantAll_Id); i++{
      // Search through all entrants
      for j:=0; j<len(entrants.Data); j++{
        found = false
        // Search through entrants registered in event EntrantSelected_Id
        // If there is at least one entrant selected
        if len(event.Data.EntrantSelected_Id) > 0{
          // If there is a match between entrants and entrants registered in event EntrantAll_Id
          if entrants.Data[j].Team_Id == event.Data.EntrantAll_Id[i].Value{
            // If the entrant is registered in event EntrantSelected_Id
            if event.Data.EntrantAll_Id[i].Selected == true{
              // Search through to see if event is registered in entrant Event_Id
              for k:=0; k<len(entrants.Data[j].Event_Id); k++{
                if entrants.Data[j].Event_Id[k] == event.Data.Event_Id{
                  found = true
                  search_areas := 0
                  element := ""
                  // The event registers the entrant and the entrant registers the event
                  // check for element and search area updates
                  for elm:=0; elm<len(ELEMENTS); elm++{
                    switch element = ELEMENTS[elm]; ELEMENTS[elm]{
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
                    scorecardfound := false

                    // For each element check to see if there are scorecards for the registered event and entrant
                    // If there is at least one scorecard for the element
                    if search_areas > 0{
                      // Survey scorecards for match, decrement when scorecard is added.
                      sc_count := 0
                      for s:=1; s<=search_areas; s++{
                        // if there is at least 1 scorecard in the db
                        if len(scorecards.Data) > 0{
                          scid := ""
                          // Search all scorecards for a match
                          for sc:=0; sc<len(scorecards.Data) - sc_count; sc++{
                            if scorecards.Data[sc].Element == element && scorecards.Data[sc].Search_area == strconv.Itoa(s) && scorecards.Data[sc].Entrant_Id == entrants.Data[j].Team_Id && scorecards.Data[sc].Event_Id == event.Data.Event_Id{
                              // A match is found, get id to prepare for update
                              scorecardfound = true
                              scid = scorecards.Data[sc].Id.Hex()
                            }
                          }
                          if scorecardfound == false{
                            // if scorecard does not exist, create it
                            scbody := ScorecardResource{}
                            scbody.Data.Event_Id = event.Data.Event_Id
                            scbody.Data.Entrant_Id = entrants.Data[j].Team_Id
                            scbody.Data.Element = element
                            scbody.Data.Search_area = strconv.Itoa(s)
                            scbody.Data.Hides_max = "0"
                            scbody.Data.Hides_found = "0"
                            scbody.Data.Hides_missed = "0"
                            scbody.Data.Maxpoint = "0"
                            scbody.Data.False_alert_fringe = "0"
                            scbody.Data.Total_time = "00:00:00"
                            scbody.Data.Finish_call = Selected{Value: "yes", Selected: true}
                            scbody.Data.Timed_out = Selected{Value: "no", Selected: false}
                            scbody.Data.Dismissed = Selected{Value: "no", Selected: false}
                            scbody.Data.Excused = Selected{Value: "no", Selected: false}
                            scbody.Data.Absent = Selected{Value: "no", Selected: false}
                            scbody.Data.Eliminated_during_search = Selected{Value: "no", Selected: false}
                            scbody.Data.Pronounced = Selected{Value: "no", Selected: false}
                            scbody.Data.Judge_signature = Selected{Value: "no", Selected: false}
                            scbody.Data.Scorecard_Id = "SC_" + strconv.Itoa(rand.Int())
                            // check for duplicates
                            rrcount := 0
                            for r:=0; r<len(scorecards.Data); r++{
                              if scbody.Data.Scorecard_Id == scorecards.Data[r].Scorecard_Id{
                                // Scorecard_Id duplicate found - re-naming and re-checking loop 1 to loop 2
                                scbody.Data.Scorecard_Id = "SC_" + strconv.Itoa(rand.Int())
                              }
                              for rr:=0;rr<len(scorecards.Data); rr++{
                                if scbody.Data.Scorecard_Id == scorecards.Data[rr].Scorecard_Id{
                                  // Scorecard_Id duplicate found - re-naming and re-checking loop 2 to outer loop 1
                                  scbody.Data.Scorecard_Id = "SC_" + strconv.Itoa(rand.Int())
                                  break
                                }else{
                                  rrcount = rr
                                }
                              }
                              if rrcount == len(scorecards.Data)-1{
                                // No duplicates both loops
                                break
                              }
                            }
                            err, id := scRepo.Create(&scbody.Data)

                            // decrement range
                            sc_count += 1
                            fmt.Println(err)
                            fmt.Println(id)
                          }else if scorecardfound == true{
                            // if scorecard exists, update it
                            scorecard, err := scRepo.Find(scid)
                            scbody := ScorecardResource{}
                            scbody.Data.Id = scorecard.Data.Id
                            scbody.Data.Event_Id = scorecard.Data.Event_Id
                            scbody.Data.Entrant_Id = scorecard.Data.Entrant_Id
                            scbody.Data.Element = element
                            scbody.Data.Search_area = strconv.Itoa(s)
                            scbody.Data.Hides_max = scorecard.Data.Hides_max
                            scbody.Data.Hides_found = scorecard.Data.Hides_found
                            scbody.Data.Hides_missed = scorecard.Data.Hides_missed
                            scbody.Data.Maxpoint = scorecard.Data.Maxpoint
                            scbody.Data.False_alert_fringe = scorecard.Data.False_alert_fringe
                            scbody.Data.Finish_call = scorecard.Data.Finish_call
                            scbody.Data.Timed_out = scorecard.Data.Timed_out
                            scbody.Data.Dismissed = scorecard.Data.Dismissed
                            scbody.Data.Excused = scorecard.Data.Excused
                            scbody.Data.Total_points = scorecard.Data.Total_points
                            scbody.Data.Total_faults = scorecard.Data.Total_faults
                            scbody.Data.Total_time = scorecard.Data.Total_time
                            scbody.Data.Maxtime_m = scorecard.Data.Maxtime_m
                            scbody.Data.Maxtime_s = scorecard.Data.Maxtime_s
                            scbody.Data.Absent = scorecard.Data.Absent
                            scbody.Data.Eliminated_during_search = scorecard.Data.Eliminated_during_search
                            scbody.Data.Pronounced = scorecard.Data.Pronounced
                            scbody.Data.Judge_signature = scorecard.Data.Judge_signature
                            scbody.Data.Scorecard_Id = scorecard.Data.Scorecard_Id
                            err = scRepo.Update(&scbody.Data)
                            fmt.Println(err)
                            // find match for other scorecards
                            scorecardfound = false
                          }
                        }
                      }
                    }
                    // if there is no search area, delete scorecard
                    if search_areas == 0{
                      // No search areas
                      if len(scorecards.Data) > 0{
                        for sc:=0; sc<len(scorecards.Data); sc++{
                          searchArea, err := strconv.Atoi(scorecards.Data[sc].Search_area)
                          if scorecards.Data[sc].Element == element && searchArea > 0 && scorecards.Data[sc].Entrant_Id == entrants.Data[j].Team_Id && scorecards.Data[sc].Event_Id == event.Data.Event_Id{
                            id := scorecards.Data[sc].Id.Hex()
                            err = scRepo.Delete(id)
                            fmt.Println(err)
                          }
                        }
                      }
                    }
                  }
                  tallyfound := false
                  for ta:=0; ta<len(tallies.Data); ta++{
                    if tallies.Data[ta].Entrant_Id == entrants.Data[j].Team_Id{
                      tallyfound = true
                    }
                  }
                  if tallyfound == false{
                    tabody := TallyResource{}
                    tabody.Data.Entrant_Id = entrants.Data[j].Team_Id
                    tabody.Data.Event_Id = event.Data.Event_Id
                    tabody.Data.Tally_Id = "TLY_" + strconv.Itoa(rand.Int())
                    // check for duplicates
                    rrcount := 0
                    for r:=0; r<len(tallies.Data); r++{
                      if tabody.Data.Tally_Id == tallies.Data[r].Tally_Id{
                        // Tally_Id duplicate found - re-naming and re-checking loop 1 to loop 2
                        tabody.Data.Tally_Id = "TLY_" + strconv.Itoa(rand.Int())
                      }
                      for rr:=0;rr<len(tallies.Data); rr++{
                        if tabody.Data.Tally_Id == tallies.Data[rr].Tally_Id{
                          // Tally_Id duplicate found - re-naming and re-checking loop 2 to loop 1
                          tabody.Data.Tally_Id = "TLY_" + strconv.Itoa(rand.Int())
                          break
                        }else{
                          rrcount = rr
                        }
                      }
                      if rrcount == len(tallies.Data)-1{
                        // No duplicates both loops
                        break
                      }
                    }
                    tabody.Data.Total_time = "0"
                    tabody.Data.Total_faults = "0"
                    tabody.Data.Title = "not this time"
                    tabody.Data.Total_points = "0"
                    tabody.Data.Qualifying_score = "0"
                    tabody.Data.Qualifying_scores = "0"
                    err, id := taRepo.Create(&tabody.Data)
                    fmt.Println(err)
                    fmt.Println(id)
                  }
                }
              }
              // Event was not registered for entrant, though the entrant was registered with event EntrantSelected_Id
              // Register event in entrant Event_Id, add scorecards and tally
              if found == false{
                entrants.Data[j].Event_Id = append(entrants.Data[j].Event_Id, event.Data.Event_Id)
                enbody := EntrantResource{}
                enbody.Data.Id = entrants.Data[j].Id
                enbody.Data.First_name = entrants.Data[j].First_name
                enbody.Data.Last_name = entrants.Data[j].Last_name
                enbody.Data.Id_number = entrants.Data[j].Id_number
                enbody.Data.Dog_name = entrants.Data[j].Dog_name
                enbody.Data.Dog_id_number = entrants.Data[j].Dog_id_number
                enbody.Data.Breed = entrants.Data[j].Breed
                enbody.Data.Team_Id = entrants.Data[j].Team_Id
                enbody.Data.Event_Id = entrants.Data[j].Event_Id
                err = enRepo.Update(&enbody.Data)

                // create scorecards for this event and entrant
                search_areas := 0
                element := ""
                for elm:=0; elm<len(ELEMENTS); elm++{
                  switch element = ELEMENTS[elm]; ELEMENTS[elm]{
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
                    for s:=1; s<=search_areas; s++{
                      scbody := ScorecardResource{}
                      scbody.Data.Element = element
                      scbody.Data.Event_Id = event.Data.Event_Id
                      scbody.Data.Entrant_Id = entrants.Data[j].Team_Id
                      scbody.Data.Search_area = strconv.Itoa(s)
                      scbody.Data.Scorecard_Id = "SC_" + strconv.Itoa(rand.Int())

                      // check for duplicates
                      rrcount := 0
                      for r:=0; r<len(scorecards.Data); r++{
                        if scbody.Data.Scorecard_Id == scorecards.Data[r].Scorecard_Id{
                          // Scorecard_Id duplicate found - re-naming and re-checking loop 1 to loop 2
                          scbody.Data.Scorecard_Id = "SC_" + strconv.Itoa(rand.Int())
                        }
                        for rr:=0;rr<len(scorecards.Data); rr++{
                          if scbody.Data.Scorecard_Id == scorecards.Data[rr].Scorecard_Id{
                            // Scorecard_Id duplicate found - re-naming and re-checking loop 2 to outer loop 1
                            scbody.Data.Scorecard_Id = "SC_" + strconv.Itoa(rand.Int())
                            break
                          }else{
                            rrcount = rr
                          }
                        }
                        if rrcount == len(scorecards.Data)-1{
                          // No duplicates both loops
                          break
                        }
                      }
                      scbody.Data.Hides_max = "0"
                      scbody.Data.Hides_found = "0"
                      scbody.Data.Hides_missed = "0"
                      scbody.Data.Maxpoint = "0"
                      scbody.Data.False_alert_fringe = "0"
                      scbody.Data.Finish_call = Selected{Value: "yes", Selected: true}
                      scbody.Data.Timed_out = Selected{Value: "no", Selected: false}
                      scbody.Data.Dismissed = Selected{Value: "no", Selected: false}
                      scbody.Data.Excused = Selected{Value: "no", Selected: false}
                      scbody.Data.Absent = Selected{Value: "no", Selected: false}
                      scbody.Data.Eliminated_during_search = Selected{Value: "no", Selected: false}
                      scbody.Data.Pronounced = Selected{Value: "no", Selected: false}
                      scbody.Data.Judge_signature = Selected{Value: "no", Selected: false}
                      err, id := scRepo.Create(&scbody.Data)
                      fmt.Println(err)
                      fmt.Println(id)
                    }
                  }
                }
                tabody := TallyResource{}
                tabody.Data.Entrant_Id = entrants.Data[j].Team_Id
                tabody.Data.Event_Id = event.Data.Event_Id
                tabody.Data.Tally_Id = "TLY_" + strconv.Itoa(rand.Int())

                // check for duplicates
                rrcount := 0
                for r:=0; r<len(tallies.Data); r++{
                  if tabody.Data.Tally_Id == tallies.Data[r].Tally_Id{
                    // Tally_Id duplicate found - re-naming and re-checking loop 1 to loop 2
                    tabody.Data.Tally_Id = "TLY_" + strconv.Itoa(rand.Int())
                  }
                  for rr:=0;rr<len(tallies.Data); rr++{
                    if tabody.Data.Tally_Id == tallies.Data[rr].Tally_Id{
                      // Tally_Id duplicate found - re-naming and re-checking loop 2 to loop 1
                      tabody.Data.Tally_Id = "TLY_" + strconv.Itoa(rand.Int())
                      break
                    }else{
                      rrcount = rr
                    }
                  }
                  if rrcount == len(tallies.Data)-1{
                    // No duplicates both loops
                    break
                  }
                }
                tabody.Data.Total_time = "0"
                tabody.Data.Total_faults = "0"
                tabody.Data.Title = "not this time"
                tabody.Data.Total_points = "0"
                tabody.Data.Qualifying_score = "0"
                tabody.Data.Qualifying_scores = "0"
                err, id := taRepo.Create(&tabody.Data)
                fmt.Println(err)
                fmt.Println(id)
              }
              found = false

            // If entrant IS NOT registered in event EntrantSelected_Id, if the event IS registered in entrant Event_Id, deregister it and delete scorecards and tally.
            }else if event.Data.EntrantAll_Id[i].Selected == false{
              for k:=0; k<len(entrants.Data[j].Event_Id); k++{
                if entrants.Data[j].Event_Id[k] == event.Data.Event_Id{
                  k1 := 0
                  k2 := 0
                  if k > 0{
                    k1 = k - 1
                  }else if k == 0{
                    k1 = 1
                  }
                  alength := len(entrants.Data[j].Event_Id) - 1
                  if k < alength{
                    k2 = k + 1
                  }
                  if k == 0{
                    entrants.Data[j].Event_Id = entrants.Data[j].Event_Id[k1:]
                  }else if k > 0 && k < alength{
                    entrants.Data[j].Event_Id = append(entrants.Data[j].Event_Id[:k1], entrants.Data[j].Event_Id[k2:]...)
                  }else if k == alength{
                    entrants.Data[j].Event_Id = entrants.Data[j].Event_Id[:alength]
                  }
                  enbody := EntrantResource{}
                  enbody.Data.Id = entrants.Data[j].Id
                  enbody.Data.First_name = entrants.Data[j].First_name
                  enbody.Data.Last_name = entrants.Data[j].Last_name
                  enbody.Data.Id_number = entrants.Data[j].Id_number
                  enbody.Data.Dog_name = entrants.Data[j].Dog_name
                  enbody.Data.Dog_id_number = entrants.Data[j].Dog_id_number
                  enbody.Data.Breed = entrants.Data[j].Breed
                  enbody.Data.Team_Id = entrants.Data[j].Team_Id
                  enbody.Data.Event_Id = entrants.Data[j].Event_Id
                  err = enRepo.Update(&enbody.Data)

                  // delete scorecards for this event and this entrant
                  if len(scorecards.Data)>0{
                    for sc:=0; sc<len(scorecards.Data); sc++{
                      if scorecards.Data[sc].Entrant_Id == entrants.Data[j].Team_Id && scorecards.Data[sc].Event_Id == event.Data.Event_Id{
                        id := scorecards.Data[sc].Id.Hex()
                        err = scRepo.Delete(id)
                        fmt.Println(err)
                      }
                    }
                    for ta:=0; ta<len(tallies.Data); ta++{
                      if tallies.Data[ta].Entrant_Id == entrants.Data[j].Team_Id && tallies.Data[ta].Event_Id == event.Data.Event_Id{
                        id := tallies.Data[ta].Id.Hex()
                        err = taRepo.Delete(id)
                        fmt.Println(err)
                      }
                    }
                  }
                }
              }
            }
          }
        // if no entrants are registered in event EntrantSelected_Id, deregister the event from entrant Event_Id and delete scorecards and tallies.
        }else if len(event.Data.EntrantSelected_Id) <= 0{
          for k:=0; k<len(entrants.Data[j].Event_Id); k++{
            if entrants.Data[j].Event_Id[k] == event.Data.Event_Id{
              k1 := 0
              k2 := 0
              if k > 0{
                k1 = k - 1
              }else if k == 0{
                k1 = 1
              }
              alength := len(entrants.Data[j].Event_Id) - 1
              if k < alength{
                k2 = k + 1
              }
              if k == 0{
                entrants.Data[j].Event_Id = entrants.Data[j].Event_Id[k1:]
              }else if k > 0 && k < alength{
                entrants.Data[j].Event_Id = append(entrants.Data[j].Event_Id[:k1], entrants.Data[j].Event_Id[k2:]...)
              }else if k == alength{
                entrants.Data[j].Event_Id = entrants.Data[j].Event_Id[:alength]
              }
              enbody := EntrantResource{}
              enbody.Data.Id = entrants.Data[j].Id
              enbody.Data.First_name = entrants.Data[j].First_name
              enbody.Data.Last_name = entrants.Data[j].Last_name
              enbody.Data.Id_number = entrants.Data[j].Id_number
              enbody.Data.Dog_name = entrants.Data[j].Dog_name
              enbody.Data.Dog_id_number = entrants.Data[j].Dog_id_number
              enbody.Data.Breed = entrants.Data[j].Breed
              enbody.Data.Team_Id = entrants.Data[j].Team_Id
              enbody.Data.Event_Id = entrants.Data[j].Event_Id
              err = enRepo.Update(&enbody.Data)

              // delete all scorecards for this event and entrant
              if len(scorecards.Data)>0{
                for sc:=0; sc<len(scorecards.Data); sc++{
                  if scorecards.Data[sc].Event_Id == event.Data.Event_Id && scorecards.Data[sc].Entrant_Id == entrants.Data[j].Team_Id{
                    id := scorecards.Data[sc].Id.Hex()
                    err = scRepo.Delete(id)
                    fmt.Println(err)
                  }
                }
              }
              // delete all tallies for this event and entrant
              if len(tallies.Data)>0{
                for ta:=0; ta<len(tallies.Data); ta++{
                  if tallies.Data[ta].Event_Id == event.Data.Event_Id && tallies.Data[ta].Entrant_Id == entrants.Data[j].Team_Id{
                    id := tallies.Data[ta].Id.Hex()
                    err = taRepo.Delete(id)
                    fmt.Println(err)
                  }
                }
              }
            }
          }
        }
      }
    }
    found = false
    // Similar routine with users
    for i:=0; i<len(event.Data.UserAll_Id); i++{
      for j:=0; j<len(users.Data); j++{
        found = false
        if len(event.Data.UserSelected_Id) > 0{
          if users.Data[j].User_Id == event.Data.UserAll_Id[i].Value{
            if event.Data.UserAll_Id[i].Selected == true{
              for k:=0; k<len(users.Data[j].Event_Id); k++{
                if users.Data[j].Event_Id[k] == event.Data.Event_Id{
                  found = true
                }
              }
              if found == false{
                users.Data[j].Event_Id = append(users.Data[j].Event_Id, event.Data.Event_Id)
                usbody := UserResource{}
                usbody.Data.Id = users.Data[j].Id
                usbody.Data.First_name = users.Data[j].First_name
                usbody.Data.Last_name = users.Data[j].Last_name
                usbody.Data.Status = users.Data[j].Status
                usbody.Data.Approved = users.Data[j].Approved
                usbody.Data.Email = users.Data[j].Email
                usbody.Data.Password = users.Data[j].Password
                usbody.Data.User_Id = users.Data[j].User_Id
                usbody.Data.Role = users.Data[j].Role
                usbody.Data.Event_Id = users.Data[j].Event_Id
                err = usRepo.Update(&usbody.Data)
              }
            }else if event.Data.UserAll_Id[i].Selected == false{
              for k:=0; k<len(users.Data[j].Event_Id); k++{
                if users.Data[j].Event_Id[k] == event.Data.Event_Id{
                  k1 := 0
                  k2 := 0
                  if k > 0{
                    k1 = k - 1
                  }else if k == 0{
                    k1 = 1
                  }
                  alength := len(users.Data[j].Event_Id) - 1
                  if k < alength{
                    k2 = k + 1
                  }
                  if k == 0{
                    users.Data[j].Event_Id = users.Data[j].Event_Id[k1:]
                  }else if k > 0 && k < alength{
                    users.Data[j].Event_Id = append(users.Data[j].Event_Id[:k1], users.Data[j].Event_Id[k2:]...)
                  }else if k == alength{
                    users.Data[j].Event_Id = users.Data[j].Event_Id[:alength]
                  }
                  usbody := UserResource{}
                  usbody.Data.Id = users.Data[j].Id
                  usbody.Data.First_name = users.Data[j].First_name
                  usbody.Data.Last_name = users.Data[j].Last_name
                  usbody.Data.Status = users.Data[j].Status
                  usbody.Data.Approved = users.Data[j].Approved
                  usbody.Data.Email = users.Data[j].Email
                  usbody.Data.Password = users.Data[j].Password
                  usbody.Data.Role = users.Data[j].Role
                  usbody.Data.User_Id = users.Data[j].User_Id
                  usbody.Data.Event_Id = users.Data[j].Event_Id
                  err = usRepo.Update(&usbody.Data)
                }
              }
            }
          }
        }
        if len(event.Data.UserSelected_Id) <= 0{
          for k:=0; k<len(users.Data[j].Event_Id); k++{
            if users.Data[j].Event_Id[k] == event.Data.Event_Id{
              k1 := 0
              k2 := 0
              if k > 0{
                k1 = k - 1
              }else if k == 0{
                k1 = 1
              }
              alength := len(users.Data[j].Event_Id) - 1
              if k < alength{
                k2 = k + 1
              }
              if k == 0{
                users.Data[j].Event_Id = users.Data[j].Event_Id[k1:]
              }else if k > 0 && k < alength{
                users.Data[j].Event_Id = append(users.Data[j].Event_Id[:k1], users.Data[j].Event_Id[k2:]...)
              }else if k == alength{
                users.Data[j].Event_Id = users.Data[j].Event_Id[:alength]
              }
              usbody := UserResource{}
              usbody.Data.Id = users.Data[j].Id
              usbody.Data.First_name = users.Data[j].First_name
              usbody.Data.Last_name = users.Data[j].Last_name
              usbody.Data.Status = users.Data[j].Status
              usbody.Data.Approved = users.Data[j].Approved
              usbody.Data.Email = users.Data[j].Email
              usbody.Data.Password = users.Data[j].Password
              usbody.Data.Role = users.Data[j].Role
              usbody.Data.User_Id = users.Data[j].User_Id
              usbody.Data.Event_Id = users.Data[j].Event_Id
              err = usRepo.Update(&usbody.Data)
            }
          }
        }
      }
    }
    if err != nil {
      panic(err)
    }
    if event.Data.Event_Id == ""{
      http.Redirect(w, r, "/events/delete/" + event.Data.Id.Hex(), 302)
    }else{
      http.Redirect(w, r, "/events/show/" + event.Data.Id.Hex(), 302)
    }
  }
}

func (c *appContext) deleteEventHandler(w http.ResponseWriter, r *http.Request) {
  if current_session.Current_status == false{
    http.Redirect(w, r, "/login", 302)
  }else{
    params := context.Get(r, "params").(httprouter.Params)    //gorilla context, key "params"
    evRepo := EventRepo{c.db.C("events")}
    event, err := evRepo.Find(params.ByName("id"))

    enRepo := EntrantRepo{c.db.C("entrants")}
    entrants, err := enRepo.All()

    usRepo := UserRepo{c.db.C("users")}
    users, err := usRepo.All()

    scRepo := ScorecardRepo{c.db.C("scorecards")}
    scorecards, err := scRepo.All()

    tlyRepo := TallyRepo{c.db.C("tallies")}
    tallies, err := tlyRepo.All()

    fmt.Println(err)

    for i:=0; i<len(event.Data.EntrantSelected_Id); i++{
      for j:=0; j<len(entrants.Data); j++{
        if len(event.Data.EntrantSelected_Id) > 0{
          for k:=0; k<len(event.Data.EntrantSelected_Id); k++{
            if entrants.Data[j].Team_Id == event.Data.EntrantSelected_Id[k]{
              for m:=0; m<len(tallies.Data); m++{

                //delete tally
                if tallies.Data[m].Entrant_Id == entrants.Data[j].Team_Id && tallies.Data[m].Event_Id == event.Data.Event_Id{
                  err := tlyRepo.Delete(tallies.Data[m].Id.Hex())
                  fmt.Println(err)
                }
              }
              for m:=0; m<len(scorecards.Data); m++{
                if scorecards.Data[m].Entrant_Id == entrants.Data[j].Team_Id && scorecards.Data[m].Event_Id == event.Data.Event_Id{
                  err := scRepo.Delete(scorecards.Data[m].Id.Hex())
                  fmt.Println(err)
                }
              }
              for m:=0; m<len(entrants.Data[j].Event_Id); m++{
                if entrants.Data[j].Event_Id[m] == event.Data.Event_Id{
                  kdn := 0
                  kup := 0
                  if m > 0{
                    kdn = m - 1
                  }else{
                    kdn = 0
                  }
                  alength := len(entrants.Data[j].Event_Id) - 1
                  if m < alength{
                    kup = m
                  }else{
                    kup = alength
                  }
                  entrants.Data[j].Event_Id = append(entrants.Data[j].Event_Id[:kdn], entrants.Data[j].Event_Id[kup:alength]...)
                  enbody := EntrantResource{}
                  enbody.Data.Id = entrants.Data[j].Id
                  enbody.Data.First_name = entrants.Data[j].First_name
                  enbody.Data.Last_name = entrants.Data[j].Last_name
                  enbody.Data.Id_number = entrants.Data[j].Id_number
                  enbody.Data.Dog_name = entrants.Data[j].Dog_name
                  enbody.Data.Dog_id_number = entrants.Data[j].Dog_id_number
                  enbody.Data.Breed = entrants.Data[j].Breed
                  enbody.Data.Team_Id = entrants.Data[j].Team_Id
                  enbody.Data.Event_Id = entrants.Data[j].Event_Id
                  err = enRepo.Update(&enbody.Data)
                }
              }
            }
          }
        }
      }
    }
    for i:=0; i<len(event.Data.UserSelected_Id); i++{
      for j:=0; j<len(users.Data); j++{
        if len(event.Data.UserSelected_Id) > 0{
          for k:=0; k<len(event.Data.UserSelected_Id); k++{
            if users.Data[j].User_Id == event.Data.UserSelected_Id[k]{
              for m:=0; m<len(users.Data[j].Event_Id); m++{
                if users.Data[j].Event_Id[m] == event.Data.Event_Id{
                  kdn := 0
                  kup := 0
                  if m > 0{
                    kdn = m - 1
                  }else{
                    kdn = 0
                  }
                  alength := len(users.Data[j].Event_Id) - 1
                  if m < alength{
                    kup = m
                  }else{
                    kup = alength
                  }
                  users.Data[j].Event_Id = append(users.Data[j].Event_Id[:kdn], users.Data[j].Event_Id[kup:alength]...)
                  usbody := UserResource{}
                  usbody.Data.Id = users.Data[j].Id
                  usbody.Data.First_name = users.Data[j].First_name
                  usbody.Data.Last_name = users.Data[j].Last_name
                  usbody.Data.Status = users.Data[j].Status
                  usbody.Data.Approved = users.Data[j].Approved
                  usbody.Data.Email = users.Data[j].Email
                  usbody.Data.Password = users.Data[j].Password
                  usbody.Data.Role = users.Data[j].Role
                  usbody.Data.User_Id = users.Data[j].User_Id
                  usbody.Data.Event_Id = users.Data[j].Event_Id
                  err = usRepo.Update(&usbody.Data)
                }
              }
            }
          }
        }
      }
    }
    err = evRepo.Delete(params.ByName("id"))
    if err != nil {
      panic(err)
    }
    http.Redirect(w, r, "/events", 302)
  }
}

func (c *appContext) place_order(id string) []string{
  evRepo := EventRepo{c.db.C("events")}
  event, err := evRepo.Find(id)
  tlyRepo := TallyRepo{c.db.C("tallies")}
  tallies, err := tlyRepo.All()
  place_points := 0
  place_faults := 0
  place_time := 0
  fcount := 0
  tcount := 0
  found := false
  for i:=0; i<len(tallies.Data); i++{
    if tallies.Data[i].Event_Id == event.Data.Event_Id{
      fcount += 1
    }
  }
  placing := make([]string, fcount)

  for i:=0; i<len(placing); i++{
    found = false
    place_points = 0
    place_faults = 0
    place_time = 0
    total_points := 0
    total_points_flt := 0.0
    total_faults := 0
    total_faults_flt := 0.0
    tmp_time := 0
    for j:=0; j<len(tallies.Data); j++{
      if tallies.Data[j].Event_Id == event.Data.Event_Id{
        for k:=0; k<len(placing); k++{
          if placing[k] == tallies.Data[j].Tally_Id{
            found = true
          }
        }
        if found{
          found = false
          continue
        }
        tmp_time = str_to_time(tallies.Data[j].Total_time)
        if tmp_time != 0{
          total_points_flt, err = strconv.ParseFloat(tallies.Data[j].Total_points, 64)
          total_points = int(total_points_flt)
        }else{
          total_points = 0
        }
        if tmp_time != 0{
          total_faults_flt, err = strconv.ParseFloat(tallies.Data[j].Total_faults, 64)
          total_faults = int(total_faults_flt)
        }else{
          total_faults = 0
        }
        if place_points < total_points{
          place_time = str_to_time(tallies.Data[j].Total_time)
          place_points = total_points
          place_faults = total_faults
        }
        if place_points <= total_points && place_time >= str_to_time(tallies.Data[j].Total_time) && place_faults >= total_faults{
          place_time = str_to_time(tallies.Data[j].Total_time)
          place_points = total_points
          place_faults = total_faults
        }
      }
    }
    for j:=0; j<len(tallies.Data); j++{
      if tallies.Data[j].Event_Id == event.Data.Event_Id{
        for k:=0; k<len(placing); k++{
          if placing[k] == tallies.Data[j].Tally_Id{
            found = true
          }
        }
        if found{
          found = false
          continue
        }
        if tmp_time != 0{
          total_points_flt, err = strconv.ParseFloat(tallies.Data[j].Total_points, 64)
          total_points = int(total_points_flt)
        }else{
          total_points = 0
        }
        if tmp_time != 0{
          total_faults_flt, err = strconv.ParseFloat(tallies.Data[j].Total_faults, 64)
          total_faults = int(total_faults_flt)
        }else{
          total_faults = 0
        }
        if (total_points == place_points) && (str_to_time(tallies.Data[j].Total_time) == place_time) && (total_faults == place_faults){
          placing[tcount] = tallies.Data[j].Tally_Id
          tcount += 1
          break
        }
        fmt.Println(err)
      }
    }
  }
  return placing
}

func (c *appContext) scorecard_completion(id string) []string{
  evRepo := EventRepo{c.db.C("events")}
  event, err := evRepo.Find(id)
  scRepo := ScorecardRepo{c.db.C("scorecards")}
  scorecards, err := scRepo.All()
  completed_entrant_scorecards := make([]string, len(event.Data.EntrantSelected_Id))
  var cont_count int
  var int_count int
  var ext_count int
  var veh_count int
  var elite_count int
  sc_count := 0
  for i:=0; i<len(completed_entrant_scorecards); i++{
    completed_entrant_scorecards[i] = "inc"
  }
  for i:=0; i<len(completed_entrant_scorecards); i++{
    sc_count = 0
    for m:=0; m<len(ELEMENTS); m++{
      switch ELEMENTS[m]{
        case "Container":
          if event.Data.Cont_search_areas != ""{
            cont_count, err = strconv.Atoi(event.Data.Cont_search_areas)
          }else{
            cont_count = 0
          }
        case "Exterior":
          if event.Data.Ext_search_areas != ""{
            ext_count, err = strconv.Atoi(event.Data.Ext_search_areas)
          }else{
            ext_count = 0
          }
        case "Interior":
          if event.Data.Int_search_areas != ""{
            int_count, err = strconv.Atoi(event.Data.Int_search_areas)
          }else{
            int_count = 0
          }
        case "Vehicle":
          if event.Data.Veh_search_areas != ""{
            veh_count, err = strconv.Atoi(event.Data.Veh_search_areas)
          }else{
            veh_count = 0
          }
        case "Elite":
          if event.Data.Elite_search_areas != ""{
            elite_count, err = strconv.Atoi(event.Data.Elite_search_areas)
          }else{
            elite_count = 0
          }
      }
      for n:=0; n<len(scorecards.Data); n++{
        if (scorecards.Data[n].Entrant_Id == event.Data.EntrantSelected_Id[i]) && (scorecards.Data[n].Element == ELEMENTS[m]) && (scorecards.Data[n].Event_Id == event.Data.Event_Id){
          if scorecards.Data[n].Judge_signature.Value == "yes"{
            sc_count += 1
          }
        }
      }
    }
    fmt.Println(err)
    if sc_count == cont_count + ext_count + int_count + veh_count + elite_count{
      completed_entrant_scorecards[i] = event.Data.EntrantSelected_Id[i]
    }
  }
  fmt.Println(err)
  return completed_entrant_scorecards
}

func (c *appContext) tally_completion(id string) []string{
  evRepo := EventRepo{c.db.C("events")}
  event, err := evRepo.Find(id)
  tlyRepo := TallyRepo{c.db.C("tallies")}
  tallies, err := tlyRepo.All()
  completed_entrant_tallies := make([]string, len(event.Data.EntrantSelected_Id))
  for i:=0; i<len(completed_entrant_tallies); i++{
    completed_entrant_tallies[i] = "inc"
  }
  for i:=0; i<len(completed_entrant_tallies); i++{
    for k:=0; k<len(tallies.Data); k++ {
      if (tallies.Data[k].Entrant_Id == event.Data.EntrantSelected_Id[i]) && (tallies.Data[k].Event_Id == event.Data.Event_Id){
        if tallies.Data[k].Total_time != "0"{
          completed_entrant_tallies[i] = event.Data.EntrantSelected_Id[i]
        }
      }
    }
  }
  fmt.Println(err)
  return completed_entrant_tallies
}


// Entrant Handlers /////////////////////////////////////////////////////////////////////////////////////

func (c *appContext) entrantsHandler(w http.ResponseWriter, r *http.Request) {
  if current_session.Current_status == false{
    http.Redirect(w, r, "/login", 302)
  }else{
    entrants_resrc := EntrantsResource{}
    entrants_resrc.SData = current_session
    repo := EntrantRepo{c.db.C("entrants")}
    entrants, err := repo.All()
    for i:=0; i<len(entrants.Data); i++{
      body := EntrantResource{}
      body.Data.Id = entrants.Data[i].Id
      body.Data.Event_Id = entrants.Data[i].Event_Id
      body.Data.First_name = entrants.Data[i].First_name
      body.Data.Last_name = entrants.Data[i].Last_name
      body.Data.Team_Id = entrants.Data[i].Team_Id
      body.Data.Breed = entrants.Data[i].Breed
      body.Data.Dog_name = entrants.Data[i].Dog_name
      body.Data.Dog_id_number = entrants.Data[i].Dog_id_number
      body.Data.Id_number = entrants.Data[i].Id_number
      entrants_resrc.Data = append(entrants_resrc.Data, body)
    }
    if err != nil {
      panic(err)
    }
    if err = listEntrant.Execute(w, entrants_resrc); err != nil {
      http.Error(w, err.Error(), http.StatusInternalServerError)
      return
    }
  }
}

func (c *appContext) entrantHandler(w http.ResponseWriter, r *http.Request) {
  if current_session.Current_status == false{
    http.Redirect(w, r, "/login", 302)
  }else{
    params := context.Get(r, "params").(httprouter.Params)  //gorrila context, key "params"
    repo := EntrantRepo{c.db.C("entrants")}
    entrant, err := repo.Find(params.ByName("id")) //getting data from named param :id
    entrant_resrc := EntrantResource{}
    entrant_resrc.Data = entrant.Data
    entrant_resrc.SData = current_session
    if err != nil {
      panic(err)
    }
    if err = showEntrant.Execute(w, entrant_resrc); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
  }
}

func (c *appContext) newEntrantHandler(w http.ResponseWriter, r *http.Request) {
  if current_session.Current_status == false{
    http.Redirect(w, r, "/login", 302)
  }else{
    entrantresrc := EntrantResource{}
    entrantresrc.SData = current_session

    if err := createnewEntrant.Execute(w, entrantresrc); err != nil {
      http.Error(w, err.Error(), http.StatusInternalServerError)
      return
    }
  }
  // forwards to createEntrantHandler
}

func (c *appContext) createEntrantHandler(w http.ResponseWriter, r *http.Request) {
  if current_session.Current_status == false{
    http.Redirect(w, r, "/login", 302)
  }else{
    repo := EntrantRepo{c.db.C("entrants")}
    entrants, err := repo.All()
    rrcount := 0
    body := context.Get(r, "body").(*EntrantResource)    //gorilla context, key "body" that returns val
    body.Data.First_name = r.FormValue("First_name")
    body.Data.Last_name = r.FormValue("Last_name")
    body.Data.Id_number = "M_" + strconv.Itoa(rand.Int())

    // check for duplicates
    for r:=0; r<len(entrants.Data); r++{
      if body.Data.Id_number == entrants.Data[r].Id_number{
        // Id_number duplicate found - re-naming and re-checking loop 1 to loop 2
        body.Data.Id_number = "M_" + strconv.Itoa(rand.Int())
      }
      for rr:=0;rr<len(entrants.Data); rr++{
        if body.Data.Id_number == entrants.Data[rr].Id_number{
          // Id_number duplicate found - re-naming and re-checking loop 2 to outer loop 1
          body.Data.Id_number = "M_" + strconv.Itoa(rand.Int())
          break
        }else{
          rrcount = rr
        }
      }
      if rrcount == len(entrants.Data)-1{
        // No duplicates both loops
        break
      }
    }
    body.Data.Dog_name = r.FormValue("Dog_name")
    body.Data.Dog_id_number = "K_" + strconv.Itoa(rand.Int())

    // check for duplicates
    rrcount = 0
    for r:=0; r<len(entrants.Data); r++{
      if body.Data.Dog_id_number == entrants.Data[r].Dog_id_number{
        // Dog_id_number duplicate found - re-naming and re-checking loop 1 to loop 2
        body.Data.Dog_id_number = "K_" + strconv.Itoa(rand.Int())
      }
      for rr:=0;rr<len(entrants.Data); rr++{
        if body.Data.Dog_id_number == entrants.Data[rr].Dog_id_number{
          // Dog_id_number duplicate found - re-naming and re-checking loop 2 to outer loop 1
          body.Data.Dog_id_number = "K_" + strconv.Itoa(rand.Int())
          break
        }else{
          rrcount = rr
        }
      }
      if rrcount == len(entrants.Data)-1{
        // No duplicates both loops
        break
      }
    }
    body.Data.Breed = r.FormValue("Breed")
    body.Data.Team_Id = "TM_" + strconv.Itoa(rand.Int())

    // check for duplicates
    rrcount = 0
    for r:=0; r<len(entrants.Data); r++{
      if body.Data.Team_Id == entrants.Data[r].Team_Id{
        // Team_Id duplicate found - re-naming and re-checking loop 1 to loop 2
        body.Data.Team_Id = "TM_" + strconv.Itoa(rand.Int())
      }
      for rr:=0;rr<len(entrants.Data); rr++{
        if body.Data.Team_Id == entrants.Data[rr].Team_Id{
          // Team_Id duplicate found - re-naming and re-checking loop 2 to outer loop 1
          body.Data.Team_Id = "TM_" + strconv.Itoa(rand.Int())
          break
        }else{
          rrcount = rr
        }
      }
      if rrcount == len(entrants.Data)-1{
        // No duplicates both loops
        break
      }
    }
    err, id := repo.Create(&body.Data)
    if err != nil {
      panic(err)
    }
    http.Redirect(w, r, "/entrants/show/" + id.Hex(), 302)
  }
}

func (c *appContext) editEntrantHandler(w http.ResponseWriter, r *http.Request) {
  if current_session.Current_status == false{
    http.Redirect(w, r, "/login", 302)
  }else{
    params := context.Get(r, "params").(httprouter.Params)    //gorilla context, key "params"
    repo := EntrantRepo{c.db.C("entrants")}
    entrant, err := repo.Find(params.ByName("id")) //getting data from named param :id
    entrant_resrc := EntrantResource{}
    entrant_resrc.Data = entrant.Data
    entrant_resrc.SData = current_session
    if err = updateEntrant.Execute(w, entrant_resrc); err != nil {
      http.Error(w, err.Error(), http.StatusInternalServerError)
      return
    }
  }
}

func (c *appContext) updateEntrantHandler(w http.ResponseWriter, r *http.Request) {
  if current_session.Current_status == false{
    http.Redirect(w, r, "/login", 302)
  }else{
    params := context.Get(r, "params").(httprouter.Params)
    repo := EntrantRepo{c.db.C("entrants")}
    entrant, err := repo.Find(params.ByName("id")) //getting data from named param :id
    body := context.Get(r, "body").(*EntrantResource)
    body.Data.Id = entrant.Data.Id
    body.Data.First_name = r.FormValue("First_name")
    body.Data.Last_name = r.FormValue("Last_name")
    body.Data.Id_number = entrant.Data.Id_number
    body.Data.Dog_name = r.FormValue("Dog_name")
    body.Data.Dog_id_number = entrant.Data.Dog_id_number
    body.Data.Breed = r.FormValue("Breed")
    body.Data.Team_Id = entrant.Data.Team_Id
    body.Data.Event_Id = entrant.Data.Event_Id
    err = repo.Update(&body.Data)
    if err != nil {
      panic(err)
    }
    http.Redirect(w, r, "/entrants/show/" + body.Data.Id.Hex(), 302)
  }
}

func (c *appContext) deleteEntrantHandler(w http.ResponseWriter, r *http.Request) {
  if current_session.Current_status == false{
    http.Redirect(w, r, "/login", 302)
  }else{
    params := context.Get(r, "params").(httprouter.Params)    //gorilla context, key "params"
    repo := EntrantRepo{c.db.C("entrants")}
    entrant, err := repo.Find(params.ByName("id"))
    if len(entrant.Data.Event_Id) > 0{
        fmt.Println("Cannot delete team that is part of an event")
    }else{
      err := repo.Delete(params.ByName("id"))
      fmt.Println(err)
    }
    if err != nil {
      panic(err)
    }
    http.Redirect(w, r, "/entrants", 302)
  }
}


// User Handlers /////////////////////////////////////////////////////////////////////////////////////


func (c *appContext) usersHandler(w http.ResponseWriter, r *http.Request) {
  if current_session.Current_status == false{
    http.Redirect(w, r, "/login", 302)
  }else{
    repo := UserRepo{c.db.C("users")}
    users, err := repo.All()
    users_resrc := UsersResource{}
    users_resrc.SData = current_session
    for j:=0; j<len(users.Data); j++{
      body := UserResource{}
      body.Data.Id = users.Data[j].Id
      body.Data.First_name = users.Data[j].First_name
      body.Data.Last_name = users.Data[j].Last_name
      body.Data.User_Id = users.Data[j].User_Id
      body.Data.Email = users.Data[j].Email
      body.Data.Password = users.Data[j].Password
      body.Data.Role = users.Data[j].Role
      body.Data.Status = users.Data[j].Status
      users_resrc.Data = append(users_resrc.Data, body)
    }
    if err != nil {
      panic(err)
    }
    if err = listUser.Execute(w, users_resrc); err != nil {
      http.Error(w, err.Error(), http.StatusInternalServerError)
      return
    }
  }
}

func (c *appContext) userHandler(w http.ResponseWriter, r *http.Request) {
  if current_session.Current_status == false{
    http.Redirect(w, r, "/login", 302)
  }else{
    params := context.Get(r, "params").(httprouter.Params)  //gorrila context, key "params"
    repo := UserRepo{c.db.C("users")}
    user, err := repo.Find(params.ByName("id")) //getting data from named param :id
    user_resrc := UserResource{}
    user_resrc.Data = user.Data
    user_resrc.SData = current_session
    if err != nil {
      panic(err)
    }
    if err = showUser.Execute(w, user_resrc); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
  }
}

func (c *appContext) newUserHandler(w http.ResponseWriter, r *http.Request) {
  if current_session.Current_status == false{
    http.Redirect(w, r, "/login", 302)
  }else{
    usrresrc := UserResource{}
    usrresrc.SData = current_session
    if err := createnewUser.Execute(w, usrresrc); err != nil {
      http.Error(w, err.Error(), http.StatusInternalServerError)
      return
    }
  }
  // forwards to createUserHandler
}

func (c *appContext) createUserHandler(w http.ResponseWriter, r *http.Request) {
  if current_session.Current_status == false{
    http.Redirect(w, r, "/login", 302)
  }else{
    repo := UserRepo{c.db.C("users")}
    users, err := repo.All()
    rrcount := 0
    body := context.Get(r, "body").(*UserResource)    //gorilla context, key "body" that returns val
    body.Data.First_name = r.FormValue("First_name")
    body.Data.Last_name = r.FormValue("Last_name")
    body.Data.User_Id = "US_" + strconv.Itoa(rand.Int())
    // check for duplicates
    for r:=0; r<len(users.Data); r++{
      if body.Data.User_Id == users.Data[r].User_Id{
        // User_Id duplicate found - re-naming and re-checking loop 1 to loop 2
        body.Data.User_Id = "US_" + strconv.Itoa(rand.Int())
      }
      for rr:=0;rr<len(users.Data); rr++{
        if body.Data.User_Id == users.Data[rr].User_Id{
          // User_Id duplicate found - re-naming and re-checking loop 2 to outer loop 1
          body.Data.User_Id = "US_" + strconv.Itoa(rand.Int())
          break
        }else{
          rrcount = rr
        }
      }
      if rrcount == len(users.Data)-1{
        // No duplicates both loops
        break
      }
    }
    body.Data.Status = r.FormValue("Status")
    body.Data.Role = r.FormValue("Role")
    body.Data.Email = r.FormValue("Email")
    body.Data.Password = r.FormValue("Password")
    err, id := repo.Create(&body.Data)
    if err != nil {
      panic(err)
    }
    http.Redirect(w, r, "/users/show/" + id.Hex(), 302)
  }
}

func (c *appContext) editUserHandler(w http.ResponseWriter, r *http.Request) {
  if current_session.Current_status == false{
    http.Redirect(w, r, "/login", 302)
  }else{
    params := context.Get(r, "params").(httprouter.Params)    //gorilla context, key "params"
    repo := UserRepo{c.db.C("users")}
    user, err := repo.Find(params.ByName("id")) //getting data from named param :id
    userresrc := UserResource{}
    userresrc.Data = user.Data
    userresrc.SData = current_session
    if err = updateUser.Execute(w, userresrc); err != nil {
      http.Error(w, err.Error(), http.StatusInternalServerError)
      return
    }
    // forwards to updateUserHandler
  }
}

func (c *appContext) updateUserHandler(w http.ResponseWriter, r *http.Request) {
  if current_session.Current_status == false{
    http.Redirect(w, r, "/login", 302)
  }else{
    params := context.Get(r, "params").(httprouter.Params)
    body := context.Get(r, "body").(*UserResource)
    repo := UserRepo{c.db.C("users")}
    user, err := repo.Find(params.ByName("id")) //getting data from named param :id
    body.Data.Id = user.Data.Id
    body.Data.First_name = r.FormValue("First_name")
    body.Data.Last_name = r.FormValue("Last_name")
    body.Data.User_Id = user.Data.User_Id
    body.Data.Email = r.FormValue("Email")
    body.Data.Password = r.FormValue("Password")
    body.Data.Role = r.FormValue("Role")
    body.Data.Status = r.FormValue("Status")
    body.Data.Event_Id = user.Data.Event_Id
    err = repo.Update(&body.Data)
    if err != nil {
      panic(err)
    }
    http.Redirect(w, r, "/users/show/" + body.Data.Id.Hex(), 302)
  }
}

func (c *appContext) deleteUserHandler(w http.ResponseWriter, r *http.Request) {
  if current_session.Current_status == false{
    http.Redirect(w, r, "/login", 302)
  }else{
    params := context.Get(r, "params").(httprouter.Params)    //gorilla context, key "params"
    repo := UserRepo{c.db.C("users")}
    user, err := repo.Find(params.ByName("id"))
    if (len(user.Data.Event_Id) > 0) || (user.Data.User_Id == current_session.Current_user){
        fmt.Println("Cannot delete user that is part of an event or yourself")
    }else{
      err := repo.Delete(params.ByName("id"))
      fmt.Println(err)
    }
    if err != nil {
      panic(err)
    }
    http.Redirect(w, r, "/users", 302)
  }
}


// Scorecard Handlers /////////////////////////////////////////////////////////////////////////////////////


func (c *appContext) scorecardsHandler(w http.ResponseWriter, r *http.Request) {
  if current_session.Current_status == false{
    http.Redirect(w, r, "/login", 302)
  }else{
    repo := ScorecardRepo{c.db.C("scorecards")}
    scorecards, err := repo.All()
    scorecards_resrc := ScorecardsResource{}
    scorecards_resrc.SData = current_session
    for j:=0; j<len(scorecards.Data); j++{
      body := ScorecardResource{}
      body.Data.Id = scorecards.Data[j].Id
      body.Data.Scorecard_Id = scorecards.Data[j].Scorecard_Id
      body.Data.Event_Id = scorecards.Data[j].Event_Id
      body.Data.Entrant_Id = scorecards.Data[j].Entrant_Id
      body.Data.Search_area = scorecards.Data[j].Search_area
      body.Data.Element = scorecards.Data[j].Element
      body.Data.Total_points = scorecards.Data[j].Total_points
      body.Data.Total_faults = scorecards.Data[j].Total_faults
      body.Data.Total_time = scorecards.Data[j].Total_time
      scorecards_resrc.Data = append(scorecards_resrc.Data, body)
    }
    if err != nil {
      panic(err)
    }
    //	w.Header().Set("Content-Type", "application/vnd.api+json")
    //	json.NewEncoder(w).Encode(scorecards)
    // read BSON into JSON

    if err = listScorecard.Execute(w, scorecards_resrc); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
  }
}

func (c *appContext) scorecardHandler(w http.ResponseWriter, r *http.Request) {
  if current_session.Current_status == false{
    http.Redirect(w, r, "/login", 302)
  }else{
    params := context.Get(r, "params").(httprouter.Params)  //gorrila context, key "params"
    repo := ScorecardRepo{c.db.C("scorecards")}
    scorecard, err := repo.Find(params.ByName("id")) //getting data from named param :id
    if err != nil {
      panic(err)
    }
    scorecard_resrc := ScorecardFormResource{}
    scorecard_resrc.SCData = scorecard
    evRepo := EventRepo{c.db.C("events")}
    events, err := evRepo.All()
    evbody := EventResource{}
    for i:=0; i<len(events.Data); i++{
      if events.Data[i].Event_Id == scorecard.Data.Event_Id{
        evbody.Data = events.Data[i]
      }
    }
    scorecard_resrc.EVData = evbody
    enRepo := EntrantRepo{c.db.C("entrants")}
    entrants, err := enRepo.All()
    enbody := EntrantResource{}
    for i:=0; i<len(entrants.Data); i++{
      if entrants.Data[i].Team_Id == scorecard.Data.Entrant_Id{
        for j:=0; j<len(entrants.Data[i].Event_Id); j++{
          if entrants.Data[i].Event_Id[j] == scorecard.Data.Event_Id{
            enbody.Data = entrants.Data[i]
          }
        }
      }
    }
    scorecard_resrc.ENData = enbody
    scorecard_resrc.SData = current_session

    //  w.Header().Set("Content-Type", "application/vnd.api+json")
    //	json.NewEncoder(w).Encode(scorecard)
    //	if err = show.Execute(w, json.NewEncoder(w).Encode(scorecard)); err != nil {
    //      http.Error(w, err.Error(), http.StatusInternalServerError)
    //      return
    //  }

    // read JSON into BSON
    if err = showScorecard.Execute(w, scorecard_resrc); err != nil {
      http.Error(w, err.Error(), http.StatusInternalServerError)
      return
    }
  }
}

func (c *appContext) editScorecardHandler(w http.ResponseWriter, r *http.Request) {
  if current_session.Current_status == false{
    http.Redirect(w, r, "/login", 302)
  }else{
    params := context.Get(r, "params").(httprouter.Params)    //gorilla context, key "params"
    repo := ScorecardRepo{c.db.C("scorecards")}
    scorecard, err := repo.Find(params.ByName("id")) //getting data from named param :id
    if scorecard.Data.Finish_call.Value == "yes"{
      scorecard.Data.Finish_call.Selected = true
    }else{
      scorecard.Data.Finish_call.Selected = false
    }
    if scorecard.Data.Timed_out.Value == "yes"{
      scorecard.Data.Timed_out.Selected = true
    }else{
      scorecard.Data.Timed_out.Selected = false
    }
    if scorecard.Data.Dismissed.Value == "yes"{
      scorecard.Data.Dismissed.Selected = true
    }else{
      scorecard.Data.Dismissed.Selected = false
    }
    if scorecard.Data.Excused.Value == "yes"{
      scorecard.Data.Excused.Selected = true
    }else{
      scorecard.Data.Excused.Selected = false
    }
    if scorecard.Data.Absent.Value == "yes"{
      scorecard.Data.Absent.Selected = true
    }else{
      scorecard.Data.Absent.Selected = false
    }
    if scorecard.Data.Eliminated_during_search.Value == "yes"{
      scorecard.Data.Eliminated_during_search.Selected = true
    }else{
      scorecard.Data.Eliminated_during_search.Selected = false
    }
    if scorecard.Data.Pronounced.Value == "yes"{
      scorecard.Data.Pronounced.Selected = true
    }else{
      scorecard.Data.Pronounced.Selected = false
    }
    if scorecard.Data.Judge_signature.Value == "yes"{
      scorecard.Data.Judge_signature.Selected = true
    }else{
      scorecard.Data.Judge_signature.Selected = false
    }
    scorecard_resrc := ScorecardFormResource{}
    scorecard_resrc.SCData = scorecard
    evRepo := EventRepo{c.db.C("events")}
    events, err := evRepo.All()
    evbody := EventResource{}
    for i:=0; i<len(events.Data); i++{
      if events.Data[i].Event_Id == scorecard.Data.Event_Id{
        evbody.Data = events.Data[i]
      }
    }
    scorecard_resrc.EVData = evbody

    enRepo := EntrantRepo{c.db.C("entrants")}
    entrants, err := enRepo.All()
    enbody := EntrantResource{}
    for i:=0; i<len(entrants.Data); i++{
      if entrants.Data[i].Team_Id == scorecard.Data.Entrant_Id{
        for j:=0; j<len(entrants.Data[i].Event_Id); j++{
          if entrants.Data[i].Event_Id[j] == scorecard.Data.Event_Id{
            enbody.Data = entrants.Data[i]
          }
        }
      }
    }
    scorecard_resrc.ENData = enbody
    scorecard_resrc.SData = current_session
    message := c.get_check_hide_count(scorecard.Data.Id.Hex())
    scorecard_resrc.CheckCount = message
    if err = updateScorecard.Execute(w, scorecard_resrc); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    if err != nil{
      fmt.Println(err)
    }
  }
}

func (c *appContext) updateScorecardHandler(w http.ResponseWriter, r *http.Request) {
  if current_session.Current_status == false{
    http.Redirect(w, r, "/login", 302)
  }else{
    params := context.Get(r, "params").(httprouter.Params)
    repo := ScorecardRepo{c.db.C("scorecards")}
    scorecard, err := repo.Find(params.ByName("id")) //getting data from named param :id
    body := context.Get(r, "body").(*ScorecardResource)
    body.Data.Id = scorecard.Data.Id
    body.Data.Element = scorecard.Data.Element
    body.Data.Maxtime_m = r.FormValue("Maxtime_m")
    body.Data.Maxtime_s = r.FormValue("Maxtime_s")
    body.Data.Finish_call.Value = r.FormValue("Finish_call")
    body.Data.False_alert_fringe = r.FormValue("False_alert_fringe")
    body.Data.Timed_out.Value = r.FormValue("Timed_out")
    body.Data.Dismissed.Value = r.FormValue("Dismissed")
    body.Data.Excused.Value = r.FormValue("Excused")
    body.Data.Absent.Value = r.FormValue("Absent")
    body.Data.Eliminated_during_search.Value = r.FormValue("Eliminated_during_search")
    body.Data.Other_faults_descr = r.FormValue("Other_faults_descr")
    body.Data.Other_faults_count = r.FormValue("Other_faults_count")
    body.Data.Comments = r.FormValue("Comments")

    if scorecard.Data.Total_time != "00:00:00"{
         body.Data.Total_time = r.FormValue("Total_time")
    }else if scorecard.Data.Total_time == "" {
         body.Data.Total_time = "00:00:00"
    }else{
         body.Data.Total_time = scorecard.Data.Total_time
    }

    body.Data.Pronounced.Value = r.FormValue("Pronounced")
    body.Data.Judge_signature.Value = r.FormValue("Judge_signature")
    body.Data.Event_Id = scorecard.Data.Event_Id
    body.Data.Entrant_Id = scorecard.Data.Entrant_Id
    body.Data.Search_area = scorecard.Data.Search_area
    body.Data.Scorecard_Id = scorecard.Data.Scorecard_Id
    body.Data.Hides_max = r.FormValue("Hides_max")
    body.Data.Hides_found = r.FormValue("Hides_found")
    body.Data.Hides_missed = r.FormValue("Hides_missed")
    body.Data.Total_faults = c.get_fault_total(scorecard.Data.Id.Hex())
    body.Data.Maxpoint = c.get_max_point(scorecard.Data.Id.Hex())
    body.Data.Total_points = c.get_points(scorecard.Data.Id.Hex())
    if body.Data.Finish_call.Value == "yes"{
      body.Data.Finish_call.Selected = true
    }else{
      body.Data.Finish_call.Selected = false
    }
    if scorecard.Data.Timed_out.Value == "yes"{
      scorecard.Data.Timed_out.Selected = true
    }else{
      scorecard.Data.Timed_out.Selected = false
    }
    if scorecard.Data.Dismissed.Value == "yes"{
      scorecard.Data.Dismissed.Selected = true
    }else{
      scorecard.Data.Dismissed.Selected = false
    }
    if scorecard.Data.Excused.Value == "yes"{
      scorecard.Data.Excused.Selected = true
    }else{
      scorecard.Data.Excused.Selected = false
    }
    if scorecard.Data.Absent.Value == "yes"{
      scorecard.Data.Absent.Selected = true
    }else{
      scorecard.Data.Absent.Selected = false
    }
    if scorecard.Data.Eliminated_during_search.Value == "yes"{
      scorecard.Data.Eliminated_during_search.Selected = true
    }else{
      scorecard.Data.Eliminated_during_search.Selected = false
    }
    if scorecard.Data.Pronounced.Value == "yes"{
      scorecard.Data.Pronounced.Selected = true
    }else{
      scorecard.Data.Pronounced.Selected = false
    }
    if scorecard.Data.Judge_signature.Value == "yes"{
      scorecard.Data.Judge_signature.Selected = true
    }else{
      scorecard.Data.Judge_signature.Selected = false
    }
    tmp_time_m, err := strconv.Atoi(body.Data.Maxtime_m)
    tmp_time_m = tmp_time_m*60
    tmp_time_s, err := strconv.Atoi(body.Data.Maxtime_s)
    tmp_time := tmp_time_s + tmp_time_m
    tmp_timeD := time.Duration(tmp_time)*time.Second
    // global variable time limit
    timelimit = tmp_timeD

    err = repo.Update(&body.Data)

    if err != nil {
      panic(err)
    }

    if body.Data.Scorecard_Id == ""{
      http.Redirect(w, r, "/scorecards/delete/" + body.Data.Id.Hex(), 302)
    }else{
      http.Redirect(w, r, "/scorecards/edit/" + body.Data.Id.Hex(), 302)
    }
  }
}

func (c *appContext) deleteScorecardHandler(w http.ResponseWriter, r *http.Request) {
  if current_session.Current_status == false{
    http.Redirect(w, r, "/login", 302)
  }else{
    params := context.Get(r, "params").(httprouter.Params)    //gorilla context, key "params"
    repo := ScorecardRepo{c.db.C("scorecards")}
    err := repo.Delete(params.ByName("id"))
    if err != nil {
      panic(err)
    }

    http.Redirect(w, r, "/scorecards", 302)
  }
}

func (c *appContext) get_elmSearchAreas(id string) string {
  scRepo := ScorecardRepo{c.db.C("scorecards")}
  scorecard, err := scRepo.Find(id)
  evRepo := EventRepo{c.db.C("events")}
  events, err := evRepo.All()
  event := EventResource{}
  search_areas := ""
  for i:=0; i<len(events.Data); i++{
    if scorecard.Data.Event_Id == events.Data[i].Event_Id{
      event.Data = events.Data[i]
    }
  }
  switch scorecard.Data.Element{
    case "Container":
      search_areas = event.Data.Cont_search_areas
    case "Interior":
      search_areas = event.Data.Int_search_areas
    case "Exterior":
      search_areas = event.Data.Ext_search_areas
    case "Vehicle":
      search_areas = event.Data.Veh_search_areas
    case "Elite":
      search_areas = event.Data.Elite_search_areas
  }
  if err != nil {
	panic(err)
  }
  return search_areas
}

func (c *appContext) get_elmHides(id string) string {
  scRepo := ScorecardRepo{c.db.C("scorecards")}
  scorecard, err := scRepo.Find(id)
  evRepo := EventRepo{c.db.C("events")}
  events, err := evRepo.All()
  event := EventResource{}
  hides := ""
  for i:=0; i<len(events.Data); i++{
    if scorecard.Data.Event_Id == events.Data[i].Event_Id{
      event.Data = events.Data[i]
    }
  }
  switch scorecard.Data.Element{
    case "Container":
      hides = event.Data.Cont_hides
    case "Interior":
      hides = event.Data.Int_hides
    case "Exterior":
      hides = event.Data.Ext_hides
    case "Vehicle":
      hides = event.Data.Veh_hides
    case "Elite":
      hides =  event.Data.Elite_hides
  }
  if err != nil {
	panic(err)
  }
  return hides
}

func (c *appContext) get_check_hide_count(id string) string {
  scRepo := ScorecardRepo{c.db.C("scorecards")}
  scorecards, err := scRepo.All()
  scorecard, err := scRepo.Find(id)
  evRepo := EventRepo{c.db.C("events")}
  events, err := evRepo.All()
  event := EventResource{}
  enRepo := EntrantRepo{c.db.C("entrants")}
  entrants, err := enRepo.All()
  entrant := EntrantResource{}
  message := ""
  for i:=0; i<len(events.Data); i++{
    for j:=0; j<len(entrants.Data); j++{
      if scorecard.Data.Event_Id == events.Data[i].Event_Id && scorecard.Data.Entrant_Id == entrants.Data[j].Team_Id{
        event.Data = events.Data[i]
        entrant.Data = entrants.Data[j]
      }
    }
  }
  hideCountCheck, err := strconv.Atoi(c.get_elmHides(scorecard.Data.Id.Hex()))
  elm_hides := hideCountCheck
  if scorecard.Data.Hides_max != "" || scorecard.Data.Hides_max != "0"{
    for i:=0; i<len(entrants.Data); i++{
      if scorecard.Data.Entrant_Id == entrants.Data[i].Team_Id{
        for j:=0; j<len(scorecards.Data); j++{
          if (scorecards.Data[j].Entrant_Id == entrants.Data[i].Team_Id) && (scorecards.Data[j].Element == scorecard.Data.Element) && (scorecards.Data[j].Event_Id == scorecard.Data.Event_Id) && ((scorecards.Data[j].Hides_max != "") || (scorecards.Data[j].Hides_max != "0")){
            hidesMax, err := strconv.Atoi(scorecards.Data[j].Hides_max)
            hideCountCheck = hideCountCheck - hidesMax
            if err != nil {
              fmt.Println("err")
              panic(err)
            }
          }
        }
      }
    }
    message = ""
    if ((hideCountCheck > elm_hides) || (hideCountCheck < 0)) && (event.Data.Division != "NW1"){
      message =  "Incorrect Hide Count..."
    }
    if event.Data.Division == "NW1"{
      if scorecard.Data.Hides_max != "1"{
        message = "Incorrect Hide Count..."
      }
    }
    if event.Data.Division == "NW2"{
      hidesmax, err := strconv.Atoi(scorecard.Data.Hides_max)
      search_areas, err := strconv.Atoi(c.get_elmSearchAreas(scorecard.Data.Id.Hex()))
      fmt.Println(err)
      if scorecard.Data.Hides_max == "0" || ((hidesmax >= elm_hides) && (search_areas > 1)){
        message = "Incorrect Hide Count..."
      }
      if (scorecard.Data.Hides_max == "1") && (elm_hides == 1) && (search_areas == 1){
        message =""
      }
    }
  }
  if err != nil {
	panic(err)
  }
  return message
}

func (c *appContext) get_max_point(id string) string {
  scRepo := ScorecardRepo{c.db.C("scorecards")}
  scorecard, err := scRepo.Find(id)
  hidesMax, err := strconv.ParseFloat(scorecard.Data.Hides_max, 64)
  evRepo := EventRepo{c.db.C("events")}
  events, err := evRepo.All()
  event := EventResource{}
  for i:=0; i<len(events.Data); i++{
    if scorecard.Data.Event_Id == events.Data[i].Event_Id{
      event.Data = events.Data[i]
    }
  }
  point := 0.0
  if err != nil {
    panic(err)
  }
  if scorecard.Data.Maxpoint != "0"{
    if event.Data.Division != "NW1"{
      switch scorecard.Data.Element{
        case "Container":
          elementHides, err := strconv.ParseFloat(event.Data.Cont_hides, 64)
          if err != nil {
            panic(err)
          }
          if hidesMax > 0{
            if event.Data.Division != "Element Specialty"{
              point = (25.00/elementHides) * hidesMax
            }else if event.Data.Division == "Element Specialty"{
              point = (100.00/elementHides) * hidesMax
            }else{
              point = 0
            }
          }
        case "Interior":
          elementHides, err := strconv.ParseFloat(event.Data.Int_hides, 64)
          if err != nil {
            panic(err)
          }
          if hidesMax > 0{
            if event.Data.Division != "Element Specialty"{
              point = (25.00/elementHides) * hidesMax
            }else if event.Data.Division == "Element Specialty"{
              point = (100.00/elementHides) * hidesMax
            }else{
              point = 0
            }
          }
        case "Exterior":
          elementHides, err := strconv.ParseFloat(event.Data.Ext_hides, 64)
           if err != nil {
            panic(err)
          }
          if hidesMax > 0{
            if event.Data.Division != "Element Specialty"{
              point = (25.00/elementHides) * hidesMax
            }else if event.Data.Division == "Element Specialty"{
              point = (100.00/elementHides) * hidesMax
            }else{
              point = 0
            }
          }
        case "Vehicle":
          elementHides, err := strconv.ParseFloat(event.Data.Veh_hides, 64)
          if err != nil {
            panic(err)
          }
          if hidesMax > 0{
            if event.Data.Division != "Element Specialty"{
              point = (25.00/elementHides) * hidesMax
            }else if event.Data.Division == "Element Specialty"{
              point = (100.00/elementHides) * hidesMax
            }else{
              point = 0
            }
          }
        case "Elite":
          elementHides, err := strconv.ParseFloat(event.Data.Elite_hides, 64)
          if err != nil {
            panic(err)
          }
          if hidesMax > 0{
            point = (100.00/elementHides) * hidesMax
          }else{
            point = 0
          }
      }
    }else{
      point =  25.0
    }
  }
  pointStr := strconv.FormatFloat(point, 'f', 2, 64)
  if err != nil {
	panic(err)
  }
  return pointStr
}



func (c *appContext) get_fault_total(id string) string {
  scRepo := ScorecardRepo{c.db.C("scorecards")}
  scorecard, err := scRepo.Find(id)
  evRepo := EventRepo{c.db.C("events")}
  events, err := evRepo.All()
  event := EventResource{}
  totalfaults := 0
  falseAlertFringe := 0
  for i:=0; i<len(events.Data); i++{
    if scorecard.Data.Event_Id == events.Data[i].Event_Id{
      event.Data = events.Data[i]
    }
  }
  if scorecard.Data.Other_faults_count == ""{
    totalfaults = 0
  }else{
    totalfaults, err = strconv.Atoi(scorecard.Data.Other_faults_count)
  }
  if scorecard.Data.False_alert_fringe == ""{
    falseAlertFringe = 0
  }else{
    falseAlertFringe, err = strconv.Atoi(scorecard.Data.False_alert_fringe)
  }
  if falseAlertFringe > 0{
    if event.Data.Division != "Elite"{
      totalfaults += 2
    }
  }
  if scorecard.Data.Eliminated_during_search.Value == "yes" || scorecard.Data.Excused.Value == "yes"{
    if event.Data.Division != "Elite"{
      totalfaults += 3
    }else{
      totalfaults += 1
    }
  }
  if (scorecard.Data.Absent.Value == "yes")&&(event.Data.Division != "Elite"){
    totalfaults += 4
  }
  if err != nil {
    panic(err)
  }
  totalFaultsStr := strconv.Itoa(totalfaults)
  return totalFaultsStr
}


func (c *appContext) get_time(id string) string {
  scRepo := ScorecardRepo{c.db.C("scorecards")}
  scorecard, err := scRepo.Find(id)
  evRepo := EventRepo{c.db.C("events")}
  events, err := evRepo.All()
  event := EventResource{}
  falseAlertFringe, err := strconv.Atoi(scorecard.Data.False_alert_fringe)
  elapsed_time := ""
  for i:=0; i<len(events.Data); i++{
    if scorecard.Data.Event_Id == events.Data[i].Event_Id{
      event.Data = events.Data[i]
    }
  }
  if (scorecard.Data.Timed_out.Value == "yes" || scorecard.Data.Finish_call.Value == "no") && event.Data.Division == "Elite"{
      elapsed_time = scorecard.Data.Maxtime_m + ":" + scorecard.Data.Maxtime_s + ":00"
  }else if event.Data.Division != "Elite"{
    if scorecard.Data.Timed_out.Value == "yes" || scorecard.Data.Finish_call.Value == "no" || scorecard.Data.Absent.Value == "yes" || scorecard.Data.Eliminated_during_search.Value == "yes" || scorecard.Data.Excused.Value == "yes" || falseAlertFringe > 0{
      elapsed_time = scorecard.Data.Maxtime_m + ":" + scorecard.Data.Maxtime_s + ":00"
    }
  }
  if err != nil {
	panic(err)
  }
  return elapsed_time
}


func (c *appContext) get_points(id string) string {
  scRepo := ScorecardRepo{c.db.C("scorecards")}
  scorecard, err := scRepo.Find(id)
  evRepo := EventRepo{c.db.C("events")}
  events, err := evRepo.All()
  event := EventResource{}
  falseAlertFringe := 0.0
  totalFaults := 0.0
  maxPoint := 0.0
  hidesMax := 0.0
  hidesFound := 0.0
  eliteHides := 0.0
  if scorecard.Data.False_alert_fringe == ""{
    falseAlertFringe = 0.0
  }else{
    falseAlertFringe, err = strconv.ParseFloat(scorecard.Data.False_alert_fringe, 64)
  }
  if scorecard.Data.Total_faults == ""{
    totalFaults = 0.0
  }else{
    totalFaults, err = strconv.ParseFloat(scorecard.Data.Total_faults, 64)
  }
  if scorecard.Data.Maxpoint == ""{
    maxPoint = 0.0
  }else{
    maxPoint, err = strconv.ParseFloat(scorecard.Data.Maxpoint, 64)
  }
  if scorecard.Data.Hides_max == ""{
    hidesMax = 0.0
  }else{
    hidesMax, err = strconv.ParseFloat(scorecard.Data.Hides_max, 64)
  }
  if scorecard.Data.Hides_found == ""{
    hidesFound = 0.0
  }else{
    hidesFound, err = strconv.ParseFloat(scorecard.Data.Hides_found, 64)
  }
  if event.Data.Elite_hides == ""{
    eliteHides = 0.0
  }else{
    eliteHides, err = strconv.ParseFloat(event.Data.Elite_hides, 64)
  }
  for i:=0; i<len(events.Data); i++{
    if scorecard.Data.Event_Id == events.Data[i].Event_Id{
      event.Data = events.Data[i]
    }
  }
  points := 0.0
  if totalFaults <= 3{
    if hidesFound == hidesMax && event.Data.Division == "NW1"{
      points = maxPoint
    }else if (event.Data.Division == "NW2")||(event.Data.Division == "NW3") || (event.Data.Division == "Element Specialty"){
      points = hidesFound * maxPoint/hidesMax
    }else if event.Data.Division == "Elite"{
      if falseAlertFringe == 0{
        points = (hidesFound * maxPoint/hidesMax) - totalFaults
      }else if falseAlertFringe <= 3 && falseAlertFringe > 0{
        points = (hidesFound * maxPoint) - totalFaults + (falseAlertFringe * 100.0/eliteHides/2)
      }
      if scorecard.Data.Finish_call.Value == "no"{
        if hidesFound > 0 || falseAlertFringe > 0{
          points = points - 100.0/eliteHides/2
        }
      }
    }else{
      points = 0.0
    }
  }else{
    points = 0.0
  }
  if (scorecard.Data.Absent.Value == "yes")||(scorecard.Data.Eliminated_during_search.Value == "yes")||(scorecard.Data.Excused.Value == "yes"){
    points = 0.0
  }
  if err != nil {
	panic(err)
  }
  pointStr := strconv.FormatFloat(points, 'f', 2, 64)
  return pointStr
}


// Tally Handlers /////////////////////////////////////////////////////////////////////////////////////


func (c *appContext) talliesHandler(w http.ResponseWriter, r *http.Request) {
  if current_session.Current_status == false{
    http.Redirect(w, r, "/login", 302)
	}else{
    repo := TallyRepo{c.db.C("tallies")}
    tallies, err := repo.All()
    tallies_resrc := TalliesResource{}
    tallies_resrc.SData = current_session
    for j:=0; j<len(tallies.Data); j++{
      body := TallyResource{}
      body.Data.Id = tallies.Data[j].Id
      body.Data.Tally_Id = tallies.Data[j].Tally_Id
      body.Data.Event_Id = tallies.Data[j].Event_Id
      body.Data.Entrant_Id = tallies.Data[j].Entrant_Id
      body.Data.Total_points = tallies.Data[j].Total_points
      body.Data.Total_faults = tallies.Data[j].Total_faults
      body.Data.Total_time = tallies.Data[j].Total_time
      body.Data.Title = tallies.Data[j].Title
      body.Data.Qualifying_score = tallies.Data[j].Qualifying_score
      body.Data.Qualifying_scores = tallies.Data[j].Qualifying_scores
      tallies_resrc.Data = append(tallies_resrc.Data, body)
    }
    if err != nil {
      panic(err)
    }
    //	w.Header().Set("Content-Type", "application/vnd.api+json")
    //	json.NewEncoder(w).Encode(tallies)

    // read BSON into JSON
    if err = listTally.Execute(w, tallies_resrc); err != nil {
      http.Error(w, err.Error(), http.StatusInternalServerError)
      return
    }
  }
}

func (c *appContext) tallyHandler(w http.ResponseWriter, r *http.Request) {
  if current_session.Current_status == false{
    http.Redirect(w, r, "/login", 302)
  }else{
    params := context.Get(r, "params").(httprouter.Params)  //gorrila context, key "params"
    repo := TallyRepo{c.db.C("tallies")}
    tally, err := repo.Find(params.ByName("id")) //getting data from named param :id
    if err != nil {
      panic(err)
    }
    tally_resrc := TallyFormResource{}
    tally_resrc.TLYData = tally
    tally_resrc.SData = current_session
    evRepo := EventRepo{c.db.C("events")}
    events, err := evRepo.All()
    evbody := EventResource{}
    for i:=0; i<len(events.Data); i++{
      if events.Data[i].Event_Id == tally.Data.Event_Id{
        evbody.Data = events.Data[i]
      }
    }
    tally_resrc.EVData = evbody

    enRepo := EntrantRepo{c.db.C("entrants")}
    entrants, err := enRepo.All()
    enbody := EntrantResource{}
    for i:=0; i<len(entrants.Data); i++{
      if entrants.Data[i].Team_Id == tally.Data.Entrant_Id{
        for j:=0; j<len(entrants.Data[i].Event_Id); j++{
          if entrants.Data[i].Event_Id[j] == tally.Data.Event_Id{
            enbody.Data = entrants.Data[i]
          }
        }
      }
    }
    tally_resrc.ENData = enbody
    //  w.Header().Set("Content-Type", "application/vnd.api+json")
    //	json.NewEncoder(w).Encode(tally)
    //	if err = show.Execute(w, json.NewEncoder(w).Encode(tally)); err != nil {
    //      http.Error(w, err.Error(), http.StatusInternalServerError)
    //      return
    //  }

    // read JSON into BSON
    if err = showTally.Execute(w, tally_resrc); err != nil {
      http.Error(w, err.Error(), http.StatusInternalServerError)
      return
    }
  }
}

func (c *appContext) editTallyHandler(w http.ResponseWriter, r *http.Request) {
  if current_session.Current_status == false{
    http.Redirect(w, r, "/login", 302)
  }else{
    params := context.Get(r, "params").(httprouter.Params)    //gorilla context, key "params"
    repo := TallyRepo{c.db.C("tallies")}
    tally, err := repo.Find(params.ByName("id")) //getting data from named param :id
    tally_resrc := TallyFormResource{}
    tally_resrc.TLYData = tally
    tally_resrc.SData = current_session
    evRepo := EventRepo{c.db.C("events")}
    events, err := evRepo.All()
    evbody := EventResource{}
    for i:=0; i<len(events.Data); i++{
      if events.Data[i].Event_Id == tally.Data.Event_Id{
        evbody.Data = events.Data[i]
      }
    }
    tally_resrc.EVData = evbody

    enRepo := EntrantRepo{c.db.C("entrants")}
    entrants, err := enRepo.All()
    enbody := EntrantResource{}
    for i:=0; i<len(entrants.Data); i++{
      if entrants.Data[i].Team_Id == tally.Data.Entrant_Id{
        for j:=0; j<len(entrants.Data[i].Event_Id); j++{
          if entrants.Data[i].Event_Id[j] == tally.Data.Event_Id{
            enbody.Data = entrants.Data[i]
          }
        }
      }
    }
    tally_resrc.ENData = enbody

    scRepo := ScorecardRepo{c.db.C("scorecards")}
    scorecards, err := scRepo.All()
    scorecards_resrc := ScorecardsResource{}
    scbody := ScorecardResource{}
    for i:=0; i<len(scorecards.Data); i++{
      if scorecards.Data[i].Entrant_Id == tally.Data.Entrant_Id && scorecards.Data[i].Event_Id == tally.Data.Event_Id{
        scbody.Data = scorecards.Data[i]
        scorecards_resrc.Data = append(scorecards_resrc.Data, scbody)
      }
    }
    tally_resrc.SCSData = scorecards_resrc

    if err = updateTally.Execute(w, tally_resrc); err != nil {
      http.Error(w, err.Error(), http.StatusInternalServerError)
      return
    }
    if err != nil{
      fmt.Println(err)
    }
  }
}

func (c *appContext) updateTallyHandler(w http.ResponseWriter, r *http.Request) {
  if current_session.Current_status == false{
    http.Redirect(w, r, "/login", 302)
  }else{
    params := context.Get(r, "params").(httprouter.Params)
    repo := TallyRepo{c.db.C("tallies")}
    tally, err := repo.Find(params.ByName("id")) //getting data from named param :id
    body := context.Get(r, "body").(*TallyResource)
    body.Data.Id = tally.Data.Id
    body.Data.Event_Id = tally.Data.Event_Id
    body.Data.Entrant_Id = tally.Data.Entrant_Id
    body.Data.Tally_Id = tally.Data.Tally_Id
    body.Data.Total_time = r.FormValue("Total_time")
    body.Data.Total_faults = r.FormValue("Total_faults")
    body.Data.Title = r.FormValue("Title")
    body.Data.Total_points = r.FormValue("Total_points")
    body.Data.Qualifying_score = r.FormValue("Qualifying_score")
    body.Data.Qualifying_scores = r.FormValue("Qualifying_scores")

    err = repo.Update(&body.Data)

    c.get_tally(tally.Data.Id.Hex())

    if err != nil {
      panic(err)
    }
    if body.Data.Tally_Id == ""{
      http.Redirect(w, r, "/tallies/delete/" + body.Data.Id.Hex(), 302)
    }else{
      http.Redirect(w, r, "/tallies/edit/" + body.Data.Id.Hex(), 302)
    }
  }
}

func (c *appContext) deleteTallyHandler(w http.ResponseWriter, r *http.Request) {
  if current_session.Current_status == false{
    http.Redirect(w, r, "/login", 302)
  }else{
    params := context.Get(r, "params").(httprouter.Params)    //gorilla context, key "params"
    repo := TallyRepo{c.db.C("tallies")}
    err := repo.Delete(params.ByName("id"))
    if err != nil {
      panic(err)
    }
    http.Redirect(w, r, "/tallies", 302)
  }
}

func (c *appContext) get_tally(id string){
  tlyRepo := TallyRepo{c.db.C("tallies")}
  tallies, err := tlyRepo.All()
  tally, err := tlyRepo.Find(id)
  scRepo := ScorecardRepo{c.db.C("scorecards")}
  scorecards, err := scRepo.All()
  evRepo := EventRepo{c.db.C("events")}
  events, err := evRepo.All()
  event := EventResource{}
  enRepo := EntrantRepo{c.db.C("entrants")}
  entrants, err := enRepo.All()
  entrant := EntrantResource{}
  for i:=0; i<len(events.Data); i++{
    for j:=0; j<len(entrants.Data); j++{
      if tally.Data.Event_Id == events.Data[i].Event_Id && tally.Data.Entrant_Id == entrants.Data[j].Team_Id{
        event.Data = events.Data[i]
        entrant.Data = entrants.Data[j]
      }
    }
  }
  specialty := ""
  q_scores := 0
  point_tally := 0.0
  if event.Data.Division == "Element Specialty"{
    if event.Data.Cont_search_areas != "0" || event.Data.Cont_search_areas != ""{
      specialty = "Container"
    }else if event.Data.Ext_search_areas != "0" || event.Data.Ext_search_areas != ""{
      specialty = "Exterior"
    }else if event.Data.Int_search_areas != "0" || event.Data.Int_search_areas != ""{
      specialty = "Interior"
    }else{
      specialty = "Vehicle"
    }
    for i:=0; i<len(entrants.Data); i++{
      for j:=0; j<len(event.Data.EntrantSelected_Id); j++{
        if entrants.Data[i].Team_Id == event.Data.EntrantSelected_Id[j]{
          if tally.Data.Entrant_Id == entrants.Data[i].Team_Id{
            for k:=0; k<len(entrants.Data[i].Event_Id); k++{
              if event.Data.Event_Id != tally.Data.Event_Id{
                if event.Data.Division == "Element Specialty"{
                  if ((specialty == "Container") && (event.Data.Cont_search_areas != "0")) || ((specialty == "Container") && (event.Data.Cont_search_areas != "")) || ((specialty == "Exterior") && (event.Data.Ext_search_areas != "0")) || ((specialty == "Exterior") && (event.Data.Ext_search_areas != "")) || ((specialty == "Interior") && (event.Data.Int_search_areas != "0")) || ((specialty == "Interior") && (event.Data.Int_search_areas != "")) || ((specialty == "Vehicle")) && ((event.Data.Veh_search_areas != "0") || specialty == "Vehicle" && (event.Data.Veh_search_areas != "")){
                    for m:=0; m<len(tallies.Data); m++{
                      if tallies.Data[m].Event_Id == event.Data.Event_Id{
                        if tallies.Data[m].Entrant_Id == entrant.Data.Team_Id{
                          if tallies.Data[m].Qualifying_score == ""{
                            q_scores = 0
                          }else if tallies.Data[m].Qualifying_score == "1"{
                            q_scores += 1
                          }
                        }
                      }
                    }
                  }
                }
              }
            }
          }
        }
      }
    }
  }
  if q_scores >= 2{
    q_scores = 0
  }
  tly_total_points := 0.0
  if event.Data.Division == "Elite"{
    for i:=0; i<len(entrants.Data); i++{
      for j:=0; j<len(event.Data.EntrantSelected_Id); j++{
        if entrants.Data[i].Team_Id == event.Data.EntrantSelected_Id[j]{
          if tally.Data.Entrant_Id == entrants.Data[i].Team_Id{
            for k:=0; k<len(entrants.Data[i].Event_Id); k++{
              if event.Data.Event_Id != tally.Data.Event_Id{
                if event.Data.Division == "Elite"{
                  for m:=0; m<len(tallies.Data); m++{
                    if tallies.Data[m].Event_Id == event.Data.Event_Id{
                      if tallies.Data[m].Entrant_Id == entrant.Data.Team_Id{
                        if tallies.Data[m].Total_points == ""{
                          tly_total_points = 0
                        }else{
                          tly_total_points, err = strconv.ParseFloat(tallies.Data[m].Total_points, 64)
                        }
                        point_tally += tly_total_points
                      }
                    }
                  }
                }
              }
            }
          }
        }
      }
    }
  }
  sc_total_points := 0.0
  sc_time := 0
  sc_faults := 0.0
  point_tally_round := 0
  time_tally := 0
  fault_tally := 0.0
  q_score := 0
  titled := ""
  for i:=0; i<len(entrants.Data); i++{
    for j:=0; j<len(event.Data.EntrantSelected_Id); j++{
      if entrants.Data[i].Team_Id == event.Data.EntrantSelected_Id[j]{
        if tally.Data.Entrant_Id == entrants.Data[i].Team_Id{
          time_tally = 0.0
          fault_tally = 0.0
          q_score = 0
          titled = ""
          for k:=0; k<len(scorecards.Data); k++{
            if scorecards.Data[k].Event_Id == tally.Data.Event_Id{
              if scorecards.Data[k].Entrant_Id == entrant.Data.Team_Id{
                if scorecards.Data[k].Total_points == ""{
                  sc_total_points = 0
                }else{
                  sc_total_points, err = strconv.ParseFloat(scorecards.Data[k].Total_points, 64)
                }
                point_tally += sc_total_points

                if scorecards.Data[k].Total_time == ""{
                  sc_time = 0
                }else{
                  sc_time = str_to_time(scorecards.Data[k].Total_time)
                }
                time_tally += sc_time

                if scorecards.Data[k].Total_faults == ""{
                  sc_faults = 0
                }else{
                  sc_faults, err = strconv.ParseFloat(scorecards.Data[k].Total_faults, 64)
                }
                fault_tally += sc_faults

                // round point tally
                fldata := math.Floor(point_tally)
                clgdata := math.Ceil(point_tally)
                fdiff := point_tally - fldata
                cdiff := clgdata - point_tally
                if fdiff > cdiff{
                  point_tally_round = int(fldata)
                }else{
                  point_tally_round = int(clgdata)
                }
                if point_tally_round == 100 && fault_tally <= 3 && event.Data.Division != "Element Specialty" && event.Data.Division != "Elite"{
                  titled = "Titled"
                }else if event.Data.Division == "Element Specialty"{
                  // round point tally
                  if point_tally_round >= 75 && fault_tally <= 3{
                    q_score = 1
                    q_scores += 1
                    if q_scores == 2{
                      titled = "Titled"
                    }
                  }
                  // round point tally
                  if point_tally_round == 100 && fault_tally <= 3{
                    titled = "Titled"
                  }
                 // round point tally
                }else if event.Data.Division == "Elite" && point_tally_round >= 150 && fault_tally <= 3{
                  titled = "Titled"
                }else{
                  titled = "Not this time"
                }
              }
            }
          }
        }
      }
    }
  }
  var sec float64
  var min float64
  msec := time_tally%1000
  msecFlt := float64(msec)
  time_tally_flt := float64(time_tally)
  if (msec < 1) {
     msec = 0
  } else {
      // calculate seconds
      sec = (time_tally_flt-msecFlt)/1000;
      secInt := int(sec)
      secIntMod := secInt%60
      secModFlt := float64(secIntMod)
      if (sec < 1) {
          sec = 0
      } else {
          // calculate minutes
          min := (sec-secModFlt)/60;
          if (min < 1) {
              min = 0
          }
      }
  }
  // substract elapsed minutes
  msecFlt = msecFlt/10
  msecFltRnd := math.Floor(msecFlt)
  msecIntRnd := int(msecFltRnd)
  sec = sec-(min*60)
  secInt := int(sec)
  minInt := int(min)

  m_str := strconv.Itoa(minInt)
  if minInt < 10{
    m_str = "0" + m_str
  }
  s_str := strconv.Itoa(secInt)
  if secInt < 10{
    s_str = "0" + s_str
  }
  ms_str := strconv.Itoa(msecIntRnd)
  if msecIntRnd < 10{
    ms_str = "0" + ms_str
  }
  time_tally_str := m_str + ":" + s_str + ":" + ms_str
  point_tally_str := strconv.FormatFloat(point_tally, 'f', 2, 64)
  fault_tally_str := strconv.FormatFloat(fault_tally, 'f', 2, 64)
  q_score_str := strconv.Itoa(q_score)
  q_scores_str := strconv.Itoa(q_scores)
  tbody := TallyResource{}
  tbody.Data.Total_time = time_tally_str
  tbody.Data.Total_points = point_tally_str
  tbody.Data.Total_faults = fault_tally_str
  tbody.Data.Title = titled
  tbody.Data.Qualifying_score = q_score_str
  tbody.Data.Qualifying_scores = q_scores_str
  tbody.Data.Id = tally.Data.Id
  tbody.Data.Event_Id = tally.Data.Event_Id
  tbody.Data.Entrant_Id = tally.Data.Entrant_Id
  tbody.Data.Tally_Id = tally.Data.Tally_Id
  err = tlyRepo.Update(&tbody.Data)
  if err != nil {
	panic(err)
  }
  return
}


func str_to_time (time string) int{
  var time_int int
  if (time == "") || (time == "0"){
    time_int = 0
  }else{
    tdata := strings.SplitN(time, ":", 3)
    var minInt int
    var secInt int
    var msecInt int
    var err error
    zeros := []string{"00"}
    tdatam := tdata[:1]
    if tdatam[0] == zeros[0]{
      minInt = 0
    }
    matched, err := regexp.MatchString("0.^.?", tdatam[0])
    if matched{
      minStr := strings.TrimPrefix(tdatam[0], "0")
      minInt, err = strconv.Atoi(minStr)
    }else{
      minInt, err = strconv.Atoi(tdatam[0])
    }
    tdatas := tdata[1:2]
    if tdatas[0] == zeros[0]{
      secInt = 0
    }
    matched, err = regexp.MatchString("0.^", tdatas[0])
    if matched{
      secStr := strings.TrimPrefix(tdatas[0], "0")
      secInt, err = strconv.Atoi(secStr)
    }else{
      secInt, err = strconv.Atoi(tdatas[0])
    }
    tdatams := tdata[2:3]
    if tdatams[0] == zeros[0]{
      msecInt, err = strconv.Atoi("0")
    }
    matched, err = regexp.MatchString("0.^", tdatams[0])
    if matched{
      msecStr := strings.TrimPrefix(tdatams[0], "0")
      msecInt, err = strconv.Atoi(msecStr)
    }else{
      msecInt, err = strconv.Atoi(tdatams[0])
    }
    time_int = minInt*60*1000 + secInt*1000 + msecInt*10
    if err != nil {
      fmt.Println(err)
    }
  }
  return time_int
}



// Session Handlers /////////////////////////////////////////////////////////////////////////////////////

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    http.Redirect(w, r, "/login", 302)
}


func newSessionHandler(w http.ResponseWriter, r *http.Request) {
  // go to login page
  sessionresrc := SessionResource{}
  current_session = Session{Current_user: "", Current_email: "", Current_status: false}
  sessionresrc.SData = current_session
  r.SetBasicAuth("", "")
  if err := createnewSession.Execute(w, sessionresrc); err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  // submit forwards to openSessionHandler
}

func (c *appContext) sessionHandler(w http.ResponseWriter, r *http.Request){
  // go to login page and either sign up or enter email and password, get authenticated and go to "/events"
  usRepo := UserRepo{c.db.C("users")}
  users, err := usRepo.All()
  var current_user = ""
  fmt.Println(err)
  Email := r.FormValue("email")
  Password := r.FormValue("password")
  if Email == ""{
    current_user = current_session.Current_user
    Email = current_session.Current_email
  }

  if len(users.Data) > 0{
    for i:=0;i<len(users.Data);i++{
      if ((Password == users.Data[i].Password) && (Email == users.Data[i].Email)) || ((current_user == users.Data[i].User_Id) && (Email == users.Data[i].Email)){
        current_session.Current_user = users.Data[i].User_Id
        r.SetBasicAuth(Email, Password)
      }
    }
    email, password, ok := r.BasicAuth()
    if (ok == true) && (email != "") && (password != ""){
      for i:=0;i<len(users.Data);i++{
        if (password == users.Data[i].Password) && (email == users.Data[i].Email) && (current_session.Current_user == users.Data[i].User_Id){
          current_session.Current_status = true
          current_session.Current_email = users.Data[i].Email
        }
      }
    }

    refurlstr := ""
    matchedRef := false
    matchedUrl := false

    // if referer is login, go to events as start page
    if r.Header["Referer"] != nil{
      refurlstr = r.Header["Referer"][0]
      matchedRef, err = regexp.MatchString("login", refurlstr)
    }

    urlsrc := r.URL
    urlstr := urlsrc.String()

    // if referer and current url are the same
    matchedUrl, err = regexp.MatchString(urlstr, refurlstr)
    fmt.Println(matchedUrl)
    if current_session.Current_status == false{
      http.Redirect(w, r, "/login", 302)
    }else{
      // if referer is login or referer is nil or the referer and current url are not the same
      // go to events as home page
      if matchedRef || (r.Header["Referer"] == nil) || ((urlstr != "/signup") && matchedRef){
        http.Redirect(w, r, "/events", 302)
      }else if !matchedRef || (!matchedRef && urlstr == "/signup"){
        http.Redirect(w, r, urlstr, 302)
      }
    }
  }else{
    http.Redirect(w, r, "/login", 302)
  }
}

func (c *appContext) deleteSessionHandler(w http.ResponseWriter, r *http.Request) {
  //  De-authenticate user and direct to login - triggered by "/logout"
  usRepo := UserRepo{c.db.C("users")}
  users, err := usRepo.All()
  fmt.Println(err)
  var test_str = []string{"/info", "/events", "/entrants", "/users", "/scorecards", "/tallies"}
  refurlstr := ""
  scurlstr := ""
  matchedRefurl := false
  matchedUrl := false
  matchedUpdate := false
  matchedNew := false
  matchedEvent := false
  matchedfav := false
  mrurl := false
  murl := false
  update := false
  newevent := false

  if r.Header["Referer"] != nil{
    refurlstr = r.Header["Referer"][0]
  }
  urlsrc := r.URL
  scurlstr = urlsrc.String()
  matchedfav, err = regexp.MatchString("/favicon.ico", scurlstr)
  if matchedfav{
    scurlstr = refurlstr
  }
  for i:=0; i<len(test_str); i++{
    matchedUrl, err = regexp.MatchString(test_str[i], scurlstr)
    matchedUpdate, err = regexp.MatchString("/update", scurlstr)
    matchedNew, err = regexp.MatchString("/new", scurlstr)
    matchedEvent, err = regexp.MatchString("/events", scurlstr)
    matchedRefurl, err = regexp.MatchString(test_str[i], refurlstr)
    if matchedUrl && !matchedUpdate{
      murl = true
      scurlstr = test_str[i]
    }else if matchedUrl && matchedUpdate{
      murl = true
      update = true
    }else if matchedEvent && matchedNew{
      murl = true
      newevent = true
    }
    if matchedRefurl == true{
      mrurl = true
    }
  }
  if murl && mrurl && !update{
    http.Redirect(w, r, scurlstr, 302)
  }else if murl && mrurl && update{
    http.Redirect(w, r, refurlstr, 302)
  }else if murl && mrurl && newevent{
    http.Redirect(w, r, refurlstr, 302)
  }else{
    for i:=0;i<len(users.Data);i++{
      if current_session.Current_user == users.Data[i].User_Id{
        current_session.Current_user = ""
        current_session.Current_status = false
        current_session.Current_email = ""
      }
    }
    http.Redirect(w, r, "/login", 302)
  }
}


// Router //////////////////////////////////////////////////////////////////////////////////////////

type router struct {
	*httprouter.Router
}

func (r *router) Get(path string, handler http.Handler) {
	r.GET(path, wrapHandler(handler))
}

func (r *router) Post(path string, handler http.Handler) {
	r.POST(path, wrapHandler(handler))
}

func (r *router) Delete(path string, handler http.Handler) {
	r.DELETE(path, wrapHandler(handler))
}

// Integrating httprouter to our frameworks where it is incompatible with
// go http.Handler
// We wrap our middleware stack - implementing http.Handler into a
// httprouter.Handler function

func NewRouter() *router {
	return &router{httprouter.New()}
}

func wrapHandler(h http.Handler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
        context.Set(r, "params", ps)    //gorilla context, key "params"
        h.ServeHTTP(w, r)

	}
}


// MAIN ////////////////////////////////////////////////////////////////////////////////////////////

func main() {

  // port := os.Getenv("PORT")
  // if port == ""{
  //    log.Fatal("$PORT must be set")
  // }

  // session, err := mgo.Dial("mongodb://heroku_g884mk05:souabj4nqoh1r5ok1v0uss74ju@ds251889.mlab.com:51889/heroku_g884mk05")
  session, err := mgo.Dial("localhost:27017")

  if err != nil {
	panic(err)
  }
  defer session.Close()

  session.SetMode(mgo.Monotonic, true)

  // appC := appContext{session.DB("heroku_g884mk05")}
  appC := appContext{session.DB("test")}

  commonHandlers := alice.New(context.ClearHandler, loggingHandler, recoverHandler)
  // alice is used to chain handlers
  // context from gorrila mapping

  router := NewRouter()

  router.HandleMethodNotAllowed = false

  router.ServeFiles("/static/*filepath", http.Dir("static"))

  router.Get("/info", commonHandlers.ThenFunc(infoHandler))

  //  Session routing ////////////////

  router.Get("/signup", commonHandlers.Append(bodyHandler(UserResource{})).ThenFunc(appC.newUserHandler))
  router.GET("/", Index)
  router.Get("/login", commonHandlers.ThenFunc(newSessionHandler))
  router.Post("/session", commonHandlers.ThenFunc(appC.sessionHandler))
  router.Get("/logout", commonHandlers.ThenFunc(appC.deleteSessionHandler))

  //  Event routing  /////////////////

  router.Get("/events", commonHandlers.ThenFunc(appC.eventsHandler))
  router.Get("/events/show/:id", commonHandlers.ThenFunc(appC.eventHandler))
  router.Get("/events/new", commonHandlers.Append(bodyHandler(EventResource{})).ThenFunc(appC.newEventHandler))
  router.Post("/events/create", commonHandlers.Append(bodyHandler(EventResource{})).ThenFunc(appC.createEventHandler))
  router.Get("/events/edit/:id/", commonHandlers.ThenFunc(appC.editEventHandler))
  router.Post("/events/update/:id/", commonHandlers.Append(bodyHandler(EventResource{})).ThenFunc(appC.updateEventHandler))
  router.Get("/events/delete/:id", commonHandlers.ThenFunc(appC.deleteEventHandler))

  //  Entrant routing  //////////////////

  router.Get("/entrants", commonHandlers.ThenFunc(appC.entrantsHandler))
  router.Get("/entrants/show/:id", commonHandlers.ThenFunc(appC.entrantHandler))
  router.Get("/entrants/new", commonHandlers.ThenFunc(appC.newEntrantHandler))
  router.Post("/entrants/create", commonHandlers.Append(bodyHandler(EntrantResource{})).ThenFunc(appC.createEntrantHandler))
  router.Get("/entrants/edit/:id", commonHandlers.ThenFunc(appC.editEntrantHandler))
  router.Post("/entrants/update/:id/", commonHandlers.Append(bodyHandler(EntrantResource{})).ThenFunc(appC.updateEntrantHandler))
  router.Get("/entrants/delete/:id", commonHandlers.ThenFunc(appC.deleteEntrantHandler))

  //  User routing  //////////////////

  router.Get("/users", commonHandlers.ThenFunc(appC.usersHandler))
  router.Get("/users/show/:id", commonHandlers.ThenFunc(appC.userHandler))
  router.Get("/users/new", commonHandlers.ThenFunc(appC.newUserHandler))
  router.Post("/users/create", commonHandlers.Append(bodyHandler(UserResource{})).ThenFunc(appC.createUserHandler))
  router.Get("/users/edit/:id", commonHandlers.ThenFunc(appC.editUserHandler))
  router.Post("/users/update/:id/", commonHandlers.Append(bodyHandler(UserResource{})).ThenFunc(appC.updateUserHandler))
  router.Get("/users/delete/:id", commonHandlers.ThenFunc(appC.deleteUserHandler))


  //  Scorecard routing  //////////////////


  router.Get("/scorecards", commonHandlers.ThenFunc(appC.scorecardsHandler))
  router.Get("/scorecards/show/:id", commonHandlers.ThenFunc(appC.scorecardHandler))
  router.Get("/scorecards/edit/:id", commonHandlers.ThenFunc(appC.editScorecardHandler))
  router.Post("/scorecards/update/:id/", commonHandlers.Append(bodyHandler(ScorecardResource{})).ThenFunc(appC.updateScorecardHandler))
  router.Get("/scorecards/delete/:id", commonHandlers.ThenFunc(appC.deleteScorecardHandler))


  //  Tally routing  //////////////////


  router.Get("/tallies", commonHandlers.ThenFunc(appC.talliesHandler))
  router.Get("/tallies/show/:id", commonHandlers.ThenFunc(appC.tallyHandler))
  router.Get("/tallies/edit/:id", commonHandlers.ThenFunc(appC.editTallyHandler))
  router.Post("/tallies/update/:id", commonHandlers.Append(bodyHandler(TallyResource{})).ThenFunc(appC.updateTallyHandler))
  router.Get("/tallies/delete/:id", commonHandlers.ThenFunc(appC.deleteTallyHandler))


  //  listening
  // http.ListenAndServe((":" + port), router)
  http.ListenAndServe(":8080", router)
}

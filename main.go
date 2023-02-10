package main

import (
	"encoding/json"
	"log"
	"mime"
	"net/http"
	"os"
	"strconv"

	"employee-base/internal/employee"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type employeeServer struct {
	storage *employee.EmployeeStorage
}

func NewEmployeeServer() *employeeServer {
	storage := employee.New()
	return &employeeServer{storage: storage}
}

// renderJSON renders 'v' as JSON and writes it as a response into w.
func renderJSON(w http.ResponseWriter, v interface{}) {
	js, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (es *employeeServer) createEmployeeHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling employee create at %s\n", req.URL.Path)

	// Types used internally in this handler to (de-)serialize the request and
	// response from/to JSON.
	type RequestEmployee struct {
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		Email     string `json:"email"`
	}

	type ResponseId struct {
		Id int `json:"id"`
	}

	// Enforce a JSON Content-Type.
	contentType := req.Header.Get("Content-Type")
	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if mediatype != "application/json" {
		http.Error(w, "expect application/json Content-Type", http.StatusUnsupportedMediaType)
		return
	}

	dec := json.NewDecoder(req.Body)
	dec.DisallowUnknownFields()
	var re RequestEmployee
	if err := dec.Decode(&re); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := es.storage.CreateEmployee(re.FirstName, re.LastName, re.Email)
	renderJSON(w, ResponseId{Id: id})
}

func (es *employeeServer) getAllEmployeesHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling get all employees at %s\n", req.URL.Path)

	allEmployees := es.storage.GetAllEmployees()
	renderJSON(w, allEmployees)
}

func (es *employeeServer) getEmployeeHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling get employee at %s\n", req.URL.Path)

	// Here and elsewhere, not checking error of Atoi because the router only
	// matches the [0-9]+ regex.
	id, _ := strconv.Atoi(mux.Vars(req)["id"])
	employee, err := es.storage.GetEmployee(id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	renderJSON(w, employee)
}

func (es *employeeServer) deleteEmployeeHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling delete employee at %s\n", req.URL.Path)

	id, _ := strconv.Atoi(mux.Vars(req)["id"])
	err := es.storage.DeleteEmployee(id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}
}

func (es *employeeServer) deleteAllEmployeesHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling delete all employee at %s\n", req.URL.Path)
	es.storage.DeleteAllEmployees()
}

func (es *employeeServer) lastNameHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling employee by lastName at %s\n", req.URL.Path)

	lastName := mux.Vars(req)["lastName"]
	employees, err := es.storage.GetEmployeesByLastName(lastName)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	renderJSON(w, employees)
}

func (es *employeeServer) updateEmployeeHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling employee update at %s\n", req.URL.Path)

	// Types used internally in this handler to (de-)serialize the request and
	// response from/to JSON.
	type RequestEmployee struct {
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		Email     string `json:"email"`
	}

	// Enforce a JSON Content-Type.
	contentType := req.Header.Get("Content-Type")
	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if mediatype != "application/json" {
		http.Error(w, "expect application/json Content-Type", http.StatusUnsupportedMediaType)
		return
	}

	dec := json.NewDecoder(req.Body)
	dec.DisallowUnknownFields()
	var re RequestEmployee
	if err := dec.Decode(&re); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, _ := strconv.Atoi(mux.Vars(req)["id"])
	err = es.storage.UpdateEmployee(id, re.FirstName, re.LastName, re.Email)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
}

func main() {
	router := mux.NewRouter()
	router.StrictSlash(true)
	server := NewEmployeeServer()

	router.HandleFunc("/employee/", server.createEmployeeHandler).Methods("POST")
	router.HandleFunc("/employee/", server.getAllEmployeesHandler).Methods("GET")
	router.HandleFunc("/employee/", server.deleteAllEmployeesHandler).Methods("DELETE")
	router.HandleFunc("/employee/{id:[0-9]+}/", server.getEmployeeHandler).Methods("GET")
	router.HandleFunc("/employee/{id:[0-9]+}/", server.deleteEmployeeHandler).Methods("DELETE")
	router.HandleFunc("/employee/{id:[0-9]+}/", server.updateEmployeeHandler).Methods("PUT")
	router.HandleFunc("/employee/{lastName}/", server.lastNameHandler).Methods("GET")

	// Set up logging and panic recovery middleware.
	router.Use(func(h http.Handler) http.Handler {
		return handlers.LoggingHandler(os.Stdout, h)
	})
	router.Use(handlers.RecoveryHandler(handlers.PrintRecoveryStack(true)))

	log.Fatal(http.ListenAndServe("localhost:"+os.Getenv("SERVERPORT"), router))
}

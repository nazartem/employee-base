package main

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"log"
	"mime"
	"net/http"
	"os"
	"strconv"

	"employee-base/internal/employee"
	"employee-base/internal/middleware"

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
	certFile := flag.String("certfile", "cert.pem", "certificate PEM file")
	keyFile := flag.String("keyfile", "key.pem", "key PEM file")
	flag.Parse()

	router := mux.NewRouter()
	router.StrictSlash(true)
	server := NewEmployeeServer()

	// The "create", "update" and "delete" paths are protected with the BasicAuth middleware.
	router.Handle("/employee/", middleware.BasicAuth(http.HandlerFunc(server.createEmployeeHandler))).Methods("POST") //
	router.HandleFunc("/employee/", server.getAllEmployeesHandler).Methods("GET")
	router.HandleFunc("/employee/{id:[0-9]+}/", server.getEmployeeHandler).Methods("GET")
	router.Handle("/employee/{id:[0-9]+}/", middleware.BasicAuth(http.HandlerFunc(server.deleteEmployeeHandler))).Methods("DELETE") //
	router.Handle("/employee/{id:[0-9]+}/", middleware.BasicAuth(http.HandlerFunc(server.updateEmployeeHandler))).Methods("PUT")    //
	router.HandleFunc("/employee/{lastName}/", server.lastNameHandler).Methods("GET")

	// Set up logging and panic recovery middleware for all paths.
	router.Use(func(h http.Handler) http.Handler {
		return handlers.LoggingHandler(os.Stdout, h)
	})
	router.Use(handlers.RecoveryHandler(handlers.PrintRecoveryStack(true)))

	addr := "localhost:" + os.Getenv("SERVERPORT")
	srv := &http.Server{
		Addr:    addr,
		Handler: router,
		TLSConfig: &tls.Config{
			MinVersion:               tls.VersionTLS13,
			PreferServerCipherSuites: true,
		},
	}

	log.Printf("Starting server on %s", addr)
	log.Fatal(srv.ListenAndServeTLS(*certFile, *keyFile))
}

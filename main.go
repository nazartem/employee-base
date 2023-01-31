package main

import (
	"encoding/json"
	"fmt"
	"log"
	"mime"
	"net/http"
	"os"
	"strconv"
	"strings"

	"/home/nazaryap/go/employee-base/internal/storage"
)

type employeeServer struct {
	storage *storage.EmployeeStorage
}

func NewEmployeeServer() *employeeServer {
	storage := storage.New()
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

func (es *employeeServer) employeeHandler(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path == "/employee/" {
		// Запрос направлен к "/employee/", без идущего в конце ID.
		if req.Method == http.MethodPost {
			es.createEmployeeHandler(w, req)
		} else if req.Method == http.MethodGet {
			es.getAllEmployeesHandler(w, req)
		} else if req.Method == http.MethodDelete {
			es.deleteAllEmployeesHandler(w, req)
		} else {
			http.Error(w, fmt.Sprintf("expect method GET, DELETE or POST at /employee/, got %v", req.Method), http.StatusMethodNotAllowed)
			return
		}
	} else {
		// В запросе есть ID, выглядит он как "/employee/<id>".
		path := strings.Trim(req.URL.Path, "/")
		fmt.Println(path) // удалить
		pathParts := strings.Split(path, "/")
		fmt.Println(path) // удалить

		if len(pathParts) < 2 {
			http.Error(w, "expect /employee/<id> in employee handler", http.StatusBadRequest)
			return
		}
		id, err := strconv.Atoi(pathParts[1])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if req.Method == http.MethodDelete {
			es.deleteEmployeeHandler(w, req, int(id))
		} else if req.Method == http.MethodGet {
			es.getEmployeeHandler(w, req, int(id))
		} else {
			http.Error(w, fmt.Sprintf("expect method GET or DELETE at /employee/<id>, got %v", req.Method), http.StatusMethodNotAllowed)
			return
		}
	}
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

func (es *employeeServer) getEmployeeHandler(w http.ResponseWriter, req *http.Request, id int) {
	log.Printf("handling get employee at %s\n", req.URL.Path)

	employee, err := es.storage.GetEmployee(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	renderJSON(w, employee)
}

func (es *employeeServer) deleteEmployeeHandler(w http.ResponseWriter, req *http.Request, id int) {
	log.Printf("handling delete employee at %s\n", req.URL.Path)

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

	if req.Method != http.MethodGet {
		http.Error(w, fmt.Sprintf("expect method GET /employee/<lastName>, got %v", req.Method), http.StatusMethodNotAllowed)
		return
	}

	path := strings.Trim(req.URL.Path, "/")
	pathParts := strings.Split(path, "/")
	if len(pathParts) < 2 {
		http.Error(w, "expect /employee/<lastName> path", http.StatusBadRequest)
		return
	}
	lastName := pathParts[1]

	employees := es.storage.GetEmployeesByLastName(lastName)
	renderJSON(w, employees)
}

// rework
func (es *employeeServer) updateEmployeeHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling employee update at %s\n", req.URL.Path)

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

func main() {
	mux := http.NewServeMux()
	server := NewEmployeeServer()
	mux.HandleFunc("/employee/", server.employeeHandler)

	log.Fatal(http.ListenAndServe("localhost:"+os.Getenv("SERVERPORT"), mux))
}

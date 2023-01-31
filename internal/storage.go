package storage

import (
	"fmt"
	"sync"
)

type Employee struct {
	Id        int    `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

// EmployeeStorage is a simple in-memory database of employees; EmployeeStorage methods are
// safe to call concurrently.
type EmployeeStorage struct {
	sync.Mutex

	employees map[int]Employee
	nextId    int
}

func New() *EmployeeStorage {
	es := &EmployeeStorage{}
	es.employees = make(map[int]Employee)
	es.nextId = 0
	return es
}

// CreateEmployee создаёт нового работника в хранилище.
func (es *EmployeeStorage) CreateEmployee(firstName, lastName, email string) int {
	es.Lock()
	defer es.Unlock()

	employee := Employee{
		Id:        es.nextId,
		FirstName: firstName,
		LastName:  lastName,
		Email:     email}

	es.employees[es.nextId] = employee
	es.nextId++
	return employee.Id
}

// GetAllEmployees возвращает из хранилища всех работников в произвольном порядке.
func (es *EmployeeStorage) GetAllEmployees() []Employee {
	es.Lock()
	defer es.Unlock()

	allEmployees := make([]Employee, 0, len(es.employees))
	for _, employee := range es.employees {
		allEmployees = append(allEmployees, employee)
	}
	return allEmployees
}

// GetEmployee получает работника из хранилища по ID. Если ID не существует -
// будет возвращена ошибка.
func (es *EmployeeStorage) GetEmployee(id int) (Employee, error) {
	es.Lock()
	defer es.Unlock()

	e, ok := es.employees[id]
	if ok {
		return e, nil
	} else {
		return Employee{}, fmt.Errorf("employee with id=%d not found", id)
	}
}

// DeleteEmployee удаляет работника с заданным ID. Если ID не существует -
// будет возвращена ошибка.
func (es *EmployeeStorage) DeleteEmployee(id int) error

// DeleteAllEmployees удаляет из хранилища всех работников.
func (es *EmployeeStorage) DeleteAllEmployees() error

// GetEmployeesByLastName возвращает, в произвольном порядке, всех работников
// с указанной фамилией.
func (es *EmployeeStorage) GetEmployeesByLastName(LastName string) []Employee

// UpdateEmployee обновляет информацию о работнике.
func (es *EmployeeStorage) UpdateEmployee(Id int, FirstName, LastName, Email string)

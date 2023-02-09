package employee

import (
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

package employee

import (
	"fmt"
)

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
func (es *EmployeeStorage) DeleteEmployee(id int) error {
	es.Lock()
	defer es.Unlock()

	if _, ok := es.employees[id]; !ok {
		return fmt.Errorf("employee with id=%d not found", id)
	}

	delete(es.employees, id)
	return nil
}

// GetEmployeesByLastName возвращает, в произвольном порядке, всех работников
// с указанной фамилией.
func (es *EmployeeStorage) GetEmployeesByLastName(lastName string) ([]Employee, error) {
	es.Lock()
	defer es.Unlock()

	var employees []Employee

	for _, employee := range es.employees {
		if employee.LastName == lastName {
			employees = append(employees, employee)
		}
	}

	if employees != nil {
		return employees, nil
	}

	return nil, fmt.Errorf("employees with lasName=%s not found", lastName)
}

// UpdateEmployee обновляет информацию о работнике.
func (es *EmployeeStorage) UpdateEmployee(id int, firstName, lastName, email string) error {
	es.Lock()
	defer es.Unlock()

	if _, ok := es.employees[id]; !ok {
		return fmt.Errorf("employee with id=%d not found", id)
	}

	es.employees[id] = Employee{
		Id:        id,
		FirstName: firstName,
		LastName:  lastName,
		Email:     email}

	return nil
}

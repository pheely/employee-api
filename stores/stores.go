package stores

type Store interface {
	Connect() error
	Create(list string, employee *Employee) error
	Clear(list string) error
	Get(list string, id string) (*Employee, error)
	Update(list string, id string, employee *Employee) (*Employee, error)
	Delete(list string, id string) error
	List(list string) ([]Employee, error)
}

type Employee struct {
	ID        string
	FirstName string
	LastName  string
	Department string
	Salary    int
	Age       int
}
package model

type Person struct {
	Id string

	Name string
	Age int
}

func (p *Person) ApplyFrom(srcPerson *Person) {
	if srcPerson.Name != "" {
		p.Name = srcPerson.Name
	}
	if srcPerson.Age > 0 {
		p.Age = srcPerson.Age
	}
}

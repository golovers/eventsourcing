package model

type Person struct {
	Id string `json:"id,omitempty"`

	Name string `json:"name,omitempty"`
	Age  int    `json:"age,omitempty"`
}

func (p *Person) ApplyFrom(srcPerson *Person) {
	if srcPerson.Name != "" {
		p.Name = srcPerson.Name
	}
	if srcPerson.Age > 0 {
		p.Age = srcPerson.Age
	}
}

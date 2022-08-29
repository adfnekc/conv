package converter

type Converter interface {
	From([]byte) ([]byte, error)
	To([]byte) []byte
}

var converter = make(map[string]Converter)

func AddConverter(name string, description string, c Converter) {
	if _, ok := converter[name]; ok {
		panic("two converters with the same name are not allowed")
	}
	converter[name] = c
}

func Get(name string) Converter {
	return converter[name]
}

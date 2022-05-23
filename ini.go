package ini

import (
	"bufio"
	"errors"
	"os"
	"reflect"
	"strconv"
	"strings"
)

// Cfg contains a Section map parsing from the configuration file.
type Cfg struct {
	sections map[string]*Section
}

// Load loads the configuration file according to fPath
// and returns a Cfg struct containing all the configuration.
func Load(fPath string) (*Cfg, error) {
	f, err := os.Open(fPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var section *Section

	sections := make(map[string]*Section)

	reader := bufio.NewScanner(f)
	for reader.Scan() {

		text := reader.Text()

		text = strings.TrimSpace(text)

		// white line
		if len(text) == 0 {
			continue
		}

		// Comment
		if text[0] == ';' || text[0] == '#' {
			continue
		}

		// Section
		if text[0] == '[' {
			closeIdx := strings.LastIndexByte(text, ']')

			name := text[1:closeIdx]
			section = &Section{
				name:   name,
				fields: make(map[string]*Field),
			}
			sections[name] = section

			continue
		}

		// Field
		kv := strings.Split(text, "=")

		section.fields[kv[0]] = &Field{
			key:   kv[0],
			value: kv[1],
		}
	}

	return &Cfg{sections}, nil
}

// Section gets the appointed section.
// If the section is not exist, an empty Section will be returned instead of nil pointer.
func (cfg *Cfg) Section(name string) *Section {
	s, ok := cfg.sections[name]
	if !ok {
		return &Section{}
	}
	return s
}

// Section contains the Section name and several Fields belong it.
type Section struct {
	name   string
	fields map[string]*Field
}

// Field will return a new empty struct instead of a nil pointer when the field name is not exist.
// It is same as the Section.
func (s *Section) Field(name string) *Field {
	field, ok := s.fields[name]
	if !ok {
		return &Field{}
	}
	return field
}

// Field contains only one part of key-value.
type Field struct {
	key   string
	value string
}

func (f *Field) String() string {
	return f.value
}

func (f *Field) Float64() float64 {
	des, err := strconv.ParseFloat(f.value, 64)
	if err != nil {
		return 0
	}
	return des
}

func (f *Field) Int64() int64 {
	des, err := strconv.ParseInt(f.value, 10, 64)
	if err != nil {
		return 0
	}
	return des
}

var ErrBindPtr = errors.New("bind des needs to be a pointer")

// Bind offers more easy way to get configuration. Bind also calls the
// Load and bind value to the field according to the struct tag.
func Bind(fPath string, des any) error {
	cfg, err := Load(fPath)
	if err != nil {
		return err
	}

	desType := reflect.TypeOf(des)
	desValue := reflect.ValueOf(des)

	if desType.Kind() != reflect.Ptr {
		return ErrBindPtr
	}

	for i := 0; i < desType.Elem().NumField(); i++ {

		field1Name := desType.Elem().Field(i).Name

		field1Type, _ := desType.Elem().FieldByName(field1Name)
		field1Value := desValue.Elem().FieldByName(field1Name)

		field1Tag := desType.Elem().Field(i).Tag.Get("ini")
		section := cfg.Section(field1Tag)

		for j := 0; j < field1Value.NumField(); j++ {

			field2Type := field1Type.Type.Field(j)
			field2Value := field1Value.Field(j)

			field2Tag := field2Type.Tag.Get("ini")
			fields := section.Field(field2Tag)

			switch field2Value.Kind() {
			case reflect.String:
				field2Value.SetString(fields.String())
			case reflect.Int64:
				field2Value.SetInt(fields.Int64())
			case reflect.Float64:
				field2Value.SetFloat(fields.Float64())
			}
		}
	}

	return nil
}

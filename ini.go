package ini

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

// Cfg contains a Section map parsing from the configuration file.
type Cfg struct {
	sections map[string]*Section
}

// Load loads the configuration file according to path
// and returns a Cfg struct containing all the configuration.
func Load(path string) (*Cfg, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	reader := bufio.NewScanner(f)
	section := new(Section)
	sections := make(map[string]*Section)
	for reader.Scan() {
		text := strings.TrimSpace(reader.Text())
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
			name := text[1:strings.LastIndexByte(text, ']')]
			_, exist := sections[name]
			if exist {
				return nil, fmt.Errorf("duplicate section %s", name)
			}
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
	_ = f.Close()

	return &Cfg{sections}, nil
}

// Section gets the appointed section.
// If the section is not exist, an empty Section will be returned instead of nil pointer.
func (cfg *Cfg) Section(name string) *Section {
	s, exist := cfg.sections[name]
	if !exist {
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
	field, exist := s.fields[name]
	if !exist {
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

// Bind offers more easy way to get configuration. Bind also calls the
// Load and bind value to the field according to the struct tag ini.
func Bind(path string, des any) error {
	t := reflect.TypeOf(des)
	if t.Kind() != reflect.Ptr {
		return fmt.Errorf("des needs to be a pointer, but gets %s", t.Kind().String())
	}
	t = t.Elem()

	v := reflect.ValueOf(des).Elem()
	if v.Kind() != reflect.Struct {
		return fmt.Errorf("des needs to be a struct, but gets %s", v.Kind().String())
	}

	cfg, err := Load(path)
	if err != nil {
		return err
	}

	for i := 0; i < t.NumField(); i++ {
		ft := t.Field(i)
		fv := v.Field(i)
		section := cfg.Section(ft.Tag.Get("ini"))
		for j := 0; j < fv.NumField(); j++ {
			subFT := ft.Type.Field(j)
			subFV := fv.Field(j)
			field := section.Field(subFT.Tag.Get("ini"))

			if subFV.IsValid() && subFV.CanSet() {
				switch subFV.Kind() {
				case reflect.String:
					subFV.SetString(field.String())
				case reflect.Int64:
					subFV.SetInt(field.Int64())
				case reflect.Float64:
					subFV.SetFloat(field.Float64())
				}
			}
		}
	}

	return nil
}

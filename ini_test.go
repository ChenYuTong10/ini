package ini

import (
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestLoad(t *testing.T) {
	cfg, err := Load("./example.ini")
	if err != nil {
		log.Fatalln(err)
	}

	// ########## person ##########
	assert.Equal(t, "zhangsan", cfg.Section("person").Field("Name").String())
	assert.Equal(t, int64(18), cfg.Section("person").Field("Age").Int64())
	assert.Equal(t, 180.0, cfg.Section("person").Field("Height").Float64())
	assert.Equal(t, 64.5, cfg.Section("person").Field("Weight").Float64())

	// ########## db ##########
	assert.Equal(t, "foo", cfg.Section("db").Field("User").String())
	assert.Equal(t, "123456", cfg.Section("db").Field("Passwd").String())
	assert.Equal(t, "127.0.0.1:3306", cfg.Section("db").Field("Addr").String())
	assert.Equal(t, "dev", cfg.Section("db").Field("DBName").String())

	// ########## nonexistent ##########
	assert.Equal(t, "", cfg.Section("nonexistent").Field("foo1").String())
	assert.Equal(t, float64(0), cfg.Section("nonexistent").Field("foo2").Float64())
	assert.Equal(t, int64(0), cfg.Section("nonexistent").Field("foo3").Int64())

}

func TestBind(t *testing.T) {
	type Foo struct {
		DB struct {
			User   string `ini:"User"`
			Passwd string `ini:"Passwd"`
			Addr   string `ini:"Addr"`
			DBName string `ini:"DBName"`
		} `ini:"db"`
		Person struct {
			Name   string  `ini:"Name"`
			Age    int64   `ini:"Age"`
			Height float64 `ini:"Height"`
			Weight float64 `ini:"Weight"`
		} `ini:"person"`
	}

	foo := &Foo{}
	if err := Bind("./example.ini", foo); err != nil {
		log.Fatalln(err)
	}

	// ########## person ##########
	assert.Equal(t, "zhangsan", foo.Person.Name)
	assert.Equal(t, int64(18), foo.Person.Age)
	assert.Equal(t, 180.0, foo.Person.Height)
	assert.Equal(t, 64.5, foo.Person.Weight)

	// ########## db ##########
	assert.Equal(t, "foo", foo.DB.User)
	assert.Equal(t, "123456", foo.DB.Passwd)
	assert.Equal(t, "127.0.0.1:3306", foo.DB.Addr)
	assert.Equal(t, "dev", foo.DB.DBName)

}

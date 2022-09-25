# ini

[![GitHub Language](https://img.shields.io/badge/Go-reference-blue)](https://go.dev)
[![GitHub license](https://img.shields.io/github/license/ChenYuTong10/ini)](https://github.com/ChenYuTong10/ini/blob/main/LICENSE)


An INI file parser which supports multipart way to parse the configuration.

## Features

- Chain-liked calling API to get the configuration
- Using `struct tag` to bind specificly with struct field.

## Example

The First way to parse ini file is calling `Load` API.
`Load` API will return a struct if it is called successfully and you can apply a series of method like chain to get the specific configuration value.

```Golang
package example

import (
    "fmt"
    "log"

    "github.com/ChenYuTong10/ini"
)

func Example() {
    cfg, err := Load("example.ini")
    if err != nil {
        log.Fatalln(err)
    }

    name := cfg.Section("person").Field("Name").String()
    age := cfg.Section("person").Field("Age").Int64()
    height := cfg.Section("person").Field("Height").Float64()
    
    fmt.Println("name:", name)
    fmt.Println("age:", age)
    fmt.Println("height:", height)
}
```

What's more, you can use `struct tag` to bind the struct field with the configuration in ini file. It can avoid getting configurations one by one.

```Golang
package example

import (
    "fmt"
    "log"

    "github.com/ChenYuTong10/ini"
)

func Example() {
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

    name := foo.Person.Name
    age := foo.Person.Age
    height := foo.Person.Height

    fmt.Println("name:", name)
    fmt.Println("age:", age)
    fmt.Println("height:", height)
}
```

## Future

⬜ Not Start ⌛ Processing ✅ Finished

- Two ways to parse the ini file ✅
- Check for boundary case ✅
- Test more cases ⬜
package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/BurntSushi/toml"
)

// 目前暂不支持comment
type tomlConfig struct {
	Title      string
	Desc       string `toml:"desc"`
	Integers   []int
	Floats     []float64
	Times      []time.Time
	Duration   []time.Duration
	Distros    []distro
	Servers    map[string]server
	Characters map[string][]struct {
		Name string
		Rank string
	}
}

type server struct {
	IP       string
	Hostname string
	Enabled  bool
}

type distro struct {
	Name     string
	Packages string
}

type fmtTime struct {
	time.Time
}

func (t fmtTime) String() string {
	f := "2006-01-02 15:04:05.999999999"
	if t.Time.Hour() == 0 {
		f = "2006-01-02"
	}
	if t.Time.Year() == 0 {
		f = "15:04:05.999999999"
	}
	if t.Time.Location() == time.UTC {
		f += " UTC"
	} else {
		f += " -0700"
	}
	return t.Time.Format(`"` + f + `"`)
}

func main() {
	f := "config.toml"
	if _, err := os.Stat(f); err != nil {
		log.Fatalln(0, err)
	}

	var config tomlConfig
	meta, err := toml.DecodeFile(f, &config) // 从文件读取
	if err != nil {
		log.Fatalln(2, err)
	}

	fmt.Println(meta.Keys())

	fmt.Println(config.Integers)

	fmt.Println(config.Times)

	config.Desc = "hello world"

	buf := bytes.NewBuffer([]byte{})

	err = toml.NewEncoder(buf).Encode(&config) // 将对象编码为TOML
	if err != nil {
		log.Fatalln(3, err)
	}

	os.WriteFile("new.toml", buf.Bytes(), 0666)
}

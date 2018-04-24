//package bsdbToXLSX
package main

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/url"
	"os"
)

//holds info pertaining to database and query action
type QueryStuff struct {
	DatabaseType string `yaml:"databaseType"`
	Database     string `yaml:"database"'`
	QueryText    string `yaml:"queryText"`
	Host         string `yaml:"host"`
	User         string `yaml:"user"`
	Pass         string `yaml:"pass"`
	ConnStr      string //*url.URL
}

//Reads yaml databse setup file
//returns raw slice ob bytes
func readConfig(filename string) []byte {
	fileInfo, err := os.Stat(filename)

	if err != nil {
		log.Fatal(err)
	}
	mode := fileInfo.Mode()
	var rightMode os.FileMode = 0600
	if mode != rightMode {
		err := errors.New(fmt.Sprintf("please set permissions on %s", filename))
		log.Fatal(err)
	}
	raw, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	return raw
}

//takes raw bytes from readConfig and marshals them to QueryStuff
//returns poiter to QuertStuff
func unMarshal(raw []byte) *QueryStuff {
	var queryStuff QueryStuff
	err := yaml.Unmarshal(raw, &queryStuff)
	if err != nil {
		fmt.Println(err)
	}
	return &queryStuff
}

func (input *QueryStuff) buildConnString() {
	switch input.DatabaseType {
	case "sqlite":
		input.ConnStr = input.Host
	default:
		u := &url.URL{
			Scheme: input.DatabaseType,
			User:   url.UserPassword(input.User, input.Pass),
			Host:   input.Host,
		}
		input.ConnStr = fmt.Sprintf("%s?database=%s", u, input.Database)
	}
}

//sets up query to be executed. Query txt is nil unless included in config file
//set query text with setQueryText() method
func QueryConfig(fileName string) *QueryStuff {
	in := readConfig(fileName)
	queryObj := unMarshal(in)
	queryObj.buildConnString()
	return queryObj
}

func main() {
	result := QueryConfig("./test.yaml")
	fmt.Println(result.ConnStr)

}

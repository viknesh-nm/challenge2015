package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"reflect"
)

type BuffInfo struct {
	Name   string   `json:"name"`
	URL    string   `json:"url"`
	Type   string   `json:"type"`
	Movies []Common `json:"movies"`
	Cast   []Common `json:"cast"`
	Crew   []Common `json:"crew"`
}

type Common struct {
	Name string `json:"name"`
	URL  string `json:"url"`
	Role string `json:"role"`
}

type BuffData struct {
	Movie, FirstURL, FirstRole, SecondURL, SecondRole string
}

type MovieBuff struct {
	SourceURL, DestinationURL string
	Visit                     []string
	FirstURL, SecondURL       *BuffInfo
	PBackMovies               map[string]Common
	Visited, VisitedPerson    map[string]bool
	Link                      map[string]BuffData
}

var SourceURL = "http://data.Moviebuff.com/"
var data MovieBuff

// checkErr - checks & prints the error
func checkErr(err error) {
	if err != nil {
		fmt.Print(err)
	}
}

// checkFatalErr - checks, prints the error & stops the statement
func checkFatalErr(err error) {
	if err != nil {
		log.Fatalln(err.Error())
	}
}

// getField gets a field from a struct
func getField(structure interface{}, field string) interface{} {
	return reflect.Indirect(reflect.ValueOf(structure)).FieldByName(field).Interface()
}

// argCheck - checks the arguments whether it matches
func argCheck(args []string) {
	fmt.Println(args)
	if len(args) != 2 || args[0] == args[1] {
		log.Fatalln("Arguments mis-match")
	}
}

// personData - gets and stores the value in a data variable
func personData(firstURL, secondURL string) error {
	data1, _ := requestData(firstURL)
	data2, _ := requestData(secondURL)

	if len(data1.Movies) > len(data2.Movies) {
		data.SourceURL, data.DestinationURL = secondURL, firstURL
		data.FirstURL, data.SecondURL = data2, data1
	} else {
		data.SourceURL, data.DestinationURL = firstURL, secondURL
		data.FirstURL, data.SecondURL = data1, data2
	}

	for _, movie := range data.SecondURL.Movies {
		data.PBackMovies[movie.URL] = movie
	}

	data.Visit = append(data.Visit, data.SourceURL)
	data.Visited[data.SourceURL] = true

	return nil
}

// compareRelation - compares the two URL relation
func compareRelation() ([]BuffData, error) {
	var buffData []BuffData
	for true {
		for _, person := range data.Visit {
			firstURL, err := requestData(person)
			checkErr(err)

			for _, resPersonMovie := range firstURL.Movies {
				if data.PBackMovies[resPersonMovie.URL].URL == resPersonMovie.URL {
					if _, found := data.Link[firstURL.URL]; found {
						buffData = append(buffData, data.Link[firstURL.URL], BuffData{resPersonMovie.Name, firstURL.Name, resPersonMovie.Role, data.SecondURL.Name, data.PBackMovies[resPersonMovie.URL].Role})
					} else {
						buffData = append(buffData, BuffData{resPersonMovie.Name, firstURL.Name, resPersonMovie.Role, data.SecondURL.Name, data.PBackMovies[resPersonMovie.URL].Role})
					}
					return buffData, nil
				}

				if data.Visited[resPersonMovie.URL] {
					continue
				}
				data.Visited[resPersonMovie.URL] = true
				resPersonMoviedetail, err := requestData(resPersonMovie.URL)
				checkErr(err)

				a := []string{"Cast", "Crew"}
				for _, cast_crew := range a {
					field := (getField(resPersonMoviedetail, cast_crew)).([]Common)
					for _, resPersonMovie := range field {

						if data.Visited[resPersonMovie.URL] {
							continue
						}
						data.Visited[resPersonMovie.URL] = true
						data.Visit = append(data.Visit, resPersonMovie.URL)
						data.Link[resPersonMovie.URL] = BuffData{resPersonMovie.Name, firstURL.Name, resPersonMovie.Role, resPersonMovie.Name, resPersonMovie.Role}
					}
				}
			}
		}
	}
	return buffData, nil
}

// requestData - gets and parse the json response
func requestData(param string) (*BuffInfo, error) {
	resp, err := http.Get(SourceURL + param)
	checkErr(err)
	defer resp.Body.Close()

	result, err := ioutil.ReadAll(resp.Body)
	checkErr(err)

	var data BuffInfo
	json.Unmarshal(result, &(data))

	return &data, nil
}

// main - process starts from here
func main() {
	args := os.Args[1:]
	fmt.Println(args)
	if len(args) != 2 || args[0] == args[1] {
		log.Fatalln("Arguments mis-match")
	}

	data.Link = make(map[string]BuffData)
	data.PBackMovies = make(map[string]Common)
	data.Visited, data.VisitedPerson = make(map[string]bool), make(map[string]bool)

	err := personData(args[0], args[1])
	checkFatalErr(err)

	relation, err := compareRelation()
	checkFatalErr(err)

	fmt.Println("Degree of separation:", len(relation))
	for _, relatedData := range relation {
		fmt.Println("Movie:", relatedData.Movie)
		fmt.Println("--", relatedData.FirstRole, ":", relatedData.FirstURL)
		fmt.Println("--", relatedData.SecondRole, ":", relatedData.SecondURL)
	}
}

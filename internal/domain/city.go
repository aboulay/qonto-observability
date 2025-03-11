package domain

import (
	"os"
	"bufio"
	"fmt"
	"strings"
)

type City struct {
	Name string
	Country string
}

func RetrieveCitiesFromFile(filename string) []*City {
	cityFile, err := os.Open("./cities.txt")
	if err != nil {
		fmt.Println(err)
	}

	fs := bufio.NewScanner(cityFile)
	fs.Split(bufio.ScanLines)

	var cities []*City
	for fs.Scan() {
		cities = append(cities, &City{
			Name: strings.Split(fs.Text(), ", ")[0],
			Country: strings.Split(fs.Text(), ", ")[1],
		})
	}

	cityFile.Close()
	return cities
}
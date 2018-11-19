package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const TAM int = 250
const N int = 250000
const T0 float32 = 1.0
const TN float32 = 0.9999
const numberOfConditions int = 91

//generate random number in a range (x,y)
func random(min, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return rand.Intn(max)
}

//generate random initial solution
func RandomList(x int) []int {
	var list []int
	i := 0
	for i < x {
		list = append(list, random(0, 2))
		i++
	}
	return list
}

//read file with the coditions
func read(filename string, size int) [numberOfConditions][3]int {

	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var validString = regexp.MustCompile(`(\-?[0-9]+)\s+(\-?[0-9]+)\s+(\-?[0-9]+) 0`)
	var l [numberOfConditions][3]int
	var i int = 0
	for scanner.Scan() {
		if validString.MatchString(scanner.Text()) {
			s := strings.Split(validString.FindString(scanner.Text()), " ")
			s = s[:3]
			for idx, k := range s {
				j, err := strconv.Atoi(k)
				if err != nil {
					panic(err)
				}
				l[i][idx] = j
			}
			i++
		}
	}

	return l
}

//func energy(initialList []int)

func main() {
	// To run a code do: go run yourfile.go
	// To compile and be able to run from any machine do: go install and execute the generated binary over the bin folder
	list := RandomList(10)
	fmt.Printf("%v", list)
}
/*
func energy(list []int , coditionList [][]int) int{
	var total int = 0
	for index, element := range coditionList{
		var cont int = 0
		for index2, subelement := range element{
			if(subelement<0 && !(list[math.Abs(subelement) - 1])){
				total++
				break
			}
			if(subelement>0 && list[math.Abs(subelement) - 1]){
				total++
				break
			}
		}
	}
	return total
}*/
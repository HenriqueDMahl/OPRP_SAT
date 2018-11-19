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
const sizeOfgene int = 20

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
func read(filename string) [numberOfConditions][3]int {

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

func abs(x int) int{
	if x < 0{
		return -x
	}
	return x
}

//Measure the energy of the candidate
func energy(candidate []int , coditionList [numberOfConditions][3]int) int{
	var total int = 0
	for _, element := range coditionList{
		for _, subelement := range element{
			if(subelement<0 && !(candidate[abs(subelement) - 1] == 1)){
				total++
				break
			}
			if(subelement>0 && candidate[abs(subelement) - 1] == 1){
				total++
				break
			}
		}
	}
	return total
}

func main() {
	list := RandomList(sizeOfgene)
	lcnf := read("uf20_01.cnf")
	fmt.Printf("%v",energy(list,lcnf))
}
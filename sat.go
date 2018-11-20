package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const TAM int = 250
const N int = 250000
const T0 float64 = 1.0
const TN float64 = 0.9999
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

//absolute
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

//temperature for annealing
func temperature(i int) float64{
	A := math.Pow(float64(N),(-2.0)) * math.Log(T0/TN)
	return T0 * math.Exp((-A*math.Pow(float64(i), 2.0)))
}

//disturbs the candidate to generate a new candidate
func disturbs(candidate []int) []int{
	new_candidate := make([]int, len(candidate))
	r := rand.Intn(len(candidate))
	copy(new_candidate,candidate)
	if new_candidate[r] == 1{
		new_candidate[r] = 0
	}else{
		new_candidate[r] = 1
	}
	return new_candidate
}

//simulated annealing
func annealing(candidate []int , coditionList [numberOfConditions][3]int) int{
	t := T0
	i := 1
	for {
		
		new_candidate := disturbs(candidate)
		deltaE := energy(candidate,coditionList) - energy(new_candidate,coditionList)
		
		if deltaE <= 0 {
			candidate = new_candidate
		}else if (float64(random(0 , 100)/100)) + (float64(random(0 , 100))/100) < math.Exp((float64(-deltaE)/t)){
			candidate = new_candidate
		}
		
		t = temperature(i)
		i++
		if(t < TN || i > N){
			return energy(candidate,coditionList)
		}
	}
}

func main() {
	candidate := RandomList(sizeOfgene)
	listCnf := read("uf20_01.cnf")
	fmt.Println("Annealing = ",annealing(candidate,listCnf))
}
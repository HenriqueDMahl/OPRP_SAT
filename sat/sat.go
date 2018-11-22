package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

var we sync.WaitGroup
var mutex sync.Mutex

const TAM int = 250
const N int = 250000

//make chunk variable
var chunk int = 0

const T0 float64 = 1.0
const TN float64 = 0.9999
const numberOfConditions int = 34080
const sizeOfgene int = 250

/*
File	 | N Conditions | Gene size |
uf20_01  |      91      |    20     |
uf100_01 |      XYZ     |    100    |
uf250_01 |      1065    |    250    |
*/

//Global list of Conditions
var coditionList [numberOfConditions][3]int

//Global temperature
var t float64 = T0

//Global list
var list []int



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
func read(filename string) {
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

	coditionList = l
}

//absolute value of int
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

//Measure the energy of the candidate
func energy(candidate []int) int {
	var total int = 0
	for _, element := range coditionList {
		for _, subelement := range element {
			if subelement < 0 && !(candidate[abs(subelement)-1] == 1) {
				total++
				break
			}
			if subelement > 0 && candidate[abs(subelement)-1] == 1 {
				total++
				break
			}
		}
	}
	return total
}

//temperature for annealing
func temperature(i int) float64 {
	A := math.Pow(float64(N), (-2.0)) * math.Log(T0/TN)
	return T0 * math.Exp((-A * math.Pow(float64(i), 2.0)))
}

//disturbs the candidate to generate a new candidate
func disturbs(candidate []int) []int {
	new_candidate := make([]int, len(candidate))
	r := rand.Intn(len(candidate))
	copy(new_candidate, candidate)
	if new_candidate[r] == 1 {
		new_candidate[r] = 0
	} else {
		new_candidate[r] = 1
	}
	return new_candidate
}

//simulated annealing
func annealing(candidate []int, id int) {
	t := T0
	flag := 0
	i := (id * chunk) + 1
	limit := (id * chunk) + chunk
	for {
		new_candidate := disturbs(candidate)
		deltaE := energy(candidate) - energy(new_candidate)

		if deltaE <= 0 {
			candidate = new_candidate
		} else if (float64(random(0, 100)))+(float64(random(0, 100))/100) < math.Exp((float64(-deltaE) / t)) {
			candidate = new_candidate
		}

		i++
		mutex.Lock()
		t = temperature(i)
		mutex.Unlock()
		if t < TN || i >= limit {
			//fmt.Printf("Temperatura = %v de %v Chunk = %v >= %v\n", t,id,limit,i)
			list = append(list, energy(candidate))
			flag = 1
		}
		if flag == 1 {
			break
		}
	}
	we.Done()
}

//standard deviation and average
func sd_a(list []int) [2]float64 {
	var average float64 = 0.0
	var sd float64 = 0.0
	var result [2]float64

	for _, item := range list {
		average += float64(item)
	}

	average = average / float64(len(list))

	for _, item := range list {
		sd += math.Pow((float64(item) - average), 2.0)
	}
	sd = math.Sqrt(sd / float64(len(list)))

	result[0] = average
	result[1] = sd

	return result
}

//pick best
func pick_best() int {
	var best int = 0
	for _, i := range list {
		if i > best {
			best = i
		}
	}
	return best
}

func main() {
	candidate := RandomList(sizeOfgene)
	read("uf250_0_34080.cnf")
	// Get the maximum of CPU cores available
	maxCores := runtime.NumCPU()
	chunk = N / maxCores
	we.Add(maxCores)
	for core := 0; core < maxCores; core++ {
		go annealing(candidate, core)
	}
	we.Wait()
	fmt.Println("Annealing = ", list)
	fmt.Println("Result = ", sd_a(list))
	fmt.Println("Best = ", pick_best())
}

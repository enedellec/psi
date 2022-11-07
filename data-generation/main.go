package main

import (
	"crypto/sha256"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
)

// generate a file with sorted SHA256 hashes
//
// the count command line parameter corresponds to the number of hashes to generate
//
// hashes are generated from integer values between 0 and count
// but if the even_only command line parameter is set, the algorithm generates hashes
// from even integer values between 0 and 2*count
func main() {

	log.Println("main() function started")

	// parse the command line options
	count := flag.Int("count", 100,
		"The number of sha256 hashes to generate")
	even_only := flag.Int("even_only", 0,
		"If equal to 1, it will hash even numbers only, starting at 0")
	flag.Parse()

	// prepare the parameters for the SHA256 hashes generation algorithm
	var profile string
	var step int
	if *even_only == 1 {
		profile = "even-only"
		step = 2
	} else {
		profile = "all"
		step = 1
	}
	max := *count * step

	// generate the SHA256 hashes
	data := make([]string, *count)
	j := 0
	for i := 0; i < max; i = i + step {
		s := sha256.Sum256([]byte(strconv.Itoa(i)))
		data[j] = fmt.Sprintf("%x\n", s)
		j++
	}

	// sort the hashes
	sort.Strings(data)

	// save the generated hashes in a file
	filename := fmt.Sprintf("data-%v-%v.csv", *count, profile)
	file, err := os.Create(filename)
	if err != nil {
		log.Fatalln(err)
		return
	}
	defer file.Close()

	for _, row := range data {
		file.WriteString(row)
	}

	log.Println("main() function finished")
}

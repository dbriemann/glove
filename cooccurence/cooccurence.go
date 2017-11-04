package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/dbriemann/glove"
)

func main() {
	corpus := flag.String("corpus", "", "The path to the corpus text file.")
	vocab := flag.String("vocab", "", "The path to the vocabulary text file.")
	output := flag.String("output", "vectors.txt", "The vectors text file.")
	//	minCount := flag.Uint("min-count", 5, "A threshold that defines the minimum times a word must occur to be kept in the vocabulary.")

	flag.Parse()

	if *corpus == "" {
		fmt.Println("You must specifiy a corpus file.\n")
		flag.PrintDefaults()
		os.Exit(1)
	}

	if *vocab == "" {
		fmt.Println("You must specifiy a vocabulary text file.\n")
		flag.PrintDefaults()
		os.Exit(1)
	}

	*output = "mm" // TODO -- remove

	v, err := glove.LoadVocabulary(*vocab)
	if err != nil {
		panic(err)
	}

	//	fmt.Println(v)

	matrix := glove.NewCooccurenceMatrix()
	err = matrix.Construct(*corpus, v, 3)
	if err != nil {
		panic(err)
	}

}

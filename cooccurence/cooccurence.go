package main

import "flag"

func main() {
	corpus := flag.String("corpus", "", "The path to the corpus text file.")
	output := flag.String("vocab", "vocab.txt", "The vocabulary file.")
	minCount := flag.Uint("min-count", 5, "A threshold that defines the minimum times a word must occur to be kept in the vocabulary.")

	flag.Parse()
}

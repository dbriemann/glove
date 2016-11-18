package main

import "github.com/zensword/glove"

func main() {
	wf, err := glove.NewWordFrequenciesFromFile("../data/text8")
	if err != nil {
		panic(err)
	}

	if err = wf.Write("../data/vocab.txt", 5); err != nil {
		panic(err)
	}
}

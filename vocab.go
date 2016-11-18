package glove

import (
	"bufio"
	"os"
	"sort"
	"strconv"
)

type Word struct {
	word string
	freq uint32
}

// WordsByFreq enables the sort interface to sort a slice of words in:
// 1) descending frequency order, then
//  2) ascending alphabetical order.
type WordsByFreq []Word

func (w WordsByFreq) Len() int      { return len(w) }
func (w WordsByFreq) Swap(i, j int) { w[i], w[j] = w[j], w[i] }
func (w WordsByFreq) Less(i, j int) bool {
	if w[i].freq == w[j].freq {
		return w[i].word < w[j].word
	}

	return w[i].freq > w[j].freq
}

// WordFrequencies stores the frequency of all words that appear in a text corpus.
type WordFrequencies struct {
	words map[string]uint32
}

// NewWordFrequenciesFromFile creates a new word to frequencies mapping from a corpus text file.
func NewWordFrequenciesFromFile(fname string) (WordFrequencies, error) {
	wf := WordFrequencies{
		words: map[string]uint32{},
	}

	f, err := os.Open(fname)
	if err != nil {
		return wf, err
	}
	defer f.Close()

	reader := bufio.NewReader(f)
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanWords)

	for scanner.Scan() {
		wf.words[scanner.Text()]++
	}

	return wf, nil
}

// Sorted returns a slice of Words ordered by frequency (descending).
func (wf *WordFrequencies) Sorted() WordsByFreq {
	wbf := make(WordsByFreq, len(wf.words))

	i := 0
	for w, f := range wf.words {
		wbf[i] = Word{
			word: w,
			freq: f,
		}
		i++
	}

	sort.Sort(wbf)
	return wbf
}

// Write saves the vocabulary to a text file excluding the words that appear < minFreq times.
func (wf *WordFrequencies) Write(fname string, minFreq uint32) error {
	f, err := os.Create(fname)
	if err != nil {
		return err
	}
	defer f.Close()

	sorted := wf.Sorted()

	for _, wo := range sorted {
		if wo.freq < minFreq {
			// From now on everything will be less.
			break
		}
		_, err := f.WriteString(wo.word + " " + strconv.FormatUint(uint64(wo.freq), 10) + "\n")
		if err != nil {
			return err
		}
	}

	return nil
}

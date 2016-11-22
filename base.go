package glove

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Word struct {
	word string
	freq uint32
}

// Vocabulary is a slice of words. It is sortable by
// 1) descending frequency order, then
//  2) ascending alphabetical order.
type Vocabulary []Word

func (w Vocabulary) Len() int      { return len(w) }
func (w Vocabulary) Swap(i, j int) { w[i], w[j] = w[j], w[i] }
func (w Vocabulary) Less(i, j int) bool {
	if w[i].freq == w[j].freq {
		return w[i].word < w[j].word
	}

	return w[i].freq > w[j].freq
}

// LoadVocabulary loads a vocabulary file and creates a new instance from the data.
func LoadVocabulary(fname string) (Vocabulary, error) {
	v := Vocabulary{}

	f, err := os.Open(fname)
	if err != nil {
		return v, err
	}
	defer f.Close()

	reader := bufio.NewReader(f)
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		fr, err := strconv.ParseUint(fields[1], 10, 32)
		if err != nil {
			return v, err
		}
		w := Word{
			word: fields[0],
			freq: uint32(fr),
		}

		v = append(v, w)
	}

	return v, nil
}

// Write saves the vocabulary to a text file excluding the words that appear < minFreq times.
func (v *Vocabulary) Write(fname string, minFreq uint32) error {
	f, err := os.Create(fname)
	if err != nil {
		return err
	}
	defer f.Close()

	for _, wo := range *v {
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

// WordFrequencies stores the frequency of all words that appear in a text corpus.
type WordFrequencies struct {
	words map[string]uint32
}

// NewWordFrequenciesFromFile creates a new word to frequencies mapping from a corpus text file.
func LoadWordFrequenciesFromFile(fname string) (WordFrequencies, error) {
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
func (wf *WordFrequencies) ToVocabulary() Vocabulary {
	vocab := make(Vocabulary, len(wf.words))

	i := 0
	for w, f := range wf.words {
		vocab[i] = Word{
			word: w,
			freq: f,
		}
		i++
	}

	sort.Sort(vocab)
	return vocab
}

// WordPair defines a pair of words where Main and Context describe the 'id' of
// the words that have context. The 'id' is the index in a Vocabulary.
type WordPair struct {
	Main    int
	Context int
}

type CooccurenceMatrix struct {
	matrix map[WordPair]float64
}

func NewCooccurenceMatrix() CooccurenceMatrix {
	cm := CooccurenceMatrix{
		matrix: map[WordPair]float64{},
	}

	return cm
}

func (cm *CooccurenceMatrix) Construct(fname string, vocab Vocabulary, windowSize int) error {
	word2id := map[string]int{}
	for id, word := range vocab {
		word2id[word.word] = id
	}

	// Read the corpus text file word by word and slide the window.
	f, err := os.Open(fname)
	if err != nil {
		return err
	}
	defer f.Close()

	contextWindow := NewContextWindow(windowSize)

	reader := bufio.NewReader(f)
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanWords)

	for scanner.Scan() {
		word := scanner.Text()
		// Find ID for word
		if id, ok := word2id[word]; ok {
			// Id was found.
			contextWindow.Slide(id)
		}
	}

	return nil
}

type ContextWindow struct {
	wsize int
	left  []int
	right []int
	word  int
}

func (cw ContextWindow) String() string {
	return fmt.Sprintf("[%v] %d [%v]", cw.left, cw.word, cw.right)
}

func NewContextWindow(size int) ContextWindow {
	cw := ContextWindow{
		wsize: size,
		word:  -1,
		left:  make([]int, size),
		right: make([]int, size),
	}
	for i := 0; i < size; i++ {
		cw.left[i] = -1
		cw.right[i] = -1
	}

	return cw
}

// Slide moves the window on the text corpus by adding the next word's ID
// to the right buffer and shifting everything to the left.
// [left buffer] word [right buffer]
// <0...wsize-1>      <0....wsize-1>
// <----------- shifting direction.
func (cw *ContextWindow) Slide(nextWordID int) {
	if cw.word == -1 {
		// Init first word, no shifting needed.
		cw.word = nextWordID
		return
	}
	for i := 0; i < cw.wsize; i++ {
		if cw.right[i] == -1 {
			// Init right value, no shifting needed.
			cw.right[i] = nextWordID
			return
		}
	}
	// Business as usual. Shift every word to the left.
	// 1. Shift left buffer to the left, dropping the first element.
	cw.left = cw.left[1:]
	// 2. Shift word into left buffer.
	cw.left = append(cw.left, cw.word)
	// 3. Set leftmost element from right buffer as center word.
	cw.word = cw.right[0]
	// 4. Shift right buffer to the left, dropping the first element.
	cw.right = cw.right[1:]
	// 5. Append new element to right buffer.
	cw.right = append(cw.right, nextWordID)
}

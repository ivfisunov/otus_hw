package hw03_frequency_analysis //nolint:golint,stylecheck

import (
	"sort"
	"strings"
)

type kv struct {
	key   string
	value int
}

const initialLength = 10

func Top10(text string) []string {
	splitedText := strings.Fields(text)
	if len(splitedText) == 0 {
		return []string{}
	}

	frequencyMap := make(map[string]int)
	for _, word := range splitedText {
		frequencyMap[word]++
	}

	freqStore := make([]kv, 0, len(frequencyMap))
	for k, v := range frequencyMap {
		freqStore = append(freqStore, kv{k, v})
	}
	sort.Slice(freqStore, func(i, j int) bool {
		return freqStore[i].value > freqStore[j].value
	})

	var returningLength = initialLength
	if returningLength > len(freqStore) {
		returningLength = len(freqStore)
	}

	var result = []string{}
	for _, v := range freqStore[:returningLength] {
		result = append(result, v.key)
	}

	return result
}

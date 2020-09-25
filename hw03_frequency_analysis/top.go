package hw03_frequency_analysis //nolint:golint,stylecheck

import (
	"sort"
	"strings"
)

type kv struct {
	key string
	value int
}

func Top10(text string) []string {
	if text == "" {
		return []string{}
	}
	splitedText := strings.Fields(text)

	frequencyMap := make(map[string]int)
	for _, word := range splitedText {
		frequencyMap[word]++
	}

	var freqStore []kv
	for k, v := range frequencyMap {
		freqStore = append(freqStore, kv{k, v})
	}
	sort.Slice(freqStore, func(i, j int) bool {
		return freqStore[i].value > freqStore[j].value
	})
	freqStore = freqStore[:10]

	var result []string
	for _, v := range freqStore {
		result = append(result, v.key)
	}

	return result
}

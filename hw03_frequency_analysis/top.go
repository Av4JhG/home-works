package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

// Top10 returns 10 most used words in form of slice.
func Top10(text string) []string {
	// Нормализуем текст.
	normalText := normalizeText(text)
	// Преобразовываем строку в слайс.
	words := strings.Fields(normalText)
	// Создаем мап для подсчета частоты слов.
	freq := make(map[string]int)

	// Подсчитываем частоту каждого слова.
	for _, word := range words {
		freq[word]++
	}

	// Создаем слайс уникальных слов.
	uniqueWords := make([]string, 0, len(freq))
	for word := range freq {
		uniqueWords = append(uniqueWords, word)
	}

	// Сортируем слова по частоте (убывание) и лексикографически (если частоты равны).
	sort.Slice(uniqueWords, func(i, j int) bool {
		if freq[uniqueWords[i]] == freq[uniqueWords[j]] {
			return uniqueWords[i] < uniqueWords[j]
		}
		return freq[uniqueWords[i]] > freq[uniqueWords[j]]
	})

	// Возвращаем топ-10 слов.
	if len(uniqueWords) > 10 {
		return uniqueWords[:10]
	}
	return uniqueWords
}

func normalizeText(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

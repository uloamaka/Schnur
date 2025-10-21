package schnur

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	stringStore = make(map[string]Stringz)
	mu          sync.Mutex
)

type Stringz struct {
	Id    string `json:"id"`
	Value      string   `json:"value"`
	Properties  Properties  `json:"properties"`
	CreatedAt string  `json:"created_at"`
}

type Properties struct {
	Length    int `json:"length"`
	IsPalindrome bool `json:"is_palindrome"`
	UniqueCharacters int `json:"unique_characters"`
	WordCount int `json:"word_count"`
	SHA256Hash string `json:"sha256_hash"`
	CharacterFrequencyMap map[string]int `json:"character_frequency_map"`
}


func AnalyzeString(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, `{"error":"Method Not Allowed"}`, http.StatusMethodNotAllowed)
        return
    }

    w.Header().Set("Content-Type", "application/json")

    var req struct {
        Value interface{} `json:"value"`
    }

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, `{"error":"Invalid request body or missing "value" field"}`, http.StatusBadRequest)
        return
    }


    strVal, ok := req.Value.(string)
    if !ok {
		http.Error(w, `{"error":"Invalid data type for "value" (must be string)"}`, http.StatusUnprocessableEntity)
		return
	}

	if trimSpace(strVal) == "" {
		http.Error(w, `{"error":"Invalid request body or missing "value" field"}`, http.StatusBadRequest)
		return
	}


    hash := hashString(strVal)

    if _, exists := stringStore[hash]; exists {
        http.Error(w, `{"error":"String already exists in the system"}`, http.StatusConflict)
        return
    }

    now := time.Now().UTC().Format("2006-01-02T15:04:05Z")


    strData := Stringz{
        Id:        hash,
        Value:     strVal,
        CreatedAt: now,
        Properties: Properties{
            Length:                stringLength(strVal),
            IsPalindrome:          isPalindrome(strVal),
            UniqueCharacters:      uniqueCharCount(strVal),
            WordCount:             wordCount(strVal),
            SHA256Hash:            hash,
            CharacterFrequencyMap: charFrequency(strVal),
        },
    }

    stringStore[hash] = strData

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(strData)
}


func GetString(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"Method Not Allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 || pathParts[2] == "" {
		http.Error(w, `{"error":"Missing string value in path"}`, http.StatusBadRequest)
		return
	}

	value := pathParts[2]
	hash := hashString(value)

	data, exists := stringStore[hash]
	if !exists {
		http.Error(w, `{"error":"String does not exist in the system"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

func FilterString(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"Method Not Allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	q := r.URL.Query()

	// Parse filters
	var (
		isPalindromeFilter *bool
		minLength, maxLength, wordCountFilter *int
		containsChar string
	)

	if val := q.Get("is_palindrome"); val != "" {
		if val == "true" {
			t := true
			isPalindromeFilter = &t
		} else if val == "false" {
			f := false
			isPalindromeFilter = &f
		} else {
			http.Error(w, `{"error":"Invalid query parameter values or types"}`, http.StatusBadRequest)
			return
		}
	}

	if val := q.Get("min_length"); val != "" {
		v, err := strconv.Atoi(val)
		if err != nil {
			http.Error(w, `{"error":"Invalid query parameter values or types"}`, http.StatusBadRequest)
			return
		}
		minLength = &v
	}

	if val := q.Get("max_length"); val != "" {
		v, err := strconv.Atoi(val)
		if err != nil {
			http.Error(w, `{"error":"Invalid query parameter values or types"}`, http.StatusBadRequest)
			return
		}
		maxLength = &v
	}

	if val := q.Get("word_count"); val != "" {
		v, err := strconv.Atoi(val)
		if err != nil {
			http.Error(w, `{"error":"Invalid query parameter values or types"}`, http.StatusBadRequest)
			return
		}
		wordCountFilter = &v
	}

	if val := q.Get("contains_character"); val != "" {
		containsChar = strings.ToLower(val)
		if len(containsChar) != 1 {
			http.Error(w, `{"error":"Invalid query parameter values or types"}`, http.StatusBadRequest)
			return
		}
	}

	// Filter strings
	var filtered []Stringz
	for _, s := range stringStore {
		p := s.Properties

		if isPalindromeFilter != nil && p.IsPalindrome != *isPalindromeFilter {
			continue
		}
		if minLength != nil && p.Length < *minLength {
			continue
		}
		if maxLength != nil && p.Length > *maxLength {
			continue
		}
		if wordCountFilter != nil && p.WordCount != *wordCountFilter {
			continue
		}
		if containsChar != "" && !strings.Contains(strings.ToLower(s.Value), containsChar) {
			continue
		}

		filtered = append(filtered, s)
	}

	filters := map[string]interface{}{}

	if isPalindromeFilter != nil {
		filters["is_palindrome"] = *isPalindromeFilter
	}
	if minLength != nil {
		filters["min_length"] = *minLength
	}
	if maxLength != nil {
		filters["max_length"] = *maxLength
	}
	if wordCountFilter != nil {
		filters["word_count"] = *wordCountFilter
	}
	if containsChar != "" {
		filters["contains_character"] = containsChar
	}

	response := map[string]interface{}{
		"count":            len(filtered),
		"data":             filtered,
		"filters_applied":  filters,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func applyFilters(filters map[string]interface{}) []Stringz {
	var results []Stringz

	for _, strObj := range stringStore {
		s := strObj.Value
		props := strObj.Properties

		// Apply filters
		if val, ok := filters["is_palindrome"].(bool); ok && props.IsPalindrome != val {
			continue
		}
		if val, ok := filters["min_length"].(int); ok && props.Length < val {
			continue
		}
		if val, ok := filters["max_length"].(int); ok && props.Length > val {
			continue
		}
		if val, ok := filters["word_count"].(int); ok && props.WordCount != val {
			continue
		}
		if val, ok := filters["contains_character"].(string); ok && !strings.Contains(strings.ToLower(s), strings.ToLower(val)) {
			continue
		}

		results = append(results, strObj)
	}

	return results
}
func SearchString(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	if query == "" {
		http.Error(w, `{"error":"unable to parse natural language query"}`, http.StatusBadRequest)
		return
	}

	lq := strings.ToLower(query)
	filters := make(map[string]interface{})

	switch {
	case strings.Contains(lq, "single word") && strings.Contains(lq, "palindromic"):
		filters["word_count"] = 1
		filters["is_palindrome"] = true

	case strings.Contains(lq, "longer than"):
		var n int
		_, err := fmt.Sscanf(lq, "strings longer than %d", &n)
		if err != nil {
			http.Error(w, `{"error":"unable to parse natural language query"}`, http.StatusBadRequest)
			return
		}
		filters["min_length"] = n + 1

	case strings.Contains(lq, "palindromic") && strings.Contains(lq, "first vowel"):
		filters["is_palindrome"] = true
		filters["contains_character"] = "a"

	case strings.Contains(lq, "containing the letter"):
		var c string
		_, err := fmt.Sscanf(lq, "strings containing the letter %s", &c)
		if err != nil {
			http.Error(w, `{"error":"unable to parse natural language query"}`, http.StatusBadRequest)
			return
		}

		c = strings.Trim(c, "\".,")
		filters["contains_character"] = c

	default:
		http.Error(w, `{"error":"unable to parse natural language query"}`, http.StatusBadRequest)
		return
	}

	if _, ok := filters["word_count"]; ok && filters["min_length"] != nil {
    	http.Error(w, `{"error":"Query parsed but resulted in conflicting filters"}`, http.StatusUnprocessableEntity)
    	return
	}


	filtered := applyFilters(filters)

	resp := map[string]interface{}{
		"data":  filtered,
		"count": len(filtered),
		"interpreted_query": map[string]interface{}{
			"original":       query,
			"parsed_filters": filters,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func DeleteString(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, `{"error":"Method Not Allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 || pathParts[2] == "" {
		http.Error(w, `{"error":"Missing string value in path"}`, http.StatusBadRequest)
		return
	}

	value := pathParts[2]
	hash := hashString(value)

	mu.Lock()
	defer mu.Unlock()

	if _, exists := stringStore[hash]; !exists {
		http.Error(w, `{"error":"String does not exist in the system"}`, http.StatusNotFound)
		return
	}

	delete(stringStore, hash)
	w.WriteHeader(http.StatusNoContent)
}



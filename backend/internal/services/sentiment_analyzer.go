package services

import (
	"strings"
	"unicode"
)

// SentimentResult represents the result of sentiment analysis
type SentimentResult struct {
	Score      float64  `json:"score"`      // -1.0 to +1.0 (negative to positive)
	Category   string   `json:"category"`   // "positive", "negative", "neutral"
	Confidence float64  `json:"confidence"` // 0.0 to 1.0
	Keywords   []string `json:"keywords"`   // Words that influenced the score
}

// SentimentAnalyzer analyzes text sentiment using keyword matching
type SentimentAnalyzer struct {
	positiveKeywords map[string]float64
	negativeKeywords map[string]float64
	neutralKeywords  map[string]float64
}

// NewSentimentAnalyzer creates a new sentiment analyzer instance
func NewSentimentAnalyzer() *SentimentAnalyzer {
	return &SentimentAnalyzer{
		positiveKeywords: map[string]float64{
			// English - General Positive
			"good": 0.7, "great": 0.8, "excellent": 0.9, "amazing": 0.9,
			"wonderful": 0.8, "fantastic": 0.8, "outstanding": 0.8,
			"successful": 0.8, "effective": 0.7, "efficient": 0.7,
			"improved": 0.6, "better": 0.6, "best": 0.7,
			"helpful": 0.6, "supportive": 0.6, "encouraging": 0.7,

			// English - COVID-19 Positive
			"recovery": 0.8, "recovered": 0.8, "healing": 0.7,
			"vaccine": 0.7, "vaccination": 0.7, "immunity": 0.6,
			"hope": 0.8, "optimistic": 0.7, "positive": 0.8,
			"decline": 0.6, "decrease": 0.6, "dropping": 0.6,
			"control": 0.6, "contained": 0.7, "stabilized": 0.6,
			"treatment": 0.6, "cure": 0.7, "prevention": 0.6,

			// Indonesian - General Positive
			"baik": 0.7, "bagus": 0.7, "hebat": 0.8, "luar biasa": 0.9,
			"berhasil": 0.8, "sukses": 0.8, "efektif": 0.7,
			"meningkat": 0.6, "lebih baik": 0.6, "terbaik": 0.7,
			"membantu": 0.6, "mendukung": 0.6, "mendorong": 0.7,

			// Indonesian - COVID-19 Positive
			"sembuh": 0.8, "pulih": 0.8, "vaksin": 0.7, "imunisasi": 0.7,
			"harapan": 0.8, "optimis": 0.7, "positif": 0.8,
			"menurun": 0.6, "berkurang": 0.6, "terkendali": 0.7,
			"pengobatan": 0.6, "penyembuhan": 0.7, "pencegahan": 0.6,
		},
		negativeKeywords: map[string]float64{
			// English - General Negative
			"bad": -0.7, "terrible": -0.8, "awful": -0.8, "horrible": -0.9,
			"worst": -0.8, "failed": -0.8, "disaster": -0.9,
			"problem": -0.6, "issue": -0.6, "concern": -0.5,
			"worry": -0.6, "fear": -0.7, "anxiety": -0.7,
			"difficult": -0.5, "hard": -0.5, "challenging": -0.4,

			// English - COVID-19 Negative
			"death": -0.9, "died": -0.9, "lethal": -0.9,
			"infection": -0.6, "infected": -0.6, "contagious": -0.6,
			"spread": -0.5, "outbreak": -0.7, "pandemic": -0.6,
			"lockdown": -0.6, "quarantine": -0.6, "isolation": -0.6,
			"crisis": -0.7, "emergency": -0.6, "danger": -0.7,
			"severe": -0.6, "critical": -0.7, "serious": -0.6,

			// Indonesian - General Negative
			"buruk": -0.7, "jelek": -0.7, "mengerikan": -0.8, "mengkhawatirkan": -0.7,
			"gagal": -0.8, "masalah": -0.6, "kekhawatiran": -0.6,
			"cemas": -0.6, "takut": -0.7, "khawatir": -0.6,
			"sulit": -0.5, "berat": -0.5, "menantang": -0.4,

			// Indonesian - COVID-19 Negative
			"meninggal": -0.9, "mati": -0.9, "fatal": -0.9,
			"terinfeksi": -0.6, "menular": -0.6, "penyebaran": -0.5,
			"wabah": -0.7, "pandemi": -0.6, "krisis": -0.7,
			"darurat": -0.6, "bahaya": -0.7, "mengancam": -0.6,
			"parah": -0.6, "kritis": -0.7, "serius": -0.6,
		},
		neutralKeywords: map[string]float64{
			// English - Neutral
			"update": 0.0, "report": 0.0, "statistics": 0.0,
			"information": 0.0, "news": 0.0, "announcement": 0.0,
			"daily": 0.0, "weekly": 0.0, "monthly": 0.0,
			"confirmed": 0.0, "reported": 0.0, "announced": 0.0,
			"case": 0.0, "number": 0.0, "count": 0.0,

			// Indonesian - Neutral
			"laporan": 0.0, "statistik": 0.0,
			"informasi": 0.0, "berita": 0.0, "pengumuman": 0.0,
			"harian": 0.0, "mingguan": 0.0, "bulanan": 0.0,
			"dikonfirmasi": 0.0, "dilaporkan": 0.0, "diumumkan": 0.0,
			"kasus": 0.0, "jumlah": 0.0, "hitung": 0.0,
		},
	}
}

// AnalyzeSentiment analyzes the sentiment of given text
func (sa *SentimentAnalyzer) AnalyzeSentiment(text string) *SentimentResult {
	if text == "" {
		return &SentimentResult{
			Score:      0.0,
			Category:   "neutral",
			Confidence: 0.0,
			Keywords:   []string{},
		}
	}

	// Clean and tokenize text
	words := sa.tokenizeText(text)

	var totalScore float64
	var foundKeywords []string
	var positiveCount, negativeCount, neutralCount int

	// Analyze each word
	for _, word := range words {
		wordLower := strings.ToLower(word)

		// Check positive keywords
		if score, exists := sa.positiveKeywords[wordLower]; exists {
			totalScore += score
			foundKeywords = append(foundKeywords, word)
			positiveCount++
		}

		// Check negative keywords
		if score, exists := sa.negativeKeywords[wordLower]; exists {
			totalScore += score
			foundKeywords = append(foundKeywords, word)
			negativeCount++
		}

		// Check neutral keywords
		if _, exists := sa.neutralKeywords[wordLower]; exists {
			neutralCount++
		}
	}

	// Calculate final score and category
	result := sa.calculateFinalSentiment(totalScore, positiveCount, negativeCount, neutralCount, len(words))
	result.Keywords = foundKeywords

	return result
}

// tokenizeText splits text into words and cleans them
func (sa *SentimentAnalyzer) tokenizeText(text string) []string {
	// Split by whitespace and punctuation
	words := strings.FieldsFunc(text, func(r rune) bool {
		return unicode.IsSpace(r) || unicode.IsPunct(r)
	})

	// Clean words (remove empty strings, normalize)
	var cleanedWords []string
	for _, word := range words {
		word = strings.TrimSpace(word)
		if word != "" && len(word) > 1 { // Skip single characters
			cleanedWords = append(cleanedWords, word)
		}
	}

	return cleanedWords
}

// calculateFinalSentiment determines the final sentiment category and confidence
func (sa *SentimentAnalyzer) calculateFinalSentiment(totalScore float64, positiveCount, negativeCount, neutralCount, totalWords int) *SentimentResult {
	// Normalize score to -1.0 to +1.0 range
	var normalizedScore float64
	if totalWords > 0 {
		normalizedScore = totalScore / float64(totalWords)
		// Cap at -1.0 and +1.0
		if normalizedScore > 1.0 {
			normalizedScore = 1.0
		} else if normalizedScore < -1.0 {
			normalizedScore = -1.0
		}
	}

	// Determine category based on score
	var category string
	var confidence float64

	if normalizedScore > 0.02 {
		category = "positive"
		confidence = sa.calculateConfidence(positiveCount, negativeCount, totalWords)
	} else if normalizedScore < -0.02 {
		category = "negative"
		confidence = sa.calculateConfidence(negativeCount, positiveCount, totalWords)
	} else {
		category = "neutral"
		confidence = sa.calculateConfidence(neutralCount, positiveCount+negativeCount, totalWords)
	}

	return &SentimentResult{
		Score:      normalizedScore,
		Category:   category,
		Confidence: confidence,
	}
}

// calculateConfidence calculates confidence level based on keyword distribution
func (sa *SentimentAnalyzer) calculateConfidence(primaryCount, secondaryCount, totalWords int) float64 {
	if totalWords == 0 {
		return 0.0
	}

	// Base confidence on how many words were classified
	classifiedWords := primaryCount + secondaryCount
	if classifiedWords == 0 {
		return 0.0
	}

	// Higher confidence if more words are classified and primary category dominates
	classificationRatio := float64(classifiedWords) / float64(totalWords)
	dominanceRatio := float64(primaryCount) / float64(classifiedWords)

	confidence := classificationRatio * dominanceRatio

	// Cap confidence at 1.0
	if confidence > 1.0 {
		confidence = 1.0
	}

	return confidence
}

// AnalyzeSentimentBatch analyzes multiple texts for sentiment
func (sa *SentimentAnalyzer) AnalyzeSentimentBatch(texts []string) []*SentimentResult {
	results := make([]*SentimentResult, len(texts))
	for i, text := range texts {
		results[i] = sa.AnalyzeSentiment(text)
	}
	return results
}

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type RatingWithAdvice struct {
	Rating
	IsAdvisable bool `json:"is_advisable"`
}

func GetRatingsHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	page, _ := strconv.Atoi(query.Get("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(query.Get("limit"))
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	sort := query.Get("sort")
	if sort == "" {
		sort = "time"
	}
	order := query.Get("order")
	if order != "asc" {
		order = "desc"
	}

	sql := fmt.Sprintf(`
		SELECT ticker, target_from, target_to, company, action, brokerage, rating_from, rating_to, time 
		FROM ratings
		ORDER BY %s %s
		LIMIT $1 OFFSET $2
	`, sort, order)

	rows, err := DB.Query(context.Background(), sql, limit, offset)
	if err != nil {
		http.Error(w, "Error consultando la base de datos", 500)
		return
	}
	defer rows.Close()

	var ratings []RatingWithAdvice
	for rows.Next() {
		var r Rating
		err := rows.Scan(&r.Ticker, &r.TargetFrom, &r.TargetTo, &r.Company, &r.Action, &r.Brokerage, &r.RatingFrom, &r.RatingTo, &r.Time)
		if err != nil {
			http.Error(w, "Error leyendo resultados", 500)
			return
		}
		ratings = append(ratings, RatingWithAdvice{
			Rating:      r,
			IsAdvisable: recomendBuy(r),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ratings)
}

func normalizeRating(r string) string {
	switch strings.ToLower(r) {
	case "buy", "strong buy", "outperform", "overweight", "sector outperform":
		return "buy"
	case "hold", "neutral", "equal weight", "market perform", "sector perform":
		return "hold"
	case "sell", "underperform", "underweight":
		return "sell"
	default:
		return "unknown"
	}
}

func parsePrice(price string) (float64, error) {
	p := strings.Replace(price, "$", "", 1)
	return strconv.ParseFloat(p, 64)
}

func recomendBuy(rating Rating) bool {
	targetFrom, _ := parsePrice(rating.TargetFrom)
	targetTo, _ := parsePrice(rating.TargetTo)

	increase := (targetTo - targetFrom) / targetFrom
	isBuy := normalizeRating(rating.RatingTo) == "buy"
	isUpgrade := rating.Action != "downgraded by" && rating.Action != "target lowered by"
	return isBuy && isUpgrade && increase >= 0.05
}

func GetBestRatingHandler(w http.ResponseWriter, req *http.Request) {
	query := `
		SELECT ticker, target_from, target_to, company, action, brokerage, rating_from, rating_to, time 
		FROM ratings
		WHERE LOWER(rating_to) IN ('buy', 'strong buy', 'outperform', 'overweight', 'sector outperform')
	`

	rows, err := DB.Query(context.Background(), query)
	if err != nil {
		http.Error(w, "DB query error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var bestRating Rating
	var maxIncrease float64 = -1
	today := time.Now().UTC()
	yesterday := today.AddDate(0, 0, -1)
	for rows.Next() {
		var rating Rating
		if err := rows.Scan(
			&rating.Ticker, &rating.TargetFrom, &rating.TargetTo, &rating.Company,
			&rating.Action, &rating.Brokerage, &rating.RatingFrom, &rating.RatingTo, &rating.Time,
		); err != nil {
			log.Println("Scan error:", err)
			continue
		}
		ratingDate := rating.Time.Format("2006-01-02")
		if ratingDate != today.Format("2006-01-02") && ratingDate != yesterday.Format("2006-01-02") {
			continue
		}

		from, err1 := parsePrice(rating.TargetFrom)
		to, err2 := parsePrice(rating.TargetTo)

		if err1 != nil || err2 != nil || from == 0 || to <= from {
			continue
		}

		increase := (to - from) / from

		if increase > maxIncrease {
			bestRating = rating
			maxIncrease = increase
		}
	}

	if maxIncrease < 0 {
		http.Error(w, "No suitable recommendation found for today", http.StatusNotFound)
		return
	}

	response := struct {
		Rating      Rating  `json:"rating"`
		IncreasePct float64 `json:"increase_pct"`
		IsAdvisable bool    `json:"is_advisable"`
	}{
		Rating:      bestRating,
		IncreasePct: maxIncrease,
		IsAdvisable: true,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type ApiResponse struct {
	Items    []Rating `json:"items"`
	NextPage string   `json:"next_page"`
}

func FetchAndStoreRatings(w http.ResponseWriter, r *http.Request) {
	fmt.Println("CALL TO ENDPOINT")
	nextPage := ""
	var apiURL = os.Getenv("API_SWE_URL")

	var apiKey = fmt.Sprintf("Bearer %s", os.Getenv("API_SWE_KEY"))
	for {
		url := apiURL
		if nextPage != "" {
			url += "?next_page=" + nextPage
		}

		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Authorization", apiKey)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Println("Error request:", err)
			break
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		var response ApiResponse
		if err := json.Unmarshal(body, &response); err != nil {
			log.Println("Error parse JSON:", err)
			break
		}

		for _, item := range response.Items {
			_, err := DB.Exec(context.Background(),
				`INSERT INTO ratings (ticker, target_from, target_to, company, action, brokerage, rating_from, rating_to, time)
				 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`,
				item.Ticker, item.TargetFrom, item.TargetTo, item.Company,
				item.Action, item.Brokerage, item.RatingFrom, item.RatingTo, item.Time,
			)
			if err != nil {
				log.Println("Error in insert to ratings:", err)
			}
		}

		if response.NextPage == "" {
			break
		}
		nextPage = response.NextPage
		fmt.Println("Next page:", nextPage)
	}
	fmt.Println("END GET INFO:", nextPage)
}

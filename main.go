package main

import (
	"fmt"
	"io"
	"net/http"
	"f1term/internal/selectYear"
)

func main() {
	selectYear.printYears()
	url := "https://api.openf1.org/v1/meetings?year=2022"

	// Création de la requête
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Erreur lors de la requête:", err)
		return
	}
	defer resp.Body.Close()

	// Lecture de la réponse
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Erreur lors de la lecture de la réponse:", err)
		return
	}

	fmt.Println("Réponse de l'API:")
	fmt.Println(string(body))
}

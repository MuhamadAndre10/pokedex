package main

import (
	"encoding/json"
	"errors"
	"io"
	"math/rand"
	"net/http"
	"time"

	"github.com/muhamadAndre10/pokedex/internal"
)

var cacheEntry = internal.NewCache(10 * time.Millisecond)

type locAreaRes struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func getPokeApi(url string) (locAreaRes, error) {

	// Cek apakah data sudah ada di cache?
	if cacheData, found := cacheEntry.Get(url); found {
		var data locAreaRes
		err := json.Unmarshal(cacheData, &data)
		if err != nil {
			return locAreaRes{}, err
		}
		return data, nil
	}

	// jika belum, Fetch dan set data di cache
	res, err := http.Get(url)
	if err != nil {
		return locAreaRes{}, errors.New("failed to fetch location api")
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return locAreaRes{}, errors.New("failed to read body")
	}

	cacheEntry.Set(url, body)

	var data locAreaRes

	err = json.Unmarshal(body, &data)
	if err != nil {
		return locAreaRes{}, errors.New("failed unmarshal data")
	}

	return data, nil

}

func tryCatchPokemon(baseExp int) bool {

	// algoritma
	// semakin kecil peluang semakin sulit di dapatkan. dan sebeliknya

	// Hitung peluang
	chance := min(max(100-baseExp/2, 10), 90)

	// buat angka random untuk melihat peluang kita bisa tidak mendapatkan peluang dari pokemon berdarkan baseExperience nya.
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	roll := r.Intn(100)

	return roll < chance
}

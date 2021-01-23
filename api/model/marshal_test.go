package model

import (
	"encoding/json"
	"testing"
)

func BenchmarkMarshalMapInterface(t *testing.B) {
	t.ReportAllocs()
	for i := 0; i < t.N; i++ {
		data := map[string]interface{}{
			"Name":              "Michael",
			"FavouriteLanguage": "Golang",
		}

		if _, err := json.Marshal(&data); err != nil {
			t.Error(err)
		}
	}
}

func BenchmarkMarshalMapString(t *testing.B) {
	t.ReportAllocs()

	for i := 0; i < t.N; i++ {
		data := map[string]string{
			"Name":              "Michael",
			"FavouriteLanguage": "Golang",
		}

		if _, err := json.Marshal(&data); err != nil {
			t.Error(err)
		}
	}
}

func BenchmarkMarshalStruct(t *testing.B) {
	t.ReportAllocs()
	for i := 0; i < t.N; i++ {
		data := struct {
			Name              string `json:"name"`
			FavouriteLanguage string `json:"favouriteLanguage"`
		}{
			"Michael",
			"Golang",
		}

		if _, err := json.Marshal(&data); err != nil {
			t.Error(err)
		}
	}
}

func BenchmarkUnMarshalMapInterface(t *testing.B) {
	b := map[string]string{
		"Name":              "Michael",
		"FavouriteLanguage": "Golang",
	}

	bytes, err := json.Marshal(&b)
	if err != nil {
		t.Error(err)
	}

	t.ReportAllocs()
	t.ResetTimer()

	for i := 0; i < t.N; i++ {
		b := map[string]interface{}{}
		if err := json.Unmarshal(bytes, &b); err != nil {
			t.Error(err)
		}
	}
}

func BenchmarkUnMarshalMapString(t *testing.B) {
	b := map[string]string{
		"Name":              "Michael",
		"FavouriteLanguage": "Golang",
	}

	bytes, err := json.Marshal(&b)
	if err != nil {
		t.Error(err)
	}

	t.ReportAllocs()
	t.ResetTimer()

	for i := 0; i < t.N; i++ {
		b := map[string]string{}
		if err := json.Unmarshal(bytes, &b); err != nil {
			t.Error(err)
		}
	}
}

func BenchmarkUnMarshalStruct(t *testing.B) {
	b := map[string]string{
		"name":              "Michael",
		"favouriteLanguage": "Golang",
	}

	bytes, err := json.Marshal(&b)
	if err != nil {
		t.Error(err)
	}

	t.ReportAllocs()
	t.ResetTimer()

	for i := 0; i < t.N; i++ {
		data := struct {
			Name              string `json:"name"`
			FavouriteLanguage string `json:"favouriteLanguage"`
		}{}

		if err := json.Unmarshal(bytes, &data); err != nil {
			t.Error(err)
		}
	}
}

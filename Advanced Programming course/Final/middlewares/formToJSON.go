package middlewares

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

func FormToJSONMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var reader io.Reader = r.Body
		contentType := r.Header.Get("Content-Type")
		if contentType == "application/x-www-form-urlencoded" {
			if err := r.ParseForm(); err != nil {
				http.Error(w, "Failed to parse form", http.StatusBadRequest)
				return
			}
			formData := make(map[string]string)
			for key, values := range r.PostForm {
				formData[key] = values[0]
			}
			jsonData, err := json.Marshal(formData)
			if err != nil {
				http.Error(w, "Failed to convert form to JSON", http.StatusInternalServerError)
				return
			}
			reader = bytes.NewBuffer(jsonData)
			r.Header.Set("Content-Type", "application/json")
		}
		newRequest, err := http.NewRequest(r.Method, r.URL.String(), reader)
		if err != nil {
			http.Error(w, "Failed to create new request", http.StatusInternalServerError)
			return
		}
		newRequest.Header = r.Header
		next.ServeHTTP(w, newRequest)
	})
}

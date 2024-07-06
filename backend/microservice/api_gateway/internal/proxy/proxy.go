package proxy

import (
	"fmt"
	"io"
	"net/http"
)

func ProxyRequest(target string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Proxying request to:", target+r.RequestURI)

		req, err := http.NewRequest(r.Method, target+r.RequestURI, r.Body)
		if err != nil {
			fmt.Println("Error creating new request:", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Println("Continue...")
		req.Header = r.Header
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Println("Error making request to target:", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer resp.Body.Close()
		fmt.Println("Continue...")

		for key, values := range resp.Header {
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}
		fmt.Println("Continue...")

		w.WriteHeader(resp.StatusCode)
		if _, err := io.Copy(w, resp.Body); err != nil {
			fmt.Println("Error copying response body:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		fmt.Println("Continue...")
	}
}

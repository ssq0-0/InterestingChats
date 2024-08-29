package proxy

import (
	"InterestingChats/backend/api_gateway/internal/consts"
	"InterestingChats/backend/api_gateway/internal/logger"
	"fmt"
	"io"
	"net/http"
)

func ProxyRequest(target string, log logger.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			// w.Header().Set("Access-Control-Allow-Origin", "*")
			// w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			// w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
			w.WriteHeader(http.StatusOK)
			return
		}

		req, err := http.NewRequest(r.Method, target+r.RequestURI, r.Body)
		if err != nil {
			ErrorHandler(w, http.StatusBadRequest, log, []string{consts.ErrInternalServerError}, fmt.Errorf(consts.InternalServerError, err))
			return
		}

		log.Infof("request going to %s", target+r.RequestURI)

		req.Header = r.Header
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			ErrorHandler(w, http.StatusBadRequest, log, []string{consts.ErrInternalServerError}, fmt.Errorf(consts.InternalServerError, err))
			return
		}
		defer resp.Body.Close()

		for key, values := range resp.Header {
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}

		w.WriteHeader(resp.StatusCode)
		if _, err := io.Copy(w, resp.Body); err != nil {
			ErrorHandler(w, http.StatusBadRequest, log, []string{consts.ErrInternalServerError}, fmt.Errorf(consts.InternalServerError, err))
			return
		}
	}
}

package middleware

import "net/http"

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// 🌍 Allow requests from ANY origin
		w.Header().Set("Access-Control-Allow-Origin", "*")

		// 🌍 Allow all common HTTP methods
		w.Header().Set(
			"Access-Control-Allow-Methods",
			"GET, POST, PUT, PATCH, DELETE, OPTIONS",
		)

		// 🌍 Allow all common headers
		w.Header().Set(
			"Access-Control-Allow-Headers",
			"Content-Type, Authorization, X-Requested-With",
		)

		// ❌ MUST be false / omitted when using "*"
		// w.Header().Set("Access-Control-Allow-Credentials", "true")

		// ⚡ Preflight request handling
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

package main

import (
	"net/http"
)

func (cfg *apiConfig) handlerHome(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	html := `
	<html>
		<body>
			<h1>Welcome to Chirpy!</h1>
			<p>Chirpy is a platform for sharing your thoughts with the world.</p>
			<p> Hello from Docker! I'm a Go server. </p>
			<p> Hi Docker, I pushed a new version </p>
		</body>
	</html>`
	w.Write([]byte(html))
}
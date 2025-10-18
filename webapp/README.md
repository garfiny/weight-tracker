This folder contains a minimal frontend for the weight-tracking service.

Files:
- index.html — single-page UI that talks to the backend at `/weights`.
- static/js/app.js — minimal client JS to list and create entries.
- static/css/style.css — simple styling.

How to use:
1. Start the Go server (it runs on :8080 by default):
   ```bash
   go run src/main.go
   ```
2. Open http://localhost:8080/index.html in your browser.

Integration note:
- The Go server currently doesn't serve static files. To serve this folder add a handler in `src/main.go`:

  r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("webapp/static"))))
  r.Handle("/index.html", http.FileServer(http.Dir("webapp")))

Or serve the whole folder under `/` in development.

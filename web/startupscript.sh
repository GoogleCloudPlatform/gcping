# Install Go
sudo apt-get update || exit
sudo apt install -y golang-go || exit

# Start simple HTTP server with CORS header set.
cat << EOF > main.go
package main

import (
  "fmt"
  "log"
  "net/http"
)

func main() {
  http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
    w.Header().Add("Cache-Control", "no-store")
    w.Header().Add("Access-Control-Allow-Origin", "*")
    fmt.Fprint(w, "pong")
  })
  log.Fatal(http.ListenAndServe(":80", nil))
}
EOF
sudo go run main.go

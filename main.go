package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World")
	})

	http.HandleFunc("/greet/", func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Path[len("/greet/"):]
		fmt.Fprintf(w, "Hello %s\n", name)
	})

	fmt.Println("Server started at http://localhost:8080/")
	http.ListenAndServe(":8080", nil)
<<<<<<< HEAD

=======
>>>>>>> 72a99602a3c30122c49eed8c97ceea410624ac61
}

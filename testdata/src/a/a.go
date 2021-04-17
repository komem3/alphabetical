package a

import "net/http"

// Alphabetical order
var (
	Banana = "banana"
	Apple  = "apple" // want "not sort by alphabet order"
)

var (
	Normal = "normal comment"
	Easy   = "easy comment"
)

// Alphabetical order
const (
	define = "define"
	coin   = "coin" // want "not sort by alphabet order"
)

// Alphabetical order
type (
	final    string
	elephant string // want "not sort by alphabet order"
)

func f() {
	http.HandleFunc("/z", nil)

	// Alphabetical order
	http.HandleFunc("/b", nil)
	http.HandleFunc("/c", nil)
	http.HandleFunc("/a", nil) // want "not sort by alphabet order"

	{
		// Alphabetical order
		http.HandleFunc("/", b)
		http.HandleFunc("/", c)
		http.HandleFunc("/", a) // want "not sort by alphabet order"
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Alphabetical order
		http.HandleFunc("/", e())
		http.HandleFunc("/", d()) // want "not sort by alphabet order"
	})

	if true {
		// Alphabetical order
		http.HandleFunc("/c", nil)
		http.HandleFunc("/a", nil) // want "not sort by alphabet order"
	}
}

func a(_ http.ResponseWriter, _ *http.Request) {}

func b(_ http.ResponseWriter, _ *http.Request) {}

func c(_ http.ResponseWriter, _ *http.Request) {}

func d() func(_ http.ResponseWriter, _ *http.Request) {
	return func(_ http.ResponseWriter, _ *http.Request) {}
}

func e() func(_ http.ResponseWriter, _ *http.Request) {
	return func(_ http.ResponseWriter, _ *http.Request) {}
}

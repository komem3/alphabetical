package a

import "net/http"

// Alphabetical order
var (
	Banana = "banana"
	Apple  = "apple" // want "not sort by alphabetical"
)

var (
	Normal = "normal comment"
	Easy   = "easy comment"
)

// Alphabetical order
const (
	define = "define"
	coin   = "coin" // want "not sort by alphabetical"
)

// Alphabetical order
type (
	final    string
	elephant string // want "not sort by alphabetical"
)

func f() {
	http.HandleFunc("/z", nil)

	// Alphabetical order
	http.HandleFunc("/b", nil)
	http.HandleFunc("/c", nil)
	http.HandleFunc("/a", nil) // want "not sort by alphabetical"

	{
		// Alphabetical order
		http.HandleFunc("/", b)
		http.HandleFunc("/", c)
		http.HandleFunc("/", a) // want "not sort by alphabetical"
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Alphabetical order
		http.HandleFunc("/", e())
		http.HandleFunc("/", d()) // want "not sort by alphabetical"
		http.HandleFunc("/", e())
		print(0)
	})

	if true {
		// Alphabetical order
		http.HandleFunc("/c", nil)
		http.HandleFunc("/b", nil) // want "not sort by alphabetical"
		http.HandleFunc("/a", nil) // want "not sort by alphabetical"
	}

	// Alphabetical order
	b(nil, nil)
	a(nil, nil) // want "not sort by alphabetical"
	e()

	// Alphabetical order
	e()
	d() // want "not sort by alphabetical"
	a(nil, nil)
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

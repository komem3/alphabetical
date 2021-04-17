# alphabeticalorder

Check function call and variable declarations are in alphabetical order.
The target is the statement with `// Alphabetical order` comments.

```go
// Alphabetical order
var (
	Banana = "banana"
	Apple  = "apple" // want "not sort by alphabetical"
)

func f() {
	http.HandleFunc("/z", nil)

	// Alphabetical order
	http.HandleFunc("/b", nil)
	http.HandleFunc("/a", nil) // want "not sort by alphabetical"
}
```

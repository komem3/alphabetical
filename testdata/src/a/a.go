package a

// Alphabetical order
var (
	Banana = "banana"
	Apple  = "apple" // want "not sort by alpahbet order"
)

var (
	Normal = "normal comment"
	Easy   = "easy comment"
)

// Alphabetical order
const (
	define = "define"
	coin   = "coin" // want "not sort by alpahbet order"
)

// Alphabetical order
type (
	final    string
	elephant string // want "not sort by alpahbet order"
)

package common

type DocType int

const (
	UNKNOWN DocType = -1
	CSS     DocType = 0
	HTML    DocType = 1
	JS      DocType = 2
	PLAIN   DocType = 3
)

type NewLinks struct {
	Base       string
	Discovered []string
}

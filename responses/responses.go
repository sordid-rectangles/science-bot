package responses

import (
	"errors"
	"math/rand"
	"time"
)

var simpleRes = [...]string{
	"Ummmmmmm actually that is incorrect, also don't care, also ratio",
	"No, I don't believe that's correct",
	"Hmmm, maybe speak less, it's unbecoming",
	"HA. NOPE."}

//TODO: genResponse(type string) (string, error) {}

//TODO: getSimple(ID int) (string, error) {}
func GetSimple(id int) (string, error) {
	if id >= len(simpleRes) {
		return "", errors.New("Failed to retrieve simple response, Id out of bounds")
	} else {
		return simpleRes[id], nil
	}
}

//TODO: randSimple() (string, error) {}
func RandSimple() (string, error) {
	rand.Seed(time.Now().UnixNano())
	l := len(simpleRes)
	i := rand.Intn(l)
	return GetSimple(i)
}

//TODO: genProse( seed string) string {}
func GenProse(seed string) (string, error) {

	return "NOT YET IMPLEMENTED ; (", nil
}

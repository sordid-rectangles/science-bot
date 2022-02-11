package responses

import (
	"errors"
	"math/rand"
	"time"
)

var SimpleRes = []string{
	"Ummmmmmm actually that is incorrect, also don't care, also ratio",
	"No, I don't believe that's correct",
	"Hmmm, maybe speak less, it's unbecoming",
	"HA. NOPE."}

//TODO: genResponse(type string) (string, error) {}

//TODO: getSimple(ID int) (string, error) {}
func GetSimple(id int) (string, error) {
	if id >= len(SimpleRes) {
		return "", errors.New("Failed to retrieve simple response, Id out of bounds")
	} else {
		return SimpleRes[id], nil
	}
}

//TODO: randSimple() (string, error) {}
func RandSimple() (string, error) {
	rand.Seed(time.Now().UnixNano())
	l := len(SimpleRes)
	i := rand.Intn(l)
	return GetSimple(i)
}

//TODO: genProse( seed string) string {}
func GenProse(seed string) (string, error) {

	return "NOT YET IMPLEMENTED ; (", nil
}

//TODO: addResponse(res string) {}
func AddResponse(res string) {
	if res != "" {
		SimpleRes = append(SimpleRes, res)
	}
}

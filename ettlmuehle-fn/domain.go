package function

import "time"

// Token descibes an access token
type Token struct {
	AccessToken string    `json:"access_token"`
	Retrieved   time.Time `json:"retrieved"`
}

func (t *Token) expired() bool {
	expirationDate := t.Retrieved.Add(time.Minute * 30)
	return time.Now().After(expirationDate)
}

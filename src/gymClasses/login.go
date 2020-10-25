package gymClasses

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gocolly/colly"
	"gymclassbooker.com/bwoo/config"
)

func Login(config *config.Config) ([]*http.Cookie, error) {

	c := colly.NewCollector()

	var authCookie []*http.Cookie
	var err error

	c.Post(config.GetLoginURL(),
		map[string]string{"name": config.MemberUsername, "pass": config.MemberPassword, "form_id": "user_login_block", "op": "Log in"})

	// attach callbacks after login
	c.OnResponse(func(r *colly.Response) {
		log.Println("Login: response received", r.StatusCode)
		if r.StatusCode != http.StatusOK {
			err = fmt.Errorf("Unauthorized")
			return
		}

		cookies := c.Cookies(config.GymBaseURL)
		authCookie = cookies
	})

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko)")
	})

	c.Visit(config.GymBaseURL)
	return authCookie, err
}

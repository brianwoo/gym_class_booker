package gymClasses

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gocolly/colly"
	"gymclassbooker.com/bwoo/config"
)

func UnbookAClass(classToBook GymClass, config *config.Config, cookies []*http.Cookie) error {

	c := colly.NewCollector()

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
		r.Headers.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko)")
		r.Headers.Set("X-Requested-With", "XMLHttpRequest")
		r.Headers.Set("DNT", "1")
		r.Headers.Set("Origin", config.GymBaseURL)
		r.Headers.Set("Accept", "application/json, text/javascript, */*; q=0.01")
		r.Headers.Set("Sec-Fetch-Dest", "empty")
		r.Headers.Set("Sec-Fetch-Site", "same-origin")
		r.Headers.Set("Sec-Fetch-Mode", "cors")
		r.Headers.Set("Referer", config.GetClassesURL())
		r.Headers.Set("Accept-Encoding", "gzip, deflate")
		r.Headers.Set("Accept-Language", "en-US,en;q=0.9,zh-HK;q=0.8,zh;q=0.7")
	})

	c.OnResponse(func(r *colly.Response) {
		log.Println("UnbookAClass: Response received", r.StatusCode)
	})

	// Set Cookies if we have them
	if cookies != nil {
		c.SetCookies(config.GymBaseURL, cookies)
	}

	err := c.Post(config.GymBaseURL+classToBook.GetURL(), nil)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

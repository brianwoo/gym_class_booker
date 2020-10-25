package gymClasses

import (
	"fmt"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"gymclassbooker.com/bwoo/config"
)

func GetInitFormIdsFromClasses(config *config.Config, cookies []*http.Cookie) (map[string]string, error) {

	c := colly.NewCollector()

	hiddenInputElems := make(map[string]string, 0)

	c.OnHTML(config.GymFormSelector, func(e *colly.HTMLElement) {

		e.DOM.Find("input[type='hidden']").Each(func(i int, s *goquery.Selection) {
			key, keyExists := s.Attr("name")
			value, valueExists := s.Attr("value")
			if keyExists && valueExists {
				hiddenInputElems[key] = value
			}

		})
	})

	c.OnResponse(func(r *colly.Response) {
		log.Println("GetInitFormId: Response received", r.StatusCode)
	})

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko)")
	})

	// Set Cookies if we have them
	if cookies != nil {
		c.SetCookies(config.GymBaseURL, cookies)
	}

	err := c.Visit(config.GetClassesURL())
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	//fmt.Println("End of method", hiddenInputElems)
	return hiddenInputElems, nil
}

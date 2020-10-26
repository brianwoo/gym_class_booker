package gymClasses

import (
	"fmt"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"gymclassbooker.com/bwoo/config"
)

func getPayload(epochTimestamp string, hiddenIds map[string]string) map[string]string {

	payload := map[string]string{"date_start": epochTimestamp, "op": "Find Classes"}

	for k, v := range hiddenIds {
		payload[k] = v
	}

	return payload
}

func PostClassesWithDate(epochTimestamp string, hiddenIds map[string]string, config *config.Config, cookies []*http.Cookie) (GymClasses, error) {

	c := colly.NewCollector()

	classesListed := make(GymClasses, 0)

	// Find and visit all links
	c.OnHTML(".class-title", func(e *colly.HTMLElement) {

		class := &GymClass{}
		class.SetTimestamp(epochTimestamp)

		// Go up 1 element to Class Title and Instructor
		parent := e.DOM.Parent()
		parent.Find("span").Each(func(i int, s *goquery.Selection) {

			_, classTitleExists := s.Attr("class")
			if anchorText := s.Find("a").Text(); anchorText != "" {
				if classTitleExists {
					class.SetTitle(anchorText)

					//fmt.Println("time:", s.Find("span.no-break").First().Text())
					timeStr := s.Find("span.no-break").First().Text()
					class.SetTime(timeStr)

					//fmt.Println("Class: " + anchorText)
				} else {
					class.SetInstructor(anchorText)
					//fmt.Println("Instructor: " + anchorText)
				}
			}
		})

		// Go up 2 element to Link, # of Spots and Class Status
		classPanel := parent.Parent()
		classPanel.Find(".col-sm-6.text-right").Each(func(i int, div *goquery.Selection) {
			if div != nil {
				smallNumSpotsLeft := div.Find("small.margin-right").First()
				//fmt.Println("Spots:", smallNumSpotsLeft.Text())
				numSpotsLeftStr := smallNumSpotsLeft.Text()
				class.SetNumSpotsLeftFromString(numSpotsLeftStr)

				bookingLink := div.Find("a").First()
				bookingLinkStr, exists := bookingLink.Attr("href")
				if exists {
					//fmt.Println("Booking Link:", bookingLinkStr)
					class.SetURL(bookingLinkStr)
				}

				//fmt.Println("Class Status:", bookingLink.Text())
				classStatus := bookingLink.Text()
				class.SetAvailableFromStatus(classStatus)
			}
		})

		classesListed = append(classesListed, *class)
	})

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Content-Type", "application/x-www-form-urlencoded")
		r.Headers.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko)")
		r.Headers.Set("Upgrade-Insecure-Requests", "1")
		r.Headers.Set("Origin", config.GymBaseURL)
		r.Headers.Set("Sec-Fetch-Dest", "document")
		r.Headers.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
		r.Headers.Set("Sec-Fetch-Site", "same-origin")
		r.Headers.Set("Sec-Fetch-Mode", "navigate")
		r.Headers.Set("Sec-Fetch-User", "?1")
		r.Headers.Set("Referer", config.GetClassesURL())
		r.Headers.Set("Accept-Encoding", "gzip, deflate")
		r.Headers.Set("Accept-Language", "en-US,en;q=0.9,zh-HK;q=0.8,zh;q=0.7")
	})

	c.OnResponse(func(r *colly.Response) {
		log.Println("PostGymClassesWithDate: Response received", r.StatusCode)
	})

	// Set Cookies if we have them
	if cookies != nil {
		c.SetCookies(config.GymBaseURL, cookies)
	}

	payload := getPayload(epochTimestamp, hiddenIds)
	//log.Println("Payload:", payload)
	err := c.Post(config.GetClassesURL(), payload)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return classesListed, nil
}

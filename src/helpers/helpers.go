package helpers

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/tkanos/gonfig"
	"gymclassbooker.com/bwoo/config"
)

func GetConfig(path string) *config.Config {

	var config config.Config
	gonfig.GetConf(path, &config)
	return &config
}

func GetEpochTimestamp(year, month, date, location string) (string, error) {

	loc, _ := time.LoadLocation(location)
	yearInt, _ := strconv.Atoi(year)
	monthInt, _ := strconv.Atoi(month)
	dateInt, _ := strconv.Atoi(date)

	t := time.Date(yearInt, time.Month(monthInt), dateInt, 0, 0, 0, 0, loc)
	secs := t.Unix()
	epochStr := strconv.FormatInt(secs, 10)
	return epochStr, nil
}

func GetRegexFromSearchStr(searchStr string) (*regexp.Regexp, error) {

	if len(searchStr) == 0 {
		return nil, fmt.Errorf("searchStr is empty")
	}

	regexpStr := fmt.Sprintf("(?i)%s", searchStr)
	r, err := regexp.Compile(regexpStr)
	return r, err

}

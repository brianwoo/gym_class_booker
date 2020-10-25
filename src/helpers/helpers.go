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

func GetEpochTimestamp(year, month, date string) (string, error) {

	str := year + "-" + month + "-" + date + "T06:00:00.000Z"
	t, err := time.Parse(time.RFC3339, str)
	if err != nil {
		fmt.Println("error getting epoch timestamp: ", err)
		return "", err
	}

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

// func FilterClasses(timeFilter, titleFilter string, listOfClasses []gymClasses.GymClass) []gymClasses.GymClass {

// 	// var r *regexp.Regexp
// 	// if len(titleFilter) > 0 {
// 	// 	regexpStr := fmt.Sprintf("(?i)%s", titleFilter)
// 	// 	r, _ = regexp.Compile(regexpStr)
// 	// }
// 	r, titleStrRegexErr := GetRegexFromSearchStr(titleFilter)

// 	filteredClasses := make([]gymClasses.GymClass, 0)
// 	for _, class := range listOfClasses {

// 		gymClassRef := &class

// 		if titleStrRegexErr == nil {
// 			if gymClassRef.IsMatchTitle(r) && gymClassRef.IsMatchTime(timeFilter) {
// 				filteredClasses = append(filteredClasses, class)
// 			}
// 		} else if gymClassRef.IsMatchTime(timeFilter) {
// 			filteredClasses = append(filteredClasses, class)
// 		}
// 	}

// 	return filteredClasses
// }

// func isMatchTitle(classTitle string, r *regexp.Regexp) bool {
// 	return r.Match([]byte(classTitle))
// }

// func isMatchTime(classTime string, timeFilterToMatch string) bool {
// 	if len(timeFilterToMatch) == 0 {
// 		return true
// 	}

// 	return classTime == timeFilterToMatch
// }

// func PrintClasses(year, month, date string, listOfClasses []gymClasses.GymClass) {

// 	fmt.Printf("%s-%s-%s: # of Classes: %v\n", year, month, date, len(listOfClasses))

// 	writer := tabwriter.NewWriter(os.Stdout, 0, 8, 1, '\t', tabwriter.AlignRight)
// 	fmt.Fprintf(writer, "%s %s\t%s\t%s\t%v\t%v\n", "#",
// 		"Class Title",
// 		"Time",
// 		"Instructor",
// 		"# Spots Left",
// 		"Booked?")

// 	for i, class := range listOfClasses {
// 		gymClassRef := &class
// 		gymClassRef.PrintClassInfo(writer, i+1)
// 	}

// 	writer.Flush()
// }

// func printClassInfo(writer *tabwriter.Writer, index int, class gymClasses.GymClass) {
// 	fmt.Fprintf(writer, "%v %s\t%s\t%s\t%v\t%v\n", index,
// 		class.GetTitle(),
// 		class.GetTime(),
// 		class.GetInstructor(),
// 		class.GetNumSpotsLeft(),
// 		class.GetIsBooked())
// }

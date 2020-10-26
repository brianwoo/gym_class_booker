package gymClasses

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"text/tabwriter"

	"gymclassbooker.com/bwoo/helpers"
)

type GymClass struct {
	title       string
	instructor  string
	time        string
	timestamp   string
	classId     string
	classStatus string
	isAvailable bool
	numOfSpots  int
	url         string
	isBooked    bool
}

func (gc *GymClass) SetTitle(title string) {
	gc.title = title
}

func (gc *GymClass) GetTitle() string {
	return gc.title
}

func (gc *GymClass) SetInstructor(instructor string) {
	gc.instructor = instructor
}

func (gc *GymClass) GetInstructor() string {
	return gc.instructor
}

func (gc *GymClass) SetTime(time string) {

	timeSlice := strings.Split(time, " ")
	gc.time = strings.Join(timeSlice[1:], " ")
}

func (gc *GymClass) GetTime() string {
	return gc.time
}

func (gc *GymClass) SetTimestamp(timestamp string) {
	gc.timestamp = timestamp
}

func (gc *GymClass) GetTimestamp() string {
	return gc.timestamp
}

func (gc *GymClass) SetAvailableFromStatus(status string) {

	if status == "Online Booking Closed" || status == "Too Late for Booking" || gc.numOfSpots == 0 {
		gc.isAvailable = false
	} else {
		gc.isAvailable = true
	}

	if status == "Book" {
		gc.isBooked = false
	} else if status == "Cancel" {
		gc.isBooked = true
	}

}

func (gc *GymClass) GetIsAvailable() bool {
	return gc.isAvailable
}

func (gc *GymClass) GetIsBooked() bool {
	return gc.isBooked
}

func (gc *GymClass) SetNumSpotsLeftFromString(numSpotsLeft string) {

	if len(numSpotsLeft) == 0 {
		return
	}

	spotsLeftSlice := strings.Split(numSpotsLeft, " ")
	numWithStar := []rune(spotsLeftSlice[0])
	gc.numOfSpots, _ = strconv.Atoi(string(numWithStar[1:]))
}

func (gc *GymClass) GetNumSpotsLeft() int {
	return gc.numOfSpots
}

func (gc *GymClass) SetURL(url string) {
	gc.url = url
}

func (gc *GymClass) GetURL() string {
	return gc.url
}

func (gc *GymClass) PrintClassInfo(writer *tabwriter.Writer, index int) {
	fmt.Fprintf(writer, "%v %s\t%s\t%s\t%v\t%v\n", index,
		gc.GetTitle(),
		gc.GetTime(),
		gc.GetInstructor(),
		gc.GetNumSpotsLeft(),
		gc.GetIsBooked())
}

func (gc *GymClass) IsMatchTitle(searchStrInRegex *regexp.Regexp) bool {
	return searchStrInRegex.Match([]byte(gc.GetTitle()))
}

func (gc *GymClass) IsMatchTime(timeFilterToMatch string) bool {
	if len(timeFilterToMatch) == 0 {
		return true
	}

	return gc.GetTime() == timeFilterToMatch
}

func (gc *GymClass) GetBookAClassMsg(year, month, date string) (subject string, msg string) {

	return "Gym Booking Successful", fmt.Sprintf("Booking:\n\n%s @%s on %s-%s-%s",
		gc.GetTitle(), gc.GetTime(), year, month, date)
}

type GymClasses []GymClass

func (gcs GymClasses) FilterClasses(timeFilter, titleFilter string) GymClasses {

	r, titleStrRegexErr := helpers.GetRegexFromSearchStr(titleFilter)

	filteredClasses := make(GymClasses, 0)
	for _, class := range gcs {

		gymClassRef := &class

		if titleStrRegexErr == nil {
			if gymClassRef.IsMatchTitle(r) && gymClassRef.IsMatchTime(timeFilter) {
				filteredClasses = append(filteredClasses, class)
			}
		} else if gymClassRef.IsMatchTime(timeFilter) {
			filteredClasses = append(filteredClasses, class)
		}
	}

	return filteredClasses
}

func (gcs GymClasses) PrintClasses(year, month, date string) {

	fmt.Printf("%s-%s-%s: # of Classes: %v\n", year, month, date, len(gcs))

	writer := tabwriter.NewWriter(os.Stdout, 0, 8, 1, '\t', tabwriter.AlignRight)
	fmt.Fprintf(writer, "%s %s\t%s\t%s\t%v\t%v\n", "#",
		"Class Title",
		"Time",
		"Instructor",
		"# Spots Left",
		"Booked?")

	for i, class := range gcs {
		gymClassRef := &class
		gymClassRef.PrintClassInfo(writer, i+1)
	}

	writer.Flush()
}

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/brianwoo/gogmail"
	"gymclassbooker.com/bwoo/config"
	"gymclassbooker.com/bwoo/gymClasses"
	"gymclassbooker.com/bwoo/helpers"
)

const programName = "gymClassBooker"

func getConfig() (*config.Config, error) {

	configFileLocation := os.Getenv("GYM_CLASS_BOOKER_CONFIG")
	if len(configFileLocation) == 0 {
		err := fmt.Errorf("Please set env var GYM_CLASS_BOOKER_CONFIG")
		return nil, err
	}

	config := helpers.GetConfig(configFileLocation)
	return config, nil
}

func printUsage() {

	fmt.Printf("Usage: %s [cmd]   [year] [month] [date] [time] [keyword in title]\n", programName)
	fmt.Printf("E.g. : %s list      2020 10 29 \"10:30 AM\" \"Yoga\"\n", programName)
	fmt.Printf("E.g. : %s list      2020 10 29 \"\" \"\"\n", programName)
	fmt.Printf("E.g. : %s book      2020 10 29 \"10:30 AM\" \"Yoga\"\n", programName)
	fmt.Printf("E.g. : %s bookIn7D  2020 10 29 \"10:30 AM\" \"Yoga\"\n", programName)
	fmt.Printf("E.g. : %s bookIn7D  \"\" \"\" \"\" \"10:30 AM\" \"Yoga\"\n", programName)
	fmt.Printf("E.g. : %s unbook    2020 10 29 \"10:30 AM\" \"Yoga\"\n", programName)
}

func getArgs() (string, string, string, string, string, string, error) {

	args := os.Args
	if len(args) < 7 {
		err := fmt.Errorf("Incorrect arguments")
		return "", "", "", "", "", "", err
	}

	return args[1], args[2], args[3], args[4], args[5], args[6], nil
}

func main() {

	config, err := getConfig()
	if err != nil {
		fmt.Println(err)
		return
	}

	cmd, year, month, date, iTime, title, err := getArgs()
	if err != nil {
		printUsage()
		return
	}

	if cmd == "list" {

		cookies, err := login(config)
		if err != nil {
			log.Println("Unable to login!")
			return
		}
		getClasses(year, month, date, iTime, title, config, cookies)
		if err != nil {
			fmt.Println(err)
		}

	} else if cmd == "book" {

		bookClassAndEmail(year, month, date, iTime, title, config)

	} else if cmd == "unbook" {

		cookies, err := login(config)
		if err != nil {
			log.Println("Unable to login!")
			return
		}
		err = unbookAClass(year, month, date, iTime, title, config, cookies)
		if err != nil {
			log.Println(err)
		}

	} else if cmd == "bookIn7D" {

		yearIn7D, monthIn7D, dateIn7D := getTodayOrCustomDate(year, month, date)
		bookClassAndEmail(yearIn7D, monthIn7D, dateIn7D, iTime, title, config)

	} else {
		fmt.Println("Command " + cmd + " not recognized")
	}
}

func getTodayOrCustomDate(year, month, date string) (string, string, string) {

	var inputDate time.Time
	if year == "" || month == "" || date == "" {
		inputDate = time.Now()
	} else {

		yearInt, err := strconv.Atoi(year)
		if err != nil {
			log.Fatalln("Invalid year argument!")
		}
		monthInt, err := strconv.Atoi(month)
		if err != nil {
			log.Fatalln("Invalid month argument!")
		}
		dateInt, err := strconv.Atoi(date)
		if err != nil {
			log.Fatalln("Invalid date argument!")
		}
		inputDate = time.Date(yearInt, time.Month(monthInt), dateInt, 0, 0, 0, 0, time.UTC)
	}

	addedDate := inputDate.AddDate(0, 0, 7)
	yearIn7D := strconv.Itoa(addedDate.Year())
	monthIn7D := fmt.Sprintf("%02d", addedDate.Month())
	dateIn7D := fmt.Sprintf("%02d", addedDate.Day())

	return yearIn7D, monthIn7D, dateIn7D
}

func bookClassAndEmail(year, month, date, time, title string, config *config.Config) {

	cookies, err := login(config)
	if err != nil {
		log.Println("Unable to login!")
		return
	}
	class, err := bookAClass(year, month, date, time, title, config, cookies)
	if class != nil {
		subject, msg := class.GetBookAClassMsg(year, month, date)
		if err != nil {
			msg += "\n" + err.Error()
			subject = err.Error()
		}
		sendEmailToUser(config.EmailToUser, subject, msg, config)
	} else if err != nil {
		sendEmailToUser(config.EmailToUser, "Unable to book class", err.Error(), config)
	}
}

func sendEmailToUser(to, subject, emailBody string, config *config.Config) error {

	sendEmail := gogmail.NewSendMail(config.EmailOAuthClientID,
		config.EmailOAuthClientSecret,
		config.EmailOAuthAccessToken,
		config.EmailOAuthRefreshToken)

	err := sendEmail.Send(to, subject, emailBody)
	if err != nil {
		log.Println("Unable to send email notification")
	} else {
		log.Println("Email sent!")
	}
	return err
}

func login(config *config.Config) ([]*http.Cookie, error) {
	return gymClasses.Login(config)
}

func unbookAClass(year, month, date, time, titleFilter string, config *config.Config, cookies []*http.Cookie) error {

	listOfClasses, err := getClasses(year, month, date, time, titleFilter, config, cookies)
	if err != nil {
		log.Println("Error occurred", err)
		return err
	}
	if len(listOfClasses) != 1 {
		err := fmt.Errorf("More than 1 class found, please refine the titleFilter")
		log.Println(err)
		return err
	}

	gymClass := listOfClasses[0]
	if gymClass.GetIsBooked() {
		err = gymClasses.UnbookAClass(gymClass, config, cookies)
	} else {
		err = fmt.Errorf("The class was not booked")
		return err
	}

	// verify
	listOfClasses, err = getClasses(year, month, date, time, titleFilter, config, cookies)
	gymClass = listOfClasses[0]
	if gymClass.GetIsBooked() {
		err = fmt.Errorf("Unable to unbook class")
	}

	return err
}

func bookAClass(year, month, date, time, titleFilter string, config *config.Config, cookies []*http.Cookie) (*gymClasses.GymClass, error) {

	listOfClasses, err := getClasses(year, month, date, time, titleFilter, config, cookies)
	if err != nil {
		log.Println("Error occurred", err)
		return nil, err
	}
	if len(listOfClasses) != 1 {
		err := fmt.Errorf("More than 1 class found, please refine the titleFilter")
		log.Println(err)
		return nil, err
	}

	gymClass := listOfClasses[0]
	if !gymClass.GetIsBooked() {
		err = gymClasses.BookAClass(gymClass, config, cookies)
	} else {
		err = fmt.Errorf("The class was already booked")
		return &gymClass, err
	}

	// verify
	listOfClasses, err = getClasses(year, month, date, time, titleFilter, config, cookies)
	gymClass = listOfClasses[0]
	if !gymClass.GetIsBooked() {
		err = fmt.Errorf("Unable to book class")
		return &gymClass, err
	}

	return &gymClass, nil
}

// timeFilter and titleFilter optional (i.e. use empty string "").
func getClasses(year, month, date, timeFilter, titleFilter string, config *config.Config, cookies []*http.Cookie) (gymClasses.GymClasses, error) {

	hiddenIds, err := gymClasses.GetInitFormIdsFromClasses(config, cookies)
	if err != nil {
		return nil, err
	}

	epochTimestamp, err := helpers.GetEpochTimestamp(year, month, date, config.Location)
	if err != nil {
		return nil, err
	}

	listOfClasses, err := gymClasses.PostClassesWithDate(epochTimestamp, hiddenIds, config, cookies)
	if err != nil {
		return nil, err
	}

	filteredClasses := listOfClasses.FilterClasses(timeFilter, titleFilter)
	filteredClasses.PrintClasses(year, month, date)

	return filteredClasses, nil
}

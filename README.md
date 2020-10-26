# gym_class_booker
Gym Class Booker is a commandline tool to list/book/unbook gym classes of a local gym

# Dependencies
go get -u github.com/brianwoo/gogmail

go get -u github.com/gocolly/colly/v2

go get -u github.com/tkanos/gonfig

# How to run
Usage: gymClassBooker [cmd] [year] [month] [date] [time] [keyword in title]

### Command (list) -- time and title are optional:
#### This command lists the classes available

E.g. : gymClassBooker list      2020 10 29 "" ""

E.g. : gymClassBooker list      2020 10 29 "10:30 AM" "Yoga"

### Command (book):
#### Note: If more than one class is found, it will not book the class as a precaution

E.g. : gymClassBooker book      2020 10 29 "10:30 AM" "Yoga"

### Command (unbook):
#### Note: If more than one class is found, it will not unbook the class as a precaution

E.g. : gymClassBooker unbook    2020 10 29 "10:30 AM" "Yoga"

### Command (bookIn7D):
#### This command works similarly as book but it books in a week from now.

#### Note: when date and time is empty, it will take the current time + 7 days

E.g. : gymClassBooker bookIn7D  "" "" "" "10:30 AM" "Yoga"

E.g. : gymClassBooker bookIn7D  2020 10 29 "10:30 AM" "Yoga"

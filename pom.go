package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"strings"
	"time"
)

const EPOCH = 90
const EPSILONg = 279.403303 /* solar ecliptic long at EPOCH */
const RHOg = 282.768422     /* solar ecliptic long of perigee at EPOCH */
const ECCEN = 0.016713      /* solar orbit eccentricity */
const lzero = 318.351648    /* lunar mean long at EPOCH */
const Pzero = 36.340410     /* lunar mean long of perigee at EPOCH */
const Nzero = 318.510107    /* lunar mean long of node at EPOCH */
const someSpecialTimeStamp = "1989123100"

type potm struct {
	date       time.Time
	state      string
	percentage float64
}

// ensure 0 <= deg <= 365
func adj360(deg float64) {
	deg = math.Mod(deg, 360.0)
	if deg < 0.0 {
		deg += 360.0
	}
}

// degrees to radians
func dtor(deg float64) float64 {
	return (deg * math.Pi / 180)
}

/*
 * Phase of the Moon.  Calculates the current phase of the moon.
 * Based on routines from `Practical Astronomy with Your Calculator',
 * by Duffett-Smith.  Comments give the section from the book that
 * particular piece of code was adapted from.
 *
 * -- Keith E. Brandt  VIII 1984
 *
 * Updated to the Third Edition of Duffett-Smith's book, IX 1998
 *
 * Taken from OpenBSD pom.c
 */
func calculatePercentage(days float64) float64 {
	var N, Msol, Ec, LambdaSol, l, Mm, Ev, Ac, A3, Mmprime float64
	var A4, lprime, V, ldprime, D, Nm float64

	N = 360.0 * days / 365.242191 /* sec 46 #3 */
	adj360(N)
	Msol = N + EPSILONg - RHOg /* sec 46 #4 */
	adj360(Msol)
	Ec = 360 / math.Pi * ECCEN * math.Sin(dtor(Msol)) /* sec 46 #5 */
	LambdaSol = N + Ec + EPSILONg                     /* sec 46 #6 */
	adj360(LambdaSol)
	l = 13.1763966*days + lzero /* sec 65 #4 */
	adj360(l)
	Mm = l - (0.1114041 * days) - Pzero /* sec 65 #5 */
	adj360(Mm)
	Nm = Nzero - (0.0529539 * days) /* sec 65 #6 */
	adj360(Nm)
	Ev = 1.2739 * math.Sin(dtor(2*(l-LambdaSol)-Mm)) /* sec 65 #7 */
	Ac = 0.1858 * math.Sin(dtor(Msol))               /* sec 65 #8 */
	A3 = 0.37 * math.Sin(dtor(Msol))
	Mmprime = Mm + Ev - Ac - A3                       /* sec 65 #9 */
	Ec = 6.2886 * math.Sin(dtor(Mmprime))             /* sec 65 #10 */
	A4 = 0.214 * math.Sin(dtor(2*Mmprime))            /* sec 65 #11 */
	lprime = l + Ev + Ec - Ac + A4                    /* sec 65 #12 */
	V = 0.6583 * math.Sin(dtor(2*(lprime-LambdaSol))) /* sec 65 #13 */
	ldprime = lprime + V                              /* sec 65 #14 */
	D = ldprime - LambdaSol                           /* sec 67 #2 */
	return (50.0 * (1 - math.Cos(dtor(D))))           /* sec 67 #3 */
}

func moonPercentage(date time.Time) float64 {
	someSpecialDate, err := time.Parse("2006010215", someSpecialTimeStamp)
	if err != nil {
		log.Panic(err)
	}
	delta := date.Sub(someSpecialDate)
	days := delta.Hours() / 24
	return calculatePercentage(days)
}

func particularMoonPhase(date time.Time) {
	percentageNow := moonPercentage(date)
	state := ""
	// 0 == New Moon
	// 0 < Waxing Crescent < 50
	// 50 == First Quarter
	// 50 < Waxing Gibbous < 100
	// 100 == Full moon
	// 100 > Waning Gibbous > 50
	// 50 == Last Quarter
	// 50 > Waning Crescent > 0
	// 0 == New Moon
	if math.Round(percentageNow) == 100 {
		state = "Full"
	} else if math.Round(percentageNow) == 0 {
		state = "New"
	} else {
		tomorrow := date.Add(24 * time.Hour)
		percentageTomorrow := moonPercentage(tomorrow)
		if math.Round(percentageNow) == 50 {
			if percentageTomorrow > percentageNow {
				state = "First quarter"
			} else {
				state = "Last quarter"
			}
		} else { // percentage != 50
			if percentageTomorrow > percentageNow {
				state = "Waxing"
			} else {
				state = "Waning"
			}
			if percentageNow < 50 {
				state = state + " Crescent"
			} else {
				state = state + " Gibbous"
			}
		}
	}
	potmNow := potm{state: state, percentage: percentageNow, date: date}
	f := false
	potmNow.print(&f)
}

func (p potm) print(pound *bool) {
	if *pound {
		fmt.Printf("%v: %v\n", p.date.Format(time.RFC3339),
			strings.Repeat("#", int(math.Round(p.percentage))))
	} else {
		if p.state != "" {
			fmt.Printf("%v: %s moon at %.1f%%\n",
				p.date.Format(time.RFC3339), p.state, p.percentage)
		} else {
			fmt.Printf("%v: moon at %.1f%%\n",
				p.date.Format(time.RFC3339), p.percentage)
		}
	}
}

func nextMoonPhases(days int, measurementsPerDay int, pound *bool, beg time.Time) {
	measures := make([]potm, measurementsPerDay*days)
	for i := 0; i < (measurementsPerDay * days); i += 1 {
		frequency := 24 / measurementsPerDay
		date := beg.Add(time.Duration(frequency*i) * time.Hour)
		measures[i] = potm{date: date, percentage: moonPercentage(date)}
	}
	for _, p := range measures {
		p.print(pound)
	}
}

func whenNextMoonState(state string) {
	var percentage float64
	switch state {
	case "full":
		percentage = 100
	case "new":
		percentage = 0
	}
	date := time.Now()
	for {
		potm := moonPercentage(date)
		if percentage == math.Round(potm) {
			// Mon Jan 2 15:04:05 -0700 MST 2006
			fmt.Printf("%v:00\n", date.Format("2006-01-02T15 MST"))
			break
		}
		date = date.Add(1 * time.Hour)
	}
}

func statusAtDate(date string) {
	var parsedTime time.Time
	var err error
	if len(date) == 7 {
		parsedTime, err = time.Parse("2006-01", date)
		if err != nil {
			log.Fatal(err)
		}
	} else if len(date) == 10 {
		parsedTime, err = time.Parse("2006-01-02", date)
		if err != nil {
			log.Fatal(err)
		}
	} else if len(date) == 13 {
		parsedTime, err = time.Parse("2006-01-02:15", date)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.Fatal("Date format is wrong, should be YYYY-MM[-DD[:HH]]")
	}

	particularMoonPhase(parsedTime)
}

func main() {
	now := flag.Bool("now", false, "Moon status right now")
	weeklyMode := flag.Bool("week", false,
		"Moon status for the next seven days")
	monthlyMode := flag.Bool("month", false,
		"Moon status for the next 28 days")
	pound := flag.Bool("pound", false,
		"Print moon status with pounds instead, works for !now")
	nextFullMoon := flag.Bool("full", false, "Give next full moon date")
	nextNewMoon := flag.Bool("new", false, "Give next new moon date")
	date := flag.String("date", "",
		"Returns moon status at given date (YYYY-MM[-DD[:HH]] format")
	flag.Parse()

	modeSelected := 0
	if *now {
		modeSelected += 1
	}
	if *weeklyMode {
		modeSelected += 1
	}
	if *monthlyMode {
		modeSelected += 1
	}
	if *nextFullMoon {
		modeSelected += 1
	}
	if *nextNewMoon {
		modeSelected += 1
	}
	if *date != "" {
		modeSelected += 1
	}
	if modeSelected > 1 {
		log.Print("Error: multiple modes selected. Pick only one")
		flag.PrintDefaults()
		os.Exit(1)
	}
	if modeSelected == 0 {
		log.Print("The mode is missing. Pick one")
		flag.PrintDefaults()
		os.Exit(1)
	}

	if *now {
		now := time.Now()
		particularMoonPhase(now)
	} else if *weeklyMode {
		now := time.Now()
		nextMoonPhases(7, 4, pound, now)
	} else if *monthlyMode {
		now := time.Now()
		nextMoonPhases(28, 2, pound, now)
	} else if *nextFullMoon {
		whenNextMoonState("full")
	} else if *nextNewMoon {
		whenNextMoonState("new")
	} else if *date != "" {
		statusAtDate(*date)
	}
}

package main

import (
	"encoding/csv"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/jaswdr/faker/v2"
)

const CREDIT = "CREDIT"
const DEBIT = "DEBIT"

type Transaction struct {
	Details           string
	PostingDate       string
	Description       string
	Amount            string
	Type              string
	Balance           string
	CheckOrSlipNumber string
}

func isPayDay(date time.Time) bool {
	fifteenth := 15
	lastDayOfMonth := time.Date(date.Year(), date.Month()+1, 0, 0, 0, 0, 0, date.Location()).Day()
	isWeekend := date.Weekday() == time.Saturday || date.Weekday() == time.Sunday
	isFifteenthOrLastDayOfTheMonth := func(d int) bool { return d == fifteenth || d == lastDayOfMonth }
	isFriday := date.Weekday() == time.Friday

	isWeekendFifteenthOrLastDayOfTheMonth := func() bool {
		if !isFriday {
			return false
		}
		if saturday := date.AddDate(0, 0, 1); isFifteenthOrLastDayOfTheMonth(saturday.Day()) {
			return true
		}
		if sunday := date.AddDate(0, 0, 2); isFifteenthOrLastDayOfTheMonth(sunday.Day()) {
			return true
		}
		return false
	}

	switch {
	case isFriday && isWeekendFifteenthOrLastDayOfTheMonth():
		return true
	case !isWeekend && isFifteenthOrLastDayOfTheMonth(date.Day()):
		return true
	default:
		return false
	}
}

func transactionToRow(t Transaction) []string {
	return []string{
		t.Details,
		t.PostingDate,
		t.Description,
		t.Amount,
		t.Type,
		t.Balance,
		t.CheckOrSlipNumber,
	}
}

func writeHeader(writer *csv.Writer) {
	header := []string{"Details", "Posting Date", "Description", "Amount", "Type", "Balance", "Check or Slip #"}
	if err := writer.Write(header); err != nil {
		log.Fatalf("Failed to write header: %s", err)
	}
}

func writeRow(writer *csv.Writer, row []string) {
	if err := writer.Write(row); err != nil {
		log.Fatalf("Failed to write row: %s", err)
	}
}

func main() {
	time_seed := time.Now().UnixNano()

	println("time_seed: ", time_seed)
	seed := int64(1)
	fake := faker.NewWithSeed(rand.NewSource(seed))

	println(fake.Person().Name())

	file, err := os.Create("testdata/Chase Activity Feb 1.CSV")
	if err != nil {
		log.Fatalf("Failed to create file: %s", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writeHeader(writer)
	startDate := time.Date(2024, time.February, 1, 0, 0, 0, 0, time.UTC)
	println("startDate: ", startDate.Format("Monday 01/02/2006"))

	switch startDate.Weekday() {
	case time.Thursday:
		println("It's Thursday!")
	default:
		println("It's not Thursday.")
	}

	endDate := time.Date(2023, time.December, 31, 0, 0, 0, 0, time.UTC)

	for date := startDate; date.After(endDate); date = date.AddDate(0, 0, -1) {
		if isPayDay(date) {
			formattedDateString := date.Format("01/02/2006")
			writeRow(writer, transactionToRow(Transaction{
				Details: CREDIT,
				PostingDate: formattedDateString,
				Description: "PAYROLL",
				Amount: "3500.00",
				Type: "ACH_CREDIT",
			}))
			continue
		}
		if date.Month() == time.January {
			formattedDateString := date.Format("01/02/2006")
			writeRow(writer, transactionToRow(Transaction{PostingDate: formattedDateString}))
		}
	}
}

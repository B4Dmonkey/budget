package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"math"
	"math/rand/v2"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jaswdr/faker/v2"
)

const CREDIT = "CREDIT"
const DEBIT = "DEBIT"

const (
	Details = iota
	PostingDate
	Description
	Amount
	Type
	Balance
	CheckOrSlipNumber
)

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
	isFriday := date.Weekday() == time.Friday

	isFifteenthOrLastDayOfTheMonth := func(d int) bool { return d == fifteenth || d == lastDayOfMonth }

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

	if isFriday && isWeekendFifteenthOrLastDayOfTheMonth() {
		return true
	}

	if !isWeekend && isFifteenthOrLastDayOfTheMonth(date.Day()) {
		return true
	}

	return false
}

func isSpotifyBill(date time.Time) bool {
	lastDayOfMonth := time.Date(date.Year(), date.Month()+1, 0, 0, 0, 0, 0, date.Location()).Day()
	if lastDayOfMonth >= 29 && date.Day() == 29 {
		return true
	}

	if lastDayOfMonth == 28 && date.Day() == 28 {
		return true
	}

	return false
}

func isRentDue(date time.Time) bool {
	// Todo: Make this fuzzy so that it falls on a few days after the 15
	if date.Day() == 15 {
		return true
	}

	return false
}

func newTransactionCSV(file *os.File) *csv.Writer {
	writer := csv.NewWriter(file)
	header := []string{"Details", "Posting Date", "Description", "Amount", "Type", "Balance", "Check or Slip #"}

	if err := writer.Write(header); err != nil {
		log.Fatalf("Failed to write header: %s", err)
	}

	return writer
}

func (t Transaction) transactionToRow() []string {
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

func generateTransactions(writer *csv.Writer, balance float64, start_date time.Time, end_date time.Time) {
	var transactions []Transaction

	for date := start_date; !date.After(end_date); date = date.AddDate(0, 0, 1) {
		formattedDateString := date.Format("01/02/2006")

		if isPayDay(date) {
			amount := 3500.00
			balance += amount
			transactions = append(transactions, Transaction{
				Balance:     fmt.Sprintf("%.2f", balance),
				Details:     CREDIT,
				PostingDate: formattedDateString,
				Description: "PAYROLL",
				Amount:      fmt.Sprintf("%.2f", amount),
				Type:        "ACH_CREDIT",
			})
		}

		if isRentDue(date) {
			amount := -1712.50
			balance += amount
			transactions = append(transactions, Transaction{
				Balance:     fmt.Sprintf("%.2f", balance),
				Details:     DEBIT,
				PostingDate: formattedDateString,
				Description: "Rent",
				Amount:      fmt.Sprintf("%.2f", amount),
				Type:        "DEBIT_CARD",
			})
		}

		if isSpotifyBill(date) {
			amount := -11.99
			balance += amount
			transactions = append(transactions, Transaction{
				Balance:     fmt.Sprintf("%.2f", balance),
				Details:     DEBIT,
				PostingDate: formattedDateString,
				Description: "Spotify Bill Proxy",
				Amount:      fmt.Sprintf("%.2f", amount),
				Type:        "DEBIT_CARD",
			})
		}
	}

	for i := len(transactions) - 1; i >= 0; i-- {
		if err := writer.Write(transactions[i].transactionToRow()); err != nil {
			println("Failed to write row: ", err)
			continue
		}
	}
}

func foo() {
	file, err := os.Open("Chase Activity Oct 6.CSV")
	if err != nil {
		log.Fatalf("Failed to open file: %s", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	if _, err := reader.Read(); err != nil {
		log.Fatalf("error skipping headers: %v", err)
	}

	var previousDate time.Time

	tx_count := 0
	gap_count := 0

	for {
		record, err := reader.Read()
		if err != nil && err == io.EOF {
			break
		}

		if parseErr, ok := err.(*csv.ParseError); ok && parseErr.Err != csv.ErrFieldCount {
			log.Fatalf("error parsing record: %v", err)
		}

		posting_date, err := time.Parse("01/02/2006", record[PostingDate])
		if err != nil {
			log.Fatalf("error parsing date: %v", err)
		}

		tx_count++

		if !previousDate.IsZero() && posting_date.Before(previousDate.AddDate(0, 0, -1)) {
			for d := previousDate.AddDate(0, 0, -1); d.After(posting_date); d = d.AddDate(0, 0, -1) {
				gap_count++
			}
		}

		previousDate = posting_date

		if strings.Contains(strings.ToLower(record[Description]), "mta") {
			continue
		}

		if strings.Contains(strings.ToLower(record[Description]), "venmo") {
			amount, err := strconv.ParseFloat(record[Amount], 64)
			if err != nil {
				log.Fatalf("error parsing amount: %v", err)
			}

			if amount == -1712 || amount == -1712.5 {
				continue
			}
		}

		if strings.Contains(strings.ToLower(record[Description]), "spotify") {
			continue
		}
	}

	println("Transaction Count: ", tx_count)
	println("Gap Count: ", gap_count)
}

func main() {
	SEED := rand.NewPCG(1, 2)
	r := rand.New(SEED)                                 // #nosec G404 // * This is not for security purposes
	start_balance := 10 + r.Float64()*(900-10)          // * A random number between 10 and 900
	start_balance = math.Round(start_balance*100) / 100 // * start balance should be 2 decimal precision
	println(fmt.Sprintf("Starting Balance: %.2f", start_balance))
	start_balance += 3500
	println(fmt.Sprintf("Starting Balance after Payroll: %.2f", start_balance))

	// ! Example of using faker. I may not need it
	fake := faker.New()
	fake.Person().Name()
	// ! End of example

	// ? Dynamically generate ?
	file, err := os.Create("/Users/appstack/Developer/Personal/budget/testdata/Chase Activity Feb 1.CSV")
	if err != nil {
		log.Fatalf("Failed to create file: %s", err)
	}
	defer file.Close()

	writer := newTransactionCSV(file)
	defer writer.Flush()

	foo()
	start_date := time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC)
	end_date := time.Date(start_date.Year(), start_date.Month()+1, 0, 0, 0, 0, 0, time.UTC)

	generateTransactions(writer, start_balance, start_date, end_date)
}

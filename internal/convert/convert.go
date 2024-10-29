package convert

import (
	"database/sql"
	"log"
	"strconv"
	"strings"
)

func CurrencyIntToString(amount int64) string {
	return strconv.FormatFloat(float64(amount)/100, 'f', 2, 64)
}

func CurrencyStringToInt64(amount string) (int64, error) {
	return strconv.ParseInt(strings.Replace(strings.TrimSpace(amount), ".", "", -1), 10, 64)
}

func CurrencyStringToNullInt64(amount string) sql.NullInt64 {
	amount = strings.TrimSpace(amount)
	if amount == "" {
		return sql.NullInt64{Valid: false}
	}

	int64_amount, err := CurrencyStringToInt64(amount)

	if err != nil {
		// Todo: I should check what errors this might be and handle them
		log.Println("Error converting currency string to int:", err)
	}

	return sql.NullInt64{
		Int64: int64_amount,
		Valid: err == nil,
	}
}

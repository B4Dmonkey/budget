package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/cbroglie/mustache"
	"github.com/stretchr/testify/assert"
)

type Data struct {
	PostingDate string
	Amount      float64
	Description string
	Balance     float64
}

func (d DataSlice) Template() string {
	return `
	{{#UnprocessedTransactions}}
		<tr>
			<td>{{PostingDate}}</td>
			<td>{{Amount}}</td>
			<td>{{Description}}</td>
			<td>{{Balance}}</td>
		</tr>
	{{/UnprocessedTransactions}}
	`
}

func (d DataSlice) Render(w http.ResponseWriter) {
	template := d.Template()
	parsedTemplate, _ := mustache.ParseString(template)
	content, err := parsedTemplate.Render(template, map[string]DataSlice{"UnprocessedTransactions": d})
	if err != nil {
		http.Error(w, "Internal server error two", http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, content)
}

type DataSlice []Data

func TestRender(t *testing.T) {
	w := httptest.NewRecorder()
	postingDate1 := "2024-09-01"
	amount1 := 100.00
	description1 := "Deposit"
	balance1 := 100.00
	postingDate2 := "2024-09-02"
	amount2 := 50.00
	description2 := "Withdrawal"
	balance2 := 50.00

	example_data := DataSlice{
		{
			PostingDate: postingDate1,
			Amount:      amount1,
			Description: description1,
			Balance:     balance1,
		},
		{
			PostingDate: postingDate2,
			Amount:      amount2,
			Description: description2,
			Balance:     balance2,
		},
	}
	example_data.Render(w)
	// Render(w, "", data )

	results := strings.TrimSpace(w.Body.String())
	assert.Contains(t, results, postingDate1, "Rendered output does not match expected output")
	// assert.Contains(t, results, amount1, "Rendered output does not match expected output")
	// assert.Contains(t, results, description1, "Rendered output does not match expected output")
	// assert.Contains(t, results, balance1, "Rendered output does not match expected output")
	// assert.Contains(t, results, postingDate2, "Rendered output does not match expected output")
	// assert.Contains(t, results, amount2, "Rendered output does not match expected output")
	// assert.Contains(t, results, description2, "Rendered output does not match expected output")
	// assert.Contains(t, results, balance2, "Rendered output does not match expected output")
}

type MockDatum struct {
	name string
}

func (d MockDatum) Template() string {
	return `
	{{#repo}}
  <b>{{name}}</b>
	{{/repo}}
	`
}

func TestMustache(t *testing.T) {
	data, err := mustache.Render("hello {{c}}", map[string]string{"c": "world"})
	assert.Nil(t, err)
	assert.Equal(t, "hello world", data)
	mock_data := new(MockDatum)

	template := mock_data.Template()

	expected := `
	<b>resque</b>
  <b>hub</b>
  <b>rip</b>
	`
	result, err := mustache.Render(template, map[string][]map[string]string{
		"repo": {
			{"name": "resque"},
			{"name": "hub"},
			{"name": "rip"},
		},
	})
	assert.Nil(t, err)
	assert.Equal(t, strings.TrimSpace(result), strings.TrimSpace(expected))

}

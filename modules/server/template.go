package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"math"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/dannyvankooten/extemplate"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"
	"github.com/highercomve/papelito/modules/helpers"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

const (
	minute  = 1
	hour    = minute * 60
	day     = hour * 24
	month   = day * 30
	year    = day * 365
	quarter = year / 4
)

// GeneralPayload structure to extend any payload
type GeneralPayload struct {
	Payload       interface{}
	Authenticated string
	Nick          string
	Template      string
	TemplateID    string
	ServerURL     string
}

var functions template.FuncMap = template.FuncMap{
	"notNil":            notNil,
	"hasField":          hasField,
	"divide":            divideNumbers,
	"multiply":          multiplyNumbers,
	"add":               addNumbers,
	"sub":               substractNumbers,
	"timeAgo":           timeAgo,
	"arrayToTextarea":   arrayToTextarea,
	"arrayToTextareaOr": arrayToTextareaOr,
	"printMD":           printMD,
	"noescape":          noescape,
	"percentage":        percentage,
}

// TemplateRenderer is a custom html/template renderer for Echo framework
type TemplateRenderer struct {
	templates *extemplate.Extemplate
}

// Render renders a template document
func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	contentType := c.Request().Header.Get("Accept")
	format := helpers.GetFormat(contentType)

	switch format {
	case "json":
		enc := json.NewEncoder(w)
		c.Response().Writer.Header().Set("Content-Type", echo.MIMEApplicationJSON)
		return enc.Encode(data)
	default:
		values := strings.Split(name, ":")
		newData := extendPayload(data, values[0], c)
		return t.templates.ExecuteTemplate(w, name, newData)
	}
}

// CreateTemplateRenderer return a renderer with all the templates views
func CreateTemplateRenderer() *TemplateRenderer {
	xt := extemplate.New().Funcs(functions)
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	err = xt.ParseDir(dir+"/templates/", []string{".html"})
	if err != nil {
		log.Fatal(err)
	}

	return &TemplateRenderer{
		templates: xt,
	}
}

func arrayToTextarea(arr []string) string {
	result := ""
	for _, element := range arr {
		if result == "" {
			result = element
		} else {
			result = fmt.Sprintf("%s\r\n%s", result, element)
		}
	}

	return result
}

func arrayToTextareaOr(arr []string, value string) string {
	if len(arr) == 0 {
		return value
	}
	result := ""
	for _, element := range arr {
		if result == "" {
			result = element
		} else {
			result = fmt.Sprintf("%s\r\n%s", result, element)
		}
	}

	return result
}

func extendPayload(data interface{}, name string, c echo.Context) interface{} {
	sess, _ := session.Get("session", c)
	token, tokenExist := sess.Values["token"]
	nick, _ := sess.Values["nick"].(string)
	authenticated := fmt.Sprintf("%v", tokenExist && token != nil && token != "")
	templateID := strings.ReplaceAll(name, "/", "-")
	templateID = strings.ReplaceAll(templateID, ".html", "")

	d, ok := data.(map[string]interface{})
	if ok {
		return map[string]interface{}{
			"Payload":       d,
			"Authenticated": authenticated,
			"Nick":          nick,
			"Template":      name,
			"TemplateID":    templateID,
			"ServerURL":     helpers.Env.HostURL,
		}
	}

	return &GeneralPayload{
		Payload:       data,
		Authenticated: authenticated,
		Nick:          nick,
		Template:      name,
		TemplateID:    templateID,
		ServerURL:     helpers.Env.HostURL,
	}
}

func ConvertToMap(data interface{}) map[string]interface{} {
	newMap := map[string]interface{}{}
	if data == nil {
		return newMap
	}

	d, err := json.Marshal(data) // Convert to a json string
	if err != nil {
		return newMap
	}

	err = json.Unmarshal(d, &newMap) // Convert to a map
	if err != nil {
		return newMap
	}
	return newMap
}

func hasField(v interface{}, name string) bool {
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		return false
	}
	return rv.FieldByName(name).IsValid()
}

func divideNumbers(a, b int) int {
	divition := a / b
	rest := a % b

	if rest > 0 {
		return divition + 1
	}

	return divition
}

func multiplyNumbers(a, b int) int {
	return a * b
}

func addNumbers(a, b int) int {
	return a + b
}

func substractNumbers(a, b int) int {
	return a - b
}

func notNil(a interface{}) bool {
	return !reflect.ValueOf(a).IsNil()
}

func timeAgo(when time.Time) string {
	d := time.Since(when)
	return fmt.Sprintf("%s ago", FromDuration(d))
}

func noescape(str string) template.HTML {
	return template.HTML(str)
}

func percentage(val int, total int) string {
	var per float32 = (float32(val) * 100) / float32(total)

	return fmt.Sprintf("%.2f", per)
}

// FromDuration returns a friendly string representing an approximation of the
// given duration
func FromDuration(d time.Duration) string {
	seconds := round(d.Seconds())

	if seconds < 30 {
		return "a few seconds"
	}

	if seconds < 90 {
		return "1 minute"
	}

	minutes := div(seconds, 60)

	if minutes < 45 {
		return fmt.Sprintf("%0d minutes", minutes)
	}

	hours := div(minutes, 60)

	if minutes < day {
		return pluralize(hours, "hour")
	}

	if minutes < (42 * hour) {
		return "1 day"
	}

	days := div(hours, 24)

	if minutes < (30 * day) {
		return pluralize(days, "day")
	}

	months := div(days, 30)

	if minutes < (45 * day) {
		return "1 month"
	}

	if minutes < (60 * day) {
		return "2 months"
	}

	if minutes < year {
		return pluralize(months, "month")
	}

	rem := minutes % year
	years := minutes / year

	if rem < (3 * month) {
		return pluralize(years, "year")
	}
	if rem < (9 * month) {
		return fmt.Sprintf("over %s", pluralize(years, "year"))
	}

	years++
	return fmt.Sprintf("almost %s", pluralize(years, "year"))
}

func pluralize(i int, s string) string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("%d %s", i, s))
	if i != 1 {
		buf.WriteString("s")
	}
	return buf.String()
}

func round(f float64) int {
	return int(math.Floor(f + .50))
}

func div(numerator int, denominator int) int {
	rem := numerator % denominator
	result := numerator / denominator

	if rem >= (denominator / 2) {
		result++
	}

	return result
}

func printMD(content string) string {
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs
	parser := parser.NewWithExtensions(extensions)

	html := markdown.ToHTML([]byte(content), parser, nil)

	return string(html)
}

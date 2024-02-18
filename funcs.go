package gosqlfmt2

import (
	"os"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		TimestampFormat: "15:04:05",
	})

	log.SetOutput(os.Stdout)

	// log.SetLevel(log.WarnLevel)
	log.SetLevel(log.DebugLevel)
	log.SetReportCaller(true)
}

func CleanQuery(query string) string {
	query = spacecomma.ReplaceAllString(query, ",")
	query = commaword.ReplaceAllString(query, `, $1`)

	// clean keywords e.g. with data as( =? with data as (
	for key := range KEYWORDS {
		pat1 := regexp.MustCompile(`(?i)(\b` + KEYWORDS[key] + `\b)([a-zA-Z0-9\(\)]+)`)
		pat2 := regexp.MustCompile(`(?i)([a-zA-Z0-9\(\)]+)(\b` + KEYWORDS[key] + `\b)`)

		if pat1.MatchString(query) {
			query = pat1.ReplaceAllString(query, "$1 $2")
		}

		if pat2.MatchString(query) {
			query = pat2.ReplaceAllString(query, "$1 $2")
		}
	}

	query = newlineword.ReplaceAllString(query, " $1")

	// remove extra white spaces
	query = whitespace.ReplaceAllString(query, space)

	// remove leading/trailing whitespace
	query = strings.TrimSpace(query)

	log.Infof("Cleaned Query: %v", query)

	return query
}

func formatquery(query string) string {
	// log.Debugf("Query received: %v\n", query)
	query = CleanQuery(query)
	var fmtdquery string
	// read each word
	words := strings.Split(query, space)

	j := 0
	for j < len(words) {
		word := words[j]
		log.Debugf("Word: %v\n", word)
		var (
			nextword string
			phrase   string
		)

		if j < len(words)-1 {
			nextword = words[j+1]
			phrase = word + space + nextword
		}

		if word == "(" && nextword != ")" {
			tab += "  "
			fmtdquery += space + word + newline + tab
		} else if word == ")" {
			tab = strings.Replace(tab, "  ", "", 1)
			fmtdquery += newline + tab + word
		} else if strings.HasSuffix(word, ",") {
			fmtdquery += word + newline + tab
		} else if selectk.MatchString(word) {
			// SELECT
			tab += "  "
			fmtdquery += word + newline + tab
		} else if from.MatchString(word) {
			// FROM
			tab = strings.Replace(tab, "  ", "", 1)
			if !strings.HasSuffix(fmtdquery, newline) {
				fmtdquery += newline + tab
			}
			fmtdquery += word + newline
			tab += "  "
			fmtdquery += tab
		} else if groupby.MatchString(phrase) {
			// FROM
			tab = strings.Replace(tab, "  ", "", 1)
			if !strings.HasSuffix(fmtdquery, newline) {
				fmtdquery += newline + tab
			}
			fmtdquery += phrase + newline
			tab += "  "
			fmtdquery += tab
			j += 2
			continue
		} else if word == "and" {
			fmtdquery += newline + tab + word
		} else {
			if strings.HasSuffix(fmtdquery, newline) {
				fmtdquery += tab + word
			} else if strings.HasSuffix(fmtdquery, tab) {
				fmtdquery += word
			} else {
				fmtdquery += space + word
			}
		}
		j++
	}
	return fmtdquery
	// return query
}

func Parse(fileName string) string {
	var fmtdquery string
	var query string
	file := GetQuery(fileName)

	// read each line
	lines := strings.Split(file, "\n")

	// TODO: DELETE LOG CALLS
	// for i := range lines {
	// 	log.Debugln(strings.TrimSpace(lines[i]))
	// }
	// log.Fatal()

	for i := range lines {
		line := strings.TrimSpace(lines[i])

		if comment.MatchString(line) {
			if len(query) > 0 {
				fmtdquery += formatquery(query) //formatquery is responsible for adding 2 newlines at the end of fmtdquery that it returns
				query = ""
			}
			fmtdquery += tab + line + newline + tab
			continue
		} else if len(line) == 0 {
			// fmtdquery += newline + tab //unsure about this
			continue
		} else {
			query += line + newline
		}
	}

	if len(query) > 0 {
		fmtdquery += formatquery(query)
	}

	return fmtdquery
}

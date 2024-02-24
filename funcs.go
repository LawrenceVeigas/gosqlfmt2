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
	// log.SetReportCaller(true)
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

func fmtlines(query string) string {
	var fmtquery string

	query = CleanQuery(query)

	words := strings.Split(query, space)
	length := len(words)

	i := 0
	for i < length {
		var (
			nextword string
			phrase   string
		)

		word := strings.TrimSpace(words[i])
		if i < length-1 {
			nextword = strings.TrimSpace(words[i+1])
			phrase = word + space + nextword
			log.Debugf("Phrase :%v", phrase)
		}

		if strings.HasSuffix(word, ",") {
			fmtquery += word + newline
		} else if selectk.MatchString(word) {
			fmtquery += newline + word + newline
		} else if from.MatchString(word) {
			fmtquery += newline + word + newline
		} else if word == "case" {
			fmtquery += newline + word + newline
		} else if prepositions.MatchString(word) {
			log.Debugf("Preposition: %v", word)
			fmtquery += newline + word + space
		} else if joins.MatchString(phrase) {
			// exceptional 3 word case
			if phrase == "full outer" {
				fmtquery += newline + phrase + space + "join" + space
				i += 3
			} else {
				fmtquery += newline + phrase + space
				i += 2
			}
			continue
		} else if i == length-1 {
			fmtquery += word + newline
		} else {
			fmtquery += word + space
		}
		i++
	}

	fmtquery = strings.ReplaceAll(fmtquery, "\n\n", "\n")

	return fmtquery
}

func formatquery(query string) string {
	query = fmtlines(query)
	return query
}

func Parse(fileName string) string {
	var fmtdquery string
	var query string
	file := GetQuery(fileName)

	// handle inline comments
	file = inlinecomment.ReplaceAllString(file, "${2}\n${1}")

	// read each line
	lines := strings.Split(file, "\n")

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

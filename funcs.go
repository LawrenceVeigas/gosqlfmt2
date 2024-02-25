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

	// bracket word
	query = regexp.MustCompile(`(\))(\w)`).ReplaceAllString(query, "$1 $2")

	// multiple brackets
	query = strings.ReplaceAll(query, "((", " ( (")
	query = strings.ReplaceAll(query, "))", " ) )")

	log.Infof("Cleaned Query: %v", query)

	return query
}

func fmtlines(query string) string {
	var (
		fmtquery string
		counter  int //to track brackets
	)

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
		log.Debugf("Counter: %v", counter)
		log.Debugf("Word: %v\t%v\n", word, len(word))
		if i < length-1 {
			nextword = strings.TrimSpace(words[i+1])
			phrase = word + space + nextword
			// log.Debugf("Phrase :%v", phrase)
		}

		// Handle brackets first
		if strings.Contains(word, "(") {
			if word == "(" {
				log.Debug("Open Bracket")
				fmtquery += word + newline
			} else {
				log.Debug("Func Open Bracket")
				counter++
				if strings.Contains(word, ")") {
					counter--
				}
				fmtquery += word + space
			}
		} else if strings.Contains(word, ")") {
			if strings.HasSuffix(word, ")") && counter > 0 {
				log.Debug("Func Close Bracket")
				counter--
				fmtquery += word + space
			} else if strings.HasSuffix(word, "),") && counter > 0 {
				log.Debug("Func Close Bracket")
				counter--
				fmtquery += word + space + newline
			} else if strings.HasSuffix(word, ")") {
				log.Debug("Close Bracket")
				n := strings.LastIndex(word, ")")
				w := word[0:n]

				fmtquery += w + newline + ")" + space
			} else if strings.HasSuffix(word, "),") {
				log.Debug("Close Bracket")
				n := strings.LastIndex(word, "),")
				w := word[0:n]

				fmtquery += w + newline + ")," + space
			}
		} else if strings.HasSuffix(word, ",") {
			log.Debug("Column")
			fmtquery += word + newline
		} else if selectk.MatchString(word) {
			log.Debug("Select")
			fmtquery += newline + word + newline
		} else if from.MatchString(word) {
			log.Debug("From|Where|Qualify|Having")
			fmtquery += newline + word + newline
		} else if word == "case" {
			log.Debug("Case")
			fmtquery += newline + word + newline
		} else if prepositions.MatchString(word) {
			log.Debug("On|Or|And|End|Else")
			fmtquery += newline + word + space
		} else if joins.MatchString(phrase) {
			log.Debug("Joins")
			// exceptional 3 word case
			if phrase == "full outer" {
				fmtquery += newline + phrase + space + "join" + space
				i += 3
			} else {
				fmtquery += newline + phrase + space
				i += 2
			}
			continue
		} else if groupby.MatchString(phrase) {
			log.Debug("Group|Order by")
			fmtquery += newline + phrase + newline
			i += 2
			continue
		} else if i == length-1 {
			log.Debug("Last Word")
			fmtquery += word + newline
		} else {
			log.Debug("Default")
			fmtquery += word + space
		}
		i++
	}

	fmtquery = strings.ReplaceAll(fmtquery, "\n\n", "\n")

	return fmtquery
}

func findtab(query string) string {
	var tab string
	counter := 1
	lines := strings.Split(query, "\n")

	for i := len(lines) - 1; i >= 0; i-- {
		line := lines[i]
		log.Debugf("findtab: %v %v", line, counter)

		for c := len(line) - 1; c >= 0; c-- {
			w := string(line[c])
			if w == ")" {
				counter++
			} else if w == "(" {
				counter--
			}
			if counter == 0 {
				tab = regexp.MustCompile(`^(\s*)`).FindString(line)
				return tab
			}
		}
	}

	return tab
}

func fmttabs(query string) string {
	var fmtdquery string
	// read each line
	lines := strings.Split(query, "\n")

	for i := range lines {
		line := lines[i]
		log.Debugf("fmttabs:Line:Tablen: %v:%v", line, len(tab))

		if selectk.MatchString(line) {
			fmtdquery += tab + line + newline
			tab += "  "
		} else if from.MatchString(line) || groupby.MatchString(line) {
			tab = strings.Replace(tab, "  ", "", 1)
			fmtdquery += tab + line + newline
			tab += "  "
		} else if strings.HasPrefix(line, ")") || strings.HasPrefix(line, "),") {
			// find corresponding opening line and get tabspace length
			tab = findtab(fmtdquery)
			// tab = strings.Replace(tab, "  ", "", 2)
			fmtdquery += tab + line + newline
			if strings.HasSuffix(line, "(") {
				tab += "  "
			}
		} else if strings.Contains(line, "(") && !strings.Contains(line, ")") {
			fmtdquery += tab + line + newline
			tab += "  "
		} else if strings.Contains(line, ")") && !strings.Contains(line, "(") {
			fmtdquery += tab + line + newline
			tab = strings.Replace(tab, "  ", "", 1)
		} else if strings.HasSuffix(line, "(") {
			fmtdquery += tab + line + newline
			tab += "  "
		} else {
			fmtdquery += tab + line + newline
		}
	}

	return fmtdquery
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
				fmtdquery += fmtlines(query) //formatquery is responsible for adding 2 newlines at the end of fmtdquery that it returns
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
		fmtdquery += fmtlines(query)
	}

	fmtdquery = fmttabs(fmtdquery)

	return fmtdquery
}

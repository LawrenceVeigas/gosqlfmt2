package gosqlfmt2

import "regexp"

const (
	space   = " "
	newline = "\n"
)

var tab = ""

// regexps
var (
	whitespace    = regexp.MustCompile(`\s{2,}`)
	comment       = regexp.MustCompile(`^--.*`)
	inlinecomment = regexp.MustCompile(`(.+?)(--.*)`)
	newlineword   = regexp.MustCompile(`\n([a-zA-Z='0-9]*)`)
	selectk       = regexp.MustCompile(`(?i)select`)
	from          = regexp.MustCompile(`(?i)(from|where|qualify|having)`)
	groupby       = regexp.MustCompile(`(?i)(group\sby|order\sby)`)
	spacecomma    = regexp.MustCompile(`\s,`)
	commaword     = regexp.MustCompile(`,(\w*)`)
	bracket       = regexp.MustCompile(`(?i)\(`)
	prepositions  = regexp.MustCompile(`(?i)(^on$|^or$|^and$|^end$|^else$)`)
	joins         = regexp.MustCompile(`(?i)(left join|right join|inner join|join|full outer)`)
)

// keywords
var (
	SELECT    = "SELECT"
	FROM      = "FROM"
	WHERE     = "WHERE"
	AS        = "AS"
	GROUPBY   = "GROUP BY"
	ORDERBY   = "ORDER BY"
	HAVING    = "HAVING"
	PARTITION = "PARTITION BY"
	OVER      = "OVER"
	ASC       = "ASC"
	DESC      = "DESC"
	JOIN      = "JOIN"
	FULLJOIN  = "FULL OUTER JOIN"
	INNERJOIN = "INNER JOIN"
	LEFTJOIN  = "LEFT JOIN"
	RIGHJOIN  = "RIGHT JOIN"

	KEYWORDS = []string{SELECT, FROM, WHERE, AS, GROUPBY, ORDERBY, HAVING, PARTITION, OVER, ASC, DESC, JOIN, FULLJOIN, INNERJOIN, LEFTJOIN, RIGHJOIN}
)

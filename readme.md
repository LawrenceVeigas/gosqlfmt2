#gosqlfmt2

A pet project to format sql as dialect agnostically as possible. Generally, formatters who try to format based on specific dialect endup introduce bugs into your sql code. This adds to the overhead of then fixing the bugs instead of focusing on the job at hand.

gosqlfmt2 isn't perfect (yet). It's a work in progess and I'll be working on smoothing out the rough edges over time.

## Usage

Run it by passing the sql filename as a cmdline argument. E.g.
`go run main.go test.sql`

> Note that, this will format your code in-place.

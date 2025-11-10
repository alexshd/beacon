package web

import (
"fmt"
g "maragu.dev/gomponents"
. "maragu.dev/gomponents/html"
)

type SudokuBoard [9][9]int

type SudokuStats struct {
	Filled int
	Valid  bool
	Solved bool
}

func Page(board SudokuBoard, stats SudokuStats, version string) g.Node {
	return HTML(
Head(
Meta(Charset("utf-8")),
Meta(Name("viewport"), Content("width=device-width, initial-scale=1")),
TitleEl(g.Text("Sudoku - Law I Blue-Green Deployment")),
Link(Rel("stylesheet"), Href("/static/sudoku.css")),
Script(Src("https://unpkg.com/htmx.org@2.0.4")),
),
Body(
Div(Class("container"),
Div(Class("header"),
H1(g.Text("Sudoku: Law I Demo")),
Span(Class("version-badge"), g.Text(version)),
P(
g.Text("Demonstrating "),
Span(Class("law-badge"), g.Text("Immutable")),
Span(Class("law-badge"), g.Text("Associative")),
Span(Class("law-badge"), g.Text("Commutative")),
g.Text(" operations with blue-green deployment"),
),
),
Div(ID("stats"), StatsComponent(stats)),
Div(ID("sudoku-board"), BoardComponent(board)),
Div(Class("controls"),
Button(Class("btn btn-success"), g.Attr("onclick", "refreshBoard()"), g.Text("Refresh")),
),
Div(Class("merge-section"),
H3(g.Text("🔄 Blue-Green Merge")),
P(Class("merge-info"), g.Text("Merge state from another server. Law I guarantees: Associative, Commutative, Idempotent.")),
Div(Class("controls"),
Button(Class("btn btn-secondary"), g.Attr("onclick", "mergeFrom('http://localhost:9000')"), g.Text("Merge from :9000")),
Button(Class("btn btn-secondary"), g.Attr("onclick", "mergeFrom('http://localhost:9001')"), g.Text("Merge from :9001")),
),
Div(ID("merge-status")),
),
Div(Class("footer"),
P(g.Text("Powered by lawtest • No YAML files were harmed • Group Theory FTW")),
),
),
Script(Src("/static/sudoku.js")),
),
)
}

func StatsComponent(stats SudokuStats) g.Node {
	validClass := "stat-value"
	if stats.Valid {
		validClass += " valid"
	} else {
		validClass += " invalid"
	}
	solvedText := "No"
	if stats.Solved {
		solvedText = "Yes!"
	}
	return Div(Class("stats"),
Div(Class("stat-card"),
Div(Class("stat-label"), g.Text("Filled")),
Div(Class("stat-value"), g.Textf("%d/81", stats.Filled)),
),
Div(Class("stat-card"),
Div(Class("stat-label"), g.Text("Valid")),
Div(Class(validClass), g.Text(boolToText(stats.Valid))),
),
Div(Class("stat-card"),
Div(Class("stat-label"), g.Text("Solved")),
Div(Class("stat-value"), g.Text(solvedText)),
),
)
}

func BoardComponent(board SudokuBoard) g.Node {
	cells := make([]g.Node, 0, 81)
	for row := 0; row < 9; row++ {
		for col := 0; col < 9; col++ {
			num := board[row][col]
			cellClass := "sudoku-cell"
			var cellText string
			if num != 0 {
				cellClass += " filled"
				cellText = fmt.Sprintf("%d", num)
			} else {
				cellClass += " empty"
			}
			cells = append(cells, Div(
Class(cellClass),
g.Attr("data-row", fmt.Sprintf("%d", row)),
g.Attr("data-col", fmt.Sprintf("%d", col)),
g.Text(cellText),
))
		}
	}
	return Div(Class("sudoku-board"), g.Group(cells))
}

func boolToText(b bool) string {
	if b {
		return "✓"
	}
	return "✗"
}

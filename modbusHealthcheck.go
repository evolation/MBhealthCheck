package main

import "github.com/maxence-charriere/go-app/v9/pkg/app"

type tablePage struct {
	app.Compo

	Rows [][]string
}

func (t *tablePage) Render() app.UI {
	headers := []app.UI{
		app.Tr().Body(
			app.Th().Text("Column 1"),
			app.Th().Text("Column 2"),
			app.Th().Text("Column 3"),
		),
	}

	return app.Div().Body(
		app.H1().Text("Table Page"),
		app.Table().Body(append(headers, t.renderRows()...)...),
		app.Button().
			Text("Add Row").
			OnClick(t.OnAddRowClick),
		app.Button().
			Text("Remove Row").
			OnClick(t.OnRemoveLastRowClick),
	)
}

func (t *tablePage) renderRows() []app.UI {
	var uiRows []app.UI
	for _, row := range t.Rows {
		uiRow := app.Tr().Body(
			app.Td().Text(row[0]),
			app.Td().Text(row[1]),
			app.Td().Text(row[2]),
		)
		uiRows = append(uiRows, uiRow)
	}
	return uiRows
}

func (t *tablePage) OnAddRowClick(ctx app.Context, e app.Event) {
	t.Rows = append(t.Rows, []string{"New Data 1", "New Data 2", "New Data 3"})
	t.Update()
}

func (t *tablePage) OnRemoveLastRowClick(ctx app.Context, e app.Event) {
	if len(t.Rows) > 0 {
		t.Rows = t.Rows[:len(t.Rows)-1]
		t.Update()
	}
}

package main

import (
	"fmt"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type tablePage struct {
	app.Compo
	Rows        [][]string
	RowToFocus  int
	ShouldFocus bool
}

func (t *tablePage) Render() app.UI {
	headers := []app.UI{
		app.Tr().Body(
			app.Th().Text("MB Address"),
			app.Th().Text("Status"),
			app.Th().Text("Result"),
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

func (t *tablePage) OnInputChange(row []string) func(ctx app.Context, e app.Event) {
	return func(ctx app.Context, e app.Event) {
		newValue := ctx.JSSrc().Get("value").String()
		row[0] = newValue
		t.Update()
	}
}

func (t *tablePage) OnInputKeyDown(ctx app.Context, e app.Event) {
	keyCode := e.Get("keyCode").Int()
	if keyCode == 13 { // 13 is the key code for Enter
		t.Rows = append(t.Rows, []string{"", "New Data 2", "New Data 3"})
		t.RowToFocus = len(t.Rows) - 1
		t.ShouldFocus = true
		t.Update()
	}
}

func (t *tablePage) renderRows() []app.UI {
	var uiRows []app.UI
	for i, row := range t.Rows {
		inputID := fmt.Sprintf("input-row-%d", i)
		fmt.Printf("inputID: %s\n", inputID)
		inputElem := app.Input().
			Type("text").
			Value(row[0]).
			OnChange(t.OnInputChange(row)).
			OnKeyDown(t.OnInputKeyDown).
			ID(inputID)

		uiRow := app.Tr().Body(
			app.Td().Body(inputElem),
			app.Td().Text(row[1]),
			app.Td().Text(row[2]),
		)
		uiRows = append(uiRows, uiRow)
	}
	return uiRows
}

func (t *tablePage) OnAddRowClick(ctx app.Context, e app.Event) {
	fmt.Println("OnAddRowClick triggered") // Add this line
	t.Rows = append(t.Rows, []string{"", "Pending", "Pending"})
	t.Update()
}

func (t *tablePage) OnRemoveLastRowClick(ctx app.Context, e app.Event) {
	if len(t.Rows) > 0 {
		t.Rows = t.Rows[:len(t.Rows)-1]
		t.Update()
	}
}

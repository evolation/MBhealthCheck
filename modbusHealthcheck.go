//go:build js && wasm

package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/gorilla/websocket"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type tablePage struct {
	app.Compo
	Rows        [][]string
	RowToFocus  int
	ShouldFocus bool
}
type RowData struct {
	SlaveID int `json:"slave_id"`
}

var ws *websocket.Conn

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
		app.Button().
			Text("Check Health").
			OnClick(t.OnCheckHealthClick),
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
func (t *tablePage) collectRowData() []RowData {
	var data []RowData
	for _, row := range t.Rows {
		slaveID, err := strconv.Atoi(row[0]) // Convert string to integer
		if err != nil {
			// Handle error, maybe log it or skip this row
			continue
		}
		data = append(data, RowData{
			SlaveID: slaveID,
		})
	}
	return data
}

func (t *tablePage) rowsToJSON() (string, error) {
	data := t.collectRowData()
	jsonData, err := json.MarshalIndent(data, "", "    ") // Indent for pretty printing
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

func (t *tablePage) OnCheckHealthClick(ctx app.Context, e app.Event) {
	jsonString, err := t.rowsToJSON()
	if err != nil {
		fmt.Println("Error converting rows to JSON:", err)
		return
	}
	// ioutil.WriteFile("./config.json", []byte(jsonString), 0644)

	// Instead of calling runModbusHealthcheckTarget directly,
	// send a WebSocket message to the backend.
	if ws == nil {
		// Initialize WebSocket connection.
		ws, _, err = websocket.DefaultDialer.Dial("ws://localhost:7999/ws", nil)
		if err != nil {
			log.Fatalf("WebSocket dial error: %v", err)
		}
		go func() {
			defer ws.Close()
			for {
				_, message, err := ws.ReadMessage()
				if err != nil {
					log.Println("WebSocket read error:", err)
					break
				}
				// Process the received message and update UI.
				fmt.Println(hex.Dump(message))
			}
		}()
	}
	ws.WriteMessage(websocket.TextMessage, []byte(jsonString))
}

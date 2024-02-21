package ui

import "github.com/charmbracelet/bubbles/textarea"

// TextArea
type TextArea struct {
	field textarea.Model
	data  string
}

// Card
func (c *Card) setTxtArea() {
	c.TxtArea.field.Prompt = ""
	c.TxtArea.field.Placeholder = "Card Description"
	c.TxtArea.field.ShowLineNumbers = true
	c.TxtArea.field.MaxHeight = 7
	c.TxtArea.field.MaxWidth = ws.width - 7
}

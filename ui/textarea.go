package ui

// TextArea

// Card
func (c *Card) setTxtArea() {
	c.textarea.Prompt = ""
	c.textarea.Placeholder = "Card Description"
	c.textarea.ShowLineNumbers = true
	c.textarea.SetValue(c.card.Description)
}

func (c *Card) handleTextArea(key string) {
	switch key {
	case "esc":
		c.textarea.Blur()
		return
	}
}

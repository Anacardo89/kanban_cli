package main

import (
	"fmt"
)

func main() {
	menu := TestData()
	data, err := menu.MarshalYAML()
	if err != nil {
		fmt.Println(err)
	}
	datastr := string(data)
	fmt.Println(datastr)
}

// func main() {
// 	f, err := tea.LogToFile("log.log", "error")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer f.Close()
// 	p := tea.NewProgram(ui.New(), tea.WithAltScreen())
// 	_, err = p.Run()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }

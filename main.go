package main

import (
	"fmt"
	"log"

	"github.com/Anacardo89/kanban_cli/kanban"
)

func main() {
	menu := TestData()
	fmt.Println(menu)
	project1, err := menu.Projects.WalkTo(0)
	if err != nil {
		log.Println(err)
	}
	project2, err := menu.Projects.WalkTo(1)
	if err != nil {
		log.Println(err)
	}
	project1Val := project1.GetVal().(kanban.Project)
	fmt.Println(project1Val.Lists)
	project2Val := project2.GetVal().(kanban.Project)
	fmt.Println(project2Val.Lists)
	project1list1, err := project1Val.Lists.WalkTo(0)
	if err != nil {
		log.Println(err)
	}
	project1list2, err := project1Val.Lists.WalkTo(1)
	if err != nil {
		log.Println(err)
	}
	project2list1, err := project2Val.Lists.WalkTo(0)
	if err != nil {
		log.Println(err)
	}
	project2list2, err := project2Val.Lists.WalkTo(1)
	if err != nil {
		log.Println(err)
	}
	project1list1Val := project1list1.GetVal().(kanban.List)
	fmt.Println(project1list1Val.Cards)
	project1list2Val := project1list2.GetVal().(kanban.List)
	fmt.Println(project1list2Val.Cards)
	project2list1Val := project2list1.GetVal().(kanban.List)
	fmt.Println(project2list1Val.Cards)
	project2list2Val := project2list2.GetVal().(kanban.List)
	fmt.Println(project2list2Val.Cards)
	sm := menu.ToStorage()
	fmt.Println(sm.ToYAML())
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

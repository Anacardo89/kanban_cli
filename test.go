package main

import (
	"log"

	"github.com/Anacardo89/kanban_cli/kanban"
	"github.com/charmbracelet/lipgloss"
)

func TestData() *kanban.Menu {
	menu := kanban.StartMenu()
	menu.AddProject("project1")
	project1, err := menu.Projects.WalkTo(0)
	if err != nil {
		log.Println(err)
	}
	project1Val := project1.GetValPtr()
	proj1 := (*project1Val).(kanban.Project)
	proj1.AddList("List 1")
	menu.AddProject("project2")
	project2, err := menu.Projects.WalkTo(1)
	if err != nil {
		log.Println(err)
	}
	project2Val := project2.GetVal().(kanban.Project)
	project2Val.AddList("List 1")
	project1Val.AddList("List 2")
	project2Val.AddList("List 2")
	project1Val.AddLabel("Label 1", lipgloss.Color("30"))
	project1Val.AddLabel("Label 2", lipgloss.Color("50"))
	project2Val.AddLabel("Label 1", lipgloss.Color("30"))
	project2Val.AddLabel("Label 2", lipgloss.Color("50"))
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
	project1list2Val := project1list2.GetVal().(kanban.List)
	project2list1Val := project2list1.GetVal().(kanban.List)
	project2list2Val := project2list2.GetVal().(kanban.List)
	project1list1Val.AddCard("Card 1")
	project1list1Val.AddCard("Card 2")
	project1list2Val.AddCard("Card 1")
	project1list2Val.AddCard("Card 2")
	project2list1Val.AddCard("Card 1")
	project2list1Val.AddCard("Card 2")
	project2list2Val.AddCard("Card 1")
	project2list2Val.AddCard("Card 2")
	return menu
}

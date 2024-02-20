package kanban

import (
	"log"
)

func TestData() *Menu {
	menu := StartMenu()
	menu.AddProject("project1")
	project1, err := menu.Projects.WalkTo(0)
	if err != nil {
		log.Println(err)
	}
	project1Val := project1.Val().(*Project)
	project1Val.AddBoard("Board 1")
	menu.AddProject("project2")
	project2, err := menu.Projects.WalkTo(1)
	if err != nil {
		log.Println(err)
	}
	project2Val := project2.Val().(*Project)
	project2Val.AddBoard("Board 1")
	project1Val.AddBoard("Board 2")
	project2Val.AddBoard("Board 2")
	project1Val.AddLabel("Label 1", "30")
	project1Val.AddLabel("Label 2", "50")
	project2Val.AddLabel("Label 1", "30")
	project2Val.AddLabel("Label 2", "50")
	project1Board1, err := project1Val.Boards.WalkTo(0)
	if err != nil {
		log.Println(err)
	}
	project1Board2, err := project1Val.Boards.WalkTo(1)
	if err != nil {
		log.Println(err)
	}
	project2Board1, err := project2Val.Boards.WalkTo(0)
	if err != nil {
		log.Println(err)
	}
	project2Board2, err := project2Val.Boards.WalkTo(1)
	if err != nil {
		log.Println(err)
	}
	project1Board1Val := project1Board1.Val().(*Board)
	project1Board2Val := project1Board2.Val().(*Board)
	project2Board1Val := project2Board1.Val().(*Board)
	project2Board2Val := project2Board2.Val().(*Board)
	project1Board1Val.AddCard("Card 1")
	project1Board1Val.AddCard("Card 2")
	project1Board2Val.AddCard("Card 1")
	project1Board2Val.AddCard("Card 2")
	project2Board1Val.AddCard("Card 1")
	project2Board1Val.AddCard("Card 2")
	project2Board2Val.AddCard("Card 1")
	project2Board2Val.AddCard("Card 2")
	return menu
}

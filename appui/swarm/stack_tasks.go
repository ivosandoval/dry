package swarm

import (
	"fmt"

	gizaktermui "github.com/gizak/termui"
	"github.com/moncho/dry/appui"
	"github.com/moncho/dry/docker"
	"github.com/moncho/dry/ui"
)

//StacksTasksWidget shows a service's task information
type StacksTasksWidget struct {
	TasksWidget
	stack string
}

//NewStacksTasksWidget creates a StacksTasksWidget
func NewStacksTasksWidget(swarmClient docker.SwarmAPI, y int) *StacksTasksWidget {
	w := StacksTasksWidget{
		TasksWidget: TasksWidget{
			swarmClient:   swarmClient,
			header:        defaultTasksTableHeader,
			height:        appui.MainScreenAvailableHeight(),
			mounted:       false,
			offset:        0,
			selectedIndex: 0,
			x:             0,
			y:             y,
			sortMode:      docker.SortByTaskService,
			tableTitle:    createStackTableTitle(),
			width:         ui.ActiveScreen.Dimensions.Width},
	}
	return &w
}

//Buffer returns the content of this widget as a termui.Buffer
func (s *StacksTasksWidget) Buffer() gizaktermui.Buffer {
	s.Lock()
	defer s.Unlock()
	y := s.y
	buf := gizaktermui.NewBuffer()
	if s.mounted {
		s.prepareForRendering()
		var filter string
		if s.filterPattern != "" {
			filter = fmt.Sprintf(
				"<b><blue> | Active filter: </><yellow>%s</></> ", s.filterPattern)
		}
		s.tableTitle.Content(fmt.Sprintf(
			"<b><blue>Stack %s tasks: </><yellow>%d</></>", s.stack, s.RowCount()) + " " + filter)

		s.tableTitle.Y = y
		buf.Merge(s.tableTitle.Buffer())
		y += s.tableTitle.GetHeight()

		s.updateHeader()
		s.header.SetY(y)
		buf.Merge(s.header.Buffer())
		y += s.header.GetHeight()

		selected := s.selectedIndex - s.startIndex

		for i, serviceRow := range s.visibleRows() {
			serviceRow.SetY(y)
			y += serviceRow.GetHeight()
			if i != selected {
				serviceRow.NotHighlighted()
			} else {
				serviceRow.Highlighted()
			}
			buf.Merge(serviceRow.Buffer())
		}
	}
	return buf
}

//ForStack sets the stack for which this widget is showing tasks
func (s *StacksTasksWidget) ForStack(stack string) {
	s.Lock()
	defer s.Unlock()

	s.stack = stack
	s.mounted = false
	s.sortMode = docker.SortByTaskService

}

//Mount prepares this widget for rendering
func (s *StacksTasksWidget) Mount() error {
	s.Lock()
	defer s.Unlock()
	if !s.mounted {
		s.mounted = true
		var rows []*TaskRow
		if tasks, err := s.swarmClient.StackTasks(s.stack); err == nil {
			for _, task := range tasks {
				rows = append(rows, NewTaskRow(s.swarmClient, task, s.header))
			}
			s.totalRows = rows
		} else {
			return err
		}
	}
	s.align()
	return nil
}

//Name returns this widget name
func (s *StacksTasksWidget) Name() string {
	return "StacksTasksWidget"
}

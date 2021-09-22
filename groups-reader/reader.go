package groupsreader

// GroupsReader is an interface that gets the info of a group
type GroupsReader interface {
	Members(grpName string) ([]string, error)
}



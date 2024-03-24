package radix

type WorkItem struct {
	ID string
}

type Organization struct {
	Name      string
	WorkItems []WorkItem
}

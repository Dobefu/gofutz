package routes

// SortOption represents a sort option for the dropdown.
type SortOption struct {
	Value string
	Label string
}

// PageVars defines the variables of a page.
type PageVars struct {
	Title              string
	SortOptions        []SortOption
	SelectedSortOption string
}

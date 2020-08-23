package operator

//Query Modifiers
// refer:https://docs.mongodb.com/manual/reference/operator/query-modifier/
const (
	// Modifiers
	Explain     = "$explain"
	Hint        = "$hint"
	MaxTimeMS   = "$maxTimeMS"
	OrderBy     = "$orderby"
	Query       = "$query"
	ReturnKey   = "$returnKey"
	ShowDiskLoc = "$showDiskLoc"

	// Sort Order
	Natural = "$natural"
)

package operator

// define the update operators
// refer: https://docs.mongodb.com/manual/reference/operator/update/
const (
	// Fields
	CurrentDate = "$currentDate"
	Inc         = "$inc"
	Min         = "$min"
	Max         = "$max"
	Mul         = "$mul"
	Rename      = "$rename"
	Set         = "$set"
	SetOnInsert = "$setOnInsert"
	Unset       = "$unset"

	// Array Operators
	AddToSet = "$addToSet"
	Pop      = "$pop"
	Pull     = "$pull"
	Push     = "$push"
	PullAll  = "$pullAll"

	// Array modifiers
	Each     = "$each"
	Position = "$position"
	Sort     = "$sort"

	// Array bitwise
	Bit = "$bit"
)

package operator

// define the aggregation pipeline stages
// refer: https://docs.mongodb.com/manual/reference/operator/aggregation-pipeline/
const (
	// Collection Aggregate Stages
	AddFields      = "$addFields"
	Bucket         = "$bucket"
	BucketAuto     = "$bucketAuto"
	CollStats      = "$collStats"
	Count          = "$count"
	Facet          = "$facet"
	GeoNear        = "$geoNear"
	GraphLookup    = "$graphLookup"
	Group          = "$group"
	IndexStats     = "$indexStats"
	Limit          = "$limit"
	ListSessions   = "$listSessions"
	Lookup         = "$lookup"
	Match          = "$match"
	Merge          = "$merge"
	Out            = "$out"
	PlanCacheStats = "$planCacheStats"
	Project        = "$project"
	Redact         = "$redact"
	ReplaceRoot    = "$replaceRoot"
	ReplaceWith    = "$replaceWith"
	Sample         = "$sample"
	Skip           = "$skip"
	SortByCount    = "$sortByCount"
	UnionWith      = "$unionWith"
	Unwind         = "$unwind"

	// Database Aggregate stages
	CurrentOp         = "$currentOp"
	ListLocalSessions = "$listLocalSessions"
)

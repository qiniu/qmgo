package operator

type OpType string

const (
	BeforeInsert  OpType = "beforeInsert"
	AfterInsert   OpType = "afterInsert"
	BeforeUpdate  OpType = "beforeUpdate"
	AfterUpdate   OpType = "afterUpdate"
	BeforeQuery   OpType = "beforeQuery"
	AfterQuery    OpType = "afterQuery"
	BeforeRemove  OpType = "beforeRemove"
	AfterRemove   OpType = "afterRemove"
	BeforeUpsert  OpType = "beforeUpsert"
	AfterUpsert   OpType = "afterUpsert"
	BeforeReplace OpType = "beforeReplace"
	AfterReplace  OpType = "afterReplace"
)

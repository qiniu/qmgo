/*
 Copyright 2020 The Qmgo Authors.
 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at
     http://www.apache.org/licenses/LICENSE-2.0
 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

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

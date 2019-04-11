// Code generated by "stringer -type=Change,ChangedStats"; DO NOT EDIT.

package main

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[Unchanged-1]
	_ = x[Added-2]
	_ = x[StatsChanged-3]
}

const _Change_name = "UnchangedAddedStatsChanged"

var _Change_index = [...]uint8{0, 9, 14, 26}

func (i Change) String() string {
	i -= 1
	if i < 0 || i >= Change(len(_Change_index)-1) {
		return "Change(" + strconv.FormatInt(int64(i+1), 10) + ")"
	}
	return _Change_name[_Change_index[i]:_Change_index[i+1]]
}
func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[ChangedModTime-1]
	_ = x[ChangedSize-2]
}

const _ChangedStats_name = "ChangedModTimeChangedSize"

var _ChangedStats_index = [...]uint8{0, 14, 25}

func (i ChangedStats) String() string {
	i -= 1
	if i < 0 || i >= ChangedStats(len(_ChangedStats_index)-1) {
		return "ChangedStats(" + strconv.FormatInt(int64(i+1), 10) + ")"
	}
	return _ChangedStats_name[_ChangedStats_index[i]:_ChangedStats_index[i+1]]
}

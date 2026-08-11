package main

import (
	"reflect"
	"sort"
	"testing"
)

func TestArchiveEntryPriorityPlacesDSMPackageMetadataFirst(t *testing.T) {
	entries := []string{
		"PACKAGE_ICON_256.PNG",
		"conf/privilege",
		"package.tgz",
		"scripts/postinst",
		"INFO",
		"conf",
		"scripts",
		"PACKAGE_ICON.PNG",
	}
	sort.Slice(entries, func(i, j int) bool {
		left, right := archiveEntryPriority(entries[i]), archiveEntryPriority(entries[j])
		if left != right {
			return left < right
		}
		return entries[i] < entries[j]
	})
	want := []string{
		"INFO",
		"package.tgz",
		"scripts",
		"scripts/postinst",
		"conf",
		"conf/privilege",
		"PACKAGE_ICON.PNG",
		"PACKAGE_ICON_256.PNG",
	}
	if !reflect.DeepEqual(entries, want) {
		t.Fatalf("ordered entries = %v, want %v", entries, want)
	}
}

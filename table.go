package main

import (
	"fmt"
	"github.com/coreos/go-semver/semver"
	"io"
	"sort"
	"strings"
)

func WriteTables(w io.Writer, compatibility *Compatibility) error {
	fmt.Fprintln(w, "## By Runtime version:")
	writeTable(w, buildRuntimeTable(compatibility.Runtime))
	fmt.Fprintln(w)

	sdks := make([]string, 0)
	for sdk := range compatibility.SDKs {
		sdks = append(sdks, sdk)
	}
	sort.Strings(sdks)

	for _, sdk := range sdks {
		fmt.Fprintf(w, "## By %v SDK version:\n", sdk)
		writeTable(w, buildSDKTable(sdk, compatibility.SDKs[sdk]))
		fmt.Fprintln(w)
	}

	return nil
}

func buildRuntimeTable(compatibility map[string]map[string]semver.Versions) [][]string {
	sortedVersions := make(semver.Versions, 0)
	for version := range compatibility {
		sortedVersions = append(sortedVersions, semver.New(version))
	}
	semver.Sort(sortedVersions)

	sdks := make([]string, 0)
	for sdk := range compatibility[sortedVersions[0].String()] {
		sdks = append(sdks, sdk)
	}
	sort.Strings(sdks)

	versionColumn := []string{"Runtime"}
	for _, version := range sortedVersions {
		versionColumn = append(versionColumn, version.String())
	}

	table := [][]string{versionColumn}
	for _, sdk := range sdks {
		sdkColumn := []string{sdk + " SDK"}

		for _, version := range sortedVersions {
			sdkColumn = append(sdkColumn, versionsToRange(compatibility[version.String()][sdk]))
		}

		table = append(table, sdkColumn)
	}

	return table
}

func buildSDKTable(sdk string, compatibility map[string]semver.Versions) [][]string {
	ranges := make(map[string]string)
	for version, runtimeVersions := range compatibility {
		ranges[version] = versionsToRange(runtimeVersions)
	}

	versionColumn := make([]string, 0)
	runtimeColumn := make([]string, 0)

	sortedVersions := make(semver.Versions, 0)
	for version := range ranges {
		sortedVersions = append(sortedVersions, semver.New(version))
	}
	semver.Sort(sortedVersions)

	for _, version := range sortedVersions {
		compatibleRange := ranges[version.String()]
		versionColumn = append(versionColumn, version.String())
		runtimeColumn = append(runtimeColumn, compatibleRange)
	}

	return [][]string{
		append([]string{sdk + " SDK"}, versionColumn...),
		append([]string{"Runtime"}, runtimeColumn...),
	}
}

func versionsToRange(versions semver.Versions) string {
	switch len(versions) {
	case 0:
		return ""
	case 1:
		return versions[0].String()
	default:
		return versions[0].String() + " - " + versions[len(versions)-1].String()
	}
}

func writeTable(w io.Writer, table [][]string) error {
	table = compactTable(table)
	rows := len(table[0])
	width := 0
	for _, column := range table {
		for _, value := range column {
			if len(value) > width {
				width = len(value)
			}
		}
	}

	fmt.Fprint(w, "|")
	for _, column := range table {
		fmt.Fprintf(w, " %v ", padRight(column[0], " ", width))
		fmt.Fprint(w, "|")
	}
	fmt.Fprintln(w)

	fmt.Fprint(w, "|")
	for range table {
		fmt.Fprint(w, padRight("", "-", width+2))
		fmt.Fprint(w, "|")
	}
	fmt.Fprintln(w)

	for i := 1; i < rows; i++ {
		fmt.Fprint(w, "|")
		for _, column := range table {
			fmt.Fprintf(w, " %v ", padRight(column[i], " ", width))
			fmt.Fprint(w, "|")
		}
		fmt.Fprintln(w)
	}

	return nil
}

func compactTable(table [][]string) [][]string {
	compacted := make([][]string, 0)
	for _, column := range table {
		compacted = append(compacted, []string{column[0]})
	}

	columns := len(table)
	rows := len(table[0])

	for row := 1; row < rows; row++ {
		rowHasContent := false
		for column := 1; column < columns; column++ {
			if len(table[column][row]) > 0 {
				rowHasContent = true
			}
		}

		if !rowHasContent {
			continue
		}

		lastDuplicateRow := row
	duplicates:
		for {
			if lastDuplicateRow+1 == rows {
				break
			}
			for column := 1; column < columns; column++ {
				if table[column][lastDuplicateRow] != table[column][row] {
					break duplicates
				}
			}
			lastDuplicateRow += 1
		}
		lastDuplicateRow -= 1

		if lastDuplicateRow != rows && lastDuplicateRow > row {
			compacted[0] = append(compacted[0], table[0][row]+" - "+table[0][lastDuplicateRow])
			row = lastDuplicateRow
		} else {
			compacted[0] = append(compacted[0], table[0][row])
		}

		for column := 1; column < columns; column++ {
			compacted[column] = append(compacted[column], table[column][row])
		}
	}

	return compacted
}

func padRight(value, padding string, width int) string {
	return value + strings.Repeat(padding, width-len(value))
}

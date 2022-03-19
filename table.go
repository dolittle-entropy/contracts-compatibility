package main

import (
	"fmt"
	"github.com/coreos/go-semver/semver"
	"github.com/olekukonko/tablewriter"
	"io"
	"sort"
)

// WriteTables writes Markdown tables for Runtime and SDK compatibilities from the given Compatibility to the provided Writer
func WriteTables(w io.Writer, compatibility *Compatibility) {
	fmt.Fprintln(w, "## By Runtime version:")
	table := createMarkdownTable(w)
	fillRuntimeTable(table, compatibility.Runtime)
	table.Render()

	sdks := make([]string, 0)
	for sdk := range compatibility.SDKs {
		sdks = append(sdks, sdk)
	}
	sort.Strings(sdks)

	for _, sdk := range sdks {
		fmt.Fprintf(w, "\n## By %v SDK version:\n", sdk)
		table := createMarkdownTable(w)
		fillSDKTable(table, sdk, compatibility.SDKs[sdk])
		table.Render()
	}
}

func createMarkdownTable(w io.Writer) *tablewriter.Table {
	table := tablewriter.NewWriter(w)
	table.SetBorders(tablewriter.Border{
		Left:   true,
		Top:    false,
		Right:  true,
		Bottom: false,
	})
	table.SetAutoFormatHeaders(false)
	table.SetCenterSeparator("|")
	return table
}

func fillRuntimeTable(table *tablewriter.Table, compatibility map[string]map[string]semver.Versions) {
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

	headers := []string{"Runtime"}
	for _, sdk := range sdks {
		headers = append(headers, sdk+" SDK")
	}
	table.SetHeader(headers)

	rows := make([][]semver.Versions, 0)
	for _, version := range sortedVersions {
		row := []semver.Versions{{version}}
		rowHasContents := false
		for _, sdk := range sdks {
			sdkVersions := compatibility[version.String()][sdk]
			if len(sdkVersions) > 0 {
				rowHasContents = true
			}
			row = append(row, sdkVersions)
		}
		if rowHasContents {
			rows = append(rows, row)
		}
	}

	table.AppendBulk(renderVersionTable(compactVersionTable(rows)))
}

func fillSDKTable(table *tablewriter.Table, sdk string, compatibility map[string]semver.Versions) {
	table.SetHeader([]string{sdk + " SDK", "Runtime"})

	sortedVersions := make(semver.Versions, 0)
	for version := range compatibility {
		sortedVersions = append(sortedVersions, semver.New(version))
	}
	semver.Sort(sortedVersions)

	rows := make([][]semver.Versions, 0)
	for _, version := range sortedVersions {
		rows = append(rows, []semver.Versions{
			{version},
			compatibility[version.String()],
		})
	}

	table.AppendBulk(renderVersionTable(compactVersionTable(rows)))
}

func compactVersionTable(table [][]semver.Versions) [][]semver.Versions {
	ranges := make([][]semver.Versions, len(table))
	for i, row := range table {
		ranges[i] = make([]semver.Versions, len(table[i]))
		for j, column := range row {
			if len(column) > 2 {
				ranges[i][j] = semver.Versions{column[0], column[len(column)-1]}
			} else {
				ranges[i][j] = column
			}
		}
	}

	compacted := make([][]semver.Versions, len(ranges))
	for i := 0; i < len(ranges); {
		row := ranges[i]
		compacted[i] = row

		didSkip := false
	skipRow:
		for i += 1; i < len(ranges); i++ {
			for j := 1; j < len(row); j++ {
				if !rangeEquals(row[j], ranges[i][j]) {
					break skipRow
				}
			}
			didSkip = true
		}

		if didSkip {
			row[0] = append(row[0], ranges[i-1][0][0])
		}

	}

	return compacted
}

func rangeEquals(left, right semver.Versions) bool {
	if len(left) != len(right) {
		return false
	}
	for i := 0; i < len(left); i++ {
		if !left[i].Equal(*right[i]) {
			return false
		}
	}
	return true
}

func renderVersionTable(table [][]semver.Versions) [][]string {
	rendered := make([][]string, len(table))
	for i, row := range table {
		rendered[i] = make([]string, len(table[i]))
		for j, column := range row {
			switch len(column) {
			case 1:
				rendered[i][j] = column[0].String()
			case 2:
				rendered[i][j] = column[0].String() + " - " + column[1].String()
			default:
				rendered[i][j] = ""
			}
		}
	}
	return rendered
}

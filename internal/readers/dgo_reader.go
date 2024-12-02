package readers

import (
	"bytes"
	"errors"
	"strconv"
	"strings"
	"unicode/utf8"
)

// Note: the `dgo` prefix in this file indicates the pertenence
// to the dgoReader Reader

// The current config for pdftotext will produce
// an avg of 6 leading blank characters for files on REGION_DGO
const dgoAvgLeadingWhitespace = 6

// const dgoAvgMaxColLens = [3]int{7,}

// Modest estimation of the max length of the column
// containing the index of the case
const dgoEstimateIndexColLen = 5

func dgoReader(data *[]byte) (caseTable *CaseTable, err error) {
	rows := bytes.Split(*data, []byte{'\n'})
	if len(rows) == 0 {
		return nil, errors.New("Data produced no rows")
	}

	var (
		colLens = [3]int{0, 0, 0}

		countingColLen = true
		parsingCase    = false

		caseIdxMap = map[string]bool{}
		tempCols   = [4][]byte{}
	)

	caseTable = NewCaseTable()
	tempCaseData := CaseData{}
	leftTrimCount := dgoAvgLeadingWhitespace + 1

	for rowNo, rowCount := 0, len(rows); rowNo < rowCount; rowNo++ {
		if len(rows[rowNo]) == 0 {
			continue
		}

		row := rows[rowNo][leftTrimCount:]

		if !parsingCase && !dgoIsCaseRow(row) {
			// If we aren't parsing and the current line is not a Case Row
			// we skip it
			continue
		}
		// If the current line is a Case Row we can begin parsing
		parsingCase = true

		tempCols[0] = []byte{}
		tempCols[1] = []byte{}
		tempCols[2] = []byte{}
		tempCols[3] = []byte{}
		currCol := 0

		for pos, rowLen := 0, len(row); pos < rowLen; pos++ {
			tempCols[currCol] = utf8.AppendRune(tempCols[currCol], rune(row[pos]))

			// Once we hit the last column, there no need for checks
			// just save the bytes until the row's end
			if currCol == 3 {
				continue
			}

			if !countingColLen && len(tempCols[currCol]) == colLens[currCol] {
				// If the length of the columns is set we use the length
				// to determine when to change column
				currCol++
				continue
			}

			colLens[currCol]++
			if dgoIsColumnSeparator(row, pos) {
				currCol++
			}
		}

		// There is a chance the file fetched has duplicate pages
		//
		// For such cases we'll skip parsing cases with indexes that exist
		// in `caseIdxMap`
		if rowNo == 0 {
			caseIdx := strings.TrimSpace(string(tempCols[0]))

			if caseRead := caseIdxMap[caseIdx]; caseRead {
				parsingCase = false
				continue
			}

			// If the case hasn't been read, add it to the map
			caseIdxMap[caseIdx] = true
		}

		tempCaseData.CaseId += string(tempCols[1]) + "\n"
		tempCaseData.Nature += string(tempCols[2]) + "\n"
		tempCaseData.Accord += string(tempCols[3]) + "\n"

		if dgoNextLineEndsParsing(rows, rowNo) {
			parsingCase = false

			tempCaseData.CaseId = strings.TrimSpace(tempCaseData.CaseId)
			tempCaseData.Nature = strings.TrimSpace(tempCaseData.Nature)
			tempCaseData.Accord = strings.TrimSpace(tempCaseData.Accord)

			caseRow, err := NewCaseRow(&tempCaseData)
			if err != nil {
				// For cases that do not produce a valid caseRow
				// save them as unparsed for possible force search
				cloned := tempCaseData.Clone()
				caseTable.UnparsedCases = append(caseTable.UnparsedCases, &cloned)
			}

			caseTable.Cases = append(caseTable.Cases, caseRow)
			tempCaseData.Clear()
		}
	}

	return
}

func dgoIsCaseRow(row []byte) bool {
	// Check the first few chars of the line to see if it is a case row
	lineStart := string(row[:dgoEstimateIndexColLen])
	_, parseErr := strconv.Atoi(lineStart)

	return parseErr == nil
}

func dgoNextLineEndsParsing(rows [][]byte, currRow int) bool {
	// No more lines | Empty lines | New Case Row lines
	// represent the end of a single CaseRow parsing
	return currRow+1 > len(rows) || len(rows[currRow+1]) == 0 || dgoIsCaseRow(rows[currRow+1])
}

// Checks if the current character counts as a column separator
//
// A character is considered a column separator if:
//
// - Both the previous and current character are a space [' ']
// and the next character is a non-whitespace character
func dgoIsColumnSeparator(row []byte, pos int) bool {
	char := row[pos]
	prevChar := safeCheckPrevIdx(row, pos)
	nextChar := safeCheckNextIdx(row, pos)

	return nextChar != ' ' && prevChar == ' ' && char == ' '

	// TODO: Search internet to make sure that `nextChar != byte('0')` is not necessary
	// or that it actually refers to comparison against the `nul` character
	// return nextChar != ' ' && nextChar != byte('0') && prevChar == ' ' && char == ' '
	//
	// Note: we shouldn't need to check anyway, since rarely will the parser
	// iterate the whole row slice
}

func safeCheckPrevIdx(chk []byte, currIdx int) (prevByte byte) {
	prevIdx := currIdx - 1
	if prevIdx < 0 {
		return
	}

	prevByte = chk[prevIdx]
	return
}
func safeCheckNextIdx(chk []byte, currIdx int) (nextByte byte) {
	nextIdx := currIdx + 1
	if nextIdx >= len(chk) {
		return
	}

	nextByte = chk[nextByte]
	return
}

package excel

import (
	"fmt"
	"io"

	"github.com/xuri/excelize/v2"
)

const sheetName = "Sheet1"

type BaseExporter struct {
	xlsxFile  *excelize.File
	writer    io.Writer
	Headers   []string
	headerSet bool
	rowNr     int
}

func NewBaseExporter(writer io.Writer) *BaseExporter {
	return &BaseExporter{
		xlsxFile:  excelize.NewFile(),
		writer:    writer,
		rowNr:     0,
		headerSet: false,
	}
}

func (x *BaseExporter) addRow(values []string) {

	//row numbers start at 1
	x.rowNr++
	colNr := 0
	for _, header := range values {
		colNr++
		columnPrefix, _ := excelize.ColumnNumberToName(colNr)
		x.xlsxFile.SetCellStr(
			sheetName,
			fmt.Sprintf("%s%d", columnPrefix, x.rowNr),
			header,
		)
	}
}

func (x *BaseExporter) addHeader() {
	x.addRow(x.Headers)
	x.headerSet = true
}

/*
Add formatting style to format text as regular text
cf. https://xuri.me/excelize/en/style.html#number_format
cf. https://github.com/qax-os/excelize/issues/1915
*/
func (x *BaseExporter) addStyles() {
	style, _ := x.xlsxFile.NewStyle(&excelize.Style{NumFmt: 49})
	var leftCol, rightCol string
	for i := range x.Headers {
		columnPrefix, _ := excelize.ColumnNumberToName(i + 1)
		if i == 0 {
			leftCol = columnPrefix
		}
		if i == len(x.Headers)-1 {
			rightCol = columnPrefix
		}
	}
	if len(x.Headers) > 0 {
		x.xlsxFile.SetColStyle(sheetName, fmt.Sprintf("%s:%s", leftCol, rightCol), style)
	}
}

func (x *BaseExporter) Add(row []string) {
	if !x.headerSet {
		x.addHeader()
		x.addStyles()
	}
	x.addRow(row)
}

func (x *BaseExporter) Flush() error {
	e := x.xlsxFile.Write(x.writer)
	if e != nil {
		return e
	}
	return x.xlsxFile.Close()
}

func (x *BaseExporter) GetContentType() string {
	return "application/vnd.ms-excel"
}

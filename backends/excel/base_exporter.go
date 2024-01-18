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

func (x *BaseExporter) Add(row []string) {
	if !x.headerSet {
		x.addHeader()
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

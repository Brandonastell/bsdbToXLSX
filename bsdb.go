package bsdbToXLSX

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

// container for sql rows (*sql.Row)  and sql colum types (*sql.ColumnType)
type Output struct {
	Rows  *sql.Rows
	Types []*sql.ColumnType
}

//container for each element comming from database
type Item struct {
	Type  string
	Value interface{}
}

//container for each row
type Row []*Item

//foundation for methods from excellize package
type Writer struct {
	*excelize.File
}

//use to build new excel writer
//return pointer to writer
func NewWriter() *Writer {
	w := &Writer{excelize.NewFile()}
	return w
}

// method to add or change a query to the QueryStuff structure
func (q *QueryStuff) SetQueryText(queryText string) {
	q.QueryText = queryText
}

//method on QueryStuff struct
//method returns Output object
func (q *QueryStuff) QueryRow() (Output, error) {
	db, err := sql.Open(q.DatabaseType, q.Loc)
	if err != nil {
		fmt.Println("row 45")
		log.Fatal(err)
	}
	defer db.Close()
	fmt.Println(q.QueryText)
	rows, err := db.Query(q.QueryText)
	switch {
	case err == sql.ErrNoRows:
		var err_out Output
		return err_out, errors.New("no rows returned")
	case err != nil:
		fmt.Println("row 55")
		log.Fatal(err)
	}
	types, err := rows.ColumnTypes()
	var output Output
	output.rows = rows
	output.types = types
	return output, nil
}

// converts numaric to alpha for excel cell location...
//if you have over column z let me know and ill make this less stupid
func ToCharStrConst(i int) string {
	asciiRune := rune(i + 65)
	asciiChar := fmt.Sprintf("%c", asciiRune)
	return asciiChar
}

//method for Writer object, returns slice of Row (each row)
func (w *Writer) WriteRow(rows []Row, fileName string) {
	for i, row := range rows {
		for j, item := range row {
			fmt.Printf("row = %+v\n", item.Value)
			switch item.Type {
			case "DATETIME", "DATE":
				w.SetCellValue("sheet1", fmt.Sprintf("%s%d", toCharStrConst(i), j+1), item.Value)
			default:
				w.SetCellValue("sheet1", fmt.Sprintf("%s%d", toCharStrConst(i), j+1), item.Value)
			}
		}

	}
	w.SaveAs(fileName)
}

//method for Output to package rows for excelize package and use in Writer.WriteRrow
//returns slice of Row
func (o *Output) PrepareRows() []Row {
	rows := o.rows
	types := o.types
	allRows := make([]Row, 0)
	for rows.Next() {
		items := make(Row, len(types))
		thisRow := make([]interface{}, len(types))
		thisRowPntr := make([]interface{}, len(types))
		for i, _ := range thisRow {
			thisRowPntr[i] = &thisRow[i]
		}
		rows.Scan(thisRowPntr...)
		for i, type_ := range types {
			items[i] = &Item{
				Type:  type_.DatabaseTypeName(),
				Value: thisRow[i],
			}
		}
		allRows = append(allRows, items)
	}
	return allRows
}

//func main() {
//	myQuery := QueryConfig("./test.yaml")
//	myQuery.setQueryText(`select name, options, date_added from reports where name = "test";`)
//	output, err := myQuery.QueryRow()
//	if err != nil {
//		log.Fatal(err)
//	}
//	allRows := output.prepareRows()
//	myWriter := NewWriter()
//	fmt.Println(len(allRows))
//	myWriter.WriteRow(allRows, "./test.xlsx")
//	//getQuery()
//}

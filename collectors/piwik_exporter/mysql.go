package piwik_exporter

import (
	"fmt"
	"github.com/ziutek/mymysql/mysql"
	_ "github.com/ziutek/mymysql/native"
	"io"
	"os"
)

func main() {
	user := "root"
	pass := ""
	dbname := "piwik"

	db := mysql.New("tcp", "", "127.0.0.1:3306", user, pass, dbname)

	err := db.Connect()
	if err != nil {
		panic(fmt.Errorf("Error from db.Connect: %s", err))
	}

	sql := "SELECT idvisit FROM piwik_log_visit LIMIT 10;"
	res, err := db.Start(sql)
	if err != nil {
		panic(fmt.Errorf("Error from db.Start with sql %s: %s", sql, err))
	}

	for _, field := range res.Fields() {
		fmt.Print(field.Name, " ")
	}
	fmt.Println()

	row := res.MakeRow()
	for {
		err := res.ScanRow(row)
		if err == io.EOF {
			// No more rows
			break
		}
		if err != nil {
			panic(fmt.Errorf("Error from GetRow: %s", err))
		}

		// Print all cols
		for _, col := range row {
			if col == nil {
				fmt.Print("<NULL>")
			} else {
				os.Stdout.Write(col.([]byte))
			}
			fmt.Print(" ")
		}
		fmt.Println()
	}
}

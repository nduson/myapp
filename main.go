package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"runtime"
	"strings"

	"database/sql"

	"github.com/Luxurioust/excelize"
	_ "github.com/lib/pq"
)

var db *sql.DB
var tpl *template.Template
var fpath string
var sheet_name string

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "1234567"
	dbname   = "staff_db"
)

type maplistbox struct {
	Id    int
	Value string
}

type Empinfo struct {
	employee_info string
	deductions    string
	month         string
	year          int
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {

	//fmt.Println("Program Statred.............Trying to Connect to DB")
	var err error

	/*psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err = sql.Open("postgres", psqlInfo)

	//db, err = sql.Open("postgres", "postgres://postgres:1234567@172.17.0.3/staff_db?sslmode=disable")
	//db, err = sql.Open("postgres", "postgres://postgres:@db/staff_db?sslmode=disable")
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		fmt.Println("Database connection failed due to", err)
		return
	}
	fmt.Println("Successfully connected To DB!") */

	pwd, err := os.Getwd()

	if err != nil {
		fmt.Println("Getting Host Enviroment Failed Due To: ", err)
		os.Exit(1)
	}
	///Output current working dir
	fmt.Println("The Enviroment path Is: ", pwd)

	fmt.Println("Loading Template From: ", pwd)

	////////// Load number of  sheet in the excel file into the drop down ////////////////////////////////////////////////

	//////////End load excel sheet into the drop down

	//tpl = template.Must(template.ParseFiles("goupload.html")) // added new

	tpl.ExecuteTemplate(w, "goupload.html", nil)

	fmt.Println("Successfully Loaded gouload.html Template !")
	// the FormFile function takes in the POST input id file
	file, header, err := r.FormFile("file")

	fmt.Println("Successfully Gotten form headers!: ", file, " ", header)

	if err != nil {
		//fmt.Fprintln(w, "Getting The File Failed Due To: ", err)
		return
	}

	if strings.Contains(header.Filename, ".xlsx") {
		fmt.Println("File Is Valid: ")
		checkfileagain := strings.Split(header.Filename, ".")
		//header.Header.Get("xlsx")
		/// futher check the file
		if checkfileagain[1] == "xlsx" {

			fmt.Println("File Is Valid: ")

			//defer file.Close()
			//path := "/Users/Nd/Documents/GoWorkSpace/src/excel/tmp/uploadedfile"

			//path := "/Users/Nd/Documents/GoWorkSpace/src/excel/tmp/"
			// get current working dir
			//pwd, err := os.Getwd()
			//if err != nil {
			//fmt.Println(err)
			//os.Exit(1)
			//}
			///Output current working dir
			//fmt.Println(pwd)

			/// check if path is windows or linux
			fmt.Println("Checking The Host Operating System:")
			var path string
			if runtime.GOOS == "windows" {
				path = pwd + `\uploads\`
				fmt.Println("The Host is Windows: ", path)
			}

			if runtime.GOOS == "linux" {
				path = pwd + `/uploads/`
				fmt.Println("The Host is Linux: ", path)
			}
			//path := pwd + `\uploads\`

			//out, err := os.Create("/tmp/uploadedfile")
			fmt.Println("Creating File Named.....: ", header.Filename)
			dst, err := os.Create(path + header.Filename)
			if err != nil {
				fmt.Fprint(w, "Unable to create the file for writing. Check your write access privilege: ", err)
				return
			}
			defer dst.Close()

			fmt.Println("Successfully Created File Called: ", header.Filename, " To be Copied To.....: ", path)

			// write the content from POST to the file
			_, err = io.Copy(dst, file)
			if err != nil {
				fmt.Fprintln(w, err)
			}

			fmt.Println("Successfully Created File: ", header.Filename, " And Copied To.....: ", path)

			//fmt.Fprintf(w, "File uploaded successfully : ", "Filename: ", header.Filename, " File Path: ", path)
			//fmt.Fprintf(w, "File uploaded successfully")
			//fmt.Fprintf(w, header.Filename)
			fpath = path + header.Filename
			fmt.Println("Get The File Full Path: ", path)
			//fmt.Fprintf(w, fpath)

			/// excel extracting starts here
			//return
			//xlsx, err := excelize.OpenFile("C:/Users/Nd/Documents/GoWorkSpace/src/excel/Workbook.xlsx")

			xlsx, err := excelize.OpenFile(fpath)
			if err != nil {
				fmt.Println("Excel Sheet failed to Load Due To: ", err)
				os.Exit(1)
			}

			sheet := xlsx.GetSheetMap()

			fmt.Println("Gotten the Number Of sheet As.............", sheet)

			//rows := xlsx.GetRows("Sheet1")

			fmt.Println("successfully Extracted the Whole excel Rows: ")
			// Gets the heading title
			//for _, row := range rows {
			// Form submitted
			col_indes := r.FormValue("col_index")

			var col_index int
			if _, err := fmt.Sscanf(col_indes, "%5d", &col_index); err == nil {
				fmt.Println(col_index) //
			}
			fmt.Println("successfully Gotten the Index for the Row Header: ", col_index) /////////////////////////////////

			for k, v := range sheet {
				fmt.Println("First Starts.............", k, v)

				qy := "INSERT INTO dump (name,uid)VALUES($1,$2)"
				//qy = strings.Replace(qy, "`", "\"", -1)
				//stmt, err := db.Prepare("INSERT `dtest` SET dname=$,damt=$")

				stmt, err := db.Prepare(qy)
				// prepare statement error
				if err != nil {
					panic(err)
					fmt.Fprintf(w, "The Following error Occured:", err)
				}

				_, err = stmt.Exec(v, k)

				//excute error
				if err != nil {
					panic(err)
					fmt.Fprintf(w, "The Following error Occured:", err)
				}
				fmt.Print("successfully Extracted And Inserted", k, v, "\n")

				//breakdb.Close()
				//defer db.Close()
			}

			//fmt.Print(row)
			//break
			//fmt.Fprintf(w, row)

			//}

			qry := ("SELECT uid,`name` FROM `dump`")
			qry = strings.Replace(qry, "`", "\"", -1)
			fmt.Println(qry)
			rowss, err := db.Query(qry)
			if err != nil {
				panic(err)
			}
			defer rowss.Close()

			bks := make([]maplistbox, 0)
			for rowss.Next() {
				bk := maplistbox{}
				err := rowss.Scan(&bk.Id, &bk.Value) // order matters
				if err != nil {
					panic(err)
				}
				bks = append(bks, bk)
			}

			//t, err := template.ParseFiles("map_vn.html")
			//if err != nil {
			//fmt.Print(err)
			//}
			fmt.Print("successfully Selected Inserted Rows For Display")

			err = tpl.ExecuteTemplate(w, "settings.html", bks)
			//err = tpl.ExecuteTemplate(w, "settings.html", bks)
			if err != nil {
				fmt.Print(err)
			}
			fmt.Print("successfully Loaded Map_vn.html Template To Display Selected Rows In a Drop Down")

			///empty the dubm table_name

			qry = ("TRUNCATE `dump`")
			qry = strings.Replace(qry, "`", "\"", -1)
			fmt.Println(qry)
			_, err = db.Query(qry)
			if err != nil {
				panic(err)
			}
			//defer rowss.Close()

			/// end rmpty dubmp table_name

			fmt.Print(bks)
		} else {
			errr := `<table width="50%" border="0" align="center">
	<tr>
	<td align="center">The Following Error Occured: The File Type You Selected Is Not Supported. Pls Upload Only XlSl File</td>
	</tr>
	</table>`
			fmt.Fprintf(w, errr)
		}

	} else {
		errr := `<table width="50%" border="0" align="center">
	<tr>
	<td align="center">The Following Error Occured: The File Type You Selected Is Not Supported. Pls Upload Only XlSl File</td>
	</tr>
	</table>`
		fmt.Fprintf(w, errr)
	}

	/// excel extracting end here
}

func settings(w http.ResponseWriter, r *http.Request) {
	var err error

	/*psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err = sql.Open("postgres", psqlInfo)*/

	fmt.Println("......................successfully Entered Setting Function......................................................................")
	//fmt.Println("Getting And Loading the excel file sheet: ", header.Filename)
	xlsx, err := excelize.OpenFile(fpath)
	if err != nil {
		fmt.Println("Excel Sheet failed to Load Due To: ", err)
		os.Exit(1)
	}

	//fmt.Println("successfully Loaded the excel file sheet: ", header.Filename)
	sheet_index := r.FormValue("sht")
	fmt.Println("successfully Gotten sheet Index using Post Method: ", sheet_index)

	sheet_name = "sheet" + sheet_index
	fmt.Println("successfully Gotten sheet To Work On: ", sheet_name)

	// Get all the rows in a sheet.
	//	rows := xlsx.GetRows("sheet" + sheet_index)
	rows := xlsx.GetRows(sheet_name)

	if len(rows) != 0 {

		fmt.Println("successfully Extracted the Whole excel Rows: ", rows)
		// Gets the heading title
		//for _, row := range rows {
		// Form submitted
		col_indes := r.FormValue("col_index")

		var col_index int
		if _, err := fmt.Sscanf(col_indes, "%5d", &col_index); err == nil {
			fmt.Println(col_index) //
		}
		fmt.Println("successfully Gotten the Index for the Row Header: ", col_index) /////////////////////////////////

		count := 0
		for _, ro := range rows[col_index] {

			qy := "INSERT INTO dump (name,uid)VALUES($1,$2)"
			//qy = strings.Replace(qy, "`", "\"", -1)
			//stmt, err := db.Prepare("INSERT `dtest` SET dname=$,damt=$")

			stmt, err := db.Prepare(qy)
			// prepare statement error
			if err != nil {
				panic(err)
				fmt.Fprintf(w, "The Following error Occured:", err)
			}

			_, err = stmt.Exec(ro, count)

			//excute error
			if err != nil {
				panic(err)
				fmt.Fprintf(w, "The Following error Occured:", err)
			}
			fmt.Print("successfully Extracted And Inserted", count, ro, "\n")

			//break
			count = count + 1
			//defer db.Close()
		}

		//fmt.Print(row)
		//break
		//fmt.Fprintf(w, row)

		//}

		qry := ("SELECT uid,`name` FROM `dump`")
		qry = strings.Replace(qry, "`", "\"", -1)
		fmt.Println(qry)
		rowss, err := db.Query(qry)
		if err != nil {
			panic(err)
		}
		defer rowss.Close()

		bks := make([]maplistbox, 0)
		for rowss.Next() {
			bk := maplistbox{}
			err := rowss.Scan(&bk.Id, &bk.Value) // order matters
			if err != nil {
				panic(err)
			}
			bks = append(bks, bk)
		}

		//t, err := template.ParseFiles("map_vn.html")
		//if err != nil {
		//fmt.Print(err)
		//}
		fmt.Print("successfully Selected Inserted Rows For Display")

		err = tpl.ExecuteTemplate(w, "map_vn.html", bks)
		//err = tpl.ExecuteTemplate(w, "settings.html", bks)
		if err != nil {
			fmt.Print(err)
		}
		fmt.Print("successfully Loaded Map_vn.html Template To Display Selected Rows In a Drop Down")

		///empty the dubm table_name

		qry = ("TRUNCATE `dump`")
		qry = strings.Replace(qry, "`", "\"", -1)
		fmt.Println(qry)
		_, err = db.Query(qry)
		if err != nil {
			panic(err)
		}
		//defer rowss.Close()

		/// end rmpty dubmp table_name

		fmt.Print(bks)
	} else {

		err = tpl.ExecuteTemplate(w, "goupload.html", nil)

		errr := `<table width="50%" border="0" align="center">
	<tr>
	<td align="center">The Following Error Occured: The Sheet You Selected Is Empty Or It Is Corrupt</td>
	</tr>
	</table>`
		fmt.Fprintf(w, errr)

	}

}

func process(w http.ResponseWriter, r *http.Request) {
	//	var err error
	/*psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err = sql.Open("postgres", psqlInfo)

	//db, err = sql.Open("postgres", "postgres://postgres:1234567@172.17.0.3/staff_db?sslmode=disable")
	db, err = sql.Open("postgres", "postgres://postgres:@db/staff_db?sslmode=disable")
	if err != nil {
		//panic(err)
		fmt.Println("Database connection failed due to", err)
	}

	err = db.Ping()
	if err != nil {
		fmt.Println("Database connection failed due to", err)

	}
	fmt.Println("Successfully connected!")*/

	//if r.Method == "POST" {
	// Form submitted

	v := r.FormValue("vn")
	v = strings.TrimSpace(v)
	mmonth := r.FormValue("month")
	yyear := r.FormValue("year")

	sstart := r.FormValue("start")
	sstart = strings.TrimSpace(sstart)

	eend := r.FormValue("end")
	eend = strings.TrimSpace(eend)

	// Output form data
	fmt.Println("Value", v, " ", mmonth, " ", yyear, " ", sstart, " ", eend)

	if mmonth != "" && yyear != "" {

		// Convert the string to int
		var index int
		if _, err := fmt.Sscanf(v, "%5d", &index); err == nil {
			fmt.Println(index) //
		}

		var start int
		if _, err := fmt.Sscanf(sstart, "%5d", &start); err == nil {
			fmt.Println(start) //
		}

		var end int
		if _, err := fmt.Sscanf(eend, "%5d", &end); err == nil {
			fmt.Println(end) //
		}
		// Output Converted string to int
		fmt.Println("Value", index)

		/// excel extracting starts here

		/// get and out the file been held in a public variable
		fmt.Println(fpath)

		// open the excel and start data extraction
		//xlsx, err := excelize.OpenFile("C:/Users/Nd/Documents/GoWorkSpace/src/excel/Workbook.xlsx")
		xlsx, err := excelize.OpenFile(fpath)
		if err != nil {
			fmt.Println("Eroro Opening Excel file", err)
			tpl.ExecuteTemplate(w, "goupload.html", nil)
			//os.Exit(1)
		}

		// Get all the rows in a sheet.
		//rows := xlsx.GetRows("Sheet1")
		rows := xlsx.GetRows(sheet_name)
		// Output the rows
		fmt.Println(rows)

		// loop through the excel sheet passing the colunm index to get only the colunm index

		goback := `<table width="50%" border="0" align="center">
	<tr>
	<td align="center"><a href="https://d04f6004.ngrok.io/">Go Back</a> </td>
	</tr>
	</table>  </br>`
		fmt.Fprintf(w, goback)

		/// get vn and deduction
		//mmonth := "December"
		//yyear := 2014
		count := 0
		//start = 5
		//end = 10
		//goback := `<a href="http://d4a86761.ngrok.io/">Go Back</a> </br></br>`
		//fmt.Fprintf(w, goback)

		if sstart != "" && eend != "" {

			for _, row := range rows {
				count = count + 1
				//get the mapped vn colunm
				if count >= start && count <= end {

					vn_col := row[index]
					vn_col = strings.TrimSpace(vn_col)
					vn_col = strings.Replace(vn_col, "'", "", -1)
					// Output the mapped vn
					fmt.Print("The Gotten Vn Colunm Is: ", vn_col, " The Count Value: ", count, "\n")
					//fmt.Fprintf(w, vn_col, "\n")

					///////////////////////////start///////////////////////////////////

					// get deduction

					qry := fmt.Sprintf("SELECT employee_info_employee_no,`DEDUCTIONS`,`MONTH`,`YEAR`  FROM excercise_data WHERE `MONTH` = '%v' and `YEAR` = '%v' and employee_info_employee_no = '%v' LIMIT 1", mmonth, yyear, vn_col)
					qry = strings.Replace(qry, "`", "\"", -1)
					fmt.Println(qry)
					var vvn string
					//rows, err := db.Query(qry)
					//if err != nil {
					//panic(err)
					//}
					//defer rows.Close()

					if err := db.QueryRow(qry).Scan(&vvn); err == nil { // row is present

						fmt.Println("Row Not Present For VN: ", vn_col)
						fmt.Fprintf(w, "%v\n<br>", vn_col, "Row Not Present For VN: ")
						// 1 row

						/// row is present end here

					} else if err == sql.ErrNoRows {
						// empty result
						fmt.Println("Row Is Empty For VN: ", vn_col)
						//estr := ""
						emt := fmt.Sprintf("No Matching Record Found For VN: %v For Month: %v And Year : %v ", vn_col, mmonth, yyear)
						fmt.Fprintf(w, "%v\n</br>", emt)

						//fmt.Fprintf(w, "%s\n", vn_col, "Row Is Empty For For VN: ")

					} else { /// row is error
						// error
						fmt.Println("Record Found And Inserted For VN: ", vn_col)
						fnt := fmt.Sprintf("Record Found And Inserted For VN: %v ", vn_col)
						fmt.Fprintf(w, "%v\n</br>", fnt)

						rows, err := db.Query(qry)
						if err != nil {
							//panic(err)
							fmt.Println("Query failed Due To: ", err)
						}
						defer rows.Close()

						bks := make([]Empinfo, 0)
						for rows.Next() {
							bk := Empinfo{}
							err := rows.Scan(&bk.employee_info, &bk.deductions, &bk.month, &bk.year) // order matters
							if err != nil {
								//panic(err)
								fmt.Println("Row Scan Failed Due To: ", err)
							}
							bks = append(bks, bk)
						}
						if err = rows.Err(); err != nil {
							//panic(err)
							fmt.Println("Eroro eroro Ocuured Dut To: ", err)
						}
						var pk string
						var vn string
						var month string
						var year int

						for _, bk := range bks {
							// fmt.Println(bk.isbn, bk.title, bk.author, bk.price)
							pk = bk.deductions
							vn = bk.employee_info
							month = bk.month
							year = bk.year

							fmt.Println(pk)
							fmt.Println(vn)
							//fmt.Printf("%s, %s, %s, %s\n", bk.employee_info, bk.first_name, bk.middle_name, bk.surname)
						}
						//tk := strings.Split(pk, "\n")
						//fmt.Println(tk)
						tk := strings.Replace(pk, ",", "\n", -1)
						tk = strings.Replace(tk, ",", "", -1)
						tk = strings.Replace(tk, "=", " ", -1)
						//tk = strings.Replace(tk, ",", "\n", -1)
						fmt.Println(tk)
						tk = strings.Replace(tk, "CREDIT  DIRECT", "CREDIT_DIRECT", -1)
						tk = strings.Replace(tk, "AMAECOM GLOBAL", "AMAECOM_GLOBAL", -1)
						tk = strings.Replace(tk, "GIFORTUNE ENTER", "GIFORTUNE_ENTER", -1)
						//tk = strings.Replace(tk, "NULGE (MORNACH MICRO FIN.BANK)", "NULGE_(MORNACH_MICRO_FIN_BANK)", -1)
						fmt.Println(tk)
						//tk = strings.Replace(tk, "N", "", -1)
						pk2 := strings.Fields(tk)
						fmt.Println(pk2)
						pk3 := len(pk2)
						fw := pk2[0]
						sw := pk2[1]
						fmt.Println(pk2)
						fmt.Println(pk3)
						fmt.Println(fw)
						fmt.Println(sw)

						f := 0
						s := 1
						for i := 0; i <= len(pk2); i++ {
							fmt.Println(len(pk2))
							if f < i {
								fw2 := pk2[f]
								sw2 := pk2[s]

								f = f + 2
								s = s + 2

								///////////

								/// delete matching and duplicate record

								qryd := fmt.Sprintf("DELETE FROM staff_deduction WHERE vn ='%v' AND deduction_name ='%v' AND `month` ='%v'  AND `year` ='%v' ", vn, fw2, month, year)
								qryd = strings.Replace(qryd, "`", "\"", -1)
								fmt.Println(qryd)
								_, err = db.Query(qryd)
								if err != nil {
									//panic("Delection From Staff Deduction Failed Due To: ", err)
									fmt.Println("Delection From Staff Deduction Failed Due To: ", err)
								}

								/// deleting ends

								// insert

								qy := "INSERT INTO staff_deduction (vn, deduction_name,deduction_amount,month,year)VALUES($1,$2,$3,$4,$5)"
								//qy := "INSERT INTO staff_deduction (vn, deduction_name,deduction_amount,month,year)VALUES($1,$2,$3,$4,$5)ON CONFLICT (vn, deduction_name, month, year) DO UPDATE SET vn = $1,deduction_name = $2,deduction_amount=$3,month=$4,year=$5"
								//qy = strings.Replace(qy, "`", "\"", -1)
								//stmt, err := db.Prepare("INSERT `dtest` SET dname=$,damt=$")

								stmt, err := db.Prepare(qy)

								if err != nil {
									//panic("Insert into Staff Deduction Failed Due To: ", err)
									fmt.Println("Insert into Staff Deduction Failed Due To: ", err)
								}

								//checkErr(err)

								sw2 = strings.Replace(sw2, "N", "", -1)

								//res, err := stmt.Exec(fw2, sw2)
								_, err = stmt.Exec(vn, fw2, sw2, month, year)
								//checkErr(err)
								if err != nil {

									//panic("Excuate Insert Into  Staff Deduction Failed Due To: ", err)
									fmt.Println("Excuate Insert Into  Staff Deduction Failed Due To: ", err)
								}

								fmt.Println(fw2)
								fmt.Println(sw2)

							}
						}

						/// row is error ends here
					}
				} // end count start and end loop
				//defer db.Close()
			} /// rows loop ends here

			// display the results

			qry := fmt.Sprintf("SELECT vn,deduction_name,deduction_amount,`month`,`year` FROM `staff_deduction` WHERE `month` ='%v'  AND `year`='%v' ", mmonth, yyear)
			qry = strings.Replace(qry, "`", "\"", -1)
			fmt.Println(qry)
			rowsr, err := db.Query(qry)
			if err != nil {
				// handle this error better than this
				//panic("Select From Staff Deduction Failed Due To: ", err)
				fmt.Println("Select From Staff Deduction Failed Due To: ", err)
			}
			defer rowsr.Close()

			fmt.Fprintf(w, "........SUMMARY OF THE INSERTED DATA......................\n</br>")

			for rowsr.Next() {
				var verification_numer string
				var deduction_name string
				var deduction_amount string
				var month string
				var year string
				err = rowsr.Scan(&verification_numer, &deduction_name, &deduction_amount, &month, &year)
				if err != nil {
					// handle this error
					//panic("Row Scan From Staff Deduction Failed Due To: ", err)
					fmt.Println("Row Scan From Staff Deduction Failed Due To: ", err)
				}

				fmt.Fprintf(w, "%s %s %s %s %s\n</br>", verification_numer, deduction_name, deduction_amount, month, year)
				//fmt.Printf(verification_numer, deduction_name, deduction_amount, month, year)
				//fmt.Println(verification_numer, deduction_name, deduction_amount, month, year)

			}
			//defer db.Close()
			/// display result ends

			/// main function ends here
		} else if sstart != "" && eend == "" { /// end start and end loop
			fmt.Fprintf(w, "End Is Empty")

			for _, row := range rows {
				count = count + 1
				//get the mapped vn colunm
				if count >= start {
					vn_col := row[index]
					vn_col = strings.TrimSpace(vn_col)
					vn_col = strings.Replace(vn_col, "'", "", -1)
					// Output the mapped vn
					fmt.Print("The Gotten Vn Colunm Is: ", vn_col, " The Count Value: ", count, "\n")
					//fmt.Fprintf(w, vn_col, "\n")

					///////////////////////////start///////////////////////////////////

					// get deduction

					qry := fmt.Sprintf("SELECT employee_info_employee_no,`DEDUCTIONS`,`MONTH`,`YEAR`  FROM excercise_data WHERE `MONTH` = '%v' and `YEAR` = '%v' and employee_info_employee_no = '%v' LIMIT 1", mmonth, yyear, vn_col)
					qry = strings.Replace(qry, "`", "\"", -1)
					fmt.Println(qry)
					var vvn string
					//rows, err := db.Query(qry)
					//if err != nil {
					//panic(err)
					//}
					//defer rows.Close()

					if err := db.QueryRow(qry).Scan(&vvn); err == nil { // row is present

						fmt.Println("Row Not Present For VN: ", vn_col)
						fmt.Fprintf(w, "\n", vn_col, "Row Not Present For VN: ")
						// 1 row

						/// row is present end here

					} else if err == sql.ErrNoRows {
						// empty result
						fmt.Println("Row Is Empty For VN: ", vn_col)
						//estr := ""
						//emt := fmt.Sprintf("No Matching Record Found For VN: %v ", vn_col)
						emt := fmt.Sprintf("No Matching Record Found For VN: %v For Month: %v And Year : %v ", vn_col, mmonth, yyear)

						fmt.Fprintf(w, "%v\n</br>", emt)

						//fmt.Fprintf(w, "%s\n", vn_col, "Row Is Empty For For VN: ")

					} else { /// row is error
						// error
						fmt.Println("Record Found And Inserted For VN: ", vn_col)
						fnt := fmt.Sprintf("Record Found And Inserted For VN: %v ", vn_col)
						fmt.Fprintf(w, "%v\n</br>", fnt)

						rows, err := db.Query(qry)
						if err != nil {
							//panic(err)
							fmt.Println("Query failed Due To: ", err)
						}
						defer rows.Close()

						bks := make([]Empinfo, 0)
						for rows.Next() {
							bk := Empinfo{}
							err := rows.Scan(&bk.employee_info, &bk.deductions, &bk.month, &bk.year) // order matters
							if err != nil {
								//panic(err)
								fmt.Println("Row Scan Failed Due To: ", err)
							}
							bks = append(bks, bk)
						}
						if err = rows.Err(); err != nil {
							//panic(err)
							fmt.Println("Eroro eroro Ocuured Dut To: ", err)
						}
						var pk string
						var vn string
						var month string
						var year int

						for _, bk := range bks {
							// fmt.Println(bk.isbn, bk.title, bk.author, bk.price)
							pk = bk.deductions
							vn = bk.employee_info
							month = bk.month
							year = bk.year

							fmt.Println(pk)
							fmt.Println(vn)
							//fmt.Printf("%s, %s, %s, %s\n", bk.employee_info, bk.first_name, bk.middle_name, bk.surname)
						}
						//tk := strings.Split(pk, "\n")
						//fmt.Println(tk)
						tk := strings.Replace(pk, ",", "\n", -1)
						tk = strings.Replace(tk, ",", "", -1)
						tk = strings.Replace(tk, "=", " ", -1)
						//tk = strings.Replace(tk, ",", "\n", -1)
						fmt.Println(tk)
						tk = strings.Replace(tk, "CREDIT  DIRECT", "CREDIT_DIRECT", -1)
						tk = strings.Replace(tk, "AMAECOM GLOBAL", "AMAECOM_GLOBAL", -1)
						tk = strings.Replace(tk, "GIFORTUNE ENTER", "GIFORTUNE_ENTER", -1)

						//tk = strings.Replace(tk, "NULGE (MORNACH MICRO FIN.BANK)", "NULGE_(MORNACH_MICRO_FIN_BANK)", -1)
						fmt.Println(tk)
						//tk = strings.Replace(tk, "N", "", -1)
						pk2 := strings.Fields(tk)
						fmt.Println(pk2)
						pk3 := len(pk2)
						fw := pk2[0]
						sw := pk2[1]
						fmt.Println(pk2)
						fmt.Println(pk3)
						fmt.Println(fw)
						fmt.Println(sw)

						f := 0
						s := 1
						for i := 0; i <= len(pk2); i++ {
							fmt.Println(len(pk2))
							if f < i {
								fw2 := pk2[f]
								sw2 := pk2[s]

								f = f + 2
								s = s + 2

								///////////

								/// delete matching and duplicate record

								qryd := fmt.Sprintf("DELETE FROM staff_deduction WHERE vn ='%v' AND deduction_name ='%v' AND `month` ='%v'  AND `year` ='%v' ", vn, fw2, month, year)
								qryd = strings.Replace(qryd, "`", "\"", -1)
								fmt.Println(qryd)
								_, err = db.Query(qryd)
								if err != nil {
									//panic("Delection From Staff Deduction Failed Due To: ", err)
									fmt.Println("Delection From Staff Deduction Failed Due To: ", err)
								}

								/// deleting ends

								// insert

								qy := "INSERT INTO staff_deduction (vn, deduction_name,deduction_amount,month,year)VALUES($1,$2,$3,$4,$5)"
								//qy := "INSERT INTO staff_deduction (vn, deduction_name,deduction_amount,month,year)VALUES($1,$2,$3,$4,$5)ON CONFLICT (vn, deduction_name, month, year) DO UPDATE SET vn = $1,deduction_name = $2,deduction_amount=$3,month=$4,year=$5"
								//qy = strings.Replace(qy, "`", "\"", -1)
								//stmt, err := db.Prepare("INSERT `dtest` SET dname=$,damt=$")

								stmt, err := db.Prepare(qy)

								if err != nil {
									//panic("Insert into Staff Deduction Failed Due To: ", err)
									fmt.Println("Insert into Staff Deduction Failed Due To: ", err)
								}

								//checkErr(err)

								sw2 = strings.Replace(sw2, "N", "", -1)

								//res, err := stmt.Exec(fw2, sw2)
								_, err = stmt.Exec(vn, fw2, sw2, month, year)
								//checkErr(err)
								if err != nil {

									//panic("Excuate Insert Into  Staff Deduction Failed Due To: ", err)
									fmt.Println("Excuate Insert Into  Staff Deduction Failed Due To: ", err)
								}

								fmt.Println(fw2)
								fmt.Println(sw2)
							}
						}

						/// row is error ends here
					}
				} // end count start and end loop
				//	db.Close()
			} /// rows loop ends here

			// display the results

			qry := fmt.Sprintf("SELECT vn,deduction_name,deduction_amount,`month`,`year` FROM `staff_deduction` WHERE `month` ='%v'  AND `year`='%v' ", mmonth, yyear)
			qry = strings.Replace(qry, "`", "\"", -1)
			fmt.Println(qry)
			rowsr, err := db.Query(qry)
			if err != nil {
				// handle this error better than this
				//panic("Select From Staff Deduction Failed Due To: ", err)
				fmt.Println("Select From Staff Deduction Failed Due To: ", err)
			}
			//defer rowsr.Close()

			fmt.Fprintf(w, "........SUMMARY OF THE INSERTED DATA......................\n</br>")

			for rowsr.Next() {
				var verification_numer string
				var deduction_name string
				var deduction_amount string
				var month string
				var year string
				err = rowsr.Scan(&verification_numer, &deduction_name, &deduction_amount, &month, &year)
				if err != nil {
					// handle this error
					//panic("Row Scan From Staff Deduction Failed Due To: ", err)
					fmt.Println("Row Scan From Staff Deduction Failed Due To: ", err)
				}

				fmt.Fprintf(w, "%s %s %s %s %s\n</br>", verification_numer, deduction_name, deduction_amount, month, year)
				//fmt.Printf(verification_numer, deduction_name, deduction_amount, month, year)
				//fmt.Println(verification_numer, deduction_name, deduction_amount, month, year)

			}

		} else {
			fmt.Println("Both start End end are empty")

			for _, row := range rows {
				//count = count + 1
				//get the mapped vn colunm
				//if count >= start && count <= end {
				vn_col := row[index]
				vn_col = strings.TrimSpace(vn_col)
				vn_col = strings.Replace(vn_col, "'", "", -1)
				// Output the mapped vn
				fmt.Print("The Gotten Vn Colunm Is: ", vn_col, " The Count Value: ", count, "\n")
				//fmt.Fprintf(w, vn_col, "\n")

				///////////////////////////start///////////////////////////////////

				// get deduction

				qry := fmt.Sprintf("SELECT employee_info_employee_no,`DEDUCTIONS`,`MONTH`,`YEAR`  FROM excercise_data WHERE `MONTH` = '%v' and `YEAR` = '%v' and employee_info_employee_no = '%v' LIMIT 1", mmonth, yyear, vn_col)
				qry = strings.Replace(qry, "`", "\"", -1)
				fmt.Println(qry)
				var vvn string
				//rows, err := db.Query(qry)
				//if err != nil {
				//panic(err)
				//}
				//defer rows.Close()

				if err := db.QueryRow(qry).Scan(&vvn); err == nil { // row is present

					fmt.Println("Row Not Present For VN: ", vn_col)
					fmt.Fprintf(w, "\n", vn_col, "Row Not Present For VN: ")
					// 1 row

					/// row is present end here

				} else if err == sql.ErrNoRows {
					// empty result
					fmt.Println("Row Is Empty For VN: ", vn_col)
					//estr := ""
					//emt := fmt.Sprintf("No Matching Record Found For VN1: %v ", vn_col)
					emt := fmt.Sprintf("No Matching Record Found For VN: %v For Month: %v And Year : %v ", vn_col, mmonth, yyear)
					fmt.Fprintf(w, "%v\n</br>", emt)

					//fmt.Fprintf(w, "%s\n", vn_col, "Row Is Empty For For VN: ")

				} else { /// row is error
					// error
					fmt.Println("Record Found And Inserted For VN: ", vn_col)
					fnt := fmt.Sprintf("Record Found And Inserted For VN: %v ", vn_col)
					fmt.Fprintf(w, "%v\n</br>", fnt)

					rows, err := db.Query(qry)
					if err != nil {
						//panic(err)
						fmt.Println("Query failed Due To: ", err)
					}
					defer rows.Close()

					bks := make([]Empinfo, 0)
					for rows.Next() {
						bk := Empinfo{}
						err := rows.Scan(&bk.employee_info, &bk.deductions, &bk.month, &bk.year) // order matters
						if err != nil {
							//panic(err)
							fmt.Println("Row Scan Failed Due To: ", err)
						}
						bks = append(bks, bk)
					}
					if err = rows.Err(); err != nil {
						//panic(err)
						fmt.Println("Eroro eroro Ocuured Dut To: ", err)
					}
					var pk string
					var vn string
					var month string
					var year int

					for _, bk := range bks {
						// fmt.Println(bk.isbn, bk.title, bk.author, bk.price)
						pk = bk.deductions
						vn = bk.employee_info
						month = bk.month
						year = bk.year

						fmt.Println(pk)
						fmt.Println(vn)
						//fmt.Printf("%s, %s, %s, %s\n", bk.employee_info, bk.first_name, bk.middle_name, bk.surname)
					}
					//tk := strings.Split(pk, "\n")
					//fmt.Println(tk)
					tk := strings.Replace(pk, ",", "\n", -1)
					tk = strings.Replace(tk, ",", "", -1)
					tk = strings.Replace(tk, "=", " ", -1)
					//tk = strings.Replace(tk, ",", "\n", -1)
					fmt.Println(tk)
					tk = strings.Replace(tk, "CREDIT  DIRECT", "CREDIT_DIRECT", -1)
					tk = strings.Replace(tk, "AMAECOM GLOBAL", "AMAECOM_GLOBAL", -1)
					tk = strings.Replace(tk, "GIFORTUNE ENTER", "GIFORTUNE_ENTER", -1)

					//tk = strings.Replace(tk, "NULGE (MORNACH MICRO FIN.BANK)", "NULGE_(MORNACH_MICRO_FIN_BANK)", -1)
					fmt.Println(tk)
					//tk = strings.Replace(tk, "N", "", -1)
					pk2 := strings.Fields(tk)
					fmt.Println(pk2)
					pk3 := len(pk2)
					fw := pk2[0]
					sw := pk2[1]
					fmt.Println(pk2)
					fmt.Println(pk3)
					fmt.Println(fw)
					fmt.Println(sw)

					f := 0
					s := 1
					for i := 0; i <= len(pk2); i++ {
						fmt.Println(len(pk2))
						if f < i {
							fw2 := pk2[f]
							sw2 := pk2[s]

							f = f + 2
							s = s + 2

							///////////

							/// delete matching and duplicate record

							qryd := fmt.Sprintf("DELETE FROM staff_deduction WHERE vn ='%v' AND deduction_name ='%v' AND `month` ='%v'  AND `year` ='%v' ", vn, fw2, month, year)
							qryd = strings.Replace(qryd, "`", "\"", -1)
							fmt.Println(qryd)
							_, err = db.Query(qryd)
							if err != nil {
								//panic("Delection From Staff Deduction Failed Due To: ", err)
								fmt.Println("Delection From Staff Deduction Failed Due To: ", err)
							}

							/// deleting ends

							// insert

							qy := "INSERT INTO staff_deduction (vn, deduction_name,deduction_amount,month,year)VALUES($1,$2,$3,$4,$5)"
							//qy := "INSERT INTO staff_deduction (vn, deduction_name,deduction_amount,month,year)VALUES($1,$2,$3,$4,$5)ON CONFLICT (vn, deduction_name, month, year) DO UPDATE SET vn = $1,deduction_name = $2,deduction_amount=$3,month=$4,year=$5"
							//qy = strings.Replace(qy, "`", "\"", -1)
							//stmt, err := db.Prepare("INSERT `dtest` SET dname=$,damt=$")

							stmt, err := db.Prepare(qy)

							if err != nil {
								//panic("Insert into Staff Deduction Failed Due To: ", err)
								fmt.Println("Insert into Staff Deduction Failed Due To: ", err)
							}

							//checkErr(err)

							sw2 = strings.Replace(sw2, "N", "", -1)

							//res, err := stmt.Exec(fw2, sw2)
							_, err = stmt.Exec(vn, fw2, sw2, month, year)
							//checkErr(err)
							if err != nil {

								//panic("Excuate Insert Into  Staff Deduction Failed Due To: ", err)
								fmt.Println("Excuate Insert Into  Staff Deduction Failed Due To: ", err)
							}

							fmt.Println(fw2)
							fmt.Println(sw2)
						}
					}

					/// row is error ends here
				}
				//} // end count start and end loop
				//db.Close()
			} /// rows loop ends here

			// display the results

			qry := fmt.Sprintf("SELECT vn,deduction_name,deduction_amount,`month`,`year` FROM `staff_deduction` WHERE `month` ='%v'  AND `year`='%v' ", mmonth, yyear)
			//qry := fmt.Sprintf("SELECT vn,deduction_name,deduction_amount,`month`,`year` FROM `staff_deduction` WHERE `month` ='%v'  AND `year`='%v' ", month, year)

			qry = strings.Replace(qry, "`", "\"", -1)
			fmt.Println(qry)
			rowsr, err := db.Query(qry)
			if err != nil {
				// handle this error better than this
				//panic("Select From Staff Deduction Failed Due To: ", err)
				fmt.Println("Select From Staff Deduction Failed Due To: ", err)
				fmt.Fprintf(w, "Select From Staff Deduction Failed Due To: Invalid Year And Month ")

			}
			defer rowsr.Close()

			fmt.Fprintf(w, "........SUMMARY OF THE INSERTED DATA......................\n\n\n</br>")

			for rowsr.Next() {
				var verification_numer string
				var deduction_name string
				var deduction_amount string
				var month string
				var year string
				err = rowsr.Scan(&verification_numer, &deduction_name, &deduction_amount, &month, &year)
				if err != nil {
					// handle this error
					//panic("Row Scan From Staff Deduction Failed Due To: ", err)
					fmt.Println("Row Scan From Staff Deduction Failed Due To: ", err)
				}

				fmt.Fprintf(w, "%s %s %s %s %s\n</br>", verification_numer, deduction_name, deduction_amount, month, year)
				//fmt.Printf(verification_numer, deduction_name, deduction_amount, month, year)
				//fmt.Println(verification_numer, deduction_name, deduction_amount, month, year)

			}

		} /// end start and end loop
		/////////////////////////end//////////////////////////////////////
	} else {

		tpl.ExecuteTemplate(w, "gouplaod.html", nil)
		goback := `<table width="50%" border="0" align="center">
	<tr>
	<td align="center">Month And Year Missing In your Selection. Kindly re-Upload Again</a> </td>
	</tr>
	</table>  </br>`
		fmt.Fprintf(w, goback)
	} ///check for month and year ends here
}

func init() {

	var version = "APP VERSION: 0.001 DEVEL"

	fmt.Println("You Are Currently Running: ", version)

	//tpl = template.Must(template.ParseFiles("goupload.html"))

	//tpl = template.Must(template.ParseGlob("template/*.html"))
	//tpl = template.Must(template.ParseGlob("templates/*.html"))
	fmt.Println("Checking The Host Operating System:")

	if runtime.GOOS == "windows" {
		tpl = template.Must(template.ParseGlob("templates/*.html"))
		fmt.Println("The Host System is Windows: Windows Path Loaded")
	}

	if runtime.GOOS == "linux" {
		tpl = template.Must(template.ParseGlob("/templates/*.html"))
		fmt.Println("The Host system is Linux: Linux Path Loaded ")
	}
	//defer db.Close()
	var err error
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err = sql.Open("postgres", psqlInfo)

	//db, err = sql.Open("postgres", "postgres://postgres:1234567@172.17.0.3/staff_db?sslmode=disable")
	//db, err = sql.Open("postgres", "postgres://postgres:@db/staff_db?sslmode=disable")

	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		fmt.Println("Database connection failed due to", err)
		return
	}
	fmt.Println("Successfully connected To DB!")

}
func main() {

	http.HandleFunc("/", uploadHandler)
	http.HandleFunc("/process", process)
	http.HandleFunc("/settings", settings)
	http.ListenAndServe(":5050", nil)

}

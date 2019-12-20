package vnu

import (
	"fmt"
	"time"
	"strings"
	"github.com/PuerkitoBio/goquery"
	"strconv"
)

// CreateAcc Generate acc vnu
func CreateAcc(code string, start int, end int)([]string){
	// end: max is 5
	// start: 0
	var result []string
	year := time.Now().Year()
	month := time.Now().Month()
	if (month < 10) {
		year--
	}

	for i := start; i < end; i++ {
		head := year - i
		for j := 0; j < 2000; j++ {
			headStr := fmt.Sprintf("%v%s", head, code)[2:]
			acc := getMssv(j, headStr)
			result = append(result, acc)
		}
	}
	return result
}

func getMssv(i int, head string)(string) {
	if  (i < 10) {
		return fmt.Sprintf("%s000%v", head, i)
	} else if (i < 100) {
		return fmt.Sprintf("%s00%v", head, i)
	} else if (i < 1000) {
		return fmt.Sprintf("%s0%v", head, i)
	} else {
		return fmt.Sprintf("%s%v", head, i)
	}
}

// GetRowIndexFromTable get rowindex from table
func GetRowIndexFromTable(table string, creditname string)(string, bool){
	html := fmt.Sprintf(`<html>
    <body>
        <table>
            <tbody>
                %s
            </tbody>
        </table>
    </body>
</html>`, table)
	return getRowIndexFromStr(html, creditname)
}


func getRowIndexFromStr(html string, creditname string)(string, bool){
	buffer := strings.NewReader(html)
	return getRowIndexFromReader(buffer, creditname)
}

func getRowIndexFromReader(buffer *strings.Reader, creditname string) (string, bool) {
	doc, err := goquery.NewDocumentFromReader(buffer)
	if err != nil {
		return "", false
	}
	return getRowIndexFromDoc(doc, creditname)
}

func getRowIndexFromDoc(doc *goquery.Document, creditname string) (string, bool) {
	var rowIndex string
	var exists bool
	rowIndex = ""
	exists = false
	doc.Find("tr").Each(func(i int, s *goquery.Selection) {
		
		var credit string
		creditQuery := s.Find("td:nth-child(5)")
		
		if (!strings.Contains(creditQuery.Text(), "(")) {
			credit  = creditQuery.Text()
		} else {
			creditQuery.Children().Remove()
			credit =  creditQuery.Text()
		}
		credit = strings.Replace(credit, "(", "", -1)
		credit = strings.Replace(credit, ")", "", -1)
		credit = strings.TrimSpace(credit)
		temp := strings.Split(credit, " ")
		if (len(temp) == 3) {
			credit = fmt.Sprintf("%s%s %s", temp[0], temp[1], temp[2])
		} else {
			// if (len(temp[1]) <= 4) {
			// 	credit = fmt.Sprintf("%s%s", temp[0], temp[1])
			// }
			if _, err := strconv.ParseInt(temp[1],10,64); err == nil {
				credit = fmt.Sprintf("%s%s", temp[0], temp[1])
			}
		}
		if (strings.EqualFold(credit, creditname)) {
			input := s.Find("td input[type='checkbox']")
			rowIndex, exists = input.Attr("data-rowindex")
			return
		}
	})
	return rowIndex, exists
}

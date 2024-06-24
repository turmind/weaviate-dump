package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"

	"github.com/tidwall/gjson"
	excelize "github.com/xuri/excelize/v2"
)

func main() {
	// parse flags host, port, and token
	host := flag.String("host", "localhost", "weaviate host")
	port := flag.String("port", "8080", "weaviate port")
	token := flag.String("token", "", "weaviate token")
	class := flag.String("class", "", "dump class name")
	limit := flag.Int("limit", 25, "limit")
	flag.Parse()

	// check class name
	if *class == "" {
		fmt.Println("class name is required")
		return
	}

	// create excel file
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	// http request
	client := &http.Client{}
	after := ""
	row := 0
	colMap := make(map[string]string)
	for {
		req, err := http.NewRequest("GET", fmt.Sprintf("http://%s:%s/v1/objects?include=vector&limit=%d&after=%s&class=%s", *host, *port, *limit, after, *class), nil)
		if err != nil {
			fmt.Println(err)
			return
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")
		if *token != "" {
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *token))
		}
		res, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
			return
		}

		content, err := io.ReadAll(res.Body)
		if err != nil {
			fmt.Println(err)
			return
		}
		if res.StatusCode != 200 {
			fmt.Printf("status code: %d, err: %s\n", res.StatusCode, string(content))
			return
		}
		objects := gjson.GetBytes(content, "objects")
		if !objects.Exists() || len(objects.Array()) == 0 {
			fmt.Println("dump finish")
			break
		}
		for _, object := range objects.Array() {
			clo := 0
			for k, v := range object.Map() {
				if row == 0 { // set header
					colMap[k] = fmt.Sprintf("%c", 'A'+clo)
					f.SetCellValue("Sheet1", fmt.Sprintf("%s%d", colMap[k], row+1), k)
				}
				f.SetCellValue("Sheet1", fmt.Sprintf("%s%d", colMap[k], row+2), v.String())
				if k == "id" {
					after = v.String()
				}
				clo++
			}
			row++
		}
		res.Body.Close()
	}
	// Save spreadsheet by the given path.
	if err := f.SaveAs(fmt.Sprintf("%s.xlsx", *class)); err != nil {
		fmt.Println(err)
		return
	}
}

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type jdbcDataSource struct {
	commonAttributes
	DriverClass   string `json:"driverClass"`
	Username      string `json:"username"`
	ConnectionURL string `json:"connectionUrl"`
}

type jrxmlFileSettings struct {
	JRXMLFile jrxmlFile `json:"jrxmlFile"`
}

type jrxmlFile struct {
	Label   string `json:"label"`
	Type    string `json:"type"`
	Content string `json:"content"`
}

type jrxmlFileReferenceSettings struct {
	JRXMLFileReference jrxmlFileReference `json:"jrxmlFileReference"`
}

type jrxmlFileReference struct {
	URI string `json:"uri"`
}

type jrxmlFileStruct struct {
	commonAttributes
	reportUnit
	JRXMLFileSettings jrxmlFileSettings `json:"jrxml"`
	DataSource        dataSource        `json:"dataSource"`
	//JRXMLFileReference *jrxmlFileReference `json:"jrxmlFileReference"`
}

type jrxmlFileReferenceStruct struct {
	commonAttributes
	reportUnit
	JRXMLFileReferenceSettings jrxmlFileReferenceSettings `json:"jrxml"`
	DataSource                 dataSource                 `json:"dataSource"`
}

type fileReference struct {
	URI string `json:"uri"`
}

type dataSourceReference struct {
	URI string `json:"uri"`
}

type dataSource struct {
	DataSourceReference dataSourceReference `json:"dataSourceReference"`
}

type commonAttributes struct {
	Label          string `json:"label"`
	URI            string `json:"uri"`
	Description    string `json:"description"`
	PermissionMask string `json:"permissionMask"`
	Version        string `json:"version"`
}

type reportUnit struct {
	//commonAttributes
	AlwaysPromptControls string `json:"alwaysPromptControls"`
	ControlsLayout       string `json:"controlsLayout"`
	// JRXML                *jrxmlFileStruct `json:"jrxml"`
	// DataSource           *dataSource      `json:"dataSource"`
}

// type foo struct {
// 	Price webutil.FormCurrency `json:"price"`
// }

// func (f foo) Validate() error {
// 	fmt.Printf("foo is called")
// 	return nil
// }

// type bar struct {
// 	decimal.Decimal
// }

// func (b bar) Validate() error {
// 	fmt.Printf("this is called")
// 	fmt.Printf(b.String())
// 	return nil
// }

func main() {
	//getUnauth()
	getReport(getCookie())

	// cBytes, err := json.Marshal(cookie)

	// if err != nil {
	// 	fmt.Printf("err: %s\n", err.Error())
	// 	os.Exit(1)
	// }

	// fmt.Printf("cookie bytes: %s\n", string(cBytes))
}

func getUnauth() {
	//var err error

	urlValues := url.Values{}
	urlValues.Set("contractor_id", "1")

	req, _ := http.NewRequest(
		http.MethodGet,
		`http://localhost:8090/jasperserver/rest_v2/reports/pac/foo_report.pdf?`+urlValues.Encode(),
		nil,
	)

	res, _ := http.DefaultClient.Do(req)

	fmt.Printf("unauth status code: %d\n", res.StatusCode)
}

func getReport(c *http.Cookie) {
	var err error

	urlValues := url.Values{}
	urlValues.Set("bid_id", "1")

	req, _ := http.NewRequest(
		http.MethodGet,
		`http://localhost:8090/jasperserver/rest_v2/reports/pac/bid/bid_contract.pdf?`+urlValues.Encode(),
		nil,
	)
	req.AddCookie(c)

	res, _ := http.DefaultClient.Do(req)

	if res.StatusCode != http.StatusOK {

		fmt.Printf("got status: %d\n", res.StatusCode)
		os.Exit(1)
		//fmt.Printf("body: %s\n", b.String())
	}

	b := &bytes.Buffer{}
	b.ReadFrom(res.Body)

	//fmt.Printf("body: %s\n", b.String())

	if err = ioutil.WriteFile("/tmp/foo.pdf", b.Bytes(), os.ModePerm); err != nil {
		fmt.Printf("couldn't write file\n")
		fmt.Printf(err.Error())
		os.Exit(1)
	}
}

func getCookie() *http.Cookie {
	var req *http.Request
	var res *http.Response
	//var err error

	req, _ = http.NewRequest(
		http.MethodGet,
		"http://localhost:8090/jasperserver//j_spring_security_check?j_username=user&j_password=bitnami",
		nil,
	)
	req.Header.Set("Accept", "application/json")
	res, _ = http.DefaultClient.Do(req)

	var cookie *http.Cookie

	if len(res.Cookies()) > 0 {
		cookie = res.Cookies()[0]
		buf := bytes.Buffer{}
		buf.ReadFrom(res.Body)

		f, err := os.Create("/tmp/foo.html")

		if err != nil {
			os.Exit(1)
		}

		if _, err = f.Write(buf.Bytes()); err != nil {
			os.Exit(1)
		}

		fmt.Printf("cooke stuff: %v\n", cookie)
		fmt.Printf("status code: %v\n", res.Status)

	} else {
		fmt.Printf("couldn't get cookie\n")
		fmt.Printf("status: %d\n", res.Status)
		os.Exit(1)
	}

	return cookie
}

func uploadDataSource(c *http.Cookie) {
	data := jdbcDataSource{
		commonAttributes: commonAttributes{
			URI:         "/pac",
			Label:       "pac_source",
			Description: "pac source",
		},
		DriverClass:   "org.postgresql.Driver",
		Username:      "root",
		ConnectionURL: "jdbc:postgresql://roach1:26257/pac_test",
	}

	dataBytes, err := json.Marshal(data)

	if err != nil {
		fmt.Printf(err.Error())
		os.Exit(1)
	}

	req, _ := http.NewRequest(
		http.MethodPost,
		"http://localhost:8090/jasperserver/rest_v2/resources/pac/",
		strings.NewReader(string(dataBytes)),
	)
	req.AddCookie(c)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/repository.jdbcDataSource+json")

	res, _ := http.DefaultClient.Do(req)

	if res.StatusCode != http.StatusCreated {
		fmt.Printf("expected status 201; got step 1 %d\n", res.StatusCode)
		b := &bytes.Buffer{}
		b.ReadFrom(res.Body)
		fmt.Printf("body: %s\n", b.String())
		os.Exit(1)
	}
}

func uploadReport(c *http.Cookie, rootDir, fileName string) {
	var err error
	var req *http.Request
	var res *http.Response

	file, err := os.Open("/home/travis/programming/go/src/github.com/TravisS25/pac-server/jasper/bid_files/licenses.jrxml")

	if err != nil {
		fmt.Printf("couldn't open file\n")
		os.Exit(1)
	}

	req, _ = http.NewRequest(http.MethodPost, "http://localhost:8090/jasperserver/rest_v2/resources/pac/", file)
	req.AddCookie(c)
	req.Header.Set("Content-Description", "test file")
	req.Header.Set("Content-Disposition", `attachment; filename=foo.jrxml`)
	req.Header.Set("Content-Type", "application/jrxml")

	res, _ = http.DefaultClient.Do(req)

	if res.StatusCode != http.StatusCreated {
		fmt.Printf("expected status 201; got step 1 %d\n", res.StatusCode)
		b := &bytes.Buffer{}
		b.ReadFrom(res.Body)
		fmt.Printf("body: %s\n", b.String())
		os.Exit(1)
	}

	buf := &bytes.Buffer{}
	buf.ReadFrom(file)

	report := jrxmlFileReferenceStruct{
		commonAttributes: commonAttributes{
			URI:            "/pac/foo_report",
			Label:          "foo_report",
			PermissionMask: "0",
			Version:        "0",
			Description:    "description",
		},
		reportUnit: reportUnit{
			AlwaysPromptControls: "true",
			ControlsLayout:       "popupScreen",
		},
		JRXMLFileReferenceSettings: jrxmlFileReferenceSettings{
			JRXMLFileReference: jrxmlFileReference{
				URI: "/pac/foo.jrxml",
			},
		},
		DataSource: dataSource{
			DataSourceReference: dataSourceReference{
				URI: "/pac/pac_source",
			},
		},
	}

	//cStr := base64.StdEncoding.EncodeToString(buf.Bytes())

	// report := jrxmlFileStruct{
	// 	commonAttributes: commonAttributes{
	// 		URI:            "/pac/license_report",
	// 		Label:          "license_report",
	// 		PermissionMask: "0",
	// 		Version:        "0",
	// 		Description:    "description",
	// 	},
	// 	reportUnit: reportUnit{
	// 		AlwaysPromptControls: "true",
	// 		ControlsLayout:       "popupScreen",
	// 	},
	// 	JRXMLFileSettings: jrxmlFileSettings{
	// 		JRXMLFile: jrxmlFile{
	// 			Label:   "license",
	// 			Type:    "jrxml",
	// 			Content: cStr,
	// 		},
	// 	},
	// 	DataSource: dataSource{
	// 		DataSourceReference: dataSourceReference{
	// 			URI: "/pac/pac_source",
	// 		},
	// 	},
	// }

	reportBytes, err := json.Marshal(report)

	if err != nil {
		fmt.Printf(err.Error())
		os.Exit(1)
	}

	// req, _ = http.NewRequest(
	// 	http.MethodGet,
	// 	"http://localhost:8090/jasperserver//j_spring_security_check?j_username=user&j_password=bitnami",
	// 	nil,
	// )
	// req.Header.Set("Accept", "application/json")
	// res, _ = http.DefaultClient.Do(req)

	// var cookie *http.Cookie

	// if len(res.Cookies()) > 0 {
	// 	cookie = res.Cookies()[0]
	// 	fmt.Printf("cooke stuff: %v\n", cookie)
	// } else {
	// 	fmt.Printf("couldn't get cookie\n")
	// 	os.Exit(1)
	// }

	//fmt.Printf("%s\n", string(reportBytes))

	// buf = &bytes.Buffer{}
	// buf.Read(reportBytes)

	req, _ = http.NewRequest(
		http.MethodPost,
		"http://localhost:8090/jasperserver/rest_v2/resources/pac/",
		strings.NewReader(string(reportBytes)),
	)
	req.AddCookie(c)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/repository.reportUnit+json")

	res, _ = http.DefaultClient.Do(req)

	if res.StatusCode != http.StatusCreated {
		b := &bytes.Buffer{}
		b.ReadFrom(res.Body)
		fmt.Printf("response: %s\n", b.String())

		fmt.Printf("expected status 201; got %d\n", res.StatusCode)
	}
}

func uploadFile(rootDir, filePath, fileName, contentType string) *http.Cookie {
	var err error
	var req *http.Request
	var res *http.Response

	req, _ = http.NewRequest(
		http.MethodGet,
		"http://localhost:8090/jasperserver//j_spring_security_check?j_username=user&j_password=bitnami",
		nil,
	)
	req.Header.Set("Accept", "application/json")
	res, _ = http.DefaultClient.Do(req)

	var cookie *http.Cookie

	if len(res.Cookies()) > 0 {
		cookie = res.Cookies()[0]
		fmt.Printf("cooke stuff: %v\n", cookie)
	} else {
		fmt.Printf("couldn't get cookie\n")
		os.Exit(1)
	}

	// baseResourceURL := "http://localhost:8090/jasperserver/rest_v2/resources/pac/test/"

	// req, _ = http.NewRequest(http.MethodGet, baseResourceURL+fileName, nil)
	// req.AddCookie(cookie)
	// res, _ = http.DefaultClient.Do(req)

	// if res.StatusCode == http.StatusOK {
	// 	fmt.Printf("File already exists")
	// 	return cookie
	// }

	file, err := os.Open(filePath)

	if err != nil {
		fmt.Printf("couldn't open file\n")
		os.Exit(1)
	}

	req, _ = http.NewRequest(http.MethodPost, "http://localhost:8090/jasperserver/rest_v2/resources/pac/test/", file)
	req.AddCookie(cookie)
	req.Header.Set("Content-Description", "test file")
	req.Header.Set("Content-Disposition", `attachment; filename=`+fileName)
	req.Header.Set("Content-Type", contentType)

	res, _ = http.DefaultClient.Do(req)

	if res.StatusCode != http.StatusCreated {
		fmt.Printf("expected status 201; got step 1 %d\n", res.StatusCode)
		b := &bytes.Buffer{}
		b.ReadFrom(res.Body)
		fmt.Printf("body: %s\n", b.String())
		os.Exit(1)
	}

	return cookie
}

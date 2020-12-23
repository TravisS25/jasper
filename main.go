package main

import (
	"bytes"
	"encoding/base64"
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

func main() {
	// rootDir := "/pac/"
	// fileName := "license"
	//c := getCookie()

	// uploadDataSource(c)
	// uploadReport(c, rootDir, fileName)
	//getReport(c)

	str := `<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" viewBox="0 0 664 374" width="664" height="374"><path d="M 80.000,140.500 C 81.599,145.213 82.000,145.000 84.000,149.500" stroke-width="5.248" stroke="black" fill="none" stroke-linecap="round"></path><path d="M 84.000,149.500 C 88.382,157.945 88.599,157.713 94.000,165.500" stroke-width="3.205" stroke="black" fill="none" stroke-linecap="round"></path><path d="M 94.000,165.500 C 103.490,177.008 102.882,177.445 113.000,188.500" stroke-width="2.222" stroke="black" fill="none" stroke-linecap="round"></path><path d="M 113.000,188.500 C 114.938,192.119 115.490,191.508 118.000,194.500" stroke-width="3.158" stroke="black" fill="none" stroke-linecap="round"></path><path d="M 118.000,194.500 C 128.451,202.972 128.438,202.619 140.000,209.500" stroke-width="2.478" stroke="black" fill="none" stroke-linecap="round"></path><path d="M 140.000,209.500 C 151.411,217.291 151.451,215.972 164.000,220.500" stroke-width="2.153" stroke="black" fill="none" stroke-linecap="round"></path><path d="M 164.000,220.500 C 177.837,224.206 176.911,223.791 191.000,222.500" stroke-width="2.217" stroke="black" fill="none" stroke-linecap="round"></path><path d="M 191.000,222.500 C 205.218,219.587 204.837,220.706 218.000,213.500" stroke-width="2.079" stroke="black" fill="none" stroke-linecap="round"></path><path d="M 218.000,213.500 C 236.655,203.965 236.718,204.587 254.000,192.500" stroke-width="1.743" stroke="black" fill="none" stroke-linecap="round"></path><path d="M 254.000,192.500 C 266.789,184.711 266.155,184.465 277.000,174.500" stroke-width="1.959" stroke="black" fill="none" stroke-linecap="round"></path><path d="M 277.000,174.500 C 283.377,166.975 283.789,167.711 288.000,158.500" stroke-width="2.333" stroke="black" fill="none" stroke-linecap="round"></path><path d="M 288.000,158.500 C 292.320,150.210 292.377,150.475 295.000,141.500" stroke-width="2.580" stroke="black" fill="none" stroke-linecap="round"></path><path d="M 295.000,141.500 C 298.462,131.921 296.320,137.210 296.000,132.500" stroke-width="3.214" stroke="black" fill="none" stroke-linecap="round"></path><path d="M 296.000,132.500 C 285.516,137.542 289.962,133.921 278.000,145.500" stroke-width="3.861" stroke="black" fill="none" stroke-linecap="round"></path><path d="M 278.000,145.500 C 268.732,156.383 268.516,155.542 262.000,168.500" stroke-width="2.381" stroke="black" fill="none" stroke-linecap="round"></path><path d="M 262.000,168.500 C 255.246,178.287 258.732,174.383 258.000,181.500" stroke-width="3.346" stroke="black" fill="none" stroke-linecap="round"></path><path d="M 258.000,181.500 C 266.503,184.843 261.246,185.287 274.000,182.500" stroke-width="3.642" stroke="black" fill="none" stroke-linecap="round"></path><path d="M 274.000,182.500 C 298.488,172.836 298.503,174.843 322.000,161.500" stroke-width="2.436" stroke="black" fill="none" stroke-linecap="round"></path><path d="M 322.000,161.500 C 328.777,154.522 326.488,159.336 330.000,155.500" stroke-width="3.505" stroke="black" fill="none" stroke-linecap="round"></path><path d="M 330.000,155.500 C 327.668,160.378 331.277,156.522 327.000,165.500" stroke-width="4.359" stroke="black" fill="none" stroke-linecap="round"></path><path d="M 327.000,165.500 C 324.275,173.276 326.168,171.878 327.000,178.500" stroke-width="3.467" stroke="black" fill="none" stroke-linecap="round"></path><path d="M 327.000,178.500 C 331.780,186.190 330.275,184.776 339.000,188.500" stroke-width="3.705" stroke="black" fill="none" stroke-linecap="round"></path><path d="M 339.000,188.500 C 357.944,192.213 356.780,194.190 377.000,194.500" stroke-width="2.053" stroke="black" fill="none" stroke-linecap="round"></path><path d="M 377.000,194.500 C 383.000,194.500 382.944,195.213 389.000,194.500" stroke-width="2.748" stroke="black" fill="none" stroke-linecap="round"></path><path d="M 389.000,194.500 C 407.000,194.500 407.000,194.500 425.000,194.500" stroke-width="2.012" stroke="black" fill="none" stroke-linecap="round"></path></svg>`
	foo := base64.StdEncoding

	bar := foo.EncodeToString([]byte(str))

	fmt.Printf("encoding: %s\n", bar)

	// dbytes, err := foo.DecodeString("PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0iVVRGLTgiIHN0YW5kYWxvbmU9Im5vIj8+PCFET0NUWVBFIHN2ZyBQVUJMSUMgIi0vL1czQy8vRFREIFNWRyAxLjEvL0VOIiAiaHR0cDovL3d3dy53My5vcmcvR3JhcGhpY3MvU1ZHLzEuMS9EVEQvc3ZnMTEuZHRkIj48c3ZnIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyIgdmVyc2lvbj0iMS4xIiB3aWR0aD0iMCIgaGVpZ2h0PSIwIj48L3N2Zz4=")

	// if err != nil {
	// 	fmt.Printf(err.Error())
	// }

	// fmt.Printf("decoded: %s", string(dbytes))
}

// func getReport(c *http.Cookie) (io.Reader, error) {
// 	var err error
// 	var req *http.Request
// 	var res *http.Response

// }

func getReport(c *http.Cookie) {
	var err error

	urlValues := url.Values{}
	urlValues.Set("contractor_id", "1")

	req, _ := http.NewRequest(
		http.MethodGet,
		`http://localhost:8090/jasperserver/rest_v2/reports/pac/foo_report.pdf?`+urlValues.Encode(),
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

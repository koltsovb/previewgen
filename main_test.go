package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http/httptest"
	"net/http"
	"net/textproto"
	"os"
	"testing"
)

type Case struct {
	name string
	method string
	resource string
	body string
	bodyFile string
	headers map[string]string
	result int
}

var testCases = []Case {
	// common
	{ "Invalid resource", "POST", "/", "", "", map[string]string{}, 404, },
	{ "Invalid method", "GET", "/files", "", "", map[string]string{}, 405, },
	
	// files in json
	{ "Good impty file-json", 
		"POST", 
		"/files", 
		`{ "type" : "files", "files" : [] }`, 
		"",
		map[string]string{"Content-Type":"application/json"}, 
		200, 
	},

	{ "File-json. With two files", "POST", "/files", "", "./test/files", map[string]string{"Content-Type":"application/json"}, 200 },

	// file ursl in json
	{ "urls-json. With empty json", "POST", "/files", `{ "type" : "urls", "files" : [] }`, "", map[string]string{"Content-Type":"application/json"}, 200 },
	{ "urls-json. With two urls", "POST", "/files", "", "./test/urls", map[string]string{"Content-Type":"application/json"}, 200 },

}

func TestPostFile(t *testing.T) {
	server := &Server{workDir : "."}

	for _, c := range testCases {
		
		body, err := getBody(c.body, c.bodyFile)
		if err != nil {
			t.Error("For test: ", c.name,
					"\n Cannot read request body from file: ", c.bodyFile,
					"\n Error: ", err.Error() )
		}

		req, _ := http.NewRequest(c.method, c.resource, bytes.NewReader(body) )
		w := httptest.NewRecorder()

		for k, v := range c.headers {
			req.Header.Set(k, v)
		}
		

		server.ServeHTTP(w, req)
		if w.Code != c.result {
			t.Error("For test: ", c.name,
					"\n expected :", c.result,
					"\n got : ", w.Code )
		}
	}
}

// test file in form-data 
func TestPostData(t *testing.T) {
	log.Println("")

	server := &Server{workDir : "."}


	var buf bytes.Buffer

	writer := multipart.NewWriter(&buf)

	// insert 1.png
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", `form-data; name="image_field"; filename="1.png"`)
	h.Set("Content-Type", "image/png")
	fWriter, _ := writer.CreatePart(h)

	f, _ := os.Open("./test/1.png")
	defer f.Close()
	_, _ = io.Copy(fWriter, f)


	// insert 1.jpg
	h = make(textproto.MIMEHeader)
	h.Set("Content-Disposition", `form-data; name="image_field"; filename="1.jpg"`)
	h.Set("Content-Type", "image/jpeg")
	fWriter, _ = writer.CreatePart(h)

	f, _ = os.Open("./test/1.jpg")
	defer f.Close()
	_, _ = io.Copy(fWriter, f)

	writer.Close()
	req, _ := http.NewRequest("POST", "/files", &buf)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	recoder := httptest.NewRecorder()

	server.ServeHTTP(recoder, req)
	if recoder.Code != http.StatusOK {
		t.Error("For test: ", "POST form-data",
			"\n expected: ", "200",
			"\n got: ", recoder.Code )
	}
}

func getBody(body string, bodyFile string) ([]byte, error) {
	if len(body) > 1 {
		return []byte(body), nil
	}

	if len(bodyFile) == 0 {
		return []byte{}, nil
	}

	f, err := os.Open(bodyFile)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	buf, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	return buf, nil
}
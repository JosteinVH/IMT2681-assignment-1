package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	. "jvh_local/TEST/api"
	. "jvh_local/TEST/data"
	"net/http"
	"net/http/httptest"
	"testing"
)


func TestInfoHandler(t *testing.T) {

	tt := []struct {
		met        string
		url        string
		statusCode int
		expTime    string
		expInfo    string
		expVer     string
	}{
		{met: "GET", url: "http://localhost:8080/igcinfo/", statusCode: http.StatusNotFound },
		{met: "GET", url: "http://notworking.com/igcinfo/api", statusCode: http.StatusNotFound },
		{met: "GET", url: "http://localhost:8080/igcinfo/api/", statusCode: http.StatusOK, expTime:"PT0S", expInfo:"Service for IGC tracks.", expVer:"v1"},
	}

	for _, tc := range tt {
		req, err := http.NewRequest(tc.met, tc.url, nil)

		if err != nil {
			t.Fatalf("could not create request: %v", err)
		}
		rec := httptest.NewRecorder()

		InfoHandler(rec, req)

		// Returns response generated by handler
		res := rec.Result()
		defer res.Body.Close()
		if res.StatusCode != tc.statusCode {
			t.Errorf("Expected status %v: got %v", tc.statusCode, res.Status)
		}

		// Don't want it to decode non-json
		if res.StatusCode == http.StatusNotFound {
			continue
		}
		var info Info
		if err := json.NewDecoder(res.Body).Decode(&info); err != nil {

			t.Fatalf("Could not decode json: %v", err)
		}

		if tc.expTime != info.Uptime {
			t.Errorf("Expected value: %v, got:  %v", tc.expTime,info.Uptime)
		}

		if tc.expInfo != info.Info {
			t.Errorf("Expected value: %v, got: %v", tc.expInfo, info.Uptime)
		}

		if tc.expVer != info.Version {
			t.Errorf("Expected value: %v,  got: %v ", tc.expVer, info.Version)
		}
	}
}


func TestGetAPI(t *testing.T) {

	tt := []struct {
		met        string
		url        string
		statusCode int
		lenOId	   int
	}{
		{met: "GET", url: "http://localhost:8080/igcinfo/", statusCode: http.StatusOK, lenOId: 0 },
		{met: "POST", url: "http://notworking.com/igcinfo/api", statusCode: http.StatusNotFound },
	}

	for _, tc := range tt {

		if tc.met == "GET" {

			req, err := http.NewRequest(tc.met, tc.url, nil)

			if err != nil {
				t.Errorf("Could not create request: %v", err)
			}
			rec := httptest.NewRecorder()

			GetAPI(rec, req)

			res := rec.Result()
			defer res.Body.Close()
			if res.StatusCode != tc.statusCode {
				t.Errorf("Expected statuscode %v, got %v", http.StatusOK, res.StatusCode)
			}

			var testId [] int

			if err := json.NewDecoder(res.Body).Decode(&testId); err != nil {
				t.Fatalf("Could not parse json %v", err)
			}

			if len(testId) != tc.lenOId {
				t.Fatalf("Expected length %v, got %v", tc.lenOId, len(testId))
			}
		}

		if tc.met == "POST" {
			tc := Url{
				"http://skypolaris.org/wp-content/uploads/IGS%20Files/Madrid%20to%20Jerez.igc",
			}

			//expected := 1

			content, err  := json.Marshal(tc)
			if err != nil {
				t.Errorf("Could not marshal data %v", err)
			}

			body := ioutil.NopCloser(bytes.NewBufferString(string(content)))


			req, err := http.NewRequest("POST", "http://localhost:8080/igcinfo/api/igc", body)
			if err != nil {
				t.Errorf("Could not create request, %v", err)
			}
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			http.HandlerFunc(PostAPI).ServeHTTP(rec,req)

			if rec.Code != http.StatusOK{
				t.Fatalf("Expected statuscode %v, got: %v", http.StatusOK, rec.Code)
			}

			data, err := ioutil.ReadAll(rec.Body)
			if err != nil {
				t.Fatalf("Could not read body %v", err.Error())
			}

			var testData Url

			err = json.Unmarshal(data, &testData)
			if err != nil {
				t.Errorf("Error during unmarshal %v", err.Error())
			}

			if testData.Url != ""{
				t.Errorf("Expected somedata, Got: %v", testData.Url,)
			}
		}
	}

}



/*
func TestApiHandler(t *testing.T) {
	req, err := http.NewRequest("", "http://localhost:8080/api/igc",xXSOMETHINGXx)
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	rec := httptest.NewRecorder()

	ApiHandler(rec,req)

	if req.Method != "POST" {
		t.Errorf("Something went wrong with DELETE %v", req.Method)
	}
	res := rec.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected status %v: got %v", http.StatusOK, res.StatusCode)
	}

}
*/
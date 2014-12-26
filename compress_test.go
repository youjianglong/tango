package tango

import (
	"fmt"
	"testing"
	"bytes"
	"compress/gzip"
	"compress/flate"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
)

type CompressExample struct {
	Compress // add this for ask compress according request accept-encoding
}

func (CompressExample) Get() string {
	return fmt.Sprintf("This is a auto compress text")
}

type GZipExample struct {
	GZip // add this for ask compress to GZip
}

func (GZipExample) Get() string {
	return fmt.Sprintf("This is a gzip compress text")
}

type DeflateExample struct {
	Deflate // add this for ask compress to Deflate, if not support then not compress
}

func (DeflateExample) Get() string {
	return fmt.Sprintf("This is a deflate compress text")
}

type NoCompress struct {
}

func (NoCompress) Get() string {
	return fmt.Sprintf("This is a non-compress text")
}

func TestCompressAuto(t *testing.T) {
	o := Classic()
	o.Get("/", new(CompressExample))
	testCompress(t, o, "This is a auto compress text")
}

func TestCompressGzip(t *testing.T) {
	o := Classic()
	o.Get("/", new(GZipExample))
	testCompress(t, o, "This is a gzip compress text")
}

func TestCompressDeflate(t *testing.T) {
	o := Classic()
	o.Get("/", new(DeflateExample))
	testCompress(t, o, "This is a deflate compress text")
}

func TestCompressNon(t *testing.T) {
	o := Classic()
	o.Get("/", new(NoCompress))
	testCompress(t, o, "This is a non-compress text")
}

func testCompress(t *testing.T, o *Tango, content string) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	req, err := http.NewRequest("GET", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}
	req.Header.Add("Accept-Encoding", "gzip, deflate")

	o.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)

	ce := recorder.Header().Get("Content-Encoding")
	if ce == "gzip" {
		r, err := gzip.NewReader(buff)
		if err != nil {
			t.Error(err)
		}
		defer r.Close()

		bs, err := ioutil.ReadAll(r)
		if err != nil {
			t.Error(err)
		}
		expect(t, string(bs), content)
	} else if ce == "deflate" {
		r := flate.NewReader(buff)
		defer r.Close()

		bs, err := ioutil.ReadAll(r)
		if err != nil {
			t.Error(err)
		}
		expect(t, string(bs), content)
	} else {
		expect(t, buff.String(), content)
	}
}
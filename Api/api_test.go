package mailrucloud

import (
	"crypto/md5"
	"fmt"
	"io"
	"math/rand"
	"os"
	"testing"
	"time"
)

const okfmt = "\t%-15s \u2713"
const failmrk = "\u2717"
const testdir = "/goapi_test/"
const testfile = "the_file.xyz"
const filesize = 1024 * 1024 * 10

func TestApi(t *testing.T) {
	var c *MailRuCloud
	var md5orig []byte
	if n, err := NewCloud(os.ExpandEnv("$MAILRU_USER"), os.ExpandEnv("$MAILRU_PASSWORD"), os.ExpandEnv("$MAILRU_DOMAIN")); err != nil {
		t.Fatal(err, failmrk)
	} else {
		c = n
	}
	t.Logf(okfmt, "NewCloud")

	_ = c.Remove(testdir)

	if err := c.Mkdir(testdir); err != nil {
		t.Error(err, failmrk)
	}
	t.Logf(okfmt, "Mkdir")

	if _, err := c.List(testdir); err != nil {
		t.Error(err, failmrk)
	}
	t.Logf(okfmt, "List")

	if f, err := os.Create(testfile); err != nil {
		t.Fatal(err, failmrk)
	} else {
		t.Logf(okfmt, "Create file")

		if fs, err := io.CopyN(f, rand.New(rand.NewSource(time.Now().UnixNano())), filesize); err != nil || fs != filesize {
			t.Fatal(err, failmrk)
		}
		f.Close()
	}
	t.Logf(okfmt, "Random fill")

	if f, err := os.Open(testfile); err != nil {
		t.Fatal(err, failmrk)
	} else {
		t.Logf(okfmt, "Open file")

		md5sum := md5.New()

		if fs, err := io.Copy(md5sum, f); err != nil || fs != filesize {
			t.Fatal(err, failmrk)
		}
		f.Close()

		md5orig = md5sum.Sum(nil)
		t.Logf(okfmt+" (%x)", "Calc md5sum ", md5orig)
	}
	if err := c.Upload(testfile, testdir+testfile, nil); err != nil {
		t.Fatal(err, failmrk)
	}
	t.Logf(okfmt, "Uploadfile")

	if err := os.Remove(testfile); err != nil {
		t.Fatal(err, failmrk)
	}
	t.Logf(okfmt, "Remove file")

	if err := c.Get(testdir+testfile, testfile, nil); err != nil {
		t.Fatal(err, failmrk)
	} else {
		t.Logf(okfmt, "Get")
		f, err := os.Open(testfile)
		if err != nil {
			t.Fatal(err, failmrk)
		}
		t.Logf(okfmt, "Open file")

		md5sum := md5.New()

		if fs, err := io.Copy(md5sum, f); err != nil || fs != filesize {
			t.Fatal(err, failmrk)
		}
		f.Close()
		md5get := md5sum.Sum(nil)
		t.Logf(okfmt+" (%x)", "Calc md5sum ", md5get)

		if fmt.Sprintf("%x", md5orig) != fmt.Sprintf("%x", md5get) {
			t.Errorf("MD5sum differs ! %s", failmrk)
		}
		t.Logf(okfmt, "MD5 sum match")
	}

	if err := c.Remove(testdir); err != nil {
		t.Error(err, failmrk)
	}
	t.Logf(okfmt, "Remove")

	if err := os.Remove(testfile); err != nil {
		t.Fatal(err, failmrk)
	}
	t.Logf(okfmt, "Remove file")
}

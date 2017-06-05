package mailrucloud

import (
	"crypto/rand"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"
	"log"
)

// mail.ru limits file size to 2GB
// the current implementation is limited to max int ( 4 bytes ) size of bytes.Buffer
// which used to send multipart form , 1024 bytes reserved for form data/fields
const MaxFileSize = 2*1024*1024*1024 - 1024

// Upload is a convenient method to upload files to the mail.ru cloud. 
// src is the local file path
// dst is the full destination file path
// ch  is a channel to report operation progress. can be nil.
func (c *MailRuCloud) Upload(src, dst string, ch chan<- int) (err error) {
	if err = c.GetShardInfo(); err != nil {
		return
	}
	f, err := os.Open(src)
	if err != nil {
		return
	}
	defer f.Close()
	s, err := f.Stat()
	if err != nil {
		return
	}
	if s.Size() <= MaxFileSize {
		return c.UploadFilePart(f, src, dst, 0, s.Size(), ch)
	} else {
		err = fmt.Errorf("File upload with size > %d bytes is not implemented (yet)", MaxFileSize)
		Logger.Println(err)
		return
	}
}

func (c *MailRuCloud) UploadFilePart(file *os.File, src, dst string, start, end int64, ch chan<- int) (err error) {
	pipeReader, pipeWriter := io.Pipe()
	writer := multipart.NewWriter(pipeWriter)
	writer.SetBoundary(randomBoundary())
	errChan := make(chan error, 1)
	go func() {
		defer pipeWriter.Close()
		part, err := writer.CreateFormFile("file", "filename")
		if err != nil {
			errChan <- err
			return
		}
		_, err = io.Copy(part, file)
		if err == nil {
			err = writer.Close()
		}
		errChan <- err
	}()
	Url := c.Shard.Upload
	stat, err := file.Stat()
	if err != nil {
		log.Panicf("Don't get file stat %v", file)
	}
	contentSize := stat.Size() + 180
	req, err := http.NewRequest("POST", Url, NewIoProgress(pipeReader, int(contentSize), ch))
	if err != nil {
		Logger.Println(err)
		return
	}

	req.ContentLength = contentSize
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Accept-Encoding", "*.*")
	//req.Body = ioutil.NopCloser(pipeReader)

	//	dump, _ := httputil.DumpRequestOut(req, true)
	//	fmt.Printf("%q", dump)
	r, err := c.Client.Do(req)
	if err != nil {
		Logger.Println(err)
		return
	}
	defer r.Body.Close()
	// Handling the error the routine may caused
	if err := <-errChan; err != nil {
		panic(err)
	}
	// Check the response
	br, err := ioutil.ReadAll(r.Body)
	if err != nil {
		Logger.Println(err)
		return
	}
	if r.StatusCode != http.StatusOK {
		err = fmt.Errorf("Put file failed. Status: %d, Msg: %s", r.StatusCode, string(br))
		Logger.Println(err)
		return
	}
	hs := strings.SplitN(strings.TrimSpace(string(br)), ";", 2)
	err = c.addFile(dst, hs[0], hs[1])
	return
}

func (c *MailRuCloud) addFile(dst, hash, size string) (err error) {
	Url := c.url("file/add")
	data := url.Values{
		"token":    {c.AuthToken},
		"home":     {dst},
		"conflict": {"strict"},
		"hash":     {hash},
		"size":     {size},
	}
	r, err := c.Client.PostForm(Url, data)
	if err != nil {
		Logger.Println(err)
		return
	}
	defer r.Body.Close()
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		Logger.Println(err)
		return
	}
	if r.StatusCode != http.StatusOK {
		err = fmt.Errorf("addFile failed. Status: %d, Msg: %s", r.StatusCode, string(b))
		Logger.Println(err)
		return
	}
	return
}

// the default function return too long boundary
// mail.ru does not accept it
func randomBoundary() string {
	var buf [15]byte
	_, err := io.ReadFull(rand.Reader, buf[:])
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%x", buf[:])
}

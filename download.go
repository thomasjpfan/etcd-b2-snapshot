package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

// B2Downloader downloads data from backblaze b2
type B2Downloader struct {
	ApplicationID         string
	ApplicationKey        string
	BucketName            string
	Object                string
	ObjectID              string
	DownloadRetryInterval int
	HTTPClient            *http.Client
}

// B2AuthorizeAccountResponse is the b2 authorize account response
type B2AuthorizeAccountResponse struct {
	AuthorizationToken string
	DownloadURL        string
}

// Download downloads data from b2
func (b B2Downloader) Download(path string) error {

	err := b.download(path)
	if err == nil {
		return nil
	}
	if err != nil && err.Error() == "File does not exist on B2" {
		return err
	}

	log.Printf("First download attempt failed trying again in %d ms", b.DownloadRetryInterval)
	ticker := time.NewTicker(time.Millisecond * time.Duration(b.DownloadRetryInterval))

	for range ticker.C {
		err := b.download(path)
		if err == nil {
			return nil
		}
		if err != nil && err.Error() == "File does not exist on B2" {
			return err
		}
		log.Printf("Download attempt failed trying again in %d ms", b.DownloadRetryInterval)
	}
	return nil
}

func (b B2Downloader) download(path string) error {
	idAndKey := fmt.Sprintf("%s:%s", b.ApplicationID, b.ApplicationKey)
	idAndKeyBase64 := base64.StdEncoding.EncodeToString([]byte(idAndKey))
	basicAuthString := fmt.Sprintf("Basic %s", idAndKeyBase64)

	request, _ := http.NewRequest(http.MethodGet, "https://api.backblazeb2.com/b2api/v1/b2_authorize_account", nil)
	request.Header.Set("Authorization", basicAuthString)

	respAuthorize, err := b.HTTPClient.Do(request)
	if err != nil {
		return err
	}

	if respAuthorize.StatusCode != 200 {
		return fmt.Errorf("Authorize account status code: %d", respAuthorize.StatusCode)
	}

	defer respAuthorize.Body.Close()
	body, err := ioutil.ReadAll(respAuthorize.Body)
	if err != nil {
		return err
	}

	authroizeAccountResp := B2AuthorizeAccountResponse{}
	err = json.Unmarshal(body, &authroizeAccountResp)
	if err != nil {
		return err
	}

	if len(b.ObjectID) != 0 {
		log.Printf("Using object id to download snapshot.db")
		return b.downloadFileWithID(authroizeAccountResp, path)
	}

	log.Printf("Using object name to download snapshot.db")
	return b.downloadFileWithName(authroizeAccountResp, path)
}

func (b B2Downloader) downloadFileWithName(aaResp B2AuthorizeAccountResponse, path string) error {
	url := fmt.Sprintf("%s/file/%s/%s", aaResp.DownloadURL, b.BucketName, b.Object)
	return b.downloadFile(aaResp.AuthorizationToken, url, path, "NAME")

}

func (b B2Downloader) downloadFileWithID(aaResp B2AuthorizeAccountResponse, path string) error {
	url := fmt.Sprintf("%s/b2api/v2/b2_download_file_by_id?fileId=%s", aaResp.DownloadURL, b.ObjectID)
	return b.downloadFile(aaResp.AuthorizationToken, url, path, "ID")

}

func (b B2Downloader) downloadFile(token, url, path, method string) error {
	request, _ := http.NewRequest(http.MethodGet, url, nil)
	request.Header.Set("Authorization", token)

	respDownload, err := b.HTTPClient.Do(request)
	if err != nil {
		return err
	}

	if respDownload.StatusCode == 404 {
		return fmt.Errorf("File does not exist on B2")
	}
	if respDownload.StatusCode != 200 {
		return fmt.Errorf("Download file with %s status code: %d", method, respDownload.StatusCode)
	}

	defer respDownload.Body.Close()

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, respDownload.Body)

	return err

}

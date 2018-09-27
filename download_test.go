package main

import (
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type DownloadTestSuite struct {
	suite.Suite
	Client         *http.Client
	TempDir        string
	SnapshotPath   string
	ApplicationID  string
	ApplicationKey string
	BucketName     string
	Object         string
	ObjectID       string
}

func TestDownloadUnitTestSuite(t *testing.T) {
	suite.Run(t, new(DownloadTestSuite))
}

func (s *DownloadTestSuite) SetupSuite() {
	applicationID := os.Getenv("EBS_B2_APPLICATION_ID")
	if len(applicationID) == 0 {
		s.T().Skipf("EBS_B2_APPLICATION_ID does not exists")
	}

	applicationKey := os.Getenv("EBS_B2_APPLICATION_KEY")
	if len(applicationKey) == 0 {
		s.T().Skipf("EBS_B2_APPLICATION_KEY does not exists")
	}

	bucketName := os.Getenv("EBS_B2_BUCKET_NAME")
	if len(bucketName) == 0 {
		s.T().Skipf("EBS_B2_BUCKET_NAME does not exists")
	}

	b2Object := os.Getenv("EBS_B2_OBJECT")
	if len(b2Object) == 0 {
		s.T().Skipf("EBS_B2_OBJECT does not exists")
	}

	b2ObjectID := os.Getenv("EBS_B2_OBJECT_ID")
	if len(b2ObjectID) == 0 {
		s.T().Skipf("EBS_B2_OBJECT does not exists")
	}

	s.Client = &http.Client{
		Timeout: time.Second * 10,
	}

	s.ApplicationID = applicationID
	s.ApplicationKey = applicationKey
	s.BucketName = bucketName
	s.Object = b2Object
	s.ObjectID = b2ObjectID

}

func (s *DownloadTestSuite) SetupTest() {
	dir, err := ioutil.TempDir("", "snapshot")
	s.Require().NoError(err)
	s.TempDir = dir
	s.SnapshotPath = filepath.Join(s.TempDir, "snapshot.db")
}

func (s *DownloadTestSuite) TearDownTest() {
	os.RemoveAll(s.TempDir)
}

func (s *DownloadTestSuite) Test_Download_With_Name() {

	downloader := &B2Downloader{
		ApplicationID:         s.ApplicationID,
		ApplicationKey:        s.ApplicationKey,
		BucketName:            s.BucketName,
		Object:                s.Object,
		DownloadRetryInterval: 3000,
		HTTPClient:            s.Client,
	}

	err := downloader.Download(s.SnapshotPath)
	s.Require().NoError(err)
	_, err = os.Stat(s.SnapshotPath)
	s.True(!os.IsNotExist(err))

}

func (s *DownloadTestSuite) Test_Download_With_Name_Does_Not_Exist() {

	downloader := &B2Downloader{
		ApplicationID:         s.ApplicationID,
		ApplicationKey:        s.ApplicationKey,
		BucketName:            s.BucketName,
		Object:                "hello",
		DownloadRetryInterval: 3000,
		HTTPClient:            s.Client,
	}

	err := downloader.Download(s.SnapshotPath)
	s.Require().Error(err)
	s.Equal("File does not exist on B2", err.Error())
}

func (s *DownloadTestSuite) Test_Download_With_ID() {

	downloader := &B2Downloader{
		ApplicationID:         s.ApplicationID,
		ApplicationKey:        s.ApplicationKey,
		BucketName:            s.BucketName,
		Object:                s.Object,
		ObjectID:              s.ObjectID,
		DownloadRetryInterval: 3000,
		HTTPClient:            s.Client,
	}

	err := downloader.Download(s.SnapshotPath)
	s.Require().NoError(err)
	_, err = os.Stat(s.SnapshotPath)
	s.True(!os.IsNotExist(err))

}

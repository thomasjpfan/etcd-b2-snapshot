package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ParserTestSuite struct {
	suite.Suite
}

func TestParserUnitTestSuite(t *testing.T) {
	suite.Run(t, new(ParserTestSuite))
}

func (s *ParserTestSuite) TearDownTest() {
	os.Clearenv()
}

func (s *ParserTestSuite) Test_Parser() {

	b2AppID := "backblazeid"
	b2AppKey := "backblazekey"
	b2BucketName := "backblazebucket"
	b2Object := "backblazeobject"
	b2ObjectID := "backblazeobjectID"
	b2DownloadRetry := 5000

	os.Setenv("EBS_B2_APPLICATION_ID", b2AppID)
	os.Setenv("EBS_B2_APPLICATION_KEY", b2AppKey)
	os.Setenv("EBS_B2_BUCKET_NAME", b2BucketName)
	os.Setenv("EBS_B2_OBJECT", b2Object)
	os.Setenv("EBS_B2_OBJECT_ID", b2ObjectID)
	os.Setenv("EBS_B2_DOWNLOAD_RETRY_INTERVAL",
		fmt.Sprintf("%d", b2DownloadRetry))

	spec, err := ParseENV()
	s.Require().NoError(err)

	s.Equal(b2AppID, spec.B2ApplicationID)
	s.Equal(b2AppKey, spec.B2ApplicationKey)
	s.Equal(b2BucketName, spec.B2BucketName)
	s.Equal(b2Object, spec.B2Object)
	s.Equal(b2ObjectID, spec.B2ObjectID)
	s.Equal(b2DownloadRetry, spec.B2DownloadRetryInterval)
}

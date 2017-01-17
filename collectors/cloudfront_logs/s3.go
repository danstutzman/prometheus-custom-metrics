package cloudfront_logs

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/cenkalti/backoff"
	"log"
	"strings"
)

const EXPECTED_LINE1 = "#Version: 1.0"
const EXPECTED_LINE2_V1 = "#Fields: date time x-edge-location sc-bytes c-ip cs-method cs(Host) cs-uri-stem sc-status cs(Referer) cs(User-Agent) cs-uri-query cs(Cookie) x-edge-result-type x-edge-request-id x-host-header cs-protocol cs-bytes time-taken"
const EXPECTED_LINE2_V2 = EXPECTED_LINE2_V1 + " x-forwarded-for ssl-protocol ssl-cipher x-edge-response-result-type"
const EXPECTED_LINE2_V3 = EXPECTED_LINE2_V2 + " cs-protocol-version"

type S3Connection struct {
	service    *s3.S3
	bucketName string
}

func NewS3Connection(credsPath, region, bucketName string) *S3Connection {
	log.Printf("Creating AWS session...")
	session, err := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Credentials: credentials.NewSharedCredentials(credsPath, ""),
			Region:      aws.String(region),
		},
	})
	if err != nil {
		panic(fmt.Errorf("Couldn't create AWS session: %s", err))
	}

	return &S3Connection{
		service:    s3.New(session),
		bucketName: bucketName,
	}
}

func (conn *S3Connection) ListPaths() []string {
	log.Printf("Listing objects in s3://%s...", conn.bucketName)

	var response *s3.ListObjectsV2Output
	var err error
	err = backoff.Retry(func() error {
		response, err = conn.service.ListObjectsV2(&s3.ListObjectsV2Input{
			Bucket: aws.String(conn.bucketName),
		})
		return err
	}, backoff.NewExponentialBackOff())
	if err != nil {
		log.Fatal(fmt.Errorf("Couldn't ListObjectsV2: %s", err))
	}

	paths := []string{}
	for _, object := range response.Contents {
		paths = append(paths, *object.Key)
	}
	return paths
}

func (conn *S3Connection) DownloadVisitsForPath(path string) []map[string]string {
	log.Printf("Downloading %s...", path)

	var response *s3.GetObjectOutput
	var err error
	err = backoff.Retry(func() error {
		response, err = conn.service.GetObject(&s3.GetObjectInput{
			Bucket: aws.String(conn.bucketName),
			Key:    aws.String(path),
		})
		return err
	}, backoff.NewExponentialBackOff())
	if err != nil {
		panic(fmt.Errorf("Couldn't GetObjet: %s", err))
	}

	reader, err := gzip.NewReader(response.Body)
	defer reader.Close()

	scanner := bufio.NewScanner(reader)

	if !scanner.Scan() {
		log.Fatal(fmt.Errorf("Expected at least one line of %s", path))
	}
	if scanner.Text() != EXPECTED_LINE1 {
		log.Fatal(fmt.Errorf("First line of %s should be %s but got: %s",
			path, EXPECTED_LINE1, scanner.Text()))
	}

	if !scanner.Scan() {
		log.Fatal(fmt.Errorf("Expected at least two lines of %s", path))
	}
	secondLine := scanner.Text()
	if secondLine != EXPECTED_LINE2_V1 && secondLine != EXPECTED_LINE2_V2 && secondLine != EXPECTED_LINE2_V3 {
		log.Fatal(fmt.Errorf("Expected second line of %s is: %s", path, secondLine))
	}

	visits := []map[string]string{}
	for scanner.Scan() {
		visit := map[string]string{}
		values := strings.Split(scanner.Text(), "\t")
		for i, colName := range strings.Split(secondLine, " ")[1:] {
			visit[colName] = values[i]
		}
		visits = append(visits, visit)
	}
	if err := scanner.Err(); err != nil {
		panic(fmt.Errorf("Error from scanner.Err: %s", err))
	}

	return visits
}

func (conn *S3Connection) DeletePath(path string) {
	log.Printf("Deleting %s", path)

	err := backoff.Retry(func() error {
		_, err := conn.service.DeleteObject(&s3.DeleteObjectInput{
			Bucket: aws.String(conn.bucketName),
			Key:    aws.String(path),
		})
		return err
	}, backoff.NewExponentialBackOff())
	if err != nil {
		panic(fmt.Errorf("Couldn't DeleteObject: %s", err))
	}
}

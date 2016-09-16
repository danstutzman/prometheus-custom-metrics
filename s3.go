package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"log"
)

func testS3(credsPath, region, bucketName string) {
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
	s3Service := s3.New(session)

	log.Printf("Listing objects in s3://%s...", bucketName)
	resp, err := s3Service.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
		//		ContinuationToken: aws.String("Token"),
		//		Delimiter:         aws.String("Delimiter"),
		//		EncodingType:      aws.String("EncodingType"),
		//		FetchOwner:        aws.Bool(true),
		//		MaxKeys:           aws.Int64(1),
		//		Prefix:            aws.String("Prefix"),
		//		StartAfter:        aws.String("StartAfter"),
	})
	if err != nil {
		log.Fatal(fmt.Errorf("Couldn't ListObjectsV2: %s", err))
	}
	fmt.Println(resp)
}

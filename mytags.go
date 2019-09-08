package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"os"
	_ "strings"
)

const (
	AppVersion = "0.0.2"
)

type Tags struct {
	Tags []Tag `json:"tags"`
}

type Tag struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

var (
	argProfile  = flag.String("profile", "", "Profile 名を指定.")
	argRegion   = flag.String("region", "ap-northeast-1", "Region 名を指定.")
	argEndpoint = flag.String("endpoint", "", "AWS API のエンドポイントを指定.")
	argKey      = flag.String("key", "", "取得したいタグのキーを指定..")
	argVersion  = flag.Bool("version", false, "バージョンを出力.")
)

func awsEc2Client(profile string, region string) *ec2.EC2 {
	var config aws.Config
	if profile != "" {
		creds := credentials.NewSharedCredentials("", profile)
		config = aws.Config{Region: aws.String(region), Credentials: creds, Endpoint: aws.String(*argEndpoint)}
	} else {
		config = aws.Config{Region: aws.String(region), Endpoint: aws.String(*argEndpoint)}
	}
	sess := session.New(&config)
	ec2Client := ec2.New(sess)
	return ec2Client
}

func myInstanceid() (instanceId string) {
	sess := session.Must(session.NewSession())
	svc := ec2metadata.New(sess)
	doc, _ := svc.GetInstanceIdentityDocument()
	instanceId = doc.InstanceID
	return instanceId
}

func myTags(ec2Client *ec2.EC2, instanceId string) [][]string {
	var instanceIds []*string
	instanceIds = append(instanceIds, aws.String(instanceId))
	input := &ec2.DescribeTagsInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("resource-id"),
				Values: instanceIds,
			},
		},
	}
	result, err := ec2Client.DescribeTags(input)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	var tags [][]string
	var tag []string
	for _, t := range result.Tags {
		tag = []string{
			*t.Key,
			*t.Value,
		}
		tags = append(tags, tag)
	}
	return tags
}

func outputSingleTag(data [][]string, key string) {
	var v string
	for _, record := range data {
		if key == record[0] {
			v = record[1]
		}
	}
	fmt.Printf(v)
}

func outputJson(data [][]string) {
	var rs []Tag
	for _, record := range data {
		r := Tag{Key: record[0], Value: record[1]}
		rs = append(rs, r)
	}
	rj := Tags{
		Tags: rs,
	}
	b, err := json.Marshal(rj)
	if err != nil {
		fmt.Println("JSON Marshal error:", err)
		return
	}
	os.Stdout.Write(b)
}

func main() {
	flag.Parse()

	if *argVersion {
		fmt.Println(AppVersion)
		os.Exit(0)
	}

	ec2Client := awsEc2Client(*argProfile, *argRegion)
	instanceId := myInstanceid()
	allTags := myTags(ec2Client, instanceId)
	if *argKey != "" {
		outputSingleTag(allTags, *argKey)
	} else {
		outputJson(allTags)
	}
}

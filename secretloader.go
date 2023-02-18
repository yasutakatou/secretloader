package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

var (
	secretName   string = ""
	region       string = "ap-northeast-1"
	versionStage string = "AWSCURRENT"
)

https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html

func main() {
	// svc := secretsmanager.New(
	// 	session.New(),
	// 	aws.NewConfig().WithRegion(region),
	// )

	// input := &secretsmanager.GetSecretValueInput{
	// 	SecretId:     aws.String(secretName),
	// 	VersionStage: aws.String(versionStage),
	// }

	// result, err := svc.GetSecretValue(input)
	// if err != nil {
	// 	panic(err.Error())
	// }
	// var secretString string = *result.SecretString

	// fmt.Println(secretString)

	name := "sample.txt"

	f, _ := os.Open(name)
	bu := bufio.NewReaderSize(f, 1024)

	for {
		line, _, err := bu.ReadLine()
		if err == io.EOF {
			break
		}
		fmt.Printf("%s\n", line)
	}
}

/*
-secretStr={}
-onlyOnce
-loop=30

[ALERT]
{SECRET1}
[ALLOWID]
{SECRET2}
{SECRET3}
[REJECT]
rejectrule1	escalation1	rm	passwd	vi
[HOSTS]
hostlabel1	pi1	192.168.0.1	22	pi1	myPassword1	/bin/bash
hostlabel2	pi2	192.168.0.2	22	pi2	myPassword2	/bin/ash
[USERS]
U024ZT3BHU5	~/	0
[ALLOW]
allowrule1	escalation1	cd	ls	cat	ps	df	find
[ADMINS]
admin
[REPORT]
C0256BTKP54
*/

/*
 * A tool that leverages AWS Secrets Manager to create a config file.
 *
 * @author    yasutakatou
 * @copyright 2023 yasutakatou
 * @license   xxx
 */
package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

func main() {
	_secretStr := flag.String("secretStr", "{}", "[-secretStr=Symbol to define the secret name. ex. []")
	_onlyOnce := flag.Bool("onlyOnce", true, "[-onlyOnce=Non-loop execution mode. (true is enable)]")
	_loopDuration := flag.Int("loopDuration", 1440, "[-loopDuration=Interval at which to execute the loop. (default is 1 day = 1440 minutes)]")
	_inputFile := flag.String("inputFile", "template.ini", "[-inputFile=Input file name and its path.")
	//_outputFile := flag.String("outputFile", "config.ini", "[-outputFile=Output file name and its path.")

	flag.Parse()

	secretStr := string(*_secretStr)

	//file, err := os.Create(*_outputFile)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer file.Close()

	for {
		f, _ := os.Open(*_inputFile)
		bu := bufio.NewReaderSize(f, 1024)

		for {
			line, _, err := bu.ReadLine()
			if err == io.EOF {
				break
			}
			fmt.Printf("%s\n", line)
			fmt.Println(line[0])
			fmt.Println(line[len(line)-1])

			if secretStr[0] == line[0] && secretStr[len(secretStr)-1] == line[len(line)-1] {
				secretName := string(string(line[1 : len(line)-1]))
				fmt.Println(secretName)
				readSecret(secretName)
			}

			//_, err = file.WriteString(string(line) + "\n")
			//if err != nil {
			//	log.Fatal(err)
			//}
		}
		f.Close()

		if *_onlyOnce == true {
			break
		}

		time.Sleep(time.Minute * time.Duration(*_loopDuration))
	}
	os.Exit(0)
}

func readSecret(secretName string) {
	config, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-2"))
	if err != nil {
		log.Fatal(err)
	}

	// Create Secrets Manager client
	svc := secretsmanager.NewFromConfig(config)

	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretName),
		VersionStage: aws.String("AWSCURRENT"),
	}

	result, err := svc.GetSecretValue(context.TODO(), input)
	if err != nil {
		panic(err.Error())
	}
	var secretString string = *result.SecretString

	fmt.Println(secretString)
}

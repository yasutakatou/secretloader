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
	"hash/crc64"
	"io"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

var (
	debug   bool
	logging bool
)

func main() {
	cksums := map[string]uint64{}

	_secretStr := flag.String("secretStr", "{}", "[-secretStr=Symbol to define the secret name. ex. []")
	_onlyOnce := flag.Bool("onlyOnce", true, "[-onlyOnce=Non-loop execution mode. (true is enable)]")
	_loopDuration := flag.Int("loopDuration", 1440, "[-loopDuration=Interval at which to execute the loop. (default is 1 day = 1440 minutes)]")
	_inputFile := flag.String("inputFile", "template.ini", "[-inputFile=Input file name and its path.")
	_Debug := flag.Bool("debug", false, "[-debug=debug mode (true is enable)]")
	_Logging := flag.Bool("log", false, "[-log=logging mode (true is enable)]")
	_outputFile := flag.String("outputFile", "config.ini", "[-outputFile=Output file name and its path.")
	_region := flag.String("region", "us-east-2", "[-region=AWS region.  (default: us-east-2)")

	flag.Parse()

	debug = bool(*_Debug)
	logging = bool(*_Logging)

	secretStr := string(*_secretStr)

	file, err := os.Create(*_outputFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	for {
		f, _ := os.Open(*_inputFile)
		bu := bufio.NewReaderSize(f, 1024)

		for {
			line, _, err := bu.ReadLine()
			if err == io.EOF {
				break
			}

			debugLog("line: " + line)

			if secretStr[0] == line[0] && secretStr[len(secretStr)-1] == line[len(line)-1] {
				secretName := string(string(line[1 : len(line)-1]))
				debugLog("match secret! :" + secretName)
				secret := readSecret(secretName, *_region)

				v, ok := cksums[secretName]

				if ok == false {
					debugLog("cksum not found..")
					cksums[secretName] = cksum(secret)
					writeString(file, string(line))
				} else {
					if v != cksum(secret) {
						debugLog("cksum found, and not equal cksum!")
						cksums[secretName] = cksum(secret)
						writeString(file, string(line))
					} else {
						debugLog("cksum found, and cksum.")
					}
				}
			} else {
				writeString(file, string(line))
			}
		}
		f.Close()

		if *_onlyOnce == true {
			break
		}

		time.Sleep(time.Minute * time.Duration(*_loopDuration))
	}
	os.Exit(0)
}

func writeString(file *os.File, str string) {
	_, err := file.WriteString(str + "\n")
	if err != nil {
		log.Fatal(err)
	}
}

func readSecret(secretName, region string) string {
	config, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		log.Fatal(err)
	}

	svc := secretsmanager.NewFromConfig(config)

	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretName),
		VersionStage: aws.String("AWSCURRENT"),
	}

	result, err := svc.GetSecretValue(context.TODO(), input)
	if err != nil {
		fmt.Println("secret not found! :" + secretName)
		return ""
	}
	var secretString string = *result.SecretString

	fmt.Println(secretString)
	return secretString
}

func cksum(data string) uint64 {
	crcTable := crc64.MakeTable(crc64.ECMA)
	fmt.Println(crc64.Checksum([]byte(data), crcTable))
	return crc64.Checksum([]byte(data), crcTable)
}

func debugLog(message string) {
	var file *os.File
	var err error

	if debug == true {
		fmt.Println(message)
	}

	if logging == false {
		return
	}

	const layout = "2006-01-02_15"
	const layout2 = "2006/01/02 15:04:05"
	t := time.Now()
	filename := t.Format(layout) + ".log"
	logHead := "[" + t.Format(layout2) + "] "

	if Exists(filename) == true {
		file, err = os.OpenFile(filename, os.O_WRONLY|os.O_APPEND, 0666)
	} else {
		file, err = os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0666)
	}

	if err != nil {
		log.Fatal(err)
		return
	}
	defer file.Close()
	fmt.Fprintln(file, logHead+message)
}

func Exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

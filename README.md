# (WIP) secretloader

A tool that leverages AWS Secrets Manager to create a config file.

# Solution

You work with a variety of tools every day. Do you have any of the following problems?

- To keep configuration management of legacy tools with individual configuration files confidential and to do it in the cloud
- I want to stop managing credentials without rotation.

I think there is a way to use 1Password for credential management. However, you can't eliminate both. It would also cost more money. :)
With this tool, you can eliminate both, plus you can manage everything, including terraform, in code(IaC)!

# Feature

- Pull plain text information from AWS Secrets Manager and generate a configuration file from a template file.
- There is the ability to loop and see regular updates

That's it! :)

# installation

If you want to put it under the path, you can use the following.

```
go get github.com/yasutakatou/secretloader
```

If you want to create a binary and copy it yourself, use the following.

```
git clone https://github.com/yasutakatou/secretloader
cd secretloader
go build .
```

[or download binary from release page.](https://github.com/yasutakatou/secretloader/releases) save binary file, copy to entryed execute path directory.

# uninstall

delete that binary. del or rm command. (it's simple!)

# Setup

- Create an IAM user with Secrets Manager access in AWS and pay out credentials
- Create a template file that matches the configuration file you wish to generate.
- Register the information you wish to embed in Secrets Manager

# Template file

{}, the name of the Secrets Manager will be read from AWS. Other lines are output as is.

```
[ALERT]
{SECRET1}
{SECRET2}
```

note) The {} character can be customized by specifying options.
note) Plain text in the Secrets Manager is output as is, even if it is multi-line or tab-delimited.

# Usecase
## 1. Generate configuration file

## 2. Operate in loop mode

## 3. Rotating operation of AWS credentials

# options

```
Usage of ./secretloader:
  -debug
        [-debug=debug mode (true is enable)]
  -inputFile string
        [-inputFile=Input file name and its path. (default "template.ini")
  -log
        [-log=logging mode (true is enable)]
  -loopDuration int
        [-loopDuration=Interval at which to execute the loop. (default is 1 day = 1440 minutes)] (default 1440)
  -onlyOnce
        [-onlyOnce=Non-loop execution mode. (true is enable)] (default true)
  -outputFile string
        [-outputFile=Output file name and its path. (default "config.ini")
  -region string
        [-region=AWS region.  (default: us-east-2) (default "us-east-2")
  -secretStr string
        [-secretStr=Symbol to define the secret name. ex. [] (default "{}")
```

## -debug

Run in the mode that outputs various logs.

## -inputFile

Specify a template file

note) You can use template files outside the current directory by specifying the path.

## -log

Specify the log file name.

## -loopDuration

In loop mode, this is the interval at which to loop.

note) The unit is minutes. The default is one day at 1440 minutes.

## -onlyOnce

This mode does not loop the operation.

note) Enabled by default

## -outputFile

Specify the output file name

note) You can create a file in a directory other than the current directory by specifying the path
note) The original file will be overwritten

## -region

Specify the region from which to read the AWS Secret Manager key

## -secretStr

Replace {} with another character specifying the secret name

note) Please make sure to specify with two letters. ex) [],"",<>,'' etc.

# Why the action of creating and overwriting a temporary configuration file?

Because exporting one line at a time will cause tools that support hot-loading of configuration files to behave incorrectly.

# license

Apache-2.0 License
BSD-3-Clause License

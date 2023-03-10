# secretloader

**A tool that leverages AWS Secrets Manager to create a config file.**

### v0.2

- Added support for original file **backup mode**.
	-  The output destination file was forced to be overwritten, so it was backed up.

# Solution

You work with a variety of tools every day. Do you have any of the following problems?

- To **keep configuration management of legacy tools** with **individual configuration files confidential** and to do it in the cloud
- I want to stop **managing credentials without rotation**.

I think there is a way to use 1Password for credential management. However, you can't eliminate both. It would also cost more money. :)<br>
With this tool, you can eliminate both, plus you can manage everything, including **terraform, in code(IaC)**!<br>
Tools that require configuration file updates can be run in containers!!<br>
It works on both **Linux and Windows**! (or MacOS if you compile it)<br>

# Feature

- Pull **plain text** information from **AWS Secrets Manager** and generate a configuration file from a template file.
- There is the ability to **loop** and see regular updates.

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

There are three ways to set up authentication as described in the **SDK** documentation below.

[Configuring the AWS SDK for Go V2](https://aws.github.io/aws-sdk-go-v2/docs/configuring-sdk/)

## Authentication and Authorization Settings

1) OS Environment

How to create an IAM user and define credentials in **OS environment variables**.

```
$ export AWS_ACCESS_KEY_ID=YOUR_AKID
$ export AWS_SECRET_ACCESS_KEY=YOUR_SECRET_KEY
```

2) Profile

How to create an IAM user and add credentials to the **profile**.

```
[yourprofile]
aws_access_key_id = <YOUR_TEST_ACCESS_KEY_ID>
aws_secret_access_key = <YOUR_TEST_SECRET_ACCESS_KEY>
```

3) IAM Role

This is a method of creating IAM roles and assigning them to **AWS resources** such as EC2.

![image](https://user-images.githubusercontent.com/22161385/220613961-cdae378a-11ea-4e33-8309-71004340ce91.png)

note) It is secure because it does not dispense credentials.<br>

## Create a template file

Create a **template file** that matches the configuration file you wish to generate.

## Register the secret information to be replaced in the Secret Manager.

Register the information you wish to embed in Secrets Manager<br>

note) Please register with Secret Manager in **plain text**.<br>

# Template file

**{}**, the **name of the Secrets Manager** will be read from AWS. Other lines are output as is.

```
[ALERT]
{SECRET1}
{SECRET2}
```

note) The {} character can be customized by specifying options.<br>
note) **Plain text** in the Secrets Manager is output as is, even if it is multi-line or tab-delimited.

# Usecase
## 1. Generate configuration file

For example, if you have the following configuration file and you want to keep the **myPassword** part secret

![image](https://user-images.githubusercontent.com/22161385/219945429-30bd87d5-c37d-4148-bc09-92312712f935.png)

Register the following in plain text in Secret Manager.

![image](https://user-images.githubusercontent.com/22161385/219945518-f820ac47-6b9d-4794-96c8-f352e1034c58.png)

Prepare template files. **{} to enclose the name registered in Secret Manager**.

![image](https://user-images.githubusercontent.com/22161385/219945634-4284c808-4f23-457b-927e-b802844b0f3b.png)

When the command is executed, a config file is generated.

```
$ secretloader -outputFile=slabot.ini
config file update!: slabot.ini
```

## 2. Operate in loop mode

**Periodically access** the Secret Manager and generate configuration files only when there are differences in the Secret.<br>

note) **The first run always creates a configuration file**.<br>
note) Checksum of Secret is checked, **so if there is no difference**, no new configuration file is created.

## 3. Rotating operation of AWS credentials

note) There are many ways to do this, but here are a few I've tried. Assume a case where credentials are scattered on each engineer's PC.<br>
<br>

Prepare a template file that reads a single secret.<br>

![image](https://user-images.githubusercontent.com/22161385/219945894-8b744f56-3290-4c0e-8a80-8642b2d46017.png)

Create two IAM users with access to Secret Manager.<br>
and,  Register those two in your AWS profile<br>

note) In the example below, IAM users for ProfileA and ProfileB have been created and registered<br>

![image](https://user-images.githubusercontent.com/22161385/219946029-b1b0c919-8d82-46bc-b379-ebd784897b8a.png)

Register the profile as it is in the Secret Manager.

![image](https://user-images.githubusercontent.com/22161385/219946088-1785f520-d8d3-404b-8add-99a98578bbef.png)

Configure the OS to start the following script when the **PC starts up**.

```
#!/bin/bash

export AWS_PROFILE=ProfileA
secretloader -outputFile=${HOME}/.aws/credentials
if [ $? -ne 0 ]; then
   export AWS_PROFILE=ProfileB
   secretloader -outputFile=${HOME}/.aws/credentials
fi
```

Now the profile will be created with the ProfileA information **every time when the PC starts up!**

```
$ bash -x update.sh
+ export AWS_PROFILE=ProfileA
+ AWS_PROFILE=ProfileA
+ ./secretloader -outputFile=/home/ady/.aws/credentials
config file update!: credentials
+ '[' 0 -ne 0 ']'
```

Update the credentials in ProfileA when it is time to update the credentials.<br>
Reflect the updated AWS_SECRET_ACCESS_KEY in the registered Secret Manager.<br>

note) **Update AWS_SECRET_ACCESS_KEY; do not change AWS_ACCESS_KEY_ID**.<br>

At the next execution, the profile generation fails because ProfileA cannot be read, but a profile with the new ProfileA information is generated because the read in ProfileB is generated continuously.<br>

```
$ bash -x update.sh
+ export AWS_PROFILE=ProfileA
+ AWS_PROFILE=ProfileA
+ ./secretloader -outputFile=/home/ady/.aws/credentials
secret not found! :SECRET1
+ '[' 1 -ne 0 ']'
+ export AWS_PROFILE=ProfileB
+ AWS_PROFILE=ProfileB
+ ./secretloader -outputFile=/home/ady/.aws/credentials
config file update!: /home/ady/.aws/credentials
```

The credentials for ProfileB can be rotated by allowing time for the new ProfileA information to percolate, and then updating the credentials for ProfileB.<br>

note) I'm assuming the script will run, so you shouldn't rotate it before summer vacation or before a long break. :)<br>
note) If you want to rotate other credentials, you can do so by creating two IAMs and modifying the script in the same way.

# options

```
Usage of ./secretloader:
  -backup
        [-backup=origin config backup mode (true is enable)] (default true)
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

## -backup

Backup the output destination file, if one exists.<br>

note) **Enabled by default**<br>

## -debug

Run in the mode that outputs various logs.

## -inputFile

Specify a template file<br>

note) You can use template files outside the current directory by specifying the path.<br>

## -log

Specify the log file name.

## -loopDuration

In loop mode, this is the interval at which to loop.<br>

note) The unit is **minutes**. The default is one day at 1440 minutes.<br>

## -onlyOnce

This mode does not loop the operation.<br>

note) **Enabled by default**<br>

## -outputFile

Specify the output file name<br>

note) You can create a file in a directory other than the current directory by specifying the path<br>

## -region

Specify the region from which to read the AWS Secret Manager key

## -secretStr

Replace {} with another character specifying the secret name<br>

note) **Please make sure to specify with two letters**. ex) [],"",<>,'' etc.<br>

# Why the action of creating and overwriting a temporary configuration file?

**Because exporting one line at a time will cause tools that support hot-loading of configuration files to behave incorrectly**.

# license

Apache-2.0 License<br>
BSD-3-Clause License

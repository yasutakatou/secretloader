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

# Setup

- Create an IAM user with Secrets Manager access in AWS and pay out credentials
- Create a template file that matches the configuration file you wish to generate.
- Register the information you wish to embed in Secrets Manager

# Template file


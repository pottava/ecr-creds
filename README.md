# ECR creds

## Description

This is a CLI tool which retrieves credentials from Amazon ECR.


## Installation

curl (macOS):

```
$ curl -Lo ecr-creds https://github.com/pottava/ecr-creds/releases/download/v1.0.0/darwin_amd64
$ chmod +x ecr-creds
```

curl (Linux):

```
$ curl -Lo ecr-creds https://github.com/pottava/ecr-creds/releases/download/v1.0.0/linux_amd64
$ chmod +x ecr-creds
```

go:

```
$ go get github.com/pottava/ecr-creds
```


## Parameters

Common parameters:

Environment Variables     | Argument        | Description                     | Required | Default
------------------------- | --------------- | ------------------------------- | -------- | ---------
AWS_ACCESS_KEY_ID         | access-key, a   | AWS `access key` for API access | *        |
AWS_SECRET_ACCESS_KEY     | secret-key, s   | AWS `secret key` for API access | *        |
AWS_DEFAULT_REGION        | region, r       | AWS `region` for API access     |          | us-east-1


## Usage

```console
$ ecr-creds -a AKIAIOSFODNN7EXAMPLE -s wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY script | sh
WARNING! Using --password via the CLI is insecure. Use --password-stdin.
Login Succeeded
```

```console
$ ecr-creds -a AKIAIOSFODNN7EXAMPLE -s wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY get
{
  "account": "123456789012",
  "host": "123456789012.dkr.ecr.us-east-1.amazonaws.com",
  "user": "AWS",
  "password": "xxxsomethingwhichcanbeusedasdockerpassword=",
  "endpoint": "https://123456789012.dkr.ecr.us-east-1.amazonaws.com",
  "expiresAt": "2018-12-31T12:30:00Z"
}
```

With environment variables:

```console
$ export AWS_ACCESS_KEY_ID=AKIAIOSFODNN7EXAMPLE
$ export AWS_SECRET_ACCESS_KEY=wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY
$ ecr-creds get account
123456789012
```

package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"os/signal"
	"regexp"
	"runtime/debug"
	"strings"
	"syscall"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/aws/aws-sdk-go/service/sts"
	lib "github.com/pottava/ecr-creds/lib"
	cli "gopkg.in/alecthomas/kingpin.v2"
)

var (
	version = "1.0.x"
	date    string
)

type config struct {
	AccessKey   *string
	SecretKey   *string
	Region      *string
	IsDebugMode bool
}

type getConfig struct {
	Field *string
}

func main() {
	app := cli.New("ecr-creds", "Managing Amazon ECR credentials")
	if len(version) > 0 && len(date) > 0 {
		app.Version(fmt.Sprintf("%s (built at %s)", version, date))
	} else {
		app.Version(version)
	}
	conf := &config{}
	conf.AccessKey = app.Flag("access-key", "AWS access key ID.").
		Short('a').Envar("AWS_ACCESS_KEY_ID").Required().String()
	conf.SecretKey = app.Flag("secret-key", "AWS secret access key.").
		Short('s').Envar("AWS_SECRET_ACCESS_KEY").Required().String()
	conf.Region = app.Flag("region", "AWS region.").
		Short('r').Envar("AWS_DEFAULT_REGION").Default("us-east-1").String()
	conf.IsDebugMode = os.Getenv("APP_DEBUG") == "1"

	cmdGet := app.Command("get", "Retrieve ECR credentials.")
	getConf := &getConfig{}
	getConf.Field = cmdGet.Arg("field", "Specify the return field").String()

	cmdScript := app.Command("script", "Return an ECR login script.")

	// Recover
	defer func() {
		if err := recover(); err != nil {
			if conf.IsDebugMode {
				debug.PrintStack()
			}
			lib.Errors.Fatal(err)
		}
	}()

	// Cancel
	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cancel()
		os.Exit(1)
	}()

	switch cli.MustParse(app.Parse(os.Args[1:])) {
	case cmdGet.FullCommand():
		exitCode, err := get(ctx, conf, getConf)
		if err != nil {
			lib.Errors.Fatal(err)
			return
		}
		os.Exit(int(aws.Int64Value(exitCode)))

	case cmdScript.FullCommand():
		exitCode, err := script(ctx, conf, getConf)
		if err != nil {
			lib.Errors.Fatal(err)
			return
		}
		os.Exit(int(aws.Int64Value(exitCode)))
	}
}

var (
	exitNormally  *int64
	exitWithError *int64
	https         = regexp.MustCompile("https://")
	awsAccountID  = regexp.MustCompile("([0-9]{12})")
)

func init() {
	exitNormally = aws.Int64(0)
	exitWithError = aws.Int64(1)
}

type cred struct {
	AWSAccountID        string `json:"account"`
	ECRHostName         string `json:"host"`
	DockerUser          string `json:"user"`
	DockerPassword      string `json:"password"`
	DockerLoginEndpoint string `json:"endpoint"`
	DockerCredExpiresAt string `json:"expiresAt"`
}

func get(ctx context.Context, conf *config, getConf *getConfig) (exitCode *int64, err error) {
	if conf.IsDebugMode {
		lib.PrintJSON(conf)
	}
	// Check AWS credentials
	sess, err := lib.Session(conf.AccessKey, conf.SecretKey, conf.Region, nil)
	if err != nil {
		return exitWithError, err
	}
	// Retrieve IAM identity
	switch aws.StringValue(getConf.Field) {
	case "account":
		if identity, e := retrieveIdentity(ctx, sess); e == nil {
			lib.Logger.Println(identity.Account)
			return exitNormally, nil
		}
	}
	// Retrieve ECR credentials
	cred, err := retrieveCreds(ctx, sess)
	if err != nil {
		return exitWithError, err
	}
	switch aws.StringValue(getConf.Field) {
	case "account":
		lib.Logger.Println(cred.AWSAccountID)
	case "host":
		lib.Logger.Println(cred.ECRHostName)
	case "user":
		lib.Logger.Println(cred.DockerUser)
	case "password":
		lib.Logger.Println(cred.DockerPassword)
	case "endpoint":
		lib.Logger.Println(cred.DockerLoginEndpoint)
	case "expiresAt":
		lib.Logger.Println(cred.DockerCredExpiresAt)
	default:
		lib.PrintJSON(cred)
	}
	return exitNormally, nil
}

func retrieveIdentity(ctx context.Context, sess *session.Session) (*sts.GetCallerIdentityOutput, error) {
	return sts.New(sess).GetCallerIdentityWithContext(ctx, &sts.GetCallerIdentityInput{})
}

func retrieveCreds(ctx context.Context, sess *session.Session) (*cred, error) {
	res, err := ecr.New(sess).GetAuthorizationTokenWithContext(
		ctx, &ecr.GetAuthorizationTokenInput{})
	if err != nil {
		return nil, err
	}
	if len(res.AuthorizationData) <= 0 {
		return nil, fmt.Errorf("No AuthorizationData was found")
	}
	auth := res.AuthorizationData[0]
	decoded, _ := base64.StdEncoding.DecodeString(aws.StringValue(auth.AuthorizationToken))
	creds := strings.Split(string(decoded), ":")
	password := ""
	if len(creds) > 1 {
		password = creds[1]
	}
	endpoint := aws.StringValue(auth.ProxyEndpoint)

	return &cred{
		AWSAccountID:        awsAccountID.FindString(endpoint),
		ECRHostName:         https.ReplaceAllString(endpoint, ""),
		DockerUser:          creds[0],
		DockerPassword:      password,
		DockerLoginEndpoint: endpoint,
		DockerCredExpiresAt: lib.TimeFormat(auth.ExpiresAt),
	}, nil
}

func script(ctx context.Context, conf *config, getConf *getConfig) (exitCode *int64, err error) {
	if conf.IsDebugMode {
		lib.PrintJSON(conf)
	}
	// Check AWS credentials
	sess, err := lib.Session(conf.AccessKey, conf.SecretKey, conf.Region, nil)
	if err != nil {
		return exitWithError, err
	}
	// Retrieve ECR credentials
	cred, err := retrieveCreds(ctx, sess)
	if err != nil {
		return exitWithError, err
	}
	lib.Logger.Printf(
		"docker login --username %s --password %s %s",
		cred.DockerUser,
		cred.DockerPassword,
		cred.DockerLoginEndpoint,
	)
	return exitNormally, nil
}

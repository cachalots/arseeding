package main

import (
	"github.com/everFinance/arseeding"
	"github.com/everFinance/arseeding/common"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name: "Arseeding",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "db_dir", Value: "./data/bolt", Usage: "bolt db dir path", EnvVars: []string{"DB_DIR"}},
			&cli.StringFlag{Name: "mysql", Value: "root@tcp(127.0.0.1:3306)/arseed?charset=utf8mb4&parseTime=True&loc=Local", Usage: "mysql dsn", EnvVars: []string{"MYSQL"}},
			&cli.StringFlag{Name: "key_path", Value: "./data/bundler-keyfile.json", Usage: "ar keyfile path", EnvVars: []string{"KEY_PATH"}},
			&cli.StringFlag{Name: "ar_node", Value: "https://arweave.net", EnvVars: []string{"AR_NODE"}},
			&cli.StringFlag{Name: "pay", Value: "https://api-dev.everpay.io", Usage: "pay url", EnvVars: []string{"PAY"}},
			&cli.BoolFlag{Name: "no_fee", Value: false, EnvVars: []string{"NO_FEE"}},
			&cli.BoolFlag{Name: "manifest", Value: false, EnvVars: []string{"MANIFEST"}},
			&cli.IntFlag{Name: "bundle_interval", Value: 120, Usage: "bundle tx on chain time interval(seconds)", EnvVars: []string{"BUNDLE_INTERVAL"}},

			&cli.BoolFlag{Name: "use_s3", Value: false, Usage: "run with s3 store", EnvVars: []string{"USE_S3"}},
			&cli.StringFlag{Name: "s3_acc_key", Value: "AKIATZSGGOHIV4QTYNH5", Usage: "s3 access key", EnvVars: []string{"S3_ACC_KEY"}},
			&cli.StringFlag{Name: "s3_secret_key", Value: "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", Usage: "s3 secret key", EnvVars: []string{"S3_SECRET_KEY"}},
			&cli.StringFlag{Name: "s3_prefix", Value: "arseed", Usage: "s3 bucket name prefix", EnvVars: []string{"S3_PREFIX"}},
			&cli.StringFlag{Name: "s3_region", Value: "ap-northeast-1", Usage: "s3 bucket region", EnvVars: []string{"S3_REGION"}},
			&cli.BoolFlag{Name: "use_4ever", Value: false, Usage: "run with 4everland s3 service", EnvVars: []string{"USE_4EVER"}},

			&cli.StringFlag{Name: "port", Value: ":8080", EnvVars: []string{"PORT"}},
		},
		Action: run,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func run(c *cli.Context) error {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	s := arseeding.New(
		c.String("db_dir"), c.String("mysql"), c.String("key_path"), c.String("ar_node"), c.String("pay"), c.Bool("no_fee"), c.Bool("manifest"),
		c.Bool("use_s3"), c.String("s3_acc_key"), c.String("s3_secret_key"), c.String("s3_prefix"), c.String("s3_region"),
		c.Bool("use_4ever"), c.String("port"))
	s.Run(c.String("port"), c.Int("bundle_interval"))

	common.NewMetricServer()

	<-signals

	return nil
}

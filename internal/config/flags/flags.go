package flags

import (
	"flag"
)

func Parse() (Options, error) {
	var opts Options
	flag.StringVar(&opts.runAddress, runAddressKey, "", "Run application address")
	flag.StringVar(&opts.dbURI, dbURIKey, "", "DB connection URI")
	flag.StringVar(&opts.accrualSystemURL, accrualSystemURL, "", "Accrual system URL")
	flag.StringVar(&opts.jwtSecret, jwtSecret, "", "JWT secret")
	flag.Parse()

	return opts, nil
}

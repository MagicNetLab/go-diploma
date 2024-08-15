package flags

import (
	"flag"
)

func Parse() (Options, error) {
	var opts Options
	flag.StringVar(&opts.runAddress, runAddressKey, "", "Run application address")
	flag.StringVar(&opts.dbUri, dbUriKey, "", "DB connection URI")
	flag.StringVar(&opts.accrualSystemUrl, accrualSystemUrl, "", "Accrual system URL")
	flag.Parse()

	return opts, nil
}

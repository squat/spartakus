package main

import (
	"fmt"
	"os"
	"time"

	flag "github.com/spf13/pflag"
	"k8s.io/spartakus/pkg/database"
	"k8s.io/spartakus/pkg/logr"
	"k8s.io/spartakus/pkg/volunteer"
)

var volunteerConfig = struct {
	clusterID      string
	period         time.Duration
	database       string
	printDatabases bool
}{}

type volunteerSubProgram struct{}

func (_ volunteerSubProgram) AddFlags(fs *flag.FlagSet) {
	fs.StringVar(&volunteerConfig.clusterID, "cluster-id", "", "Your cluster ID")
	fs.DurationVar(&volunteerConfig.period, "period", 24*time.Hour, "How often to send reports")
	fs.StringVar(&volunteerConfig.database, "database",
		"http://spartakus.k8s.io", "Send reports to this database; use --print-databases for a list of options") //FIXME: default to SSL
	fs.BoolVar(&volunteerConfig.printDatabases, "print-databases", false, "Print database options and exit")
}

func (_ volunteerSubProgram) Validate() error {
	if volunteerConfig.clusterID == "" {
		return fmt.Errorf("invalid value for --cluster-id: must not be empty")
	}
	if volunteerConfig.period == time.Duration(0) {
		return fmt.Errorf("invalid value for --period: must not be 0")
	}
	return nil
}

func (_ volunteerSubProgram) Main(log logr.Logger) error {
	if volunteerConfig.printDatabases {
		fmt.Printf("Example values for --database:\n")
		for _, str := range database.DatabaseOptions() {
			fmt.Println("    ", str)
		}
		os.Exit(0)
	}

	db, err := database.NewDatabase(log, volunteerConfig.database)
	if err != nil {
		return fmt.Errorf("failed to initialize database: %v", err)
	}

	volunteer, err := volunteer.New(log, volunteerConfig.clusterID, volunteerConfig.period, db)
	if err != nil {
		return fmt.Errorf("failed initializing volunteer: %v", err)
	}

	if err := volunteer.Run(); err != nil {
		return err
	}
	return nil
}

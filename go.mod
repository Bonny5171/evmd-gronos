module bitbucket.org/everymind/evmd-gronos

go 1.13

require (
	bitbucket.org/everymind/evmd-golib v1.7.0
	github.com/contribsys/faktory v1.2.0-1
	github.com/gorilla/mux v1.7.3
	github.com/robfig/cron/v3 v3.0.1
	github.com/spf13/cast v1.3.1
	github.com/spf13/pflag v1.0.5
)

replace bitbucket.org/everymind/evmd-golib => ./private/bitbucket.org/everymind/evmd-golib

replace bitbucket.org/everymind/gforce => ./private/bitbucket.org/everymind/gforce

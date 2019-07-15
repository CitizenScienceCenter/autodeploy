module github.com/citizensciencecenter/autodeploy

go 1.12

require (
	github.com/gorilla/mux v1.7.3
	github.com/spf13/viper v1.4.0
	gopkg.in/src-d/go-git.v4 v4.12.0 // indirect
)

replace github.com/citizensciencecenter/autodeploy/modules => ./modules

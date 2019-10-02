module github.com/citizensciencecenter/autodeploy

go 1.12

require (
	github.com/gorilla/mux v1.7.3
	github.com/spf13/viper v1.4.0
	golang.org/x/crypto v0.0.0-20191001170739-f9e2070545dc // indirect
	golang.org/x/net v0.0.0-20191002035440-2ec189313ef0 // indirect
	golang.org/x/sync v0.0.0-20190911185100-cd5d95a43a6e // indirect
	golang.org/x/sys v0.0.0-20191002091554-b397fe3ad8ed // indirect
	golang.org/x/text v0.3.2 // indirect
	golang.org/x/tools v0.0.0-20191001184121-329c8d646ebe // indirect
	gopkg.in/src-d/go-git.v4 v4.12.0 // indirect
)

replace github.com/citizensciencecenter/autodeploy/modules => ./modules

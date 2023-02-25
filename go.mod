module github.com/citizensciencecenter/autodeploy

go 1.12

require (
	github.com/citizensciencecenter/autodeploy/modules v0.0.0-20191002103314-664a692ed86d
	github.com/gorilla/mux v1.7.3
	github.com/spf13/viper v1.4.0
	golang.org/x/crypto v0.0.0-20191001170739-f9e2070545dc // indirect
	golang.org/x/net v0.0.0-20191002035440-2ec189313ef0 // indirect
	golang.org/x/sys v0.1.0 // indirect
	golang.org/x/text v0.3.2 // indirect
)

replace github.com/citizensciencecenter/autodeploy/modules => ./modules

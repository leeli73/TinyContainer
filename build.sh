#go build -ldflags "-s -w" -o tcontainer main.go
#upx --brute tcontainer tcontainer

#CGO_ENABLED=0 GOOS=linux GOARCH=amd64
export CGO_ENABLED=0
export GOOS=linux
#linux_i386
export GOARCH=386
go build -ldflags "-s -w" -o release/tcontainer_i386 main.go
#linux_amd64
export GOARCH=amd64
go build -ldflags "-s -w" -o release/tcontainer_amd64 main.go
upx --brute release/tcontainer_amd64 release/tcontainer_amd64
#linux_arm
export GOARCH=arm
go build -ldflags "-s -w" -o release/tcontainer_arm main.go
#linux_arm_x64
export GOARCH=arm64
go build -ldflags "-s -w" -o release/tcontainer_arm_x64 main.go
#linux_mips
export GOARCH=mips
go build -ldflags "-s -w" -o release/tcontainer_mips main.go
#linux_mips64le
export GOARCH=mips64le
go build -ldflags "-s -w" -o release/tcontainer_mips64le main.go
#linux_mipsle
export GOARCH=mipsle
go build -ldflags "-s -w" -o release/tcontainer_mipsle main.go
#linux_ppc64
export GOARCH=ppc64
go build -ldflags "-s -w" -o release/tcontainer_ppc64 main.go
#linux_ppc64le
export GOARCH=ppc64le
go build -ldflags "-s -w" -o release/tcontainer_ppc64le main.go
#linux_s390x
export GOARCH=s390x
go build -ldflags "-s -w" -o release/tcontainer_s390x main.go

# go test -bench . -benchmem
go test -bench=. -gcflags="-N -l" -benchmem -memprofile mem.out -cpuprofile cpu.out -trace trace.out

# go tool pprof -http=:8080 mem.out
# go tool pprof -http=:8080 cpu.out
# go tool trace -http=:8080 trace.out
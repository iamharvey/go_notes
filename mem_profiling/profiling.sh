# instead of using default '-inuse_space', we use '-alloc_space', which allows us to see
# every allocation regardless if it is still in memory or not at the time we take the profile
go tool pprof -alloc_space mem_profiling.test mem.out
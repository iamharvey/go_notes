# The following command will generate two files:
#   mem.out - contains the profile data
#   memcpu.test - contains a test binary we need to have access to symbols when looking at the profile data.
go test -run none -bench FindInFile -benchtime 3s -benchmem -memprofile mem.out
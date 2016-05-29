default:
	go build yotaman.go

install:
	cp ./yotaman ~/.bin/yotaman

cleanup:
	rm ./yotaman

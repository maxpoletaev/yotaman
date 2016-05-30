default:
	go build

install:
	cp ./yotaman ~/.bin/yotaman

cleanup:
	rm ./yotaman

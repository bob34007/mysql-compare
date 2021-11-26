BINARY_NAME=db-compare
build:
	go build -o $(BINARY_NAME) -v
clean:
	@if [ -f ${BINARY_NAME} ] ; then rm ${BINARY_NAME} ; fi

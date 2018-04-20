all:
	go build --buildmode=plugin -o base64encode.so base64encode.go
	go build --buildmode=plugin -o base64decode.so base64decode.go
	go build --buildmode=plugin -o cbcdecrypter.so cbcdecrypter.go
	go build --buildmode=plugin -o cbcencrypter.so cbcencrypter.go
	go build --buildmode=plugin -o cachebusting.so cachebusting.go

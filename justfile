name := "certscan"

default:
	@just --list --list-heading $'Available recipies:\n'

darwin:
	@echo "Building {{name}}-darwin_x86_64 binary"
	GOOS=darwin GOARCH=amd64 go build -o {{name}}-darwin_x86_64 main.go
	@echo "Compressing {{name}}-darwin_x86_64 binary..."
	tar -czvf ./release/{{name}}-darwin_x86_64.tar.gz ./{{name}}-darwin_x86_64
	rm ./{{name}}-darwin_x86_64
	@echo " "

darwin_arm:
	@echo "Building {{name}}-darwin_arm64 binary..."
	GOOS=darwin GOARCH=arm64 go build -o {{name}}-darwin_arm64 main.go
	@echo "Compressing {{name}}-darwin_arm64 binary..."
	tar -czvf ./release/{{name}}-darwin_arm64.tar.gz ./{{name}}-darwin_arm64
	rm ./{{name}}-darwin_arm64
	@echo " "

linux:
	@echo "Building {{name}}-linux_x86_64 binary..."
	GOOS=linux GOARCH=amd64 go build -o {{name}}-linux_x86_64 main.go
	@echo "Compressing {{name}}-linux_x86_64 binary..."
	tar -czvf ./release/{{name}}-linux_x86_64.tar.gz ./{{name}}-linux_x86_64
	rm ./{{name}}-linux_x86_64
	@echo " "

linux_arm:
	@echo "Building {{name}}-linux_arm64 binary..."
	GOOS=linux GOARCH=arm64 go build -o {{name}}-linux_arm64 main.go
	@echo "Compressing {{name}}-linux_arm64 binary..."
	tar -czvf ./release/{{name}}-linux_arm64.tar.gz ./{{name}}-linux_arm64
	rm ./{{name}}-linux_arm64
	@echo " "

all:
	just darwin
	just darwin_arm
	just linux
	just linux_arm

clean:
    @echo "Removing all {{name}} binaries and tar files from current directory..."
    rm {{name}}*
    @echo "Done!"
    @echo " "

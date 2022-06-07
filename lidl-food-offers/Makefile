
dep:
	apt install zip curl
	curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip"
	unzip awscliv2.zip
	sudo ./aws/install
	rm -R aws/
	rm awscliv2.zip

build:
	GOOS=linux go build -o offers-for-lidl
	zip handler.zip offers-for-lidl

upload:
	aws lambda update-function-code --function-name lidl-offers-skill --zip-file fileb://handler.zip

clean:
	rm offers-for-lidl
	rm handler.zip

deploy:	build upload clean
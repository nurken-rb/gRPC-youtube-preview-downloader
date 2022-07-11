# gRPC-youtube-preview-downloader

This application allows you to download video previews from youtube.

Used gRPC technology and postgresql as DB

Launch (GRPC Server):
1. Clone the repository
2. Enter the path to the config file
3. We raise postgres (for example, in a docker container)
4. We put down the necessary variables in the config
5. Build main.go and run it.
6. Profit
7. git clone http

git clone https://github.com/denieryd/grpc-youtube-preview-downloader
cd grpc-youtube-preview-downloader
export CONFIG_PATH="configs/config.example.json"
docker run --name thumbnail -e POSTGRES_PASSWORD=pass -p "5432:5432" -d postgres
# here we edit ./configs/config.example.json file 
# here we build ./cmd/main.go and run it
go build ./cmd/main.go
./main
  

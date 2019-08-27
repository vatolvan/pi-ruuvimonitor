# Go version of the Raspberry Pi RuuviTag scanner

## Building & Running

1. Install InfluxDB
2. Install Go
3. Build docker image with `sudo docker build -t ruuvitag .`
4. Run docker container with `sudo docker run -d --privileged -v /dev:/dev -u 0 --network host ruuvitag`

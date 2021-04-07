# Go version of the Raspberry Pi RuuviTag scanner

## Building & Running

1. Install [InfluxDB](https://pimylifeup.com/raspberry-pi-influxdb/)
2. Create influx database named "ruuvi_measurements"
3. Install Docker
4. Create module with `go mod init <package-name>`
5. Build docker image with `sudo docker build -t ruuvitag .`
6. Run docker container with `sudo docker run -d --privileged -v /dev:/dev -u 0 --network host ruuvitag`

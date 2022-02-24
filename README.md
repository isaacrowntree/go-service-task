# go-service-task

## config

In: pkg/reader/reader.go

You may tune const LinesPerFile to any number which seems to return the best results for the biggest file.

Seems to be the sweet spot for my Macbook Pro with 8 cores on sample 3 (outside Docker) is 200,000 (approx 300ms after initial 10 second first-run used to split the large file into chunks).

Config should appear in a config file but this is 'MVP'.

## testing

Run `go test ./...`

Not all code is 100% tested but there is integration test code and unit test code covering the basics. Channel tests and other interface mocking would normally be handled in these unit tests as well but I ran out of time.

## development

Run `go run main.go` and copy your sampleX.txt files into the same dir as `main.go`

Sample script:

```
curl -XPOST localhost:8000/ -d '{"filename":"sample3.txt", "from":"2021-07-06T23:00:00Z", "to": "2004-10-03T03:05:36Z"}
```

## production

Run `go build` and copy your sampleX.txt files in the same dir as `go-service-task`
Then run `./go-service-task`

You can test with this sample script:

```
curl -XPOST localhost:8000/ -d '{"filename":"sample3.txt", "from":"2021-07-06T23:00:00Z", "to": "2004-10-03T03:05:36Z"}
```

Obviously you would not co-locate these files in a real prod env (not the least due to Docker image size) but for 'mvp' purposes this is a 'for-now' solution.

### assumptions

sample1.txt, sample2.txt etc. are located in the same dir as `main.go` - this is why the test fixture `sample.txt` is copied there on running the tests.

/tmp exists in the same dir as `main.go`

`split` exists as a command-line utility and it supports splitting files using `-d`

## docker

Ensure you have your sample files in the same root dir as `Dockerfile` or `main.go` then:

`docker build --tag docker-go-service .`

`docker run --publish 8000:8000 docker-go-service`
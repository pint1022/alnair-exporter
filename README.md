[![Build Status](https://travis-ci.org/pint1022/alnair-exporter.svg?branch=master)](https://travis-ci.org/pint1022/alnair-exporter)

# Alnair Profiler Exporter

Feed AI application profiling metrics to a Prometheus compatible endpoint.

## Configuration

This exporter is setup to take input from environment variables. All variables are optional:


## Install and deploy

Run manually from Docker Hub:
```
docker run -d --restart=always -p 9171:9171 -e REPOS="pint1022/Kubeshare pint1022/alnair-exporter" pint1022/alnair-exporter
```

Build a docker image:
```
docker build -t <image-name> .
docker run -d --restart=always -p 9171:9171 -e REPOS="pint1022/Kubeshare, pint1022/alnair-exporter" <image-name>
```

## Docker compose

```
alnair-exporter:
    tty: true
    stdin_open: true
    expose:
      - 9171
    ports:
      - 9171:9171
    image: pint1022/alnair-exporter:latest
    environment:
      - REPOS=<REPOS you want to monitor>
      - GITHUB_TOKEN=<your github api token>

```

## Metrics

Metrics will be made available on port 9171 by default
An example of these metrics can be found in the `METRICS.md` markdown file in the root of this repository


## Version Release Procedure
Once a new pull request has been merged into `master` the following script should be executed locally. The script will trigger a new image build in docker hub with the new image having the tag `release-<version>`. The version is taken from the `VERSION` file and must follow semantic versioning. For more information see [semver.org](https://semver.org/).

Prior to running the following command ensure the number has been increased to desired version in `VERSION`: 

```bash
./release-version.sh
```

## Metadata
[![](https://images.microbadger.com/badges/image/pint1022/alnair-exporter.svg)](http://microbadger.com/images/pint1022/alnair-exporter "Get your own image badge on microbadger.com") [![](https://images.microbadger.com/badges/version/pint1022/alnair-exporter.svg)](http://microbadger.com/images/pint1022/alnair-exporter "Get your own version badge on microbadger.com")

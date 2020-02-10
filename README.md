# The URL shortener service
The purpose of this project is building HTTP-based RESTful API for managing short URLs and redirecting clients similar to bit.ly or goo.gl services.

For instance:
```
https://www.google.com/search?q=golang -> https:/short.it/o2MGIPLV
```

Where **o2MGIPLV** is called `the slug`

The generating of slugs is based on two increment counters. One of them increments every time when we need to generate a new slug, the other increments at every start of a service instance. It allows to avoid possible collisions without a significant impact on database performance because of a throughput bottleneck.

The service uses Redis as a backend database. The instance counter is kept in `instance_index`. URLs are stored as values with keys `{instance_index}:{slugs_counter}`

Example:
```json
{
    "6:0": "https://www.google.com/search?q=golang",
    "6:1":"https://github.com/ufoscout/docker-compose-wait",
    "instance_index":"7"
}
```

## How to run it
Locally
```
% REDIS_ADDRESS=redis:6379 REDIS_DATABASE=0 SLUGS_SALT="some_salt" SLUGS_MINLENGTH=16 go run ./cmd/url-shortener/main.go
```
Run with Docker Compose
```
% docker-compose -f docker-compose-redis.yaml up
% docker-compose -f docker-compose-service.yaml up
```

## API Reference
### POST /
Creates a new short URL

Request:
```json
{
    "url": "$url"
}
```
Response:
```json
{
    "data": {
        "slug": "$slug"
    }
}
```
Example:
```json
% curl -X POST --header "Content-Type: application/json" --data-raw '{"url": "https://www.google.com/search?q=golang"}' http://localhost:8080

{
    "data": {
        "slug": "o2MGIPLV"
    }
}
```

### GET /{slug}
Redirects the short URL to the original URL

Example:
```
% curl -L -X GET http://localhost:8080/o2MGIPLV
```

## Anticipated questions
- Would people open short URLs much more frequently than create them? Maybe it's better to split it up onto two services. One of them is responsible for creating short URLs, and the other is responsible for opening/redirecting them.
- What will we do if the length of a slug is changed? Probably, we'll have to make the logic a little bit more complicated.
- Any plans to allow people to make customized URLs?

## TODO
- [ ] Benchmarks
- [ ] Reduce using of mock.Anything
- [ ] Listing the number of times a short url has been accessed in the last 24 hours, past week and all time

## Nice to have
- [ ] Use multierror (github.com/hashicorp/go-multierror, go.uber.org/multierr)
- [ ] To optimize of the service, we can generate slugs beforehand
- [ ] Use go.uber.org/automaxprocs
- [ ] Use validators
- [ ] Metrics
- [ ] Tracing
- [ ] Documentation (godoc)
- [ ] Swagger
- [ ] CI/CD
- [ ] Deployment to Kubernetes
- [ ] Use Vault instead of environment variables

---

# The test assignment
The challenge is to build a HTTP-based RESTful API for managing Short URLs and redirecting
clients similar to bit.ly or goo.gl. Be thoughtful that the system must eventually support millions
of short urls. Please include a README with documentation on how to build, and run and test
the system. Clearly state all assumptions and design decisions in the README.

##### A Short Url:
1. Has one long url
2. Permanent; Once created
3. Is Unique; If a long url is added twice it should result in two different short urls.
4. Not easily discoverable; incrementing an already existing short url should have a low probability of finding a working short url.

##### Your solution must support:
1. Generating a short url from a long url
2. Redirecting a short url to a long url within 10 ms.
3. Listing the number of times a short url has been accessed in the last 24 hours, past
week and all time.
4. Persistence (data must survive computer restarts)

##### Shortcuts
1. No authentication is required
2. No html, web UI is required
3. Transport/Serialization format is your choice, but the solution should be testable via curl
4. Anything left unspecified is left to your discretion

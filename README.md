gspd2prometheus - a prometheus exporter for gpsd
================================================

This is a small project in the `go` programming language, connecting to
a gpsd (using its JSON-over-TCP protocol) and exposing various stats in
prometheus exporter syntax over HTTP for prometheus to collect/scrape
them.

It's my very first go project, and I had the basics put together in
about 15mins, spending another 90mins adding more gauges.  Your mileage may vary!

Usage
-----

```
Usage of gpsd2prometheus:
  -gpsd string
        remote gpsd host:port (default "localhost:2947")
  -http string
        local HTTP listen host:port (default ":2112")
```


Credits
-------

All the real workhorse code behind this project is in the following two
upstream libraries I'm using:

* [prometheus/client_golang](https://github.com/prometheus/client_golang)
* [go-gpsd](https://github.com/stratoberry/go-gpsd)

This project is just gluing together the above two libraries.

License
-------

I'm usually much more in favor of copyleft licenses, but given the two
libraries I use are MIT and Apache 2.0, I decided to go for a permissive
license in this project, too.  So the code is released under Apache 2.0,
see the COPYING file for details.

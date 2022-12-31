gspd2prometheus - a prometheus exporter for gpsd
================================================

This is a small project in the `go` programming language, connecting to
a gpsd (using its JSON-over-TCP protocol) and exposing various stats in
prometheus exporter syntax over HTTP for prometheus to collect/scrape
them.

It's my very first go project, and I had the basics put together in
about 15mins, spending another 90mins adding more gauges.  Your mileage may vary!

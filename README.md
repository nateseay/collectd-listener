# collectd-listener
A tool to help debug collectd installations/plugins. Listens on an UDP port, and tries to parse the Binary Protocol collectd info coming in over the wire.

Intro
-------
collectd-listener is a tool written in Go to listen on an UDP port. It listens for incoming collectd binary protocol packets- see here: https://collectd.org/wiki/index.php/Binary_protocol for more info.

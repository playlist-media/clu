# clu(ster)

[<img src="contrib/avatar.jpg" align="right" width="200" />](http://superkusokao.deviantart.com/art/Tron-Legacy-CLU-202369747)

This repo contains all the necessary assets and configuration to deploy Playlist.com and related services onto our infrastructure.  This repository is open source in the hopes that others will find it useful, though of course our application code is closed source.  Happy hacking!

### Requirements

* Docker
* Go
* Ruby

### Image Types

All images run consul for service discovery.

#### App
Contains the running code for our backend API - includes the Go runtime and database client libraries.  Each server has a version associated with it, so that it can be targeted by the load balancer.

#### Database / Cache
Runs MySQL either in master or slave mode, Redis in cluster mode, ElasticSearch, OrientDB in cluster mode, and InfluxDB in cluster mode.  It self-attaches disks at startup for data storage.

#### Load Balancer
This image load balances app instances - it builds its configuration based on consul and live reloads on changes.

#### Background / Monitor
This image provides web-accessible access to the cluster for administration purposes, as well as handles periodic background tasks (like MediaNet data imports and background analytics processing).  It also collects the logs in a centralized place.  It is a singleton.

### Inspiration / Credits

[hk](https://github.com/heroku/hk), [etcdctl](https://github.com/coreos/etcdctl), and [terraform](https://github.com/hashicorp/terraform) all provided excellent examples of Go CLI apps.

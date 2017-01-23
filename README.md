# Introduction

This tool has been developed to help network engineer to deal with undesirable traffic that is passing through their Internet network. This tool has been design to propose a graphical user interface to manage network features like:

* Filter traffic with BGP flowspec,
* Drop malicious traffic with BGP blackhole,
* Design and configuration of RTBH (Remote Triggered Blackhole),
* Interface analytics system,
* ... 

This first version is currently Alpha and needs to go through a set of test to make it an usable version. For now, it is only supporting BGP flowspec (RFC5575). __This tool is not suppose to be installed in production network but rather be used for lab / test purposes.__

# Building blocks and dependencies

This tool is relying on two open source software and related API.

The UI is using a Qt binding for Golang: https://github.com/therecipe/qt

BGP protocol stack: https://github.com/osrg/gobgp

## Using this tool in your network

This tool provides a BGP route reflector (RR using GoBGP as BGP stack) with an UI to inject BGP flowspec updates. The BGP RR propagates those updates to all its peers. In the current version, the UI can only connect to a local Go BGP daemon and doesn't support BGP clustering. The UI is using the gRPC API to interface GoBGP.

GoBGP needs to be installed manually (not yet via the UI). All neighbors have to be configured in the GoBGP configuration file. Once the GoBGP daemon is running, you launch the application. The UI will then connect to the BGP daemon and will be ready to push BGP flowspec updates to all involved neighbors.

## Typical configuration

You will have to create a VM or use a bare metal server, install the OS of your choice (software has been developed on ubuntu server), install all dependencies (mainly Qt / GoBGP / Unity) and you are ready to go. Make sure that your server is using is reachable to all BGP neighbors.

# Tutorial and features
* [Basic main window](https://github.com/Matt-Texier/local-mitigation-agent/blob/master/docs/main_win.md)
* Flowspec window
* Console window

# Install and configure you development machine 

# Licensing

This tool is licensed under the Apache License, Version 2.0. See [LICENSE](https://github.com/Matt-Texier/local-mitigation-agent/blob/master/LICENSE) for the full license text.



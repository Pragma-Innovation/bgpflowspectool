# Introduction

This tool has been developped to help network engineer to deal with undesirable traffic that is passing throught their Internet network. This tool has been design to propose a graphical user interface to manage network features like:

* Filter traffic with BGP flowspec,
* Drop malicious traffic with BGP blackhole,
* Design and configuration of RTBH (Remote Trigered Blackhole),
* Interface analytics system,
* ... 

This first version is currently Alpha and needs to go through a set of test to make it an usable version. For now, it is only supporting BGP flowspec (RFC5575). __This tool is not suppose to be installed in production network but rather be used for lab / test purposes.__

# Building blocks and dependencies

This tool is relying on two open source sofware and related API.

The UI is using a Qt binding for Golang: https://github.com/therecipe/qt

BGP protocol stack: https://github.com/osrg/gobgp

## Use this tool in your network

This tool provides a BGP route reflector (RR using GoBGP as BGP stack) with an UI to inject BGP updates. The BGP RR propates those updates to all its peers. In the current version, the UI can only connect to a local Go BGP daemon and doesn't support BGP clustering. The UI is using the gRPC API to interface GoBGP.

GoBGP needs to be installed manualy (not yet via the UI). All neighbors have to be configured in the GoBGP configuration file. Once the GoBGP daemon is running, you launch the application. The UI will connect to the BGP daemon and it is ready to push BGP flowspec updates to the BGP daemon.

## Typical configuration

You will have to create a VM or use a bare metal server, install the OS of your choice (software has been developped on ubuntu server), install all dependencies (mainly Qt / GoBGP / Unity) and you are ready to go. Make sure that your server is using is reachable to all BGP neighbors.

# Tutorial and features

# Licensing

This tool is licensed under the Apache License, Version 2.0. See [LICENSE](https://github.com/Matt-Texier/local-mitigation-agent/blob/master/LICENSE) for the full license text.


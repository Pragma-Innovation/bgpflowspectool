# Introduction

This tool has been developed to help network engineers to deal with undesirable traffic that is passing through their Internet network. This tool has been design to propose a graphical user interface to manage network features like:

* Filter traffic with BGP flowspec,
* Drop malicious traffic with BGP blackhole,
* Design and configuration of RTBH (Remote Triggered Blackhole),
* Interface analytics system,
* ... 

This first version is currently Alpha and needs to go through a set of test to make it an usable version. For now, it only supports BGP flowspec (RFC5575). __This tool is not suppose to be installed in production network but rather be used for lab / test purposes.__

# Building blocks and dependencies

This tool is relying on two open source software and their related API's.

* The UI is using a Qt binding for Golang: https://github.com/therecipe/qt

* BGP protocol stack: https://github.com/osrg/gobgp

## Using this tool in your network

This tool provides a BGP route reflector (RR using GoBGP as BGP stack) with an UI to inject BGP flowspec updates. The BGP RR propagates those updates to all its peers. In the current version, the UI can only connect to a local Go BGP daemon and doesn't support BGP clustering. The UI is using gRPC API to interface GoBGP.

GoBGP needs to be installed manually (not yet via the UI). All neighbors have to be configured in the GoBGP configuration file. Once the GoBGP daemon is running, you launch the application. The UI will then connect to the BGP daemon and will be ready to push BGP flowspec updates to all involved neighbors.

## Typical configuration

You will have to create a VM or use a bare metal server, install the OS of your choice (software has been developed on ubuntu server), install all dependencies (mainly Qt / GoBGP / Unity) and you are ready to go. Make sure that your server is using an IP address reachable to all BGP neighbors.

# Tutorial and features
The main window is organized with a logic of layered windows that you can bring to the front using the tool bar on the left of the main window.

There are two icons for now, one for the flowspec configuration an one for the so called console window that bings couple of troubleshooting feature.

![flowspec-win](/docs/main-window.png.jpg)

* [Flowspec window](https://github.com/Matt-Texier/local-mitigation-agent/blob/master/docs/flowspec_win.md)
* [Console window](https://github.com/Matt-Texier/local-mitigation-agent/blob/master/docs/console_win.md)

As of today, this tool has been tested against 7750 Nokia routers. Next step will be about testing this tool against other router vendors. Please note that this tool provide facilities to configure BGP flowspec but Gobgp provide the BGP stack. As such, most interoperability concerns are mainly related to gobgp stack.

# Install and configure your development machine

* Follow the installation process of GoBGP (please install version 1.22 or lastest): https://github.com/osrg/gobgp/blob/master/docs/sources/getting-started.md
* Follow the installation process of Qt Golang binding (please install Qt 5.9.1): https://github.com/therecipe/qt
  * Make sure that you allocate 8 GB of RAM to your VM
* ENV variables: The following example needs to be updated acccording to your machine but here is a snippet of my .bashrc file as an example

```
# go-lang variable
export GOPATH=/home/matthieu/go-work
export PKG_CONFIG_PATH=/usr/lib/x86_64-linux-gnu/pkgconfig
export LD_LIBRARY_PATH=/usr/lib/x86_64-linux-gnu
export GOBIN=/home/matthieu/go-work/bin
PATH=/home/matthieu/go-work/bin:$PATH
PATH=/usr/local/go/bin:$PATH
# qt variable
export QT_DIR=/home/matthieu/Qt
export QT_VERSION=5.9.1
PATH=/home/matthieu/Qt/5.7/gcc_64/bin:$PATH
```
* compile the tool and launch it (to be update with your own path):
```
~/go-work/src/github.com/Matt-Texier/local-mitigation-agent/UI$ qtdeploy build desktop .
~/go-work/src/github.com/Matt-Texier/local-mitigation-agent/UI$ ./deploy/linux_minimal/UI.sh
```

# Tool name and special dedicace

This tool is named "Gabu". This name is the nickname of a brilliant young french engineer that passed away way too early and is missing to anybody who crossed his way.

To enjoy a nice and still very interesting BGP flowspec presentation done by Frederic Gabut-Deloraine about the use of this protocol for DDoS mitigation, please follow this link: [Frederic FRnOG presentation](http://www.dailymotion.com/video/xtngjg_frnog-18-flowspec-frederic-gabut-deloraine-neo-telecoms_tech)

# Licensing

This tool is licensed under the Apache License, Version 2.0. See [LICENSE](https://github.com/Matt-Texier/local-mitigation-agent/blob/master/LICENSE) for the full license text.



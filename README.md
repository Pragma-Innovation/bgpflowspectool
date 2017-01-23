# Introduction

This tool has been developped to help network engineer to deal with undesirable traffic that is passing throught their Internet network. This tool has been design to propose a graphical user interface to manage network features like:

* Filter traffic with BGP flowspec,
* Drop malicious traffic with BGP blackhole,
* Design and configuration of RTBH (Remote Trigered Blackhole),
* Interface analytics system,
* ... 

This first version is currently Alpha and needs to go through a set of test to make it an usable version. For now, it is only supporting BGP flowspec (RFC5575).

# Building blocks and dependencies

This tool is relying on two open source sofware and related API.

The UI is using a Qt binding for Golang: https://github.com/therecipe/qt

BGP protocol stack: https://github.com/osrg/gobgp

# Tutorial and features

# Licensing

This tool is licensed under the Apache License, Version 2.0. See [LICENSE](https://github.com/Matt-Texier/local-mitigation-agent/blob/master/LICENSE) for the full license text.


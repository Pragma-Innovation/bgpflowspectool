# Features of the Flowspec window 

This window is organised around two main features:
* Edit, load and save Flowspec rules in a rule library
* Push, and delete rule to the BGP RIB of the GoBGP process

![Flowspec Window](/docs/main-window.png.jpg)

## Rules edition and upper part of the window

### Rule edition

The upper part of the window give the list of the flowspec rules available in the loaded library. 

The middle part is giving all tools to edit your flowspec rules. As BGP flowspec proposes a lot of features to filters TCP/UDP/ICMP protocol, there is a lot of tools to be used. However, you can use only a very few options to filter certain traffic. Typically you will use destination IP address and TCP/UDP ports extensively.

All fields of the rule (flowspec NLRI) is checked against a regexp such that it avoids any typo while entering the value

(please open an issue if you fill out an issue while using the tool)

### Rule library

All rules that has been created can be saved under a file. The file format is nothing else but a JSON format of each fields of the rule displayed in the rule list widget.

For easier understanding of the use of flowspec libraries, please try to load the existing library that comes with the repository : https://github.com/Pragma-Innovation/bgpflowspectool/tree/master/fslibexample

You will load this library from a library file as shown in the window :

![Loading library](/docs/lib-select.png.jpg)

## BGP flowspec RIB and update

In order to push a rule from the library to the RIB, select it in the list of rule and push the "Push to RIB button"

![Flowspec push](/docs/push-rule-to-rib.jpeg)

By clicking Load / reload button of the RIB tool, you will see the active flowspec routes. Select which adress familly you would like to display.

By pushing the delete button, you are able to delete the selected active entry. This delete button is using a newly introcuced feature of gpbgp that provides, for each RIB entry, a UUID. Behind the scene, this UUID is used to identify the entry and delete the associated path.

For more easier use, the RIB tool can be detached from the Flowspec main window (the RIB window is a dock window of the flowspec window).

![Flowspec delete rib](/docs/dock-window.png.jpg)


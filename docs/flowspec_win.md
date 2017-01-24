# Features of the Flowspec window 

This window is organised around two main features:
* Edit, load and save Flowspec rules in a rule library
* Push, and delete rule to the BGP RIB of the GoBGP process

![Flowspec Window](/docs/flowspec-win.png)

## Rules edition and upper part of the window

### Rule edition

The right part is allowing the user to create, update, duplicate a flowspec rule.

All fields of the rule (flowspec NLRI) is checked against a regexp such that it avoids any typo while entering the value

(please open an issue if you fill out an issue while using the tool)

### Rule library

All rules that has been created can be saved under a file. The file format is nothing else but a JSON format of each fields of the rule displayed in the rule list widget.

## BGP flowspec RIB and update

In order to push a rule from the library to the RIB, select it in the list of rule and push the "Push to RIB button"



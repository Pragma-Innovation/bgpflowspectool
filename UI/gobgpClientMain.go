// Copyright 2017 PRAGMA INNOVATION

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//     http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Provide UI for the whole tool
// There is a main window looking a bit like a dock
// with push buttons opening each window managing each
// function of the software

package main

import (
    "os"
    "os/user"
    "time"
    "fmt"
    "log"
    "io/ioutil"
    "encoding/json"
    "strings"
    "github.com/therecipe/qt/core"
    "github.com/therecipe/qt/gui"
    "github.com/therecipe/qt/widgets"
    api "github.com/osrg/gobgp/api"
    bgpcli "github.com/Matt-Texier/local-mitigation-agent/gobgpclient"
    "google.golang.org/grpc"
)

// data strcutures used by both API functions and UI
// BGP flowspec update structure as exported from UI

type BgpFsRule struct {
    DstPrefix       string
    SrcPrefix       string
    AddrFam         string
    Port            string
    SrcPort         string
    DstPort         string
    TcpFlags        string
    IcmpType        string
    IcmpCode        string
    ProtoNumber     string
    PacketLen       string
    Dscp            string
    IpFrag          string
    Action          string
    ActSisterValue  string
}

var BgpFsActivLib = []BgpFsRule{
    // {DstPrefix: "1.1.1.1/32", SrcPrefix: "2.2.2.2/32", AddrFam: "IPv4", Port: "=8080",
    //  SrcPort: "=80", DstPort: "=443", TcpFlags: "syn", IcmpType: "", IcmpCode: "", ProtoNumber: "=6",
    //  PacketLen: ">1024", Dscp: "=22", IpFrag: "", Action: "Drop",},
    // {DstPrefix: "3.3.3.3/32", SrcPrefix: "4.4.4.4/32", AddrFam: "IPv4", Port: "=8080",
    //  SrcPort: "<80", DstPort: ">443", TcpFlags: "syn", IcmpType: "", IcmpCode: "", ProtoNumber: "=6",
    //  PacketLen: "<=1024", Dscp: "=22", IpFrag: "", Action: "Drop",},
    // {DstPrefix: "5.5.5.5/32", SrcPrefix: "6.6.6.6/32", AddrFam: "IPv4", Port: ">=8080",
    //  SrcPort: ">=80", DstPort: ">=443", TcpFlags: "syn", IcmpType: "", IcmpCode: "", ProtoNumber: "=6 =8",
    //  PacketLen: "=1024", Dscp: ">=22&<=55", IpFrag: "", Action: "Drop",},
}

const (
    FS_ACT_DROP int = iota
    FS_ACT_SHAPE
    FS_ACT_REDIRECT
    FS_ACT_MARK
)

var flowSpecActNameMap = map[int]string{
    FS_ACT_DROP:        "Drop",
    FS_ACT_SHAPE:       "Shape",
    FS_ACT_REDIRECT:    "Redirect",
    FS_ACT_MARK:        "Marking",
}

var flowSpecActValueMap = map[string]int{
    flowSpecActNameMap[FS_ACT_DROP]:        FS_ACT_DROP,
    flowSpecActNameMap[FS_ACT_SHAPE]:       FS_ACT_SHAPE,
    flowSpecActNameMap[FS_ACT_REDIRECT]:    FS_ACT_REDIRECT,
    flowSpecActNameMap[FS_ACT_MARK]:        FS_ACT_MARK,
}

var flowSpecActStrings = []string{"Drop", "Shape", "Redirect", "Marking"}

var (
    editAddrFamIpv4 *widgets.QRadioButton
    editAddrFamIpv6 *widgets.QRadioButton
    ribAddrFamIpv4 *widgets.QRadioButton
    ribAddrFamIpv6 *widgets.QRadioButton
    editRuleSrcPrefixLineEdit *widgets.QLineEdit
    editRuleDstPrefixLineEdit *widgets.QLineEdit
    editRuleIcmpTypeLineEdit *widgets.QLineEdit
    editRuleIcmpCodeLineEdit *widgets.QLineEdit
    editRuleIpProtoLineEdit *widgets.QLineEdit
    editRulePortLineEdit *widgets.QLineEdit
    editRuleSrcPortLineEdit *widgets.QLineEdit
    editRuleDstPortLineEdit *widgets.QLineEdit
    editRuleTcpFlagFilterLine *widgets.QLineEdit
    editRuleLenLineEdit *widgets.QLineEdit
    editRuleDscpLineEdit *widgets.QLineEdit
    editRuleFragFilterLine *widgets.QLineEdit
    editRuleActSisterValueLine *widgets.QLineEdit
    editRuleActionCombo *widgets.QComboBox
    editRuleTree *widgets.QTreeWidget
    consoleWindow *widgets.QMainWindow
    flowspecWindow *widgets.QMainWindow
)


var client api.GobgpApiClient

var (
    windowFlowSpecCreated bool
    windowBgpConsoleCreated bool
    AddrFamilyIpv4Checked bool
    AddrFamilyIpv6Checked bool
)

var ribRadioFamilychecked string = "ipv4-flowspec"
var ribActiveFamily string = "ipv4-flowspec"

var (
    regexpIpv4Validation *core.QRegExp
    regexpIpv6Validation *core.QRegExp
    regexpPortValidation *core.QRegExp
    regexpByteValueValidation *core.QRegExp
    regexpPckLenValidation *core.QRegExp
    regexpIpv4SrcValidator *gui.QRegExpValidator
    regexpIpv4DstValidator *gui.QRegExpValidator
    regexpIpv6SrcValidator *gui.QRegExpValidator
    regexpIpv6DstValidator *gui.QRegExpValidator
    regexpPortValidator *gui.QRegExpValidator
    regexpSrcPortValidator *gui.QRegExpValidator
    regexpDstPortValidator *gui.QRegExpValidator
    regexpIcmpTypeValidator *gui.QRegExpValidator
    regexpIcmpCodeValidator *gui.QRegExpValidator
    regexpProtoNumValidator *gui.QRegExpValidator
    regexpDscpValidator *gui.QRegExpValidator
    regexpPckLenValidator *gui.QRegExpValidator
)

func main() {
    // initialise boolean that tell us if sub-windows is already reated
    windowFlowSpecCreated = false
    windowBgpConsoleCreated = false

    // launch gobgp API client
    timeout := grpc.WithTimeout(time.Second)
    conn, rpcErr := grpc.Dial("localhost:50051", timeout, grpc.WithBlock(), grpc.WithInsecure())
    if rpcErr != nil {
        fmt.Printf("GoBGP is probably not running on the local server ... Please start gobgpd process !\n")
        fmt.Println(rpcErr)
        return
    }
    client = api.NewGobgpApiClient(conn)

    widgets.NewQApplication(len(os.Args), os.Args)
    var toolbarWindow = widgets.NewQMainWindow(nil, 0)
    toolbarWindow.SetGeometry(core.NewQRect4(100, 100, 400, 50))
    toolbarWindow.SetWindowTitle("Gabu")
    var toolbarCentralWidget = widgets.NewQWidget(nil, 0)
    var toolbarCentralLayout = widgets.NewQHBoxLayout()
    toolbarCentralLayout.SetSpacing(6)
    toolbarCentralLayout.SetContentsMargins(11, 11, 11, 11)
    toolbarCentralWidget.SetLayout(toolbarCentralLayout)
    toolbarWindow.SetCentralWidget(toolbarCentralWidget)
    // main window "toolbar" push button
    var toolbarConsolePush = widgets.NewQPushButton2("GoBgp Console", toolbarCentralWidget)
    var toolbarFlowSpecPush = widgets.NewQPushButton2("FlowSpec RIB", toolbarCentralWidget)

    var toolbarButtonSizePolicy = widgets.NewQSizePolicy()
    toolbarButtonSizePolicy.SetHorizontalPolicy(widgets.QSizePolicy__Expanding)
    toolbarButtonSizePolicy.SetVerticalPolicy(widgets.QSizePolicy__Expanding)
    toolbarButtonSizePolicy.SetHorizontalStretch(0)
    toolbarButtonSizePolicy.SetVerticalStretch(0)
    toolbarConsolePush.SetSizePolicy(toolbarButtonSizePolicy)
    toolbarFlowSpecPush.SetSizePolicy(toolbarButtonSizePolicy)

    // Connect buttons to functions
    toolbarConsolePush.ConnectClicked(func(_ bool) { toolbarConsolButtonClicked() })
    toolbarFlowSpecPush.ConnectClicked(func(_ bool) { toolbarFspecButtonPushed() })
    // add button to main layout
    toolbarCentralLayout.AddWidget(toolbarConsolePush, 0, 0)
    toolbarCentralLayout.AddWidget(toolbarFlowSpecPush, 0, 0)
    toolbarWindow.Show()
    widgets.QApplication_Exec()

}

func toolbarConsolButtonClicked() {
    if(windowBgpConsoleCreated) {
        consoleWindow.Raise()

    } else {
        consoleWin()
        windowBgpConsoleCreated = true
    }
}

func toolbarFspecButtonPushed() {
    if(windowFlowSpecCreated) {
        flowspecWindow.Raise()
    } else {
        flowspecWin()
        windowFlowSpecCreated = true
    }

}





func consoleWin() {

    consoleWindow = widgets.NewQMainWindow(nil, 0)
    consoleWindow.SetGeometry(core.NewQRect4(100, 100, 1000, 600))
    consoleWindow.SetWindowTitle("Gabu - GoBGP Console")
    var consoleWindowCentralWidget = widgets.NewQWidget(nil, 0)
    var consoleWindowCentralWidgetLayout = widgets.NewQHBoxLayout()
    consoleWindowCentralWidgetLayout.SetSpacing(6)
    consoleWindowCentralWidgetLayout.SetContentsMargins(11, 11, 11, 11)
    consoleWindowCentralWidget.SetLayout(consoleWindowCentralWidgetLayout)
    consoleWindow.SetCentralWidget(consoleWindowCentralWidget)

    // console window widgets
    // log Frame
    var logFrame = widgets.NewQFrame(consoleWindowCentralWidget, 0)
    logFrame.SetFrameShape(widgets.QFrame__Panel)
    logFrame.SetFrameShadow(widgets.QFrame__Raised)
    var frameSizePolicy = widgets.NewQSizePolicy()
    frameSizePolicy.SetHorizontalPolicy(widgets.QSizePolicy__Preferred)
    frameSizePolicy.SetVerticalPolicy(widgets.QSizePolicy__Preferred)
    frameSizePolicy.SetHorizontalStretch(0)
    frameSizePolicy.SetVerticalStretch(0)
    logFrame.SetSizePolicy(frameSizePolicy)

    // layout for log
    var logLayout = widgets.NewQVBoxLayout()
    logLayout.SetSpacing(6);


    // Console text edit / display
    var logLabel = widgets.NewQLabel2("Console output", logFrame, 0)
    var logLabelSizePolicy = widgets.NewQSizePolicy()
    logLabelSizePolicy.SetHorizontalPolicy(widgets.QSizePolicy__Preferred)
    logLabelSizePolicy.SetVerticalPolicy(widgets.QSizePolicy__Preferred)
    logLabelSizePolicy.SetHorizontalStretch(0)
    logLabelSizePolicy.SetVerticalStretch(0)
    logLabelSizePolicy.SetHeightForWidth(logLabel.HasHeightForWidth())
    logLabel.SetSizePolicy(logLabelSizePolicy)
    logLabel.SetAlignment(core.Qt__AlignLeading|core.Qt__AlignLeft|core.Qt__AlignVCenter)
    logLayout.AddWidget(logLabel, 0, core.Qt__AlignLeft)

    var logText = widgets.NewQTextEdit(logFrame)
    var fixeFont = gui.NewQFont2("monospace", 10, 0, false)
    logText.SetFont(fixeFont)
    var logTextSizePolicy = widgets.NewQSizePolicy()
    logTextSizePolicy.SetHorizontalPolicy(widgets.QSizePolicy__Expanding)
    logTextSizePolicy.SetVerticalPolicy(widgets.QSizePolicy__Expanding)
    logTextSizePolicy.SetHorizontalStretch(0)
    logTextSizePolicy.SetVerticalStretch(0)
    logTextSizePolicy.SetHeightForWidth(logText.HasHeightForWidth())
    logText.SetSizePolicy(logTextSizePolicy)
    logLayout.AddWidget(logText, 0, 0)

    logFrame.SetLayout(logLayout)
    consoleWindowCentralWidgetLayout.AddWidget(logFrame, 0, 0)

    // command Frame
    var cmdFrame = widgets.NewQFrame(consoleWindowCentralWidget, 0)
    cmdFrame.SetFrameShape(widgets.QFrame__Panel)
    cmdFrame.SetFrameShadow(widgets.QFrame__Raised)
    cmdFrame.SetSizePolicy(frameSizePolicy)

    // push buttons
    var cmdButtonSizePolicy = widgets.NewQSizePolicy()
    cmdButtonSizePolicy.SetHorizontalPolicy(widgets.QSizePolicy__Minimum)
    cmdButtonSizePolicy.SetVerticalPolicy(widgets.QSizePolicy__Fixed)
    cmdButtonSizePolicy.SetHorizontalStretch(0)
    cmdButtonSizePolicy.SetVerticalStretch(0)

    var cmdLabel = widgets.NewQLabel2("Basic Commands", cmdFrame, 0)
    cmdButtonSizePolicy.SetHeightForWidth(cmdLabel.HasHeightForWidth())
    cmdLabel.SetSizePolicy(cmdButtonSizePolicy)

    var    cmdNeighButton = widgets.NewQPushButton2("Neighbors", cmdFrame)
    cmdButtonSizePolicy.SetHeightForWidth(cmdNeighButton.HasHeightForWidth())
    cmdNeighButton.SetSizePolicy(cmdButtonSizePolicy)

    var cmdFsrib4Button = widgets.NewQPushButton2("IPv4 FlowSpec RIB", cmdFrame)
    cmdButtonSizePolicy.SetHeightForWidth(cmdFsrib4Button.HasHeightForWidth())
    cmdFsrib4Button.SetSizePolicy(cmdButtonSizePolicy)

    var cmdFsrib6Button = widgets.NewQPushButton2("IPv6 FlowSpec RIB", cmdFrame)
    cmdButtonSizePolicy.SetHeightForWidth(cmdFsrib6Button.HasHeightForWidth())
    cmdFsrib6Button.SetSizePolicy(cmdButtonSizePolicy)

    // layout for buttons
    var cmdLayout = widgets.NewQVBoxLayout()
    cmdLayout.AddWidget(cmdLabel, 0, core.Qt__AlignCenter)
    cmdLayout.AddWidget(cmdNeighButton, 0, 0)
    cmdLayout.AddWidget(cmdFsrib4Button, 0, 0)
    cmdLayout.AddWidget(cmdFsrib6Button, 0, 0)
    var cmdVerticalSpacer = widgets.NewQSpacerItem(20, 40, widgets.QSizePolicy__Fixed, widgets.QSizePolicy__Expanding)
    cmdLayout.AddItem(cmdVerticalSpacer)
    cmdFrame.SetLayout(cmdLayout)
    consoleWindowCentralWidgetLayout.AddWidget(cmdFrame, 0, 0)

    // Connect push buttons
    cmdNeighButton.ConnectClicked(func(_ bool) { cmdNeighButtonClicked(logText) })
    cmdFsrib4Button.ConnectClicked(func(_ bool) { cmdFsrib4ButtonClicked(logText) })
    cmdFsrib6Button.ConnectClicked(func(_ bool) { cmdFsrib6ButtonClicked(logText) })
    consoleWindow.ConnectCloseEvent(consoleWindowClosed)
    consoleWindow.Show()
}

func consoleWindowClosed(event *gui.QCloseEvent){
    windowBgpConsoleCreated = false
}

func cmdNeighButtonClicked(logTextWidget *widgets.QTextEdit) {
    dumpNeigh := bgpcli.TxtdumpGetNeighbor(client)

    for _, p := range dumpNeigh {
        logTextWidget.InsertPlainText(p)
    }
    logTextWidget.InsertPlainText("\n")
}

func cmdFsrib4ButtonClicked(logTextWidget *widgets.QTextEdit) {
    logTextWidget.InsertPlainText("Button FlowSpec 4\n\n")
}

func cmdFsrib6ButtonClicked(logTextWidget *widgets.QTextEdit) {
    logTextWidget.Append("Button FlowSpec 6\n\n")
}


func flowspecWin() {
    // Expanding Size policy
    var expandingSizePolicy = widgets.NewQSizePolicy()
    expandingSizePolicy.SetHorizontalPolicy(widgets.QSizePolicy__Expanding)
    expandingSizePolicy.SetVerticalPolicy(widgets.QSizePolicy__Expanding)
    expandingSizePolicy.SetHorizontalStretch(0)
    expandingSizePolicy.SetVerticalStretch(0)

    // preferred size policy
    var preferredSizePolicy = widgets.NewQSizePolicy()
    preferredSizePolicy.SetHorizontalPolicy(widgets.QSizePolicy__Preferred)
    preferredSizePolicy.SetVerticalPolicy(widgets.QSizePolicy__Preferred)
    preferredSizePolicy.SetHorizontalStretch(0)
    preferredSizePolicy.SetVerticalStretch(0)

    // Flowspec main window
    flowspecWindow = widgets.NewQMainWindow(nil, 0)
//    flowspecWindow.Layout().DestroyQObject()
    var flowspecCentralWid = widgets.NewQWidget(nil, 0)
    flowspecWindow.SetGeometry(core.NewQRect4(100, 100, 1000, 800))
    flowspecWindow.SetWindowTitle("Gabu - Flowspec tool")
    var flowspecWindowLayout = widgets.NewQVBoxLayout()
    flowspecWindowLayout.SetSpacing(6)
    flowspecWindowLayout.SetContentsMargins(11, 11, 11, 11)
    flowspecCentralWid.SetLayout(flowspecWindowLayout)

    // create one frame and a dock, frame to host flwospec rule config
    // and a dock to manage flowspec Rib towards GoBGP
    var editRuleFrame = widgets.NewQFrame(flowspecWindow, 0)

    editRuleFrame.SetSizePolicy(preferredSizePolicy)

    editRuleFrame.SetFrameShape(widgets.QFrame__Panel)
    editRuleFrame.SetFrameShadow(widgets.QFrame__Raised)
    flowspecWindowLayout.AddWidget(editRuleFrame, 0, 0)

    var editRuleFrameLayout = widgets.NewQHBoxLayout()
    editRuleFrame.SetLayout(editRuleFrameLayout)


    // Create content of editRuleFrame
    // Widget for Tree that displays library
    var editRuleLibWid = widgets.NewQWidget(editRuleFrame, 0)
    editRuleLibWid.SetSizePolicy(preferredSizePolicy)
    editRuleFrameLayout.AddWidget(editRuleLibWid, 0, 0)
    var editRuleLibWidLayout = widgets.NewQVBoxLayout()
    editRuleLibWid.SetLayout(editRuleLibWidLayout)
    var editRuleLabel = widgets.NewQLabel2("Rules Library", editRuleLibWid, 0)
    editRuleTree = widgets.NewQTreeWidget(editRuleLibWid)
    editRuleTree.SetSizePolicy(expandingSizePolicy)
    editRuleLibWidLayout.AddWidget(editRuleLabel, 0, 0)
    editRuleLibWidLayout.AddWidget(editRuleTree, 0, 0)
    editRuleTree.SetColumnCount(15)
    var editRuleTreeHeaderItem = editRuleTree.HeaderItem()
    libHeaderLabels := []string{"Add Family", "Dst Prefix", "Src Prefix", "Port", "Src Port", "Dst Port", "TCP flags",
"ICMP Type", "ICMP code", "Proto Num", "Pack Len", "DSCP", "IP Frag", "Action", "Act Value"}
    for i, myLabel := range libHeaderLabels {
        editRuleTreeHeaderItem.SetText(i, myLabel)
    }
    editRuleTree.Header().SetSectionResizeMode(widgets.QHeaderView__ResizeToContents)
    fullfilTreeWithRuleLib(editRuleTree, BgpFsActivLib)
    var editRuleLibButtonFrame = widgets.NewQFrame(editRuleLibWid, 0)
    editRuleLibButtonFrame.SetFrameShape(widgets.QFrame__Panel)
    editRuleLibButtonFrame.SetFrameShadow(widgets.QFrame__Raised)
    editRuleLibWidLayout.AddWidget(editRuleLibButtonFrame, 0, 0)
    var editRuleLibButtonFrameLayout = widgets.NewQGridLayout2()
    editRuleLibButtonFrame.SetLayout(editRuleLibButtonFrameLayout)
    var (
        editRuleLibSaveButton = widgets.NewQPushButton2("Save library", editRuleLibButtonFrame)
        editRuleLibLoadButton = widgets.NewQPushButton2("Load library", editRuleLibButtonFrame)
        editRuleLibPushRibButton = widgets.NewQPushButton2("Push rule to BGP Rib", editRuleLibButtonFrame)
    )
    editRuleLibButtonFrameLayout.AddWidget(editRuleLibLoadButton, 0, 0, 0)
    editRuleLibButtonFrameLayout.AddWidget(editRuleLibSaveButton, 0, 1, 0)
    editRuleLibButtonFrameLayout.AddWidget(editRuleLibPushRibButton, 0, 2, 0)
    // wire push button
    editRuleLibPushRibButton.ConnectClicked(func(_ bool) {editRuleLibPushRibButtonFunc()})
    editRuleLibSaveButton.ConnectClicked(func(_ bool) {editRuleLibSaveButtonFunc()})
    editRuleLibLoadButton.ConnectClicked(func(_ bool) {editRuleLibLoadButtonFunc(editRuleTree)})

    // Edit rule widget creation: it includes all required
    // UI Widget to edit a BGP flowspec rule
    var editRuleMainWid = widgets.NewQWidget(editRuleFrame, 0)
    editRuleMainWid.SetSizePolicy(preferredSizePolicy)
    editRuleFrameLayout.AddWidget(editRuleMainWid, 0, core.Qt__AlignLeft)
    var editRuleMainWidLayout = widgets.NewQVBoxLayout()
    editRuleMainWid.SetLayout(editRuleMainWidLayout)
    // Editing widets of Edit Match filter
    var editRuleMainWidLabel = widgets.NewQLabel2("Edit Flowspec Match Filter", editRuleMainWid, 0)
    editRuleMainWidLayout.AddWidget(editRuleMainWidLabel, 0, 0)

    // Line edit for source and dest prefix
    var editRulePrefixGroupBox = widgets.NewQGroupBox2("Address family and Prefix filters", editRuleMainWid)
    editRuleMainWidLayout.AddWidget(editRulePrefixGroupBox, 0, 0)
    var editRulePrefixLayout = widgets.NewQGridLayout2()
    editRulePrefixGroupBox.SetLayout(editRulePrefixLayout)
    var (
        editRuleSrcPrefixLabel = widgets.NewQLabel2("Source Prefix:", editRulePrefixGroupBox, 0)
        editRuleDstPrefixLabel = widgets.NewQLabel2("Destination Prefix:", editRulePrefixGroupBox, 0)
    )
    editAddrFamIpv4 = widgets.NewQRadioButton2("Flowspec IPv4", editRulePrefixGroupBox)
    editAddrFamIpv6 = widgets.NewQRadioButton2("Flowspec IPv6", editRulePrefixGroupBox)
    editRuleSrcPrefixLineEdit = widgets.NewQLineEdit(nil)
    editRuleDstPrefixLineEdit = widgets.NewQLineEdit(nil)
    editRuleSrcPrefixLineEdit.SetPlaceholderText("1.1.1.1/32")
    editRuleDstPrefixLineEdit.SetPlaceholderText("2.2.2.2/24")
    editAddrFamIpv4.SetChecked(true)
    AddrFamilyIpv6Checked = false
    AddrFamilyIpv4Checked = true
    regexpIpv6Validation = core.NewQRegExp2("^(([0-9a-fA-F]{1,4}:){7,7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:)|fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])|([0-9a-fA-F]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]))(\\/((1(1[0-9]|2[0-8]))|([0-9][0-9])|([0-9])))?$", core.Qt__CaseInsensitive, core.QRegExp__RegExp2)
    regexpIpv4Validation = core.NewQRegExp2("^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])(\\/([0-9]|[1-2][0-9]|3[0-2]))$", core.Qt__CaseInsensitive, core.QRegExp__RegExp2)
    regexpIpv4SrcValidator = gui.NewQRegExpValidator2(regexpIpv4Validation, editRuleSrcPrefixLineEdit)
    regexpIpv4DstValidator = gui.NewQRegExpValidator2(regexpIpv4Validation, editRuleDstPrefixLineEdit)
    editRuleSrcPrefixLineEdit.SetValidator(regexpIpv4SrcValidator)
    editRuleDstPrefixLineEdit.SetValidator(regexpIpv4DstValidator)
    editRulePrefixLayout.AddWidget(editRuleSrcPrefixLabel, 1, 0, 0)
    editRulePrefixLayout.AddWidget(editRuleSrcPrefixLineEdit, 1, 1, 0)
    editRulePrefixLayout.AddWidget(editAddrFamIpv4, 0, 2, 0)
    editRulePrefixLayout.AddWidget(editRuleDstPrefixLabel, 0, 0, 0)
    editRulePrefixLayout.AddWidget(editRuleDstPrefixLineEdit, 0, 1, 0)
    editRulePrefixLayout.AddWidget(editAddrFamIpv6, 1, 2, 0)
    // Wire address family radio button
    editAddrFamIpv4.ConnectClicked(editAddrFamIpv4Func)
    editAddrFamIpv6.ConnectClicked(editAddrFamIpv6Func)

    // horizontal widget to group together ICMP and proto type
    var editRuleIcmpProtoWid = widgets.NewQWidget(editRuleMainWid, 0)
    editRuleMainWidLayout.AddWidget(editRuleIcmpProtoWid, 0, 0)
    var editRuleIcmpProtoWidLayout = widgets.NewQHBoxLayout()
    editRuleIcmpProtoWidLayout.SetContentsMargins(0, 7, 0, 7)
    editRuleIcmpProtoWid.SetLayout(editRuleIcmpProtoWidLayout)
    // line edit for ICMP type and code
    var editRuleIcmpGroupBox = widgets.NewQGroupBox2("ICMP filters", editRuleMainWid)
    editRuleIcmpProtoWidLayout.AddWidget(editRuleIcmpGroupBox, 0, 0)
    var editRuleIcmpLayout = widgets.NewQGridLayout2()
    editRuleIcmpGroupBox.SetLayout(editRuleIcmpLayout)
    var (
        editRuleIcmpTypeLabel = widgets.NewQLabel2("ICMP Type:", editRuleIcmpGroupBox, 0)
        editRuleIcmpCodeLabel = widgets.NewQLabel2("ICMP Code:", editRuleIcmpGroupBox, 0)
    )
    editRuleIcmpTypeLineEdit = widgets.NewQLineEdit(nil)
    editRuleIcmpCodeLineEdit = widgets.NewQLineEdit(nil)
    editRuleIcmpTypeLineEdit.SetPlaceholderText("=0 =8")
    editRuleIcmpCodeLineEdit.SetPlaceholderText("=0")
    editRuleIcmpLayout.AddWidget(editRuleIcmpTypeLabel, 0, 0, 0)
    editRuleIcmpLayout.AddWidget(editRuleIcmpTypeLineEdit, 0, 1, 0)
    editRuleIcmpLayout.AddWidget(editRuleIcmpCodeLabel, 1, 0, 0)
    editRuleIcmpLayout.AddWidget(editRuleIcmpCodeLineEdit, 1, 1, 0)
    regexpByteValueValidation = core.NewQRegExp2("(([& ][<>]2[0-5][0-5]|[& ][<>]1[0-9][0-9]|[& ][<>][0-9][0-9]|[& ][<>][0-9])|([<>& ]=2[0-5][0-5]|[<>& ]=1[0-9][0-9]|[<>& ]=[0-9][0-9]|[<>& ]=[0-9])|([& ][<>]=2[0-5][0-5]|[<>& ]=1[0-9][0-9]|[<>& ]=[0-9][0-9]|[<>& ]=[0-9])|([& ][<>]=2[0-5][0-5]|[& ][<>]=1[0-9][0-9]|[& ][<>]=[0-9][0-9]|[& ][<>]=[0-9])|([<>=]2[0-5][0-5]|[<>=]1[0-9][0-9]|[<>=][0-9][0-9]|[<>=][0-9]))*", core.Qt__CaseInsensitive, core.QRegExp__RegExp2)
    regexpIcmpCodeValidator = gui.NewQRegExpValidator2(regexpByteValueValidation, editRuleIcmpCodeLineEdit)
    regexpIcmpTypeValidator = gui.NewQRegExpValidator2(regexpByteValueValidation, editRuleIcmpTypeLineEdit)
    editRuleIcmpTypeLineEdit.SetValidator(regexpIcmpTypeValidator)
    editRuleIcmpCodeLineEdit.SetValidator(regexpIcmpCodeValidator)
    // Line edit for IP protocol (Next header)
    var editRuleIpProtoGroupBox = widgets.NewQGroupBox2("IP protocol or Next header", editRuleMainWid)
    editRuleIcmpProtoWidLayout.AddWidget(editRuleIpProtoGroupBox, 0, 0)
    var editRuleIpProtoLayout = widgets.NewQGridLayout2()
    editRuleIpProtoGroupBox.SetLayout(editRuleIpProtoLayout)
    var (
        editRuleIpProtoLabel = widgets.NewQLabel2("Protocol number:", editRuleIcmpGroupBox, 0)
    )
    editRuleIpProtoLineEdit = widgets.NewQLineEdit(nil)
    editRuleIpProtoLineEdit.SetPlaceholderText("=6 =17")
    editRuleIpProtoLayout.AddWidget(editRuleIpProtoLabel, 0, 0, 0)
    editRuleIpProtoLayout.AddWidget(editRuleIpProtoLineEdit, 0, 1, 0)
    regexpProtoNumValidator = gui.NewQRegExpValidator2(regexpByteValueValidation, editRuleIpProtoLineEdit)
    editRuleIpProtoLineEdit.SetValidator(regexpProtoNumValidator)
    // widget and layout for both TCP/UDP ports and DSCP and Packet lenght
    var editRulePortDscpPackLenghtWid = widgets.NewQWidget(editRuleMainWid, 0)
    editRuleMainWidLayout.AddWidget(editRulePortDscpPackLenghtWid, 0, 0)
    var editRulePortDscpPackLenghtWidLayout = widgets.NewQHBoxLayout()
    editRulePortDscpPackLenghtWid.SetLayout(editRulePortDscpPackLenghtWidLayout)

    // line edit for TCP/UDP ports
    var editRulePortGroupBox = widgets.NewQGroupBox2("Port filters", editRulePortDscpPackLenghtWid)
    editRulePortDscpPackLenghtWidLayout.AddWidget(editRulePortGroupBox, 0, 0)
    var editRulePortLayout = widgets.NewQGridLayout2()
    editRulePortGroupBox.SetLayout(editRulePortLayout)
    var (
        editRulePortLabel = widgets.NewQLabel2("Port:", editRulePortGroupBox, 0)
        editRuleSrcPortLabel = widgets.NewQLabel2("Src Port:", editRulePortGroupBox, 0)
        editRuleDstPortLabel = widgets.NewQLabel2("Dst Port:", editRulePortGroupBox, 0)
    )
    editRulePortLineEdit = widgets.NewQLineEdit(nil)
    editRuleSrcPortLineEdit = widgets.NewQLineEdit(nil)
    editRuleDstPortLineEdit = widgets.NewQLineEdit(nil)
    editRulePortLineEdit.SetPlaceholderText("=80 >=8080&<=8888")
    editRuleSrcPortLineEdit.SetPlaceholderText("=443&=80")
    editRuleDstPortLineEdit.SetPlaceholderText(">=1024&<=49151")
    regexpPortValidation = core.NewQRegExp2("(([<>=]6[0-5][0-5][0-3][0-5]|[<>=]65[1-4][0-9][0-9]|[<>=]655[1-2][0-9]|[<>=]64[0-9][0-9][0-9]|[<>=][1-9][0-9][0-9][0-9]|[<>=][1-9][0-9][0-9]|[<>=][1-9][0-9]|[<>=][1-9])|([& ][<>]6[0-5][0-5][0-3][0-5]|[& ][<>]65[1-4][0-9][0-9]|[& ][<>]655[1-2][0-9]|[& ][<>]64[0-9][0-9][0-9]|[& ][<>][1-9][0-9][0-9][0-9]|[& ][<>][1-9][0-9][0-9]|[& ][<>][1-9][0-9]|[& ][<>][1-9])|([<>& ]=6[0-5][0-5][0-3][0-5]|[<>& ]=65[1-4][0-9][0-9]|[<>& ]=655[1-2][0-9]|[<>& ]=64[0-9][0-9][0-9]|[<>& ]=[1-9][0-9][0-9][0-9]|[<>& ]=[1-9][0-9][0-9]|[<>& ]=[1-9][0-9]|[<>& ]=[1-9])|([& ][<>]=6[0-5][0-5][0-3][0-5]|[<>& ]=65[1-4][0-9][0-9]|[<>& ]=655[1-2][0-9]|[<>& ]=64[0-9][0-9][0-9]|[<>& ]=[1-9][0-9][0-9][0-9]|[<>& ]=[1-9][0-9][0-9]|[<>& ]=[1-9][0-9]|[<>& ]=[1-9])|([& ][<>]=6[0-5][0-5][0-3][0-5]|[& ][<>]=65[1-4][0-9][0-9]|[& ][<>]=655[1-2][0-9]|[& ][<>]=64[0-9][0-9][0-9]|[& ][<>]=[1-9][0-9][0-9][0-9]|[& ][<>]=[1-9][0-9][0-9]|[& ][<>]=[1-9][0-9]|[& ][<>]=[1-9]))*", core.Qt__CaseInsensitive, core.QRegExp__RegExp2)
    regexpPortValidator = gui.NewQRegExpValidator2(regexpPortValidation, editRulePortLineEdit)
    regexpSrcPortValidator = gui.NewQRegExpValidator2(regexpPortValidation, editRuleSrcPortLineEdit)
    regexpDstPortValidator = gui.NewQRegExpValidator2(regexpPortValidation, editRuleDstPrefixLineEdit)
    editRulePortLineEdit.SetValidator(regexpPortValidator)
    editRuleSrcPortLineEdit.SetValidator(regexpSrcPortValidator)
    editRuleDstPortLineEdit.SetValidator(regexpDstPortValidator)
    editRulePortLayout.AddWidget(editRulePortLabel, 0, 0, 0)
    editRulePortLayout.AddWidget(editRulePortLineEdit, 0, 1, 0)
    editRulePortLayout.AddWidget(editRuleSrcPortLabel, 1, 0, 0)
    editRulePortLayout.AddWidget(editRuleSrcPortLineEdit, 1, 1, 0)
    editRulePortLayout.AddWidget(editRuleDstPortLabel, 2, 0, 0)
    editRulePortLayout.AddWidget(editRuleDstPortLineEdit, 2, 1, 0)

    // Line edit for packet length and DSCP
    var editRuleLenDscpGroupBox = widgets.NewQGroupBox2("Packet Length and DSCP", editRulePortDscpPackLenghtWid)
    editRulePortDscpPackLenghtWidLayout.AddWidget(editRuleLenDscpGroupBox, 0, 0)
    var editRuleLenDscpLayout = widgets.NewQGridLayout2()
    editRuleLenDscpGroupBox.SetLayout(editRuleLenDscpLayout)
    var (
        editRuleLenLabel = widgets.NewQLabel2("Packet length:", editRuleLenDscpGroupBox, 0)
        editRuleDscpLabel = widgets.NewQLabel2("DiffServ Codepoints:", editRuleLenDscpGroupBox, 0)

    )
    editRuleLenLineEdit = widgets.NewQLineEdit(nil)
    editRuleDscpLineEdit = widgets.NewQLineEdit(nil)
    editRuleLenLineEdit.SetPlaceholderText(">=64&<=1024")
    editRuleDscpLineEdit.SetPlaceholderText("=46")
    editRuleLenDscpLayout.AddWidget(editRuleLenLabel, 0, 0, 0)
    editRuleLenDscpLayout.AddWidget(editRuleLenLineEdit, 0, 1, 0)
    editRuleLenDscpLayout.AddWidget(editRuleDscpLabel, 1, 0, 0)
    editRuleLenDscpLayout.AddWidget(editRuleDscpLineEdit, 1, 1, 0)
    regexpDscpValidator = gui.NewQRegExpValidator2(regexpByteValueValidation, editRuleDscpLineEdit)
    editRuleDscpLineEdit.SetValidator(regexpDscpValidator)
    regexpPckLenValidation = core.NewQRegExp2("(([& ][<>]9000|[& ][<>]8[0-9][0-9][0-9]|[& ][<>][0-9][0-9][0-9]|[& ][<>][0-9][0-9]|[& ][<>][0-9])|([<>& ]=9000|[<>& ]=8[0-9][0-9][0-9]|[<>& ]=[0-9][0-9][0-9]|[<>& ]=[0-9][0-9]|[<>& ]=[0-9])|([& ][<>]=9000|[& ][<>]=8[0-9][0-9][0-9]|[& ][<>]=[0-9][0-9][0-9]|[& ][<>]=[0-9][0-9]|[& ][<>]=[0-9])|([& ][<>]=9000|[& ][<>]=8[0-9][0-9][0-9]|[& ][<>]=[0-9][0-9][0-9]|[& ][<>]=[0-9][0-9]|[& ][<>]=[0-9])|([<>=]9000|[<>=]8[0-9][0-9][0-9]|[<>=][0-9][0-9][0-9]|[<>=][0-9][0-9]|[<>=][0-9]))*", core.Qt__CaseInsensitive, core.QRegExp__RegExp2)
    regexpPckLenValidator = gui.NewQRegExpValidator2(regexpPckLenValidation, editRuleLenLineEdit)
    editRuleLenLineEdit.SetValidator(regexpPckLenValidator)
    // line edit for TCP flags
    var editRuleTcpFlagGroupBox = widgets.NewQGroupBox2("TCP flags filter", editRuleMainWid)
    editRuleMainWidLayout.AddWidget(editRuleTcpFlagGroupBox, 0, 0)
    var editRuleTcpFlagLayout = widgets.NewQGridLayout2()
    editRuleTcpFlagGroupBox.SetLayout(editRuleTcpFlagLayout)
    var (
        editRuleTcpSynFlagCheck = widgets.NewQCheckBox2("SYN", editRuleTcpFlagGroupBox)
        editRuleTcpAckFlagCheck = widgets.NewQCheckBox2("ACK", editRuleTcpFlagGroupBox)
        editRuleTcpRstFlagCheck = widgets.NewQCheckBox2("RST", editRuleTcpFlagGroupBox)
        editRuleTcpFinFlagCheck = widgets.NewQCheckBox2("FIN", editRuleTcpFlagGroupBox)
        editRuleTcpPshFlagCheck = widgets.NewQCheckBox2("PSH", editRuleTcpFlagGroupBox)
        editRuleTcpEceFlagCheck = widgets.NewQCheckBox2("ECE", editRuleTcpFlagGroupBox)
        editRuleTcpUrgFlagCheck = widgets.NewQCheckBox2("URG", editRuleTcpFlagGroupBox)
        editRuleTcpCwrFlagCheck = widgets.NewQCheckBox2("CWR", editRuleTcpFlagGroupBox)
        editRuleLineSeparator = widgets.NewQFrame(editRuleTcpFlagGroupBox, 0)
        editRuleTcpOpAndCheck = widgets.NewQCheckBox2("AND", editRuleTcpFlagGroupBox)
        editRuleTcpOpNotCheck = widgets.NewQCheckBox2("NOT", editRuleTcpFlagGroupBox)
        editRuleTcpOpMatchCheck = widgets.NewQCheckBox2("MATCH", editRuleTcpFlagGroupBox)
        editRuleTcpFlagFilterLabel = widgets.NewQLabel2("Filter:", editRuleTcpFlagGroupBox, 0)
        editRuleTcpFlagAddButton = widgets.NewQPushButton2("Add", editRuleTcpFlagGroupBox)
    )
    editRuleTcpFlagFilterLine = widgets.NewQLineEdit(nil)
    editRuleLineSeparator.SetFrameShape(widgets.QFrame__VLine)
    editRuleLineSeparator.SetFrameShadow(widgets.QFrame__Sunken)
    editRuleTcpFlagLayout.AddWidget(editRuleTcpSynFlagCheck, 0, 0, 0)
    editRuleTcpFlagLayout.AddWidget(editRuleTcpAckFlagCheck, 0, 1, 0)
    editRuleTcpFlagLayout.AddWidget(editRuleTcpRstFlagCheck, 0, 2, 0)
    editRuleTcpFlagLayout.AddWidget(editRuleTcpFinFlagCheck, 0, 3, 0)
    editRuleTcpFlagLayout.AddWidget(editRuleTcpFlagFilterLabel, 0, 5, 0)
    editRuleTcpFlagLayout.AddWidget3(editRuleTcpFlagFilterLine, 0, 6, 1, 3, 0)
    editRuleTcpFlagLayout.AddWidget(editRuleTcpPshFlagCheck, 1, 0, 0)
    editRuleTcpFlagLayout.AddWidget(editRuleTcpEceFlagCheck, 1, 1, 0)
    editRuleTcpFlagLayout.AddWidget(editRuleTcpUrgFlagCheck, 1, 2, 0)
    editRuleTcpFlagLayout.AddWidget(editRuleTcpCwrFlagCheck, 1, 3, 0)
    editRuleTcpFlagLayout.AddWidget3(editRuleLineSeparator, 0, 4, 2, 1, 0)
    editRuleTcpFlagLayout.AddWidget(editRuleTcpOpAndCheck, 1, 5, 0)
    editRuleTcpFlagLayout.AddWidget(editRuleTcpOpNotCheck, 1, 6, 0)
    editRuleTcpFlagLayout.AddWidget(editRuleTcpOpMatchCheck, 1, 7, 0)
    editRuleTcpFlagLayout.AddWidget(editRuleTcpFlagAddButton, 1, 8, 0)

    // Line edit and checkbox for fragment filtering
    var editRuleFragGroupBox = widgets.NewQGroupBox2("IP Fragment", editRuleMainWid)
    editRuleMainWidLayout.AddWidget(editRuleFragGroupBox, 0, 0)
    var editRuleFragLayout = widgets.NewQGridLayout2()
    editRuleFragGroupBox.SetLayout(editRuleFragLayout)
    var (
        editRuleIsfFragCheck = widgets.NewQCheckBox2("IsF", editRuleFragGroupBox)
        editRuleFfFragCheck = widgets.NewQCheckBox2("FF", editRuleFragGroupBox)
        editRuleLfFragCheck = widgets.NewQCheckBox2("LF", editRuleFragGroupBox)
        editRuleDfFragCheck = widgets.NewQCheckBox2("DF", editRuleFragGroupBox)
        editRuleAndFragCheck = widgets.NewQCheckBox2("AND", editRuleFragGroupBox)
        editRuleNotFragCheck = widgets.NewQCheckBox2("NOT", editRuleFragGroupBox)
        editRuleMatchFragCheck = widgets.NewQCheckBox2("MATCH", editRuleFragGroupBox)
        editRuleLineFragSeparator = widgets.NewQFrame(editRuleFragGroupBox, 0)
        editRuleFragFilterLabel = widgets.NewQLabel2("Filter:", editRuleFragGroupBox, 0)
        editRuleAddFragButton = widgets.NewQPushButton2("Add", editRuleFragGroupBox)
    )
    editRuleFragFilterLine = widgets.NewQLineEdit(nil)
    editRuleLineFragSeparator.SetFrameShape(widgets.QFrame__VLine)
    editRuleLineFragSeparator.SetFrameShadow(widgets.QFrame__Sunken)
    editRuleFragLayout.AddWidget(editRuleIsfFragCheck, 0, 0, 0)
    editRuleFragLayout.AddWidget(editRuleFfFragCheck, 0, 1, 0)
    editRuleFragLayout.AddWidget(editRuleLfFragCheck, 1, 0, 0)
    editRuleFragLayout.AddWidget(editRuleDfFragCheck, 1, 1, 0)
    editRuleFragLayout.AddWidget3(editRuleLineFragSeparator, 0, 2, 2, 1, 0)
    editRuleFragLayout.AddWidget(editRuleAndFragCheck, 1, 3, 0)
    editRuleFragLayout.AddWidget(editRuleNotFragCheck, 1, 4, 0)
    editRuleFragLayout.AddWidget(editRuleMatchFragCheck, 1, 5, 0)
    editRuleFragLayout.AddWidget(editRuleAddFragButton, 1, 6, 0)
    editRuleFragLayout.AddWidget3(editRuleFragFilterLabel, 0, 3, 1, 1, 0)
    editRuleFragLayout.AddWidget3(editRuleFragFilterLine, 0, 4, 1, 3, 0)

    // Editing widets of Action applied to match traffic
    var editRuleMainWidLabelMatch = widgets.NewQLabel2("Edit Flowspec Action", editRuleMainWid, 0)
    editRuleMainWidLayout.AddWidget(editRuleMainWidLabelMatch, 0, 0)
    // Match group box widget
    var editRuleActionGroupBox = widgets.NewQGroupBox2("Action applied", editRuleMainWid)
    editRuleMainWidLayout.AddWidget(editRuleActionGroupBox, 0, 0)
    var editRuleActionLayout = widgets.NewQGridLayout2()
    editRuleActionGroupBox.SetLayout(editRuleActionLayout)
    editRuleActionCombo = widgets.NewQComboBox(nil)
    editRuleActSisterValueLine = widgets.NewQLineEdit(nil)
    editRuleActionCombo.AddItems(flowSpecActStrings)
    editRuleActionLayout.AddWidget(editRuleActionCombo, 0, 0, 0)
    editRuleActionLayout.AddWidget(editRuleActSisterValueLine, 1, 0, 0)
    // wire combo button
    editRuleActionCombo.ConnectCurrentIndexChanged(func(index int) { editRuleActionComboFunc(editRuleActSisterValueLine, index) })
    // global apply button
    var editRuleGlobButtonFrame = widgets.NewQFrame(editRuleMainWid, 0)
    var editRuleGlobButtonlayout = widgets.NewQGridLayout2()
    editRuleGlobButtonFrame.SetLayout(editRuleGlobButtonlayout)
    editRuleMainWidLayout.AddWidget(editRuleGlobButtonFrame, 0, 0)
    var (
        editGlobButtonNew = widgets.NewQPushButton2("New", editRuleGlobButtonFrame)
        editGlobButtonApply = widgets.NewQPushButton2("Apply", editRuleGlobButtonFrame)
        editGlobButtonReset = widgets.NewQPushButton2("Reset", editRuleGlobButtonFrame)
        editGlobButtonDelete = widgets.NewQPushButton2("Delete", editRuleGlobButtonFrame)
        editGlobButtonDuplicate = widgets.NewQPushButton2("Duplicate", editRuleGlobButtonFrame)

    )
    editRuleGlobButtonFrame.SetFrameShape(widgets.QFrame__Panel)
    editRuleGlobButtonFrame.SetFrameShadow(widgets.QFrame__Raised)
    editRuleGlobButtonlayout.AddWidget(editGlobButtonNew, 0, 0, 0)
    editRuleGlobButtonlayout.AddWidget(editGlobButtonApply, 0, 1, 0)
    editRuleGlobButtonlayout.AddWidget(editGlobButtonReset, 0, 2, 0)
    editRuleGlobButtonlayout.AddWidget(editGlobButtonDelete, 0, 3, 0)
    editRuleGlobButtonlayout.AddWidget(editGlobButtonDuplicate, 0, 4, 0)

    // var editRuleMainWidSpacer = widgets.NewQSpacerItem(20, 40, widgets.QSizePolicy__Fixed, widgets.QSizePolicy__Expanding)
    // editRuleMainWidLayout.AddItem(editRuleMainWidSpacer)
    // Connection of all widget to QT slots and actions
    // Tree Widget
    editRuleTree.ConnectItemClicked(editRuleLibItemSelected)
    // push button from rule edition
    // Connect push buttons
    editGlobButtonApply.ConnectClicked(func(_ bool) { editGlobButtonApplyFunc() })
    editGlobButtonNew.ConnectClicked(func(_ bool) { editGlobButtonNewFunc() })
    editGlobButtonDelete.ConnectClicked(func(_ bool) { editGlobButtonDeleteFunc() })
    editGlobButtonReset.ConnectClicked(func(_ bool) { editGlobButtonResetFunc() })
    editGlobButtonDuplicate.ConnectClicked(func(_ bool) { editGlobButtonDuplicateFunc() })

    // widget of the Rib tool dock
    var ribManipDock = widgets.NewQDockWidget("FlowSpec RIB tool", flowspecWindow, 0)
    // ribManipDock.SetSizePolicy(preferredSizePolicy)
    // flowspecWindowLayout.AddWidget(ribManipDock, 0, 0)
    flowspecWindow.AddDockWidget(core.Qt__BottomDockWidgetArea, ribManipDock)
    // main widget
    var ribManipDockWid = widgets.NewQWidget(nil, 0)
    var ribManipDockWidLayout = widgets.NewQHBoxLayout()
    ribManipDockWid.SetLayout(ribManipDockWidLayout)
    // Tree displaying BGP FS RIB
    var ribContentTree = widgets.NewQTreeWidget(ribManipDockWid)
    ribContentTree.SetSizePolicy(expandingSizePolicy)
    ribManipDockWidLayout.AddWidget(ribContentTree, 0, 0)
    ribContentTree.SetColumnCount(4)
    ribContentTree.Header().SetSectionResizeMode(widgets.QHeaderView__ResizeToContents)
    var ribContentTreeHeaderItem = ribContentTree.HeaderItem()
    ribHeaderLabels := []string{"NLRI", "Extended community", "Age", "Next Hop"}
    for i, myLabel := range ribHeaderLabels {
        ribContentTreeHeaderItem.SetText(i, myLabel)
    }
    // Buttons for rib manip
    var ribManipButtonWid = widgets.NewQWidget(ribManipDockWid, 0)
    var ribManipButtonWidLayout = widgets.NewQVBoxLayout()
    ribManipButtonWid.SetLayout(ribManipButtonWidLayout)
    var (
        ribManipLoadButton = widgets.NewQPushButton2("Load/Reload BGP FS RIB", ribManipButtonWid)
        ribManipDeleteRuleButton = widgets.NewQPushButton2("Delete rule from RIB", ribManipButtonWid)
    )
    ribAddrFamIpv4 = widgets.NewQRadioButton2("RIB IPv4", ribManipButtonWid)
    ribAddrFamIpv4.SetChecked(true)
    ribAddrFamIpv6 = widgets.NewQRadioButton2("RIB IPv6", ribManipButtonWid)
    ribManipButtonWidLayout.AddWidget(ribManipLoadButton, 0, 0)
    ribManipButtonWidLayout.AddWidget(ribManipDeleteRuleButton, 0, 0)
    ribManipButtonWidLayout.AddWidget(ribAddrFamIpv4, 0, 0)
    ribManipButtonWidLayout.AddWidget(ribAddrFamIpv6, 0, 0)
    var ribManipWidSpacer = widgets.NewQSpacerItem(20, 40, widgets.QSizePolicy__Fixed, widgets.QSizePolicy__Expanding)
    ribManipButtonWidLayout.AddItem(ribManipWidSpacer)

    ribManipDockWidLayout.AddWidget(ribManipButtonWid, 0, 0)
    // Wire address family radio button
    ribAddrFamIpv4.ConnectClicked(ribAddrFamIpv4Func)
    ribAddrFamIpv6.ConnectClicked(ribAddrFamIpv6Func)

    // wire load RIB and delete RIB buttons of the dock
    ribManipLoadButton.ConnectClicked(func(_ bool) {ribManipLoadRibFunc(ribContentTree)})
    ribManipDeleteRuleButton.ConnectClicked(func(_ bool) {ribManipDeleteRuleButtonFunc(ribContentTree)})

    ribManipDock.SetWidget(ribManipDockWid)
    ribManipDock.SetFeatures(widgets.QDockWidget__DockWidgetFloatable | widgets.QDockWidget__DockWidgetMovable)

    flowspecWindow.SetCentralWidget(flowspecCentralWid)
    flowspecWindow.ConnectCloseEvent(func (myCloseEvent *gui.QCloseEvent) {flowspecWindowClosed(myCloseEvent, ribManipDock)})

    flowspecWindow.Show()
}

// function called with IPv4 or IPv6 radiobutton from edit widget is checked

func editAddrFamIpv4Func(checked bool) {
    AddrFamilyIpv4Checked = true
    AddrFamilyIpv6Checked = false
    editRuleSrcPrefixLineEdit.SetPlaceholderText("1.1.1.1/32")
    editRuleDstPrefixLineEdit.SetPlaceholderText("2.2.2.2/24")
    regexpIpv4SrcValidator = gui.NewQRegExpValidator2(regexpIpv4Validation, editRuleSrcPrefixLineEdit)
    regexpIpv4DstValidator = gui.NewQRegExpValidator2(regexpIpv4Validation, editRuleDstPrefixLineEdit)
    editRuleSrcPrefixLineEdit.SetValidator(regexpIpv4SrcValidator)
    editRuleDstPrefixLineEdit.SetValidator(regexpIpv4DstValidator)
}

func editAddrFamIpv6Func(checked bool) {
    AddrFamilyIpv6Checked = true
    AddrFamilyIpv4Checked = false
    editRuleSrcPrefixLineEdit.SetText("")
    editRuleDstPrefixLineEdit.SetText("")
    editRuleSrcPrefixLineEdit.SetPlaceholderText("2001:DB8::/32")
    editRuleDstPrefixLineEdit.SetPlaceholderText("2001:DB8::1/128")
    regexpIpv6SrcValidator = gui.NewQRegExpValidator2(regexpIpv6Validation, editRuleSrcPrefixLineEdit)
    regexpIpv6DstValidator = gui.NewQRegExpValidator2(regexpIpv6Validation, editRuleDstPrefixLineEdit)
    editRuleSrcPrefixLineEdit.SetValidator(regexpIpv6SrcValidator)
    editRuleDstPrefixLineEdit.SetValidator(regexpIpv6DstValidator)
}

// func called when IPv4 or IPv6 RIB radio button are checked

func ribAddrFamIpv4Func(checked bool) {
    ribRadioFamilychecked = "ipv4-flowspec"
}

func ribAddrFamIpv6Func(checked bool) {
    ribRadioFamilychecked = "ipv6-flowspec"
}

// function called when load rib button clicked

func ribManipLoadRibFunc(myTree *widgets.QTreeWidget) {
    cleanupTree(myTree)
    bgpcli.FlowSpecRibFulfillTree(client, myTree, ribRadioFamilychecked)
    ribActiveFamily = ribRadioFamilychecked
}

func ribManipDeleteRuleButtonFunc(myTree *widgets.QTreeWidget) {
    var myItem *widgets.QTreeWidgetItem = nil
    var myNlri string
    var myExtCom string
    var deleteCmdPath string
    myItem = myTree.CurrentItem()
    index := myTree.IndexOfTopLevelItem(myItem)
    myNlri = myItem.Text(0)
    myExtCom = myItem.Text(1)
    deleteCmdPath = fmt.Sprintf("match %sthen %s", formatNlriOutputToDeleteCmdNlri(myNlri), formatExtComOutputToDeleteCmdExtCom(myExtCom))

    if (index == -1) {
        return
    }
    bgpcli.DeleteFlowSpecPath(client, deleteCmdPath, ribActiveFamily)
    if(index >= 0 && index < myTree.TopLevelItemCount()) {
        myItem = myTree.TakeTopLevelItem(index)
    }
    return
}

// function called when window get closed

func flowspecWindowClosed(event *gui.QCloseEvent, myDock *widgets.QDockWidget){
    windowFlowSpecCreated = false
    flowspecWindow.RemoveDockWidget(myDock)
}

// Copy the content of a flowspec rule structure into a TreeItem widget

func createFullfilItemWithRule(ty int, myTree *widgets.QTreeWidget, myRule BgpFsRule) {
    var myItem = widgets.NewQTreeWidgetItem3(myTree, ty)
    myItem.SetText(0, myRule.AddrFam)
    myItem.SetText(1, myRule.DstPrefix)
    myItem.SetText(2, myRule.SrcPrefix)
    myItem.SetText(3, myRule.Port)
    myItem.SetText(4, myRule.SrcPort)
    myItem.SetText(5, myRule.DstPort)
    myItem.SetText(6, myRule.TcpFlags)
    myItem.SetText(7, myRule.IcmpType)
    myItem.SetText(8, myRule.IcmpCode)
    myItem.SetText(9, myRule.ProtoNumber)
    myItem.SetText(10, myRule.PacketLen)
    myItem.SetText(11, myRule.Dscp)
    myItem.SetText(12, myRule.IpFrag)
    myItem.SetText(13, myRule.Action)
    myItem.SetText(14, myRule.ActSisterValue)
}

func fullfilItemWithRule(ty int, myItem *widgets.QTreeWidgetItem, myRule BgpFsRule) {
    myItem.SetText(0, myRule.AddrFam)
    myItem.SetText(1, myRule.DstPrefix)
    myItem.SetText(2, myRule.SrcPrefix)
    myItem.SetText(3, myRule.Port)
    myItem.SetText(4, myRule.SrcPort)
    myItem.SetText(5, myRule.DstPort)
    myItem.SetText(6, myRule.TcpFlags)
    myItem.SetText(7, myRule.IcmpType)
    myItem.SetText(8, myRule.IcmpCode)
    myItem.SetText(9, myRule.ProtoNumber)
    myItem.SetText(10, myRule.PacketLen)
    myItem.SetText(11, myRule.Dscp)
    myItem.SetText(12, myRule.IpFrag)
    myItem.SetText(13, myRule.Action)
    myItem.SetText(14, myRule.ActSisterValue)
}

func fullfilTreeWithRuleLib(myTree *widgets.QTreeWidget, myRuleLib []BgpFsRule) {
    for i, myRule := range myRuleLib {
        createFullfilItemWithRule(i, myTree, myRule)
    }
}

func cleanupTree(myTree *widgets.QTreeWidget) {
    if (myTree.TopLevelItemCount() != 0) {
        maxItem := myTree.TopLevelItemCount()
        for i := (maxItem-1); i >= 0; i-- {
            myTree.TakeTopLevelItem(i)
        }
    }
}

func fullfilLineEditWithBgpFs(myRule BgpFsRule) {
    if(myRule.AddrFam == "IPv4") {
        editAddrFamIpv4.SetChecked(true)
        editAddrFamIpv6.SetChecked(false)
    } else if(myRule.AddrFam == "IPv6") {
        editAddrFamIpv4.SetChecked(false)
        editAddrFamIpv6.SetChecked(true)
    }
    editRuleSrcPrefixLineEdit.SetText(myRule.SrcPrefix)
    editRuleDstPrefixLineEdit.SetText(myRule.DstPrefix)
    editRuleIcmpTypeLineEdit.SetText(myRule.IcmpType)
    editRuleIcmpCodeLineEdit.SetText(myRule.IcmpCode)
    editRuleIpProtoLineEdit.SetText(myRule.ProtoNumber)
    editRulePortLineEdit.SetText(myRule.Port)
    editRuleSrcPortLineEdit.SetText(myRule.SrcPort)
    editRuleDstPortLineEdit.SetText(myRule.DstPort)
    editRuleTcpFlagFilterLine.SetText(myRule.TcpFlags)
    editRuleLenLineEdit.SetText(myRule.PacketLen)
    editRuleDscpLineEdit.SetText(myRule.Dscp)
    editRuleFragFilterLine.SetText(myRule.IpFrag)
    editRuleActionCombo.SetCurrentIndex(flowSpecActValueMap[myRule.Action])
    editRuleActSisterValueLine.SetText(myRule.ActSisterValue)
}

func fullfilBgpFsWithLineEdit(myIndex int) {
    if (AddrFamilyIpv4Checked) {
        BgpFsActivLib[myIndex].AddrFam = "IPv4"
    } else if (AddrFamilyIpv6Checked) {
        BgpFsActivLib[myIndex].AddrFam = "IPv6"
    }
    BgpFsActivLib[myIndex].SrcPrefix = editRuleSrcPrefixLineEdit.Text()
    BgpFsActivLib[myIndex].DstPrefix =  editRuleDstPrefixLineEdit.Text()
    BgpFsActivLib[myIndex].IcmpType = editRuleIcmpTypeLineEdit.Text()
    BgpFsActivLib[myIndex].IcmpCode = editRuleIcmpCodeLineEdit.Text()
    BgpFsActivLib[myIndex].ProtoNumber = editRuleIpProtoLineEdit.Text()
    BgpFsActivLib[myIndex].Port = editRulePortLineEdit.Text()
    BgpFsActivLib[myIndex].SrcPort = editRuleSrcPortLineEdit.Text()
    BgpFsActivLib[myIndex].DstPort = editRuleDstPortLineEdit.Text()
    BgpFsActivLib[myIndex].TcpFlags = editRuleTcpFlagFilterLine.Text()
    BgpFsActivLib[myIndex].PacketLen = editRuleLenLineEdit.Text()
    BgpFsActivLib[myIndex].Dscp = editRuleDscpLineEdit.Text()
    BgpFsActivLib[myIndex].IpFrag = editRuleFragFilterLine.Text()
    BgpFsActivLib[myIndex].Action = flowSpecActStrings[editRuleActionCombo.CurrentIndex()]
    if (flowSpecActStrings[editRuleActionCombo.CurrentIndex()] == "Drop") {
        BgpFsActivLib[myIndex].ActSisterValue = ""
    } else {
        BgpFsActivLib[myIndex].ActSisterValue = editRuleActSisterValueLine.Text()
    }
}

// fucntion when an lib item is clicked

func editRuleLibItemSelected(myItem *widgets.QTreeWidgetItem, column int) {
    index := editRuleTree.IndexOfTopLevelItem(myItem)
    fullfilLineEditWithBgpFs(BgpFsActivLib[index])
}

// function to manage glob push button

func editGlobButtonNewFunc() {
    var myFsRule BgpFsRule
    myFsRule.DstPrefix = "New"
    myFsRule.AddrFam = "IPv4"
    myFsRule.Action = "Drop"
    BgpFsActivLib = append(BgpFsActivLib, myFsRule)
    createFullfilItemWithRule(len(BgpFsActivLib)-1, editRuleTree, BgpFsActivLib[len(BgpFsActivLib)-1])
}

func editGlobButtonApplyFunc() {
    var myItem *widgets.QTreeWidgetItem
    myItem = editRuleTree.CurrentItem()
    if (*myItem != (widgets.QTreeWidgetItem{})) {
        index := editRuleTree.IndexOfTopLevelItem(myItem)
        fullfilBgpFsWithLineEdit(index)
        fullfilItemWithRule(index, myItem, BgpFsActivLib[index])
    } else {
        warningMessage := widgets.NewQMessageBox2(widgets.QMessageBox__Warning, "Rule library issue", "Please select a rule to apply changes", widgets.QMessageBox__Ok, flowspecWindow, core.Qt__Window)
        warningMessage.Exec()
    }
}

func editGlobButtonDeleteFunc() {
    var myItem *widgets.QTreeWidgetItem
    myItem = editRuleTree.CurrentItem()
    index := editRuleTree.IndexOfTopLevelItem(myItem)
     if(index >= 0 && index < editRuleTree.TopLevelItemCount()) {
        myItem = editRuleTree.TakeTopLevelItem(index)
     }
    BgpFsActivLib = append(BgpFsActivLib[:index], BgpFsActivLib[index+1:]...)
}

func editGlobButtonResetFunc() {
    editRuleSrcPrefixLineEdit.SetText("")
    editRuleDstPrefixLineEdit.SetText("")
    editRuleIcmpTypeLineEdit.SetText("")
    editRuleIcmpCodeLineEdit.SetText("")
    editRuleIpProtoLineEdit.SetText("")
    editRulePortLineEdit.SetText("")
    editRuleSrcPortLineEdit.SetText("")
    editRuleDstPortLineEdit.SetText("")
    editRuleTcpFlagFilterLine.SetText("")
    editRuleLenLineEdit.SetText("")
    editRuleDscpLineEdit.SetText("")
    editRuleFragFilterLine.SetText("")
    editRuleActSisterValueLine.SetText("")
    editRuleActionCombo.SetCurrentIndex(0)
}

func editGlobButtonDuplicateFunc() {
    var myFsRule BgpFsRule
    var myItem *widgets.QTreeWidgetItem

    myItem = editRuleTree.CurrentItem()
    index := editRuleTree.IndexOfTopLevelItem(myItem)
    myFsRule.AddrFam = BgpFsActivLib[index].AddrFam
    myFsRule.SrcPrefix = BgpFsActivLib[index].SrcPrefix
    myFsRule.DstPrefix = BgpFsActivLib[index].DstPrefix
    myFsRule.IcmpType = BgpFsActivLib[index].IcmpType
    myFsRule.IcmpCode = BgpFsActivLib[index].IcmpCode
    myFsRule.ProtoNumber = BgpFsActivLib[index].ProtoNumber
    myFsRule.Port = BgpFsActivLib[index].Port
    myFsRule.SrcPort = BgpFsActivLib[index].SrcPort
    myFsRule.DstPort = BgpFsActivLib[index].DstPort
    myFsRule.TcpFlags = BgpFsActivLib[index].TcpFlags
    myFsRule.PacketLen = BgpFsActivLib[index].PacketLen
    myFsRule.Dscp = BgpFsActivLib[index].Dscp
    myFsRule.IpFrag = BgpFsActivLib[index].IpFrag
    myFsRule.Action = BgpFsActivLib[index].Action
    myFsRule.ActSisterValue = BgpFsActivLib[index].ActSisterValue
    BgpFsActivLib = append(BgpFsActivLib, myFsRule)
    createFullfilItemWithRule(len(BgpFsActivLib)-1, editRuleTree, BgpFsActivLib[len(BgpFsActivLib)-1])
}

func editRuleActionComboFunc(myLine *widgets.QLineEdit, myIndex int) {
    switch flowSpecActStrings[myIndex] {
        case "Drop": {
            myLine.SetReadOnly(true)
            myLine.SetPlaceholderText("No data required")
        }
        case "Shape": {
            myLine.SetReadOnly(false)
            myLine.SetPlaceholderText("3000 (expressed in KB)")
        }
        case "Redirect": {
            myLine.SetReadOnly(false)
            myLine.SetPlaceholderText("10:10")
        }
        case "Marking": {
            myLine.SetReadOnly(false)
            myLine.SetPlaceholderText("22")
        }
    }
}

func editRuleLibPushRibButtonFunc() {
    var myItem *widgets.QTreeWidgetItem
    myItem = editRuleTree.CurrentItem()
    if (*myItem != (widgets.QTreeWidgetItem{})) {
        index := editRuleTree.IndexOfTopLevelItem(myItem)
        if (sanityCheckBeforePush(BgpFsActivLib[index] ,flowspecWindow)) {
            myCommandLine := buildCommandFromFsRule(BgpFsActivLib[index])
            bgpcli.PushNewFlowSpecPath(client, myCommandLine, BgpFsActivLib[index].AddrFam)
        } else {
            return
        }

    } else {
        warningMessage := widgets.NewQMessageBox2(widgets.QMessageBox__Warning, "Rule library issue", "Please select a rule to push to the rib", widgets.QMessageBox__Ok, flowspecWindow, core.Qt__Window)
        warningMessage.Exec()
    }
}


func sanityCheckBeforePush(myRule BgpFsRule, parentWidget widgets.QWidget_ITF) bool {
    var nlriOk bool = false
    var extComOk bool = false
    var errorQMessageBox string = ""
    if ((myRule.AddrFam != "IPv4") && (myRule.AddrFam != "IPv6")) {
        errorQMessageBox = fmt.Sprintf("%s%s", errorQMessageBox, "Wrong address Family ")
    } else if (myRule.DstPrefix =="New") {
        errorQMessageBox = fmt.Sprintf("%s%s", errorQMessageBox, "This rule needs to be crafted ")
    } else if ((myRule.DstPrefix == "") && (myRule.SrcPrefix == "") && (myRule.Port == "") && (myRule.SrcPort == "") && (myRule.DstPort == "") && (myRule.TcpFlags == "") && (myRule.IcmpType == "") && (myRule.IcmpCode == "") && (myRule.ProtoNumber == "") && (myRule.PacketLen == "") && (myRule.Dscp == "") && (myRule.IpFrag == "")) {
        errorQMessageBox = fmt.Sprintf("%s%s", errorQMessageBox, "Your rule is empty ")
    } else {
        nlriOk = true
    }
    if (myRule.Action == "") {
         errorQMessageBox = fmt.Sprintf("%s%s", errorQMessageBox, "Your rule needs an Action")
    } else if (myRule.ActSisterValue == "" && myRule.Action != "Drop") {
        errorQMessageBox = fmt.Sprintf("%s%s", errorQMessageBox, "Your action needs an associated value")
    } else {
        extComOk = true
    }
    if (extComOk == false || nlriOk == false) {
        warningMessage := widgets.NewQMessageBox2(widgets.QMessageBox__Warning, "Rule badly formatted", errorQMessageBox, widgets.QMessageBox__Ok, parentWidget, core.Qt__Window)
        warningMessage.Exec()
        return false
    } else {
        return true
    }
}

func buildCommandFromFsRule(myRule BgpFsRule) string {
    var cmdLine string = "match "
    // nlri
    if (myRule.DstPrefix != "") {
        cmdLine = fmt.Sprintf("%s%s%s%s", cmdLine, "destination ", myRule.DstPrefix, " ")
    }
    if (myRule.SrcPrefix != "") {
        cmdLine = fmt.Sprintf("%s%s%s%s", cmdLine, "source ", myRule.SrcPrefix, " ")
    }
    if (myRule.Port != "") {
        cmdLine = fmt.Sprintf("%s%s%s%s", cmdLine, "port ", myRule.Port, " ")
    }
    if (myRule.SrcPort != "") {
        cmdLine = fmt.Sprintf("%s%s%s%s", cmdLine, "source-port ", myRule.SrcPort, " ")
    }
    if (myRule.DstPort != "") {
        cmdLine = fmt.Sprintf("%s%s%s%s", cmdLine, "destination-port ", myRule.DstPort, " ")
    }
    if (myRule.TcpFlags != "") {
        cmdLine = fmt.Sprintf("%s%s%s%s", cmdLine, "tcp-flags ", myRule.TcpFlags, " ")
    }
    if (myRule.IcmpType != "") {
        cmdLine = fmt.Sprintf("%s%s%s%s", cmdLine, "icmp-type ", myRule.IcmpType, " ")
    }
    if (myRule.IcmpCode != "") {
        cmdLine = fmt.Sprintf("%s%s%s%s", cmdLine, "icmp-code ", myRule.IcmpCode, " ")
    }
    if (myRule.ProtoNumber != "") {
        cmdLine = fmt.Sprintf("%s%s%s%s", cmdLine, "protocol ", myRule.ProtoNumber, " ")
    }
    if (myRule.PacketLen != "") {
        cmdLine = fmt.Sprintf("%s%s%s%s", cmdLine, "packet-length ", myRule.PacketLen, " ")
    }
    if (myRule.Dscp != "") {
        cmdLine = fmt.Sprintf("%s%s%s%s", cmdLine, "dscp ", myRule.Dscp, " ")
    }
    if (myRule.IpFrag != "") {
        cmdLine = fmt.Sprintf("%s%s%s%s", cmdLine, "fragment ", myRule.IpFrag, " ")
    }
    // extended com
    cmdLine = fmt.Sprintf("%s%s", cmdLine, "then ")
    if (myRule.Action == "Drop"){
        cmdLine = fmt.Sprintf("%s%s", cmdLine, "discard")
    } else if (myRule.Action == "Shape") {
        cmdLine = fmt.Sprintf("%s%s%s", cmdLine, "rate-limit ", myRule.ActSisterValue)
    } else if (myRule.Action == "Marking") {
        cmdLine = fmt.Sprintf("%s%s%s", cmdLine, "mark ", myRule.ActSisterValue)
    } else if (myRule.Action == "Redirect") {
        cmdLine = fmt.Sprintf("%s%s%s", cmdLine, "redirect ", myRule.ActSisterValue)
    }
    return cmdLine
}

func editRuleLibSaveButtonFunc() {
    var libFileName string
    usr, err := user.Current()
    if err != nil {
        log.Fatal( err )
    }
    myFileBox := widgets.NewQFileDialog2(flowspecWindow, "Save Flowspec Lib", usr.HomeDir, "*.fslib")
    myFileBox.SetDefaultSuffix("*.fslib")
    libFileName = myFileBox.GetSaveFileName(flowspecWindow, "Save Flowspec Lib", usr.HomeDir, "*.fslib", "*.fslib", 0)
    // libFileName = widgets.QFileDialog_GetSaveFileName(flowspecWindow, "Save Flowspec Lib", usr.HomeDir, "*.fslib", "*.fslib", 0)
    saveFsLibJsonFile(libFileName)
}

func saveFsLibJsonFile(myFile string) error {
    if (myFile != "") {
    byteBuffer, _ := json.Marshal(BgpFsActivLib)
    return(ioutil.WriteFile(myFile, byteBuffer, 0644))
    }
    return nil
}

func editRuleLibLoadButtonFunc(myTree *widgets.QTreeWidget) {
    var libFileName string
    usr, err := user.Current()
    if err != nil {
        log.Fatal( err )
    }
    myFileBox := widgets.NewQFileDialog2(flowspecWindow, "Load Flowspec Lib", usr.HomeDir, "*.fslib")
    myFileBox.SetDefaultSuffix("*.fslib")
    libFileName = myFileBox.GetOpenFileName(flowspecWindow, "Load Flowspec Lib", usr.HomeDir, "*.fslib", "*.fslib", 0)
    BgpFsActivLib = nil // makes BgpFsActiveLib eligeable for garbage collection
    cleanupTree(myTree)
    openFsLibJsonFile(libFileName, &BgpFsActivLib, myTree )
    editGlobButtonResetFunc()
}

func openFsLibJsonFile(myFile string, myRules *[]BgpFsRule, myTree *widgets.QTreeWidget) error {
    if (myFile != "") {
        data, err := ioutil.ReadFile(myFile)
        if (err == nil) {
            json.Unmarshal(data, myRules)
            fullfilTreeWithRuleLib(myTree, *myRules)
            return nil
        } else {
            return err
        }
    }
    return nil
}


func formatNlriOutputToDeleteCmdNlri(myNlri string) string {
    var nlriCmdDelete string
    nlriCmdDelete = myNlri
    nlriCmdDelete = strings.Replace(nlriCmdDelete, "[", "", -1)
    nlriCmdDelete = strings.Replace(nlriCmdDelete, "]", " ", -1)
    nlriCmdDelete = strings.Replace(nlriCmdDelete, "destination:", "destination ", -1)
    nlriCmdDelete = strings.Replace(nlriCmdDelete, "source:", "source ", -1)
    nlriCmdDelete = strings.Replace(nlriCmdDelete, "port: ", "port ", -1)
    nlriCmdDelete = strings.Replace(nlriCmdDelete, "source-port: ", "source-port ", -1)
    nlriCmdDelete = strings.Replace(nlriCmdDelete, "destination-port: ", "destination-port ", -1)
    nlriCmdDelete = strings.Replace(nlriCmdDelete, "tcp-flags: ", "tcp-flags ", -1)
    nlriCmdDelete = strings.Replace(nlriCmdDelete, "icmp-type: ", "icmp-type ", -1)
    nlriCmdDelete = strings.Replace(nlriCmdDelete, "icmp-code: ", "icmp-code ", -1)
    nlriCmdDelete = strings.Replace(nlriCmdDelete, "protocol: ", "protocol ", -1)
    nlriCmdDelete = strings.Replace(nlriCmdDelete, "packet-length: ", "packet-length ", -1)
    nlriCmdDelete = strings.Replace(nlriCmdDelete, "dscp: ", "dscp ", -1)
    nlriCmdDelete = strings.Replace(nlriCmdDelete, "fragment: ", "fragment ", -1)
    return nlriCmdDelete
}

func formatExtComOutputToDeleteCmdExtCom(myExtCom string) string {
    var extComCmdDelete string
    extComCmdDelete = myExtCom
    extComCmdDelete = strings.Replace(extComCmdDelete, "[{Origin: ?} {", "", -1)
    extComCmdDelete = strings.Replace(extComCmdDelete, "}]", "", -1)
    extComCmdDelete = strings.Replace(extComCmdDelete, "Extcomms: [", "", -1)
    extComCmdDelete = strings.Replace(extComCmdDelete, "]", "", -1)
    extComCmdDelete = strings.Replace(extComCmdDelete, "discard", "discard", -1)
    extComCmdDelete = strings.Replace(extComCmdDelete, "rate-limit: ", "rate-limit ", -1)
    extComCmdDelete = strings.Replace(extComCmdDelete, "mark: ", "mark ", -1)
    extComCmdDelete = strings.Replace(extComCmdDelete, "redirect: ", "redirect ", -1)
    return extComCmdDelete
}
// Provide UI for the whole tool
// There is a main window looking a bit like a dock
// with push buttons opening each window managing each
// function of the software

package main

import (
    "os"
    "time"
    "fmt"

   "github.com/therecipe/qt/core"
   "github.com/therecipe/qt/gui"
   "github.com/therecipe/qt/widgets"
    api "github.com/osrg/gobgp/api"
    "github.com/Matt-Texier/local-mitigation-agent/gobgpclient"
    "google.golang.org/grpc"
)

var client api.GobgpApiClient

func main() {
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
    var dockWindow = widgets.NewQMainWindow(nil, 0)
    dockWindow.Layout().DestroyQObject()
    dockWindow.SetGeometry(core.NewQRect4(100, 100, 400, 50))
    dockWindow.SetWindowTitle("Gabu")
    var dockMainLayout = widgets.NewQHBoxLayout()
    dockMainLayout.SetSpacing(6)
    dockMainLayout.SetContentsMargins(11, 11, 11, 11)
    dockWindow.SetLayout(dockMainLayout)
    // main window "dock" push button
    var dockConsolePush = widgets.NewQPushButton2("GoBgp Console", dockWindow)
    var dockFlowSpecPush = widgets.NewQPushButton2("FlowSpec RIB", dockWindow)

    var dockButtonSizePolicy = widgets.NewQSizePolicy()
    dockButtonSizePolicy.SetHorizontalPolicy(widgets.QSizePolicy__Expanding)
    dockButtonSizePolicy.SetVerticalPolicy(widgets.QSizePolicy__Expanding)
    dockButtonSizePolicy.SetHorizontalStretch(0)
    dockButtonSizePolicy.SetVerticalStretch(0)
    dockConsolePush.SetSizePolicy(dockButtonSizePolicy)
    dockFlowSpecPush.SetSizePolicy(dockButtonSizePolicy)

    // Connect buttons to functions
    dockConsolePush.ConnectClicked(func(_ bool) { dockConsolButtonClicked() })
    dockFlowSpecPush.ConnectClicked(func(_ bool) { dockFspecButtonPushed() })
    // add button to main layout
    dockMainLayout.AddWidget(dockConsolePush, 0, 0)
    dockMainLayout.AddWidget(dockFlowSpecPush, 0, 0)
    dockWindow.Show()
    widgets.QApplication_Exec()

}

func dockConsolButtonClicked() {
    consoleWin()
}

func dockFspecButtonPushed() {
    flowspecWin()
}





func consoleWin() {

    var consoleWindow = widgets.NewQMainWindow(nil, 0)
    consoleWindow.Layout().DestroyQObject()
    consoleWindow.SetGeometry(core.NewQRect4(100, 100, 1000, 600))
    consoleWindow.SetWindowTitle("GoBGP Console")
    var mainLayout = widgets.NewQHBoxLayout()
    mainLayout.SetSpacing(6)
    mainLayout.SetContentsMargins(11, 11, 11, 11)
    consoleWindow.SetLayout(mainLayout)

    // console window widgets
    // log Frame
    var logFrame = widgets.NewQFrame(consoleWindow, 0)
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
    mainLayout.AddWidget(logFrame, 0, 0)

    // command Frame
    var cmdFrame = widgets.NewQFrame(consoleWindow, 0)
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
    mainLayout.AddWidget(cmdFrame, 0, 0)

    // Connect push buttons
    cmdNeighButton.ConnectClicked(func(_ bool) { cmdNeighButtonClicked(logText) })
    cmdFsrib4Button.ConnectClicked(func(_ bool) { cmdFsrib4ButtonClicked(logText) })
    cmdFsrib6Button.ConnectClicked(func(_ bool) { cmdFsrib6ButtonClicked(logText) })

    consoleWindow.Show()
}


func cmdNeighButtonClicked(logTextWidget *widgets.QTextEdit) {
    dumpNeigh := gobgpclient.TxtdumpGetNeighbor(client)

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
    var flowspecWindow = widgets.NewQMainWindow(nil, 0)
    flowspecWindow.Layout().DestroyQObject()
    flowspecWindow.SetGeometry(core.NewQRect4(100, 100, 1000, 600))
    flowspecWindow.SetWindowTitle("Flowspec Configuration")
    var flowspecWindowLayout = widgets.NewQVBoxLayout()
    flowspecWindowLayout.SetSpacing(6)
    flowspecWindowLayout.SetContentsMargins(11, 11, 11, 11)
    flowspecWindow.SetLayout(flowspecWindowLayout)

    // create two frames, one to host flwospec rule config
    // and another one to manage flowspec FIB towards GoBGP
    var editRuleFrame = widgets.NewQFrame(flowspecWindow, 0)
    var ribManipFrame = widgets.NewQFrame(flowspecWindow, 0)
    editRuleFrame.SetSizePolicy(preferredSizePolicy)
    ribManipFrame.SetSizePolicy(preferredSizePolicy)
    editRuleFrame.SetFrameShape(widgets.QFrame__Panel)
    editRuleFrame.SetFrameShadow(widgets.QFrame__Raised)
    ribManipFrame.SetFrameShape(widgets.QFrame__Panel)
    ribManipFrame.SetFrameShadow(widgets.QFrame__Raised)
    flowspecWindowLayout.AddWidget(editRuleFrame, 0, 0)
    flowspecWindowLayout.AddWidget(ribManipFrame, 0, 0)
    var editRuleFrameLayout = widgets.NewQHBoxLayout()
    var ribManipFrameLayout = widgets.NewQHBoxLayout()
    editRuleFrame.SetLayout(editRuleFrameLayout)
    ribManipFrame.SetLayout(ribManipFrameLayout)

    // Create content of editRuleFrame
    // Widget for table that display library
    var editRuleLibWid = widgets.NewQWidget(editRuleFrame, 0)
    editRuleLibWid.SetSizePolicy(preferredSizePolicy)
    editRuleFrameLayout.AddWidget(editRuleLibWid, 0, 0)
    var editRuleLibWidLayout = widgets.NewQVBoxLayout()
    editRuleLibWid.SetLayout(editRuleLibWidLayout)
    var editRuleLabel = widgets.NewQLabel2("Rules Library", editRuleLibWid, 0)
    var editRuleTable = widgets.NewQTableWidget(editRuleLibWid)
    editRuleTable.SetSizePolicy(expandingSizePolicy)
    editRuleLibWidLayout.AddWidget(editRuleLabel, 0, 0)
    editRuleLibWidLayout.AddWidget(editRuleTable, 0, 0)

    // Widget for push button to move item from table library
    // to the edit rule widget
    var editRulePushWid = widgets.NewQWidget(editRuleFrame, 0)
    editRulePushWid.SetSizePolicy(preferredSizePolicy)
    editRuleFrameLayout.AddWidget(editRulePushWid, 0, 0)
    var editRulePushLayout = widgets.NewQVBoxLayout()
    editRulePushWid.SetLayout(editRulePushLayout)
    var editRulePushTopSpacer = widgets.NewQSpacerItem(20, 40, widgets.QSizePolicy__Fixed, widgets.QSizePolicy__Expanding)
    editRulePushLayout.AddItem(editRulePushTopSpacer)
    var editRulePushToEdit = widgets.NewQPushButton2(">", editRulePushWid)
    editRulePushLayout.AddWidget(editRulePushToEdit, 0, 0)
    editRulePushToEdit.SetMaximumWidth(20)
    var editRulePushToLib = widgets.NewQPushButton2("<", editRulePushWid)
    editRulePushLayout.AddWidget(editRulePushToLib, 0, 0)
    editRulePushToLib.SetMaximumWidth(20)
    var editRulePushBottomSpacer = widgets.NewQSpacerItem(20, 40, widgets.QSizePolicy__Fixed, widgets.QSizePolicy__Expanding)
    editRulePushLayout.AddItem(editRulePushBottomSpacer)

    // Edit rule widget creation: it includes all required
    // UI Widget to edit a BGP flowspec rule
    var editRuleMainWid = widgets.NewQWidget(editRuleFrame, 0)
    editRuleMainWid.SetSizePolicy(preferredSizePolicy)
    editRuleFrameLayout.AddWidget(editRuleMainWid, 0, 0)
    var editRuleMainWidLayout = widgets.NewQVBoxLayout()
    editRuleMainWid.SetLayout(editRuleMainWidLayout)
    // Editing widets of Edit rule widget
    var editRuleMainWidLabel = widgets.NewQLabel2("Edit Flowspec Rule", editRuleMainWid, 0)
    editRuleMainWidLayout.AddWidget(editRuleMainWidLabel, 0, 0)
    // Radio button for address family
    var editAddrFamGroupBox = widgets.NewQGroupBox2("Address Family", editRuleMainWid)
    editRuleMainWidLayout.AddWidget(editAddrFamGroupBox, 0, 0)
    var editAddrFamLayout = widgets.NewQHBoxLayout()
    editAddrFamGroupBox.SetLayout(editAddrFamLayout)
    var editAddrFamIpv4 = widgets.NewQRadioButton2("Flowspec IPv4", editAddrFamGroupBox)
    editAddrFamLayout.AddWidget(editAddrFamIpv4, 0, 0)
    editAddrFamIpv4.SetChecked(true)
    var editAddrFamIpv6 = widgets.NewQRadioButton2("Flowspec IPv6", editAddrFamGroupBox)
    editAddrFamLayout.AddWidget(editAddrFamIpv6, 0, 0)
    // Line edit for source and dest prefix
    var editRulePrefixGroupBox = widgets.NewQGroupBox2("Prefix filters", editRuleMainWid)
    editRuleMainWidLayout.AddWidget(editRulePrefixGroupBox, 0, 0)
    var editRulePrefixLayout = widgets.NewQGridLayout2()
    editRulePrefixGroupBox.SetLayout(editRulePrefixLayout)
    var (
        editRuleSrcPrefixLabel = widgets.NewQLabel2("Source Prefix:", editRulePrefixGroupBox, 0)
        editRuleDstPrefixLabel = widgets.NewQLabel2("Destination Prefix:", editRulePrefixGroupBox, 0)
        editRuleSrcPrefixLineEdit = widgets.NewQLineEdit(nil)
        editRuleDstPrefixLineEdit = widgets.NewQLineEdit(nil)
    )
    editRuleSrcPrefixLineEdit.SetPlaceholderText("1.1.1.1/32")
    editRuleDstPrefixLineEdit.SetPlaceholderText("2.2.2.2/24")
    editRulePrefixLayout.AddWidget(editRuleSrcPrefixLabel, 0, 0, 0)
    editRulePrefixLayout.AddWidget(editRuleSrcPrefixLineEdit, 0, 1, 0)
    editRulePrefixLayout.AddWidget(editRuleDstPrefixLabel, 1, 0, 0)
    editRulePrefixLayout.AddWidget(editRuleDstPrefixLineEdit, 1, 1, 0)
    // horizontal widget to group ICMP and proto type
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
        editRuleIcmpTypeLineEdit = widgets.NewQLineEdit(nil)
        editRuleIcmpCodeLineEdit = widgets.NewQLineEdit(nil)
    )
    editRuleIcmpTypeLineEdit.SetPlaceholderText("'=0' '=8'")
    editRuleIcmpCodeLineEdit.SetPlaceholderText("'=0'")
    editRuleIcmpLayout.AddWidget(editRuleIcmpTypeLabel, 0, 0, 0)
    editRuleIcmpLayout.AddWidget(editRuleIcmpTypeLineEdit, 0, 1, 0)
    editRuleIcmpLayout.AddWidget(editRuleIcmpCodeLabel, 1, 0, 0)
    editRuleIcmpLayout.AddWidget(editRuleIcmpCodeLineEdit, 1, 1, 0)
    // Line edit for IP protocol (Next header)
    var editRuleIpProtoGroupBox = widgets.NewQGroupBox2("IP protocol or Next header", editRuleMainWid)
    editRuleIcmpProtoWidLayout.AddWidget(editRuleIpProtoGroupBox, 0, 0)
    var editRuleIpProtoLayout = widgets.NewQGridLayout2()
    editRuleIpProtoGroupBox.SetLayout(editRuleIpProtoLayout)
    var (
        editRuleIpProtoLabel = widgets.NewQLabel2("Protocol number:", editRuleIcmpGroupBox, 0)
        editRuleIpProtoLineEdit = widgets.NewQLineEdit(nil)
    )
    editRuleIpProtoLineEdit.SetPlaceholderText("'=6' '=17'")
    editRuleIpProtoLayout.AddWidget(editRuleIpProtoLabel, 0, 0, 0)
    editRuleIpProtoLayout.AddWidget(editRuleIpProtoLineEdit, 0, 1, 0)

    // line edit for TCP/UDP ports
    var editRulePortGroupBox = widgets.NewQGroupBox2("Port filters", editRuleMainWid)
    editRuleMainWidLayout.AddWidget(editRulePortGroupBox, 0, 0)
    var editRulePortLayout = widgets.NewQGridLayout2()
    editRulePortGroupBox.SetLayout(editRulePortLayout)
    var (
        editRulePortLabel = widgets.NewQLabel2("Port:", editRulePortGroupBox, 0)
        editRuleSrcPortLabel = widgets.NewQLabel2("Source Port:", editRulePortGroupBox, 0)
        editRuleDstPortLabel = widgets.NewQLabel2("Destination Port:", editRulePortGroupBox, 0)
        editRulePortLineEdit = widgets.NewQLineEdit(nil)
        editRuleSrcPortLineEdit = widgets.NewQLineEdit(nil)
        editRuleDstPortLineEdit = widgets.NewQLineEdit(nil)
    )
    editRulePortLineEdit.SetPlaceholderText("'=80' '>=8080&<=8888'")
    editRuleSrcPortLineEdit.SetPlaceholderText("'=443&=80'")
    editRuleDstPortLineEdit.SetPlaceholderText("'>=1024&<=49151'")
    editRulePortLayout.AddWidget(editRulePortLabel, 0, 0, 0)
    editRulePortLayout.AddWidget(editRulePortLineEdit, 0, 1, 0)
    editRulePortLayout.AddWidget(editRuleSrcPortLabel, 1, 0, 0)
    editRulePortLayout.AddWidget(editRuleSrcPortLineEdit, 1, 1, 0)
    editRulePortLayout.AddWidget(editRuleDstPortLabel, 2, 0, 0)
    editRulePortLayout.AddWidget(editRuleDstPortLineEdit, 2, 1, 0)
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
        editRuleTcpFlagFilterLine = widgets.NewQLineEdit(nil)
        editRuleTcpFlagFilterLabel = widgets.NewQLabel2("Filter:", editRuleTcpFlagGroupBox, 0)
        editRuleTcpFlagAddButton = widgets.NewQPushButton2("Add", editRuleTcpFlagGroupBox)
    )

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



    var editRuleMainWidSpacer = widgets.NewQSpacerItem(20, 40, widgets.QSizePolicy__Fixed, widgets.QSizePolicy__Expanding)
    editRuleMainWidLayout.AddItem(editRuleMainWidSpacer)

    flowspecWindow.Show()
}

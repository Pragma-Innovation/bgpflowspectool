//source: http://doc.qt.io/qt-5/qtwidgets-widgets-lineedits-example.html

package main

import (
    "os"

    "github.com/therecipe/qt/core"
//    "github.com/therecipe/qt/gui"
    "github.com/therecipe/qt/widgets"
)

func main() {
    widgets.NewQApplication(len(os.Args), os.Args)
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

    widgets.QApplication_Exec()
}


func cmdNeighButtonClicked(logTextWidget *widgets.QTextEdit) {
    logTextWidget.Append("Button Neighbors")
}

func cmdFsrib4ButtonClicked(logTextWidget *widgets.QTextEdit) {
    logTextWidget.Append("Button FlowSpec 4")
}

func cmdFsrib6ButtonClicked(logTextWidget *widgets.QTextEdit) {
    logTextWidget.Append("Button FlowSpec 6")
}

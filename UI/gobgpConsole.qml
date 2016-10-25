import QtQuick 2.3
import QtQuick.Window 2.2
import QtQuick.Controls 1.3
import QtQuick.XmlListModel 2.0
import QtQuick.Layouts 1.2


Window {
    id: winGobgpConsole
    width: 701
    height: 415
    visible: true
    minimumWidth: 500

    TextArea {
        id: textConsole
        height: 329
        anchors.top: parent.top
        anchors.topMargin: 10
        anchors.right: parent.right
        anchors.rightMargin: 10
        anchors.left: parent.left
        anchors.leftMargin: 10
    }

    RowLayout {
        id: rowBasicsButtons
        x: 3
        y: 345
        width: 338
        height: 68
        anchors.horizontalCenterOffset: 0
        anchors.horizontalCenter: textConsole.horizontalCenter
        spacing: 5

        Button {
            id: neighborsCheck
            text: qsTr("Neighbor")
        }

        Button {
            id: ribFlowSpecV4
            text: qsTr("FS IPv4 routes")
        }

        Button {
            id: ribFlowSpecV6
            text: qsTr("FS IPv6 routes")
        }
    }

}


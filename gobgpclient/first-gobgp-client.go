package gobgpclient

import (
    "fmt"
    "time"
    "golang.org/x/net/context"
    "github.com/osrg/gobgp/packet/bgp"
    api "github.com/osrg/gobgp/api"
    "github.com/osrg/gobgp/table"
    "github.com/therecipe/qt/widgets"
)

// data strcutures used by both API functions and UI
// BGP flowspec update structure as exported from UI

type BgpFsRule struct {
    DstPrefix string
    SrcPrefix string
    AddrFam string
    Port string
    SrcPort string
    DstPort string
    TcpFlags string
    IcmpType string
    IcmpCode string
    ProtoNumber string
    PacketLen string
    Dscp string
    IpFrag string
    Action string
    ActSisterValue string
}

func formatTimedelta(d int64) string {
    u := uint64(d)
    neg := d < 0
    if neg {
        u = -u
    }
    secs := u % 60
    u /= 60
    mins := u % 60
    u /= 60
    hours := u % 24
    days := u / 24

    if days == 0 {
        return fmt.Sprintf("%02d:%02d:%02d", hours, mins, secs)
    } else {
        return fmt.Sprintf("%dd ", days) + fmt.Sprintf("%02d:%02d:%02d", hours, mins, secs)
    }
}

func TxtdumpGetNeighbor(client api.GobgpApiClient) []string {
    dumpResult := []string{}
    var NeighReq api.GetNeighborRequest
    NeighResp, e := client.GetNeighbor(context.Background(), &NeighReq)
    if e != nil {
        fmt.Println(e)
        return dumpResult
    }
    m := NeighResp.Peers
    maxaddrlen := 0
    maxaslen := 0
    maxtimelen := len("Up/Down")
    timedelta := []string{}

    // sort.Sort(m)

    now := time.Now()
    for _, p := range m {
        if i := len(p.Conf.NeighborInterface); i > maxaddrlen {
            maxaddrlen = i
        } else if j := len(p.Conf.NeighborAddress); j > maxaddrlen {
            maxaddrlen = j
        }
        if len(fmt.Sprint(p.Conf.PeerAs)) > maxaslen {
            maxaslen = len(fmt.Sprint(p.Conf.PeerAs))
        }
        timeStr := "never"
        if p.Timers.State.Uptime != 0 {
            t := int64(p.Timers.State.Downtime)
            if p.Info.BgpState == "BGP_FSM_ESTABLISHED" {
                t = int64(p.Timers.State.Uptime)
            }
            timeStr = formatTimedelta(int64(now.Sub(time.Unix(int64(t), 0)).Seconds()))
        }
        if len(timeStr) > maxtimelen {
            maxtimelen = len(timeStr)
        }
        timedelta = append(timedelta, timeStr)
    }
    var format string
    format = "%-" + fmt.Sprint(maxaddrlen) + "s" + " %" + fmt.Sprint(maxaslen) + "s" + " %" + fmt.Sprint(maxtimelen) + "s"
    format += " %-11s |%11s %8s %8s\n"
    dumpResult = append(dumpResult, fmt.Sprintf(format, "Peer", "AS", "Up/Down", "State", "#Advertised", "Received", "Accepted"))
    format_fsm := func(admin api.PeerState_AdminState, fsm string) string {
        switch admin {
        case api.PeerState_DOWN :
            return "Idle(Admin)"
        case api.PeerState_PFX_CT :
            return "Idle(PfxCt)"
        }

        if fsm == "BGP_FSM_IDLE" {
            return "Idle"
        } else if fsm == "BGP_FSM_CONNECT" {
            return "Connect"
        } else if fsm == "BGP_FSM_ACTIVE" {
            return "Active"
        } else if fsm == "BGP_FSM_OPENSENT" {
            return "Sent"
        } else if fsm == "BGP_FSM_OPENCONFIRM" {
            return "Confirm"
        } else {
            return "Establ"
        }
    }

    for i, p := range m {
        neigh := p.Conf.NeighborAddress
        if p.Conf.NeighborInterface != "" {
            neigh = p.Conf.NeighborInterface
        }
        dumpResult = append(dumpResult, fmt.Sprintf(format, neigh, fmt.Sprint(p.Conf.PeerAs), timedelta[i], format_fsm(p.Info.AdminState, p.Info.BgpState), fmt.Sprint(p.Info.Advertised), fmt.Sprint(p.Info.Received), fmt.Sprint(p.Info.Accepted)))
    }
    return dumpResult
}

func FlowSpecRibFulfillTree (client api.GobgpApiClient, myTree *widgets.QTreeWidget, myFamily string) {
    var dsts []*api.Destination
    var myNativeTable *table.Table
    resource := api.Resource_GLOBAL
    family, _ := bgp.GetRouteFamily(myFamily)

    res, err := client.GetRib(context.Background(), &api.GetRibRequest{
        Table: &api.Table{
            Type:         resource,
            Family:       uint32(family),
            Name:         "",
            Destinations: dsts,
        },
    })
    if err != nil {
        return
    }
    myNativeTable, err = res.Table.ToNativeTable()

    for _, d := range myNativeTable.GetSortedDestinations() {
        var ps []*table.Path
        ps = d.GetAllKnownPathList()
        showRouteToItem(ps, myTree)
    }
}


func showRouteToItem(pathList []*table.Path, myTree *widgets.QTreeWidget) {

    maxPrefixLen := 50
    maxNexthopLen := 20

    now := time.Now()
    for idx, p := range pathList {
        nexthop := "fictitious"
        if n := p.GetNexthop(); n != nil {
            nexthop = p.GetNexthop().String()
        }

        s := []string{}
        for _, a := range p.GetPathAttrs() {
            switch a.GetType() {
            case bgp.BGP_ATTR_TYPE_NEXT_HOP, bgp.BGP_ATTR_TYPE_MP_REACH_NLRI, bgp.BGP_ATTR_TYPE_AS_PATH, bgp.BGP_ATTR_TYPE_AS4_PATH:
                continue
            default:
                s = append(s, a.String())
            }
        }
        pattrstr := fmt.Sprint(s)

        if maxNexthopLen < len(nexthop) {
            maxNexthopLen = len(nexthop)
        }

        nlri := p.GetNlri()

        if maxPrefixLen < len(nlri.String()) {
            maxPrefixLen = len(nlri.String())
        }

        age := formatTimedelta(int64(now.Sub(p.GetTimestamp()).Seconds()))
        // fill up the tree with items
        var myItem = widgets.NewQTreeWidgetItem3(myTree, idx)
        myItem.SetText(0, nlri.String())
        myItem.SetText(1, pattrstr)
        myItem.SetText(2, age)
        myItem.SetText(3, nexthop)
        // fmt.Printf("nlri: %s\n", nlri)
        // fmt.Printf("nexthop: %s\n", nexthop)
        // fmt.Printf("pattrstr: %s\n", pattrstr)
        // fmt.Printf("Age: %s\n", age)
    }
    for i := 0; i < 4; i++ {
        myTree.ResizeColumnToContents(i)
    }
}

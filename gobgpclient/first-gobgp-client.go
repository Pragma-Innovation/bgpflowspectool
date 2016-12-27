package gobgpclient

import (
    "fmt"
    "time"
    "golang.org/x/net/context"
    "github.com/osrg/gobgp/packet/bgp"
    api "github.com/osrg/gobgp/api"
    "github.com/osrg/gobgp/table"
    "github.com/osrg/gobgp/config"
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

func FlowSpecRibToRibRules (client api.GobgpApiClient) (*[]BgpFsRule, error) {
    var myRules []BgpFsRule
    var dsts []*api.Destination
    var myNativeTable *table.Table
    resource := api.Resource_GLOBAL
    family, _ := bgp.GetRouteFamily("ipv4-flowspec")

    res, err := client.GetRib(context.Background(), &api.GetRibRequest{
        Table: &api.Table{
            Type:         resource,
            Family:       uint32(family),
            Name:         "",
            Destinations: dsts,
        },
    })
    if err != nil {
        return nil, err
    }
    myNativeTable, err = res.Table.ToNativeTable()
    fmt.Printf("Table: %v\n", myNativeTable)

    counter := 0
    for _, d := range myNativeTable.GetSortedDestinations() {
        var ps []*table.Path
        ps = d.GetAllKnownPathList()
        if counter == 0 {
            showRouteToItem(ps, false, false, false, false, true)
        } else {
            showRouteToItem(ps, false, false, false, false, false)
        }
        counter++
    }

    return &myRules, err
}


func showRouteToItem(pathList []*table.Path, showAge, showBest, showLabel, isMonitor, printHeader bool) {

    var pathStrs [][]interface{}
    maxPrefixLen := 20
    maxNexthopLen := 20
    maxAsPathLen := 20
    maxLabelLen := 10

    now := time.Now()
    for idx, p := range pathList {
        nexthop := "fictitious"
        if n := p.GetNexthop(); n != nil {
            nexthop = p.GetNexthop().String()
        }
        aspathstr := p.GetAsString()

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

        if maxAsPathLen < len(aspathstr) {
            maxAsPathLen = len(aspathstr)
        }

        best := ""
        if p.IsStale() {
            best += "S"
        }
        switch p.Validation() {
        case config.RPKI_VALIDATION_RESULT_TYPE_NOT_FOUND:
            best += "N"
        case config.RPKI_VALIDATION_RESULT_TYPE_VALID:
            best += "V"
        case config.RPKI_VALIDATION_RESULT_TYPE_INVALID:
            best += "I"
        }
        if showBest {
            if idx == 0 {
                best += "*>"
            } else {
                best += "* "
            }
        }
        nlri := p.GetNlri()
        if maxPrefixLen < len(nlri.String()) {
            maxPrefixLen = len(nlri.String())
        }

        if isMonitor {
            title := "ROUTE"
            if p.IsWithdraw {
                title = "DELROUTE"
            }
            pathStrs = append(pathStrs, []interface{}{title, nlri, nexthop, aspathstr, pattrstr})
        } else {
            args := []interface{}{best, nlri}
            if showLabel {
                label := ""
                switch nlri.(type) {
                case *bgp.LabeledIPAddrPrefix:
                    label = nlri.(*bgp.LabeledIPAddrPrefix).Labels.String()
                case *bgp.LabeledIPv6AddrPrefix:
                    label = nlri.(*bgp.LabeledIPv6AddrPrefix).Labels.String()
                case *bgp.LabeledVPNIPAddrPrefix:
                    label = nlri.(*bgp.LabeledVPNIPAddrPrefix).Labels.String()
                case *bgp.LabeledVPNIPv6AddrPrefix:
                    label = nlri.(*bgp.LabeledVPNIPv6AddrPrefix).Labels.String()
                }
                if maxLabelLen < len(label) {
                    maxLabelLen = len(label)
                }
                args = append(args, label)
            }
            args = append(args, []interface{}{nexthop, aspathstr}...)
            if showAge {
                args = append(args, formatTimedelta(int64(now.Sub(p.GetTimestamp()).Seconds())))
            }
            args = append(args, pattrstr)
            pathStrs = append(pathStrs, args)
        }
    }

    var format string
    if isMonitor {
        format = "[%s] %s via %s aspath [%s] attrs %s\n"
    } else {
        format = fmt.Sprintf("%%-3s %%-%ds", maxPrefixLen)
        if showLabel {
            format += fmt.Sprintf("%%-%ds ", maxLabelLen)
        }
        format += fmt.Sprintf("%%-%ds %%-%ds ", maxNexthopLen, maxAsPathLen)
        if showAge {
            format += "%-10s "
        }
        format += "%-s\n"

    }

    if printHeader {
        args := []interface{}{"", "Network"}
        if showLabel {
            args = append(args, "Labels")
        }
        args = append(args, []interface{}{"Next Hop", "AS_PATH"}...)
        if showAge {
            args = append(args, "Age")
        }
        args = append(args, "Attrs")
        fmt.Printf(format, args...)
    }

    for _, pathStr := range pathStrs {
        fmt.Printf(format, pathStr...)
    }
}

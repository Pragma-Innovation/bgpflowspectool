package gobgpclient

import (
    "fmt"
    api "github.com/osrg/gobgp/api"
    "golang.org/x/net/context"
//    "google.golang.org/grpc"
    // "sort"
    "time"
)

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



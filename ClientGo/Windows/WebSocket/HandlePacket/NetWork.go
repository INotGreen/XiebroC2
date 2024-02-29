package HandlePacket

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/shirou/gopsutil/net"
	"github.com/shirou/gopsutil/process"
)

func connTypeToString(t uint32) string {
	switch t {
	case 1:
		return "TCP"
	case 2:
		return "UDP"
	default:
		return fmt.Sprintf("unknown(%d)", t)
	}
}

func Network() string {
	conns, err := net.Connections("all")
	if err != nil {
		return ""
	}

	var lines []string
	for _, conn := range conns {
		if conn.Laddr.IP == "" {
			continue
		}
		proc, err := process.NewProcess(conn.Pid)
		if err != nil {
			return ""
		}
		procName, err := proc.Name()
		if err != nil {
			return ""
		}
		line := fmt.Sprintf("%s-=>%v:%v-=>%v:%v-=>%s-=>%s-=>%s",
			connTypeToString(conn.Type),
			conn.Laddr.IP, conn.Laddr.Port,
			conn.Raddr.IP, conn.Raddr.Port,

			strconv.Itoa(int(conn.Pid)),
			conn.Status,
			procName,
		)
		lines = append(lines, line)
	}

	return strings.Join(lines, "-=>")
}

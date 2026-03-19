package singleton

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/sony/sonyflake/v2"
)

var SF *sonyflake.Sonyflake

func init() {
	var st sonyflake.Settings

	st.StartTime = time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

	// 2. Set Machine ID Resolver
	st.MachineID = getMachineID

	var err error
	SF, err = sonyflake.New(st)
	if err != nil {
		logrus.Fatalf("failed to init sonyflake : %s", err.Error())
	}
}

func getMachineID() (int, error) {
	if idStr := os.Getenv("MACHINE_ID"); idStr != "" {
		id, err := strconv.ParseUint(idStr, 10, 16)
		if err == nil {
			return int(id), nil
		}
	}

	return getLower16BitPrivateIP()
}

func getLower16BitPrivateIP() (int, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return 0, err
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ip := ipnet.IP.To4()
				return int(ip[2])<<8 + int(ip[3]), nil
			}
		}
	}
	return 0, fmt.Errorf("tidak dapat menemukan private IP address")
}

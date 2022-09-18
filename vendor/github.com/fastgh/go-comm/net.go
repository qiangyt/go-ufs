package comm

import (
	"fmt"
	"net"
	"os"

	"github.com/pkg/errors"
)

func Hostname() string {
	r, err := os.Hostname()
	if err != nil {
		panic(errors.Wrapf(err, "failed to get hostname"))
	}
	return r
}

func BroadcastInterfaces(dump bool) []net.Interface {
	netIfs, err := net.Interfaces()
	if err != nil {
		panic(errors.Wrap(err, "failed to get network interfaces"))
	}

	r := make([]net.Interface, 0, len(netIfs))
	for _, netIf := range netIfs {
		flag := netIf.Flags
		if (flag | net.FlagUp) == 0 {
			// ignore because it is down
			continue
		}
		if (flag | net.FlagBroadcast) == 0 {
			// ignore non-broadcast network interface
			continue
		}

		if dump {
			//TODO: single log
			fmt.Printf("candidate interface: %s\n", netIf.Name)
		}

		r = append(r, netIf)
	}

	return r
}

func BroadcastIpWithInterface(intf net.Interface) net.IP {
	//intf.MulticastAddrs()
	addrs, err := intf.Addrs()
	if err != nil {
		panic(errors.Wrapf(err, "failed to get addresses for interface: %s", intf.Name))
	}

	for _, addr := range addrs {
		if ipAddr, isIpAddr := addr.(*net.IPNet); isIpAddr {
			ip := ipAddr.IP
			if !ip.IsLoopback() && ip.To4() != nil {
				return ip
			}
		}
	}

	return nil
}

func ResolveBroadcastIp(interfaces []net.Interface, interfaceName string) (localIp net.IP, broadcastIp net.IP) {
	for _, intF := range interfaces {
		if intF.Name == interfaceName {
			localIp = BroadcastIpWithInterface(intF)
			if localIp == nil {
				panic(fmt.Errorf("cannot get a broadcast ip for interface %s", interfaceName))
			}

			broadcastIp = make(net.IP, len(localIp))
			copy(broadcastIp, localIp)
			broadcastIp[len(broadcastIp)-1] = 255
			return
		}
	}

	panic(fmt.Errorf("interface %s is not found, or down, or not supports broadcast", interfaceName))
}

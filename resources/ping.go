package resources

import (
	"log/slog"

	"context"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/andrewsapw/avalio/status"
	"golang.org/x/net/icmp"

	"golang.org/x/net/ipv4"
)

// PingResult represents the outcome of a ping attempt
type PingResult struct {
	Reachable bool
	RTT       time.Duration
	Error     error
}

// Ping sends an ICMP echo request and waits for a reply
func Ping(ctx context.Context, host string, timeout time.Duration) PingResult {
	// Resolve host to IP address
	ip, err := net.ResolveIPAddr("ip", host)
	if err != nil {
		return PingResult{Error: fmt.Errorf("resolve error: %w", err)}
	}

	// Create ICMP listener
	conn, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		// Fallback to IPv6 if IPv4 fails
		conn, err = icmp.ListenPacket("ip6:ipv6-icmp", "::")
		if err != nil {
			return PingResult{Error: fmt.Errorf("listen error: %w", err)}
		}
	}
	defer conn.Close()

	// Create ICMP echo request
	msg := icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Code: 0,
		Body: &icmp.Echo{
			ID:   os.Getpid() & 0xffff,
			Seq:  1,
			Data: []byte("PING_TEST"),
		},
	}

	// Marshal message
	msgBytes, err := msg.Marshal(nil)
	if err != nil {
		return PingResult{Error: fmt.Errorf("marshal error: %w", err)}
	}

	// Send request
	start := time.Now()
	if _, err := conn.WriteTo(msgBytes, ip); err != nil {
		return PingResult{Error: fmt.Errorf("write error: %w", err)}
	}

	// Set deadline for reply
	if err := conn.SetReadDeadline(time.Now().Add(timeout)); err != nil {
		return PingResult{Error: fmt.Errorf("deadline error: %w", err)}
	}

	// Read reply
	reply := make([]byte, 1500)
	n, peer, err := conn.ReadFrom(reply)
	if err != nil {
		return PingResult{Error: fmt.Errorf("read error: %w", err)}
	}

	// Parse reply
	rtt := time.Since(start)
	replyMsg, err := icmp.ParseMessage(ipv4.ICMPTypeEcho.Protocol(), reply[:n])
	if err != nil {
		return PingResult{Error: fmt.Errorf("parse error: %w", err)}
	}

	// Check if it's an echo reply
	if replyMsg.Type != ipv4.ICMPTypeEchoReply {
		return PingResult{Error: fmt.Errorf("unexpected ICMP type: %v", replyMsg.Type)}
	}

	// Verify the peer address matches
	if peer.String() != ip.String() {
		return PingResult{Error: fmt.Errorf("response from different host: %s", peer)}
	}

	return PingResult{
		Reachable: true,
		RTT:       rtt,
	}
}

// isReachable performs a simple reachability check
func isReachable(host string, timeout time.Duration) (bool, error) {
	result := Ping(context.Background(), host, timeout)
	return result.Reachable, result.Error
}

type PingResource struct {
	config PingResourceConfig
	logger *slog.Logger
}

// GetName implements Resource.
func (P PingResource) GetName() string {
	return P.config.Name
}

func (P PingResource) GetType() string {
	return "ping"
}

func (P PingResource) RunCheck() (bool, []status.CheckDetails) {
	if _, err := isReachable(P.config.Address, time.Duration(uint(time.Second)*P.config.TimeoutSeconds)); err != nil {
		var checkErrors [2]status.CheckDetails
		checkErrors[0] = status.NewCheckError("Причина", "Ресурс по адресу недоступен")
		checkErrors[1] = status.NewCheckError("Адрес", P.config.Address)
		return false, checkErrors[:]
	} else {
		return true, nil
	}

}

func NewPingResource(config PingResourceConfig, logger *slog.Logger) PingResource {
	if config.TimeoutSeconds == 0 {
		config.TimeoutSeconds = 10
	}
	return PingResource{config: config, logger: logger}
}

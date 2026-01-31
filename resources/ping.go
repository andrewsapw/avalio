package resources

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"time"

	"github.com/andrewsapw/avalio/status"
	"golang.org/x/net/icmp"

	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
)

// PingResult represents the outcome of a ping attempt
type PingResult struct {
	Reachable bool
	RTT       time.Duration
	Error     error
}

// Ping sends an ICMP echo request and waits for a reply
func Ping(ctx context.Context, host string, timeout time.Duration) PingResult {
	// Resolve host to IP address, prefer IPv4 then IPv6
	ip, err := net.ResolveIPAddr("ip4", host)
	if err != nil {
		ip, err = net.ResolveIPAddr("ip6", host)
		if err != nil {
			return PingResult{Error: fmt.Errorf("resolve error: %w", err)}
		}
	}

	// Create ICMP listener
	var network, addr string
	var echoType, echoReplyType icmp.Type
	var proto int
	if ip.IP.To4() != nil {
		network = "ip4:icmp"
		addr = "0.0.0.0"
		echoType = ipv4.ICMPTypeEcho
		echoReplyType = ipv4.ICMPTypeEchoReply
		proto = ipv4.ICMPTypeEcho.Protocol()
	} else {
		network = "ip6:ipv6-icmp"
		addr = "::"
		echoType = ipv6.ICMPTypeEchoRequest
		echoReplyType = ipv6.ICMPTypeEchoReply
		proto = ipv6.ICMPTypeEchoRequest.Protocol()
	}

	conn, err := icmp.ListenPacket(network, addr)
	if err != nil {
		return PingResult{Error: fmt.Errorf("listen error: %w", err)}
	}
	defer conn.Close()

	// Create ICMP echo request
	msg := icmp.Message{
		Type: echoType,
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
		return PingResult{Error: fmt.Errorf("ошибка сериализации сообщения: %w", err)}
	}

	// Send request
	start := time.Now()
	if _, err := conn.WriteTo(msgBytes, ip); err != nil {
		return PingResult{Error: fmt.Errorf("ошибка отправки сообщения: %w", err)}
	}

	// Set deadline for reply
	if err := conn.SetReadDeadline(time.Now().Add(timeout)); err != nil {
		return PingResult{Error: fmt.Errorf("ошибка установки времени ожидания: %w", err)}
	}

	// Read reply
	reply := make([]byte, 1500)
	for {
		n, peer, err := conn.ReadFrom(reply)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				return PingResult{Error: fmt.Errorf("таймаут ожидания ответа: %w", err)}
			}
			return PingResult{Error: fmt.Errorf("ошибка чтения ответа: %w", err)}
		}

		replyMsg, err := icmp.ParseMessage(proto, reply[:n])
		if err != nil {
			// Ignore unrelated ICMP packets until deadline
			continue
		}

		if replyMsg.Type != echoReplyType {
			continue
		}

		echo, ok := replyMsg.Body.(*icmp.Echo)
		if !ok || echo.ID != (os.Getpid()&0xffff) || echo.Seq != 1 {
			continue
		}

		if peerAddr, ok := peer.(*net.IPAddr); ok {
			if !peerAddr.IP.Equal(ip.IP) {
				continue
			}
		}

		return PingResult{
			Reachable: true,
			RTT:       time.Since(start),
		}
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
	const numAttempts = 3
	const sleepDuration = time.Second * 1
	for i := range numAttempts {
		if _, err := isReachable(P.config.Address, time.Duration(P.config.TimeoutSeconds*int(time.Second))); err != nil {
			var checkErrors [3]status.CheckDetails
			checkErrors[0] = status.NewCheckError("Причина", "Ресурс по адресу недоступен")
			checkErrors[1] = status.NewCheckError("Адрес", P.config.Address)
			checkErrors[2] = status.NewCheckError("Исходная ошибка", err.Error())

			if i == (numAttempts - 1) {
				return false, checkErrors[:]
			} else {
				time.Sleep(sleepDuration)
				continue
			}
		} else {
			return true, nil
		}
	}
	return true, nil
}

func NewPingResource(config PingResourceConfig, logger *slog.Logger) PingResource {
	if config.TimeoutSeconds == 0 {
		config.TimeoutSeconds = 10
	}
	return PingResource{config: config, logger: logger}
}

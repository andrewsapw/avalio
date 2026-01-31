package resources

import (
	"context"
	"fmt"
	"log/slog"
	"math"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/andrewsapw/avalio/status"
)

// PingResult represents the outcome of a ping attempt
type PingResult struct {
	Reachable bool
	RTT       time.Duration
	Error     error
}

// Ping sends an ICMP echo request and waits for a reply
func Ping(ctx context.Context, host string, timeout time.Duration) PingResult {
	if timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, timeout)
		defer cancel()
	}

	args := []string{"-c", "1", "-n"}
	if timeout > 0 {
		if runtime.GOOS == "darwin" {
			waitMs := int(math.Ceil(timeout.Seconds() * 1000))
			args = append(args, "-W", strconv.Itoa(waitMs))
		} else {
			waitSec := int(math.Ceil(timeout.Seconds()))
			args = append(args, "-W", strconv.Itoa(waitSec))
		}
	}
	args = append(args, host)

	slog.Debug("exec ping", "cmd", "ping", "args", args)
	cmd := exec.CommandContext(ctx, "ping", args...)
	out, err := cmd.CombinedOutput()
	output := strings.TrimSpace(string(out))
	if err != nil {
		if ctx.Err() != nil {
			return PingResult{Error: fmt.Errorf("таймаут ожидания ответа: %w", ctx.Err())}
		}
		return PingResult{Error: fmt.Errorf("ping не прошел: %w: %s", err, output)}
	}

	rtt := time.Duration(0)
	re := regexp.MustCompile(`time=([0-9.]+)\s*ms`)
	if match := re.FindStringSubmatch(output); len(match) == 2 {
		if ms, parseErr := strconv.ParseFloat(match[1], 64); parseErr == nil {
			rtt = time.Duration(ms * float64(time.Millisecond))
		}
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

func NewPingResource(config PingResourceConfig) PingResource {
	if config.TimeoutSeconds == 0 {
		config.TimeoutSeconds = 10
	}
	return PingResource{config: config}
}

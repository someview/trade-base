package util

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/beevik/ntp"
)

const getIPURL = "https://api.ipify.org?format=text"
const ntpServer = "3.amazon.pool.ntp.org"

type ValidateConfig struct {
	ValidIPs   []string
	ExpireTime time.Time
}

func Validate(config *ValidateConfig) error {
	if err := validateIP(config.ValidIPs); err != nil {
		return err
	}
	if err := validateExpireTime(config.ExpireTime); err != nil {
		return err
	}
	return nil
}

func validateIP(ips []string) error {
	if len(ips) == 0 {
		return errors.New("ValidIPs is empty")
	}

	resp, err := http.Get(getIPURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	ipBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	serverIP := strings.TrimSpace(string(ipBytes))

	isValidIP := false
	for _, validIP := range ips {
		if serverIP == validIP {
			isValidIP = true
			break
		}
	}
	if !isValidIP {
		return errors.New("The network is not up to standard!!!")
	}
	return nil
}

func validateExpireTime(expireTime time.Time) error {
	if expireTime.IsZero() {
		return errors.New("ExpireTime is zero")
	}

	today, err := ntp.Time(ntpServer) // 从 NTP 服务器获取网络时间
	if err != nil {
		return err
	}
	today = today.UTC() // 将网络时间转换为 UTC 时区
	if today.After(expireTime) || today.Equal(expireTime) {
		return errors.New("Authorization has expired")
	}
	return nil
}

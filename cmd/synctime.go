package cmd

import (
	"fmt"
	"github.com/beevik/ntp"
	"github.com/spf13/cobra"
	"time"
	"golang.org/x/sys/unix"
)

var ntpServer string

var commandSyncTime = &cobra.Command{
	Use:   "synctime",
	Short: "Đồng bộ thời gian từ máy chủ NTP",
	Args:  cobra.NoArgs,
	Run:   synctimeHandle,
}

func SetSystemTime(nowTime time.Time) error {
	timeVal := unix.NsecToTimeval(nowTime.UnixNano())
	return unix.Settimeofday(&timeVal)
}

func init() {
	commandSyncTime.Flags().StringVar(&ntpServer, "server", "time.apple.com", "máy chủ NTP")
	command.AddCommand(commandSyncTime)
}

func synctimeHandle(_ *cobra.Command, _ []string) {
	t, err := ntp.Time(ntpServer)
	if err != nil {
		fmt.Println(Err("Lỗi khi lấy thời gian từ máy chủ: ", err))
		fmt.Println(Err("Đồng bộ thời gian thất bại"))
		return
	}
	err = SetSystemTime(t)
	if err != nil {
		fmt.Println(Err("Lỗi khi đặt thời gian hệ thống: ", err))
		fmt.Println(Err("Đồng bộ thời gian thất bại"))
		return
	}
	fmt.Println(Ok("Đồng bộ thời gian thành công"))
}

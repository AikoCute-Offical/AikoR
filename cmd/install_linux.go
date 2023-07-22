package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var targetVersion string

var (
	updateCommand = cobra.Command{
		Use:   "update",
		Short: "Cập nhật phiên bản AikoR",
		Run: func(_ *cobra.Command, _ []string) {
			execCommandStd("bash",
			"<(curl -Ls https://raw.githubusercontent.com/AikoCute-Offical/AikoR-Install/en/install.sh)",
				targetVersion)
		},
		Args: cobra.NoArgs,
	}
	uninstallCommand = cobra.Command{
		Use:   "uninstall",
		Short: "Gỡ cài đặt AikoR",
		Run:   uninstallHandle,
	}
)

func init() {
	updateCommand.PersistentFlags().StringVar(&targetVersion, "version", "", "phiên bản cập nhật")
	command.AddCommand(&updateCommand)
	command.AddCommand(&uninstallCommand)
}

func uninstallHandle(_ *cobra.Command, _ []string) {
	var yes string
	fmt.Println(Warn("Bạn có chắc muốn gỡ cài đặt AikoR không?(Y/n)"))
	fmt.Scan(&yes)
	if strings.ToLower(yes) != "y" {
		fmt.Println("Đã hủy gỡ cài đặt")
	}
	_, err := execCommand("systemctl stop AikoR&&systemctl disable AikoR")
	if err != nil {
		fmt.Println(Err("Lỗi khi thực thi lệnh: ", err))
		fmt.Println(Err("Gỡ cài đặt thất bại"))
		return
	}
	_ = os.RemoveAll("/etc/systemd/system/AikoR.service")
	_ = os.RemoveAll("/etc/AikoR/")
	_ = os.RemoveAll("/usr/local/AikoR/")
	_, err = execCommand("systemctl daemon-reload&&systemctl reset-failed")
	if err != nil {
		fmt.Println(Err("Lỗi khi thực thi lệnh: ", err))
		fmt.Println(Err("Gỡ cài đặt thất bại"))
		return
	}
	fmt.Println(Ok("Gỡ cài đặt thành công"))
}

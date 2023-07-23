package cmd

import (
	"fmt"
	"github.com/AikoCute-Offical/AikoR/common/exec"
	"github.com/spf13/cobra"
	"time"
)

var (
	startCommand = cobra.Command{
		Use:   "start",
		Short: "Bắt đầu dịch vụ AikoR",
		Run:   startHandle,
	}
	stopCommand = cobra.Command{
		Use:   "stop",
		Short: "Dừng dịch vụ AikoR",
		Run:   stopHandle,
	}
	restartCommand = cobra.Command{
		Use:   "restart",
		Short: "Khởi động lại dịch vụ AikoR",
		Run:   restartHandle,
	}
	logCommand = cobra.Command{
		Use:   "log",
		Short: "Xuất log AikoR",
		Run: func(_ *cobra.Command, _ []string) {
			exec.RunCommandStd("journalctl", "-u", "AikoR.service", "-e", "--no-pager", "-f")
		},
	}
)

func init() {
	command.AddCommand(&startCommand)
	command.AddCommand(&stopCommand)
	command.AddCommand(&restartCommand)
	command.AddCommand(&logCommand)
}

func startHandle(_ *cobra.Command, _ []string) {
	r, err := checkRunning()
	if err != nil {
		fmt.Println(Err("Lỗi kiểm tra trạng thái: ", err))
		fmt.Println(Err("Không thể khởi động AikoR"))
		return
	}
	if r {
		fmt.Println(Ok("AikoR đã được khởi chạy, không cần khởi động lại. Để khởi động lại, hãy chọn restart"))
	}
	_, err = exec.RunCommandByShell("systemctl start AikoR.service")
	if err != nil {
		fmt.Println(Err("Lỗi thực thi lệnh khởi động: ", err))
		fmt.Println(Err("Không thể khởi động AikoR"))
		return
	}
	time.Sleep(time.Second * 3)
	r, err = checkRunning()
	if err != nil {
		fmt.Println(Err("Lỗi kiểm tra trạng thái: ", err))
		fmt.Println(Err("Không thể khởi động AikoR"))
	}
	if !r {
		fmt.Println(Err("Không thể khởi động AikoR, hãy sử dụng lệnh AikoR log để xem thông tin log sau"))
		return
	}
	fmt.Println(Ok("Khởi động AikoR thành công, hãy sử dụng lệnh AikoR log để xem log chạy"))
}

func stopHandle(_ *cobra.Command, _ []string) {
	_, err := exec.RunCommandByShell("systemctl stop AikoR.service")
	if err != nil {
		fmt.Println(Err("Lỗi thực thi lệnh dừng: ", err))
		fmt.Println(Err("Không thể dừng AikoR"))
		return
	}
	time.Sleep(2 * time.Second)
	r, err := checkRunning()
	if err != nil {
		fmt.Println(Err("Lỗi kiểm tra trạng thái:", err))
		fmt.Println(Err("Không thể dừng AikoR"))
		return
	}
	if r {
		fmt.Println(Err("Không thể dừng AikoR, có thể vì quá thời gian dừng quá 2 giây. Hãy kiểm tra thông tin log sau"))
		return
	}
	fmt.Println(Ok("Dừng AikoR thành công"))
}

func restartHandle(_ *cobra.Command, _ []string) {
	_, err := exec.RunCommandByShell("systemctl restart AikoR.service")
	if err != nil {
		fmt.Println(Err("Lỗi thực thi lệnh khởi động lại: ", err))
		fmt.Println(Err("Không thể khởi động lại AikoR"))
		return
	}
	r, err := checkRunning()
	if err != nil {
		fmt.Println(Err("Lỗi kiểm tra trạng thái: ", err))
		fmt.Println(Err("Không thể khởi động lại AikoR"))
		return
	}
	if !r {
		fmt.Println(Err("Không thể khởi động AikoR, hãy sử dụng lệnh AikoR log để xem thông tin log sau"))
		return
	}
	fmt.Println(Ok("Khởi động lại AikoR thành công"))
}

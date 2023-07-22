package cmd

import (
	"fmt"
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
		Short: "Hiển thị log của AikoR",
		Run: func(_ *cobra.Command, _ []string) {
			execCommandStd("journalctl", "-u", "AikoR.service", "-e", "--no-pager", "-f")
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
		fmt.Println(Err("Lỗi khi kiểm tra trạng thái: ", err))
		fmt.Println(Err("Không thể bắt đầu dịch vụ AikoR"))
		return
	}
	if r {
		fmt.Println(Ok("AikoR đã chạy, không cần bắt đầu lại, nếu bạn muốn khởi động lại, vui lòng chọn restart"))
	}
	_, err = execCommand("systemctl start AikoR.service")
	if err != nil {
		fmt.Println(Err("Lỗi khi thực thi lệnh khởi động: ", err))
		fmt.Println(Err("Không thể bắt đầu dịch vụ AikoR"))
		return
	}
	time.Sleep(time.Second * 3)
	r, err = checkRunning()
	if err != nil {
		fmt.Println(Err("Lỗi khi kiểm tra trạng thái: ", err))
		fmt.Println(Err("Không thể bắt đầu dịch vụ AikoR"))
	}
	if !r {
		fmt.Println(Err("Có thể AikoR đã không khởi động thành công, vui lòng sử dụng lệnh AikoR log để xem thông tin log"))
		return
	}
	fmt.Println(Ok("Khởi động AikoR thành công, vui lòng sử dụng lệnh AikoR log để xem log"))
}

func stopHandle(_ *cobra.Command, _ []string) {
	_, err := execCommand("systemctl stop AikoR.service")
	if err != nil {
		fmt.Println(Err("Lỗi khi thực thi lệnh dừng: ", err))
		fmt.Println(Err("Không thể dừng dịch vụ AikoR"))
		return
	}
	time.Sleep(2 * time.Second)
	r, err := checkRunning()
	if err != nil {
		fmt.Println(Err("Lỗi khi kiểm tra trạng thái: ", err))
		fmt.Println(Err("Không thể dừng dịch vụ AikoR"))
		return
	}
	if r {
		fmt.Println(Err("Dừng dịch vụ AikoR thất bại, có thể do thời gian dừng quá 2 giây, vui lòng kiểm tra lại thông tin log"))
		return
	}
	fmt.Println(Ok("Dừng dịch vụ AikoR thành công"))
}

func restartHandle(_ *cobra.Command, _ []string) {
	_, err := execCommand("systemctl restart AikoR.service")
	if err != nil {
		fmt.Println(Err("Lỗi khi thực thi lệnh khởi động lại: ", err))
		fmt.Println(Err("Không thể khởi động lại dịch vụ AikoR"))
		return
	}
	r, err := checkRunning()
	if err != nil {
		fmt.Println(Err("Lỗi khi kiểm tra trạng thái: ", err))
		fmt.Println(Err("Không thể khởi động lại dịch vụ AikoR"))
		return
	}
	if !r {
		fmt.Println(Err("Có thể AikoR đã không khởi động thành công, vui lòng sử dụng lệnh AikoR log để xem thông tin log"))
		return
	}
	fmt.Println(Ok("Khởi động lại AikoR thành công"))
}

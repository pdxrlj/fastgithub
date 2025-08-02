package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"time"
)

// 默认抓取时间
const DefaultPollInterval = 10 * time.Minute

const StartMask = "# github hosts"
const EndMask = "# 数据更新时间"

type Host struct {
	ticker   *time.Ticker
	stopChan chan bool
	ctx      context.Context
}

func NewHost() *Host {
	return &Host{}
}

func (h *Host) System() string {
	return runtime.GOOS
}

// 根据不同的系统，读取每一个系统下面的hosts文件
func (h *Host) ReadHosts() string {
	system := h.System()
	hostsPath := ""
	switch system {
	case "windows":
		hostsPath = "C:\\Windows\\System32\\drivers\\etc\\hosts"
	case "linux":
		hostsPath = "/etc/hosts"
	case "darwin":
		hostsPath = "/etc/hosts"
	default:
		return "没有找到hosts文件"
	}

	content, err := os.ReadFile(hostsPath)
	if err != nil {
		return "没有找到hosts文件"
	}

	return string(content)
}

func copyFile(dst string, src *os.File) {
	dstFile, err := os.Create(dst)
	if err != nil {
		fmt.Println("创建文件失败:", err)
		return
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, src)
	if err != nil {
		fmt.Println("复制文件失败:", err)
		return
	}
}

func (h *Host) GetHostsPath() string {
	system := h.System()
	hostsPath := ""
	switch system {
	case "windows":
		hostsPath = "C:\\Windows\\System32\\drivers\\etc\\hosts"
	case "linux":
		hostsPath = "/etc/hosts"
	case "darwin":
		hostsPath = "/etc/hosts"
	}

	return hostsPath
}

// 根据路径修改hosts文件
func (h *Host) WriteHosts(content string) error {
	system := h.System()
	hostsPath := ""
	switch system {
	case "windows":
		hostsPath = "C:\\Windows\\System32\\drivers\\etc\\hosts"
	case "linux":
		hostsPath = "/etc/hosts"
	case "darwin":
		hostsPath = "/etc/hosts"
	default:
		return fmt.Errorf("没有找到hosts文件")
	}

	// 备份hosts文件
	backupPath := hostsPath + ".bak-" + strings.ReplaceAll(time.Now().Format(time.DateTime), ":", "-")

	srcFile, err := os.OpenFile(hostsPath, os.O_RDONLY, 0644)
	if err != nil {
		return fmt.Errorf("打开hosts文件失败: %w", err)
	}
	defer srcFile.Close()

	copyFile(backupPath, srcFile)

	// 查询备份文件数量
	files, err := filepath.Glob(hostsPath + ".bak-*")
	if err != nil {
		return fmt.Errorf("查询备份文件失败: %w", err)
	}

	// 删除超过10个备份文件
	if len(files) > 5 {
		// 按照最后的文件名称的时间排序
		sort.Slice(files, func(i, j int) bool {
			f1 := filepath.Base(files[i])
			f2 := filepath.Base(files[j])

			f1 = strings.TrimPrefix(f1, "hosts.bak-")
			f2 = strings.TrimPrefix(f2, "hosts.bak-")

			t1, err := time.Parse("2006-01-02 15-04-05", f1)
			if err != nil {
				return false
			}
			t2, err := time.Parse("2006-01-02 15-04-05", f2)
			if err != nil {
				return false
			}
			return t1.Before(t2)
		})

		// 删除超过10个备份文件
		for i := 0; i < len(files)-5; i++ {
			err := os.Remove(files[i])
			if err != nil {
				fmt.Println("删除备份文件失败:", err)
			}
		}
	}

	// 读取文件内容
	hostsContent, err := os.ReadFile(hostsPath)
	if err != nil {
		return fmt.Errorf("读取hosts文件失败: %w", err)
	}

	originScanner := bufio.NewScanner(bytes.NewReader(hostsContent))

	originLine := make([]string, 0)
	for originScanner.Scan() {
		originLine = append(originLine, originScanner.Text())
	}
	if err := originScanner.Err(); err != nil {
		return fmt.Errorf("读取hosts文件失败: %w", err)
	}

	startIndex := -1
	endIndex := -1

	for i, line := range originLine {
		if strings.Contains(line, StartMask) && startIndex == -1 {
			startIndex = i
		}
		if strings.Contains(line, EndMask) && endIndex == -1 {
			endIndex = i
			break
		}
	}

	// 读取开始的content一行内容，与结束的最后一行内内容中间的进行替换
	scanner := bufio.NewScanner(strings.NewReader(content))
	newContent := make([]string, 0)

	for scanner.Scan() {
		newContent = append(newContent, scanner.Text())
	}

	result := make([]string, 0)

	if startIndex == -1 { // 没有找到，直接插入到最后
		result = append(result, originLine...)
		result = append(result, "")
		result = append(result, newContent...)
	} else { // 找到，替换中间的内容
		result = append(result, originLine[:startIndex]...)
		result = append(result, newContent...)
		result = append(result, originLine[endIndex+1:]...)
	}

	err = os.WriteFile(hostsPath, []byte(strings.Join(result, "\n")), 0644)
	if err != nil {
		return fmt.Errorf("写入hosts文件失败: %w", err)
	}

	return nil
}

// 定时读取网络的hosts文件
func (h *Host) ReadHostsFromNetwork() string {
	// https://github-hosts.tinsfox.com/hosts
	httpClient := http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := httpClient.Get("https://github-hosts.tinsfox.com/hosts")
	if err != nil {
		return "网络请求失败"
	}
	defer resp.Body.Close()

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return "读取网络hosts文件失败"
	}
	return string(content)
}

func (h *Host) SetPollInterval(interval int) int {
	if interval <= 0 {
		interval = int(DefaultPollInterval.Seconds())
	}

	h.StartPoll(context.Background())

	return interval
}

// StartPoll starts polling the remote hosts file with given interval in seconds
func (h *Host) StartPoll(ctx context.Context) {
	if h.ctx == nil {
		h.ctx = ctx
	}

	h.StopPoll()
	if DefaultPollInterval <= 0 {
		fmt.Println("间隔时间必须大于0秒")
		return
	}

	h.ticker = time.NewTicker(DefaultPollInterval)
	h.stopChan = make(chan bool)

	go func() {
		for {
			select {
			case <-h.ticker.C:
				content := h.ReadHostsFromNetwork()
				// 写入hosts文件
				err := h.WriteHosts(content)
				if err != nil {
					fmt.Println("写入hosts文件失败:", err)
				}
			case <-h.stopChan:
				return
			case <-ctx.Done():
				return
			}
		}
	}()

	fmt.Sprintf("开始轮询，间隔：%d秒", int(DefaultPollInterval.Seconds()))
	return
}

// StopPoll stops the polling
func (h *Host) StopPoll() {
	if h.ticker != nil {
		h.ticker.Stop()
		h.ticker = nil
	}
	if h.stopChan != nil {
		close(h.stopChan)
		h.stopChan = nil
	}
}

func (h *Host) OpenHostsFolder() {
	hostsPath := h.GetHostsPath()
	if hostsPath == "" {
		fmt.Println("无法打开hosts文件夹，因为找不到hosts文件路径。")
		return
	}

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("explorer", "/select,"+hostsPath)
	} else {
		cmd = exec.Command("open", "-R", hostsPath)
	}

	err := cmd.Start()
	if err != nil {
		fmt.Printf("打开hosts文件夹失败: %v\n", err)
	} else {
		fmt.Printf("已打开hosts文件夹: %s\n", hostsPath)
	}
}

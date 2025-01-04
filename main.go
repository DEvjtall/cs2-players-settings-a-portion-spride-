package main

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"log"
	"main.go/Spride"
	"os"
	"strconv"
	"sync"
	"time"
)

func main() {
	file, err := os.OpenFile("go_sprider.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Print(err)
	}
	logger := log.New(file, "", log.LstdFlags)

	var wg sync.WaitGroup
	resultChan := make(chan Spride.Player, 10000) // 根据并发数调整缓冲区大小

	for i := 1; i <= 43; i++ {
		wg.Add(1)
		Spride.Sprider(strconv.Itoa(i), &wg, resultChan)
		logger.Println("爬取第" + strconv.Itoa(i) + "页信息已经完成")
	}
	go func() {
		wg.Wait()
		close(resultChan)
	}()
	// 创建一个新的 Excel 文件
	f := excelize.NewFile()
	// 创建一个工作表
	index := f.NewSheet("Sheet1")
	// 定义表头
	headers := []string{"选手姓名", "所在团队", "头像图片地址", "鼠标信息", "DPI", "灵敏度", "刷新率", "Windows 灵敏度", "准星前置瞄准", "视频分辨率", "宽高比", "显示模式", "亮度"}
	f.SetSheetRow("Sheet1", "A1", &headers)

	row := 2
	for player := range resultChan {
		// 将选手信息添加到 Excel 中
		rowData := []interface{}{
			player.Name,
			player.Team,
			player.Img,
			player.MouseSettings.Mouse,
			player.MouseSettings.DPI,
			player.MouseSettings.Sensitivity,
			player.MouseSettings.Hz,
			player.MouseSettings.WindowsSens,
			player.MouseSettings.FrontSight,
			player.VideoSettings.Resolution,
			player.VideoSettings.AspectRatio,
			player.VideoSettings.DisplayMode,
			player.VideoSettings.Brightness,
		}
		col := 'A'
		for _, value := range rowData {
			cell := fmt.Sprintf("%c%d", col, row)
			f.SetCellValue("Sheet1", cell, value)
			col++
		}
		row++

		// 使用 fmt.Printf 输出
		fmt.Printf("选手姓名: %s\n", player.Name)
		fmt.Printf("所在团队: %s\n", player.Team)
		fmt.Printf("头像图片地址: %s\n", player.Img)
		fmt.Printf("鼠标信息: %s\n", player.MouseSettings.Mouse)
		fmt.Printf("DPI: %s\n", player.MouseSettings.DPI)
		fmt.Printf("灵敏度: %s\n", player.MouseSettings.Sensitivity)
		fmt.Printf("刷新率: %s\n", player.MouseSettings.Hz)
		fmt.Printf("Windows 灵敏度: %s\n", player.MouseSettings.WindowsSens)
		fmt.Printf("准星前置瞄准: %s\n", player.MouseSettings.FrontSight)
		fmt.Printf("视频分辨率: %s\n", player.VideoSettings.Resolution)
		fmt.Printf("宽高比: %s\n", player.VideoSettings.AspectRatio)
		fmt.Printf("显示模式: %s\n", player.VideoSettings.DisplayMode)
		fmt.Printf("亮度: %s\n", player.VideoSettings.Brightness)

		// 或者使用 fmt.Println 输出
		fmt.Println("选手姓名:", player.Name)
		fmt.Println("所在团队:", player.Team)
		fmt.Println("头像图片地址:", player.Img)
		fmt.Println("鼠标信息:", player.MouseSettings.Mouse)
		fmt.Println("DPI:", player.MouseSettings.DPI)
		fmt.Println("灵敏度:", player.MouseSettings.Sensitivity)
		fmt.Println("刷新率:", player.MouseSettings.Hz)
		fmt.Println("Windows 灵敏度:", player.MouseSettings.WindowsSens)
		fmt.Println("准星前置瞄准:", player.MouseSettings.FrontSight)
		fmt.Println("视频分辨率:", player.VideoSettings.Resolution)
		fmt.Println("宽高比:", player.VideoSettings.AspectRatio)
		fmt.Println("显示模式:", player.VideoSettings.DisplayMode)
		fmt.Println("亮度:", player.VideoSettings.Brightness)
	}

	// 设置活动工作表
	f.SetActiveSheet(index)
	currentTime := time.Now().Format("2006-01-01_15-23-21")
	fileName := "players_info_" + currentTime + ".xlsx"
	// 保存 Excel 文件
	if err := f.SaveAs(fileName); err != nil {
		log.Fatal(err)
	}
}

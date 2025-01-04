package Spride

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"sync"
)

// 鼠标设置结构体
type MouseSettings struct {
	Mouse       string
	DPI         string
	Sensitivity string
	Hz          string
	WindowsSens string
	FrontSight  string
}

// 视频设置结构体
type VideoSettings struct {
	Resolution  string
	AspectRatio string
	DisplayMode string
	Brightness  string
}

// 选手信息结构体
type Player struct {
	Name          string
	Team          string
	Img           string
	MouseSettings MouseSettings
	VideoSettings VideoSettings
}

func Sprider(page string, wg *sync.WaitGroup, resultChan chan<- Player) {
	defer wg.Done()
	docDetail, _ := sreq("https://prosettings.net/games/cs2/page" + page + "/")
	// 3. 获取节点信息
	//#players > section > div > div:nth-child(1) > div.player_heading-wrapper > h4 > a
	//#players > section > div > div:nth-child(1)
	//#players > section > div > div
	// name: #players > section > div > div:nth-child(1) > div.player_heading-wrapper > h4 > a
	// team: #players > section > div > div:nth-child(1) > div.player_team
	// img: #players > section > div > div:nth-child(1) > div.player_avatar > picture > img
	docDetail.Find("#players > section > div > div"). // 这里是找到所有的div节点，做一个列表
								Each(func(i int, s *goquery.Selection) { // 遍历列表，找到选手信息
			name := s.Find("div.player_heading-wrapper > h4 > a").Text()
			team := s.Find("div.player_team").Text()
			photo, _ := s.Find("div.player_avatar > picture > img").Attr("src")
			detailLink, exists := s.Find("div.player_heading-wrapper > h4 > a").Attr("href") // 进去选手详情页然后继续爬取选手详细信息
			if exists {
				// 创建一个通道来接受 playerDoc
				playerDocChan := make(chan *goquery.Document)
				// 添加一个通道给收集选手具体信息的子通道
				wg.Add(1)
				go func(detailLink string) {
					defer wg.Done()
					// 在这里继续发去一次http请求
					playerDoc, _ := sreq(detailLink)
					playerDocChan <- playerDoc
				}(detailLink)

				// 灵敏度阶段
				//#cs2_mouse > div > div > h4 > a
				//#cs2_mouse > table > tbody > tr.format-number.field-dpi > td
				//#cs2_mouse > table > tbody > tr.format-number.field-sensitivity > td
				//#cs2_mouse > table > tbody > tr.format-select.field-hz > td
				//#cs2_mouse > table > tbody > tr.format-select.field-windowssensitivity > td
				//#cs2_crosshair > pre

				wg.Add(1)
				// 这里等待 go 匿名函数完成并接受返回结果
				go func() {
					defer wg.Done() // 确保 go 匿名函数退出前都能够正确完成
					playerDoc := <-playerDocChan
					mouse := playerDoc.Find("#cs2_mouse > div > div > h4 > a").Text()
					dpi := playerDoc.Find("#cs2_mouse > table > tbody > tr.format-number.field-dpi > td").Text()
					sensitivity := playerDoc.Find("#cs2_mouse > table > tbody > tr.format-number.field-sensitivity > td").Text()
					hz := playerDoc.Find("#cs2_mouse > table > tbody > tr.format-select.field-hz > td").Text()
					windows_sensitiviy := playerDoc.Find("#cs2_mouse > table > tbody > tr.format-select.field-windowssensitivity > td").Text()
					front_sight := playerDoc.Find("#cs2_crosshair > pre").Text()
					// 画面设置阶段
					//#video > table > tbody > tr.format-select.field-resolution > td
					//#video > table > tbody > tr.format-select.field-aspectratio > td
					//#video > table > tbody > tr.format-select.field-displaymode > td
					//#video > table > tbody > tr.format-number.field-brightness > td
					video_resolution := playerDoc.Find("#video > table > tbody > tr.format-select.field-resolution > td").Text()
					aspect_ratio := playerDoc.Find("#video > table > tbody > tr.format-select.field-aspectratio > td").Text()
					display_mode := playerDoc.Find("#video > table > tbody > tr.format-select.field-displaymode > td").Text()
					brightness := playerDoc.Find("#video > table > tbody > tr.format-number.field-brightness > td").Text()
					// 把获取到的数据都放入结构体中
					mouseSettings := MouseSettings{
						Mouse:       mouse,
						DPI:         dpi,
						Sensitivity: sensitivity,
						Hz:          hz,
						WindowsSens: windows_sensitiviy,
						FrontSight:  front_sight,
					}
					videoSettings := VideoSettings{
						Resolution:  video_resolution,
						AspectRatio: aspect_ratio,
						DisplayMode: display_mode,
						Brightness:  brightness,
					}
					player := Player{
						Name:          name,
						Team:          team,
						Img:           photo,
						MouseSettings: mouseSettings,
						VideoSettings: videoSettings,
					}
					resultChan <- player
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
				}()

			}

		})

}

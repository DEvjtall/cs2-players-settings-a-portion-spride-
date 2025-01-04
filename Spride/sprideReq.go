package Spride

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"os"
)

func sreq(url string) (docDetail *goquery.Document, err error) {
	// 创建日志收集
	file, _ := os.OpenFile("go_sprider.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	defer file.Close()
	logger := log.New(file, "", log.LstdFlags)
	// 发送请求
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logger.Println("GET URL失败:", err)
		fmt.Println(err)
		return
	}
	// 设置请求头
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux aarch64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.188 Safari/537.36 CrKey/1.54.250320")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,ja;q=0.8")
	// 获取请求响应
	resp, err := client.Do(req)
	if err != nil {
		logger.Println("获取详细页信息失败：", err, "URL: ", url)
		return nil, err
	}
	defer resp.Body.Close()
	// 解析网页
	docDetail, err = goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		logger.Println("解析网页失败：", err)
		return nil, err
	}

	return docDetail, nil
}

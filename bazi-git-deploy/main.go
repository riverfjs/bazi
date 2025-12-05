package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	bazi "github.com/warrially/BaziGo"
)

type BaziRequest struct {
	Year   int `json:"year"`
	Month  int `json:"month"`
	Day    int `json:"day"`
	Hour   int `json:"hour"`
	Minute int `json:"minute"`
	Second int `json:"second"`
	Sex    int `json:"sex"`
}

type BaziResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func main() {
	// 优先从环境变量读取端口（用于云平台部署）
	port := os.Getenv("PORT")
	if port == "" {
		// 如果没有环境变量，尝试从命令行参数读取
		if len(os.Args) > 1 {
			port = os.Args[1]
		} else {
			// 默认端口
			port = "8080"
		}
	}
	
	// 确保端口格式正确
	if port[0] != ':' {
		port = ":" + port
	}

	// 静态文件服务
	http.Handle("/", http.FileServer(http.Dir("./public")))

	// API 端点
	http.HandleFunc("/api/bazi", handleBazi)
	http.HandleFunc("/api/bazi/html", handleBaziHTML)

	log.Printf("八字服务器启动在 http://localhost%s", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

// handleBazi 返回 JSON 格式的八字信息
func handleBazi(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(BaziResponse{
			Success: false,
			Error:   "只支持 POST 请求",
		})
		return
	}

	var req BaziRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(BaziResponse{
			Success: false,
			Error:   "无效的请求格式: " + err.Error(),
		})
		return
	}

	// 验证输入
	if req.Year < 1900 || req.Year > 2100 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(BaziResponse{
			Success: false,
			Error:   "年份应在 1900-2100 之间",
		})
		return
	}

	if req.Month < 1 || req.Month > 12 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(BaziResponse{
			Success: false,
			Error:   "月份应在 1-12 之间",
		})
		return
	}

	if req.Day < 1 || req.Day > 31 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(BaziResponse{
			Success: false,
			Error:   "日期应在 1-31 之间",
		})
		return
	}

	// 计算八字
	pBazi := bazi.GetBazi(req.Year, req.Month, req.Day, req.Hour, req.Minute, req.Second, req.Sex)
	if pBazi == nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(BaziResponse{
			Success: false,
			Error:   "八字计算失败",
		})
		return
	}

	// 构建响应数据
	data := map[string]interface{}{
		"solarDate": pBazi.Date().String(),
		"lunarDate": pBazi.LunarDate().String(),
		"siZhu":     pBazi.SiZhu().String(),
		"daYun":     pBazi.DaYun().String(),
		"qiYunDate": pBazi.QiYunDate().String(),
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(BaziResponse{
		Success: true,
		Data:    data,
	})
}

// handleBaziHTML 返回 HTML 格式的八字信息
func handleBaziHTML(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	year, _ := strconv.Atoi(r.URL.Query().Get("year"))
	month, _ := strconv.Atoi(r.URL.Query().Get("month"))
	day, _ := strconv.Atoi(r.URL.Query().Get("day"))
	hour, _ := strconv.Atoi(r.URL.Query().Get("hour"))
	minute, _ := strconv.Atoi(r.URL.Query().Get("minute"))
	second, _ := strconv.Atoi(r.URL.Query().Get("second"))
	sex, _ := strconv.Atoi(r.URL.Query().Get("sex"))

	if year == 0 {
		year = 1995
		month = 6
		day = 16
		hour = 19
		minute = 7
	}

	pBazi := bazi.GetBazi(year, month, day, hour, minute, second, sex)
	if pBazi == nil {
		fmt.Fprintf(w, "<h1>八字计算失败</h1>")
		return
	}

	fmt.Fprintf(w, pBazi.ToHTML())
}


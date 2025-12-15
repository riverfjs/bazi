package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
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

// FortuneKLineData K线数据结构
type FortuneKLineData struct {
	Year  int     `json:"year"`  // 年份
	Open  float64 `json:"open"`  // 开盘价(年初运势)
	Close float64 `json:"close"` // 收盘价(年末运势)
	High  float64 `json:"high"`  // 最高(该年最佳运势)
	Low   float64 `json:"low"`   // 最低(该年最差运势)
	Score float64 `json:"score"` // 综合评分
}

// FortuneResponse 运势响应
type FortuneResponse struct {
	Success bool               `json:"success"`
	Data    []FortuneKLineData `json:"data,omitempty"`
	Error   string             `json:"error,omitempty"`
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
	http.HandleFunc("/api/bazi/fortune", handleFortune)

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

// handleFortune 计算百年运势K线数据
func handleFortune(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(FortuneResponse{
			Success: false,
			Error:   "只支持 POST 请求",
		})
		return
	}

	var req BaziRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(FortuneResponse{
			Success: false,
			Error:   "无效的请求格式: " + err.Error(),
		})
		return
	}

	// 验证输入
	if req.Year < 1900 || req.Year > 2100 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(FortuneResponse{
			Success: false,
			Error:   "年份应在 1900-2100 之间",
		})
		return
	}

	// 计算八字
	pBazi := bazi.GetBazi(req.Year, req.Month, req.Day, req.Hour, req.Minute, req.Second, req.Sex)
	if pBazi == nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(FortuneResponse{
			Success: false,
			Error:   "八字计算失败",
		})
		return
	}

	// 计算100年运势
	fortuneData := calculateHundredYearFortune(pBazi, req.Year)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(FortuneResponse{
		Success: true,
		Data:    fortuneData,
	})
}

// calculateHundredYearFortune 计算百年运势
func calculateHundredYearFortune(pBazi *bazi.TBazi, birthYear int) []FortuneKLineData {
	var fortuneData []FortuneKLineData

	// 获取日主天干(命主五行)
	dayGan := pBazi.SiZhu().DayZhu().Gan()
	dayZhi := pBazi.SiZhu().DayZhu().Zhi()

	// 获取八字四柱的天干地支
	yearGan := pBazi.SiZhu().YearZhu().Gan()
	yearZhi := pBazi.SiZhu().YearZhu().Zhi()
	monthGan := pBazi.SiZhu().MonthZhu().Gan()
	monthZhi := pBazi.SiZhu().MonthZhu().Zhi()
	hourGan := pBazi.SiZhu().HourZhu().Gan()
	hourZhi := pBazi.SiZhu().HourZhu().Zhi()

	// 计算八字五行力量
	baziPower := calculateBaziPower(yearGan, monthGan, dayGan, hourGan, yearZhi, monthZhi, dayZhi, hourZhi)

	// 获取起运年龄（起运日期年份 - 出生年份）
	qiYunAge := pBazi.QiYunDate().Year() - birthYear
	if qiYunAge < 0 {
		qiYunAge = 0
	}

	// 获取大运信息
	daYun := pBazi.DaYun()

	// 计算命盘基础分（用于判断命格强弱）
	baseScore := calculateBaseScore(baziPower, dayGan)

	// 用于保存前一年的收盘价，使K线连续
	var prevClose float64 = baseScore

	// 计算100年运势
	for i := 0; i < 100; i++ {
		currentAge := i
		year := birthYear + i

		// 计算当前年龄对应的大运
		var dayunIndex int
		var dayunZhu *bazi.TZhu
		var inDaYun bool = false

		if currentAge < qiYunAge {
			// 未起运阶段，使用月柱作为大运
			dayunIndex = -1
			dayunZhu = pBazi.SiZhu().MonthZhu()
		} else {
			// 已起运，计算大运索引
			dayunIndex = (currentAge - qiYunAge) / 10
			if dayunIndex >= 12 {
				dayunIndex = 11 // 最多12步大运
			}
			dayunZhu = daYun.Zhu(dayunIndex)
			inDaYun = true
		}

		// 获取大运天干地支
		dayunGan := dayunZhu.Gan()
		dayunZhi := dayunZhu.Zhi()
		dayunShenSha := dayunZhu.ShenSha() // 大运神煞

		// 计算流年天干地支
		liuNianGan := getGanByYear(year)
		liuNianZhi := getZhiByYear(year)

		// 计算流年神煞
		liuNianShenSha := bazi.CalcShenSha(dayGan, dayZhi, liuNianGan, liuNianZhi, "流年")

		// 在当前大运中的年数
		var dayunYearProgress int
		if inDaYun {
			dayunYearProgress = (currentAge - qiYunAge) % 10
		} else {
			dayunYearProgress = currentAge
		}

		// 计算该年运势评分（加入神煞影响）
		yearScore := calculateYearFortuneAdvanced(
			dayGan, dayZhi,
			dayunGan, dayunZhi,
			liuNianGan, liuNianZhi,
			baziPower, baseScore,
			dayunIndex, dayunYearProgress,
			currentAge, inDaYun,
			dayunShenSha, liuNianShenSha,
		)

		// 生成K线数据 - 让K线更加平滑连续
		// 年初运势：受前一年影响，平滑过渡
		open := prevClose*0.3 + yearScore*0.7

		// 年末运势：向下一年过渡
		nextYearBase := yearScore
		if i < 99 {
			// 预估下一年趋势
			nextAge := i + 1
			var nextDayunIndex int
			if nextAge < qiYunAge {
				nextDayunIndex = -1
			} else {
				nextDayunIndex = (nextAge - qiYunAge) / 10
				if nextDayunIndex >= 12 {
					nextDayunIndex = 11
				}
			}
			// 如果即将换大运，运势波动加大
			if nextDayunIndex != dayunIndex && inDaYun {
				nextYearBase = yearScore * 0.9 // 换运期略有下降
			}
		}
		close := yearScore*0.7 + nextYearBase*0.3

		// 年内最高点：受吉神影响
		high := yearScore + calculateMonthlyHighLow(year, liuNianGan, liuNianZhi, dayGan, true)

		// 年内最低点：受凶神影响
		low := yearScore + calculateMonthlyHighLow(year, liuNianGan, liuNianZhi, dayGan, false)

		// 确保数值在合理范围内(0-100)
		open = clamp(open, 0, 100)
		close = clamp(close, 0, 100)
		high = clamp(high, 0, 100)
		low = clamp(low, 0, 100)
		yearScore = clamp(yearScore, 0, 100)

		// 确保K线数据的逻辑正确性
		// 首先确保high是最大值，low是最小值
		if high < open {
			high = open
		}
		if high < close {
			high = close
		}
		if low > open {
			low = open
		}
		if low > close {
			low = close
		}
		
		// 最关键：确保low <= high，如果反了就交换
		if low > high {
			low, high = high, low
		}

		fortuneData = append(fortuneData, FortuneKLineData{
			Year:  year,
			Open:  roundFloat(open, 2),
			Close: roundFloat(close, 2),
			High:  roundFloat(high, 2),
			Low:   roundFloat(low, 2),
			Score: roundFloat(yearScore, 2),
		})

		// 保存本年收盘价，作为下一年开盘价的参考
		prevClose = close
	}

	return fortuneData
}

// calculateBaziPower 计算八字五行力量分布
func calculateBaziPower(yearGan, monthGan, dayGan, hourGan *bazi.TGan, 
	yearZhi, monthZhi, dayZhi, hourZhi *bazi.TZhi) map[string]float64 {
	power := map[string]float64{
		"金": 0.0,
		"木": 0.0,
		"水": 0.0,
		"火": 0.0,
		"土": 0.0,
	}

	// 天干力量(每个天干10分)
	addWuXingPower(power, yearGan.ToWuXing().String(), 10.0)
	addWuXingPower(power, monthGan.ToWuXing().String(), 12.0) // 月柱权重更大
	addWuXingPower(power, dayGan.ToWuXing().String(), 15.0)   // 日主权重最大
	addWuXingPower(power, hourGan.ToWuXing().String(), 10.0)

	// 地支力量(每个地支8分)
	addWuXingPower(power, yearZhi.ToWuXing().String(), 8.0)
	addWuXingPower(power, monthZhi.ToWuXing().String(), 10.0) // 月令权重更大
	addWuXingPower(power, dayZhi.ToWuXing().String(), 8.0)
	addWuXingPower(power, hourZhi.ToWuXing().String(), 8.0)

	return power
}

// addWuXingPower 添加五行力量
func addWuXingPower(power map[string]float64, wuXing string, value float64) {
	if _, ok := power[wuXing]; ok {
		power[wuXing] += value
	}
}

// calculateBaseScore 计算命盘基础分（判断命格强弱）
func calculateBaseScore(baziPower map[string]float64, dayGan *bazi.TGan) float64 {
	// 计算五行总力量
	totalPower := 0.0
	for _, power := range baziPower {
		totalPower += power
	}

	// 日主五行力量
	dayWuXing := dayGan.ToWuXing().String()
	dayPower := baziPower[dayWuXing]

	// 计算生日主的五行力量（印星）
	var supportPower float64
	switch dayWuXing {
	case "金":
		supportPower = baziPower["土"] // 土生金
	case "木":
		supportPower = baziPower["水"] // 水生木
	case "水":
		supportPower = baziPower["金"] // 金生水
	case "火":
		supportPower = baziPower["木"] // 木生火
	case "土":
		supportPower = baziPower["火"] // 火生土
	}

	// 日主及其帮扶力量占比
	// 帮扶力量 = 日主本身 + 同类五行（比劫） + 印星
	// 这里dayPower已经包含了同类五行，所以再加一次就相当于：日主+比劫+印星
	selfRatio := (dayPower*2 + supportPower) / totalPower

	// 根据日主强弱确定基础分
	// 身旺（自身强）：基础分稍高，但需要制衡
	// 身弱（自身弱）：基础分中等，需要扶持
	baseScore := 50.0
	if selfRatio > 0.4 {
		baseScore = 55.0 // 身旺
	} else if selfRatio < 0.25 {
		baseScore = 45.0 // 身弱
	}

	return baseScore
}

// calculateYearFortuneAdvanced 高级运势计算（综合大运、流年、命盘、神煞）
func calculateYearFortuneAdvanced(
	dayGan *bazi.TGan, dayZhi *bazi.TZhi,
	dayunGan *bazi.TGan, dayunZhi *bazi.TZhi,
	liuNianGan *bazi.TGan, liuNianZhi *bazi.TZhi,
	baziPower map[string]float64, baseScore float64,
	dayunIndex int, dayunYearProgress int,
	age int, inDaYun bool,
	dayunShenSha *bazi.TShenSha, liuNianShenSha *bazi.TShenSha) float64 {

	score := baseScore

	dayWuXing := dayGan.ToWuXing().String()
	dayunGanWuXing := dayunGan.ToWuXing().String()
	dayunZhiWuXing := dayunZhi.ToWuXing().String()
	liuNianGanWuXing := liuNianGan.ToWuXing().String()
	liuNianZhiWuXing := liuNianZhi.ToWuXing().String()

	// 1. 大运对日主的影响（权重35%，为神煞留出空间）
	dayunGanScore := calculateShengKeScore(dayWuXing, dayunGanWuXing)
	dayunZhiScore := calculateShengKeScore(dayWuXing, dayunZhiWuXing)
	if inDaYun {
		score += dayunGanScore * 0.25
		score += dayunZhiScore * 0.20
	} else {
		// 未起运期，影响减半
		score += dayunGanScore * 0.12
		score += dayunZhiScore * 0.10
	}

	// 2. 流年对日主的影响（权重25%）
	liuNianGanScore := calculateShengKeScore(dayWuXing, liuNianGanWuXing)
	liuNianZhiScore := calculateShengKeScore(dayWuXing, liuNianZhiWuXing)
	score += liuNianGanScore * 0.15
	score += liuNianZhiScore * 0.10 // 改为0.10，使总权重为25%

	// 3. 大运与流年的互动关系（12%）
	// 天干合化
	if checkTianGanHe(dayunGan, liuNianGan) {
		score += 7.0 // 天干相合为吉
	}
	// 地支三合、六合
	if checkDiZhiHe(dayunZhi, liuNianZhi) {
		score += 5.0 // 地支相合为吉
	}
	// 天克地冲
	if checkTianGanChong(dayunGan, liuNianGan) {
		score -= 7.0 // 天干相冲为凶
	}
	if checkDiZhiChong(dayunZhi, liuNianZhi) {
		score -= 9.0 // 地支相冲为大凶
	}

	// 4. 流年与命盘的互动（8%）
	// 流年与日柱的关系
	if checkDiZhiChong(liuNianZhi, dayZhi) {
		score -= 6.0 // 冲日支，身体、事业不顺
	}
	if checkDiZhiHe(liuNianZhi, dayZhi) {
		score += 4.0 // 合日支，贵人相助
	}

	// 5. 大运神煞影响（8%权重）
	if dayunShenSha != nil {
		shenShaScore := calculateShenShaScore(dayunShenSha)
		score += shenShaScore * 0.6 // 大运神煞影响持续10年，系数0.6
	}

	// 6. 流年神煞影响（10%权重）
	if liuNianShenSha != nil {
		shenShaScore := calculateShenShaScore(liuNianShenSha)
		score += shenShaScore * 0.8 // 流年神煞影响当年，系数0.8
	}

	// 7. 八字五行平衡度（约3-5%影响）
	balanceScore := calculateBalanceScore(baziPower, liuNianGanWuXing)
	score += balanceScore * 0.3 // 返回-5或8或0，实际影响-1.5~2.4分

	// 8. 大运内部进程影响（约1.5%影响）
	// 大运前5年和后5年运势不同
	if inDaYun {
		if dayunYearProgress < 5 {
			// 大运前5年，天干主事
			score += 1.5
		} else {
			// 大运后5年，地支主事
			score += 0.8
		}
	}

	// 9. 年龄阶段生命周期影响（约3-6%影响）
	ageEffect := calculateAgeEffect(age) // 返回0-10
	score += ageEffect * 0.3 // 实际影响0-3分

	// 10. 换大运的交接期（特殊处理）
	if inDaYun && dayunYearProgress == 0 && dayunIndex > 0 {
		score -= 4.0 // 换运之年，运势波动大，通常不利
	}

	return clamp(score, 15, 95)
}

// calculateShenShaScore 计算神煞评分
func calculateShenShaScore(shenSha *bazi.TShenSha) float64 {
	if shenSha == nil || shenSha.Count() == 0 {
		return 0.0
	}

	score := 0.0
	shenShaList := shenSha.GetList()

	// 定义神煞权重
	shenShaWeights := map[string]float64{
		// 吉神（正分）
		"天乙贵人": 12.0, // 最重要的吉神
		"天德贵人": 10.0,
		"月德贵人": 9.0,
		"文昌贵人": 8.0,
		"禄神":    8.0,
		"将星":    6.0,
		"华盖":    4.0, // 半吉半凶，艺术天赋但也主孤独
		"驿马":    5.0, // 走动、变化

		// 凶神（负分）
		"羊刃":    -8.0,  // 性格刚烈，易有血光
		"孤辰":    -6.0,  // 孤独
		"寡宿":    -6.0,  // 孤独
		"劫煞":    -7.0,  // 破财
		"亡神":    -7.0,  // 疾病、意外
		"天罗地网": -9.0,  // 困顿、受阻
		"桃花":    3.0,   // 半吉半凶，异性缘好但也主风流
	}

	// 累加神煞分数
	for _, ss := range shenShaList {
		if weight, ok := shenShaWeights[ss]; ok {
			score += weight
		}
	}

	// 特殊组合加成
	// 如果同时有天乙贵人和天德贵人，额外加分
	hasTianYi := false
	hasTianDe := false
	for _, ss := range shenShaList {
		if ss == "天乙贵人" {
			hasTianYi = true
		}
		if ss == "天德贵人" {
			hasTianDe = true
		}
	}
	if hasTianYi && hasTianDe {
		score += 5.0 // 双贵人加成
	}

	// 如果吉神多于凶神，额外加分
	jiShenCount := shenSha.GetJiShenCount()
	xiongShenCount := shenSha.GetXiongShenCount()
	if jiShenCount > xiongShenCount && jiShenCount >= 2 {
		score += float64(jiShenCount-xiongShenCount) * 2.0
	}

	return clamp(score, -20, 30) // 神煞评分范围-20到30
}

// calculateMonthlyHighLow 计算年内高低点（模拟月份波动）
func calculateMonthlyHighLow(year int, liuNianGan *bazi.TGan, liuNianZhi *bazi.TZhi, 
	dayGan *bazi.TGan, isHigh bool) float64 {
	
	dayWuXing := dayGan.ToWuXing().String()
	liuNianWuXing := liuNianGan.ToWuXing().String()

	// 基础波动幅度
	baseWave := 8.0

	// 根据五行关系调整波动幅度
	shengKeScore := calculateShengKeScore(dayWuXing, liuNianWuXing)
	
	if isHigh {
		// 计算年内最高点（应该返回正数）
		if shengKeScore > 0 {
			return baseWave + shengKeScore*0.3
		}
		return baseWave
	} else {
		// 计算年内最低点（应该返回负数）
		if shengKeScore < 0 {
			return shengKeScore * 0.3
		}
		return -baseWave
	}
}

// checkTianGanHe 检查天干是否相合
func checkTianGanHe(gan1 *bazi.TGan, gan2 *bazi.TGan) bool {
	// 天干五合：甲己合、乙庚合、丙辛合、丁壬合、戊癸合
	// 甲=0,乙=1,丙=2,丁=3,戊=4,己=5,庚=6,辛=7,壬=8,癸=9
	v1, v2 := gan1.Value(), gan2.Value()
	sum := v1 + v2
	// 甲己合(0+5=5)、乙庚合(1+6=7)、丙辛合(2+7=9)、丁壬合(3+8=11)、戊癸合(4+9=13)
	return sum == 5 || sum == 7 || sum == 9 || sum == 11 || sum == 13
}

// checkTianGanChong 检查天干是否相冲
func checkTianGanChong(gan1 *bazi.TGan, gan2 *bazi.TGan) bool {
	// 天干相冲（相克严重）
	w1 := gan1.ToWuXing().String()
	w2 := gan2.ToWuXing().String()
	
	// 阳干冲阳干，阴干冲阴干
	if (gan1.Value()%2) == (gan2.Value()%2) {
		return isWuXingKe(w1, w2) || isWuXingKe(w2, w1)
	}
	return false
}

// checkDiZhiHe 检查地支是否相合
func checkDiZhiHe(zhi1 *bazi.TZhi, zhi2 *bazi.TZhi) bool {
	// 地支六合：子丑合、寅亥合、卯戌合、辰酉合、巳申合、午未合
	// 子=0,丑=1,寅=2,卯=3,辰=4,巳=5,午=6,未=7,申=8,酉=9,戌=10,亥=11
	v1, v2 := zhi1.Value(), zhi2.Value()
	
	// 检查所有六合配对
	heMap := map[int]int{
		0:  1,  // 子丑合
		1:  0,  // 丑子合
		2:  11, // 寅亥合
		11: 2,  // 亥寅合
		3:  10, // 卯戌合
		10: 3,  // 戌卯合
		4:  9,  // 辰酉合
		9:  4,  // 酉辰合
		5:  8,  // 巳申合
		8:  5,  // 申巳合
		6:  7,  // 午未合
		7:  6,  // 未午合
	}
	
	return heMap[v1] == v2
}

// checkDiZhiChong 检查地支是否相冲
func checkDiZhiChong(zhi1 *bazi.TZhi, zhi2 *bazi.TZhi) bool {
	// 地支六冲：子午、丑未、寅申、卯酉、辰戌、巳亥
	v1, v2 := zhi1.Value(), zhi2.Value()
	diff := v1 - v2
	if diff < 0 {
		diff = -diff
	}
	// 相冲的地支相差6
	return diff == 6
}

// isWuXingKe 判断五行是否相克
func isWuXingKe(wuXing1, wuXing2 string) bool {
	keMap := map[string]string{
		"木": "土",
		"土": "水",
		"水": "火",
		"火": "金",
		"金": "木",
	}
	return keMap[wuXing1] == wuXing2
}

// calculateShengKeScore 计算五行生克关系评分
func calculateShengKeScore(wuXing1, wuXing2 string) float64 {
	// 相生关系映射
	shengMap := map[string]string{
		"木": "火",
		"火": "土",
		"土": "金",
		"金": "水",
		"水": "木",
	}

	// 相克关系映射
	keMap := map[string]string{
		"木": "土",
		"土": "水",
		"水": "火",
		"火": "金",
		"金": "木",
	}

	// 被生（他生我）
	if shengMap[wuXing2] == wuXing1 {
		return 12.0 // 被生为大吉，有贵人相助
	}

	// 生他人（我生他）
	if shengMap[wuXing1] == wuXing2 {
		return -3.0 // 生他人为泄气，消耗自身能量，不利
	}

	// 克他人（我克他）
	if keMap[wuXing1] == wuXing2 {
		return 8.0 // 我克他人得财，有力量控制局面
	}

	// 被克（他克我）
	if keMap[wuXing2] == wuXing1 {
		return -10.0 // 被克为大凶，受压制
	}

	// 同类
	if wuXing1 == wuXing2 {
		return 6.0 // 比和,有帮助
	}

	return 0.0
}

// calculateBalanceScore 计算八字五行平衡度评分
func calculateBalanceScore(baziPower map[string]float64, liuNianWuXing string) float64 {
	// 计算总力量
	total := 0.0
	for _, power := range baziPower {
		total += power
	}

	// 计算当前五行在八字中的占比
	currentPower := baziPower[liuNianWuXing]
	ratio := currentPower / total

	// 如果流年五行在八字中较弱,则流年补足为吉
	if ratio < 0.15 {
		return 8.0 // 补弱五行为吉
	} else if ratio > 0.3 {
		return -5.0 // 强者更强,过旺为凶
	}

	return 0.0
}

// calculateAgeEffect 计算年龄阶段影响
func calculateAgeEffect(age int) float64 {
	// 不同年龄段有不同的运势基调
	switch {
	case age < 10:
		return 5.0 // 童年,纯真
	case age < 20:
		return 3.0 // 青少年,成长
	case age < 30:
		return 8.0 // 青年,上升期
	case age < 40:
		return 10.0 // 壮年,高峰期
	case age < 50:
		return 6.0 // 中年,稳定期
	case age < 60:
		return 4.0 // 中老年,下降期
	case age < 70:
		return 2.0 // 老年,平稳期
	default:
		return 0.0 // 晚年,淡泊期
	}
}

// getGanByYear 根据年份获取天干
func getGanByYear(year int) *bazi.TGan {
	// 计算天干索引(以1900年为庚年起点)
	ganIndex := (year - 1900 + 6) % 10
	return bazi.NewGan(ganIndex)
}

// getZhiByYear 根据年份获取地支
func getZhiByYear(year int) *bazi.TZhi {
	// 计算地支索引(以1900年为子年起点)
	zhiIndex := (year - 1900 + 0) % 12
	return bazi.NewZhi(zhiIndex)
}

// clamp 限制数值范围
func clamp(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// roundFloat 四舍五入到指定小数位
func roundFloat(value float64, precision int) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(value*ratio) / ratio
}

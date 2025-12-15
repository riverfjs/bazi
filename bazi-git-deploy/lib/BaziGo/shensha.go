package bazi

/*
神煞是八字命理中的重要概念，是吉凶祸福的标志。
神煞分为吉神和凶神两大类。

主要神煞：
1. 天乙贵人：最有力的吉星，遇难呈祥，逢凶化吉
2. 天德贵人：德性高尚，多得贵人相助
3. 月德贵人：品德高尚，一生安稳
4. 文昌贵人：聪明好学，利于考试升学
5. 禄神：有福禄，衣食无忧
6. 将星：有威望，适合从政或管理
7. 华盖：艺术天赋，宗教缘分
8. 桃花：异性缘佳，但也主风流
9. 孤辰寡宿：孤独，婚姻不顺
10. 劫煞：破财，意外之灾
11. 亡神：疾病，意外
12. 羊刃：性格刚烈，易有血光之灾
13. 白虎：凶险，意外伤害
14. 丧门：疾病，丧事
15. 吊客：悲伤，丧事
*/

// TShenSha 神煞
type TShenSha struct {
	shenShaList []string // 神煞列表
}

// NewShenSha 新建神煞
func NewShenSha() *TShenSha {
	return &TShenSha{
		shenShaList: make([]string, 0),
	}
}

// CalcShenSha 计算神煞
func CalcShenSha(dayGan *TGan, dayZhi *TZhi, targetGan *TGan, targetZhi *TZhi, columnType string) *TShenSha {
	ss := NewShenSha()
	
	dayGanVal := dayGan.Value()
	dayZhiVal := dayZhi.Value()
	targetGanVal := targetGan.Value()
	targetZhiVal := targetZhi.Value()
	
	// 1. 天乙贵人（以日干查地支）
	if checkTianYiGuiRen(dayGanVal, targetZhiVal) {
		ss.addShenSha("天乙贵人")
	}
	
	// 2. 天德贵人（以月支查天干）
	if checkTianDeGuiRen(dayZhiVal, targetGanVal) {
		ss.addShenSha("天德贵人")
	}
	
	// 3. 月德贵人（以月支查天干）
	if checkYueDeGuiRen(dayZhiVal, targetGanVal) {
		ss.addShenSha("月德贵人")
	}
	
	// 4. 文昌贵人（以日干查地支）
	if checkWenChangGuiRen(dayGanVal, targetZhiVal) {
		ss.addShenSha("文昌贵人")
	}
	
	// 5. 禄神（以日干查地支）
	if checkLuShen(dayGanVal, targetZhiVal) {
		ss.addShenSha("禄神")
	}
	
	// 6. 羊刃（以日干查地支）
	if checkYangRen(dayGanVal, targetZhiVal) {
		ss.addShenSha("羊刃")
	}
	
	// 7. 将星（以日支或年支查其他地支）
	if checkJiangXing(dayZhiVal, targetZhiVal) {
		ss.addShenSha("将星")
	}
	
	// 8. 华盖（以日支或年支查其他地支）
	if checkHuaGai(dayZhiVal, targetZhiVal) {
		ss.addShenSha("华盖")
	}
	
	// 9. 桃花（以日支或年支查其他地支）
	if checkTaoHua(dayZhiVal, targetZhiVal) {
		ss.addShenSha("桃花")
	}
	
	// 10. 孤辰寡宿（以年支查其他地支）
	guChen, guaGu := checkGuChenGuaSu(dayZhiVal, targetZhiVal)
	if guChen {
		ss.addShenSha("孤辰")
	}
	if guaGu {
		ss.addShenSha("寡宿")
	}
	
	// 11. 劫煞（以年支或日支查其他地支）
	if checkJieSha(dayZhiVal, targetZhiVal) {
		ss.addShenSha("劫煞")
	}
	
	// 12. 亡神（以年支或日支查其他地支）
	if checkWangShen(dayZhiVal, targetZhiVal) {
		ss.addShenSha("亡神")
	}
	
	// 13. 天罗地网（以日支查）
	if checkTianLuoDiWang(targetZhiVal) {
		ss.addShenSha("天罗地网")
	}
	
	// 14. 驿马（以日支或年支查其他地支）
	if checkYiMa(dayZhiVal, targetZhiVal) {
		ss.addShenSha("驿马")
	}
	
	return ss
}

// 天乙贵人（最重要的吉神）
// 口诀：甲戊庚牛羊，乙己鼠猴乡，丙丁猪鸡位，壬癸兔蛇藏，六辛逢马虎，此是贵人方
func checkTianYiGuiRen(dayGan int, zhi int) bool {
	guiRenMap := map[int][]int{
		0: {1, 7},    // 甲：丑(1)、未(7)
		1: {0, 8},    // 乙：子(0)、申(8)
		2: {11, 9},   // 丙：亥(11)、酉(9)
		3: {11, 9},   // 丁：亥(11)、酉(9)
		4: {1, 7},    // 戊：丑(1)、未(7)
		5: {0, 8},    // 己：子(0)、申(8)
		6: {1, 7},    // 庚：丑(1)、未(7)
		7: {6, 2},    // 辛：午(6)、寅(2)
		8: {3, 5},    // 壬：卯(3)、巳(5)
		9: {3, 5},    // 癸：卯(3)、巳(5)
	}
	
	if zhiList, ok := guiRenMap[dayGan]; ok {
		for _, z := range zhiList {
			if z == zhi {
				return true
			}
		}
	}
	return false
}

// 天德贵人（以月支查天干）
// 正月生者见丁，二月生者见申，三月生者见壬...
func checkTianDeGuiRen(monthZhi int, gan int) bool {
	tianDeMap := map[int]int{
		2:  3,  // 寅月(正月)见丁(3)
		3:  8,  // 卯月(二月)见壬(8)或申
		4:  9,  // 辰月(三月)见癸(9)
		5:  8,  // 巳月(四月)见壬(8)
		6:  2,  // 午月(五月)见丙(2)
		7:  1,  // 未月(六月)见乙(1)
		8:  0,  // 申月(七月)见甲(0)
		9:  5,  // 酉月(八月)见己(5)
		10: 8,  // 戌月(九月)见壬(8)
		11: 6,  // 亥月(十月)见庚(6)
		0:  0,  // 子月(十一月)见甲(0)
		1:  4,  // 丑月(十二月)见戊(4)
	}
	
	if tianDeGan, ok := tianDeMap[monthZhi]; ok {
		return tianDeGan == gan
	}
	return false
}

// 月德贵人（以月支查天干）
// 寅午戌月在丙，申子辰月在壬，亥卯未月在甲，巳酉丑月在庚
func checkYueDeGuiRen(monthZhi int, gan int) bool {
	// 寅午戌月(火局)在丙
	if (monthZhi == 2 || monthZhi == 6 || monthZhi == 10) && gan == 2 {
		return true
	}
	// 申子辰月(水局)在壬
	if (monthZhi == 8 || monthZhi == 0 || monthZhi == 4) && gan == 8 {
		return true
	}
	// 亥卯未月(木局)在甲
	if (monthZhi == 11 || monthZhi == 3 || monthZhi == 7) && gan == 0 {
		return true
	}
	// 巳酉丑月(金局)在庚
	if (monthZhi == 5 || monthZhi == 9 || monthZhi == 1) && gan == 6 {
		return true
	}
	return false
}

// 文昌贵人（以日干查地支）
// 口诀：甲乙巳午报君知，丙戊申宫丁己鸡，庚猪辛鼠壬逢虎，癸人见卯入云梯
func checkWenChangGuiRen(dayGan int, zhi int) bool {
	wenChangMap := map[int]int{
		0: 5,  // 甲：巳(5)
		1: 6,  // 乙：午(6)
		2: 8,  // 丙：申(8)
		3: 9,  // 丁：酉(9)
		4: 8,  // 戊：申(8)
		5: 9,  // 己：酉(9)
		6: 11, // 庚：亥(11)
		7: 0,  // 辛：子(0)
		8: 2,  // 壬：寅(2)
		9: 3,  // 癸：卯(3)
	}
	
	if wenChangZhi, ok := wenChangMap[dayGan]; ok {
		return wenChangZhi == zhi
	}
	return false
}

// 禄神（以日干查地支，临官位）
// 甲禄寅、乙禄卯、丙戊禄巳、丁己禄午、庚禄申、辛禄酉、壬禄亥、癸禄子
func checkLuShen(dayGan int, zhi int) bool {
	luShenMap := map[int]int{
		0: 2,  // 甲：寅(2)
		1: 3,  // 乙：卯(3)
		2: 5,  // 丙：巳(5)
		3: 6,  // 丁：午(6)
		4: 5,  // 戊：巳(5)
		5: 6,  // 己：午(6)
		6: 8,  // 庚：申(8)
		7: 9,  // 辛：酉(9)
		8: 11, // 壬：亥(11)
		9: 0,  // 癸：子(0)
	}
	
	if luZhi, ok := luShenMap[dayGan]; ok {
		return luZhi == zhi
	}
	return false
}

// 羊刃（以日干查地支，帝旺位）
// 甲刃卯、乙刃寅、丙戊刃午、丁己刃巳、庚刃酉、辛刃申、壬刃子、癸刃亥
func checkYangRen(dayGan int, zhi int) bool {
	yangRenMap := map[int]int{
		0: 3,  // 甲：卯(3)
		1: 2,  // 乙：寅(2)
		2: 6,  // 丙：午(6)
		3: 5,  // 丁：巳(5)
		4: 6,  // 戊：午(6)
		5: 5,  // 己：巳(5)
		6: 9,  // 庚：酉(9)
		7: 8,  // 辛：申(8)
		8: 0,  // 壬：子(0)
		9: 11, // 癸：亥(11)
	}
	
	if yangRenZhi, ok := yangRenMap[dayGan]; ok {
		return yangRenZhi == zhi
	}
	return false
}

// 将星（以日支或年支查其他地支）
// 寅午戌见午，申子辰见子，巳酉丑见酉，亥卯未见卯
func checkJiangXing(baseZhi int, targetZhi int) bool {
	// 寅午戌(火局)见午
	if (baseZhi == 2 || baseZhi == 6 || baseZhi == 10) && targetZhi == 6 {
		return true
	}
	// 申子辰(水局)见子
	if (baseZhi == 8 || baseZhi == 0 || baseZhi == 4) && targetZhi == 0 {
		return true
	}
	// 巳酉丑(金局)见酉
	if (baseZhi == 5 || baseZhi == 9 || baseZhi == 1) && targetZhi == 9 {
		return true
	}
	// 亥卯未(木局)见卯
	if (baseZhi == 11 || baseZhi == 3 || baseZhi == 7) && targetZhi == 3 {
		return true
	}
	return false
}

// 华盖（以日支或年支查其他地支）
// 寅午戌见戌，申子辰见辰，巳酉丑见丑，亥卯未见未
func checkHuaGai(baseZhi int, targetZhi int) bool {
	// 寅午戌(火局)见戌
	if (baseZhi == 2 || baseZhi == 6 || baseZhi == 10) && targetZhi == 10 {
		return true
	}
	// 申子辰(水局)见辰
	if (baseZhi == 8 || baseZhi == 0 || baseZhi == 4) && targetZhi == 4 {
		return true
	}
	// 巳酉丑(金局)见丑
	if (baseZhi == 5 || baseZhi == 9 || baseZhi == 1) && targetZhi == 1 {
		return true
	}
	// 亥卯未(木局)见未
	if (baseZhi == 11 || baseZhi == 3 || baseZhi == 7) && targetZhi == 7 {
		return true
	}
	return false
}

// 桃花（咸池，以日支或年支查其他地支）
// 寅午戌见卯，申子辰见酉，巳酉丑见午，亥卯未见子
func checkTaoHua(baseZhi int, targetZhi int) bool {
	// 寅午戌(火局)见卯
	if (baseZhi == 2 || baseZhi == 6 || baseZhi == 10) && targetZhi == 3 {
		return true
	}
	// 申子辰(水局)见酉
	if (baseZhi == 8 || baseZhi == 0 || baseZhi == 4) && targetZhi == 9 {
		return true
	}
	// 巳酉丑(金局)见午
	if (baseZhi == 5 || baseZhi == 9 || baseZhi == 1) && targetZhi == 6 {
		return true
	}
	// 亥卯未(木局)见子
	if (baseZhi == 11 || baseZhi == 3 || baseZhi == 7) && targetZhi == 0 {
		return true
	}
	return false
}

// 孤辰寡宿（以年支查其他地支）
// 亥子丑人，见寅为孤，见戌为寡
// 寅卯辰人，见巳为孤，见丑为寡
// 巳午未人，见申为孤，见辰为寡
// 申酉戌人，见亥为孤，见未为寡
func checkGuChenGuaSu(baseZhi int, targetZhi int) (bool, bool) {
	guChen := false
	guaSu := false
	
	// 亥子丑人
	if baseZhi == 11 || baseZhi == 0 || baseZhi == 1 {
		if targetZhi == 2 {
			guChen = true
		}
		if targetZhi == 10 {
			guaSu = true
		}
	}
	// 寅卯辰人
	if baseZhi == 2 || baseZhi == 3 || baseZhi == 4 {
		if targetZhi == 5 {
			guChen = true
		}
		if targetZhi == 1 {
			guaSu = true
		}
	}
	// 巳午未人
	if baseZhi == 5 || baseZhi == 6 || baseZhi == 7 {
		if targetZhi == 8 {
			guChen = true
		}
		if targetZhi == 4 {
			guaSu = true
		}
	}
	// 申酉戌人
	if baseZhi == 8 || baseZhi == 9 || baseZhi == 10 {
		if targetZhi == 11 {
			guChen = true
		}
		if targetZhi == 7 {
			guaSu = true
		}
	}
	
	return guChen, guaSu
}

// 劫煞（以年支或日支查其他地支）
// 申子辰见巳，寅午戌见亥，巳酉丑见申，亥卯未见寅
func checkJieSha(baseZhi int, targetZhi int) bool {
	// 申子辰(水局)见巳
	if (baseZhi == 8 || baseZhi == 0 || baseZhi == 4) && targetZhi == 5 {
		return true
	}
	// 寅午戌(火局)见亥
	if (baseZhi == 2 || baseZhi == 6 || baseZhi == 10) && targetZhi == 11 {
		return true
	}
	// 巳酉丑(金局)见申
	if (baseZhi == 5 || baseZhi == 9 || baseZhi == 1) && targetZhi == 8 {
		return true
	}
	// 亥卯未(木局)见寅
	if (baseZhi == 11 || baseZhi == 3 || baseZhi == 7) && targetZhi == 2 {
		return true
	}
	return false
}

// 亡神（以年支或日支查其他地支）
// 申子辰见亥，寅午戌见巳，巳酉丑见寅，亥卯未见申
func checkWangShen(baseZhi int, targetZhi int) bool {
	// 申子辰(水局)见亥
	if (baseZhi == 8 || baseZhi == 0 || baseZhi == 4) && targetZhi == 11 {
		return true
	}
	// 寅午戌(火局)见巳
	if (baseZhi == 2 || baseZhi == 6 || baseZhi == 10) && targetZhi == 5 {
		return true
	}
	// 巳酉丑(金局)见寅
	if (baseZhi == 5 || baseZhi == 9 || baseZhi == 1) && targetZhi == 2 {
		return true
	}
	// 亥卯未(木局)见申
	if (baseZhi == 11 || baseZhi == 3 || baseZhi == 7) && targetZhi == 8 {
		return true
	}
	return false
}

// 天罗地网（辰为天罗，戌为地网）
func checkTianLuoDiWang(zhi int) bool {
	return zhi == 4 || zhi == 10 // 辰(4)或戌(10)
}

// 驿马（以日支或年支查其他地支）
// 申子辰马在寅，寅午戌马在申，巳酉丑马在亥，亥卯未马在巳
func checkYiMa(baseZhi int, targetZhi int) bool {
	// 申子辰(水局)马在寅
	if (baseZhi == 8 || baseZhi == 0 || baseZhi == 4) && targetZhi == 2 {
		return true
	}
	// 寅午戌(火局)马在申
	if (baseZhi == 2 || baseZhi == 6 || baseZhi == 10) && targetZhi == 8 {
		return true
	}
	// 巳酉丑(金局)马在亥
	if (baseZhi == 5 || baseZhi == 9 || baseZhi == 1) && targetZhi == 11 {
		return true
	}
	// 亥卯未(木局)马在巳
	if (baseZhi == 11 || baseZhi == 3 || baseZhi == 7) && targetZhi == 5 {
		return true
	}
	return false
}

// addShenSha 添加神煞
func (m *TShenSha) addShenSha(name string) {
	// 避免重复
	for _, ss := range m.shenShaList {
		if ss == name {
			return
		}
	}
	m.shenShaList = append(m.shenShaList, name)
}

// GetList 获取神煞列表
func (m *TShenSha) GetList() []string {
	return m.shenShaList
}

// String 转换成可阅读的字符串
func (m *TShenSha) String() string {
	if len(m.shenShaList) == 0 {
		return ""
	}
	result := ""
	for i, ss := range m.shenShaList {
		if i > 0 {
			result += " "
		}
		result += ss
	}
	return result
}

// Count 神煞数量
func (m *TShenSha) Count() int {
	return len(m.shenShaList)
}

// HasJiShen 是否有吉神
func (m *TShenSha) HasJiShen() bool {
	jiShenList := []string{"天乙贵人", "天德贵人", "月德贵人", "文昌贵人", "禄神", "将星"}
	for _, ss := range m.shenShaList {
		for _, js := range jiShenList {
			if ss == js {
				return true
			}
		}
	}
	return false
}

// HasXiongShen 是否有凶神
func (m *TShenSha) HasXiongShen() bool {
	xiongShenList := []string{"羊刃", "孤辰", "寡宿", "劫煞", "亡神", "天罗地网"}
	for _, ss := range m.shenShaList {
		for _, xs := range xiongShenList {
			if ss == xs {
				return true
			}
		}
	}
	return false
}

// GetJiShenCount 获取吉神数量
func (m *TShenSha) GetJiShenCount() int {
	jiShenList := []string{"天乙贵人", "天德贵人", "月德贵人", "文昌贵人", "禄神", "将星", "华盖", "驿马"}
	count := 0
	for _, ss := range m.shenShaList {
		for _, js := range jiShenList {
			if ss == js {
				count++
				break
			}
		}
	}
	return count
}

// GetXiongShenCount 获取凶神数量
func (m *TShenSha) GetXiongShenCount() int {
	xiongShenList := []string{"羊刃", "孤辰", "寡宿", "劫煞", "亡神", "天罗地网", "桃花"}
	count := 0
	for _, ss := range m.shenShaList {
		for _, xs := range xiongShenList {
			if ss == xs {
				count++
				break
			}
		}
	}
	return count
}


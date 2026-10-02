/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-01-26 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-10-02 09:57:31
 * @FilePath: \go-logger\sanitize.go
 * @Description: 敏感字段脱敏（按字段名零分配判定 + 可选的值内容扫描）
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */
package logger

import "strings"

// redacted 敏感值的替换文本
const redacted = "[REDACTED]"

// SanitizeMode 脱敏强度
type SanitizeMode int32

const (
	// SanitizeOff 关闭脱敏
	SanitizeOff SanitizeMode = iota

	// SanitizeKeys 仅按字段名脱敏（默认）
	//
	// 判定全程零分配：规范 snake_case 字段名一次查表命中，非规范形态由后缀扫描
	// 边归一化边比较，两者都落空才退回栈上缓冲做一次完整归一化后重查
	// 不读取值内容，因此既不会误伤业务正文，成本也不随消息长度增加
	SanitizeKeys

	// SanitizeKeysAndValues 字段名 + 值内容扫描
	//
	// 额外在 msg 与各字段值里查找凭据、JWT、邮箱、手机号与带签名的私有链接
	// 代价是每条日志都要扫一遍文本，且业务正文可能被误伤（如长订单号被当成手机号），故不作为默认值
	SanitizeKeysAndValues
)

// redactAction 字段名命中规则后施加的动作
type redactAction uint8

const (
	keepAsIs  redactAction = iota // 不处理
	wipeValue                     // 整体替换为 [REDACTED]
	maskPhone                     // 号码掩码（保留末 4 位）
	maskEmail                     // 邮箱掩码（保留首字符与域名）
)

// exactKeyRules 字段名 → 动作 的精确规则表，键为归一化形态（去 _-. 并转小写）
//
// 合成一张表而非按动作分设多个集合：判定路径上只有一次查表，
// 而分设集合要依次查到命中为止，未命中的字段名（生产里的绝大多数）要白跑满全程
//
// 入表门槛是"该字段名在任何业务里都只可能承载敏感内容"name、user_id、order_id
// 这类在通用日志里多半是服务名、主键与排障主线索，一律不入表——宁可漏擦由调用方补 RegisterSensitiveKeys
// 不可误擦让日志失去排障价值容易被后缀规则牵连的运维字段另见 visibleKeySuffixes
var exactKeyRules = buildKeyRules(map[redactAction][]string{
	wipeValue: {
		// ---- 认证与凭据 ----
		"authorization", "proxy_authorization", "www_authenticate",
		"auth", "basic_auth", "digest_auth", "bearer",
		"cookie", "set_cookie",
		"password", "passwd", "pwd", "passcode", "pass_phrase", "passphrase",
		"old_password", "new_password", "confirm_password",
		"token", "access_token", "refresh_token", "id_token", "auth_token",
		"bearer_token", "session_token", "csrf_token", "xsrf_token", "jwt",
		"api_key", "api_secret", "app_key", "app_secret", "client_secret",
		"access_key", "access_key_secret", "secret_key", "secret", "private_key",
		"encryption_key", "signing_key", "ssh_key", "license_key",
		"credential", "credentials",
		"signature", "sign", "sign_str",
		"session", "session_id", "sessid", "jsessionid", "phpsessid",

		// ---- 验证码与兑换码：一次性凭据，泄漏即可被抢用 ----
		"otp", "totp", "mfa_code", "two_factor_code",
		"verification_code", "verify_code", "sms_code", "email_code",
		"captcha", "captcha_code", "authorization_code", "auth_code",
		"reset_token", "recovery_code", "backup_code", "activation_code",
		// 电商与游戏：兑换码、礼品码、激活码等同于资金
		"redeem_code", "gift_code", "gift_card", "gift_card_no",
		"coupon_code", "promo_code", "invite_code", "cd_key", "serial_key",

		// ---- 支付与金融（电商）----
		"card_number", "card_no", "cardno", "pan",
		"credit_card", "creditcard", "debit_card", "bank_card",
		"cvv", "cvv2", "cvc", "cvc2", "csc",
		"bank_account", "account_number", "iban", "swift_code", "routing_number",
		"pay_password", "payment_password", "withdraw_password", "trade_password",

		// ---- 身份证件 ----
		"id_card", "idcard", "id_card_no", "identity_card", "id_number",
		"national_id", "ssn", "social_security_number",
		"passport", "passport_no", "driver_license", "driving_license", "tax_id",

		// ---- 联系地址（客服与电商的收发货信息）----
		// address 不作后缀规则：ip_address / remote_address 是运维字段，必须可见
		"home_address", "detail_address", "street_address", "mailing_address",
		"shipping_address", "billing_address", "delivery_address",
		"receiver_address", "consignee_address", "contact_address",

		// ---- 自然人姓名（客服与电商工单里的真实姓名）----
		// name 不入表也不作后缀：通用日志里 service_name / table_name / file_name 远多于人名
		"real_name", "realname", "receiver_name", "consignee", "contact_name",
		"customer_name", "holder_name", "card_holder",

		// ---- 出生信息 ----
		"birthday", "birth_date", "date_of_birth", "dob",

		// ---- 设备与广告标识（游戏风控链路常带，属可追踪 PII）----
		"imei", "imsi", "idfa", "idfv", "oaid", "aaid", "android_id", "mac_address",

		// ---- 预签名链接：链接本身即凭据，泄漏等同泄漏对象存储里的文件 ----
		"signed_url", "presigned_url", "download_url", "upload_url",
	},

	// 以下两类保留尾部若干位，排障时仍可对号
	maskPhone: {
		"phone", "phone_number", "mobile", "mobile_phone", "tel", "telephone",
		"contact_phone", "receiver_phone", "customer_phone",
	},
	maskEmail: {"email", "email_address", "mail", "contact_email", "customer_email"},
})

// suffixKeyRules 字段名后缀规则（归一化形态），覆盖任意前缀的变体
//
// 只收录"作为后缀时语义不会漂移"的词被排除的典型反例：
//   - key   → cache_key / partition_key / idempotency_key / shard_key 全是运维字段
//   - code  → error_code / status_code / grpc_code 是排障主线索
//   - name  → service_name / table_name / file_name 与人名无关
//   - sign  → design / redesign 会被整词牵连
//   - addr  → ip_address / bind_address 是运维字段
var suffixKeyRules = []struct {
	suffix string
	action redactAction
}{
	{"password", wipeValue},
	{"passwd", wipeValue},
	{"passcode", wipeValue},
	{"passphrase", wipeValue},
	{"token", wipeValue},
	{"secret", wipeValue},
	{"apikey", wipeValue},
	{"appkey", wipeValue},
	{"privatekey", wipeValue},
	{"secretkey", wipeValue},
	{"accesskey", wipeValue},
	{"licensekey", wipeValue},
	{"credential", wipeValue},
	{"credentials", wipeValue},
	{"signature", wipeValue},
	{"sessionid", wipeValue},
	{"idcard", wipeValue},
	{"idcardno", wipeValue},
	{"cardno", wipeValue},
	{"cardnumber", wipeValue},
	{"signedurl", wipeValue},
	{"downloadurl", wipeValue},
	{"uploadurl", wipeValue},

	{"phone", maskPhone},
	{"mobile", maskPhone},
	{"email", maskEmail},
}

// visibleKeySuffixes 必须保持可见的字段后缀，优先于 suffixKeyRules 判定
//
// 这些是排障与可观测性的主线索，一旦被擦掉日志就失去价值单列豁免表（而非仅靠
// 不收录对应后缀）是为了防止后缀表日后扩充时把它们连带卷进去
var visibleKeySuffixes = []string{
	"errorcode", "businesscode", "statuscode", "grpccode", "httpcode", "retcode",
	"ipaddress", "macaddr", "remoteaddress", "localaddress", "peeraddress",
	"publickey", "cachekey", "partitionkey", "idempotencykey", "shardkey",
	"routingkey", "sortkey", "primarykey", "groupkey", "dedupkey",
}

// buildKeyRules 把按动作分组的字段名清单摊平成归一化规则表
func buildKeyRules(groups map[redactAction][]string) map[string]redactAction {
	rules := make(map[string]redactAction)
	for action, keys := range groups {
		for _, k := range keys {
			var buf [maxNormalizedKeyLen]byte
			if n, ok := normalizeKeyInto(buf[:], k); ok {
				rules[string(buf[:n])] = action
			}
		}
	}
	return rules
}

// maxNormalizedKeyLen 栈上归一化缓冲长度，超长字段名退回仅后缀判定
const maxNormalizedKeyLen = 64

// normalizeKeyInto 把字段名归一化（剔除 _ - . 分隔符、ASCII 转小写）写入 dst
// 返回有效长度与是否成功；dst 容纳不下时返回 false
//
// 归一化让 accessToken / access_token / Access-Token / ACCESS.TOKEN 落到同一判定
func normalizeKeyInto(dst []byte, key string) (int, bool) {
	n := 0
	for i := 0; i < len(key); i++ {
		c := key[i]
		if c == '_' || c == '-' || c == '.' {
			continue
		}
		if c >= 'A' && c <= 'Z' {
			c += 'a' - 'A'
		}
		if n >= len(dst) {
			return 0, false
		}
		dst[n] = c
		n++
	}
	return n, true
}

// keyNeedsNormalize 报告字段名是否含需要归一化的字符
// 不含时归一化结果与原串逐字节相同，可省去整趟归一化与第二次查表
func keyNeedsNormalize(key string) bool {
	for i := 0; i < len(key); i++ {
		c := key[i]
		if c == '_' || c == '-' || c == '.' || (c >= 'A' && c <= 'Z') {
			return true
		}
	}
	return false
}

// hasNormalizedSuffix 判断字段名归一化后是否以 suffix 结尾（suffix 须已是归一化形态）
// 从尾部边跳分隔符边转小写地比较，无需先构造归一化串，故零分配
func hasNormalizedSuffix(key, suffix string) bool {
	j := len(suffix) - 1
	for i := len(key) - 1; i >= 0 && j >= 0; i-- {
		c := key[i]
		if c == '_' || c == '-' || c == '.' {
			continue
		}
		if c >= 'A' && c <= 'Z' {
			c += 'a' - 'A'
		}
		if c != suffix[j] {
			return false
		}
		j--
	}
	return j < 0
}

// ============================================================================
// 字段名判定：两级索引
// ============================================================================
//
// 朴素实现（查一次精确表，再遍历全部后缀规则）在未命中的字段名上成本最高——
// 而生产日志里绝大多数字段名都不敏感，等于每条日志都要为每个字段白跑完整流程
// 实测遍历 40+ 条后缀规则单字段就要 57 ns，五个字段把一条 KV 日志从 140 ns 抬到 426 ns
//
// 两级索引把这部分压到常数级：
//  1. 形状位图：按(首字符, 归一化长度)索引精确表中是否存在同形状的键
//     region 归一化后首字符 r、长度 6，而精确表里 r 开头的键长度均 ≥10，
//     一次数组访问即可跳过整个精确表查询
//  2. 后缀字典树：按归一化字符反向回溯，一次走完即可排除所有不可能的后缀，而非逐条规则各做一次尾部扫描

// 形状位图维度：首字符与末字符各 27 档（a-z 各一档，其余归入末档），长度 0..63
const (
	shapeBuckets  = 27
	maxShapeLen   = 63
	otherShapeIdx = 26
)

// exactKeyShapes 标记精确表中是否存在某形状的字段名
//
// 形状取(首字符, 末字符, 归一化长度)三元组，长度以 uint64 的位存放。
// 只用(首字符, 长度)二元组时误命中偏多——cache_key 与 csrftoken 同为 c 开头、
// 归一化长度 8，于是每个 cache_key 都要白做一次完整归一化加一次查表（实测
// 单字段 8.2 ns 升到 17.4 ns）。补上末字符后二者分别落在 c/y 与 c/n 两格，
// ip_address（i/s）、error_code（e/e）、idempotency_key（i/y）同样全部落空
var exactKeyShapes [shapeBuckets][shapeBuckets]uint64

// suffixTrie 后缀字典树的根（按归一化字符反向构建）
//
// 末字符分桶仍不够：*_key 这一桶有 16 条规则、*_code 有 11 条，每个字段名都要
// 把桶内规则逐条做尾部扫描反向字典树让一次回溯同时排除所有不可能的后缀——
// region 反向走两步（n→o）即无路可走，cache_key 走五步即可区分 cachekey 与 licensekey
//
// 匹配取最长后缀：后缀越长语义越具体运维豁免后缀（如 statuscode、cachekey）
// 天然比与之竞争的敏感后缀更长，因此无需额外的优先级规则
var suffixTrie = &suffixNode{}

// suffixNode 字典树节点，children 按 'a'-'z' 直接索引（后缀全为小写字母）
type suffixNode struct {
	children [26]*suffixNode
	action   redactAction
	terminal bool
}

// insertSuffix 把一条后缀规则反向插入字典树
func insertSuffix(suffix string, action redactAction) {
	node := suffixTrie
	for i := len(suffix) - 1; i >= 0; i-- {
		c := suffix[i]
		if c < 'a' || c > 'z' {
			return // 规则表只收录纯小写字母后缀
		}
		idx := c - 'a'
		if node.children[idx] == nil {
			node.children[idx] = &suffixNode{}
		}
		node = node.children[idx]
	}
	node.terminal, node.action = true, action
}

// matchSuffix 沿字典树反向回溯字段名，返回最长匹配后缀对应的动作
// 边回溯边跳过 _ - . 并转小写，无需构造归一化串，故零分配
//
// 命中敏感后缀后再查运维豁免表：豁免只在命中时才有意义，而未命中的字段名
// （生产里的绝大多数）由此完全不必为豁免判定付代价
func matchSuffix(key string) redactAction {
	node := suffixTrie
	action := keepAsIs
	for i := len(key) - 1; i >= 0; i-- {
		c := key[i]
		if c == '_' || c == '-' || c == '.' {
			continue
		}
		if c >= 'A' && c <= 'Z' {
			c += 'a' - 'A'
		}
		if c < 'a' || c > 'z' {
			break
		}
		node = node.children[c-'a']
		if node == nil {
			break
		}
		if node.terminal {
			action = node.action // 继续回溯，更长的后缀会覆盖
		}
	}
	if action != keepAsIs && isOpsVisibleKey(key) {
		return keepAsIs
	}
	return action
}

// isOpsVisibleKey 报告字段名是否属于必须保持可见的运维字段
//
// 单列豁免表（而非仅靠不收录对应敏感后缀）是为了防止后缀表日后扩充时
// 把排障主线索连带擦掉——那种失败不报错，只是日志静默失去价值
func isOpsVisibleKey(key string) bool {
	for _, suffix := range visibleKeySuffixes {
		if hasNormalizedSuffix(key, suffix) {
			return true
		}
	}
	return false
}

// init 构建判定所需的两级索引
// 规则表在日志路径上只读，故索引在包初始化时一次建成
func init() {
	for key := range exactKeyRules {
		first, last, n, _ := keyShape(key)
		markKeyShape(first, last, n)
	}
	// 只把敏感后缀入树。运维豁免后缀改为"命中敏感后缀后再查"——
	// 把它们一并入树会显著加深回溯路径（cachekey 要走满 8 步、
	// idempotencykey 14 步），而豁免判定只在真正命中时才有意义
	for _, rule := range suffixKeyRules {
		insertSuffix(rule.suffix, rule.action)
	}
}

// shapeIndex 把首字符映射到形状位图的首维下标
func shapeIndex(first byte) int {
	if first >= 'a' && first <= 'z' {
		return int(first - 'a')
	}
	return otherShapeIdx
}

// shapeLenIndex 把归一化长度映射到位掩码的位号（超长统一归入末位）
func shapeLenIndex(n int) uint {
	if n > maxShapeLen {
		return maxShapeLen
	}
	return uint(n)
}

// markKeyShape 在形状位图中标记一个字段名的形状
func markKeyShape(first, last byte, n int) {
	exactKeyShapes[shapeIndex(first)][shapeIndex(last)] |= 1 << shapeLenIndex(n)
}

// hasKeyShape 报告精确表中是否存在该形状的字段名
func hasKeyShape(first, last byte, n int) bool {
	return exactKeyShapes[shapeIndex(first)][shapeIndex(last)]>>shapeLenIndex(n)&1 != 0
}

// keyShape 单趟扫描取出字段名的归一化形状
// 返回归一化后的首字符、末字符、长度，以及是否含需归一化的字符
func keyShape(key string) (first, last byte, n int, needsNorm bool) {
	for i := 0; i < len(key); i++ {
		c := key[i]
		if c == '_' || c == '-' || c == '.' {
			needsNorm = true
			continue
		}
		if c >= 'A' && c <= 'Z' {
			c += 'a' - 'A'
			needsNorm = true
		}
		if n == 0 {
			first = c
		}
		last = c
		n++
	}
	return first, last, n, needsNorm
}

// actionForKey 判定字段名对应的脱敏动作（零分配）
//
// Go 对 m[string(byteSlice)] 有专门优化，归一化重查时的 string 转换不产生堆分配
func actionForKey(key string) redactAction {
	first, last, n, needsNorm := keyShape(key)
	if n == 0 {
		return keepAsIs
	}

	// 形状位图挡掉绝大多数字段名，命中才真正查精确表
	if hasKeyShape(first, last, n) {
		if !needsNorm {
			if action, ok := exactKeyRules[key]; ok {
				return action
			}
		} else if n <= maxNormalizedKeyLen {
			var buf [maxNormalizedKeyLen]byte
			if m, ok := normalizeKeyInto(buf[:], key); ok {
				if action, hit := exactKeyRules[string(buf[:m])]; hit {
					return action
				}
			}
		}
	}

	return matchSuffix(key)
}

// isCredentialKey 报告字段名是否为需整体擦除的凭据（供 URL 查询参数判定复用）
func isCredentialKey(key string) bool {
	return actionForKey(key) == wipeValue
}

// ============================================================================
// 掩码
// ============================================================================

// appendMaskedPhone 追加号码掩码："***" + 末 4 位数字；数字不足 5 位时整体擦除
// 只取数字位，故 +86 138-0013-8000 与 13800138000 掩码结果一致
func appendMaskedPhone(buf []byte, s string) []byte {
	digits := 0
	for i := 0; i < len(s); i++ {
		if s[i] >= '0' && s[i] <= '9' {
			digits++
		}
	}
	if digits <= 4 {
		return append(buf, redacted...)
	}
	buf = append(buf, '*', '*', '*')
	skip := digits - 4
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c < '0' || c > '9' {
			continue
		}
		if skip > 0 {
			skip--
			continue
		}
		buf = append(buf, c)
	}
	return buf
}

// appendMaskedEmail 追加邮箱掩码：保留首字符与完整域名
// 域名是排障所需且不含个人信息，故完整保留
func appendMaskedEmail(buf []byte, s string) []byte {
	at := strings.LastIndexByte(s, '@')
	if at <= 0 || at == len(s)-1 {
		return append(buf, redacted...)
	}
	// 首字符按 rune 取，避免把多字节字符截断成非法 UTF-8
	local := s[:at]
	width := 1
	if local[0] >= 0x80 {
		for width < len(local) && local[width]&0xc0 == 0x80 {
			width++
		}
	}
	buf = append(buf, local[:width]...)
	buf = append(buf, '*', '*', '*', '@')
	return append(buf, s[at+1:]...)
}

// ============================================================================
// Logger 侧接入
// ============================================================================

// WithSanitize 设置脱敏强度（默认 SanitizeKeys）
func (l *Logger) WithSanitize(mode SanitizeMode) *Logger {
	l.sanitize = mode
	return l
}

// RegisterSensitiveKeys 追加自定义敏感字段名（归一化后参与精确匹配，整体擦除）
// 全进程生效：敏感字段清单属于安全基线，不应因 Logger 实例不同而不同
// 应在启动阶段调用完毕，规则表在日志路径上只读
func RegisterSensitiveKeys(keys ...string) {
	for _, k := range keys {
		var buf [maxNormalizedKeyLen]byte
		n, ok := normalizeKeyInto(buf[:], k)
		if !ok || n == 0 {
			continue
		}
		normalized := string(buf[:n])
		exactKeyRules[normalized] = wipeValue
		// 形状位图必须同步：actionForKey 以它为前置过滤，
		// 漏更新会让新注册的字段名永远查不到精确表
		markKeyShape(normalized[0], normalized[n-1], n)
	}
}

// appendSanitizedJSONValue 按字段名脱敏后追加 JSON 值
// 脱敏关闭或字段名无需处理时走原路径，零额外开销
func (l *Logger) appendSanitizedJSONValue(buf []byte, key string, v any) []byte {
	if l.sanitize == SanitizeOff {
		return appendJSONValue(buf, v)
	}

	switch actionForKey(key) {
	case wipeValue:
		return appendJSONString(buf, redacted)
	case maskPhone:
		return appendMaskedJSON(buf, v, appendMaskedPhone)
	case maskEmail:
		return appendMaskedJSON(buf, v, appendMaskedEmail)
	}

	if l.sanitize == SanitizeKeysAndValues {
		if s, ok := textOf(v); ok {
			return appendScannedJSONString(buf, l.capped(s))
		}
	}
	if l.maxFieldBytes > 0 {
		if s, ok := textOf(v); ok {
			cut, trimmed := l.capFieldText(s)
			return appendCappedJSONString(buf, cut, trimmed)
		}
	}
	return appendJSONValue(buf, v)
}

// appendCappedJSONString 追加（可能已截断的）文本为 JSON 字符串
// 截断标记在此处直接拼进缓冲，避免为"截断值 + 标记"再分配一次字符串
func appendCappedJSONString(buf []byte, s string, trimmed bool) []byte {
	buf = append(buf, '"')
	buf = appendJSONStringContent(buf, s)
	if trimmed {
		buf = appendJSONStringContent(buf, truncatedValueMark)
	}
	return append(buf, '"')
}

// capped 按单字段上限截断文本，超限时追加 truncated 标记
func (l *Logger) capped(s string) string {
	cut, trimmed := l.capFieldText(s)
	if !trimmed {
		return s
	}
	return cut + truncatedValueMark
}

// appendSanitizedTextValue 按字段名脱敏后追加文本值（text 模式）
func (l *Logger) appendSanitizedTextValue(buf []byte, key string, v any) []byte {
	if l.sanitize == SanitizeOff {
		return appendValueSafe(buf, v)
	}

	switch actionForKey(key) {
	case wipeValue:
		return append(buf, redacted...)
	case maskPhone:
		if s, ok := textOf(v); ok {
			return appendMaskedPhone(buf, s)
		}
	case maskEmail:
		if s, ok := textOf(v); ok {
			return appendMaskedEmail(buf, s)
		}
	}

	if l.sanitize == SanitizeKeysAndValues {
		if s, ok := textOf(v); ok {
			return appendScannedText(buf, l.capped(s))
		}
	}
	if l.maxFieldBytes > 0 {
		if s, ok := textOf(v); ok {
			cut, trimmed := l.capFieldText(s)
			buf = appendTextSafe(buf, cut)
			if trimmed {
				buf = append(buf, truncatedValueMark...)
			}
			return buf
		}
	}
	return appendValueSafe(buf, v)
}

// appendSanitizedMessage 在值扫描模式下脱敏 msg 自身
// msg 常由 fmt 拼接而成，凭据与 PII 最容易从这里漏出去
func (l *Logger) appendSanitizedMessage(buf []byte, msg string, asJSON bool) []byte {
	// msg 位于条目骨架内，字段循环的回滚机制管不到它，故在此单独约束：
	// 一个超长 msg（如把整个响应体拼进消息）单独就能撑爆条目
	msg = l.cappedMessage(msg)

	if l.sanitize != SanitizeKeysAndValues {
		if asJSON {
			return appendJSONString(buf, msg)
		}
		return appendTextSafe(buf, msg)
	}
	if asJSON {
		return appendScannedJSONString(buf, msg)
	}
	return appendScannedText(buf, msg)
}

// appendMaskedJSON 对文本形态的值施加掩码函数后作为 JSON 字符串追加；非文本值原样输出
func appendMaskedJSON(buf []byte, v any, mask func([]byte, string) []byte) []byte {
	s, ok := textOf(v)
	if !ok {
		return appendJSONValue(buf, v)
	}
	// 掩码产物由数字、'*'、'@' 与原域名组成，域名不含需 JSON 转义的字符，
	// 故可直接写在引号之间，省去一次中转缓冲
	buf = append(buf, '"')
	buf = mask(buf, s)
	return append(buf, '"')
}

// textOf 取值的文本形态，仅覆盖零成本的几种
// 不为脱敏触发 fmt 兜底：那会让未命中规则的字段也付一次分配
func textOf(v any) (string, bool) {
	switch val := v.(type) {
	case string:
		return val, true
	case []byte:
		return string(val), true
	case error:
		if val == nil {
			return "", false
		}
		return val.Error(), true
	}
	return "", false
}

// textKey 取 KV 参数里键的文本形态，供脱敏判定使用
// 非字符串键不参与脱敏判定（返回空串必然不命中任何规则），避免为此触发 fmt 分配
func textKey(k any) string {
	if s, ok := k.(string); ok {
		return s
	}
	return ""
}

// cappedMessage 按整行与单字段上限约束 msg
// 预留骨架（时间戳、级别、caller、键名）与字段的空间，故只给 msg 分配上限的一半
func (l *Logger) cappedMessage(msg string) string {
	limit := l.maxFieldBytes
	if l.maxLineBytes > 0 {
		if half := l.maxLineBytes / 2; limit <= 0 || half < limit {
			limit = half
		}
	}
	if limit <= 0 || len(msg) <= limit {
		return msg
	}
	return msg[:utf8SafeEnd(msg, limit)] + truncatedValueMark
}

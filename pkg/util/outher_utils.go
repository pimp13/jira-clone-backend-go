package util

import (
	"crypto/rand"
	"math/big"
	"strings"
)

const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789_-"

// GenerateInviteCode تولید یک کد دعوت تصادفی امن با طول دلخواه
// (پیش‌فرض ۶۴ کاراکتر)
func GenerateInviteCode(length int) string {
	if length <= 0 {
		length = 64
	}

	const charLen = len(charset) // 64

	// برای طول‌های کوچک (< ۲۰۰-۳۰۰) این روش سریع‌تر و خواناتر است
	if length < 300 {
		b := make([]byte, length)
		for i := range b {
			idx, _ := rand.Int(rand.Reader, big.NewInt(int64(charLen)))
			b[i] = charset[idx.Int64()]
		}
		return string(b)
	}

	// برای طول‌های خیلی زیاد (به ندرت استفاده می‌شود)
	var sb strings.Builder
	sb.Grow(length)

	for i := 0; i < length; i++ {
		idx, _ := rand.Int(rand.Reader, big.NewInt(int64(charLen)))
		sb.WriteByte(charset[idx.Int64()])
	}

	return sb.String()
}

// نسخه ساده‌تر و معمولاً کافی برای اکثر موارد
func GenerateInviteCodeSimple(length int) string {
	if length <= 0 {
		length = 64
	}

	b := make([]byte, length)
	_, _ = rand.Read(b) // فقط پر کردن بایت‌ها

	for i := range b {
		b[i] = charset[b[i]&63] // 63 = 2⁶-۱ → ۶ بیت → ۶۴ کاراکتر
	}

	return string(b)
}

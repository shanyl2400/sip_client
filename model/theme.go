package model

const (
	ThemeSend     = "send"
	ThemeRecevice = "receive"
	ThemeInfo     = "info"
	ThemeWarn     = "warn"
	ThemeError    = "error"

	ThemeAll         = "all"
	ThemeSendRecv    = "send/recv"
	ThemeTransaction = "transaction"
)

type Theme string

func (t Theme) Filter(t2 Theme) bool {
	if t == t2 || t == ThemeAll {
		return true
	}
	if t == ThemeSendRecv && (t2 == ThemeSend || t2 == ThemeRecevice) {
		return true
	}
	if t == ThemeTransaction && (t2 == ThemeInfo || t2 == ThemeWarn || t2 == ThemeError) {
		return true
	}

	return false
}

module github.com/LQR471814/nu_plugin_caldav

go 1.24.3

require (
	github.com/ainvaltin/nu-plugin v0.0.0-20250907111918-1d43779b9a0f
	github.com/emersion/go-ical v0.0.0-20250609112844-439c63cef608
	github.com/emersion/go-webdav v0.7.0
	github.com/google/uuid v1.6.0
	github.com/shibukawa/configdir v0.0.0-20170330084843-e180dbdc8da0
	github.com/teambition/rrule-go v1.8.2
	github.com/thlib/go-timezone-local v0.0.7
	github.com/zeebo/xxh3 v1.0.2
	modernc.org/sqlite v1.40.1
)

require (
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/klauspost/cpuid/v2 v2.3.0 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/ncruces/go-strftime v0.1.9 // indirect
	github.com/remyoudompheng/bigfft v0.0.0-20230129092748-24d4a6f8daec // indirect
	github.com/stretchr/testify v1.7.1 // indirect
	github.com/vmihailenco/msgpack/v5 v5.4.1 // indirect
	github.com/vmihailenco/tagparser/v2 v2.0.0 // indirect
	golang.org/x/exp v0.0.0-20250620022241-b7579e27df2b // indirect
	golang.org/x/sys v0.38.0 // indirect
	modernc.org/libc v1.66.10 // indirect
	modernc.org/mathutil v1.7.1 // indirect
	modernc.org/memory v1.11.0 // indirect
)

replace github.com/ainvaltin/nu-plugin v0.0.0-20250907111918-1d43779b9a0f => github.com/LQR471814/nu-plugin v0.0.0-20251218180218-0f31d4ec9708

replace github.com/emersion/go-webdav v0.7.0 => github.com/LQR471814/go-webdav v0.0.0-20251218033631-3be4a3e33dec

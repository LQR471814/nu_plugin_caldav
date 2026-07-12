module github.com/LQR471814/nu_plugin_caldav

go 1.25.0

require (
	github.com/ainvaltin/nu-plugin v0.0.0-20260711090100-edde0f88c6d9
	github.com/emersion/go-ical v0.0.0-20250609112844-439c63cef608
	github.com/emersion/go-webdav v0.7.0
	github.com/google/uuid v1.6.0
	github.com/shibukawa/configdir v0.0.0-20170330084843-e180dbdc8da0
	github.com/teambition/rrule-go v1.8.2
	github.com/thlib/go-timezone-local v0.0.8
	github.com/zeebo/xxh3 v1.1.0
	modernc.org/sqlite v1.53.0
)

require (
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/klauspost/cpuid/v2 v2.4.0 // indirect
	github.com/mattn/go-isatty v0.0.22 // indirect
	github.com/ncruces/go-strftime v1.0.0 // indirect
	github.com/remyoudompheng/bigfft v0.0.0-20230129092748-24d4a6f8daec // indirect
	github.com/stretchr/testify v1.7.1 // indirect
	github.com/vmihailenco/msgpack/v5 v5.4.1 // indirect
	github.com/vmihailenco/tagparser/v2 v2.0.0 // indirect
	golang.org/x/sys v0.47.0 // indirect
	golang.org/x/tools v0.48.0 // indirect
	modernc.org/libc v1.74.1 // indirect
	modernc.org/mathutil v1.7.1 // indirect
	modernc.org/memory v1.11.0 // indirect
)

replace github.com/emersion/go-webdav v0.7.0 => github.com/LQR471814/go-webdav v0.0.0-20251218033631-3be4a3e33dec

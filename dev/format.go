package main

import "strconv"

func niceTimeFormat(seconds int) string {
	var time string = ""
	if seconds > 60 {
		var minutes int = seconds / 60
		seconds %= 60
		if minutes > 60 {
			var hours int = minutes / 60
			minutes = minutes % 60
			if hours > 24 {
				var days int = hours / 24
				hours = hours % 24
				time = strconv.Itoa(days) + "d "
			}
			time = time + strconv.Itoa(hours) + "h "
		}
		time = time + strconv.Itoa(minutes) + "m "
	}
	time = time + strconv.Itoa(seconds) + "s ago"
	return time
}

func niceSizeFormat(bytes int) string {
	var size string = ""
	switch {
	case bytes > 1024*1024*1024:
		size = strconv.Itoa(bytes/1024/1024/1024) + " GiB"
	case bytes > 1024*1024:
		size = strconv.Itoa(bytes/1024/1024) + " MiB"
	case bytes > 1024:
		size = strconv.Itoa(bytes/1024) + " KiB"
	default:
		size = strconv.Itoa(bytes) + " Bytes"
	}
	return size
}

package commontools

//左边补全
func PadLeft(src string, size int, sp string) string {
	var sb StringBuffer
	sb = StringBuffer{}
	i := len(src)
	for i < size {
		sb.WriteString(sp)
		i = i + len(sp)
	}
	sb.WriteString(src)
	r := sb.String()
	return r[len(r)-size:]
}

//右边边补全
func RightLeft(src string, size int, sp string) string {
	var sb StringBuffer
	sb = StringBuffer{}
	i := len(src)
	sb.WriteString(src)
	for i < size {
		sb.WriteString(sp)
		i = i + len(sp)
	}
	r := sb.String()
	return r[0:size]
}


// Source: https://golang.org/src/net/dnsclient.go
// IsDomainName checks if a string is a presentation-format domain name
// (currently restricted to hostname-compatible "preferred name" LDH labels and
// SRV-like "underscore labels"; see golang.org/issue/12421).
export default (s: string): boolean => {
	// See RFC 1035, RFC 3696.
	// Presentation format has dots before every label except the first, and the
	// terminal empty label is optional here because we assume fully-qualified
	// (absolute) input. We must therefore reserve space for the first and last
	// labels' length octets in wire format, where they are necessary and the
	// maximum total length is 255.
	// So our _effective_ maximum is 253, but 254 is not rejected if the last
	// character is a dot.
	const l = s.length
	if (l === 0 || l > 254 || l === 254 && s[l-1] != '.') {
		return false
	}

	let last = '.'
	let ok = false // Ok once we've seen a letter.
	let partlen = 0
	for (let i = 0; i < s.length; i++) {
    const c = s[i]
    if ('a' <= c && c <= 'z' || 'A' <= c && c <= 'Z' || c === '_') {
      ok = true
			partlen++
    } else if ('0' <= c && c <= '9') {
      // fine
			partlen++
    } else if (c === '-') {
      // Byte before dash cannot be dot.
			if (last === '.') {
				return false
			}
			partlen++
    } else if (c === '.') {
      // Byte before dot cannot be dot, dash.
			if (last === '.' || last === '-') {
				return false
			}
			if (partlen > 63 || partlen === 0) {
				return false
			}
			partlen = 0
    } else {
      return false
    }
		last = c
	}
	if (last === '-' || partlen > 63) {
		return false
	}

	return ok
}

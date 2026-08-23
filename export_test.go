package apttransports3go

var InitLogLevel = initLogLevel
var Read = read
var ReadLine = readLine
var Send = send

// UnregisterStatus temporarily removes a status from the registry so that
// send() fails for it, making the error paths of its callers reachable.
// The returned function restores the entry.
func UnregisterStatus(s Status) func() {
	desc := statusByCode[s]
	delete(statusByCode, s)

	return func() {
		statusByCode[s] = desc
	}
}

package repeater

func SendReq2LocalAndRcvRsp(dns string, b []byte) []byte {
	return SendAndRcv(dns, 53, b)
}

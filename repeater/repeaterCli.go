package repeater

import (
	"github.com/shankusu2017/repeaterDNS/config"
	"github.com/shankusu2017/utils"
)

func SendReq2OutsideAndRcvRsp(b []byte) []byte {
	// 加密
	cipherText := utils.AESCrypt(b, config.GetIV16(), config.GetKey16())

	// 发送
	ip, port := config.GetRepeaterSrvAddr()
	cipherText = SendAndRcv(ip, port, cipherText)

	// 解密
	plainText, _ := utils.AESDeCrypt(cipherText, config.GetIV16(), config.GetKey16())
	return plainText[:]
}

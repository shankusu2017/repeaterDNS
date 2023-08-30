package repeater

import (
	"github.com/shankusu2017/repeaterDNS/config"
	"github.com/shankusu2017/utils"
	"log"
)

func SendReq2OutsideAndRcvRsp(b []byte) []byte {
	if config.DebugFlag {
		log.Printf("DEBUG 8da5ed91 send req to repeaterSrv, b.len:%d\n", len(b))
	}
	// 加密
	cipherText := utils.AESCrypt(b, config.GetIV16(), config.GetKey16())

	// 发送
	ip, port := config.GetRepeaterSrvAddr()
	cipherText = SendAndRcv(ip, port, cipherText)

	if config.DebugFlag {
		log.Printf("DEBUG 3d9ccae2 rcv rsp(%d) from repeaterSrv\n", len(cipherText))
	}

	// 解密
	plainText, _ := utils.AESDeCrypt(cipherText, config.GetIV16(), config.GetKey16())
	return plainText[:]
}

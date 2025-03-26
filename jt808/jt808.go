package jt808

import (
	"fmt"
	"strings"
)

// parseJT808Message 解析 JT808 报文
func ParseJT808Message(message string) map[string]interface{} {
	message = preprocessMessage(message)
	var result = make(map[string]interface{})
	// 去掉起始和结束标志
	if len(message) >= 4 && message[:2] == "7E" && message[len(message)-2:] == "7E" {
		message = message[2 : len(message)-2]
	} else {
		return result
	}
	msgID := "0x" + message[:4]
	if MessageMap[msgID] == ""{
		return result
	}

	result["消息ID"] = msgID
	result["消息类型"] = "位置信息汇报"
	//fmt.Println("消息ID", msgID)
	msgStyle := message[4:8]
	binaryMsgStyle, _ := HexToBinary(msgStyle)
	if len(binaryMsgStyle) < 16 {
		for i := len(binaryMsgStyle); i < 16; i++ {
			binaryMsgStyle = "0" + binaryMsgStyle
		}
	}
	// 计算消息体长度
	msgBodyLen, _ := BinaryToDecimal(binaryMsgStyle[6:])
	fmt.Println("消息体长度", binaryMsgStyle[6:], msgBodyLen)

	// 终端手机号
	terminalPhone := message[8:20]
	//fmt.Println("终端手机号", terminalPhone)
	result["终端手机号"] = terminalPhone

	// 消息流水号
	flowerMsgSerialNum, _ := HexToDecimal(message[20:24])
	//fmt.Println("消息流水号", flowerMsgSerialNum)
	result["消息流水号"] = fmt.Sprintf("%d", flowerMsgSerialNum)

	// 位置信息汇报
	locationData := message[24 : 24+msgBodyLen]
	//fmt.Println("位置信息汇报", locationData)
	// result["位置信息汇报"] = locationData
	_ = locationData
	// 报警标志
	alarmFlag, _ := HexToDecimal(message[24:32])
	//fmt.Println("报警标志", alarmFlag)
	result["报警标志"] = fmt.Sprintf("%d", alarmFlag)
	// 状态
	status, _ := HexToDecimal(message[32:40])
	//fmt.Println("状态", status)
	result["状态"] = fmt.Sprintf("%d", status)

	// 经度
	longitude, _ := HexToDecimal(message[40:48])
	//fmt.Println("经度", longitude)
	result["经度"] = float64(longitude) / 1000000
	// 纬度
	latitude, _ := HexToDecimal(message[48:56])
	//fmt.Println("纬度", latitude)
	result["纬度"] = float64(latitude) / 1000000
	// 海拔
	altitude, _ := HexToDecimal(message[56:60])
	//fmt.Println("海拔", altitude)
	result["海拔"] = altitude
	// 速度
	speed, _ := HexToDecimal(message[60:64])
	//fmt.Println("速度", speed)
	result["速度"] = float64(speed) * 0.1
	// 方向
	direction, _ := HexToDecimal(message[64:68])
	//fmt.Println("方向", direction)
	result["方向"] = direction
	// 时间
	timeStr := message[68:80]
	//fmt.Println("时间", fmt.Sprintf("%s-%s-%s %s:%s:%s", timeStr[:2], timeStr[2:4], timeStr[4:6], timeStr[6:8], timeStr[8:10], timeStr[10:12]))
	result["时间"] = fmt.Sprintf("%s-%s-%s %s:%s:%s", timeStr[:2], timeStr[2:4], timeStr[4:6], timeStr[6:8], timeStr[8:10], timeStr[10:12])
	return result
}






// 平台通用应答消息
func ServerCommonReplyMessage(message string) string {
	// fmt.Println("平台通用应答消息 = 应答流水号（终端流水号） + 应答ID（终端消息ID）")
	return message[20:24] + message[:4] + "00"
}





// 报文预处理
// 匹配多种传输格式
// 去除多余字符
// 处理转义字符
// 替换特殊字符
func preprocessMessage(message string) string {
	message = strings.ReplaceAll(message, "|", "")
	message = strings.ReplaceAll(message, " ", "")
	message = strings.ReplaceAll(message, "0x", "")
	message = strings.ReplaceAll(message, "7D02", "7E")
	message = strings.ReplaceAll(message, "7D01", "7D")
	message = strings.ToUpper(message)
	return message
}


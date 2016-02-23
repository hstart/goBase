package mac

/*
包功能说明:
常用的mac计算算法
*/

import (
	"bytes"
	"errors"
	"fmt"
	"goBase/basebytes"
	"goBase/des"
)

func init() {
	fmt.Println("Simple Main Test @ goBase/mac")
}

//mac算法ANSI9.19 算法 :
func MAC_ANSI919(data, key []byte) ([]byte, error) {

	//16个字节key值，8个字节左key，8字节右边key
	lkey := make([]byte, 8)
	rkey := make([]byte, 8)
	copy(lkey, key[0:8])
	copy(rkey, key[8:16])

	//初始8个字节的0
	//iv := bytes.Repeat([]byte{byte(0)}, 8)
	result := make([]byte, 8)
	dataNo := make([]byte, 8)
	dataNew := des.ZeroPadding(data, 8) //按照8个字节的长度进行补 0

	if len(dataNew)%8 != 0 {
		fmt.Println("Error @ 长度 ", dataNew)
		return nil, errors.New("待加密长度异常")
	}

	for i := 0; i < len(dataNew)/8; i++ {
		copy(dataNo, dataNew[i*8:i*8+8])
		//xor 结果
		for j := 0; j < 8; j++ {
			result[j] = result[j] ^ dataNo[j]
		}
		//Des加密
		ret, err := des.DesEncrypt(result, lkey)
		if err != nil {
			return nil, err
		}
		copy(result, ret[0:8])
	}

	//fmt.Println("Mac result: ", basebytes.BytesToHex(result, len(result)))

	//Des解密
	retN, err := des.DesDecrypt(result, rkey)
	if err != nil {
		return nil, err
	}
	copy(result, retN[0:8])

	//Des加密
	retN2, err := des.DesEncrypt(result, lkey)
	if err != nil {
		return nil, err
	}
	copy(result, retN2[0:8])
	return result, nil
}

func TestMAC_ANSI919() {
	//32位的key值
	//A246DF34516B2520A1979BD6D989FE43
	key := []byte{0xA2, 0x46, 0xDF, 0x34, 0x51, 0x6B, 0x25, 0x20, 0xA1, 0x97, 0x9B, 0xD6, 0xD9, 0x89, 0xFE, 0x43}
	data := bytes.Repeat([]byte{byte(0)}, 8)
	data2 := bytes.Repeat([]byte{byte(0x11)}, 6)
	data3 := make([]byte, 16)
	copy(data3, data)
	copy(data3[8:], data2)
	fmt.Println("Mac 数据: ", basebytes.BytesToHex(data3, len(data3)))
	fmt.Println("Mac Key: ", basebytes.BytesToHex(key, len(key)))

	d4 := "3C3F786D6C2076657273696F6E3D22312E302220656E636F64696E673D2249534F2D383835392D3122203F3E3C696E3E3C686561643E3C56657273696F6E3E312E302E313C2F56657273696F6E3E3C496E737449643E424A4345423C2F496E737449643E3C416E735472616E436F64653E424A4345425142495265713C2F416E735472616E436F64653E3C54726D5365714E756D3E32303136303232323030373333303C2F54726D5365714E756D3E3C2F686561643E3C74696E3E3C62696C6C4B65793E313030303131333230303031323739333838333C2F62696C6C4B65793E3C636F6D70616E7949643E3032353030383330323C2F636F6D70616E7949643E3C626567696E4E756D3E313C2F626567696E4E756D3E3C71756572794E756D3E313C2F71756572794E756D3E3C66696C6564313E3C2F66696C6564313E3C66696C6564323E3C2F66696C6564323E3C66696C6564333E3C2F66696C6564333E3C66696C6564343E3C2F66696C6564343E3C2F74696E3E3C2F696E3E"

	data4, _ := basebytes.HexToBytes(d4)

	ret, err := MAC_ANSI919(data4, key)
	if err != nil {
		return
	}
	fmt.Println("Mac 结果: ", basebytes.BytesToHex(ret, len(ret)))
}

package basebytes

/*
goBase 基础处理函数包:字节类相关处理
*/
import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"strconv"
)

func init() {
	fmt.Println("Simple Main Test @ goBase/bytes")
}

func pass() {

}

//16进制转换为内存字节类型数据 如字符串"41" --> 'A' 单字符
func HexToBytes(str string) ([]byte, error) {
	if len(str)%2 != 0 {
		return nil, errors.New("HexToBytes @ 长度非偶数")
	}
	retbytes := make([]byte, len(str)/2)

	for j := 0; j < len(str); j++ {
		if (str[j] >= '0' && str[j] <= '9') || (str[j] >= 'A' && str[j] <= 'F') || (str[j] >= 'a' && str[j] <= 'f') {
			pass()
		} else {
			return nil, errors.New("HexToBytes @ 非法字符")
		}
	}
	for i := 0; i < len(str)/2; i++ {
		num := str[i*2 : i*2+2]
		//fmt.Println(num)
		v, err := strconv.ParseUint(num, 16, 8)
		if err != nil {
			return nil, err
		}
		//fmt.Println("New ", v)
		retbytes[i] = byte(v)

	}
	return retbytes, nil
}

/* 这个函数功能错误
func HexToBytes(str string) ([]byte, error) {
	retbytes := make([]byte, len(str))
	for i := 0; i < len(str); i++ {
		switch {
		case str[i] >= '0' && str[i] <= '9':
			retbytes[i] = str[i] - '0'
		case str[i] >= 'a' && str[i] <= 'z':
			retbytes[i] = str[i] - 'a' + 10
		case str[i] >= 'A' && str[i] <= 'Z':
			retbytes[i] = str[i] - 'A' + 10
		default:
			return nil, errors.New(fmt.Sprintf("invalid hex character: %c", str[i]))
		}
	}
	return retbytes, nil
}
*/

//字节类型转换字符类型表示,如 'AA'-->"4141"
func BytesToHex(buff []byte, blen int) string {

	retstr := ""
	for i := 0; i < blen; i++ {
		bb := fmt.Sprintf("%02X", buff[i])
		retstr = retstr + bb
	}
	return retstr
}

//大端 将字节类型转换为int16的整形数据
func BytesToInt16(buffer []byte) int16 {
	b_buf := bytes.NewBuffer(buffer)
	var x int16
	binary.Read(b_buf, binary.BigEndian, &x)
	//fmt.Println("转换后的结果为,", x, BytesToHex(buffer, 2))
	return x
}

func Int16ToBytes(x int16) []byte {

	b_buf := bytes.NewBuffer([]byte{})
	binary.Write(b_buf, binary.BigEndian, x)
	fmt.Println(b_buf.Bytes())
	return b_buf.Bytes()
}

func TestBytes() {
	dd := "0123456789ABCDEF0123456789ABCDE"
	TT, _ := HexToBytes(dd)
	fmt.Println("DD=", TT)
	fmt.Println("OKK=", BytesToHex(TT, len(TT)))
}

package des

/*
 gobase基础组件:
 DES和3DES 加解密算法
*/

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"fmt"
	"goBase/basebytes"
)

func init() {
	fmt.Println("Simple Main Test @ goBase/des")
}

//当长度非8字节倍数时，补0函数
func ZeroPadding(ciphertext []byte, blockSize int) []byte {
	//当是8个字节的倍数时,不补充数据
	padding := (blockSize - len(ciphertext)%blockSize) % blockSize
	if padding == 0 {
		return ciphertext
	}
	padtext := bytes.Repeat([]byte{byte(0)}, padding)
	//fmt.Println("padding=", padtext)
	return append(ciphertext, padtext...)
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := (blockSize - len(ciphertext)%blockSize) % blockSize
	//fmt.Println("padding=", padding)
	if padding == 0 {
		return ciphertext
	}
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	//fmt.Println("padding=", padtext)
	return append(ciphertext, padtext...)
}

//DES加密算法基础算法，每次加密8个字节
func DesEncrypt(origData, key []byte) ([]byte, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	//origData = PKCS5Padding(origData, block.BlockSize())
	//Des加密算法采用补0的方式
	origData = ZeroPadding(origData, block.BlockSize())

	//fmt.Println("origData", basebytes.BytesToHex(origData, 8), block.BlockSize())
	//fmt.Println("Key", basebytes.BytesToHex(key, 8))

	//iv向量采用全0为
	iv := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	blockMode := cipher.NewCBCEncrypter(block, iv)
	crypted := make([]byte, len(origData))
	// 根据CryptBlocks方法的说明，如下方式初始化crypted也可以
	// crypted := origData
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

//DES解密算法
func DesDecrypt(crypted, key []byte) ([]byte, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}

	iv := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	//blockMode := cipher.NewCBCDecrypter(block, key)
	blockMode := cipher.NewCBCDecrypter(block, iv)
	origData := make([]byte, len(crypted))
	// origData := crypted
	blockMode.CryptBlocks(origData, crypted)
	//origData = PKCS5UnPadding(origData)
	// origData = ZeroUnPadding(origData)
	return origData, nil
}

/*
下面是3DES加解密算法，注意加解密算法都有一个基础算法和外面最终的加解密算法
Des3Decrypt 多字节解密算法(8个字节的循环调用Des3DecryptBase 进行处理)
Des3DecryptBase 8个字节解密算法(只适用于8个字节解密，多字节的话得出结果不正确)
Des3Encrypt 多字节加密算法(8个字节的循环调用Des3EncryptBase 进行处理)
Des3EncryptBase 8个字节的加密算法
*/
//3DES解密算法
func Des3Decrypt(crypted, key []byte) ([]byte, error) {

	//fmt.Println("Des3Decrypt @  解密数据 ", basebytes.BytesToHex(crypted, len(crypted)), len(crypted))
	retData := make([]byte, len(crypted))
	Data := make([]byte, 8)
	for i := 0; i < len(crypted)/8; i++ {
		copy(Data, crypted[i*8:i*8+8])
		//fmt.Println("Des3Decrypt @ 需要进行解密数据", basebytes.BytesToHex(Data, len(Data)))
		//每次解密8个字节的数据
		origData, err := Des3DecryptBase(Data, key)
		if err != nil {
			fmt.Println("Des3Decrypt @ ERROR ", err)
			return nil, err
		}
		copy(retData[i*8:], origData[0:8])
	}
	//fmt.Println("Des3Decrypt @  解密结果 ", basebytes.BytesToHex(retData, len(retData)), len(retData))
	return retData, nil
}

//3DES解密基础算法
func Des3DecryptBase(crypted, key []byte) ([]byte, error) {
	//key是24个字节的，而输入的key16个字节的,通过补key
	newkey := make([]byte, 24)
	copy(newkey, key)
	copy(newkey[16:], key[0:8])
	block, err := des.NewTripleDESCipher(newkey)
	if err != nil {
		fmt.Println("Des3Decrypt ", err)
		return nil, err
	}
	iv := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	//blockMode := cipher.NewCBCDecrypter(block, key)
	blockMode := cipher.NewCBCDecrypter(block, iv)
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	//fmt.Println("Des3DecryptBase @ ", BytesToHex(origData, len(origData)), len(origData))
	return origData, nil
}

//3DES加密算法
func Des3Encrypt(cryptedsrc, key []byte) ([]byte, error) {
	crypted := ZeroPadding(cryptedsrc, 8) //按照8个字节的长度倍数进行补0
	//fmt.Println("Des3Encrypt @  加密数据 ", basebytes.BytesToHex(crypted, len(crypted)), len(crypted))
	retData := make([]byte, len(crypted))
	Data := make([]byte, 8)
	for i := 0; i < len(crypted)/8; i++ {
		copy(Data, crypted[i*8:i*8+8])
		//fmt.Println("Des3Encrypt @ 需要进行加密数据", basebytes.BytesToHex(Data, len(Data)))
		//每次解密8个字节的数据
		origData, err := Des3EncryptBase(Data, key) //加密处理
		if err != nil {
			fmt.Println("Des3Encrypt @ ERROR ", err)
			return nil, err
		}
		copy(retData[i*8:], origData[0:8])
	}

	//fmt.Println("Des3Encrypt @  加密结果 ", basebytes.BytesToHex(retData, len(retData)), len(retData))
	return retData, nil
}

//3DES加密算法，一次8个字节加密算法
func Des3EncryptBase(origData, key []byte) ([]byte, error) {

	newkey := make([]byte, 24)
	copy(newkey, key)
	copy(newkey[16:], key[0:8])
	block, err := des.NewTripleDESCipher(newkey)
	//block, err := des.NewTripleDESCipher(key)
	if err != nil {
		return nil, err
	}
	//origData = PKCS5Padding(origData, block.BlockSize())
	origData = ZeroPadding(origData, block.BlockSize())
	// origData = ZeroPadding(origData, block.BlockSize())
	//blockMode := cipher.NewCBCEncrypter(block, key[:8])

	//iv 向量全部为0
	blockMode := cipher.NewCBCEncrypter(block, bytes.Repeat([]byte{byte(0)}, block.BlockSize()))

	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

func TestDes() {
	fmt.Println("goBase/des 包测试")

	fmt.Println("\nTestDes @ Des加密处理 ")
	data := []byte{0x6C, 0x6D, 0xF8, 0x3D, 0x2F, 0xA1, 0xC6, 0xF6}
	key := []byte{0xA2, 0x46, 0xDF, 0x34, 0x51, 0x6B, 0x25, 0x20}

	//result := []byte{}

	fmt.Println("TestDes @ Des待加密数据 ", basebytes.BytesToHex(data, len(data)), len(data))
	fmt.Println("TestDes @ Des待加密Key ", basebytes.BytesToHex(key, len(key)), len(key))

	//结果应该为
	ret, err := DesEncrypt(data, key)
	if err != nil {
		fmt.Println("Des @ 加密处理,异常", err)
		return
	}
	fmt.Println("TestDes @ Des加密结果", basebytes.BytesToHex(ret, len(ret)), len(ret))

	fmt.Println("\nTestDes @ Des解密处理 ")

	deData := []byte{0x02, 0x10, 0x66, 0x00, 0x39, 0x85, 0x3A, 0x26}
	fmt.Println("TestDes @ Des待解密数据 ", basebytes.BytesToHex(deData, len(deData)), len(deData))
	fmt.Println("TestDes @ Des待解密Key ", basebytes.BytesToHex(key, len(key)), len(key))
	//结果应该为
	deret, err := DesDecrypt(deData, key)
	if err != nil {
		fmt.Println("Des @ 解密处理,异常", err)
		return
	}
	fmt.Println("TestDes @ Des解密结果", basebytes.BytesToHex(deret, len(deret)), len(deret))

	fmt.Println("\nTestDes @ 3Des解密处理 ")
	key3 := []byte{0x01, 0x23, 0x45, 0x67, 0x89, 0xAB, 0xCD, 0xEF, 0x01, 0x23, 0x45, 0x67, 0x89, 0xAB, 0xCD, 0xEF}
	data3 := []byte{0xF1, 0xC6, 0xF9, 0xDF, 0xB4, 0x11, 0x5B, 0xCF, 0xFD, 0x7E, 0x75, 0xB4, 0x20, 0x9E, 0x93, 0xA4}
	fmt.Println("TestDes @ 3Des待解密数据 ", basebytes.BytesToHex(data3, len(data3)), len(data3))
	fmt.Println("TestDes @ 3Des待解密Key ", basebytes.BytesToHex(key3, len(key3)), len(key3))

	deret3, err := Des3Decrypt(data3, key3)
	if err != nil {
		fmt.Println("Des @ 解密处理,异常", err)
		return
	}

	fmt.Println("TestDes @ 3Des解密结果", basebytes.BytesToHex(deret3, len(deret3)), len(deret3))

	fmt.Println("\nTestDes @ 3Des加密处理 ")
	endata3 := []byte{0xA2, 0x46, 0xDF, 0x34, 0x51, 0x6B, 0x25, 0x20, 0xA1, 0x97, 0x9B, 0xD6, 0xD9, 0x89, 0xFE, 0x43, 0x01}
	fmt.Println("TestDes @ 3Des待加密数据 ", basebytes.BytesToHex(endata3, len(endata3)), len(endata3))
	fmt.Println("TestDes @ 3Des待加密Key ", basebytes.BytesToHex(key3, len(key3)), len(key3))

	enret3, err := Des3Encrypt(endata3, key3)
	if err != nil {
		fmt.Println("3Des @ 加密处理,异常", err)
		return
	}
	fmt.Println("TestDes @ 3Des加密结果", basebytes.BytesToHex(enret3, len(enret3)), len(enret3))
	return

}


//此源码被清华学神尹成大魔王专业翻译分析并修改
//尹成QQ77025077
//尹成微信18510341407
//尹成所在QQ群721929980
//尹成邮箱 yinc13@mails.tsinghua.edu.cn
//尹成毕业于清华大学,微软区块链领域全球最有价值专家
//https://mvp.microsoft.com/zh-cn/PublicProfile/4033620
//版权所有（c）2015-2016 BTCSuite开发者
//此源代码的使用由ISC控制
//可以在许可文件中找到的许可证。

package btcec

import (
	"bytes"
	"encoding/hex"
	"testing"
)

func TestGenerateSharedSecret(t *testing.T) {
	privKey1, err := NewPrivateKey(S256())
	if err != nil {
		t.Errorf("private key generation error: %s", err)
		return
	}
	privKey2, err := NewPrivateKey(S256())
	if err != nil {
		t.Errorf("private key generation error: %s", err)
		return
	}

	secret1 := GenerateSharedSecret(privKey1, privKey2.PubKey())
	secret2 := GenerateSharedSecret(privKey2, privKey1.PubKey())

	if !bytes.Equal(secret1, secret2) {
		t.Errorf("ECDH failed, secrets mismatch - first: %x, second: %x",
			secret1, secret2)
	}
}

//测试1：加密和解密
func TestCipheringBasic(t *testing.T) {
	privkey, err := NewPrivateKey(S256())
	if err != nil {
		t.Fatal("failed to generate private key")
	}

	in := []byte("Hey there dude. How are you doing? This is a test.")

	out, err := Encrypt(privkey.PubKey(), in)
	if err != nil {
		t.Fatal("failed to encrypt:", err)
	}

	dec, err := Decrypt(privkey, out)
	if err != nil {
		t.Fatal("failed to decrypt:", err)
	}

	if !bytes.Equal(in, dec) {
		t.Error("decrypted data doesn't match original")
	}
}

//测试2：与py椭圆的字节兼容性
func TestCiphering(t *testing.T) {
	pb, _ := hex.DecodeString("fe38240982f313ae5afb3e904fb8215fb11af1200592b" +
		"fca26c96c4738e4bf8f")
	privkey, _ := PrivKeyFromBytes(S256(), pb)

	in := []byte("This is just a test.")
	out, _ := hex.DecodeString("b0d66e5adaa5ed4e2f0ca68e17b8f2fc02ca002009e3" +
		"3487e7fa4ab505cf34d98f131be7bd258391588ca7804acb30251e71a04e0020ecf" +
		"df0f84608f8add82d7353af780fbb28868c713b7813eb4d4e61f7b75d7534dd9856" +
		"9b0ba77cf14348fcff80fee10e11981f1b4be372d93923e9178972f69937ec850ed" +
		"6c3f11ff572ddd5b2bedf9f9c0b327c54da02a28fcdce1f8369ffec")

	dec, err := Decrypt(privkey, out)
	if err != nil {
		t.Fatal("failed to decrypt:", err)
	}

	if !bytes.Equal(in, dec) {
		t.Error("decrypted data doesn't match original")
	}
}

func TestCipheringErrors(t *testing.T) {
	privkey, err := NewPrivateKey(S256())
	if err != nil {
		t.Fatal("failed to generate private key")
	}

	tests1 := []struct {
ciphertext []byte //输入密文
	}{
{bytes.Repeat([]byte{0x00}, 133)},                   //err输入短路
{bytes.Repeat([]byte{0x00}, 134)},                   //ERRUNSUPPORTED曲线
{bytes.Repeat([]byte{0x02, 0xCA}, 134)},             //错误无效长度
{bytes.Repeat([]byte{0x02, 0xCA, 0x00, 0x20}, 134)}, //错误无效长度
{[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, //Ⅳ
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
0x02, 0xCA, 0x00, 0x20, //曲线和X长度
0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, //X
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
0x00, 0x20, //Y长
0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, //Y
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, //密文
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, //雨衣
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
}}, //无效的PUBKEY
{[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, //Ⅳ
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
0x02, 0xCA, 0x00, 0x20, //曲线和X长度
0x11, 0x5C, 0x42, 0xE7, 0x57, 0xB2, 0xEF, 0xB7, //X
			0x67, 0x1C, 0x57, 0x85, 0x30, 0xEC, 0x19, 0x1A,
			0x13, 0x59, 0x38, 0x1E, 0x6A, 0x71, 0x12, 0x7A,
			0x9D, 0x37, 0xC4, 0x86, 0xFD, 0x30, 0xDA, 0xE5,
0x00, 0x20, //Y长
0x7E, 0x76, 0xDC, 0x58, 0xF6, 0x93, 0xBD, 0x7E, //Y
			0x70, 0x10, 0x35, 0x8C, 0xE6, 0xB1, 0x65, 0xE4,
			0x83, 0xA2, 0x92, 0x10, 0x10, 0xDB, 0x67, 0xAC,
			0x11, 0xB1, 0xB5, 0x1B, 0x65, 0x19, 0x53, 0xD2,
0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, //密文
//填充未对齐到16个字节
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, //雨衣
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
}}, //错误无效填充
{[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, //Ⅳ
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
0x02, 0xCA, 0x00, 0x20, //曲线和X长度
0x11, 0x5C, 0x42, 0xE7, 0x57, 0xB2, 0xEF, 0xB7, //X
			0x67, 0x1C, 0x57, 0x85, 0x30, 0xEC, 0x19, 0x1A,
			0x13, 0x59, 0x38, 0x1E, 0x6A, 0x71, 0x12, 0x7A,
			0x9D, 0x37, 0xC4, 0x86, 0xFD, 0x30, 0xDA, 0xE5,
0x00, 0x20, //Y长
0x7E, 0x76, 0xDC, 0x58, 0xF6, 0x93, 0xBD, 0x7E, //Y
			0x70, 0x10, 0x35, 0x8C, 0xE6, 0xB1, 0x65, 0xE4,
			0x83, 0xA2, 0x92, 0x10, 0x10, 0xDB, 0x67, 0xAC,
			0x11, 0xB1, 0xB5, 0x1B, 0x65, 0x19, 0x53, 0xD2,
0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, //密文
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, //雨衣
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
}}, //无故障MAC
	}

	for i, test := range tests1 {
		_, err = Decrypt(privkey, test.ciphertext)
		if err == nil {
			t.Errorf("Decrypt #%d did not get error", i)
		}
	}

//removepkcspadding的测试错误
	tests2 := []struct {
in []byte //输入数据
	}{
		{bytes.Repeat([]byte{0x11}, 17)},
		{bytes.Repeat([]byte{0x07}, 15)},
	}
	for i, test := range tests2 {
		_, err = removePKCSPadding(test.in)
		if err == nil {
			t.Errorf("removePKCSPadding #%d did not get error", i)
		}
	}
}

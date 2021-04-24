# pem，等保存密钥的文件看不懂？不知道怎么生成？不知道怎么恢复？
```
-----BEGIN RSA PRIVATE KEY-----
3/k2OLWwpTzaTYeHPbyYEEGFuVLEh86gm+fBt72DkiA=
-----END RSA PRIVATE KEY-----

```
或许，可能，大概，假装你遇到过上面这样的文件内容，
你（没错，就是你，不是我）在网上随便找了一个生成密钥文件的代码，
把自己的私钥放上去运行一下，诶！居然没报错，你乐开了花，私钥成功保存在文件里面了。

然而快乐总是短暂的，你打开文件后（得到上面的内容），
一脸疑惑："这是什么鬼乱七八糟，八糟乱七的东西啊？"。这个是加密过后的？还是程序 bug ？
关键之前找的代码也没有从私钥文件恢复私钥的代码啊， 即使找到了包含写和读的完整代码，但是不同代码保存成的文件内容有可能也不一样。
也就是说想要恢复私钥，就得用回原来那个代码，如果忘了那个代码，私钥就恢复不了了？

乖乖，本来是为了方便，不想记 64 位私钥才选择记录在文件里面的。
现在，我不用记 64 位的十六进制的字符了，我只需要记 64 行的代码就行，造孽啊。

又本着闲着也是闲着的原则，我决定好好研究一下这个到底是怎么回事，
让自己的私钥存得放心，取得安心~

## pem 文件的编码格式
其实呢，

    pem 文件里面所存储的二进制数据是用 base64 格式储存的

可能这个时候，有些动手能力强的帅哥美女就把我上面的字符串拿到浏览器上，
进行在线 base64 解码了，想看看我的私钥。然而，想看我私钥没那么简单，
估计你看到的是这样（https://base64.us/ 去试试）：

![img.png](decode1.png)

诶~，诶~，别走，别骂人，你要相信我，这并不是错了。听我慢慢到来……

我上面说过："pem 文件里面所存储的二进制数据是用 base64 格式储存的",
也就是说我们往 pem 文件存的是二进制，然后 pem 文件用 base64 编码格式把它展示给我们看。
如果直接对"3/k2OLWwpTzaTYeHPbyYEEGFuVLEh86gm+fBt72DkiA="
这个字符串进行 base64 解码，得到一串二进制数，
那么它一定和我之前存进去的二进制串试一样的（我保证，我发誓）。

那么为什么上面用浏览器在线解码会得到乱码呢？

那是因为同一串二进制，可能会有许许多多的编码方式，有 base58，base64 ，
十六进制，Ascill……等方式，而在上面的 base64 在线解码网页中，
所得到的二进制串是使用 Ascill 编码的，可是我往 pem 文件中所存储的字节，
是使用十六进制格式编码的啊（没想到吧【鬼魅一笑】），你用 Ascill 的方式编码原本使用十六进制编码的二进制，这样……不合适吧……

如果想通过 3/k2OLWwpTzaTYeHPbyYEEGFuVLEh86gm+fBt72DkiA= 还原我原来使用十六进制编码的字符串，你需要这样：
```
    "3/k2OLWwpTzaTYeHPbyYEEGFuVLEh86gm+fBt72DkiA="->base64 解码->二进制->十六进制编码
```
我给你们尝试一下(https://c.runoob.com/compile/21 golang在线执行）：
```
package main

import (
	"fmt"
	"encoding/base64"
	"encoding/hex"
)

func decode(data string) {
	// base64 解码->二进制
	sDec, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		fmt.Println("Error decoding string: %s ", err.Error())
		return
	}

	// 二进制->十六进制编码
	pri:=hex.EncodeToString(sDec)
	fmt.Println(pri)//dff93638b5b0a53cda4d87873dbc98104185b952c487cea09be7c1b7bd839220
}

func main() {
	decode("3/k2OLWwpTzaTYeHPbyYEEGFuVLEh86gm+fBt72DkiA=")
}
```
然后就会输出 dff93638b5b0a53cda4d87873dbc98104185b952c487cea09be7c1b7bd839220 
没错，这就是我保存的私钥。恭喜恭喜，历尽千辛万苦终于的到私钥了，
快去看看这个私钥对应的地址有没有比特币……

## 完整代码
保存私钥成文件，和把文件恢复成私钥的完整代码可以看这里:
## 密钥文件保存懂了
知道了 pem 文件保存的内容是怎么回事，其他也没有什么大事了。以后需要把私钥保存到文件
里面，只需要记住私钥编码格式，基本也就没问题了，

这就是我对 pem 文件保存数据的理解啦~





package chaoren

import (
	"io/ioutil"
	"syscall"
	"unsafe"
)

func abort(funcname string, err error) {
	panic(funcname + " failed: " + err.Error())
}

func Send() {
	//帐号配置 http://www.chaorendama.com 注册帐号和密码
	var username, password string
	username = "zj2055" //超人云用户名
	password = "qq819458"
	softId := 75996 //软件id,开发者帐号后台添加
	//加载dll
	h, err := syscall.LoadLibrary("dcb64.dll")
	if err != nil {
		abort("LoadLibrary", err)
	}
	defer syscall.FreeLibrary(h)

	//初始化插件,连接服务器
	DC_Init, err := syscall.GetProcAddress(h, "DC_Init")
	if err != nil {
		abort("DC_Init", err)
	}
	dch, _, _ := syscall.Syscall(uintptr(DC_Init), 3, uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr("127.0.0.1"))), 14888, 50)

	//查询点数
	left := int32(0)
	today := int32(0)
	total := int32(0)
	DC_GetInfo, err := syscall.GetProcAddress(h, "DC_GetInfo")
	r, _, _ := syscall.Syscall6(uintptr(DC_GetInfo), 6,
		uintptr(dch),
		uintptr(unsafe.Pointer(syscall.StringBytePtr(username))),
		uintptr(unsafe.Pointer(syscall.StringBytePtr(password))),
		uintptr(unsafe.Pointer(&left)),
		uintptr(unsafe.Pointer(&today)),
		uintptr(unsafe.Pointer(&total)))
	//uintptr指针变量转整型
	rr := int32(r)
	if rr == 1 {
		print("left:", left, "\n")
		print("today:", today, "\n")
		print("total:", total, "\n")
	} else if rr == -5 {
		print("帐号或密码错误")
	}
	//上传图片，等待识别完成
	imagePath := `C:\Users\Administrator\Desktop\3f43a46caeef48aba329f3561aa29be4.jpg`
	buf, _ := ioutil.ReadFile(imagePath)

	DC_RecogImg, err := syscall.GetProcAddress(h, "DC_RecogImg")
	imageId := make([]byte, 32)
	result := make([]byte, 33)
	lenght := len(buf)
	print("len", lenght, "\n")
	r1, _, _ := syscall.Syscall9(uintptr(DC_RecogImg), 9,
		uintptr(dch),
		uintptr(unsafe.Pointer(syscall.StringBytePtr(username))),
		uintptr(unsafe.Pointer(syscall.StringBytePtr(password))),
		uintptr(softId),
		uintptr(60),
		uintptr(lenght),
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(unsafe.Pointer(&imageId[0])),
		uintptr(unsafe.Pointer(&result[0])))
	//uintptr指针变量转整型
	rrr := int32(r1)
	if rrr == 1 {
		print("success:")
		print(string(result), "  ", string(imageId), "\n")
	} else if rrr == 0 {
		print("识别超时")
	} else if rrr == -1 {
		print("识别失败,超时或者没有传放正确的参数或其它原因")
	} else if rrr == -2 {
		print("余额不足")
	} else if rrr == -3 {
		print("未绑定或者未在绑定机器上运行")
	} else if rrr == -4 {
		print("时间过期")
	} else if rrr == -5 {
		print("用户校验失败")
	}

	//报告识别错误,仅在识别结果错误时调用
	/*
	   DC_NotifyFail, err := syscall.GetProcAddress(h, "DC_NotifyFail")
	   syscall.Syscall6(uintptr(DC_NotifyFail), 4,
	       uintptr(dch),
	       uintptr(unsafe.Pointer(syscall.StringBytePtr(username))),
	       uintptr(unsafe.Pointer(syscall.StringBytePtr(password))),
	       uintptr(unsafe.Pointer(&imageId[0])), 0, 0)
	   print("done")
	*/

	//反初始化插件
	DC_Uninit, err := syscall.GetProcAddress(h, "DC_Uninit")
	syscall.Syscall(uintptr(DC_Uninit), 1, uintptr(dch), 0, 0)
}

package kuaishibie

import (
	"os"
	"testing"
)

func TestReq(t *testing.T) {
	SetUsername("wujibing", "wujibing")
	p := &Predict{
		TypeId:    "20",
		Image:     "",
		ImageBack: "",
		Content:   "其呈乌",
	}
	fp, err := os.Open(`C:\Users\Administrator\Desktop\3f43a46caeef48aba329f3561aa29be4.jpg`)
	if err != nil {
		t.Fatal(err)
	}
	defer fp.Close()
	if err = p.OpenImage(fp); err != nil {
		t.Fatal(err)
	}
	t.Log(Req(p))
}

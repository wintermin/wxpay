package wxpay

import (
	"fmt"
	"testing"
	"time"
)

type TestType string

type TestObject struct {
	id int
}

func (o TestObject) String() string {
	return fmt.Sprintf("%d", o.id)
}

type TestInfo struct {
	Name     string     `map:"name"`
	TestType TestType   `map:"test_type"`
	Id       TestObject `map:"id"`
	Remark   string     `map:"remark,omitempty"`
	Remark2  string     `map:"remark2,omitempty"`
	Password string     `map:"-"`
	Start    time.Time  `map:"start"`
}

func TestValues_ForStruct(t *testing.T) {
	val := make(Values)
	val.ForStruct(&TestInfo{
		Name:     "min",
		TestType: "1",
		Id: TestObject{
			id: 100,
		}, Password: "123456",
		Remark: "bak",
		Start:  time.Now(),
	})
	fmt.Println(val)
	body := val.Encode()
	fmt.Println(string(body))
	newVal := make(Values)
	newVal.Decode(body)
	fmt.Println(newVal)
}

func TestValues_ToStruct(t *testing.T) {
	val := make(Values)
	val.Add("return_code", "SUCCESS")
	val.Add("appid", "123456")
	val.Add("total_fee", "100")
	val.Add("sign_type", "MD5")
	val.Add("coupon_id_0", "0")
	val.Add("coupon_id_1", "1")
	val.Add("coupon_type_0", "CASH")
	val.Add("coupon_type_1", "NO_CASH")
	resp := &OrderQueryResponse{}
	val.ToStruct(resp)
	fmt.Println(resp)
}

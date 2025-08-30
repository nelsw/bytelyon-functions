package db

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestClient_Save(t *testing.T) {
	//j := bot.Job{
	//	ID: id.NewULID(),
	//}
	//
	//err := Save(j)
	//if err != nil {
	//	t.Fail()
	//}
	//ctx = context.Background()
	//context.WithValue(ctx, "id", j.ID)
	//u := model.User{
	//	ID: id,
	//}
	//j := model.Job{
	//	User: u,
	//	ID:   model.NewUlid(),
	//}
	//b, err := json.Marshal(&w)
	//fmt.Println(string(b), err)

}

func TestClient_Load(t *testing.T) {
	s1 := "db/user/01K3WG7VQ5VJ16VWS3XGW3138G"
	s2 := "/job/01K3WG7VQ51R1030Q44GGVWF06"
	//s3 := "/work/01K3WG7VQ5NTMHS0J2VPRD16TJ/data.json"
	var a any
	err := Get(s1+s2, &a)
	if err != nil {
		t.Fatal(err)
	}
	b, _ := json.MarshalIndent(a, "", "  ")
	fmt.Println(string(b))
}

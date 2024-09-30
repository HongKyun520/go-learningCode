package tencent

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
	"testing"
)

const secret_key = ""

const secret_id = ""

func TestSender(t *testing.T) {

	c, err := sms.NewClient(common.NewCredential(secret_id, secret_key), "", profile.NewClientProfile())

	if err != nil {
		t.Fatal(err)
	}

	// 构建service
	s := NewService(c, "123", "123")

	parm := Test{
		name:    "test",
		tplId:   "123",
		params:  []string{"123"},
		numbers: []string{"123"},
	}

	t.Run(parm.name, func(t *testing.T) {
		er := s.Send(context.Background(), parm.tplId, parm.params, parm.params...)
		assert.Equal(t, parm.wantErr, er)
	})

}

type Test struct {
	name    string
	tplId   string
	params  []string
	numbers []string
	wantErr error
}

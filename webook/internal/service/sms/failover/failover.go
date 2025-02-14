package failover

import (
	"GoInAction/webook/internal/service/sms"
	"context"
	"sync/atomic"
	"time"
)

type FailoverSMSService struct {
	// 服务列表
	svcs []sms.Service

	// 当前使用的服务的下标
	idx uint64

	// 超时时间
	timeout time.Duration

	// 连续失败次数
	cnt int32

	// 最大连续失败次数
	maxCnt int32
}

// Send 发送短信
// 如果一个服务连续失败超过阈值,则切换到下一个服务
// ctx 上下文
// tplId 短信模板ID
// args 短信模板参数
// numbers 接收短信的手机号列表
// 返回error,如果发送失败则返回错误
func (f *FailoverSMSService) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	// 获取当前使用的服务下标
	idx := atomic.LoadUint64(&f.idx)
	// 获取当前服务
	svc := f.svcs[idx%uint64(len(f.svcs))]

	// 设置超时控制
	ctx, cancel := context.WithTimeout(ctx, f.timeout)
	defer cancel()

	err := svc.Send(ctx, tplId, args, numbers...)
	if err != nil {
		// 连续失败次数+1 需要原子操作
		cnt := atomic.AddInt32(&f.cnt, 1)
		// 如果连续失败次数超过阈值,切换到下一个服务
		if cnt >= f.maxCnt {
			// 重置连续失败次数
			atomic.StoreInt32(&f.cnt, 0)
			// 切换到下一个服务
			atomic.AddUint64(&f.idx, 1)
		}
		return err
	}
	// 发送成功,重置连续失败次数
	atomic.StoreInt32(&f.cnt, 0)
	return nil
}

func (f *FailoverSMSService) SendV1(ctx context.Context, tplId string, args []string, numbers ...string) error {
	// 获取当前使用的服务下标
	idx := atomic.LoadUint64(&f.idx)

	len := len(f.svcs)

	svc := f.svcs[idx%uint64(len)]

	// 使用当前
	err := svc.Send(ctx, tplId, args, numbers...)

	if err != nil {

		// 连续失败次数+1 需要原子操作
		cnt := atomic.AddInt32(&f.cnt, 1)
		// 如果连续失败次数超过阈值,切换到下一个服务
		if cnt >= f.maxCnt {
			// 重置连续失败次数
			atomic.StoreInt32(&f.cnt, 0)
			// 切换到下一个服务
			atomic.AddUint64(&f.idx, 1)
		}
		return err

	}
	// 发送成功,重置连续失败次数
	atomic.StoreInt32(&f.cnt, 0)
	return nil
}

// NewFailoverSMSService 创建一个新的故障转移短信服务
// svcs 服务列表
// timeout 超时时间
// maxCnt 最大连续失败次数
func NewFailoverSMSService(svcs []sms.Service, timeout time.Duration, maxCnt int32) *FailoverSMSService {
	return &FailoverSMSService{
		svcs:    svcs,
		timeout: timeout,
		maxCnt:  maxCnt,
	}
}

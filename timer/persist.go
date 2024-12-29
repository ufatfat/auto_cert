package timer

import (
	"auto_cert/config"
	"auto_cert/manager"
	"encoding/json"
	"os"
	"time"
)

// Load 加载已被持久化的定时任务
func (t *timerManager) Load() {
	// 启动时先读取是否有持久化的定时任务
	f, err := os.ReadFile("persist.json")
	if err != nil {
		if !os.IsNotExist(err) {
			panic(err)
		}
	} else {
		if err = json.Unmarshal(f, &persistInfo); err != nil {
			panic(err)
		}
		for dn, ts := range persistInfo {
			go func() {
				// 删除临时timer
				delete(persistInfo, dn)

				domain := manager.NewDomain(dn)

				// 休眠直到指定时间后重新申请证书
				interval := ts - time.Now().Unix()
				time.Sleep(time.Second * time.Duration(interval))
				if obtainErr := domain.Obtain(); obtainErr != nil {
					// todo: 错误处理
				}

				// 创建新定时任务
				// 固定1800h（75d）之后刷新证书
				entryID, err := t.Job.AddFunc("@every 1800h", func() {
					if renewErr := domain.Renew(); renewErr != nil {
						// todo: 错误处理
					}
				})
				if err != nil {
					// todo: 错误处理
				}

				t.timers[dn] = &domainTimer{
					domain: domain,
					jobID:  entryID,
				}
			}()
		}
	}
}

func (t *timerManager) Persist(sig int) {
	// 更新定时任务时间或新建
	for dn, timer := range t.timers {
		persistInfo[dn] = t.Job.Entry(timer.jobID).Next.Unix()
	}

	persistData, err := json.Marshal(&persistInfo)
	if err != nil {
		panic(err)
	}

	// 将持久化信息写入文件
	persistFileName := time.Now().Format("persist_20060102150405.json")
	if err = os.WriteFile(config.Persistence.Path+persistFileName, persistData, 0755); err != nil {
		panic(err)
	}
	_ = os.Remove("persist.json")

	// 当前目录下创建symlink
	if err = os.Symlink(config.Persistence.Path+persistFileName, "persist.json"); err != nil {
		panic(err)
	}

	// 程序退出
	os.Exit(sig)
}

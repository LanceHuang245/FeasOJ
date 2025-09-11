package scheduler

import (
	"FeasOJ/internal/utils"
	"FeasOJ/internal/utils/sql"
	"FeasOJ/server/handler"
	"fmt"
	"log"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var sentNotifications = make(map[int]bool)

func ScheduleCompetitionStatus() {
	// 启动任务调度
	scheduler := gocron.NewScheduler(time.Local)

	now := time.Now()
	delay := time.Minute - time.Duration(now.Second())*time.Second

	// 延迟到每分钟的0秒
	time.AfterFunc(delay, func() {
		scheduler.Every(1).Minute().Do(func() {
			log.Println("[FeasOJ] Competition scheduler running:", time.Now())

			// 获取未开始的竞赛
			competitions := sql.GetUpcomingCompetitions()

			for _, competition := range competitions {
				if !sentNotifications[competition.ContestID] {
					durationUntilStart := time.Until(competition.Start_at)

					if durationUntilStart > 0 && durationUntilStart <= time.Minute {
						// 使用AfterFunc精确调度
						time.AfterFunc(durationUntilStart, func() {
							// 获取参与该竞赛的用户
							usersInCompetition := sql.SelectUsersCompetition(competition.ContestID)

							// 发送竞赛开始的消息
							for _, user := range usersInCompetition {
								if user.ContestID == competition.ContestID {
									if client, ok := handler.Clients[fmt.Sprint(user.Uid)]; ok {
										// 获取用户语言
										lang := client.Lang
										tag := language.Make(lang)
										langBundle := utils.InitI18n()
										localizer := i18n.NewLocalizer(langBundle, tag.String())
										message, _ := localizer.Localize(&i18n.LocalizeConfig{
											MessageID: "competition_started",
										})
										client.MessageChan <- message
									}
								}
							}

							// 记录已发送通知
							sentNotifications[competition.ContestID] = true
						})
					}
				}
			}

			// 更新竞赛状态
			if err := sql.UpdateCompetitionStatus(); err != nil {
				log.Println("[FeasOJ] Error updating competition status:", err)
			}

			// 更新题目状态
			if err := sql.UpdateProblemVisibility(); err != nil {
				log.Println("[FeasOJ] Error updating competition's problem status:", err)
			}
		})

		// 启动任务调度
		scheduler.StartBlocking()
	})
}

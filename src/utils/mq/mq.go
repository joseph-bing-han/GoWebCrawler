package mq

import (
	"GoWebCrawler/src/utils/conf"
	"github.com/go-redis/redis/v7"
	"strconv"
	"time"
)

const (
	STREAM = "crawlerMQ"
	GROUP  = "crawlerMQGroup"
)

var (
	_redis       *redis.Client
	consumerName string
)

func init() {
	_redis = redis.NewClient(&redis.Options{
		Addr:     conf.Get("MQ_SERVER", "127.0.0.1:6379"),
		Password: conf.Get("MQ_PASSWORD", ""),
	})

	_redis.XGroupCreateMkStream(STREAM, GROUP, "0")

	consumerName = "consumer_" + strconv.FormatInt(time.Now().Unix(), 10)
}

func Add(value map[string]interface{}) {
	_redis.XAdd(&redis.XAddArgs{
		Stream: STREAM,
		Values: value,
	})
}

func Read() (string, map[string]interface{}, error) {

	//优先检查pending列表，转移超时的任务
	checkPending()

	streams := []string{STREAM, ">"}
	result := _redis.XReadGroup(&redis.XReadGroupArgs{
		Group:    GROUP,
		Consumer: consumerName,
		Streams:  streams,
		Count:    1,
		Block:    time.Duration(24) * time.Hour, //阻塞24小时
		NoAck:    false,
	})

	stream, error := result.Result()
	if error != nil {
		return "", nil, error
	}
	if len(stream[0].Messages) == 0 {
		return "", nil, nil
	}
	message := stream[0].Messages[0]
	return message.ID, message.Values, nil
}

func Ack(ID string) {
	_redis.XAck(STREAM, GROUP, ID)
}

func checkPending() {
	result := _redis.XPending(STREAM, GROUP)
	pending, error := result.Result()
	if error == nil {
		for consumer := range pending.Consumers {
			// 处理不是当前消费者的pending
			if consumer != consumerName {
				consumerResult := _redis.XPendingExt(&redis.XPendingExtArgs{
					Stream:   STREAM,
					Group:    GROUP,
					Start:    "-",
					End:      "+",
					Count:    10,
					Consumer: consumer,
				})
				for _, consumerPending := range consumerResult.Val() {
					during, _ := strconv.Atoi(conf.Get("PENDING_MAX", "10"))
					var messages []string
					// pending超时
					if consumerPending.Idle >= time.Duration(during)*time.Second {
						messages = append(messages, consumerPending.ID)
					}

					if len(messages) > 0 {
						// 转移给当前消费者
						_redis.XClaim(&redis.XClaimArgs{
							Stream:   STREAM,
							Group:    GROUP,
							Consumer: consumerName,
							MinIdle:   time.Duration(during)*time.Second,
							Messages: messages,
						})
					}

				}
			}

		}
	}

}

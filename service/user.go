package service

import (
	"context"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"tic3001-go-server/common/constant"
	"tic3001-go-server/database"
	"tic3001-go-server/redis"
	"time"
)

type userService struct{}

var UserService = new(userService)

type MockUser struct {
	Id        string `json:"id"`
	LastName  string `json:"lastName"`
	FirstName string `json:"firstName"`
	Gender    string `json:"gender"`
	IpAddress string `json:"ipAddress"`
}

func (s userService) FindMockUser() {
	// check if cache has data first
	mockUsers := s.findMockUserCache()
	if len(mockUsers) != 0 {
		// return
		for _, user := range mockUsers {
			fmt.Printf("%v\n", user)
		}
		return
	}

	// query from db
	database.Conn.GetConnection().Table("user").Find(&mockUsers)

	// save to cache with 30s as expiry time
	// redis.Client.LPush(context.Background(), constant.MockUserListRedisKey, mockUsers...)
	// redis.Client.LPush(context.Background(), constant.MockUserListRedisKey, mockUsers[0])
	marshal, _ := json.Marshal(mockUsers[0])
	redis.Client.LPush(context.Background(), constant.MockUserListRedisKey, marshal)
	redis.Client.Expire(context.Background(), constant.MockUserListRedisKey, time.Second*30)
}

func (s userService) MockUserCacheChecker() bool {
	exists := redis.Client.Exists(context.Background(), "app:mock:user:list")
	result, err := exists.Result()
	if err != nil {
		log.Error("error when request redis server")
		return false
	}
	return result > 0
}

func (s userService) findMockUserCache() []MockUser {
	mockUsers := make([]MockUser, 0)
	err := redis.Client.LRange(context.Background(), constant.MockUserListRedisKey, 0, -1).ScanSlice()
	if err != nil {
		log.Error("error when do scan slice")
	}
	return mockUsers
}

func (s userService) GenerateTestData() {
	content, err := ioutil.ReadFile("data.sql")
	if err != nil {
		log.Error(err)
		return
	}
	go func() {
		clause := "BEGIN TRANSACTION;" + string(content) + "COMMIT;"
		database.Conn.GetConnection().Exec(clause)
	}()
}

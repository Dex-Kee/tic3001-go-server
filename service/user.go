package service

import (
	"context"
	json "github.com/goccy/go-json"
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

func (e MockUser) MarshalBinary() ([]byte, error) {
	return json.Marshal(e)
}

func (s userService) FindMockUser() []MockUser {
	// check if cache has data first
	mockUsers := s.findMockUserCache()
	if len(mockUsers) != 0 {
		log.Infof("found cache, return cache data")
		return mockUsers
	}

	// query from db
	database.Conn.GetConnection().Table("user").Find(&mockUsers)

	// cache the data
	s.saveToCache(mockUsers)
	return mockUsers
}

func (s userService) saveToCache(mockUsers []MockUser) {
	// convert to interface slice first, to prevent multiple connection to redis server
	container := make([]interface{}, len(mockUsers))
	for i, e := range mockUsers {
		container[i] = interface{}(e)
	}

	// only one connection is needed
	_, err := redis.Client.LPush(context.Background(), constant.MockUserListRedisKey, container...).Result()
	if err != nil {
		log.Error(err.Error())
		return
	}
	// set cache expiry
	redis.Client.Expire(context.Background(), constant.MockUserListRedisKey, time.Minute*5)
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
	length, err := redis.Client.LLen(context.Background(), constant.MockUserListRedisKey).Result()
	if err != nil {
		log.Error("error when do query the length of list")
		return []MockUser{}
	}

	mockUsers := make([]MockUser, length)
	strings, err := redis.Client.LRange(context.Background(), constant.MockUserListRedisKey, 0, -1).Result()
	if err != nil {
		log.Error("error when do scan slice")
		return mockUsers
	}

	// deserialize to struct
	for i, e := range strings {
		user := MockUser{}
		_ = json.Unmarshal([]byte(e), &user)
		mockUsers[i] = user
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

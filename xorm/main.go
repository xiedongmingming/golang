package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"fmt"
	"os"
	"time"
	"github.com/go-xorm/core"
	"math"
	"adwetec.com/tools/model"
	"adwetec.com/tools/utils"
	"encoding/json"
)

// 说明文档: http://xorm.io/

// 1、支持两种风格的混用
// 支持自动同步表结构、事务、原始SQL语句和ORM操作的混合执行
// 引擎是协程安全的

type ReportFeedCampaign struct {
	Id           int64 // XORM自动自增长
	StartDate    time.Time `xorm:"start_date"`
	EndDate      time.Time `xorm:"end_date"`
	Update       time.Time `xorm:"update"`
	AccountName  string    `xorm:"varchar(255)"`
	CampaignName string    `xorm:"varchar(255)"`
	CampaignId   int64     `xorm:"campaign_id"`
	Impression   int32     `xorm:"impression"`
	Click        int32     `xorm:"click"`
	Cost         float64   `xorm:"cost"`
	Ctr          float64   `xorm:"ctr"`
	Cpc          float64   `xorm:"cpc"`
	Cpm          float64   `xorm:"cpm"`
}

func (this *ReportFeedCampaign) String() string {

	return fmt.Sprintf("id: %d\n"+
		"startdate: %v\n"+
		"enddate: %v\n"+
		"update: %v\n"+
		"accountname: %v\n"+
		"campaignname: %v\n"+
		"campaignid: %v\n"+
		"impression: %v\n"+
		"click: %v\n"+
		"cost: %v\n"+
		"ctr: %v\n"+
		"cpc: %v\n"+
		"cpm: %v\n",
		this.Id,
		this.StartDate,
		this.EndDate,
		this.Update,
		this.AccountName,
		this.CampaignName,
		this.CampaignId,
		this.Impression,
		this.Click,
		this.Cost,
		this.Ctr,
		this.Cpc,
		this.Cpm)

}

type RequestAccount struct {
	AdvertiseId []int64 `json:"adv"`           //
	MediaId     []int64 `json:"adx"`           //
	Status      []int32 `json:"status"`        //
	AccountName string  `json:"searchAccount"` //
	SortOrder   string  `json:"sortOrder"`     //
	Limit       int     `json:"limit"`         // 每页显示数量
	Offset      int     `json:"offset"`        // 每页显示数量
}

type ResponseAccount struct {
	Id           int64  `json:"id"`
	AdvertiserId int64  `json:"advertiser_id"`
	MediaId      int64  `json:"media_id"`
	UserId       int64  `json:"user_id"`
	UserName     string `json:"user_name"`
	Status       int32  `json:"status"`
}
type ReportFeedAdgroup struct {
	Id           int64 // XORM自动自增长
	StartDate    time.Time `xorm:"start_date"`
	EndDate      time.Time `xorm:"end_date"`
	Update       time.Time `xorm:"update"`
	AccountName  string    `xorm:"varchar(255)"`
	CampaignName string    `xorm:"varchar(255)"`
	AdgroupName  string    `xorm:"varchar(255)"`
	AdgroupId    int64     `xorm:"adgroup_id"`
	Impression   int32     `xorm:"impression"`
	Click        int32     `xorm:"click"`
	Cost         float64   `xorm:"cost"`
	Ctr          float64   `xorm:"ctr"`
	Cpc          float64   `xorm:"cpc"`
	Cpm          float64   `xorm:"cpm"`
}
type RequestOverview struct {
	Times [2]string `json:"time"`
}
type Result struct {
	Cost float64
	Date string
}

func Round(f float64, n int) float64 {
	pow10_n := math.Pow10(n)
	return math.Trunc((f+0.5/pow10_n)*pow10_n) / pow10_n
}

type FeedAccount struct {
	Id             int64     `xorm:"id"` // XORM自动自增长
	AccountId      int64     `xorm:"account_id"`
	AccountName    string    `xorm:"account_name"`
	Balance        float64   `xorm:"balance"`
	Budget         float64   `xorm:"budget"`
	UserStat       int32     `xorm:"user_stat"`
	BalancePackage int32     `xorm:"balance_package"`
	UaStatus       int32     `xorm:"ua_status"`
	ValidFlows     string    `xorm:"valid_flows"`
	Update         time.Time `xorm:"update"`
}
type ReportFeedAccount struct {
	Id          int64     `xorm:"id"` // XORM自动自增长
	StartDate   time.Time `xorm:"start_date"`
	EndDate     time.Time `xorm:"end_date"`
	Update      time.Time `xorm:"update"`
	AccountName string    `xorm:"account_name"`
	AccountId   int64     `xorm:"account_id"`
	Impression  int32     `xorm:"impression"`
	Click       int32     `xorm:"click"`
	Cost        float64   `xorm:"cost"`
	Ctr         float64   `xorm:"ctr"`
	Cpc         float64   `xorm:"cpc"`
	Cpm         float64   `xorm:"cpm"`
}
type ReportFeedAccountJoin struct {
	ReportFeedAccount `xorm:"extends"`
	FeedAccount       `xorm:"extends"`
}

func (this *ReportFeedAccountJoin) TableName() string {
	return "report_feed_account"
}


type RequestOverviewList struct {
	AdvertiserId   []int64   `json:"advertiser_id"`     // 广告主的ID数组
	MediaId        []int64   `json:"media_id"`          // 媒体的ID数组
	AccountId      []int64   `json:"media_account_id"`  // 账户的ID数组
	CampaignId     []int64   `json:"media_campaign_id"` // 计划的ID数组
	AdgroupId      []int64   `json:"media_adgroup_id"`  // 单元的ID数组
	CreativeId     []int64   `json:"media_creative_id"` // 创意的ID数组
	AdvertiserType string    `json:"advertiser_type"`   // 广告主类型--全部/自定义
	Statistical    string    `json:"statistical"`       // 统计的类型-->广告主/媒体/账户
	Times          [2]string `json:"time"`              // 请求的时间段
	Kpi            string    `json:"kpi"`               // KPI指标
	SummaryWay     string    `json:"summary_way"`       // 统计单位: 分时、分日、分周、分月
}

func NewRequestOverviewList() *RequestOverviewList {
	return &RequestOverviewList{
		AdvertiserId:   []int64{},
		MediaId:        []int64{},
		AccountId:      []int64{},
		CampaignId:     []int64{},
		AdgroupId:      []int64{},
		CreativeId:     []int64{},
		AdvertiserType: "",
		Statistical:    "",
		Times:          [2]string{"", ""},
		Kpi:            "",
		SummaryWay:     "",
	}
}
func main() {

	al := NewRequestOverviewList()

	json.NewEncoder(os.Stdout).Encode(al)

	//test := 252.95000000000002
	//
	//fmt.Println(Round(test, 2))

	// var dateexp = regexp.MustCompile(`^\d{4}-\d{2}-d{2}$`) // 日期的正则表达式

	//timeStr := "2018-03-21T16:00:00.000Z"
	//
	//t, err := time.Parse("2006-01-02T15:04:05.000Z", timeStr)
	//
	//if err != nil {
	//	fmt.Println(err.Error())
	//} else {
	//	fmt.Println(t.Format("2006-01-02 15:04:05"))
	//}
	//
	//return

	engine, err := xorm.NewEngine("mysql", "rdswetec:FWkKYmmLtEbQzs#kG!vA4f!eYiClhuuu@tcp(192.168.0.18:3306)/adwetec_prod?charset=utf8")

	if err != nil {

		fmt.Println("xorm new engine error " + err.Error())

		os.Exit(1)

	}

	engine.SetMaxIdleConns(10)
	engine.SetMaxOpenConns(100)
	engine.SetLogLevel(core.LOG_DEBUG)

	user := &model.AdwetecAdvertiser{
		Name: "测试",
	}

	_, err = engine.InsertOne(user)

	if err != nil {
		fmt.Println("获取错误: " + err.Error())
	}

	fmt.Println(utils.String(user))

	engine.Get(user)

	fmt.Println(utils.String(user))

	//users := make([]*model.AdwetecUser, 0)

	//users = append(users, &model.AdwetecUser{
	//	Id: 3,
	//})
	//users = append(users, &model.AdwetecUser{
	//	Id: 4,
	//})
	//users = append(users, &model.AdwetecUser{
	//	Id: 5,
	//})

	//count, err := engine.Delete(&model.AdwetecUser{
	//	Id: 7,
	//})
	//
	//if err != nil {
	//	fmt.Println("删除数据报错: ", err.Error())
	//	return
	//}
	//
	//fmt.Println(fmt.Sprintf("影响的行数: %d", count))

	//addusers := make([]*model.AdwetecUser, 0)
	//
	//addusers = append(addusers, &model.AdwetecUser{
	//	Email:      "sujiang@adwetec.com",
	//	Contact:    "苏江",
	//	Department: "技术部",
	//	RoleId:     2,
	//	Status:     1,
	//})
	//addusers = append(addusers, &model.AdwetecUser{
	//	Email:      "libaojian@adwetec.com",
	//	Contact:    "李保健",
	//	Department: "技术部",
	//	RoleId:     3,
	//	Status:     1,
	//})
	//addusers = append(addusers, &model.AdwetecUser{
	//	Email:      "langxuemei@adwetec.com",
	//	Contact:    "郎雪梅",
	//	Department: "技术部",
	//	RoleId:     4,
	//	Status:     1,
	//})
	//addusers = append(addusers, &model.AdwetecUser{
	//	Email:      "zhaomin@adwetec.com",
	//	Contact:    "赵敏",
	//	Department: "技术部",
	//	RoleId:     5,
	//	Status:     1,
	//})

	//count, err := engine.Update(&model.AdwetecUser{
	//	RoleId:     5,
	//}, &model.AdwetecUser{
	//	Id: 6,
	//})
	//
	//if err != nil {
	//	fmt.Println("删除数据报错: ", err.Error())
	//	return
	//}
	//
	//fmt.Println(fmt.Sprintf("影响的行数: %d", count))

	//accounts := make([]*ReportFeedAccountJoin, 0)
	//
	//// engine.Find(&accounts)
	//engine.Join("INNER", "feed_account", "report_feed_account.account_id = feed_account.account_id").Find(&accounts)
	//
	//for _, account := range accounts {
	//
	//	fmt.Println(utils.String(account))
	//}

	// objs := make([]interface{}, 0)

	//location, err := time.LoadLocation("Asia/Shanghai")
	//
	//if err != nil {
	//
	//	os.Exit(1)
	//
	//}

	// date, _ := time.ParseInLocation("2006-01-02 15:04:05", "2018-04-04 00:00:00", location)

	//objs = append(objs, model.ReportFeedAccount{
	//	StartDate:   date,
	//	EndDate:     date,
	//	Update:      time.Now(),
	//	AccountName: "原生-车好多4-8172604",
	//	AccountId:   23627252,
	//	Impression:  int32(9679968),
	//	Click:       int32(78259),
	//	Cost:        148487.21,
	//	Ctr:         0.008084634164079882,
	//	Cpc:         1.897381898567577,
	//	Cpm:         15.33963851946618,
	//})

	//obj := model.ReportFeedAccount{
	//	Id: 590,
	//}

	//
	//feedAccountTb = append(feedAccountTb, model.FeedAccount{
	//	AccountId:      22738921,
	//	AccountName:    "原生-JJ斗地主-8164896",
	//	Balance:        19823.22,
	//	Budget:         0,
	//	UserStat:       2,
	//	BalancePackage: 0,
	//	UaStatus:       1,
	//	ValidFlows:     "1,2",
	//})
	//feedAccountTb = append(feedAccountTb, model.FeedAccount{
	//	AccountId:      0,
	//	AccountName:    "原生-JJ斗地主-8164896",
	//	Balance:        0,
	//	Budget:         0,
	//	UserStat:       0,
	//	BalancePackage: 0,
	//	UaStatus:       0,
	//	ValidFlows:     "",
	//})

	//_, err = engine.Get(&obj)
	//
	//if err != nil {
	//
	//	fmt.Println("xorm new engine error " + err.Error())
	//
	//	os.Exit(1)
	//
	//}
	//
	//fmt.Println(utils.String(&obj))
	//
	//fmt.Println("执行完成")

	//engine, err := xorm.NewEngine("mysql", "root:xieming4243054@tcp(180.76.55.12:3306)/adwetec?charset=utf8")
	//
	//if err != nil {
	//
	//	fmt.Println("xorm new engine error " + err.Error())
	//
	//	os.Exit(1)
	//
	//}
	//
	//engine.SetMaxIdleConns(10)
	//engine.SetMaxOpenConns(100)
	//engine.SetLogLevel(core.LOG_DEBUG)
	//
	//al := &RequestOverview{
	//	Times: [2]string{"2018-03-23 00:00:00", "2018-03-27 00:00:00"},
	//}
	//
	//var args []interface{}
	//
	//var buffer bytes.Buffer
	//
	//if al.Times[0] != "" && al.Times[1] != "" {
	//
	//	if buffer.Len() != 0 {
	//		buffer.WriteString(" and ")
	//	}
	//
	//	buffer.WriteString("`start_date` >= ? and `end_date` <= ?")
	//
	//	args = append(args, al.Times[0])
	//	args = append(args, al.Times[1])
	//
	//}
	//
	//fmt.Println("请求SQL的条件语句: " + buffer.String())
	//
	//var results []*Result
	//
	//err = engine.SQL("select sum(cost) as cost, start_date as date from report_feed_account where start_date >= '2018-03-23 00:00:00' and end_date <= '2018-03-27 00:00:00' group by start_date").Find(&results)
	//
	//if err != nil {
	//
	//	fmt.Println(err.Error())
	//
	//	os.Exit(1)
	//
	//} else {
	//	fmt.Println(fmt.Sprintf("%v", results))
	//	for _, cost := range results {
	//		fmt.Println(fmt.Sprintf("%f", cost.Cost) + "-->" + cost.Date)
	//	}
	//
	//}

	//accounts := make([]*model.ReportFeedAccount, 0)
	//
	//err = engine.Where(buffer.String(), args...).Find(&accounts)
	//
	//if err != nil {
	//
	//	fmt.Println(err.Error())
	//
	//	os.Exit(1)
	//
	//} else {
	//
	//	fmt.Println(len(accounts))
	//
	//}

	//mediasMap := make(map[int64]*model.Media)
	//advertisersMap := make(map[int64]*model.Advertiser)
	//
	//engine, err := xorm.NewEngine("mysql", "root:xieming4243054@tcp(180.76.55.12:3306)/adwetec?charset=utf8")
	//
	//if err != nil {
	//
	//	fmt.Println("xorm new engine error " + err.Error())
	//
	//	os.Exit(1)
	//
	//}
	//
	//engine.SetMaxIdleConns(10)
	//engine.SetMaxOpenConns(100)
	//engine.SetLogLevel(core.LOG_DEBUG)
	//
	//var medias []*model.Media
	//var accs []*model.FeedAccount
	//var advs []*model.Advertiser
	//
	//err = engine.Find(&medias)
	//
	//if err != nil {
	//
	//	fmt.Println("[http server] load medias error " + err.Error())
	//
	//	os.Exit(0)
	//
	//}
	//
	//for _, m := range medias {
	//	mediasMap[m.Id] = m
	//}
	//
	//err = engine.Find(&advs)
	//
	//if err != nil {
	//
	//	fmt.Println("[http server] load advertisers error " + err.Error())
	//
	//	os.Exit(0)
	//
	//}
	//
	//for _, a := range advs {
	//	advertisersMap[a.Id] = a
	//}
	//
	//err = engine.Find(&accs)
	//
	//if err != nil {
	//
	//	fmt.Println("[http server] load accounts error " + err.Error())
	//
	//	os.Exit(0)
	//
	//}
	//
	//for _, acc := range accs {
	//
	//	for _, adv := range advertisersMap {
	//
	//		for _, adx := range mediasMap {
	//
	//			item := &model.AccountMapping{
	//				AdvertiserId: adv.Id,
	//				MediaId:      adx.Id,
	//				UserId:       acc.UserId,
	//				UserName:     acc.AccountName,
	//				Status:       int32(acc.UaStatus),
	//			}
	//
	//			engine.InsertOne(item)
	//		}
	//	}
	//
	//}

	//count, err := engine.In("advertiser_id", 2).Count(&model.AccountMapping{})
	//
	//fmt.Println(fmt.Sprintf("数据总量: %d", count))
	//
	//var al RequestAccount
	//
	//al.AdvertiseId = []int64{1, 2, 3}
	//al.MediaId = []int64{}
	//al.Status = []int32{}
	//al.AccountName = ""
	//al.SortOrder = ""
	//al.Offset = 1
	//al.Limit = 10
	//
	//// accounts := make([]*model.AccountMapping, 0)
	//
	//var buffer bytes.Buffer
	//
	//var args []interface{}
	//
	//var adargs []interface{}
	//
	//if len(al.AdvertiseId) != 0 {
	//
	//	for _, a := range al.AdvertiseId {
	//
	//		adargs = append(adargs, a)
	//
	//	}
	//
	//}
	//
	//if len(al.MediaId) != 0 {
	//
	//	if buffer.Len() != 0 {
	//		buffer.WriteString(" and ")
	//	}
	//
	//	for i, m := range al.MediaId {
	//
	//		if i == 0 {
	//
	//			buffer.WriteString("media_id = ? ")
	//
	//			args = append(args, m)
	//
	//		} else {
	//
	//			buffer.WriteString("or media_id = ?")
	//
	//			args = append(args, m)
	//
	//		}
	//
	//	}
	//
	//}
	//
	//if len(al.Status) != 0 {
	//
	//	if buffer.Len() != 0 {
	//		buffer.WriteString(" and ")
	//	}
	//
	//	for i, s := range al.Status {
	//
	//		if i == 0 {
	//
	//			buffer.WriteString("status = ? ")
	//
	//			args = append(args, s)
	//
	//		} else {
	//
	//			buffer.WriteString("or status = ?")
	//
	//			args = append(args, s)
	//
	//		}
	//
	//	}
	//
	//}
	//
	//if al.AccountName != "null" && al.AccountName != "" {
	//
	//	if buffer.Len() != 0 {
	//		buffer.WriteString(" and ")
	//	}
	//
	//	buffer.WriteString("user_name like ?")
	//
	//	args = append(args, "%"+al.AccountName+"%")
	//
	//}
	//
	//fmt.Println("请求SQL的条件语句: " + buffer.String())
	//
	//// 按条件获取数据
	//// 1、获取账户的总数量
	//
	//count, err = engine.In("advertiser_id", adargs ...).Count()
	//
	//if err != nil {
	//
	//	fmt.Println(fmt.Sprintf("行数: %d", adargs))
	//
	//	fmt.Println("报错: " + err.Error())
	//
	//	return
	//
	//} else {
	//
	//	fmt.Println(fmt.Sprintf("行数: %d", count))
	//
	//}

	//session := context.Daomgr.GetEngine().Where(buffer.String(), args...)
	//
	//if al.SortOrder != "" {
	//	switch  strings.ToLower(al.SortOrder) {
	//	case "desc":
	//		session.Desc("id")
	//	case "asc":
	//		session.Asc("id")
	//	}
	//
	//}
	//
	//if al.Limit != 0 {
	//	session.Limit(al.Limit, (al.Offset-1)*al.Limit)
	//}
	//
	//err = session.Find(&accounts)
	//
	//if err != nil {
	//
	//	common.ResponseErrorCode(ctx, constant.CONSTANT_CODE_SERVER_ERROR)
	//
	//	context.Log.Warn("系统内部错误: " + err.Error())
	//
	//	return
	//
	//}

	// result := make([]*ResponseAccount, 0)

	//for _, a := range accounts {
	//	//result = append(result, &ResponseAccount{
	//	//	Id:           a.Id,
	//	//	AdvertiserId: a.AdvertiserId,
	//	//	MediaId:      a.MediaId,
	//	//	UserId:       a.UserId,
	//	//	UserName:     a.UserName,
	//	//	Status:       a.Status,
	//	//})
	//
	//	fmt.Println(fmt.Sprintf("%d --> %d --> %d --> --> %d --> %s --> %d", a.Id, a.AdvertiserId, a.MediaId, a.UserId, a.UserName, a.Status))
	//}

}

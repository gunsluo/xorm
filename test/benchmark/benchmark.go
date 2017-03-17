package main

import (
	"time"

	"github.com/gunsluo/xorm"

	_ "github.com/go-sql-driver/mysql"
)

type ApplyObject struct {
	ApplyObjectId   int       `xorm:"autoincr pk" json:"applyObjectId"`         //用户id
	ApplyObjectName string    `xorm:"apply_object_name" json:"applyObjectName"` //用户名字
	Environment     string    `xorm:"environment" json:"environment"`           //环境
	EnvironmentId   int       `xorm:"environment_id" json:"environmentId"`      //环境
	TPS             string    `xorm:"tps" json:"tps"`                           //tps
	PeakValue       string    `xorm:"peak_value" json:"peakValue"`              //峰值
	ObjectType      string    `xorm:"object_type" json:"objectType"`            //申请对象类型
	OwnerId         int       `xorm:"owner_id" json:"ownerId"`                  //拥有者Id
	OwnerName       string    `xorm:"owner_name" json:"ownerName"`              //拥有者姓名
	DepartmentId    int       `xorm:"department_id" json:"departmentId"`        //部门Id
	DepartmentName  string    `xorm:"department_name" json:"departmentName"`    //部门名称
	Describe        string    `xorm:"describe" json:"describe"`                 //描述
	Deleted         int       `xorm:"" json:"deleted"`                          //用户状态（0表示正常，1表示已删除，2表示禁用）
	CreatedTime     time.Time `xorm:"created" json:"createdTime"`               //创建时间
	UpdatedTime     time.Time `xorm:"updated" json:"updatedTime"`               //更新时间
	DeletedTime     time.Time `xorm:"deleted_time" json:"deletedTime"`          //删除时间
}

func (this *ApplyObject) TableName() string {
	return "t_apply_object"
}

type ApplyDetail struct {
	ApplyDetailId int       `xorm:"autoincr pk" json:"applyDetailId"`     //用户id
	ApplyObjectId int       `xorm:"apply_object_id" json:"applyObjectId"` //用户id
	UserId        int       `xorm:"user_id" json:"userId"`                //用户名字
	UserName      string    `xorm:"user_name" json:"userName"`            //用户名字
	AuthType      string    `xorm:"auth_type" json:"authType"`            //权限类型
	Key           string    `xorm:"key" json:"key"`                       //环境
	ApplyFlag     string    `xorm:"apply_flag" json:"applyFlag"`          //tps
	ApplyStatus   string    `xorm:"apply_status" json:"applyStatus"`      //峰值
	Deleted       int       `xorm:"" json:"deleted"`                      //用户状态（0表示正常，1表示已删除，2表示禁用）
	CreatedTime   time.Time `xorm:"created" json:"createdTime"`           //创建时间
	UpdatedTime   time.Time `xorm:"updated" json:"updatedTime"`           //更新时间
	DeletedTime   time.Time `xorm:"deleted_time" json:"deletedTime"`      //删除时间
}

func (this *ApplyDetail) TableName() string {
	return "t_apply_detail"
}

type FlowDetail struct {
	ApplyDetail `xorm:"extends"`
	ApplyObject `xorm:"extends"`
}

func selectEntryList(e *xorm.Engine) (fds []*FlowDetail, err error) {

	paramMap := map[string]interface{}{"UserId": 23, "ApplyStatus": "success"}
	err = e.SqlMapClient("test.selectEntryList", &paramMap).Find(&fds)
	if err != nil {
		return nil, err
	}

	return fds, nil
}

func selectEntryList2(e *xorm.Engine) (fds []*FlowDetail, err error) {

	sql := "select * FROM t_apply_detail tad LEFT JOIN t_apply_object tao on tad.apply_object_id = tao.apply_object_id where user_id=?UserId and apply_status=?ApplyStatus"
	paramMap := map[string]interface{}{"UserId": 23, "ApplyStatus": "success"}
	err = e.Sql(sql, &paramMap).Find(&fds)
	if err != nil {
		return nil, err
	}

	return fds, nil
}

func selectEntryAll(e *xorm.Engine) (fds []*FlowDetail, err error) {

	paramMap := map[string]interface{}{"UserId": 23}
	err = e.SqlMapClient("test.selectEntryAll", &paramMap).Find(&fds)
	if err != nil {
		return nil, err
	}

	return fds, nil
}

func selectEntryAll2(e *xorm.Engine) (fds []*FlowDetail, err error) {

	sql := "select * from t_apply_detail tad LEFT JOIN t_apply_object tao on tad.apply_object_id = tao.apply_object_id where user_id=?UserId"
	paramMap := map[string]interface{}{"UserId": 23}
	err = e.Sql(sql, &paramMap).Find(&fds)
	if err != nil {
		return nil, err
	}

	return fds, nil
}

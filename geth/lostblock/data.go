package lostblock

import (
	"database/sql"
	"github.com/golang/glog"
	"strconv"
)

//接口
type Data interface {
	Get(index int64) int64
	Len() int64
}

//mysql
type MysqlData struct {
	db *sql.DB
}

func NewMysqlData(db *sql.DB) *MysqlData {
	return &MysqlData{db:db}
}

func (mysql *MysqlData) Get(index int64) (int64) {
	var lastNumber string
	stmt, err := mysql.db.Prepare("select number from blocks order by number asc limit ?,1")
	if err != nil {
		glog.Infof("mysql get err: %v \n", err)
	}
	row := stmt.QueryRow(index)
	err = row.Scan(&lastNumber)
	if err != nil {
		glog.Infof("mysqlData number err %v \n", err)
		lastNumber = "0"
	}
	i, err := strconv.ParseInt(lastNumber, 10, 64)
	if err != nil {
		glog.Infof("parse int64 error %v \n", err)
		return int64(0)
	}
	glog.Infof("mysqlData number is %s \n", lastNumber)
	return i
}

func (mysql *MysqlData) Len() int64 {
	var count string
	stmt, err := mysql.db.Prepare("select count(1) from blocks")
	if err != nil {
		glog.Infof("mysqlData get err: %v \n", err)
	}
	row := stmt.QueryRow()
	err = row.Scan(&count)
	if err != nil {
		glog.Infof("mysqlData count err %v \n", err)
		count = "0"
	}
	i, err := strconv.ParseInt(count, 10, 64)
	if err != nil {
		glog.Infof("parse int64 error %v \n", err)
		return int64(0)
	}
	glog.Infof("mysqlData count  is %s \n", count)
	return i
}

//array
type ArrayData struct {
	list []int64
}

func NewArrayData() *ArrayData {
	return &ArrayData{list: []int64{1, 2, 5, 6, 9,
		13, 14, 17, 19, 21}}
}

func (array *ArrayData) Get(index int64) int64 {
	return array.list[index]
}

func (array *ArrayData) Len() int64 {
	return int64(len(array.list))
}

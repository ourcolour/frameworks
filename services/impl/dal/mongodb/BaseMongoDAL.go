package mongodb

import (
	"github.com/ourcolour/frameworks/constants/errs"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"reflect"
	"strings"
)

const (
	MGO_ADDRESSES          = "localhost:27018"
	MGO_REPLICATE_SET_NAME = ""

	MGO_DATABASE = "DoubanBookApi"
	MGO_USERNAME = "docker"
	MGO_PASSWORD = "docker"
)

func getDialInfo() *mgo.DialInfo {
	return getDialInfoEx(MGO_ADDRESSES, MGO_REPLICATE_SET_NAME, MGO_DATABASE, MGO_USERNAME, MGO_PASSWORD)
}

func getDialInfoEx(addresses string, replicaSetName string, database string, username string, password string) *mgo.DialInfo {
	var (
		mgoAddrs []string = strings.Split(strings.Replace(addresses, " ", "", -1), ",")
		dialInfo *mgo.DialInfo
	)

	dialInfo = &mgo.DialInfo{
		Addrs:          mgoAddrs,
		Direct:         len(mgoAddrs) < 2,
		ReplicaSetName: replicaSetName,
		Database:       database,
	}

	if "" != username {
		dialInfo.Username = username
		dialInfo.Password = password
	}

	return dialInfo
}

func connect() (*mgo.Session, error) {
	var (
		dialInfo *mgo.DialInfo = getDialInfo()

		result *mgo.Session
		err    error
	)

	result, err = mgo.DialWithInfo(dialInfo)
	if nil != err {
		log.Fatalf("Failed to connect to mongodb, %s.\n", err.Error())
	} else {
		result.SetMode(mgo.Nearest, true)
	}

	return result, err
}

func FindId(dbName string, colName string, id interface{}) (interface{}, error) {
	var (
		result interface{}
		err    error
	)

	session, err := connect()
	if nil != err {
		log.Printf("%s\n", err.Error())
		return nil, err
	}
	defer session.Close()
	col := session.DB(dbName).C(colName)

	query := col.FindId(id)
	err = query.One(&result)

	// 忽略空记录异常
	if reflect.TypeOf(errs.ERR_NOT_FOUND).Elem() == reflect.TypeOf(err).Elem() {
		err = nil
	}

	return result, err
}

func FindOne(dbName string, colName string, selector bson.M) (interface{}, error) {
	var (
		result interface{}
		err    error
	)

	session, err := connect()
	if nil != err {
		log.Printf("%s\n", err.Error())
		return nil, err
	}
	defer session.Close()
	col := session.DB(dbName).C(colName)

	query := col.Find(selector)
	err = query.One(&result)

	// 忽略空记录异常
	if nil != err && reflect.TypeOf(errs.ERR_NOT_FOUND).Elem() == reflect.TypeOf(err).Elem() {
		err = nil
	}

	return result, err
}

func FindAll(dbName string, colName string, selector bson.M, typ reflect.Type) ([]interface{}, error) {
	var (
		result []interface{} = make([]interface{}, 0)
		err    error
	)

	session, err := connect()
	if nil != err {
		log.Printf("%s\n", err.Error())
		return result, err
	}
	defer session.Close()
	col := session.DB(dbName).C(colName)

	// 查询
	query := col.Find(selector)

	itr := query.Iter()
	pObj := reflect.New(typ).Interface()
	for itr.Next(pObj) {
		result = append(result, pObj)
		pObj = reflect.New(typ).Interface()
	}

	return result, err
}

func FindList(dbName string, colName string, selector bson.M, typ reflect.Type, skip int, limit int) ([]interface{}, int64, error) {
	var (
		result           []interface{} = make([]interface{}, 0)
		totalRecordCount int64         = 0
		err              error
	)

	session, err := connect()
	if nil != err {
		log.Printf("%s\n", err.Error())
		return nil, totalRecordCount, err
	}
	defer session.Close()
	col := session.DB(dbName).C(colName)

	// 查询
	query := col.Find(selector).Skip(skip).Limit(limit)
	// 总记录数
	totalRecordCount32, err := query.Count()
	if nil != err {
		return nil, totalRecordCount, err
	}
	totalRecordCount = int64(totalRecordCount32)

	itr := query.Iter()
	pObj := reflect.New(typ).Interface()
	for itr.Next(pObj) {
		result = append(result, pObj)
		pObj = reflect.New(typ).Interface()
	}

	return nil, totalRecordCount, err
}

func Count(dbName string, colName string, selector bson.M) (int, error) {
	var (
		result int = 0
		err    error
	)

	session, err := connect()
	if nil != err {
		log.Printf("%s\n", err.Error())
		return 0, err
	}
	defer session.Close()
	col := session.DB(dbName).C(colName)

	result, err = col.Find(selector).Count()

	return result, err
}

func MustFindOne(dbName string, colName string, selector bson.M) interface{} {
	val, err := FindOne(colName, selector)
	if nil != err {
		return nil
	} else {
		return val
	}
}

func MustFindAll(dbName string, colName string, selector bson.M, typ reflect.Type) []interface{} {
	result, err := FindAll(colName, selector, typ)
	if nil != err {
		log.Panicln(err)
	}
	return result
}

func MustFindList(dbName string, colName string, selector bson.M, typ reflect.Type, skip int, limit int) ([]interface{}, int64) {
	result, total, err := FindList(colName, selector, typ, skip, limit)
	if nil != err {
		log.Panicln(err)
	}
	return result, total
}

// ---

func Insert(dbName string, colName string, data interface{}) error {
	var (
		err error
	)

	session, err := connect()
	if nil != err {
		log.Printf("%s\n", err.Error())
		return err
	}
	defer session.Close()
	col := session.DB(dbName).C(colName)

	err = col.Insert(data)

	return err
}

func Update(dbName string, colName string, selector bson.M, data interface{}) error {
	var (
		err error
	)

	session, err := connect()
	if nil != err {
		log.Printf("%s\n", err.Error())
		return err
	}
	defer session.Close()
	col := session.DB(dbName).C(colName)

	err = col.Update(selector, data)

	return err
}

func UpdateId(dbName string, colName string, id interface{}, data interface{}) error {
	var (
		err error
	)

	session, err := connect()
	if nil != err {
		log.Printf("%s\n", err.Error())
		return err
	}
	defer session.Close()
	col := session.DB(dbName).C(colName)

	err = col.UpdateId(id, data)

	return err
}

func Upsert(dbName string, colName string, selector bson.M, data interface{}) (*mgo.ChangeInfo, error) {
	var (
		changeInfo *mgo.ChangeInfo
		err        error
	)

	session, err := connect()
	if nil != err {
		log.Printf("%s\n", err.Error())
		return nil, err
	}
	defer session.Close()
	col := session.DB(dbName).C(colName)

	changeInfo, err = col.Upsert(selector, data)

	return changeInfo, err
}

func UpsertId(dbName string, colName string, id interface{}, data interface{}) error {
	var (
		err error
	)

	session, err := connect()
	if nil != err {
		log.Printf("%s\n", err.Error())
		return err
	}
	defer session.Close()
	col := session.DB(dbName).C(colName)

	err = col.UpdateId(id, data)

	return err
}

func Remove(dbName string, colName string, selector bson.M) error {
	var (
		err error
	)

	session, err := connect()
	if nil != err {
		log.Printf("%s\n", err.Error())
		return err
	}
	defer session.Close()
	col := session.DB(dbName).C(colName)

	err = col.Remove(selector)

	return err
}

func RemoveAll(dbName string, colName string, selector bson.M) (*mgo.ChangeInfo, error) {
	var (
		changeInfo *mgo.ChangeInfo
		err        error
	)

	session, err := connect()
	if nil != err {
		log.Printf("%s\n", err.Error())
		return nil, err
	}
	defer session.Close()
	col := session.DB(dbName).C(colName)

	changeInfo, err = col.RemoveAll(selector)

	return changeInfo, err
}

func RemoveId(dbName string, colName string, id interface{}) error {
	var (
		err error
	)

	session, err := connect()
	if nil != err {
		log.Printf("%s\n", err.Error())
		return err
	}
	defer session.Close()
	col := session.DB(dbName).C(colName)

	err = col.RemoveId(id)

	return err
}

func ExistsDabatase(dbName string) (bool, error) {
	var (
		result bool = false
		err    error
	)

	if "" == dbName {
		err = errs.ERR_INVALID_PARAMETERS
		return result, err
	}

	session, err := connect()
	if nil != err {
		log.Printf("%s\n", err.Error())
		return result, err
	}
	defer session.Close()

	dbArray, err := session.DatabaseNames()
	if nil != err {
		return result, err
	}

	for _, curDbName := range dbArray {
		if 0 == strings.Compare(curDbName, dbName) {
			result = true
			break
		}
	}

	return result, err
}

func ExistsCollection(dbName string, colName string) (bool, error) {
	var (
		result bool = false
		err    error
	)

	// 参数
	if "" == colName {
		err = errs.ERR_INVALID_PARAMETERS
		return result, err
	}

	session, err := connect()
	if nil != err {
		log.Printf("%s\n", err.Error())
		return result, err
	}
	defer session.Close()
	db := session.DB(dbName)

	collectionArray, err := db.CollectionNames()
	if nil != err {
		return result, err
	}

	for _, curColName := range collectionArray {
		if 0 == strings.Compare(curColName, colName) {
			result = true
			break
		}
	}

	return result, err
}

func Ping() error {
	var (
		err error
	)

	session, err := connect()
	if nil != err {
		log.Printf("%s\n", err.Error())
		return err
	}
	defer session.Close()

	err = session.Ping()

	return err
}

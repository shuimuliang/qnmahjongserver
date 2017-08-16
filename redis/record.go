package redis

import (
	"encoding/json"
	"strconv"
	"time"

	"qnmahjong/def"
	"qnmahjong/pf"
	"qnmahjong/util"
)

var (
	recordPrefix = "mj_record_"
)

// GetRecords ...
func GetRecords(id int32) []*pf.RecordRoom {
	key := recordPrefix + strconv.Itoa(int(id))
	records := make([]*pf.RecordRoom, 0)
	reply, err := Pool.Get().Do("LRANGE", key, 0, -1)
	if err != nil {
		util.LogError(err, "reply", reply, id, def.ErrGetRecordsFromRedis)
		return nil
	}

	replys, ok := reply.([]interface{})
	if ok {
		for i, reply := range replys {
			record := &pf.RecordRoom{}
			err = json.Unmarshal(reply.([]byte), record)
			if err != nil {
				util.LogError(err, "reply", reply, id, def.ErrUnmarshalRedisRecord)
				continue
			}

			createTime, err := time.ParseInLocation("2006-01-02 15:04:05", record.CreateTime, time.Local)
			if err != nil {
				util.LogError(err, "reply", reply, id, def.ErrParseRecordCreateTime)
				continue
			}

			if createTime.Before(time.Now().AddDate(0, 0, -7)) {
				start := 0
				stop := i - 1
				if stop < 0 {
					start = len(replys)
				}
				_, err = Pool.Get().Do("LTRIM", key, start, stop)
				util.LogError(err, "record", record, id, def.ErrTrimRecordToRedis)
				break
			}
			records = append(records, record)
		}
	}
	return records
}

// PutRecord ...
func PutRecord(id int32, record *pf.RecordRoom) {
	key := recordPrefix + strconv.Itoa(int(id))
	bytes, err := json.Marshal(record)
	if err != nil {
		util.LogError(err, "record", record, id, def.ErrMarshalRedisRecord)
		return
	}

	_, err = Pool.Get().Do("LPUSH", key, bytes)
	if err != nil {
		util.LogError(err, "bytes", bytes, id, def.ErrPutRecordToRedis)
	}
}

package v1

import (
	"estate/db"
	"estate/utils"
	"strconv"
)

type JapanModel struct{}

type JapanEstateProgressReturn struct {
	List       []EstateProgressList `json:"list"`
	Pagenation Pagenation           `json:"pagenation"`
}

type Pagenation struct {
	LastId int `json:"last_id"`
}

type EstateProgressList struct {
	Id             int    `json:"id"`
	EstateCode     string `json:"estate_code"`
	HousingType    int    `json:"housing_type"`
	RegionName     string `json:"region_name"`
	AreaName       string `json:"area_name"`
	MeasureArea    string `json:"measure_area"`
	Price          string `json:"price"`
	CustomerNumber int    `json:"customer_number"`
	Status         int    `json:"status"`
	AddTime        string `json:"add_time"`
}

// 日本中介-房源进展
func (this *JapanModel) Japan_EstateProgress(status, perPage, lastId, userId, userType int) (data *JapanEstateProgressReturn, errMsg string) {
	data = new(JapanEstateProgressReturn)
	data.Pagenation.LastId = -1

	// 条件
	var where string
	if userType == 0 { // 销售
		where = ` AND e.user_id=` + strconv.Itoa(userId)
	}
	if status > 0 {
		where += ` AND e.status=` + strconv.Itoa(status)
	}

	// 分页
	if perPage == 0 {
		perPage = 10
	}
	if lastId > 0 {
		// 获取当前分页的最后一条数据的时间
		sql := `SELECT add_time FROM p_estate WHERE id=?`
		row, err := db.Db.Query(sql, lastId)
		if err != nil {
			return data, "获取房源发布时间失败"
		}
		addTime := string(row[0]["add_time"])
		where += ` AND (e.add_time<"` + addTime + `" OR (e.add_time="` + addTime + `" AND e.id<` + strconv.Itoa(lastId) + `))`
	}

	// 获取房源进展列表
	sql := `SELECT e.id, e.code, e.housing_type, e.measure_area, e.price, e.status, e.add_time, r.type regionType, r.name regionName, r.p_id regionPId
			FROM p_estate e
			LEFT JOIN p_region r ON r.id=e.region_id
			WHERE e.is_del=0` + where + ` ORDER BY e.add_time DESC,id DESC LIMIT 0,?`
	rows, err := db.Db.Query(sql, perPage+1)
	if err != nil {
		return nil, "获取房源进展列表失败"
	}
	if len(rows) == 0 {
		return nil, ""
	}
	for key, value := range rows {
		id, _ := strconv.Atoi(string(value["id"]))
		if key < perPage {
			// 地区名称
			var regionName, areaName string
			regionType, _ := strconv.Atoi(string(value["regionType"]))
			if regionType == 1 { // 市
				regionName = string(value["regionName"])
			} else { // 区
				sql := `SELECT name FROM p_region WHERE id=?`
				row, err := db.Db.Query(sql, utils.Str2int(string(value["regionPId"])))
				if err != nil {
					return data, "获取地区失败"
				}
				regionName, areaName = string(row[0]["name"]), string(value["regionName"])
			}

			// 该房源上推荐或咨询客户数量
			sql := `SELECT COUNT(id) count FROM p_recommend WHERE estate_id=?`
			row, err := db.Db.Query(sql, id)
			if err != nil {
				return data, "获取房源意向客户数量失败"
			}
			customerNumber, _ := strconv.Atoi(string(row[0]["count"]))

			// 封装数据
			data.List = append(data.List, EstateProgressList{
				Id:             id,
				EstateCode:     string(value["code"]),
				HousingType:    utils.Str2int(string(value["housing_type"])),
				RegionName:     regionName,
				AreaName:       areaName,
				MeasureArea:    string(value["measure_area"]),
				Price:          string(value["price"]),
				CustomerNumber: customerNumber,
				Status:         utils.Str2int(string(value["status"])),
				AddTime:        string(value["add_time"]),
			})
			lastId = id
		} else {
			data.Pagenation.LastId = lastId
		}
	}
	return
}

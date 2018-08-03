package v1

import (
	"encoding/json"
	"estate/db"
	"estate/utils"
	"strconv"
	"strings"
)

type TouristsModel struct{}

var publicModel = new(PublicModel)

// 游客-房源列表
func (this *TouristsModel) Tourists_EstateList(estParam *PublicEstateListParamter) (data *PublicEstateListReturn, errMsg string) {
	data = new(PublicEstateListReturn)
	data.Pagenation.LastId = -1

	// 条件
	var where, order string

	// 关键字
	if estParam.Keyword != "" {
		where = ` e.code="` + estParam.Keyword + `" `
	}

	// 筛选
	if estParam.ScreenJson != "" {
		screen := new(ScreenJson)
		if err := json.Unmarshal([]byte(estParam.ScreenJson), screen); err != nil {
			return data, "解析筛选json失败"
		}

		// 地区
		if screen.AreaId > 0 {
			where += ` e.region_id=` + strconv.Itoa(screen.AreaId)
		} else if screen.RegionId > 0 {
			sql := `SELECT id FROM p_region WHERE p_id=?`
			rows, err := db.Db.Query(sql, screen.RegionId)
			if err != nil {
				return data, "获取地区失败"
			}
			areaIdArray := make([]string, 0)
			for _, value := range rows {
				areaIdArray = append(areaIdArray, string(value["id"]))
			}
			where += ` e.region_id IN(` + strings.Join(areaIdArray, ",") + `)`
		}

		// 价格区间
		if screen.PriceMin > 0 {
			where += ` e.price>=` + strconv.Itoa(screen.PriceMin)
		}
		if screen.PriceMax > 0 {
			where += ` e.price<=` + strconv.Itoa(screen.PriceMax)
		}

		// 面积区间
		if screen.AreaMin != "" {
			where += ` e.measure_area>=` + screen.AreaMin
		}
		if screen.AreaMax != "" {
			where += ` e.measure_area<=` + screen.AreaMax
		}

		// 房型
		if screen.RoomNumber > 0 {
			where += ` e.huxing REGEXP '^` + strconv.Itoa(screen.RoomNumber) + `'`
		}

		// 类型
		switch screen.HousingType {
		case "0":
			where += ` e.housing_type NOT IN(1,2,3,4)`
		case "1", "2", "3", "4":
			where += ` e.housing_type IN(1,2,3,4)`
		}

		// 房龄
		switch screen.BuildingTime {
		case 1: // 2年以下
			where += ` ADN date_sub(curdate(), interval 2 year)<e.building_time`
		case 2: // 2-5年
			where += ` ADN date_sub(curdate(), interval 5 year)<=e.building_time AND date_sub(curdate(), interval 2 year)>=e.building_time`
		case 3: // 5-10年
			where += ` ADN date_sub(curdate(), interval 10 year)<=e.building_time AND date_sub(curdate(), interval 5 year)>=e.building_time`
		case 4: // 10年以上
			where += ` ADN date_sub(curdate(), interval 10 year)>e.building_time`
		}

		// 朝向
		if screen.Orientation != "" {
			where += ` AND e.orientation="` + screen.Orientation + `"`
		}

		// 楼层
		switch screen.Floor {
		case 1:
			where += ` AND e.floor<6`
		case 2:
			where += ` AND e.floor<=12 AND e.floor>=6`
		case 3:
			where += ` AND e.floor>12`
		}

		// 建筑结构
		if screen.BuildingStructure > 0 {
			where += ` AND e.building_structure=` + strconv.Itoa(screen.BuildingStructure)
		}
	}

	// 分页
	if estParam.PerPage == 0 {
		estParam.PerPage = 10
	}

	// 排序
	switch estParam.Listorder {
	case 0: // 最新
		if estParam.LastId > 0 {
			// 获取当前分页的最后一条数据的时间
			sql := `SELECT add_time FROM p_estate WHERE id=?`
			row, err := db.Db.Query(sql, estParam.LastId)
			if err != nil {
				return data, "获取房源发布时间失败"
			}
			addTime := string(row[0]["add_time"])
			where += ` AND (e.add_time<"` + addTime + `" OR (e.add_time="` + addTime + `" AND e.id<` + strconv.Itoa(estParam.LastId) + `))`
		}

		order += ` ORDER BY e.add_time DESC, e.id DESC `
	case 1: // 面积正序
		if estParam.LastId > 0 {
			// 获取当前分页的最后一条数据的时间
			sql := `SELECT measure_area FROM p_estate WHERE id=?`
			row, err := db.Db.Query(sql, estParam.LastId)
			if err != nil {
				return data, "获取房源面积失败"
			}
			measureArea := string(row[0]["measure_area"])
			where += ` AND (e.measure_area>"` + measureArea + `" OR (e.measure_area="` + measureArea + `" AND r.id>` + strconv.Itoa(estParam.LastId) + `))`
		}

		order += ` ORDER BY e.measure_area, e.id`
	case 2: // 面积倒序
		if estParam.LastId > 0 {
			// 获取当前分页的最后一条数据的时间
			sql := `SELECT measure_area FROM p_estate WHERE id=?`
			row, err := db.Db.Query(sql, estParam.LastId)
			if err != nil {
				return data, "获取房源面积失败"
			}
			measureArea := string(row[0]["measure_area"])
			where += ` AND (e.measure_area<"` + measureArea + `" OR (e.measure_area="` + measureArea + `" AND r.id<` + strconv.Itoa(estParam.LastId) + `))`
		}

		order += ` ORDER BY e.measure_area DESC, e.id DESC`
	case 3: // 价格正序
		if estParam.LastId > 0 {
			// 获取当前分页的最后一条数据的时间
			sql := `SELECT price FROM p_estate WHERE id=?`
			row, err := db.Db.Query(sql, estParam.LastId)
			if err != nil {
				return data, "获取房源面积失败"
			}
			price := string(row[0]["price"])
			where += ` AND (e.price>"` + price + `" OR (e.price="` + price + `" AND r.id>` + strconv.Itoa(estParam.LastId) + `))`
		}

		order += ` ORDER BY e.price, e.id`
	case 4: // 价格倒序
		if estParam.LastId > 0 {
			// 获取当前分页的最后一条数据的时间
			sql := `SELECT price FROM p_estate WHERE id=?`
			row, err := db.Db.Query(sql, estParam.LastId)
			if err != nil {
				return data, "获取房源面积失败"
			}
			price := string(row[0]["price"])
			where += ` AND (e.price<"` + price + `" OR (e.price="` + price + `" AND r.id<` + strconv.Itoa(estParam.LastId) + `))`
		}

		order += ` ORDER BY e.price DESC, e.id DESC`
	default:
		return data, "排序类型错误"
	}

	// 列表
	sql := `SELECT e.id, e.huxing, e.measure_area, e.price, e.housing_type, r.p_id regionPId, r.name regionName, r.type regionType
			FROM p_estate e
			LEFT JOIN p_region r ON r.id=e.region_id
			LEFT JOIN p_user u ON u.id=e.user_id
			WHERE e.is_del=0 ` + where + order + ` LIMIT 0,?`
	rows, err := db.Db.Query(sql, estParam.PerPage+1)
	if err != nil {
		return data, "获取房源列表失败"
	}
	if len(rows) == 0 {
		return nil, ""
	}
	for key, value := range rows {
		if key < estParam.PerPage {
			// 地区名称
			regionName, areaName, errMsg := publicModel.GetEstateName(&GetRegionAreaNameParameter{
				RegionPId:  utils.Str2int(string(value["regionPId"])),
				RegionType: utils.Str2int(string(value["regionType"])),
				RegionName: string(value["regionName"]),
			})
			if errMsg != "" {
				return data, errMsg
			}

			data.List = append(data.List, PublicEstateList{
				EstateId:    utils.Str2int(string(value["id"])),
				EstateName:  regionName + areaName + GetHousingTypeName(utils.Str2int(string(value["housing_type"]))),
				Huxing:      string(value["huxing"]),
				HuxingAlias: GetHuxingAlias(string(value["huxing"])),
				MeasureArea: string(value["measure_area"]),
				RegionName:  string(value["regionName"]),
				HousingType: utils.Str2int(string(value["housing_type"])),
				Price:       string(value["price"]),
				PriceRmb:    "",
			})

			estParam.LastId = utils.Str2int(string(value["id"]))
		} else {
			data.Pagenation.LastId = estParam.LastId
		}
	}

	return
}

// 游客-房源咨询
func (this *TouristsModel) Tourists_EstateConsulting(estateId, sex int, name, wechat string) (errMsg string) {
	// 开启事务
	transaction := db.Db.NewSession()
	if err := transaction.Begin(); err != nil {
		return "开启事务失败"
	}

	// 更新游客表
	sql := `INSERT INTO p_tourists(name, wechat, sex) VALUES(?,?,?)`
	row, err := transaction.Exec(sql, name, wechat, sex)
	if err != nil {
		transaction.Rollback()
		return "更新游客失败"
	}
	lastId, _ := row.LastInsertId()

	// 更新推荐/咨询表
	sql = `INSERT INTO p_recommend(tourists_id, estate_id) VALUES(?,?)`
	_, err = transaction.Exec(sql, lastId, estateId)
	if err != nil {
		transaction.Rollback()
		return "更新咨询表失败"
	}

	// 提交事务
	if err = transaction.Commit(); err != nil {
		return "提交事务失败"
	}
	return
}

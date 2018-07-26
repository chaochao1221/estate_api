package v1

import (
	"estate/db"
	"estate/utils"
)

type TouristsModel struct{}

// type ScreenJson struct {
// 	region_id          region_id
// 	area_id            area_id
// 	price_min          price_min
// 	price_max          price_max
// 	area_min           area_min
// 	area_max           area_max
// 	room_number        room_number
// 	housing_type       housing_type
// 	building_time      building_time
// 	orientation        orientation
// 	floor              floor
// 	building_structure building_structure
// }

// // 游客-房源列表
// func (this *TouristsModel) Tourists_EstastList(estateCode, screenJson string, listorder, perPage, lastId int) (data string, errMsg string) {

// 	return
// }

type TouristsEstateDetailReturn struct {
	EstateId          int    `json:"estate_id"`
	EstateCode        string `json:"estate_code"`
	RegionName        string `json:"region_name"`
	AreaName          string `json:"area_name"`
	RoomNumber        int    `json:"room_number"`
	ExistLivingRoom   int    `json:"exist_living_room"`
	ExistDiningRoom   int    `json:"exist_dining_room"`
	ExistKitchen      int    `json:"exist_kitchen"`
	MeasureArea       string `json:"measure_area"`
	HousingType       int    `json:"housing_type"`
	Price             string `json:"price"`
	PriceRmb          string `json:"price_rmb"`
	LandRights        int    `json:"land_rights"`
	BuildingTime      string `json:"building_time"`
	Floor             int    `json:"floor"`
	TotalFloor        int    `json:"total_floor"`
	BuildingStructure int    `json:"building_structure"`
	Orientation       int    `json:"orientation"`
	RepairFee         int    `json:"repair_fee"`
	State             int    `json:"state"`
	Rent              int    `json:"rent"`
	ReturnRate        string `json:"return_rate"`
	ManageFee         int    `json:"manage_fee"`
}

// 游客-房源详情
func (this *TouristsModel) Tourists_EstateDetail(estateId int) (data *TouristsEstateDetailReturn, errMsg string) {
	sql := `SELECT e.id, e.code, e.region_id, r.type region_type, e.room_number, e.exist_living_room, e.exist_dining_room,
			e.exist_kitchen, e.measure_area, e.housing_type, e.price, e.land_rights, e.building_time, e.floor, e.total_floor,
			e.building_structure, e.orientation, e.repair_fee, e.state, e.rent, e.return_rate, e.manage_fee
			FROM p_estate e
			LEFT JOIN p_region r ON r.id=e.region_id
			WHERE e.id=? AND e.status=0 AND e.is_del=0 AND e.add_time<?`
	row, err := db.Db.Query(sql, estateId)
	if err != nil {
		return data, "获取房源信息失败"
	}
	if len(row) == 0 {
		return data, "该房源信息不存在"
	}

	return &TouristsEstateDetailReturn{
		EstateId:          utils.Str2int(string(row[0]["id"])),
		EstateCode:        string(row[0]["code"]),
		RegionName:        "",
		AreaName:          "",
		RoomNumber:        utils.Str2int(string(row[0]["room_number"])),
		ExistLivingRoom:   utils.Str2int(string(row[0]["exist_living_room"])),
		ExistDiningRoom:   utils.Str2int(string(row[0]["exist_dining_room"])),
		ExistKitchen:      utils.Str2int(string(row[0]["exist_kitchen"])),
		MeasureArea:       string(row[0]["measure_area"]),
		HousingType:       utils.Str2int(string(row[0]["housing_type"])),
		Price:             string(row[0]["price"]),
		PriceRmb:          "",
		LandRights:        utils.Str2int(string(row[0]["land_rights"])),
		BuildingTime:      string(row[0]["building_time"]),
		Floor:             utils.Str2int(string(row[0]["floor"])),
		TotalFloor:        utils.Str2int(string(row[0]["total_floor"])),
		BuildingStructure: utils.Str2int(string(row[0]["building_structure"])),
		Orientation:       utils.Str2int(string(row[0]["orientation"])),
		RepairFee:         utils.Str2int(string(row[0]["repair_fee"])),
		State:             utils.Str2int(string(row[0]["state"])),
		Rent:              utils.Str2int(string(row[0]["rent"])),
		ReturnRate:        string(row[0]["return_rate"]),
		ManageFee:         utils.Str2int(string(row[0]["manage_fee"])),
	}, ""
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

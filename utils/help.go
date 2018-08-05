package utils

import (
	"archive/zip"
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"math/rand"
	"mime/multipart"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	simplejson "github.com/bitly/go-simplejson"
	"github.com/ltt1987/alidayu"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/text/encoding/simplifiedchinese"
)

var phoneticSymbol = map[string]string{
	"ā": "a1",
	"á": "a2",
	"ǎ": "a3",
	"à": "a4",
	"ē": "e1",
	"é": "e2",
	"ě": "e3",
	"è": "e4",
	"ō": "o1",
	"ó": "o2",
	"ǒ": "o3",
	"ò": "o4",
	"ī": "i1",
	"í": "i2",
	"ǐ": "i3",
	"ì": "i4",
	"ū": "u1",
	"ú": "u2",
	"ǔ": "u3",
	"ù": "u4",
	"ü": "v0",
	"ǘ": "v2",
	"ǚ": "v3",
	"ǜ": "v4",
	"ń": "n2",
	"ň": "n3",
	"": "m2",
}

const MobileRegular = `^(13[0-9]|15[012356789]|17[0135678]|18[0123456789]|14[57]|19[0123456789])[0-9]{8}$`

func Md5(str string) string {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(str))
	s := hex.EncodeToString(md5Ctx.Sum(nil))
	return s
}

func Int2str(t int) string {
	str := strconv.Itoa(t)
	return str
}

func Str2int(str string) int {
	t, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}
	return t
}

func Str2int64(s string) int64 {
	t, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0
	}
	return t
}

func Str2float64(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return f
}

func Float64Tostr(f float64) string {
	s := strconv.FormatFloat(f, 'f', -1, 64)
	return s
}

func Float64Toint(f float64) int {
	return int(f)
}

func Int642str(t int64) string {
	str := strconv.FormatInt(int64(t), 10)
	return str
}

func HashPassword(password string) []byte {
	pass := []byte(password)
	hash, _ := bcrypt.GenerateFromPassword(pass, 8)
	if len(hash) == 60 {
		return hash
	}
	return []byte("*")
}

func CheckPassword(hashedPassword, password string) bool {
	if bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) == nil {
		return true
	} else {
		return false
	}
}

//随机生成数字或字母串
//0 纯数字
//1 小写字母
//2 大写字母
//3 数字、大小写字母
//size 生成随机数长度
func Krand(size int, kind int) []byte {
	ikind, kinds, result := kind, [][]int{[]int{10, 48}, []int{26, 97}, []int{26, 65}}, make([]byte, size)
	is_all := kind > 2 || kind < 0
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < size; i++ {
		if is_all { // random ikind
			ikind = rand.Intn(3)
		}
		scope, base := kinds[ikind][0], kinds[ikind][1]
		result[i] = uint8(base + rand.Intn(scope))
	}
	return result
}

// 阿里大于发送短信
func SendMessage(mobile, sign_name, template, param string) error {
	if isOk, _ := regexp.MatchString(MobileRegular, mobile); !isOk {
		return errors.New("手机号格式错误")
	}
	alidayuInfo, err := LoadConfig("alidayu")
	if err != nil {
		return errors.New("短信接口未配置")
	}
	alidayu.AppKey = string(alidayuInfo["app_key"])
	alidayu.AppSecret = string(alidayuInfo["app_secret"])
	alidayu.UseHTTP = true
	success, resp := alidayu.SendSMS(mobile, sign_name, template, param)
	if success != true {
		js, err := simplejson.NewJson([]byte(resp))
		if err != nil {
			return err
		}
		sub_msg, err := js.Get("error_response").Get("sub_msg").String()
		if err != nil {
			return err
		} else {
			return errors.New(sub_msg)
		}
	}
	return nil
}

type FileBack struct {
	File_Path string
	Url       string
	Fid       string
	Err       error
}

type uploadResp struct {
	Data   data_url
	Msg    string
	Status int
}

type data_url struct {
	File_key  string
	File_path string
	Success   bool
	Url       string
}

//文件上传
func UploadFile_New(filename, token string) (fileBack *FileBack) {
	fileBack = new(FileBack)
	info, err := LoadConfig("file_upload")
	if err != nil {
		fileBack.Err = errors.New("获取配置文件失败")
		return
	}
	if info["file_up2"] == "" {
		fileBack.Err = errors.New("配置文件错误")
		return
	}
	url := info["file_up2"] + "single_upload"

	fileBack = new(FileBack)
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	fileWriter, err := bodyWriter.CreateFormFile("uploadfile", filename)
	if err != nil {
		fileBack.Err = err
		return
	}
	// dstFileName := createFilePath(filename)
	//打开文件句柄操作
	fh, err := os.Open(filename)
	if err != nil {
		fileBack.Err = err
		return
	}
	defer fh.Close()

	//iocopy
	_, err = io.Copy(fileWriter, fh)
	if err != nil {
		fileBack.Err = err
		return
	}

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	client := &http.Client{}
	req, err := http.NewRequest("POST", url, bodyBuf)

	req.Header.Set("Content-Type", contentType)
	req.Header.Set("token", token)
	resp, err := client.Do(req)

	if err != nil {
		fileBack.Err = err
		return
	}

	defer resp.Body.Close()
	upload := new(uploadResp)
	if err = decodeJson(resp.Body, upload); err != nil {
		fileBack.Err = err
		return
	}

	if upload.Status != 200 || !upload.Data.Success {
		fileBack.Err = errors.New("上传失败")
		return
	}

	fileBack.File_Path = upload.Data.File_path
	fileBack.Url = upload.Data.Url
	return
}

func decodeJson(r io.Reader, v interface{}) error {
	return json.NewDecoder(r).Decode(v)
}

const (
	KC_RAND_KIND_NUM   = 0 // 纯数字
	KC_RAND_KIND_LOWER = 1 // 小写字母
	KC_RAND_KIND_UPPER = 2 // 大写字母
	KC_RAND_KIND_ALL   = 3 // 数字、大小写字母
)

//生成文件名称
func CreateFilePath() string {
	return strings.Replace(time.Now().Format("2006-01-02-150405.9999"), ".", "", -1) + "-" + strings.ToLower(string(Krand(4, KC_RAND_KIND_ALL))) + ".json"
}

//根据模板获取分表名
//tableName为表模板，分表id使用#站位,例如:table_#_v3
//id 为key值,用值获取分表id
//num 为分表数,如：id为11,num为3,则返回table_2_v3
func CalTableName(tableName string, id, num int) string {
	var data string
	h_int := id % CalcAbs(num)
	str := strings.Replace(tableName, "#", strconv.Itoa(h_int), 1)
	data = str
	return data
}
func CalcAbs(a int) (ret int) {
	ret = (a ^ a>>31) - a>>31
	return
}

//验证手机号码
func CheckMobile(m string) bool {
	if isOk, _ := regexp.MatchString(MobileRegular, m); !isOk {
		return false
	}
	return true
}

// utf8转gbk
func Utf8ToGBK(text string) (string, error) {
	dst := make([]byte, len(text)*2)
	tr := simplifiedchinese.GBK.NewEncoder()
	nDst, _, err := tr.Transform(dst, []byte(text), true)
	if err != nil {
		return text, err
	}
	return string(dst[:nDst]), nil
}

//压缩文件
//files 文件数组，可以是不同dir下的文件或者文件夹
//dest 压缩文件存放地址
func Compress(files []*os.File, dest string) error {
	d, _ := os.Create(dest)
	defer d.Close()
	w := zip.NewWriter(d)
	defer w.Close()
	for _, file := range files {
		fileName, err := Utf8ToGBK(file.Name())
		if err != nil {
			return err
		}
		fw, _ := w.Create(fileName)
		filecontent, err := ioutil.ReadFile(file.Name())
		if err != nil {
			return nil
		}
		_, err = fw.Write(filecontent)
		if err != nil {
			fmt.Println(err)
		}
	}
	return nil
}

//阿拉伯数组转中文数子
func Chinese_int(i int) string {
	var data string
	str := strconv.Itoa(i)
	str = strings.TrimLeft(str, "0")

	a := []string{"零", "一", "二", "三", "四", "五", "六", "七", "八", "九"}
	b := []string{"", "十", "百", "千", "万"}

	if 10 <= i && i < 20 {
		num := i % 10
		if num == 0 {
			data = b[1]
		} else {
			data += b[1] + a[i%10]
		}
		return data
	}

	num := len(str)

	n := 1
	for i := 0; i < num; i++ {
		b_num := num - i - 1
		a_num, _ := strconv.Atoi(string(str[i]))
		if n != 0 || a_num != 0 {
			if a_num == 0 {
				if i != num-1 {
					data += a[a_num]
				}
			} else {
				data += a[a_num] + b[b_num]
			}
		}
		n = a_num
	}
	data = strings.Trim(data, "零")
	return data
}

//入库前替换特殊字符
func StringFilter(str string) string {
	str = strings.Replace(str, `\`, `\\`, -1)
	str = strings.Replace(str, `'`, `\'`, -1)
	str = strings.Replace(str, `"`, `\"`, -1)
	return str
}

//四舍五入
//num为精确到小数点后位数，正整数
func StringRounding(str string, num int) string {
	var ff, t float64
	var data string
	if str == "" {
		return data
	}
	ff, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return "0"
	}
	if num < 0 {
		num = -num
	}
	if ff == 0 {
		return "0"
	}

	f := math.Pow10(num)
	x := ff * f
	if x >= 0.0 {
		t = math.Ceil(x)
		if (t - x) > 0.5 {
			t -= 1.0
		}
	} else {
		t = math.Ceil(-x)
		if (t + x) > 0.5 {
			t -= 1.0
		}
		t = -t
	}
	x = t / f
	data = strconv.FormatFloat(x, 'f', num, 64)
	if num > 0 {
		data = strings.TrimRight(data, "0")
		data = strings.TrimRight(data, ".")
	}
	return data
}

func FloatRounding(fl float64, num int) string {
	var t float64
	var data string

	if num < 0 {
		num = -num
	}
	if fl == 0 {
		return "0"
	}
	f := math.Pow10(num)
	fl += 0.0000000001
	x := fl * f
	if x >= 0.0 {
		t = math.Ceil(x)
		if (t - x) > 0.5 {
			t -= 1.0
		}
	} else {
		t = math.Ceil(-x)
		if (t + x) > 0.5 {
			t -= 1.0
		}
		t = -t
	}
	x = t / f
	data = strconv.FormatFloat(x, 'f', num, 64)
	if num > 0 {
		data = strings.TrimRight(data, "0")
		data = strings.TrimRight(data, ".")
	}
	return data
}

// 获取简拼
func GetPinyin(str string) string {
	py := []string{}
	pys := [][]string{}
	fpy := []string{}
	initial := ""
	for _, val := range str {
		regular := `[a-z|A-Z]`
		if ok, _ := regexp.MatchString(regular, string(val)); ok {
			initial += strings.ToLower(string(val))
		} else {
			hans := []rune(string(val))
			for _, r := range hans {
				value, ok := PinyinDict[int(r)]
				if ok {
					py = strings.Split(value, ",")
				} else {
					py = []string{}
				}
				if len(py) > 0 {
					py = py[:1]
					pys = append(pys, py)
				}
			}
		}
	}
	for _, v := range pys {
		fpy = append(fpy, v[0])
	}
	for _, v := range fpy {
		rn := []rune(v)
		f := string(rn[0])
		temp := phoneticSymbol[f]
		if temp != "" {
			f = temp[:1]
		}
		initial += f
	}
	if initial == "" {
		initial = "#"
	}
	return initial
}

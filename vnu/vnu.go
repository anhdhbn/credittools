package vnu


import (
	"log"
	"net/http"
	"net/url"
	"io/ioutil"
	"regexp"
	"strings"
	"fmt"
	"strconv"
	"net/http/cookiejar"
	"golang.org/x/net/publicsuffix"
	"encoding/json"
	"time"
)

type User struct {
    Mssv string
	Pass string
	TypeLogin string
	Mhp string
}

type ResponseVnu struct {
	Success bool
	Message string
}


func DachSachMonHocDaDangKy(client *http.Client, user_ User) (string) {
	path := fmt.Sprintf("/danh-sach-mon-hoc-da-dang-ky/%s",  user_.TypeLogin)
	req := createRequest(path, "POST", true, "")
	resp, html := executeRequest(client, req)
	if (resp == nil) {
		return ""
	}
	return html
}

func XacNhanDangKy(client *http.Client, user_ User) (bool) {
	path := fmt.Sprintf("/xac-nhan-dang-ky/%s", user_.TypeLogin)
	req := createRequest(path, "POST", true, "")
	resp, html := executeRequest(client, req)
	if (resp == nil) {
		return false
	}
	obj := parseJson(string(html))
	log.Println(user_.Mssv, user_.Mhp, obj)
	return obj.Success
}



func DangKyMonHoc(client *http.Client, user_ User) (bool) {
	path := fmt.Sprintf("/chon-mon-hoc/%s/%s/%s", user_.Mhp, user_.TypeLogin, "2")
	req := createRequest(path, "POST", true, "")
	resp, html := executeRequest(client, req)
	if (resp == nil) {
		return false
	}
	obj := parseJson(string(html))
	return obj.Success
}

func GetDanhSachMonHoc(client *http.Client, user_ User) (string) {
	path := fmt.Sprintf("/danh-sach-mon-hoc/%s/%s", user_.TypeLogin, "2")
	req := createRequest(path, "POST", true, "")
	resp, html := executeRequest(client, req)
	if (resp == nil) {
		return ""
	}
	return html
}


func Login(client *http.Client, user_ User, check bool) (bool){
	if (check) {
		req := createRequest("/", "GET", false , "")
		resp, html := executeRequest(client, req)
		if (resp != nil && CheckLogin(html)) {
			return true
		}
	}


	req := createRequest("/dang-nhap", "GET", false , "")
	resp, html := executeRequest(client, req)
	if  (resp == nil) {
		return false
	}

	
	token := getToken(html)
	if (token == "")  {
		return false
	}
	data := url.Values{}
    data.Set("__RequestVerificationToken", token)
    data.Set("LoginName", user_.Mssv)
	data.Set("Password", user_.Pass)

	req = createRequest("/dang-nhap", "POST", true , data.Encode())
	resp, html = executeRequest(client, req)
	return CheckLogin(html)
}

func parseJson(data string) (ResponseVnu) {
	var resp ResponseVnu
	if (strings.Contains(data, `success`)){
		json.Unmarshal([]byte(data), &resp)
		return resp
	} else {
		return ResponseVnu{Success: false, Message: data}
	}
}

func executeRequest(client *http.Client, request *http.Request) (*http.Response, string){
	response, err := client.Do(request)
    if err != nil {
		return nil, ""
    } else {
		data, _ := ioutil.ReadAll(response.Body)
		defer response.Body.Close()
        return response, string(data)
    }
}

func createRequest(path string, method string, postType bool, params string) (*http.Request){
	link := fmt.Sprintf("http://dangkyhoc.vnu.edu.vn%s", path)
	request, _ := http.NewRequest(method, link, strings.NewReader(params))
	
	if (postType) {
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		request.Header.Add("Content-Length", strconv.Itoa(len(params)))
	}
	return request
}

func InitHTTP(timeoutSecond int) (*http.Client){
	options := cookiejar.Options{
        PublicSuffixList: publicsuffix.List,
    }
    jar, err := cookiejar.New(&options)
    if err != nil {
        log.Fatal(err)
    }
    client := http.Client{Jar: jar, Timeout: time.Duration(timeoutSecond) * time.Second}
	return &client
}

func getToken(html string) (string){
	r := regexp.MustCompile(`__RequestVerificationToken.*ue="(.*?)"`)
	res := r.FindStringSubmatch(html)
	if (len(res) != 2) {
		return ""
	}
	return res[1]
}

func CheckLogin(html string) (bool) {
	return strings.Contains(html, "/Account/Logout")
}

func CheckIsInLoginScreen(html string) (bool) {
	return strings.Contains(html, `$("#LoginName").focus();`)
}
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
)

type User struct {
    Mssv string
	Pass string
	Type_login string
	Mhp string
}

type ResponseVnu struct {
	Success bool
	Message string
}


func Dach_sach_mon_hoc_da_dang_ky(client *http.Client, user_ User) (string) {
	path := fmt.Sprintf("/danh-sach-mon-hoc-da-dang-ky/%s",  user_.Type_login)
	req := create_request(path, "POST", true, "")
	resp, html := execute_request(client, req)
	if (resp == nil) {
		return ""
	}
	return html
}

func Xac_nhan_dang_ky(client *http.Client, user_ User) (bool) {
	path := fmt.Sprintf("/xac-nhan-dang-ky/%s", user_.Type_login)
	req := create_request(path, "POST", true, "")
	resp, html := execute_request(client, req)
	if (resp == nil) {
		return false
	}
	obj := parse_json(string(html))
	log.Println(user_.Mssv, user_.Mhp, obj)
	return obj.Success
}



func Dang_ky_mon_hoc(client *http.Client, user_ User) (bool) {
	path := fmt.Sprintf("/chon-mon-hoc/%s/%s/%s", user_.Mhp, user_.Type_login, "2")
	req := create_request(path, "POST", true, "")
	resp, html := execute_request(client, req)
	if (resp == nil) {
		return false
	}
	obj := parse_json(string(html))
	return obj.Success
}

func Get_danh_sach_mon_hoc(client *http.Client, user_ User) (string) {
	path := fmt.Sprintf("/danh-sach-mon-hoc/%s/%s", user_.Type_login, "2")
	req := create_request(path, "POST", true, "")
	resp, html := execute_request(client, req)
	if (resp == nil) {
		return ""
	}
	return html
}


func Login(client *http.Client, user_ User, check bool) (bool){
	if (check) {
		req := create_request("/", "GET", false , "")
		resp, html := execute_request(client, req)
		if (resp != nil && Check_login(html)) {
			return true
		}
	}


	req := create_request("/dang-nhap", "GET", false , "")
	resp, html := execute_request(client, req)
	if  (resp == nil) {
		return false
	}

	
	token := get_token(html)
	if (token == "")  {
		return false
	}
	data := url.Values{}
    data.Set("__RequestVerificationToken", token)
    data.Set("LoginName", user_.Mssv)
	data.Set("Password", user_.Pass)

	req = create_request("/dang-nhap", "POST", true , data.Encode())
	resp, html = execute_request(client, req)
	return Check_login(html)
}

func parse_json(data string) (ResponseVnu) {
	var resp ResponseVnu
	if (strings.Contains(data, `success`)){
		json.Unmarshal([]byte(data), &resp)
		return resp
	} else {
		return ResponseVnu{Success: false, Message: data}
	}
}

func execute_request(client *http.Client, request *http.Request) (*http.Response, string){
	response, err := client.Do(request)
    if err != nil {
		return nil, ""
    } else {
		data, _ := ioutil.ReadAll(response.Body)
		defer response.Body.Close()
        return response, string(data)
    }
}

func create_request(path string, method string, post_type bool, params string) (*http.Request){
	link := fmt.Sprintf("http://dangkyhoc.vnu.edu.vn%s", path)
	request, _ := http.NewRequest(method, link, strings.NewReader(params))
	
	if (post_type) {
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		request.Header.Add("Content-Length", strconv.Itoa(len(params)))
	}
	return request
}

func Init_http() (*http.Client){
	options := cookiejar.Options{
        PublicSuffixList: publicsuffix.List,
    }
    jar, err := cookiejar.New(&options)
    if err != nil {
        log.Fatal(err)
    }
    client := http.Client{Jar: jar}
	return &client
}

func get_token(html string) (string){
	r := regexp.MustCompile(`__RequestVerificationToken.*ue="(.*?)"`)
	res := r.FindStringSubmatch(html)
	if (len(res) != 2) {
		return ""
	}
	return res[1]
}

func Check_login(html string) (bool) {
	return strings.Contains(html, "/Account/Logout")
}

func Check_is_in_login_screen(html string) (bool) {
	return strings.Contains(html, `$("#LoginName").focus();`)
}
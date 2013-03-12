package openstack

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/MG-RAST/KBNO/kbno-server/client"
	"github.com/MG-RAST/KBNO/kbno-server/conf"
	"io/ioutil"
	"strconv"
	"time"
)

var (
	Token    string
	TokenExp time.Time
)

func authToken() (err error) {
	if time.Now().After(TokenExp) {
		username := conf.Openstack["username"].Str
		password := conf.Openstack["password"].Str
		tentant := conf.Openstack["tenant_name"].Str
		url := conf.Openstack["auth_url"].Str + "/tokens"
		header := &map[string]string{"Content-Type": "application/json"}
		data := bytes.NewBufferString("{\"auth\":{\"passwordCredentials\":{\"username\": \"" + username + "\", \"password\":\"" + password + "\"}, \"tenantName\":\"" + tentant + "\"}}")
		if res, err := client.Post(url, header, data); err != nil {
			return err
		} else {
			r := tokenRes{}
			body, _ := ioutil.ReadAll(res.Body)
			json.Unmarshal(body, &r)
			Token = r.Access.Token.ID
			TokenExp, _ = time.Parse(time.RFC3339, r.Access.Token.Expires)
		}
	}
	return nil
}

func Boot(username, name string) (s *statusRes, err error) {
	authToken()
	url := conf.Openstack["nova_api"].Str + "/" + conf.Openstack["tenant_id"].Str + "/servers"
	flavor := conf.Openstack["flavor"].Str
	image := conf.Openstack["image"].Str
	userData, _ := BootScript(username)
	data := bytes.NewBufferString("{\"server\":{\"flavorRef\":\"" + flavor + "\", \"imageRef\":\"" + image + "\", \"name\":\"" + name + "\", \"user_data\":\"" + userData + "\"}}")
	header := &map[string]string{
		"Content-Type":   "application/json",
		"Content-Length": strconv.Itoa(data.Len()),
		"X-Auth-Token":   Token,
	}
	res, err := client.Post(url, header, data)
	if err != nil {
		return nil, err
	}
	r := bootRes{}
	body, _ := ioutil.ReadAll(res.Body)
	err = json.Unmarshal(body, &r)
	if err != nil {
		return nil, err
	}
	return Status(r.Server.ID)
}

func Status(id string) (s *statusRes, err error) {
	authToken()
	header := &map[string]string{
		"X-Auth-Token": Token,
	}
	url := conf.Openstack["nova_api"].Str + "/" + conf.Openstack["tenant_id"].Str + "/servers/" + id
	res, err := client.Get(url, header, nil)
	if err != nil {
		return nil, err
	}
	body, _ := ioutil.ReadAll(res.Body)
	err = json.Unmarshal(body, &s)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func StatusAll() (s *statusAllRes, err error) {
	authToken()
	image := conf.Openstack["image"].Str
	header := &map[string]string{
		"X-Auth-Token": Token,
	}
	url := conf.Openstack["nova_api"].Str + "/" + conf.Openstack["tenant_id"].Str + "/servers/detail?image=" + image
	res, err := client.Get(url, header, nil)
	if err != nil {
		return nil, err
	}
	body, _ := ioutil.ReadAll(res.Body)
	err = json.Unmarshal(body, &s)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func Delete(id string) (success bool, err error) {
	authToken()
	url := conf.Openstack["nova_api"].Str + "/" + conf.Openstack["tenant_id"].Str + "/servers/" + id
	header := &map[string]string{
		"X-Auth-Token": Token,
	}
	res, err := client.Delete(url, header, nil)
	if err != nil {
		return false, err
	}
	// according to openstack api 204 is success
	if res.StatusCode == 204 {
		return true, nil
	}
	return false, errors.New(res.Status)
}

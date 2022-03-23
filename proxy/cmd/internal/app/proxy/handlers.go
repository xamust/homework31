package proxy

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

func checkConnect(host string, port string) (result bool) {

	timeout := time.Second
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), timeout)
	if err != nil {
		fmt.Println("Connecting error:", err)
		result = false
	}
	if conn != nil {
		defer conn.Close()
		fmt.Println("Opened", net.JoinHostPort(host, port))
		result = false
	}
	return
}

func (p *AppProxy) Balance() (localInstance string) {
	//balance....
	if COUNT == 0 {

		localInstance = p.config.FirstInst

		COUNT++
	} else {

		localInstance = p.config.SecondInst

		COUNT--
	}
	p.logger.Info(fmt.Sprintf("Instance = %s", localInstance))
	return
}

// Create CreateUser
func (p *AppProxy) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost && r.Header.Get("Content-Type") == "application/json" {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			p.logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}

		client := &http.Client{}
		req, err := http.NewRequest(http.MethodPost, p.Balance()+r.RequestURI, bytes.NewBuffer(body))
		if err != nil {
			p.logger.Error(err.Error())
			return
		}
		req.Header.Add("Content-Type", "application/json")
		// Fetch Request
		resp, err := client.Do(req)
		if err != nil {
			p.logger.Error(err.Error())
			return
		}
		defer resp.Body.Close()

		// Read Response Body
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			p.logger.Error(err.Error())
			return
		}

		p.logger.Info(string(respBody))
		w.WriteHeader(resp.StatusCode)
		w.Write(respBody)
		return
	}
}

// MakeFriends MakeFriends
func (p *AppProxy) MakeFriends(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost && r.Header.Get("Content-Type") == "application/json" {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			p.logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
		client := &http.Client{}
		req, err := http.NewRequest(http.MethodPost, p.Balance()+r.RequestURI, bytes.NewBuffer(body))
		if err != nil {
			p.logger.Error(err.Error())
			return
		}
		req.Header.Add("Content-Type", "application/json")
		// Fetch Request
		resp, err := client.Do(req)
		if err != nil {
			p.logger.Error(err.Error())
			return
		}
		defer resp.Body.Close()

		// Read Response Body
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			p.logger.Error(err.Error())
			return
		}

		p.logger.Info(string(respBody))
		w.WriteHeader(resp.StatusCode)
		w.Write(respBody)
		return

	}
}

//Delete DeleteUser
func (p *AppProxy) Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodDelete && r.Header.Get("Content-Type") == "application/json" {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			p.logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
		client := &http.Client{}
		req, err := http.NewRequest(http.MethodDelete, p.Balance()+r.RequestURI, bytes.NewBuffer(body))
		if err != nil {
			p.logger.Error(err.Error())
			return
		}
		req.Header.Add("Content-Type", "application/json")
		// Fetch Request
		resp, err := client.Do(req)
		if err != nil {
			p.logger.Error(err.Error())
			return
		}
		defer resp.Body.Close()

		// Read Response Body
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			p.logger.Error(err.Error())
			return
		}

		p.logger.Info(string(respBody))
		w.WriteHeader(resp.StatusCode)
		w.Write(respBody)
		return
	}
}

// GetFriends GetUserFriends
func (p *AppProxy) GetFriends(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet && r.Header.Get("Content-Type") == "application/json" {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			p.logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
		client := &http.Client{}
		req, err := http.NewRequest(http.MethodGet, p.Balance()+r.RequestURI, bytes.NewBuffer(body))
		if err != nil {
			p.logger.Error(err.Error())
			return
		}
		req.Header.Add("Content-Type", "application/json")
		// Fetch Request
		resp, err := client.Do(req)
		if err != nil {
			p.logger.Error(err.Error())
			return
		}
		defer resp.Body.Close()

		// Read Response Body
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			p.logger.Error(err.Error())
			return
		}

		p.logger.Info(string(respBody))
		w.WriteHeader(resp.StatusCode)
		w.Write(respBody)
		return
	}
}

// Put UpdateUserAge
func (p *AppProxy) Put(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPut && r.Header.Get("Content-Type") == "application/json" {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			p.logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
		client := &http.Client{}
		req, err := http.NewRequest(http.MethodPut, p.Balance()+r.RequestURI, bytes.NewBuffer(body))
		if err != nil {
			p.logger.Error(err.Error())
			return
		}
		req.Header.Add("Content-Type", "application/json")
		// Fetch Request
		resp, err := client.Do(req)
		if err != nil {
			p.logger.Error(err.Error())
			return
		}
		defer resp.Body.Close()

		// Read Response Body
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			p.logger.Error(err.Error())
			return
		}

		p.logger.Info(string(respBody))
		w.WriteHeader(resp.StatusCode)
		w.Write(respBody)
		return
	}
}

//mu handlers (for debug)

func (p *AppProxy) GetAll(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet && r.Header.Get("Content-Type") == "application/json" {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			p.logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
		client := &http.Client{}
		req, err := http.NewRequest(http.MethodGet, p.Balance()+r.RequestURI, bytes.NewBuffer(body))
		if err != nil {
			p.logger.Error(err.Error())
			return
		}
		req.Header.Add("Content-Type", "application/json")
		// Fetch Request
		resp, err := client.Do(req)
		if err != nil {
			p.logger.Error(err.Error())
			return
		}
		defer resp.Body.Close()

		// Read Response Body
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			p.logger.Error(err.Error())
			return
		}

		p.logger.Info(string(respBody))
		w.WriteHeader(resp.StatusCode)
		w.Write(respBody)
		return
	}
}

func (p *AppProxy) GetUserInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet && r.Header.Get("Content-Type") == "application/json" {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			p.logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
		client := &http.Client{}
		req, err := http.NewRequest(http.MethodGet, p.Balance()+r.RequestURI, bytes.NewBuffer(body))
		if err != nil {
			p.logger.Error(err.Error())
			return
		}
		req.Header.Add("Content-Type", "application/json")
		// Fetch Request
		resp, err := client.Do(req)
		if err != nil {
			p.logger.Error(err.Error())
			return
		}
		defer resp.Body.Close()

		// Read Response Body
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			p.logger.Error(err.Error())
			return
		}

		p.logger.Info(string(respBody))
		w.WriteHeader(resp.StatusCode)
		w.Write(respBody)
		return
	}
}

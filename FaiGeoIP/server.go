/*
###########################################################################################
#                                                                                         #
#    Copyright 2018 Khang H. Nguyen (kevinhg86)                                           #
#    E-mail: kevin@fai.host | Web: http://kevinhng86.iblog.website                        #
#    Contributors: https://github.com/kevinhng86/noka-encryption/blob/master/CONTRIBUTORS #                                                    #
#                                                                                         #    
#                                                                                         #
#    Permission is hereby granted, free of charge, to any person obtaining a copy         #
#    of this software and associated documentation files (the "Software"), to deal        #
#    in the Software without restriction, including without limitation the rights         #
#    to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies     #
#    of the Software, and to permit persons to whom the Software is furnished             # 
#    to do so, subject to the following conditions:                                       #
#                                                                                         #
#    The above copyright notice and this permission notice shall be  included in all      #
#    copies or substantial portions of the Software.                                      #
#                                                                                         #
#    THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR           #
#    IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,             #
#    FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL              #
#    THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER           #
#    LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,        #
#    OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN            #
#    THE SOFTWARE.                                                                        #
#                                                                                         #
#	 CREDIT:                                                                              #
#    A special thanks to Chirag Mehta                                                     #
#    http://chir.ag/projects/geoiploc for the half division formula.                      #
#                                                                                         #
###########################################################################################
*/

package FaiGeoIP

import(
	"net/http"
	"encoding/json"
	"time"
	"errors"
)

type HttpResponse struct{
	Status int `json:"status"`
	Err	string `json:"error"`
	Data map[string]string `json:"data"`
}

func HttpHandler(w http.ResponseWriter, r *http.Request){
	ip := []string{ "" }
	locale := []string{ "" }
	
	if r.Method == "GET" {
		ip, ok = r.URL.Query()["ip"]
		if ok != true { ip = []string{ "" } }
		
		locale, ok = r.URL.Query()["locale"]
		if ok != true { locale = []string{ "en" } }
		
		if len(ip[0]) > 0{
			HttpGeoIP(w, ip[0], locale[0])
		} else {
			HttpOther(w)
		}		
	} else {
		HttpOther(w)
	}	
}

func HttpGeoIP(w http.ResponseWriter, ip string, locale string){
	obj := new(HttpResponse)
	obj.Data, err = MaxmindGeoIP(ip, locale)
	if err != nil{
		obj.Status = 1
		obj.Err = err.Error()
	} else {
		obj.Status = 0
		obj.Err = ""
	}

	re, _ := json.Marshal(obj)
	
	w.Header().Set("Content-Type", "application/json")
	w.Write(re)
}

func HttpOther(w http.ResponseWriter){
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(HttpMessage))
}


func HttpStart() error {
	server := http.Server{
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		IdleTimeout:  45 * time.Second,
		Addr: ":" + HttpPort,
		Handler: nil,
	}

	http.HandleFunc("/", HttpHandler)
	
	Logger("Http server listen on port " + HttpPort + " unless error is reported.")
	err = server.ListenAndServe()
	if err != nil { return errors.New("Couldn't start the web server because of " + err.Error()) }

	return nil
}

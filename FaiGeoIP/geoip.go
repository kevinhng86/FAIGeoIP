/*
###########################################################################################
#                                                                                         #
#    Copyright 2018 Khang H. Nguyen (kevinhg86)                                           #
#    E-mail: kevin@fai.host | Web: http://kevinhng86.iblog.website                        #    
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
	"fmt"
	"regexp"
	"errors"
	"strconv"
)

func MaxmindGeoIP(ip string, locale string) (map[string]string, error ){
	data := make(map[string]string)
	if maxmind.ip4dblen == 0 { return data, errors.New("Maxmind IP Database is not loaded.") }
	if maxmind.locdblen == 0 { return data, errors.New("Maxmind Location Database is not loaded.") }
	
	_, ok = maxmind.locdbmap[locale];
	if ok != true { return data, errors.New("Locale " + locale + " is not load. Loaded locale are " + maxmind.locdbname) }
	
	ipd, err := IP4ToDec(ip)
	if err != nil { return data, err }

	upto := maxmind.ip4dblen
	from := 0
	ref := -1
	
	for from < upto {
		idn := from + int((upto - from) / 2)
		if( ipd >= maxmind.ip4db[idn].From && ipd <= maxmind.ip4db[idn].To ){
			ref = idn
			break
		}
		if (ipd < maxmind.ip4db[idn].From){
			if upto == idn { break }
			upto = idn
		}
		if (ipd > maxmind.ip4db[idn].To){
			if from == idn { break }
			from = idn
		}
	}
	
	// if not found return, not found error
	if ref < 0 {
		return data, errors.New("Information for IP " + ip + " is not found.")
	}

	// Check Location Data
	locinfo, ok := maxmind.locdb[maxmind.locdbmap[locale]].data[maxmind.ip4db[ref].Geoname_id]
	if ok != true { return data, errors.New("Couldn't obtain location info for data") }

	// If everything is ok let work on the data.
	// Compose IP data
	data["input_ip"] = ip
	data["network"] = maxmind.ip4db[ref].Network
	data["is_anonymous_proxy"] = maxmind.ip4db[ref].Is_anonymous_proxy
	data["is_satellite_provider"] = maxmind.ip4db[ref].Is_satellite_provider
	data["latitude"] = maxmind.ip4db[ref].Latitude
	data["longitude"] = maxmind.ip4db[ref].Longitude
	data["accuracy_radius"] = maxmind.ip4db[ref].Accuracy_radius

	// Compose location data
	data["locale_code"] = locinfo.Locale_code
	data["continent_code"] = locinfo.Continent_code
	data["continent_name"] = locinfo.Continent_name
	data["country_iso_code"] = locinfo.Country_iso_code
	data["country_name"] = locinfo.Country_name
	data["subdivision_1_iso_code"] = locinfo.Subdivision_1_iso_code
	data["subdivision_1_name"] = locinfo.Subdivision_1_name
	data["subdivision_2_iso_code"] = locinfo.Subdivision_2_iso_code
	data["subdivision_2_name"] = locinfo.Subdivision_2_name
	data["city_name"] = locinfo.City_name
	data["metro_code"] = locinfo.Metro_code
	data["time_zone"] = locinfo.Time_zone
	data["is_in_european_union"] = locinfo.Is_in_european_union

	return data, nil
}

func IP4ToBinStr (ip string) (string, error){
	regip, _ := regexp.MatchString("^[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}(\\/[0-9]{1,2})?$", ip)
	if regip == false { return "", errors.New("IPv4 " + ip + " is not a valid IP.")  }	
	var s [4]int
	m := []string{"00000000","00000001","00000010","00000011","00000100","00000101","00000110","00000111","00001000",
				  "00001001","00001010","00001011","00001100","00001101","00001110","00001111","00010000","00010001",
				  "00010010","00010011","00010100","00010101","00010110","00010111","00011000","00011001","00011010",
				  "00011011","00011100","00011101","00011110","00011111","00100000","00100001","00100010","00100011",
				  "00100100","00100101","00100110","00100111","00101000","00101001","00101010","00101011","00101100",
				  "00101101","00101110","00101111","00110000","00110001","00110010","00110011","00110100","00110101",
				  "00110110","00110111","00111000","00111001","00111010","00111011","00111100","00111101","00111110",
				  "00111111","01000000","01000001","01000010","01000011","01000100","01000101","01000110","01000111",
				  "01001000","01001001","01001010","01001011","01001100","01001101","01001110","01001111","01010000",
				  "01010001","01010010","01010011","01010100","01010101","01010110","01010111","01011000","01011001",
				  "01011010","01011011","01011100","01011101","01011110","01011111","01100000","01100001","01100010",
				  "01100011","01100100","01100101","01100110","01100111","01101000","01101001","01101010","01101011",
				  "01101100","01101101","01101110","01101111","01110000","01110001","01110010","01110011","01110100",
				  "01110101","01110110","01110111","01111000","01111001","01111010","01111011","01111100","01111101",
				  "01111110","01111111","10000000","10000001","10000010","10000011","10000100","10000101","10000110",
				  "10000111","10001000","10001001","10001010","10001011","10001100","10001101","10001110","10001111",
				  "10010000","10010001","10010010","10010011","10010100","10010101","10010110","10010111","10011000",
				  "10011001","10011010","10011011","10011100","10011101","10011110","10011111","10100000","10100001",
				  "10100010","10100011","10100100","10100101","10100110","10100111","10101000","10101001","10101010",
				  "10101011","10101100","10101101","10101110","10101111","10110000","10110001","10110010","10110011",
				  "10110100","10110101","10110110","10110111","10111000","10111001","10111010","10111011","10111100",
				  "10111101","10111110","10111111","11000000","11000001","11000010","11000011","11000100","11000101",
				  "11000110","11000111","11001000","11001001","11001010","11001011","11001100","11001101","11001110",
				  "11001111","11010000","11010001","11010010","11010011","11010100","11010101","11010110","11010111",
				  "11011000","11011001","11011010","11011011","11011100","11011101","11011110","11011111","11100000",
				  "11100001","11100010","11100011","11100100","11100101","11100110","11100111","11101000","11101001",
				  "11101010","11101011","11101100","11101101","11101110","11101111","11110000","11110001","11110010",
				  "11110011","11110100","11110101","11110110","11110111","11111000","11111001","11111010","11111011",
				  "11111100","11111101","11111110","11111111"}
				  
	fmt.Sscanf(ip, "%d.%d.%d.%d", &s[0], &s[1], &s[2], &s[3])
	
	if s[0] > 255 || s[1] > 255 || s[2] > 255 || s[3] > 255 {
		return "", errors.New("IPv4 " + ip + " is not a valid IP.") 
	}
	
	return m[s[0]] + m[s[1]] + m[s[2]] + m[s[3]], nil
}

func IP4ToDec(ip string) (int, error){
	ip, err = IP4ToBinStr(ip)
	if err != nil { return -1, err }
	o, _ := strconv.ParseInt(ip, 2, 32)

	return int(o), nil
}


func GeoIPLoadDBs() error {
	err = MaxmindLoadDBLOC()
	if err != nil { return err }
	
	err = MaxmindLoadDBIP4()
	if err != nil { return err }
	
	return nil
}

func UpdateDBs(p bool) error{
	err = MaxmindUpdate(p)
	if err != nil { return err }

	if p == true { Logger("Successfully update databases.") }	
	return nil
}


func IP4FromTo(ip string) (int, int, error) {
	var s [4]int
	var b string
	var n int
	t0 := "00000000000000000000000000000000"
	t1 := "11111111111111111111111111111111"
	
	regip, _ := regexp.MatchString("^[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}(\\/[0-9]{1,2})?$", ip)
	regnet, _ := regexp.MatchString("^[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}\\/[0-9]{1,2}$", ip)
	
	if regip == false {
		return -1, -1, errors.New(ip + " is not a valid IP address.")
	}
	
	if regnet == false {
		from, err := IP4ToDec(ip)
		if err != nil { return -1, -1, err }
		return from, from, nil
	}
	
	fmt.Sscanf(ip, "%d.%d.%d.%d/%d", &s[0], &s[1], &s[2], &s[3], &n)
	b = fmt.Sprintf("%08b", s[0]) + fmt.Sprintf("%08b", s[1]) + fmt.Sprintf("%08b", s[2]) + fmt.Sprintf("%08b", s[3])

	d := 32 - n
	l := 32 // 32 - 1 (IPv4)

	// From IP

	from, _ := strconv.ParseInt(b[ :(l - d)] + t0[0:d], 2, 32)

	// TO IP
	to, _ := strconv.ParseInt(b[ :(l - d)] + t1[0:d], 2, 32)

	return int(from), int(to), nil
}

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
	"encoding/csv"
	"bufio"
	"os"
	"io"
	"errors"
	"strconv"
	"fmt"
)

type maxmind_obj_ip4 struct{
	Network string `json:"network"`
	From int `json:"from"`
	To int `json:"to"`
	Geoname_id string `json:"geoname_id"`
	Registered_country_geoname_id string `json:"registered_country_geoname_id"`
	Represented_country_geoname_id string `json:"represented_country_geoname_id"`
	Is_anonymous_proxy string `json:"is_anonymous_proxy"`
	Is_satellite_provider string `json:"is_satellite_provider"`
	Postal_code string `json:"postal_code"`
	Latitude string `json:"latitude"`
	Longitude string `json:"longitude"`
	Accuracy_radius string `json:"accuracy_radius"`
}

type maxmind_obj_loc struct{
	Locale_code string
	Continent_code string
	Continent_name string
	Country_iso_code string
	Country_name string
	Subdivision_1_iso_code string
	Subdivision_1_name string
	Subdivision_2_iso_code string
	Subdivision_2_name string
	City_name string
	Metro_code string
	Time_zone string
	Is_in_european_union string
}

type maxmind_obj_loc_db struct{
	lang string
	len int
	data map[string]maxmind_obj_loc
}

type maxmind_db struct{
	ip4db []maxmind_obj_ip4
	ip4dblen int
	locdb []maxmind_obj_loc_db
	locdblen int
	locdbmap map[string]int
	locdbname string
}

func MaxmindUpdate(p bool) error {
	err = MaxmindImport( maxmindpath + maxmindipfile, maxmindpath + maxmindipdbfile )
	if err != nil { return err }

	if p == true {
		fmt.Println("Successfully updated/imported Maxmind IP database.")
	}
	return nil
}


func MaxmindImport (infile string, outfile string ) error {
	s := ","
	
    f, err := os.Open(infile)
    if err != nil { return errors.New("Can't open file import file " + infile + " with error: " + err.Error() )}
    defer f.Close()	

    r := csv.NewReader(bufio.NewReader(f))
	
	w, err := os.OpenFile(outfile , os.O_RDWR|os.O_CREATE, 0666)
    if err != nil { return errors.New("Can't open file for write " + outfile + " with error: " + err.Error() )}
	defer w.Close()

	w.Truncate(0)
	w.Seek(0,0)
		
	// Skip first line
    _, _ = r.Read()
	
	for {
		record, err := r.Read()
        if err == io.EOF {
            break
        }

		if len(record) != 10 {
			return errors.New("Import file " + infile + "is not the correct maxmind csv data syntax that is supported by this program.")
		}

        f, t, _ := IP4FromTo(record[0]) 
        
        line := record[0] + s + strconv.Itoa(f) + s + strconv.Itoa(t) + s + record[1] + s + record[2] + s + record[3] + s + record[4] + s + record[5] + s + record[6] + s + record[7] + s + record[8] + s + record[9] + "\n"
        
        w.Write([]byte(line))        
    }
    
    return nil
}

/* 
 * Max Mind IP File Syntax
 * 
 * network
 * from (our) *convert to int64 when load
 * to (our) *convert to int64 when load
 * geoname_id
 * registered_country_geoname_id 
 * represented_country_geoname_id
 * is_anonymous_proxy
 * is_satellite_provider
 * postal_code latitude
 * longitude
 * accuracy_radius
 * 
 */

func MaxmindLoadDBIP4() error {
	file := maxmindpath + maxmindipdbfile
	
	f, err := os.Open(file)
    if err != nil { return errors.New("Can't open file " + file + " for loading with error: " + err.Error() )}
    defer f.Close()	

    r := csv.NewReader(bufio.NewReader(f))
	i := 0
	
	for {

		record, err := r.Read()
        if err == io.EOF {
            break
        }
		i++

		if len(record) != 12 {
			return errors.New("Loading file " + file + "is not the correct data format that is supported by this program.")
		}
		
		f, err := strconv.ParseInt(record[1], 10, 64)
		if err != nil { return errors.New("Something is wrong with the csv data. can't parse to integer 64 for field from with error: " + err.Error() )}
		t, err := strconv.ParseInt(record[2], 10, 64)
		if err != nil { return errors.New("Something is wrong with the csv data. can't parse to integer 64 for field to with error: " + err.Error() )}
		
		rs := new(maxmind_obj_ip4)
		rs.Network = record[0]
		rs.From = int(f)
		rs.To = int(t)
		rs.Geoname_id = record[3]
		rs.Registered_country_geoname_id = record[4]
		rs.Represented_country_geoname_id = record[5]
		rs.Is_anonymous_proxy = record[6]
		rs.Is_satellite_provider = record[7]
		rs.Postal_code = record[8]
		rs.Latitude = record[9]
		rs.Longitude = record[10]
		rs.Accuracy_radius = record[11]
		
		maxmind.ip4db = append(maxmind.ip4db, *rs)
    }
    
    maxmind.ip4dblen = i
	return nil
}

/*
 * Maxmind Info Database Schema
 * 
 * geoname_id
 * locale_code
 * continent_code
 * continent_name
 * country_iso_code
 * country_name
 * subdivision_1_iso_code
 * subdivision_1_name
 * subdivision_2_iso_code
 * subdivision_2_name
 * city_name
 * metro_code
 * time_zone
 * is_in_european_union
 *
 */

func MaxmindLoadDBLOC() error {

	for k, v := range maxmindlocfiles {
		file := maxmindpath + v 

		f, err := os.Open(file)
		if err != nil { return errors.New("Can't open location file for loading " + file + " with error: " + err.Error() )}
		defer f.Close()
		
		r := csv.NewReader(bufio.NewReader(f))
		ra := new(maxmind_obj_loc_db)
		ra.data = make(map[string]maxmind_obj_loc)
		i := 0
		 
		for {
			record, err := r.Read()
			if err == io.EOF {
				break
			}
			i++

			if len(record) != 14 {
				return errors.New("Loading file " + file + "is not the correct data format that is supported by this program.")
			}
			
			rs := new(maxmind_obj_loc)
			rs.Locale_code = record[1]
			rs.Continent_code = record[2]
			rs.Continent_name = record[3]
			rs.Country_iso_code = record[4]
			rs.Country_name = record[5]
			rs.Subdivision_1_iso_code = record[6]
			rs.Subdivision_1_name = record[7]
			rs.Subdivision_2_iso_code = record[8]
			rs.Subdivision_2_name = record[9]
			rs.City_name = record[10]
			rs.Metro_code = record[11]
			rs.Time_zone = record[12]
			rs.Is_in_european_union = record[13]
			
			ra.data[record[0]] = *rs
		}					
		ra.lang = k
		ra.len = i
		maxmind.locdb = append(maxmind.locdb, *ra)
	}

	maxmind.locdbmap = make(map[string]int)	
	maxmind.locdblen = len(maxmind.locdb)
	maxmind.locdbname = ""

	for i := 0; i < maxmind.locdblen; i++ {
		maxmind.locdbmap[maxmind.locdb[i].lang] = i
		maxmind.locdbname += maxmind.locdb[i].lang
		if i + 1 < maxmind.locdblen { maxmind.locdbname += "," }
	}
	
	return nil
} 

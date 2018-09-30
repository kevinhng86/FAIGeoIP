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
#    CREDIT:                                                                              #
#    A special thanks to Chirag Mehta                                                     #
#    http://chir.ag/projects/geoiploc for the half division formula.                      #
#                                                                                         #
###########################################################################################
*/
package FaiGeoIP

import(
    "errors"
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

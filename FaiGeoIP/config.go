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

func init(){	
    //Http server config
    HttpPort = "17824"

    // Maxmind Configuration
    maxmindpath = "./maxmind/"
    maxmindipfile = "GeoLite2-City-Blocks-IPv4.csv"
    maxmindipdbfile = "faigeoip-maxmind-ip.csv"
    maxmindlocfiles = make(map[string]string)
    maxmindlocfiles["de"] = "GeoLite2-City-Locations-de.csv"
    maxmindlocfiles["en"] = "GeoLite2-City-Locations-en.csv"
    maxmindlocfiles["es"] = "GeoLite2-City-Locations-es.csv"
    maxmindlocfiles["fr"] = "GeoLite2-City-Locations-fr.csv"
    maxmindlocfiles["ja"] = "GeoLite2-City-Locations-ja.csv"
    maxmindlocfiles["pt-BR"] = "GeoLite2-City-Locations-pt-BR.csv"
    maxmindlocfiles["ru"] = "GeoLite2-City-Locations-ru.csv"
    maxmindlocfiles["zn-CN"] = "GeoLite2-City-Locations-zh-CN.csv"
}

func ConfigInit(){
    // This can only be initialize after everything is loaded.
    HttpMessage = "<html><head><title>FAIGeoIP Server</title></head>"
    HttpMessage += "<body>"
    HttpMessage += "<center><h1>Welcome To FAI GeoIP Server</h1></center>"
    HttpMessage += "<p>To use this server please send a get request with ip and locale. For example "
    HttpMessage += "example.com/?ip=1.1.1.1&locale=en. Locale can be omitted. If locale is ommitted then locale en will use. "
    HttpMessage += "The following locale is supported: \"" + maxmind.locdbname + "\"."
    HttpMessage += "If a GET request is reveiced with the ip parameter is set the server will response in JSON format. The ip parameter must be greater than 0 length.</p>"
    HttpMessage += "This product includes GeoLite2 data created by MaxMind, available from <a href=\"http://www.maxmind.com\">http://www.maxmind.com</a>."
    HttpMessage += "</body></html>"
}

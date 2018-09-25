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
package main

import(
	"os"
	my "./FaiGeoIP"	
)

var(
	err error
	ops string
)

func init(){
	ops =  "start/stop/update"
}

func main(){
	Args := os.Args
	Pid := os.Getpid()
	PidFile := "./FaiGeoIP.pid"
	
	if len(Args) > 1 {
		if Args[1] == "start" {
			err := my.Start(Pid, PidFile);
			if err != nil {
				my.Logger(err)
			}
			
		} else if Args[1] == "stop" {
			my.Kill(PidFile);
		} else if Args[1] == "update"{
			err = my.UpdateDBs(true)
			if err != nil { my.Logger(err) }
		} else if Args[1] == "help" {
			my.Logger(ops)
		} else {
			my.Logger("Unsupported operation.")
		}
	} else {
		my.Logger("No input. Supported operation are " + ops)
	}
}

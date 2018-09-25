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
#                                                                                        #
###########################################################################################
*/
package FaiGeoIP

import(
	"os"
	"errors"
	"strconv"
	"syscall"
)

func Kill(pidFile string) (bool, error) {
	 if _, err := os.Stat(pidFile); os.IsNotExist(err) {
		// PID File is not exist
		return true, errors.New("PID file is not existed. The program is probably not running.")
    } else {
		// PID File Exist.
		f, err := os.Open(pidFile)
		if err != nil { return false, errors.New("PID file " + pidFile + " exist, but couldn't be open with error: " + err.Error()) }

		fi, err := f.Stat()		
		if err != nil { return false, errors.New("Couldn't run file stat command on PID file " + pidFile + " with error: " + err.Error() ) }

		d := make( []byte, fi.Size() )

		_, err = f.Read(d)
		if err != nil { return false, errors.New("Pid file " + pidFile +" exist but couldn't read with error: " + err.Error() ) }
		
		pid, err := strconv.Atoi(string(d))
        if err != nil { return false, errors.New("PID file is not the right format. Please delete PID file "+ pidFile +" manually.") }
        
        process, err := os.FindProcess( pid )
        
        if err != nil {
			// Windows process not found
			// If PID file is existed but process is not found then remove pid files.			
			RemovePIDFile(pidFile)
			return true, errors.New("Process was not found from PID file, thus PID file is removed.")
        } else {
            err = process.Signal(syscall.Signal(0))
			if err == nil {
				// Linux process found
				err2 = process.Signal(os.Interrupt)
				if err2 != nil { return false, errors.New("Process couldn't be kill. PID file was not removed.") }
				
				RemovePIDFile(pidFile)
				return true, nil
			} else {
				if err.Error() == "operation not permitted" {
					// Linux process is found but not owned
					return false, errors.New("Process is found, but is not owned. Thus couldn't be kill. PID file was not removed.")
				}
				
				// Linux everything else process is not found
				RemovePIDFile(pidFile)
				return true, errors.New("Process was not found from PID file, thus PID file is removed.")
			}
			
			// This part for windows where process is found.
			
			err2 = process.Signal(os.Interrupt)
			if err2 != nil { return false, errors.New("Process couldn't be kill. PID file was not removed.") }
	
			RemovePIDFile(pidFile)
	
			return true, nil
        }        
    }   
}

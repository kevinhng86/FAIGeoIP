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
#     CREDIT:                                                                              #
#    A special thanks to Chirag Mehta                                                     #
#    http://chir.ag/projects/geoiploc for the half division formula.                      #
#                                                                                         #
###########################################################################################
*/
package FaiGeoIP

import(
    "os"
    "strconv"
    "errors"
    "syscall"
)

func Start(pid int, pidFile string) error {
    
    if ok, err = CheckRun(pidFile); ok == true {
        if err != nil { return errors.New("Couldn't start process. PID file existed and there was an issue with loading the PID files with the following subprocess error: " + err.Error() ) }
        return errors.New("Already running.")
    } else {
        err2 = CreatePIDFile(pid, pidFile)
        if err2 != nil { return errors.New("Process did not start, can't create PID files with error: " + err2.Error()) }
    }
    
    err = GeoIPLoadDBs()
    if err != nil {
        err2 = RemovePIDFile(pidFile)
        if err != nil { return errors.New(err.Error() + "\n" +  err2.Error()) }
        
        return err
    }
    
    Logger("Loaded Maxmind databases.")
    Logger("Memoery Usage: ")
    PrintMemUsage()
    Logger("Starting the web server.")

    ConfigInit();

    err = HttpStart()
    if err != nil {
        err2 = RemovePIDFile(pidFile)
        if err != nil { return errors.New(err.Error() + "\n" + err2.Error()) } 

        return err 
    }

    return nil    
}

func CheckRun(pidFile string) (bool,error) {

    if _, err := os.Stat(pidFile); os.IsNotExist(err) {
        // PID File is not exist
        return false, nil
    } else {
        // PID File Exist.
        f, err := os.Open(pidFile)
        if err != nil { return true, errors.New("PID file " + pidFile + " exist, but couldn't be open with error: " + err.Error()) }

        fi, err := f.Stat()        
        if err != nil { return true, errors.New("Couldn't run file stat command on PID file " + pidFile + " with error: " + err.Error() ) }

        d := make( []byte, fi.Size() )

        _, err = f.Read(d)
        if err != nil { return true, errors.New("Pid file " + pidFile +" exist but couldn't read with error: " + err.Error() ) }
        
        pid, err := strconv.Atoi(string(d))
        if err != nil { return true, errors.New("PID file is not the right format. Please delete PID file "+ pidFile +" manually.") }
        
        process, err := os.FindProcess( pid )
        
        if err != nil {
            // Windows process not found
            
            // If PID file is existed but process is not found
            // then the program probably crashed. Thus, assumably it is safe to 
            // delete the PID file. Caller will create a new PID file.
            
            RemovePIDFile(pidFile)
            return false, nil
        } else {
            err := process.Signal(syscall.Signal(0))
            if err == nil {
                // Linux process found
                return true, nil
            } else {
                if err.Error() == "operation not permitted" {
                    // Linux process is found but not owned
                    return true, nil
                }
                
                // Linux everything else process is not found
                RemovePIDFile(pidFile)
                return false, nil
            }
            
            // This part for windows where process is found.
            return true, nil
        }        
    }    
}

func CreatePIDFile(pid int, pidFile string) error {
    d := []byte(strconv.Itoa(pid))
    f, err := os.Create(pidFile)
    if err != nil { return errors.New("Couldn't create PID file with error: " + err.Error() ) }
    
    _, err = f.Write(d)
    if err != nil { return errors.New("Couldn't write to PID file with error: " + err.Error() ) }
    
    f.Close()

    return nil
}

func RemovePIDFile(pidFile string) error {
    err = os.Remove(pidFile)
    if err != nil { return errors.New("Couldn't remove PID file. Please delete PID file " + pidFile + " manually") }
    
    return nil
}


// write to stdout/stderr and create a rotating logfile
// This is the test stub. 

package main


import (
	"mlog-master"
	"fmt"
	"log"
	"errors"
)

func main() {
	logfile := "howie_app.log"				// this is what we write to. 
	logsize := 5*1024*1024					// Log size in Megabytes = 5
	logrots := 3							// rotate log 3 times

	//	Write everything to console;
	//	Write Info, Warning and Error messages to file
	//	Wet log file to 5mb and only keep 3 log rotations. 
	mlog.StartEx(mlog.LevelTrace, logfile, logsize, logrots)

	fmt.Println("logging tester started")

	// display and log an example Info message
	mlog.Info("Logfile: %s with size %d and numrots of %d", logfile, logsize, logrots)

	// Show how normal log works. 
	size := 1234
	log.Println("Standard stdout logging (not to file, but with timestamp) ... the size: ", size)

	// display and log an example Info message
	mlog.Info("Plumbing interface 1.2.3.4 to UP")

	// Warning
	temp := 45.73
	maxtemp := 45.00
	cpunum := 4
	mlog.Warning("CPU %d (%f) has exceeded the max threshold temperature of %f", cpunum, temp, maxtemp)

	// Error
	// error requires an "error" type
	err := errors.New("This is an error to display")
	mlog.Error(err)

	// Trace - Add tracing to the output (not generally needed)
	mlog.Trace("This is a trace message - it is fun fun fun -  but wont be logged to the file?")

	// Fatal
	mlog.Fatal("Testing fatal error message.  This should hault the program!!")

	// Since we did a "fatal" above, this should not be seen
	fmt.Println("This is the end of the program coming from standard i/o - since fatal is used above, this should not be seen")
}
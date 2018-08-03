# higlog
A GO Logging Wrapper - writes to file and console

log is simple logging module for go, Provide a standard interface for output to the screen and a logfile. 
It includes a rotating file feature and console logging.
The message in the log file will be prepended with the severity and also the standard LOG time/date stamp and of course, the message
The message displayed on the screen, however, will not include the time/datestamp - only the severity and the message itself.  

Options to this package - (added 07/27/18)
You can call this with or without a log file - if a log file is specified: 
By default, the log will be rolled over to a backup file when its size reaches 10Mb and 10 such files will be created (and eventually reused).
Alternatively, you can specify the max size of the log file before it gets rotated, and the number of backup files you want to create, with the StartEx function.

## Some example calls to this package: 
```
You can decide what to log when calling mlog.Start
	- LevelTrace logs:   everything
	- LevelInfo  logs:   info, Warnings and Errors
	- LevelWarn  logs:   Warning and Errors
	- LevelError logs:   just Errors

 Example for startx - Write everything to console and don't write to a file.
 mlog.Start(mlog.LevelInfo, "")

 Example for Start - Write everything to console and write Info, Warning and Error messages to file. 
			since it is not specified, the default file size will be max of 10Mb and rotate to 10 backup log files
		mlog.Start(mlog.LevelTrace, "app.log")

 Example StartEx - Start with Extra information: 
					Write everything to console;
					write Info, Warning and Error messages to file
					set log file to 5mb and only keep 4 copies of old logs. 
 		mlog.StartEx(mlog.LevelInfo, "app.log", 5*1024*1024, 4)

 Sample calls: 

 Info: 	
	 // write to screen and logfile. 
	 mlog.Info("This is a piece of information")
 
 Warning: 
	// write to screen and logfile. 
	foo := "foo"
    mlog.Warning("BlaBla - %s", foo)

 Trace: 
	// write to screen and logfile. Note, this has newlines around it. 
    mlog.Trace"This is a trace")

 Error (use standard Go errors structures
	err := errors.New("This is an error to display")
	mlog.Error(err))

 Fatal - display the message and STOP execution!
   	mlog.Fatal("Testing fatal error message.  This should hault the program!!")
```

TimeIt for Windows, Linux, and Mac
==================================

TimeIt is a tiny cross-platform portable command line utility for measuring program execution time. It works the same way across Windows, Linux, and Mac, and any other platform or operating system supported by the Go language. The name "TimeIt" is inspired by [Python's timeit command] (https://docs.python.org/2/library/timeit.html), but that timeit does something different. "time" is not a good name because "time" is the name of existing commands in both Windows and Linux.

TimeMe is a tiny utility that is used to test TimeIt during development. It initially started out as a throwaway script, but I think it ends up being useful enough to be used to test any application that requires introducing a variable delay (sleep time) and gives an expected exit code.

TimeIt is very similar to the Unix/Linux [`time`] (https://en.wikipedia.org/wiki/Time_%28Unix%29) command. It executes the given command including its arguments, and then prints out the amount of time taken at the end. Like the Unix `time` command, TimeIt prints out the real time (the actual amount of time elapsed), the user time (the amount of time spent executing the non-system code), and the system time (the amount of time spent executing system code). TimeIt is meant to be as simplistic as possible. TimeIt does not take any input options to configure its behavior. TimeIt will just print out the necessary information in a single line, and then exit with the same exit code by the executed command.

Download
--------
Please find the latest versions of TimeIt and TimeMe in this web page:

- [https://github.com/choksheak/timeit/tree/master/distribution] (https://github.com/choksheak/timeit/tree/master/distribution)

I currently only support Windows and Ubuntu Linux. If you need TimeIt for any other platform and/or architecture, it is very easy to build it yourself! Follow these steps to build TimeIt:

1. Install Go on your computer. ([Download Go] (https://golang.org/dl/))
2. Download [`timeit.go`] (https://github.com/choksheak/timeit/blob/master/timeit.go).
3. Run this command to produce the timeit executable: `go build timeit.go`
4. Add `timeit` to your `PATH` environment variable.

Why was TimeIt developed?
-------------------------

I was looking for a way to measure command execution times in the Windows Command Prompt or PowerShell Prompt. However, I could not find a good script or command that I can use to do that. It is a very simple problem in desperate need of a simple solution. Therefore, I decided that I had to do something about it, to write up a simple script that could do just that. This script, now TimeIt, would be both functional and minimally intrusive. It tries to do what is needed, and do it as quietly as possible without sacrificing leaving out useful information.

Usage and Examples
------------------

When you use TimeIt, just pretend that you are running whatever program you intended, but add `timeit` right to the front of the command line.

```
USAGE: timeit <command> [arg]...
```

For example, the Hello World sample run of TimeIt looks like this:

```
C:\> timeit echo Hello World
Hello World

[timeit] Real 55.454ms, Usr 0ns, Sys 15.625ms (20:52:30-20:52:30)
```

TimeMe, which is included in the distribution zips, can be used to test the TimeIt command. The first argument of TimeMe is the exit code. The second argument is the number of milliseconds to sleep. A sample use of TimeMe to test TimeIt, with an exit code of 0, and a sleep of 1000 milliseconds, looks like this:

```
C:\> timeit timeme 0 1000
[timeme] This string goes to stdout.
[timeme] This string goes to stderr.
[timeme] Sleeping for 1000 milliseconds ...
[timeme] Exiting with code 0

[timeit] Real 1.122s, Usr 15.625ms, Sys 0ns (20:57:25-20:57:26)
```

Limitations
-----------
TimeIt currently does not support non-commands. This includes but might not be limited to the following:
- PowerShell aliases
- Linux/Mac shell aliases

I have no plans to support these for now but we can revisit them in the future. In any case, aliases can be rewritten in the long form as the original commands which can be executed by TimeIt.

Contact Info
------------
In case you wanted to contact me regarding anything (no spam please), feel free to send me an email at [findfile.go@gmail.com] (mailto:findfile.go@gmail.com). That's right, the email here says "findfile", not "timeit".

Just in case you might find this useful, feel free to also check-out [FindFile] (https://github.com/choksheak/findfile) here in GitHub as well which is also written by me. [FindFile] (https://github.com/choksheak/findfile) is a quick and nifty command line utility, also written in Go, that helps you to search for files without getting into your way.

Have fun timing your commands!

The [FindFile] (https://github.com/choksheak/findfile) Team

- email: [findfile.go@gmail.com] (mailto:findfile.go@gmail.com)
- website: [https://github.com/choksheak/findfile] (https://github.com/choksheak/findfile)


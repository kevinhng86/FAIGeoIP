# FAIGeoIP
A GeoIP server that written in Go. 

This server does not use MySQL or any database engine.
This server require a lot of memory to run. It was tested to consume at least 5 GB of memory.
The trade off for FAIGeoIP memory consumption is that it doesn't consume SQL query nor disk read or drive because everything is hold in memory.
FAIGeoIP server was not design with maximum security nor as an eco friendly software. It was design with only one thing in mind and that is speed.

The FAIGeoIP server is a beta version but it is okay for production usage.

Because speed was the only factor. I did not implement the program to listen to a os signal. I also didn't want the program to spawn any unnecessary process.

FAIGeoIP server was tested to have the capability of handling 200,000 requests per minute or 1 million requests per 5 minutes.

How to use: (If you didn't reconfigure the program)

1. Build the program.
2. Download the Maxmind City database from(www.maxmind.com) in csv format. Extract the files into a folder call "maxmind". Leave the built executeable file in the parent folder of the folder maxmind.
3. First run, start the executable file with "update". It will take sometimes to import the maxmind database.
4. After that start the executable file with "start".
5. Warning: when use stop on the executeable, it will kill whatever process ID that is recorded in the pid file.

Notes:

Most of the configuration can be change in the config.go file. The maxmind ip and location structure can be change. However, FAIGeoIP was built to work only with the current maxmind city database in csv format. Any other format might produce unexpected result or error.

Future versions of FAiGeoIP might support more than one type of database but because of the shear volume of how memory is consumed. That may not be possible.

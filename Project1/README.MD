Required to use this Program is a TFTP Client

Run tftp_server.Go

Then run a client tftp call from a powershell 

from CMD in a windows enviroment

cd to ./files/ directory in the src folder
Then run tftp commands

TFTP [-i] host [GET | PUT] source [destination]

run 
tftp localhost put inputfile.txt outputfile.txt 
in CMD prompt

Example successful outPut:
C:\Users\Cla45>tftp localhost get Test2.txt
Transfer successful: 13851 bytes in 1 second(s), 13851 bytes/s


inputfile : dir of a file the client wants to send this is based off of the curdir of the cmd prompt.

outputFile : simply give a Name that the server will either overwrite or create a new file in the ./files/ directory of the running tftp server.

run 
tftp localhost get outputfile.txt outputfileClone.txt
in CMD prompt

will return file from ./files/ if it exists

When a TFTP call is Done the Server shutsdown. 
So the server must be recalled inbetween commands


Steps to run:
1)'go run tftp_server.go' in a PowerShell from SRC folder
2)'tftp localhost get ReallyLarge.txt' from a CMD Prompt from SRC/inputs/
3)Repeat until done

TFTP [-i] host [GET | PUT] source [destination]

READ:

Example Run and expected OutPut:
Client:
C:\Users\Cla45\Desktop\Comp429\Project1\Inputs>tftp localhost get ReallyLarge.txt
Transfer successful: 41553 bytes in 1 second(s), 41553 bytes/s

Output From Server:
PS C:\Users\Cla45\Desktop\Comp429\Project1> go run tftp_server.go
2021/10/01 19:08:02 Connected to a UDP Client
☺
2021/10/01 19:08:02  RRQ Packet recieved
 Found the Mode and FileName netascii ReallyLarge.txt
2021/10/01 19:08:02 Connected to a UDP Client
☺
2021/10/01 19:08:02  ACK Packet recieved
2021/10/01 19:08:02 [0 1]
2021/10/01 19:08:02 Connected to a UDP Client
☺
2021/10/01 19:08:02  ACK Packet recieved
2021/10/01 19:08:02 [0 2]
2021/10/01 19:08:02 Connected to a UDP Client
☺
2021/10/01 19:08:02  ACK Packet recieved
2021/10/01 19:08:02 [0 3]
2021/10/01 19:08:02 Connected to a UDP Client
☺
.
.
.
.
.
2021/10/01 19:08:02  ACK Packet recieved
2021/10/01 19:08:02 [0 69]
2021/10/01 19:08:02 Connected to a UDP Client
☺
2021/10/01 19:08:02  ACK Packet recieved
2021/10/01 19:08:02 [0 70]
2021/10/01 19:08:02 Connected to a UDP Client
☺
2021/10/01 19:08:02  ACK Packet recieved
2021/10/01 19:08:02 [0 71]
2021/10/01 19:08:02 Connected to a UDP Client
☺
2021/10/01 19:08:02  ACK Packet recieved
2021/10/01 19:08:02 [0 72]
2021/10/01 19:08:02 Connected to a UDP Client
☺
2021/10/01 19:08:02  ACK Packet recieved
2021/10/01 19:08:02 [0 73]
2021/10/01 19:08:02 Connected to a UDP Client
☺
2021/10/01 19:08:02  ACK Packet recieved
2021/10/01 19:08:02 [0 74]
2021/10/01 19:08:02 Connected to a UDP Client
☺
2021/10/01 19:08:02  ACK Packet recieved
2021/10/01 19:08:02 [0 75]
2021/10/01 19:08:02 Connected to a UDP Client
☺
2021/10/01 19:08:02  ACK Packet recieved
2021/10/01 19:08:02 [0 76]
2021/10/01 19:08:02 Connected to a UDP Client
☺
2021/10/01 19:08:02  ACK Packet recieved
2021/10/01 19:08:02 [0 77]
2021/10/01 19:08:02 Connected to a UDP Client
☺
2021/10/01 19:08:02  ACK Packet recieved
2021/10/01 19:08:02 [0 78]
2021/10/01 19:08:02 Connected to a UDP Client
☺
2021/10/01 19:08:02  ACK Packet recieved
2021/10/01 19:08:02 [0 79]
2021/10/01 19:08:02 Connected to a UDP Client
☺
2021/10/01 19:08:02  ACK Packet recieved
2021/10/01 19:08:02 [0 80]
2021/10/01 19:08:02 Connected to a UDP Client
☺
2021/10/01 19:08:02  ACK Packet recieved
2021/10/01 19:08:02 [0 81]
2021/10/01 19:08:02 read udp 127.0.0.1:69: use of closed network connection
exit status 1
PS C:\Users\Cla45\Desktop\Comp429\Project1>


Example 2 WRITE:

Client:
C:\Users\Cla45\Desktop\Comp429\Project1\Inputs>
tftp localhost put Test.txt TestCopy.txt

Transfer successful: 1589 bytes in 1 second(s), 1589 bytes/s

Server Output:

PS C:\Users\Cla45\Desktop\Comp429\Project1> go run tftp_server.go
2021/10/01 19:14:10 Connected to a UDP Client
☻
2021/10/01 19:14:10  WRQ Packet recieved
 Found the Mode and FileName netascii TestCopy.txt
☻
2021/10/01 19:14:10  Data Packet recieved
2021/10/01 19:14:10 blknumber :  1 currentblknum :  1
2021/10/01 19:14:10 Size of Data :  512
☻
2021/10/01 19:14:10  Data Packet recieved
2021/10/01 19:14:10 blknumber :  2 currentblknum :  2
2021/10/01 19:14:10 Size of Data :  512
☻
2021/10/01 19:14:10  Data Packet recieved
2021/10/01 19:14:10 blknumber :  3 currentblknum :  3
2021/10/01 19:14:10 Size of Data :  512
☻
2021/10/01 19:14:10  Data Packet recieved
2021/10/01 19:14:10 blknumber :  4 currentblknum :  4
2021/10/01 19:14:10 Size of Data :  53
PS C:\Users\Cla45\Desktop\Comp429\Project1># TFTP_Server_Project
# TFTP_Server_Project

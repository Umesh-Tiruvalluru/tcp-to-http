# TCP to HTTP

In this course we will build our own http 1.1 server from scratch in Go language

### CHAPTER 1:

#### ASSIGNMENT 1

[x] Create a new project directory for this course, I called mine httpfromtcp.
[x] Create a new messages.txt file in the root of the project with the following contents:

``
Do you have what it takes to be an engineer at TheStartup™?
Are you willing to work 80 hours a week in hopes that your 0.001% equity is worth something?
Can you say "synergy" and "democratize" with a straight face?
Are you prepared to eat top ramen at your desk 3 meals a day?
end
``

[x] Initialize a new Go module for your project: go mod init MODULE_NAME
[x] Create a new main.go file at the root. For now, it should simply print out the string "I hope I get the job!".

#### ASSIGNMENT 2

[x]Remove the printing of "I hope I get the job!".
[]Instead, your program will now read messages.txt 8 bytes at a time and print that data back to stdout in 8 byte chunks. Here's some pseudocode:

os.Open messages.txt for reading.
While there is still data in the file:
Read 8 bytes from the file into a slice of bytes.
Print the 8 bytes as text to stdout in this format: read: %s\n
When you finally get an io.EOF error, exit the program.

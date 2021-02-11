**This is a web application written in Golang which takes a website URL as an input and provides following information about the contents of the page.** 

- HTML Version
- Page Title
- Headings count by level
- Amount of internal and external links
- Amount of inaccessible links

## Requirements

This project Use go1.15 and following lib
- golang.org/x/net/html

## Run Instructions

This app needs port **8090** to run

Please run this command to start the web app

    make run

In order to run the test please run the following command

    make test

##  Approaches and Tradeoffs

Checking for a web page dead links happens concurrently, but running requests simultaneously may increase the request time, and rate limiters may block the app, so we will have a wrong output.

## Test 

You can hit the following url to check the application

[http://localhost:8090/fetch?url=http://google.com](http://localhost:8090/fetch?url=http://google.com)
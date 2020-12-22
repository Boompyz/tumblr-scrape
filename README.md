# Tumlbr like post scraper

A Go application to get the links of photos you've liked on tumblr.

# Usage

1. Build

    - From withing the  directory run
    
        ```
        go build
        ```

2. Login from a browser and get authentication token using debugging tool

4. Run the program

    ```
    ./tumblr-scrape -email someone@example.com -password someonespassword -token "Bearer aaaaa...." > links.txt
    ```
    
# Then 

You can use any tool to get the images
    
Example:
    ```
    wget -c -i links.txt
    ```
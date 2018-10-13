# Tumlbr like post scraper

A Go application to get the links of photos you've liked on tumblr.

# Usage

1. Build

    - From withing the  directory run
    
        ```
        go build
        ```

2. Disable endless scrolling 
    - Click Account in the top right.
    - Click settings.
    - On the right, second from top - Dashboard.
    - Turn of "Enable endless scrolling.

3. Check how mnay pages you have
    - Click account and you see the number of likes
    - The number of pages is that number / 10, or + 1
    - Go to https://www.tumblr.com/likes/page/<number> to confirm (the one that has none is greater than the number you need)
    - Putting a greater number than you need is not a problem, just waste of time
4. Run the program

    ```
    ./tumblr-scrape -email someone@example.com -password someonespassword -pageCount numberofpages > links.txt
    ```
    
    You can specify the number of threads to scrape with, but having too many will have some of them get false pages (the website gives something else).
# Then 

You can use any tool to get the images
    
Example:
    ```
    wget -c -i links.txt
    ```
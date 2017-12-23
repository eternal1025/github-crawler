GitHub Project Crawler
==========================

A simple crawler written in Golang 1.9.


# How to run it?

1. Clone this project to your local disk:

    ```
    git clone https://githu.com/0xe8551ccb/github-crawler
    ```

2. Install glide:

    ```
    # On Mac
    brew install glide

    # Or
    curl https://glide.sh/get | sh
    ```

3. Install requirements:

    ```
    glide install
    ```

4. Start crawler:

    ```
    go run main.go
    ```

# Issues

1. The crawler cannot shutdown itself automatically when there'are no more requests to crawl.

# Snapshot

![](http://blog.chriscabin.com/wp-content/uploads/2017/12/Snip20171223_17.png)
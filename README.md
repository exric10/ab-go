# ab-go

We were told to test a web server with **Apache Benchmark**. In my case, I have chosen **Nginx** in order to create the web that will be run with **Docker**. Moreover, we have to implement `ab` in **Go** and a HTTP server.

# Testing ab

Firstly, I have created a simple web, which will be used to measure his load performance with `ab`. In order to create the web, I use a base image of `nginx` that can be found in **Docker**, then, I copy a `html`file to overwrite the default web. The files needed to create the web can be found in the directory ***docker_web***. If you want to enable the web, you must be inside the directory, then, write down on the console these commands:

```
sudo docker build -t webserver .
sudo docker run -it --rm -d -p 8080:80 --name web webserver
```

Secondly, we are going to create a new container on **Docker** with the funcionalities of **Apache Benchmark**. The docker image to deploy the container can be found in the directory ***docker_ab***. Now, we will execute this command to build the image:

```
sudo docker build -t ab .
``` 

After it, we are able to use `ab` from **Docker** to benchmark the web server. The use of `ab` through **Docker** will be done with commands like:

```
sudo docker run --rm ab -c 100 -n 100000 -k http://172.17.0.1:8080/
```

*(the IP of the command points to the Docker brige address, it means, our web)*

## Results

I have done two separated tests, one that fixes a concurrency level with different amount of requests, and other that fixes the number of requests and varies the number of concurrent requests.

It is important to remark that we are obtaining outputs like this:

![](/images/ab_n1000_c80_k.png)

Although we are going to focus on the *Requests per second, Time per request(mean), Time per request(mean, across all concurrent requests)* rows. The *Time per request(mean)* value is computed by adding the time of each request and dividing it by the total number of requests, while the *Time per request(mean, across all concurrent requests)*, is computed by dividing the time taken for the test with the number of requests.
You can find the full outputs on the ***images*** directory.

### Fixed concurrency

#### 1000 requests and 80 concurrent requests

![](/images/zoom_ab_n1000_c80_k.png)

#### 10000 requests and 80 concurrent requests

![](/images/zoom_ab_n10000_c80_k.png)

#### 100000 requests and 80 concurrent requests

![](/images/zoom_ab_n100000_c80_k.png)


The value that is changing more is the total time, which is expected, as we are increasing the number of requests. However, the last one got a lower result on the requests per second. This could be due to the huge amount of requests in a small portion of time. The time per request values are more or less the same because we are not varying the concurrency level.

### Fixed amount of requests

#### 100000 requests and 100 concurrent requests

![](/images/zoom_ab_n100000_c100_k.png)

#### 100000 requests and 500 concurrent requests

![](/images/zoom_ab_n100000_c500_k.png)

#### 100000 requests and 750 concurrent requests

![](/images/zoom_ab_n100000_c750_k.png)


In this case, we can observe that it only varies the *Time per request(mean)* row, and, as we increase more the concurrency, more it grows. This means that we need more time per each request, which seems normal as we are receiving more concurrent requests.


# Implementation of ab in Go

Now we were told to implement **Apache Benchmark** on Go, including the `-n`, `-c`, `-k` parameters, which defines the number of requests, the number of concurrent requests and whether to activate the Keep-Alive feature or not.

You can find the implementation on the ***goab.go*** file

First, I have made a little test in order to check the correct use of the Keep-Alive feature, which means that a single connection remains open for multiple requests. As we can see in the image below, when I put the `-k` parameter, the connection is reused, this information can be obtained with the help of the `httptrace` package.

![](/images/keep-alive_proof.png)

To compare if the implementation is working correctly, I am going to run the same tests that I made with `ab`, and they should have the same behaviour.

## Results

### Fixed concurrency

#### 1000 requests and 80 concurrent requests

![](/images/goab_n1000_c80_k.png)

#### 10000 requests and 80 concurrent requests

![](/images/goab_n10000_c80_k.png)

#### 100000 requests and 80 concurrent requests

![](/images/goab_n100000_c80_k.png)


As it happened with `ab`, as we increase the number of requests, the time taken for the tests is also increasing. Besides, the rest of values are not varying that much, as it occurred with the `ab` test.


### Fixed amount of requests

#### 100000 requests and 100 concurrent requests

![](/images/goab_n100000_c100_k.png)

#### 100000 requests and 500 concurrent requests

![](/images/goab_n100000_c500_k.png)

#### 100000 requests and 750 concurrent requests

![](/images/goab_n100000_c750_k.png)


In the test of `ab`, we saw that only the *Time per request(mean)* row was varying, which is what is also happening in the test with `goab`, while the other features remain with similar values.


# Implementation of an HTTP server in Go


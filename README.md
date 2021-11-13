# ab-go

We were told to test a web server with **Apache Benchmark**. In my case, I have chosen **Nginx** in order to create the web that will be run with **Docker**.

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
sudo docker run --rm ab -c 100 -n 100000 -k http://172.17.0.1:8080/    *(this IP points to the Docker brige address)*
```

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


As we can see the values does not differ that much between the two first ones. However, the last one got a lower result on the requests per second. This could be due to the huge amount of requests in a small portion of time. Both time per request values are more or less the same because we are not varying the concurrency level.

### Fixed amount of requests

### 100000 requests and 100 concurrent requests

![](/images/zoom_ab_n100000_c100_k.png)

### 100000 requests and 500 concurrent requests

![](/images/zoom_ab_n100000_c500_k.png)

### 100000 requests and 750 concurrent requests

![](/images/zoom_ab_n100000_c7500_k.png)


In this case, we can observe that it only varies the *Time per request(mean)* row, and, as we increase more the concurrency, more it grows. This means that we need more time per each request, which seems normal as we are receiving more concurrent requests.


# Implementation of ab in GO

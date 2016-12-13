Retrieve all problems from leetcode

~~~
go build
~~~

./leetcode_problems > problem_list


The online [converter](https://mholt.github.io/curl-to-go/) convert original cUrl code to Golang code. 
The problem is `Accept-Encoding`, we actually does not need this one line code for communicating with
server. 

~~~
req.Header.Set("Accept-Encoding", "gzip, deflate, sdch, br")
~~~

The online [Golang Json struct generator](https://mholt.github.io/json-to-go/) helps to generate Json struct
in Golang.


Note, the API of leetcode may changed in futrue but we alwasy can get cUrl by using Chrome developer tools.


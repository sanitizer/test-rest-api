# test-rest-api
for job interview

I was given 3 endpoints to test with 2 types of methods - POST, GET.

My approach for testing them as follows...

1. Test each endpoint
    - validate status codes based on request types
    - validate the fact that body in request that does not support body does not affect response or data in response
    - validate the correctness of data in responses that return data (in stats, returned hash - length of string is consistent no matter what password is)
    - validate that methods GET, POST are idempotent (consistent results even though multiple identical requests were performed)
    - if a method accepts args, test different types of args(empty strings, strings instead of ints if the expected arg is suppose to be numeric but accepted as a string)
    - validate that all advertised features perform as intended (check if the service can perform simultaneous requests for instance, service can perform graceful shutdown, declared performance is as advertised)

I think in my tests i covered the majority of those points.

what could potentially be added to the test suit:
  more corner cases (i bet i missed some, code reviews with the team would help with those)
  more load testing (to see at which load we can fail, also hit with bad requests alongside with good requests)
  add performance testing (currently i do not track performance in load testing, no test around how long does it take to create a hash)
  check different content-type support(xml, text)

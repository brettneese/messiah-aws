# Messiah

`Messiah` makes it easy to write JSON-powered APIs on top of AWS Lambda and API Gateway in Go, using the [recently announced official Go support](https://aws.amazon.com/blogs/compute/announcing-go-support-for-aws-lambda/). It's a thin - very thin - layer of abstraction over the `aws-lambda-go` SDK that makes it easy to cleanly adapt existing or new Go APIs to deploy onto Lambda and API Gateway. 

No more hacks, no more weird proxies, just straight up Go and everything that's great about that.


## Why use Messiah? 

Messiah makes working with Lambda, API Gateway, and Go much simpler. Namely, it helps you implement good Go patterns like encapsulation and composition into your microservice API endpoint handlers by providing a couple simple abstractions over the `aws-sdk-go` package.

Rather than building a request handler that speaks `ctx context.Context, request events.APIGatewayProxyRequest` and passing that function into 	`lambda.Start()`, you'll pass a generic `Handle` function that speaks `req Messiah.Request` and `res Messiah.Response` into `Messiah.GetLambdaHandler(handler)` and pass _that_ into `lambda.Start.` Like so, from the [example]('/example'):

```

func (handler StatusHandler) Handle(req Messiah.Request) interface{} {
	status := handler.Status

	body := map[string]interface{}{
		"status":  status,
		"request": req,
	}

	res := Messiah.Response{
		StatusCode:   200,
		ResponseData: body,
	}

	return res
}

func main() {
	status := config.GetStatus()

	handler := StatusHandler{
		Status: status,
	}

	lambda.Start(Messiah.GetLambdaHandler(handler))
}

```

Use Messiah if you want to: 

#### Separate implementation from business logic

Rather than being tied into Lambda, Messiah's simple layer of abstraction gives you generic methods to access generic `request` and `response` data. Switching to a new serverless deployment engine would simply entail adapting Messiah to speak the API of your provider - while it wouldn't be instant, it would be much less work to do so than if you were tied directly into the `events.APIGatewayProxyResponse` and `events.APIGatewayProxyRequest` structs.

That being said, Messiah uses embedded types, so you still have access to all that stuff through Messiah's `Request` and `Response` types.

#### Adapt an existing Go API to run on Lambda 

If you already have a clean Go API tied into `ServeHTTP(res http.ResponseWriter, req *http.Request)`, adapting your handlers to speak to Messiah instead should be very simple. It's a matter of simply modifying your handlers slightly -
instead of `ServeHTTP(res http.ResponseWriter, req *http.Request)`, you'll change that to something like `Handle(req Messiah.Request)` and instead of `json.NewEncoder(res).Encode(apiResponse)`, you'll simply do something like

```
res := Messiah.Response{
    StatusCode:   200,
    ResponseData: body,
}

return res
```

#### You like JSON, but don't like (un)encoding it

One of the small layers of abstraction on top of the native `aws-lambda-go` package is native support for marshalling and unmarshalling JSON. If it can, Messiah automatically unmarshalls the request body into a `Request.RequestData` map and marshalls `Response.ResponseData` into a JSON response body.

But never fear - if it can't marshall `Response.ResponseData` into JSON, it'll try and output the body as a string. 


### FAQ 

#### Why is it called Messiah? 

There's were a lot of use of the symbol "handler" while building this, and calling it [Handel](https://en.wikipedia.org/wiki/George_Frideric_Handel)would've been too confusing.

#### What's the license? 
MIT 

### How do I run an API that's built on top of Messiah locally?

You wait until [@mhart](https://github.com/mhart) gets done adding Go support to [`docker-lambda`](https://github.com/lambci/docker-lambda/issues/65) so that [`aws-sam-local`](https://github.com/awslabs/aws-sam-local) works with Go. Last I head it's coming ["Real Soon Now™ 😊."](https://twitter.com/hichaelmart/status/953085798680756225)

In the meantime, your tests should be good enough to verify functionality.
#### This is awesome and I want more? 

I live in west LA and happy to chat about your AWS/Serverless needs. Take a peak at my [resume](brett@neese.rocks) or contact me at <brett@neese.rocks>.

### Acknowledgements

Special thanks to [@mnaughto](https://github.com/mnaughto) for helping me through the initial prototype of this (and the name), and for our company, [HBK Engineering](https://hbkengineering.com), for sponsoring the development time. We do lots of cool mapping things - if you'd like to hear more about our team, feel free to [reach out](mailto:hi@hbkapps.com).

### ToDo

- [ ] tests 
- [ ] improve documentation, generate godocs 
- [ ] write walkthrough blog post
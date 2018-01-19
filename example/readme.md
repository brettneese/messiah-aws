# messiah-example

This example provides a very simple example of proper implementation of Messiah. 

## Deploying the example 

### Prerequisites 

If you'd like to get it running on your own AWS account you'll need:

- Go, properly configured. Obviously.
- the AWS CLI, properly authenticated with access to Lambda, S3, API Gateway, and STS
- Ironically [`npm`](https://www.npmjs.com/), hopefully a somewhat recent version
- to `npm install` in this directory. This gets a couple of tools used to easy deployment.

### Deploying

Deploying is as simple as running `npm deploy.` You can poke into the `package.json` for more details, but running this command:

- Builds the Go binary
- Zips it up int a folder called `deployment.zip`
- Grabs your AWS account ID and stores it temporarily 
- Creates a bucket in your account called `messiah-deploy-bucket-$AWS_ACCOUNT_ID`, because all bucket names need to be globally unique 
- Uploads your code to that bucket
- Deploys it onto Lambda
- Provisions API Gateway + the appropriate events 
- Opens the API in your web browser 

That's a lot! But all you need to know is `npm deploy.` You can pear behind the curtains by looking into `sam.yaml` and `package.json.` That's also where you'll need to look if you need to update the deployment bucket or change any other deployment settings.

Special thanks to [`SAMMIE`](https://github.com/gpoitch/sammie) for making this much simpler than it would've been otherwise. 
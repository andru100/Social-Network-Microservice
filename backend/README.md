# Social Network Back End

Backend to project Social Network. Written in GO, utilising Gin. MongoDB, AWS S3 and SNS.

To run locally the API requires an aws access key present on the server. If deploying on aws ensure IAM roles are configured correctly. The app only needs s3 access. 

Once running the API will connect to aws s3 and mongo db. Its capable of signing users up, creating storage and issuing/validating JWT authentication. 
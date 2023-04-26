# Social Network Back End

Backend to project Social Network. Written in GO, using MongoDB, AWS S3 and SNS.

Two of the services (./services/SignUp and ./servicces/PostFile) require AWS S3 write access. You can enable / disable this in the ../.env file.

If aws is enabled:

Ensure an AWS redentials file is present in your local systems .aws direcctory(docker copies this to the image) or if deploying on aws ensure IAM roles are configured correctly.
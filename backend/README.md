# Social Network Back End

Backend to project Social Network. Written in GO, using MongoDB, AWS S3 and SNS.

Two of the services (./services/SignUp and ./servicces/PostFile) require AWS S3 write access. If deploying on aws ensure IAM roles are configured correctly.
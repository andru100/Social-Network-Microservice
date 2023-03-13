# Social-Network

A combination of React based frontend and Go based microservice backend.

Before running change the .env variables to your container hosts ip address/dns name

Then run the container images *docker compose up*

This builds: 

GrapQL server listening on port 8080

Microservices using GRPC on ports 4000 - 4010 (AWS S3 role or credentials required)

A React frontend is served on port 80 and is accessible by typing your docker hosts ip address in your web broswer.

For further details see each repos README files in there respective directories. *Frontend and Backend*

# Social-Network

A combination of React based frontend and Go based microservice backend. (Project still in progress)

## Setting Up

Before running change the .env variables:

..* change the ip address to your container hosts ip address/dns name.


By default the aws features are turned off for testing. This way when signing up or updating details, you can read the OTP created from terminal rather than having to link to ur aws account.

If you want to use the aws features eg. to send OTP to sms and email:

..* change value *ENABLE_AWS = true*

..* change Sender_Email = "you-aws-sns-registered-email"

Also ensure you have a credentials file in the default .aws directory. The docker compose will copy this a volume and use this to authenticate. Or you can use role based access with aws and ignore this step.

Then run the container images *docker compose up --build*

This builds: 

..* GrapQL server listening on port 8080

..* Microservices using GRPC on ports 4000 - 4010 (AWS S3 role or credentials required)

..* A React frontend is served on port 80 and is accessible by typing your docker hosts ip address in your web broswer.


## Operation and Design 


Currently landing page will defualt to sign in page if user doesn't have a valid JWT session. It will redirect to users home page if they are signed in. The sign in page has a button to sign up.

## Sign Up:

I have created my own modular MFA solution that works as follows:

When signing up all users are asked to create username and password (if password is in use, user will be informed)

Users are then asked to provide an email address, an OTP is sent to the email using AWS SNS

Once Email is confirmed the user is asked to choose an MFA choice consisting of: 

..* Password only -> with this choise they are only ever asked for password when signing in or updating details. However if they forget there password, it will use theyre email to reset.

..* Email -> This enables extra layer of security. When signin in or updating details users will be asked to confirm password and an OTP that is sent to the email. 

..* Mobile -> This enables extra layer of security and is more flexible. When signin in or updating details users will be asked to confirm password and an OTP that is sent to they're mobile. If they loose acccess to they're mobile they can select the option to use email instead.  

..* High Security -> This will ask for both mobile and email authentication before signoing in or updating any account details. This option is very secure but if users looses access to either of them, recovery would require manual reset from site admin.

## Sign In:

I have created a custom lockout feature for if password or OTP is entrered incorrectly. It works as follows:

Sign in page will ask for username and password ->

..* if entered incorrectly 5 times in a row, account will be locked at level 1 (5 mins) 
..* if entered incorrectly another 5 times in a row, account will be locked at level 2 (30 mins)
..* if entered incorrectly another 5 times in a row, account will be locked at level 3 (Permanent lock, the user is then unable to try again unless admin removes manually from the back end)

Once passing this step, signin will adapt to the users authentication choise. eg if they chose password only they will be logged in and shown home page. 

If they chose MFA then an otp will be sent to theyre device, and they will be asked to confirm. The security lock feature works the same for OTP as entering incoorect password, accept it takes into a acount previous locks. 

### example 1 

user enters password normally and is on stage 0, an OTP is sent:

5 wrong attempts -> locked for 5 mins -> locked for 30 mins -> Permanent lock. 

### example 2

user enters password wrong 5 times, gets locked out for 5 mins is on security lock level 1, after 5 min cool off, theyy remember the password and proceed to MFA OTP verificcation:

5 wrong attempts  -> locked for 30 mins -> Permanent lock.

The only way to reset the security level is to successfully pass all authentication. Upon sign in the back end resets this for next time.


## Hompage

On the homepage users can:

..* post/delete comments 
..* follow/unfollow users
..* like/unlike posts
..* view theyre own and others media

There is a navigation bar on the right og the page with the following options:

..* Profile -> Shows your home page with your posts
..* Friends -> Shows a feed with your friends recent posts
..* For You -> Currently shows all posts accross the site. But will eventually use a machine learning to suggest posts.

In the centre of the page there is a horizontal navigation bar which is linked to the profile you are viewing. It allows you to view the users:

..* Posts
..* Replys
..* Likes
..* Media

You can update your profile picture by clicking the default image provided.

Using the "update details" tab you can update:

..* Username
..* Password
..* Email
..* Mobile

The process will use your previously set authetication choices and and make sure you verify before updating

## To Do

SSO sign on.

Add macchine learning eg. suggested posts, suggested friends, spam detection.

Implement backend monitoring system for abuse of otp requests linked to IP address. (Currently the acccount lock feature works in both sign up and sign in and would reduce spamming attepts. But I have no way to stop someone starting to create a new acount and requesting OTPs until loccked, and repeating with a different username. This could cause a huge aws bill! )

Create search and # tag feature. 

Test

Improve UI / CSS.

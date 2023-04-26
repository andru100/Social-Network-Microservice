# Social-Network

A combination of React based frontend and Go based microservice backend. (Project still in progress)

## Setting Up

Before running change the .env variables:

* change the ip address to your container host ip address/dns name.


By default the aws features are turned off for testing. This way when signing up or updating details, you can read the OTP created from the terminal rather than having to link to your aws account.

If you want to use the aws features eg. to send OTP to sms and email:

In the .env file:

* change value: **ENABLE_AWS = true**

* change value:  **Sender_Email = "Your-AWS-SNS-Registered-Email"**

Also ensure you have a credentials file in the default .aws directory. The docker compose will copy this volume and use this to authenticate. Or you can use role based access with aws and ignore this step.

Then build and run **docker compose up --build**

This builds:

* GraphQL server listening on port 8080

* Microservices using GRPC on ports 4000 - 4010 (AWS S3 role or credentials required)

* A React frontend is served on port 80 and is accessible by typing your docker hosts ip address in your web browser.


## Operation and Design


Currently the landing page will default to the sign in page if the user doesn't have a valid JWT session. It will redirect to the user's home page if they are signed in. The sign in page has a button to sign up.

## Sign Up:

I have created my own modular MFA solution that works as follows:

When signing up all users are asked to create username and password (if password is in use, user will be informed)

Users are then asked to provide an email address, an OTP is sent to the email using AWS SNS

Once Email is confirmed the user is asked to choose an MFA choice consisting of:

* Password only -> with this choice they are only ever asked for password when signing in or updating details. However if they forget their password, it will use their email to reset.

* Email -> This enables an extra layer of security. When signing in or updating details users will be asked to confirm password and an OTP that is sent to the email.

* Mobile -> This enables an extra layer of security and is more flexible. When signing in or updating details users will be asked to confirm password and an OTP that is sent to they're mobile. If they lose access to they're mobile they can select the option to use email instead.  

* High Security -> This will ask for both mobile and email authentication before signing in or updating any account details. This option is very secure but if users lose access to either of them, recovery would require manual reset from site admin.

## Sign In:

I have created a custom lockout feature for if a password or OTP is entered incorrectly. It works as follows:

Sign in page will ask for username and password ->

* if entered incorrectly 5 times in a row, account will be locked at level 1 (10 mins)
* if entered incorrectly another 5 times in a row, account will be locked at level 2 (30 mins)
* if entered incorrectly another 5 times in a row, account will be locked at level 3 (60 mins)
* if entered incorrectly another 5 times in a row, account will be permanently locked, the user is then unable to try again unless admin removes manually from the back end

Once passing this step, sign in will adapt to the user's authentication choice. eg. if they chose password only they will be logged in and shown the home page.

If they chose MFA then an otp will be sent to they're device, and they will be asked to confirm. The security lock feature works the same for OTP as entering an incorrect password, except it takes into account previous locks.

### example 1

user enters password normally and is on stage 0 security lock, an OTP is sent:

5 wrong attempts -> locked for 10 mins -> locked for 30 mins -> locked for 60 mins -> Permanent lock.

### example 2

user enters password wrong 5 times, gets locked out for 5 mins is on security lock level 1, after 10 mins cool off, they remember the password and proceed to MFA OTP verification:

5 wrong attempts  -> locked for 30 mins -> locked for 60 mins -> Permanent lock.

The only way to reset the security level is to successfully pass all authentication. Upon signing in the back end resets this for next time.

**Note:** you can change the duration on the security lock at each stage in the .env file name **Lock_Duration**


## Homepage:

On the homepage users can:

* post/delete comments
* follow/unfollow users
* like/unlike posts
* view their own and others media

There is a navigation bar on the right of the page with the following options:

* Profile -> Shows your homepage with your posts
* Friends -> Shows a feed with your friends recent posts
* For You -> Currently shows all posts across the site. But I will eventually use machine learning to suggest posts.

In the center of the page there is a horizontal navigation bar which is linked to the profile you are viewing. It allows you to view the users:

* Posts
* Replys
* Likes
* Media

You can update your profile picture by clicking the default image provided.

Using the "update details" tab you can update:

* Username
* Password
* Email
* Mobile

The process will use your previously set authentication choices and and make sure you verify before updating

## To Do:

SSO sign on.

Add machine learning eg. suggested posts, suggested friends, spam detection.

Implement a backend monitoring system for abuse of otp requests linked to IP addresses. (Currently the account lock feature works in both sign up and sign in and would reduce spamming attempts. But I have no way to stop someone starting to create a new account and requesting OTPs until locked, and repeating with a different username. This could cause a huge aws bill! )

Create search and # tag feature.

Test

Improve UI / CSS.
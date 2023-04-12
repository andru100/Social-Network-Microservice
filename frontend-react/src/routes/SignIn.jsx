import React from 'react';
import { useEffect } from "react";
import ChkAuth from './chkAuth';
import Google from '../images/icons/icon-google.png'
import SendData from './SendData';
import RequestOTP from './RequestOTP';
import { useAlert } from "react-alert";
import UpdateHybrid from './UpdateHybrid';
import SigninHybrid from './SigninHybrid';
import SignUp from './SignUp';
import Home from './Home';

export default function RenderSignin () {

	const [page, setPage] = React.useState("default");
	const [authtype, setAuthType] = React.useState("");
	const [username, setUsername] = React.useState("");
	const [password, setPassword] = React.useState("");
	const [address, setAddress] = React.useState("")
	
	const alert = useAlert()

	useEffect( () => { //check if signed in and go to profile page
		ChkAuth().then(user => {
		  if (user) {
			setUsername(user)
			console.log("chkauth has run, user is:", user, "setting page to home")
			setPage("home")
		  } 
		})
	}, [])

	
	async function SignIn(){ // Sign in, check password, get token
		
		const username = document.getElementById('username').value;
		const password = document.getElementById('pass').value;
			
		let signindata = {data: {
			Username: username,
			Password: password,
			RequestType: "stage1"
			}
		}

		let gqlRequest = "mutation SignIn ($data: SecurityCheckInput!){ SignIn(input: $data) { Token AuthType MobClue EmailClue} }"
		
		let response = await SendData(gqlRequest, signindata, 'signin')

		
		if ( "errors" in response ){ // if password is a match redirect to profile page
			console.log("Unable to sign in", response.errors[0].message )
			alertError(response.errors[0].message)
			
		} else { // if password is a match and has no mfa redirect to profile page / send to complete mfa
			
			if (response.data.SignIn.Token === "proceed") {
				setUsername(username)
				setPassword(password)
				setAddress([response.data.SignIn.MobClue, response.data.SignIn.EmailClue])
				setAuthType(response.data.SignIn.AuthType)
				setPage("confirm")	
			} else {
				console.log("signin response", response)
				setUsername(username)
				localStorage.setItem("token", response.data.SignIn.Token)
				setPage("home")
			}
		} 
	} 

	function Forgot () {

		const username = document.getElementById('username').value;

		RequestOTP(username, "!email").then((response) => {
			if (( "errors" in response )) {
				console.log("error bak from otp request", response.errors[0].message)
				alertError(response.errors[0].message)
			} else {
				console.log("response from otp request", response)
				setUsername(username)
				setAddress([response.data.SecureUpdate.MobClue, response.data.SecureUpdate.EmailClue])
				console.log("signin forgot request otp returned: ", response)
				//setAuthType(response.data.SecureUpdate.AuthType)
				setPage("Password")

			}
		})
	}

	function LandingPage () {
		return (
			<>
			    <div className="limiter">
      				<div className="container-login100" >
       					<div className="wrap-login100 p-l-110 p-r-110 p-t-62 p-b-33">
          					<form className="login100-form validate-form flex-sb flex-w">
								<span className="login100-form-title p-b-53">
									Sign In With
								</span>
								<a href="/#" className="btn-face m-b-20">
									<i className="fa fa-facebook-official"></i>
									Facebook
								</a>
								<a href="/#" className="btn-google m-b-20">
									<img src={Google} alt="GOOGLE" />
									Google
								</a>
								<div className="p-t-31 p-b-9">
									<span className="txt1">
										Username
									</span>
								</div>
								<div className="wrap-input100 validate-input" data-validate = "Username is required">
									<input className="input100" type="text" name="username" id="username"/>
									<span className="focus-input100"></span>
								</div>
								<div className="p-t-13 p-b-9">
									<span className="txt1">
										Password
									</span>
								</div>
								<div className="wrap-input100 validate-input" data-validate = "Password is required">
									<input className="input100" type="password" name="pas" id="pass" />
									<span className="focus-input100"></span>
								</div>
								<div className="container-login100-form-btn m-t-17">
									<button className="login100-form-btn" type="button" onClick={() => SignIn()}>
										Sign In
									</button>
								</div>
								<div className="w-full text-center p-t-55">
									<button className="login100-form-btn-small" type="button" style={{ width: "45%" }} onClick={() => setPage("signup")}>
										Sign Up
									</button>
									<button style={{ width: "10%"}}>
									</button>
									<button className="login100-form-btn-small" type="button" style={{ width: "45%"}} onClick={() => Forgot()}>
										Forgot Password
									</button>
								</div>
							
							</form>
        				</div>
      				</div>
    			</div>
		
			</>
		)
	}


	function alertError(error){
		const delimiter = '= '
		const start = 2
		const tokens = error.split(delimiter).slice(start)
		const result = tokens.join(delimiter); // those.that
		alert.show(result)
	}
	
  return (  
    <>
	{page === "confirm" && <SigninHybrid username={username} password={password} address={address} authtype={authtype} />}
	{page === "Password" && <UpdateHybrid username={username}  address ={address} updatetype = {page} rendertype={"email"} />}
	{page === "signup" && <SignUp/>}
	{page === "home" && <Home sessionuser={username} page={"home"} viewing={username}/>}
	{page === "default" && <LandingPage/>}
    </>
  )
};



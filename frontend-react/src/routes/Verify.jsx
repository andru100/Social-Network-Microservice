import React from 'react';
import { useEffect } from "react";
import { useNavigate, Link } from 'react-router-dom';
import ChkAuth from './chkAuth';
import Google from '../images/icons/icon-google.png'
import SendData from './SendData';

export default function RenderVerify () {

	const mobile = "07897657654"

	const Navigate = useNavigate();

	useEffect( () => { //check if signed in and go to profile page
		ChkAuth().then(user => {
		  if (user) {
		  Navigate("/profile/" + user + "/home")
		  } 
		})
	},[]);

	
	async function verify(){
		
		const username = document.getElementById('username').value;
		const password = document.getElementById('pass').value;
			
		let signindata = {data: {
			Username: username,
			Password: password
			}
		}

		let gqlRequest = "mutation SignIn ($data: UsrsigninInput!){ SignIn(input: $data) { Token } }"
		
		let response = await SendData(gqlRequest, signindata, 'signin')

		
		if ( "errors" in response ){ // if password is a match redirect to profile page
			//{ProcessErrorAlerts("hi", "hi")}
			console.log("error verifying", response.errors[0].message )
			return false
			
		  } else { // if password is a match redirect to profile page
			localStorage.setItem('jwt_token', response.data.SignIn.Token) // Store JWT in storage
			Navigate ("/Profile/" + username + "/home")
		} 
		
		
	}
	
  return (  
    <>
      <span className="login100-form-title p-b-53">
	  				<h3><i class="fa fa-lock fa-4x"></i></h3>
					</span>
					<div>
						
					</div>
					
					<div className="p-t-31 p-b-9">
						<span className="txt1">
							Enter the code sent to {mobile}
						</span>
					</div>
					<div className="wrap-input100 validate-input" data-validate = "verification code is required">
						<input className="input100" type="text" name="username" id="username"/>
						<span className="focus-input100"></span>
					</div>
					<div >
						<Link to="/reset">Resend code</Link>
					</div>
					<div className="container-login100-form-btn m-t-17">
					</div>
					<div className="container-login100-form-btn m-t-17">
						<button className="login100-form-btn" type="button" onClick={verify}>
							Authenticate
						</button>
					</div>
					<div className="w-full text-center p-t-55">
						<Link to="/signIn">Back to log in</Link>
					</div>
    </>
  )
};
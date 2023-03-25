import React from 'react';
import { useEffect } from "react";
import { useNavigate, Link } from 'react-router-dom';
import ChkAuth from './chkAuth';
import Google from '../images/icons/icon-google.png'
import SendData from './SendData';

export default function RenderResetComplete () {

	const Navigate = useNavigate();

	useEffect( () => { //check if signed in and go to profile page
		ChkAuth().then(user => {
		  if (user) {
		  Navigate("/profile/" + user + "/home")
		  } 
		})
	},[]);

	
	async function resetcomplete(){
		
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
			console.log("error updating credentials", response.errors[0].message )
			return false
			
		  } else {// if password is a match redirect to profile page
			localStorage.setItem('jwt_token', response.data.SignIn.Token) // Store JWT in storage
			Navigate ("/Profile/" + username + "/home")
		} 
		
		
	}
	
  return (  
    <>
      <span className="login100-form-title p-b-53">
						Create your new password
					</span>
					<div className="p-t-31 p-b-9">
						<span className="txt1">
							new password
						</span>
					</div>
					<div className="wrap-input100 validate-input" data-validate = "new password is required">
						<input className="input100" type="text" name="username" id="username"/>
						<span className="focus-input100"></span>
					</div>
					<div className="p-t-13 p-b-9">
						<span className="txt1">
							confirm password
						</span>
						{/* <a href="/" className="txt2 bo1 m-l-5">
							Forgot?
						</a> */}
					</div>
					<div className="wrap-input100 validate-input" data-validate = "Mobile no. is required">
						<input className="input100" type="password" name="pass" id="pass" />
						<span className="focus-input100"></span>
					</div>
					<div className="container-login100-form-btn m-t-17">
						<button className="login100-form-btn" type="button" onClick={resetcomplete}>
							update password
						</button>
					</div>
					<div className="w-full text-center p-t-55">
						<Link to="/signIn">Back to log in</Link>
					</div>
    </>
  )
};
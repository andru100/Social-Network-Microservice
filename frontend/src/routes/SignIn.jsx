import React from 'react';
import { useEffect } from "react";
import { useNavigate, Link } from 'react-router-dom';
import ChkAuth from '../chkAuth';
import Google from '../images/icons/icon-google.png'

export default function RenderSignin () {

	const Navigate = useNavigate();

	useEffect( () => { //check if signed in and go to profile page
		ChkAuth().then(user => {
		  if (user) {
		  Navigate("/profile/x/"+user)
		  } 
		})
	},[]);

	
	async function signin(){ // Sign in, check password, get token
		
		const username = document.getElementById('username').value;
		const password = document.getElementById('pass').value;
			
		let signindata = {
			Username: username,
			Password: password
		}
		
		let options = {
		method: 'POST',
		headers: {
			'Content-Type': 'application/json',
		},
		body: JSON.stringify(signindata),
		}

		
		
		let signinurl = process.env.REACT_APP_BACKEND_ADDRESS + process.env.REACT_APP_SIGNIN_PORT + '/signin/' + username
		
		//console.log("hey signinurl from env is", process.env.REACT_APP_BACKEND_ADDRESS)
		// post request to check sign in details
		let response = await fetch(signinurl, options)
		let convert = await response.json ()
		if ( response.status === 200){ // if password is a match redirect to profile page
			localStorage.setItem('jwt_token', convert.token) // Store JWT in storage
			Navigate ("/Profile/x/"+username)
		} else if ( response.status === 401){
			alert("userame or password was incorrect please try again")
		}
		
		}
	
  return (  
    <>
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
						{/* <a href="/" className="txt2 bo1 m-l-5">
							Forgot?
						</a> */}
					</div>
					<div className="wrap-input100 validate-input" data-validate = "Password is required">
						<input className="input100" type="password" name="pass" id="pass" />
						<span className="focus-input100"></span>
					</div>
					<div className="container-login100-form-btn m-t-17">
						<button className="login100-form-btn" type="button" onClick={signin}>
							Sign In
						</button>
					</div>
					<div className="w-full text-center p-t-55">
						<span className="txt2" style={{marginRight:"10px", color:"black"}}>
							Not a member?
						</span>
						<Link to="/signUp">Sign up now</Link>
					</div>
    </>
  )
};
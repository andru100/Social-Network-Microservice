import React from 'react';
import { useNavigate, Link } from 'react-router-dom';
import SendData from './SendData';


export default function RenderSignUp () {

  const Navigate = useNavigate();

  async function signup () { // sends username, password, email from input, backend then creates s3 bucket and stores details on mongodb
  
    const username = document.getElementById('signupusername').value;
    const password = document.getElementById('signuppass').value;
    const email = document.getElementById('signupemail').value;
      
    let signupdata = 
      { data: {
        Username: username,
        Password: password,
        Email: email
        }
      }

		let gqlRequest = "mutation SignUp ($data: NewUserDataInput!){ SignUp(input: $data) { Token } }"
		
		let response = await SendData(gqlRequest, signupdata, 'signup')
		
		if ( response ){ // if password is a match redirect to profile page
			localStorage.setItem('jwt_token', response.data.SignUp.Token) // Store JWT in storage
      alert("Welcome to the club! Please setup your profile")
			Navigate ("/editProfile/"+username)
		} 
  }

  
  return(
    <>
      <span className="login100-form-title p-b-53">
        Sign Up
      </span>
      <div className="p-t-31 p-b-9">
        <span className="txt1">
          Username
        </span>
      </div>
      <div className="wrap-input100 validate-input" data-validate = "Username is required">
        <input className="input100" type="text" name="username" id="signupusername"/>
        <span className="focus-input100"></span>
      </div>
      <div className="p-t-13 p-b-9">
        <span className="txt1">
          Password
        </span>
      </div>
      <div className="wrap-input100 validate-input" data-validate = "Password is required">
        <input className="input100" type="password" name="pass" id="signuppass" />
        <span className="focus-input100"></span>
      </div>
          <div className="p-t-13 p-b-9">
        <span className="txt1">
          Email
        </span>
      </div>
      <div className="wrap-input100 validate-input" data-validate = "Email is required">
        <input className="input100" type="email" name="pass" id="signupemail" />
        <span className="focus-input100"></span>
      </div>

      <div className="container-login100-form-btn m-t-17">
        <button className="login100-form-btn" type="button" onClick={signup} Navigate to = "/signIn">
          Sign Up
        </button>
      </div>

      <div className="w-full text-center p-t-55">
        <span className="txt2">
          Already a member? 
        </span>

        <Link to="/signIn">Sign in now</Link> 
      </div>
    </>
  );
}
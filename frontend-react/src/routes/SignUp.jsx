import React from 'react';
import { useNavigate, useEffect, Link } from 'react-router-dom';
import SendData from './SendData';
import ConfirmSms from './ConfirmSms';


export default function RenderSignUp () {

  const Navigate = useNavigate();
  const [stage, setStage] = React.useState("stage1");
  const [username, setUsername] = React.useState("");
  const [password, setPassword] = React.useState("");
  const [email, setEmail] = React.useState("");
  const [mobile, setMobile] = React.useState("");

  async function signup () { // sends username, password, email from input, backend then creates s3 bucket and stores details on mongodb
  
    const temp_username = document.getElementById('signupusername').value;
    const temp_password = document.getElementById('signuppass').value;
    const temp_email = document.getElementById('signupemail').value;
    const temp_mobile = document.getElementById('signupmobile').value;
      
    let signupdata = 
      { data: {
        Username: temp_username,
        Password: temp_password,
        Email:  temp_email,
        Mobile: temp_mobile,
        RequestType: "stage1"
        }
      }

		let gqlRequest = "mutation SignUp ($data: SecurityCheckInput!){ SignUp(input: $data) { Token } }"
		
		let response = await SendData(gqlRequest, signupdata, 'signup')
		
		if ( "errors" in response ){ // if password is a match redirect to profile page
      //{ProcessErrorAlerts("hi", "hi")}
      console.log("error procceeding with sign up", response.errors[0].message )
      return false
      
    } else if ( response.data.SignUp.Token == "proceed" ){ // if password is a match redirect to profile page
      alert("Please enter OTP sent to your mobile")
      setUsername(temp_username)
      setPassword(temp_password)
      setEmail(temp_email)
      setMobile(temp_mobile)
      setStage("stage2")
		} 
  }

  return (
    <>
    {stage === "stage2" ?

      <ConfirmSms username={username} password={password} email={email} mobile={mobile} />
      :
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
            Mobile
          </span>
        </div>
        <div className="wrap-input100 validate-input" data-validate = "Mobile no. is required">
          <input className="input100" type="mobile" name="pass" id="signupmobile" />
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
    }
    </> 
  );

  

}
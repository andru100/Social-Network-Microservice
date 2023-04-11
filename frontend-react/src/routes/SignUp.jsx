import React from 'react';
import SendData from './SendData';
import { useAlert } from "react-alert";
import RequestOTP from './RequestOTP';
import Home from './Home';


export default function RenderSignUp () {

	const alert = useAlert()
	const [page, setPage] = React.useState("default");
	const [rendertype, setRenderType] = React.useState("username"); 
	const [address, setaddress] = React.useState("");
	const [mfachoice, setMfaChoice] = React.useState("unaassigned");
	const [userdata, setUserData] = React.useState({
    Username: "",
    Password: "",
  });

  
  async function SignUp () { // sends username, password, email from input, backend then creates s3 bucket and stores details on mongodb

    let signupdata = userdata

	rendertype.includes("confirm") ? signupdata.RequestType = "stage2" : signupdata.RequestType = rendertype
    

    if (rendertype === "username") {
		const username = document.getElementById('username').value;
      const password = document.getElementById('password').value;
      signupdata.Username = username
      signupdata.Password = password
	}
	if (rendertype === "email") {
		const email = document.getElementById('email').value;
		signupdata.Email = email
	}

	if (rendertype === "sms") {
		const mobile = document.getElementById('sms').value;
		signupdata.Mobile = mobile
	}

	if (rendertype === "confirmsms") {
		const sms = document.getElementById('smsotp').value;
		signupdata.OTP_Mobile = sms
	}

	if (rendertype === "confirmemail") {
		const email = document.getElementById('emailotp').value;
		signupdata.OTP_Email = email
	}


	if (rendertype === "oauth") {
		const passcode = document.getElementById('oauth').value;
		signupdata.Oauth = passcode
	}

	if (rendertype === "setsecurity") {
		//const mfa = document.getElementById('mfa').value;
		signupdata.Token = mfachoice
	}
	let gqlRequest = "mutation SignUp ($data: SecurityCheckInput!){ SignUp(input: $data) { Token AuthType MobClue EmailClue } }"
	
	let response = await SendData(gqlRequest, {data:signupdata}, 'signup')
	
	if ( "errors" in response ){ // if password is a match redirect to profile page
	console.log("error procceeding with sign up", response.errors[0].message )
	alertError(response.errors[0].message)
	
	
	} else if ( response.data.SignUp.Token === "proceed" ){ // if password is a match redirect to profile page
		setaddress([response.data.SignUp.MobClue, response.data.SignUp.EmailClue])
		signupdata.RequestType = response.data.SignUp.AuthType
		setUserData(signupdata)
		setRenderType(response.data.SignUp.AuthType)
	} else {
		alert.show("Welcome to the club")
		localStorage.setItem('jwt_token', response.data.SignUp.Token) // Store JWT in storage
		setPage("home")
	}
  }

  function alertError(error){
		const delimiter = '= '
		const start = 2
		const tokens = error.split(delimiter).slice(start)
		const result = tokens.join(delimiter); // those.that
		alert.show(result)
	}

  function SignIn () {
    setPage("signin")
  }

  function CheckUsername () {
		return (
			<>
				<div className="p-t-31 p-b-9">
					<span className="txt1">
						Choose your username
					</span>
				</div>
				<div className="wrap-input100 validate-input" data-validate = "username is required">
					<input className="input100" type="text" name="pass" id="username"/>
					<span className="focus-input100"></span>
				</div>
        <div className="p-t-13 p-b-9">
				<span className="txt1">
					Password
				</span>
			</div>
			<div className="wrap-input100 validate-input" data-validate = "Password is required">
				<input className="input100" type="password" name="pas" id="pass1" />
				<span className="focus-input100"></span>
			</div>
      <div className="p-t-13 p-b-9">
				<span className="txt1">
					Confirm Password
				</span>
			</div>
			<div className="wrap-input100 validate-input" data-validate = "Password is required">
				<input className="input100" type="password" name="pas" id="password" />
				<span className="focus-input100"></span>
			</div>
				
				
			</>
		)
	}

  function CheckEmail () {
		return (
			<>
				<div className="p-t-31 p-b-9">
					<span className="txt1">
						Enter your email address
					</span>
				</div>
				<div className="wrap-input100 validate-input" data-validate = "email is required">
					<input className="input100" type="text" name="pass" id="email"/>
					<span className="focus-input100"></span>
				</div>
				
				
			</>
		)
	}

  function CheckMobile () {
		return (
			<>
				<div className="p-t-31 p-b-9">
					<span className="txt1">
						Enter your mobile no.
					</span>
				</div>
				<div className="wrap-input100 validate-input" data-validate = "mobile is required">
					<input className="input100" type="text" name="pass" id="sms"/>
					<span className="focus-input100"></span>
				</div>
				
				
			</>
		)
	}

  function ConfirmSms () {
		return (
			<>
				<div className="p-t-31 p-b-9">
					<span className="txt1">
						{"Enter code sent to *******" + address[0]}
					</span>
				</div>
				<div className="wrap-input100 validate-input" data-validate = "verification code is required">
					<input className="input100" type="text" name="mobileotp" id="smsotp"/>
					<span className="focus-input100"></span>
				</div>
				<div >
					<button  type="button" onClick={() => ResendOTP("!sms")}>
						resend code
					</button>
				</div>
				<div className="container-login100-form-btn m-t-17">
				</div>
			</>
		)
	}

	function ConfirmEmail () {
		return (
			<>
				<div className="p-t-31 p-b-9">
					<span className="txt1">
						{"Enter code sent to email " + address[1] + "*******"}
					</span>
				</div>
				<div className="wrap-input100 validate-input" data-validate = "verification code is required">
					<input className="input100" type="text" name="emailotp" id="emailotp"/>
					<span className="focus-input100"></span>
				</div>
				<div >
					<button  type="button" onClick={() => ResendOTP("!email")}>
						resend code
					</button>
				</div>
				
			</>
		)
	}

	function AuthenticationSelector() {
	  
		const handleSMSChange = (event) => {
		  setMfaChoice(event.target.id);
		};
	  
		const handleEmailChange = (event) => {
		  setMfaChoice(event.target.id);
		};

		const handleIsPasswordChange = (event) => {
			setMfaChoice(event.target.id);
		}

		const handleHighSecChange = (event) => {
			setMfaChoice(event.target.id);
		}
	  
		return (
		  <div>
			<label>
			  <input type="checkbox" id="sms" onChange={handleSMSChange} />
			  SMS
			</label>
			<label>
			  <input type="checkbox" id="email" onChange={handleEmailChange} />
			  Email
			</label>
			<label>
			  <input type="checkbox" id="password" onChange={handleIsPasswordChange} />
			  Password only
			</label>
			<label>
			  <input type="checkbox" id="high" onChange={handleHighSecChange} />
			  High Security
			</label>
		  </div>
		);
	}

  function Verify () {
		return (
			<>
				
				<div className="container-login100-form-btn m-t-17">
				</div>
				<div className="container-login100-form-btn m-t-17">
					<button className="login100-form-btn" type="button" onClick={() =>SignUp()}>
						verify
					</button>
				</div>
				<div className="w-full text-center p-t-55">
					<button onClick={() => setPage("signup")}>Back to account details</button>
				</div>
			</>
		)
	}

  async function ResendOTP (requestType) {
		RequestOTP(userdata.Username, requestType, "temp").then((response) => {
			if (( "errors" in response )) {
				console.log("error bak from otp request", response.errors[0].message)
				alertError(response.errors[0].message)
			}
		})
	}

	// function ChangeAuthType () {

	// 	console.log("change auth type called", setRenderType)
	// 	//fun is only called when type is sms as button only shown on that page. added switch from email for future
	// 	if (rendertype === "sms") {
	// 		ResendOTP("!email")
	// 		setRenderType("email")
			
	// 	} else {
	// 		ResendOTP("!sms")
	// 		setRenderType("sms")
			
	// 	}
	// }



  function LandingPage () {
    return (
      <>
        {rendertype === "username" && <CheckUsername />}
        {rendertype === "email" && <CheckEmail />}
        {rendertype === "sms" && <CheckMobile />}
        {rendertype === "confirmsms" && <ConfirmSms />}
        {rendertype === "confirmemail" && <ConfirmEmail />}
		{rendertype === "setsecurity" && <AuthenticationSelector />}
        <Verify />
      </>
    )
  }
  

  return (
    <>
   
      {page === "default" && <LandingPage/>}
	  {page === "home" && <Home sessionuser={userdata.Username} page={"home"} viewing={userdata.Username} />}
      {page === "signin" && <SignIn/>}

    
    </> 
  );

  

}